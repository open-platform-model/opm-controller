/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reconcile_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cuelabs.dev/go/oci/ociregistry/ocimem"
	"cuelabs.dev/go/oci/ociregistry/ociserver"
	"cuelang.org/go/mod/modconfig"
	"cuelang.org/go/mod/modregistry"
	"cuelang.org/go/mod/module"
	modzip "cuelang.org/go/mod/modzip"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	oerrors "github.com/open-platform-model/library/opm/errors"
	"github.com/open-platform-model/library/opm/kernel"

	"github.com/open-platform-model/opm-controller/test/fixtures"
)

// loginRegistry is an in-process OCI registry behind HTTP basic auth. The set
// of accepted users is mutable, so a spec can revoke a login. It records the
// user of every request that passed the check.
type loginRegistry struct {
	srv *httptest.Server

	mu      sync.Mutex
	allowed map[string]string // user -> password
	served  []string          // user of each accepted request, in order
}

func newLoginRegistry() *loginRegistry {
	r := &loginRegistry{allowed: map[string]string{}}
	backend := ociserver.New(ocimem.New(), nil)
	r.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		r.mu.Lock()
		want, known := r.allowed[user]
		accepted := ok && known && want == pass
		if accepted {
			r.served = append(r.served, user)
		}
		r.mu.Unlock()
		if !accepted {
			w.Header().Set("WWW-Authenticate", `Basic realm="opm-test"`)
			http.Error(w, "login required", http.StatusUnauthorized)
			return
		}
		backend.ServeHTTP(w, req)
	}))
	DeferCleanup(r.srv.Close)
	return r
}

func (r *loginRegistry) host() string { return strings.TrimPrefix(r.srv.URL, "http://") }

// mapping is the CUE_REGISTRY value: the fixture domain goes to this registry,
// the first-party domain to GHCR (core and the abstraction catalog, public).
func (r *loginRegistry) mapping() string {
	return "testing.opmodel.dev=" + r.host() + "+insecure," +
		"opmodel.dev=ghcr.io/open-platform-model,registry.cue.works"
}

func (r *loginRegistry) allow(user, pass string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.allowed[user] = pass
}

func (r *loginRegistry) revoke(user string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.allowed, user)
}

func (r *loginRegistry) markServed() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.served)
}

// usersSince lists the distinct users of requests accepted after mark.
func (r *loginRegistry) usersSince(mark int) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	seen := map[string]bool{}
	var users []string
	for _, u := range r.served[mark:] {
		if !seen[u] {
			seen[u] = true
			users = append(users, u)
		}
	}
	return users
}

// writeDockerConfig replaces dir/config.json with a login for host, the way
// the kubelet updates a mounted Secret: write a new file, rename it over the
// old one.
func writeDockerConfig(dir, host, user, pass string) {
	GinkgoHelper()
	body, err := json.Marshal(map[string]any{
		"auths": map[string]any{
			host: map[string]string{"auth": base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))},
		},
	})
	Expect(err).NotTo(HaveOccurred())
	tmp := filepath.Join(dir, ".config.json.new")
	Expect(os.WriteFile(tmp, body, 0o600)).To(Succeed())
	Expect(os.Rename(tmp, filepath.Join(dir, "config.json"))).To(Succeed())
}

// publishFixture publishes the fixture tree at dir to the registry as
// path@version, through a client that holds the publisher login.
func publishFixture(reg *loginRegistry, publisherConfigDir, dir, modPath, version string) {
	GinkgoHelper()
	mv, err := module.NewVersion(modPath, version)
	Expect(err).NotTo(HaveOccurred())
	var zipped bytes.Buffer
	Expect(modzip.CreateFromDir(&zipped, mv, dir)).To(Succeed())

	resolver, err := modconfig.NewResolver(&modconfig.Config{
		Env: append(os.Environ(), "CUE_REGISTRY="+reg.mapping(), "DOCKER_CONFIG="+publisherConfigDir),
	})
	Expect(err).NotTo(HaveOccurred())
	client := modregistry.NewClientWithResolver(resolver)
	Expect(client.PutModule(ctx, mv, bytes.NewReader(zipped.Bytes()), int64(zipped.Len()))).To(Succeed())
}

// ghcrReachable skips the spec when the first-party domain cannot be reached:
// the fixtures import core, which these specs pull from GHCR.
func ghcrReachable() {
	GinkgoHelper()
	c := http.Client{Timeout: 5 * time.Second}
	resp, err := c.Get("https://ghcr.io/v2/")
	if err != nil {
		registrySkip("ghcr.io is not reachable: " + err.Error())
	}
	_ = resp.Body.Close()
}

var _ = Describe("Registry logins from DOCKER_CONFIG", func() {
	const (
		user     = "opm"
		password = "s3cret-one"
	)

	var (
		reg         *loginRegistry
		configDir   string // DOCKER_CONFIG of the pulling process
		publisher   string // DOCKER_CONFIG of the publishing client
		provider    fixtures.Coordinate
		providerDir string
	)

	// newCacheDir points CUE_CACHE_DIR at an empty directory: the library serves
	// a version it already holds from that cache without a registry call, so
	// only a cold cache makes a pull reach the registry.
	newCacheDir := func() {
		GinkgoHelper()
		dir, err := os.MkdirTemp("", "opm-login-cache-")
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			// The module cache is read-only on disk.
			_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
				if err == nil && d.IsDir() {
					_ = os.Chmod(p, 0o700)
				}
				return nil
			})
			_ = os.RemoveAll(dir)
		})
		GinkgoT().Setenv("CUE_CACHE_DIR", dir)
	}

	// expectUnauthorized checks that err is a refusal by the registry's login
	// check, not some other fetch failure.
	expectUnauthorized := func(err error) {
		GinkgoHelper()
		Expect(err).To(HaveOccurred())
		fe, typed := errors.AsType[*oerrors.FetchError](err)
		Expect(typed).To(BeTrue(), "want a *FetchError, got %T: %v", err, err)
		Expect(fe.Kind).To(Equal(oerrors.FetchUnauthorized), "kind %s: %v", fe.Kind, err)
	}

	pullCatalog := func(k *kernel.Kernel, version string) error {
		_, err := k.AcquireCatalogFromRegistry(ctx, provider.ModulePath, version)
		return err
	}

	BeforeEach(func() {
		ghcrReachable()
		newCacheDir()

		configDir = GinkgoT().TempDir()
		publisher = GinkgoT().TempDir()
		// The pulling process reads logins from DOCKER_CONFIG only; an empty
		// directory keeps the developer's own ~/.docker out of the spec.
		GinkgoT().Setenv("DOCKER_CONFIG", configDir)
		GinkgoT().Setenv("DOCKER_AUTH_CONFIG", "")

		reg = newLoginRegistry()
		reg.allow("publisher", "publish-pass")
		writeDockerConfig(publisher, reg.host(), "publisher", "publish-pass")

		var err error
		provider, err = fixtures.LoadCatalog("provider")
		Expect(err).NotTo(HaveOccurred())
		root, err := fixtures.CatalogDir()
		Expect(err).NotTo(HaveOccurred())
		providerDir = filepath.Join(root, "provider")
		publishFixture(reg, publisher, providerDir, provider.ModulePath, provider.Tag())

		reg.allow(user, password)
	})

	It("refuses a catalog pull when no login is present", func() {
		err := pullCatalog(kernel.New(kernel.WithRegistry(reg.mapping())), provider.Tag())
		Expect(err).To(HaveOccurred())

		fe, typed := errors.AsType[*oerrors.FetchError](err)
		Expect(typed).To(BeTrue(), "want a *FetchError, got %T: %v", err, err)
		Expect(fe.Kind).To(Equal(oerrors.FetchUnauthorized), "kind %s: %v", fe.Kind, err)
	})

	It("pulls a catalog with the login in config.json", func() {
		writeDockerConfig(configDir, reg.host(), user, password)
		mark := reg.markServed()

		Expect(pullCatalog(kernel.New(kernel.WithRegistry(reg.mapping())), provider.Tag())).To(Succeed())
		Expect(reg.usersSince(mark)).To(Equal([]string{user}), "the registry saw the mounted login")
	})

	It("pulls a module with the login in config.json", func() {
		hello := fixtures.Must(GinkgoT(), "hello")
		modules, err := fixtures.Dir()
		Expect(err).NotTo(HaveOccurred())
		publishFixture(reg, publisher, filepath.Join(modules, "hello"), hello.ModulePath, hello.Tag())

		k := kernel.New(kernel.WithRegistry(reg.mapping()))
		_, err = k.AcquireModuleFromRegistry(ctx, hello.ModulePath, hello.Tag())
		Expect(err).To(HaveOccurred(), "no login yet")
		fe, typed := errors.AsType[*oerrors.FetchError](err)
		Expect(typed).To(BeTrue(), "want a *FetchError, got %T: %v", err, err)
		Expect(fe.Kind).To(Equal(oerrors.FetchUnauthorized))

		writeDockerConfig(configDir, reg.host(), user, password)
		mod, err := k.AcquireModuleFromRegistry(ctx, hello.ModulePath, hello.Tag())
		Expect(err).NotTo(HaveOccurred(), "the same Kernel, after the login appeared")
		Expect(mod.Metadata.Name).To(Equal("hello"))
	})

	It("reads the file again for the pull after a refused pull", func() {
		writeDockerConfig(configDir, reg.host(), user, "wrong-password")
		k := kernel.New(kernel.WithRegistry(reg.mapping()))

		expectUnauthorized(pullCatalog(k, provider.Tag())) // a wrong login is refused

		// Replace the file the way the kubelet does, then pull with the same
		// Kernel and no restart.
		writeDockerConfig(configDir, reg.host(), user, password)
		Expect(pullCatalog(k, provider.Tag())).To(Succeed())
	})

	It("keeps the old login until a pull fails (rotation needs a restart)", func() {
		const rotated = "second"
		reg.allow(rotated, "s3cret-two")
		publishFixture(reg, publisher, providerDir, provider.ModulePath, "v0.1.2")
		publishFixture(reg, publisher, providerDir, provider.ModulePath, "v0.1.3")
		publishFixture(reg, publisher, providerDir, provider.ModulePath, "v0.1.4")

		writeDockerConfig(configDir, reg.host(), user, password)
		k := kernel.New(kernel.WithRegistry(reg.mapping()))
		Expect(pullCatalog(k, provider.Tag())).To(Succeed())

		// The Secret changes while the old login still works. The next pull of a
		// version the cache lacks must still reach the registry.
		writeDockerConfig(configDir, reg.host(), rotated, "s3cret-two")
		mark := reg.markServed()
		Expect(pullCatalog(k, "v0.1.2")).To(Succeed())
		Expect(reg.usersSince(mark)).To(Equal([]string{user}),
			"the shared client keeps the login it read when it was built")

		// The registry revokes the old login. The next pull fails and drops the
		// client.
		reg.revoke(user)
		expectUnauthorized(pullCatalog(k, "v0.1.3")) // the revoked login is refused

		// The pull after it reads the file again and uses the new login.
		mark = reg.markServed()
		Expect(pullCatalog(k, "v0.1.4")).To(Succeed())
		Expect(reg.usersSince(mark)).To(Equal([]string{rotated}))
	})

	It("serves a version already cached after the login is revoked", func() {
		writeDockerConfig(configDir, reg.host(), user, password)
		k := kernel.New(kernel.WithRegistry(reg.mapping()))
		Expect(pullCatalog(k, provider.Tag())).To(Succeed())

		reg.revoke(user)
		mark := reg.markServed()
		Expect(pullCatalog(k, provider.Tag())).To(Succeed(), "the module cache answers; no registry call")
		Expect(reg.usersSince(mark)).To(BeEmpty())
	})
})
