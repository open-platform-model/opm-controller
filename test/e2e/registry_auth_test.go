//go:build e2e

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

package e2e

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/open-platform-model/library/opm/helper/platformmodule"
	"github.com/open-platform-model/library/opm/kernel"

	opmcontroller "github.com/open-platform-model/opm-controller/internal/controller"
	"github.com/open-platform-model/opm-controller/test/fixtures"
	"github.com/open-platform-model/opm-controller/test/utils"
)

// This spec proves registryAuth of the controller module on a real cluster. A
// registry that refuses anonymous pulls (htpasswd) runs in Kind and holds the
// hello fixture module. The controller, with the three parts that the module
// renders for registryAuth applied, pulls the module through the login the
// kubelet mounts. A missing Secret keeps the pod in ContainerCreating.
//
// The logins are made at run time and live only in this cluster. No real token
// and no GHCR login is used; core and the catalog are public on GHCR and the
// fixture module is served only by the test registry.
//
// Prerequisites: the cue CLI (publishes the fixture) and kubectl.
var _ = Describe("Registry login for module pulls", Ordered, func() {
	const (
		// regNamespace holds the registry that refuses anonymous pulls.
		regNamespace = "registry-auth-e2e"
		// secretName is the pull Secret in the controller namespace. It differs
		// from the volume name "registry-auth" so a FailedMount message that names
		// the Secret cannot be mistaken for one that names the volume.
		secretName = "registry-login"
		// publicRegistry maps the domains the test registry does not serve.
		publicRegistry = "opmodel.dev=ghcr.io/open-platform-model,registry.cue.works"

		instanceNS   = "default"
		instanceName = "hello"
	)

	var (
		projectDir string
		login      registryLogin
		forwardCmd *exec.Cmd
		// regHost is the registry as the controller sees it: the ClusterIP of
		// the Service and its port. An IP needs no DNS: a search domain on the
		// host network can answer a service name with a wrong address.
		regHost  string
		stuckPod string
	)

	BeforeAll(func() {
		_, err := exec.LookPath("cue")
		Expect(err).NotTo(HaveOccurred(), "the cue CLI is needed to publish the fixture module")

		projectDir, err = utils.GetProjectDir()
		Expect(err).NotTo(HaveOccurred())
		login = newRegistryLogin()

		By("deploying the registry that needs a login")
		applyStdin(registryManifest(regNamespace, login))
		waitAvailable(regNamespace, "registry", 3*time.Minute)
		ip, err := utils.Run(exec.Command("kubectl", "-n", regNamespace, "get", "svc", "registry",
			"-o", "jsonpath={.spec.clusterIP}"))
		Expect(err).NotTo(HaveOccurred())
		regHost = strings.TrimSpace(ip) + ":5000"

		By("publishing the hello fixture module to it")
		port, cmd := portForward(regNamespace, "svc/registry", 5000)
		forwardCmd = cmd
		publishFixture("hello", port, login)

		By("creating the manager namespace and the pull Secret")
		_, _ = utils.Run(exec.Command("kubectl", "create", "ns", namespace))
		applyStdin(pullSecretManifest(secretName, regHost, login))

		By("installing CRDs and deploying the manager without registryAuth")
		_, err = utils.Run(exec.Command("task", "controller:crds"))
		Expect(err).NotTo(HaveOccurred(), "Failed to install CRDs")
		_, err = utils.Run(exec.Command("task", "controller:deploy", fmt.Sprintf("IMG=%s", managerImage)))
		Expect(err).NotTo(HaveOccurred(), "Failed to deploy the manager")

		mapping := "testing.opmodel.dev=" + regHost + "+insecure," + publicRegistry
		patch := fmt.Sprintf(
			`[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--registry=%s"}]`, mapping)
		_, err = utils.Run(exec.Command("kubectl", "-n", namespace, "patch", "deployment",
			"opm-controller-manager", "--type=json", "-p", patch))
		Expect(err).NotTo(HaveOccurred(), "Failed to point the controller at the test registry")
		waitAvailable(namespace, "opm-controller-manager", 3*time.Minute)

		By("applying the cluster Platform")
		_, err = utils.Run(exec.Command("kubectl", "apply", "-f",
			filepath.Join(projectDir, "config/samples/opmodel.dev_v1alpha1_platform.yaml")))
		Expect(err).NotTo(HaveOccurred(), "Failed to apply the Platform")
		Eventually(func(g Gomega) {
			out, err := utils.Run(exec.Command("kubectl", "get", "platform", "cluster",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}"))
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(out).To(Equal("True"), "Platform not Ready yet")
		}, 5*time.Minute, 5*time.Second).Should(Succeed())
	})

	AfterAll(func() {
		// The ModuleInstance prunes through its applier ServiceAccount, so the
		// instance goes first and the RBAC after it (see podinfo_test.go).
		By("removing the hello ModuleInstance")
		_, _ = utils.Run(exec.Command("kubectl", "-n", instanceNS, "delete", "moduleinstance", instanceName,
			"--ignore-not-found", "--wait=false"))
		if _, err := utils.Run(exec.Command("kubectl", "-n", instanceNS, "wait", "--for=delete",
			"moduleinstance/"+instanceName, "--timeout=2m")); err != nil {
			_, _ = utils.Run(exec.Command("kubectl", "-n", instanceNS, "patch", "moduleinstance", instanceName,
				"--type=merge", "-p", `{"metadata":{"finalizers":null}}`))
		}
		_, _ = utils.Run(exec.Command("kubectl", "delete", "--ignore-not-found", "--wait=false", "-f",
			filepath.Join(projectDir, "test/fixtures/modules/hello/moduleinstance.yaml")))

		By("removing the cluster Platform, the Secret and the manager")
		_, _ = utils.Run(exec.Command("kubectl", "delete", "--ignore-not-found", "--wait=false", "-f",
			filepath.Join(projectDir, "config/samples/opmodel.dev_v1alpha1_platform.yaml")))
		_, _ = utils.Run(exec.Command("kubectl", "-n", namespace, "delete", "secret", secretName, "--ignore-not-found"))
		_, _ = utils.Run(exec.Command("task", "controller:uninstall", "IGNORE_NOT_FOUND=true"))
		_, _ = utils.Run(exec.Command("task", "controller:crds:remove", "IGNORE_NOT_FOUND=true"))

		By("removing the registry")
		if forwardCmd != nil && forwardCmd.Process != nil {
			_ = forwardCmd.Process.Kill()
			_ = forwardCmd.Wait()
		}
		_, _ = utils.Run(exec.Command("kubectl", "delete", "ns", regNamespace, "--ignore-not-found"))
	})

	AfterEach(func() {
		if !CurrentSpecReport().Failed() {
			return
		}
		By("dumping diagnostics for the failed spec")
		for _, args := range [][]string{
			{"-n", instanceNS, "get", "moduleinstance", instanceName, "-o", "yaml"},
			{"-n", namespace, "get", "pods", "-o", "wide"},
			{"-n", namespace, "get", "events", "--sort-by=.lastTimestamp"},
			{"-n", namespace, "logs", "-l", "control-plane=manager", "--tail=100"},
			{"-n", regNamespace, "describe", "pods"},
			{"-n", regNamespace, "logs", "deploy/registry", "--all-containers", "--tail=50"},
		} {
			if out, err := utils.Run(exec.Command("kubectl", args...)); err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "--- kubectl %s ---\n%s\n", strings.Join(args, " "), out)
			}
		}
	})

	It("cannot resolve the module while the controller has no login", func() {
		By("applying the hello ModuleInstance")
		_, err := utils.Run(exec.Command("kubectl", "apply", "-f",
			filepath.Join(projectDir, "test/fixtures/modules/hello/moduleinstance.yaml")))
		Expect(err).NotTo(HaveOccurred(), "Failed to apply the hello ModuleInstance")

		Eventually(func(g Gomega) {
			out, err := utils.Run(exec.Command("kubectl", "-n", instanceNS, "get", "moduleinstance", instanceName,
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}/{.status.conditions[?(@.type=='Ready')].reason}"))
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(out).To(Equal("False/ResolutionFailed"))
		}, 4*time.Minute, 5*time.Second).Should(Succeed())

		By("confirming the registry refused the pull, not a network fault")
		msg, err := utils.Run(exec.Command("kubectl", "-n", instanceNS, "get", "moduleinstance", instanceName,
			"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].message}"))
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.ToLower(msg)).To(ContainSubstring("unauthorized"))
	})

	It("pulls the module once the three rendered parts of registryAuth are applied", func() {
		By("rendering the controller module with registryAuth set")
		patch := renderRegistryAuthParts(projectDir, secretName)

		By("applying the env, the mount and the volume to the controller Deployment")
		_, err := utils.Run(exec.Command("kubectl", "-n", namespace, "patch", "deployment",
			"opm-controller-manager", "--type=strategic", "-p", patch))
		Expect(err).NotTo(HaveOccurred(), "Failed to apply the rendered parts")
		_, err = utils.Run(exec.Command("kubectl", "-n", namespace, "rollout", "status",
			"deployment/opm-controller-manager", "--timeout=180s"))
		Expect(err).NotTo(HaveOccurred(), "the rollout with the mounted Secret did not settle")

		By("waiting for the ModuleInstance to become Ready")
		Eventually(func(g Gomega) {
			out, err := utils.Run(exec.Command("kubectl", "-n", instanceNS, "get", "moduleinstance", instanceName,
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}"))
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(out).To(Equal("True"), "hello ModuleInstance not Ready yet")
		}, 6*time.Minute, 5*time.Second).Should(Succeed())

		By("confirming the module rendered its ConfigMap")
		out, err := utils.Run(exec.Command("kubectl", "-n", instanceNS, "get", "configmap",
			"-o", "jsonpath={.items[*].data.message}"))
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(ContainSubstring("hello from ModuleInstance CR"))
	})

	It("keeps the new pod in ContainerCreating with a FailedMount event when the Secret is missing", func() {
		By("deleting the Secret and restarting the controller")
		_, err := utils.Run(exec.Command("kubectl", "-n", namespace, "delete", "secret", secretName))
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Run(exec.Command("kubectl", "-n", namespace, "rollout", "restart",
			"deployment/opm-controller-manager"))
		Expect(err).NotTo(HaveOccurred())

		By("waiting for a pod that waits in ContainerCreating")
		Eventually(func(g Gomega) {
			stuckPod = containerCreatingPod(g)
			g.Expect(stuckPod).NotTo(BeEmpty(), "no pod in ContainerCreating yet")
		}, 2*time.Minute, 3*time.Second).Should(Succeed())

		By("waiting for a FailedMount event that names the Secret")
		Eventually(func(g Gomega) {
			out, err := utils.Run(exec.Command("kubectl", "-n", namespace, "get", "events",
				"--field-selector", "reason=FailedMount,involvedObject.name="+stuckPod, "-o", "json"))
			g.Expect(err).NotTo(HaveOccurred())
			var list struct {
				Items []struct {
					Message string `json:"message"`
				} `json:"items"`
			}
			g.Expect(json.Unmarshal([]byte(out), &list)).To(Succeed())
			msgs := make([]string, 0, len(list.Items))
			for _, e := range list.Items {
				msgs = append(msgs, e.Message)
			}
			g.Expect(strings.Join(msgs, "\n")).To(ContainSubstring(fmt.Sprintf("secret %q not found", secretName)))
		}, 3*time.Minute, 5*time.Second).Should(Succeed())

		By("confirming the pod stays in ContainerCreating")
		Consistently(func(g Gomega) {
			g.Expect(containerCreatingPod(g)).To(Equal(stuckPod))
		}, 15*time.Second, 3*time.Second).Should(Succeed())
	})
})

// registryLogin is one login made for a run. It is never printed.
type registryLogin struct{ user, pass string }

func newRegistryLogin() registryLogin {
	return registryLogin{user: "e2e-" + randHex(4), pass: randHex(16)}
}

func randHex(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	Expect(err).NotTo(HaveOccurred())
	return hex.EncodeToString(b)
}

// dockerConfig is the config.json content for one registry host.
func (l registryLogin) dockerConfig(host string) []byte {
	auth := base64.StdEncoding.EncodeToString([]byte(l.user + ":" + l.pass))
	b, err := json.Marshal(map[string]any{"auths": map[string]any{host: map[string]string{"auth": auth}}})
	Expect(err).NotTo(HaveOccurred())
	return b
}

// registryManifest is the registry that refuses anonymous pulls: the registry
// image behind htpasswd, in a namespace that enforces Pod Security restricted.
// An init container makes the bcrypt file from the run's login, so no hash and
// no login is in the tree. The images are pinned by digest.
func registryManifest(ns string, l registryLogin) []byte {
	const registryImage = "registry:2@sha256:a3d8aaa63ed8681a604f1dea0aa03f100d5895b6a58ace528858a7b332415373"
	const htpasswdImage = "httpd:2.4-alpine@sha256:070da7661f0be8a05430b2c125c0cde9e53176b4048c813024021969a0e3893c"
	return fmt.Appendf(nil, `apiVersion: v1
kind: Namespace
metadata:
  name: %[1]s
  labels:
    pod-security.kubernetes.io/enforce: restricted
---
apiVersion: v1
kind: Secret
metadata:
  name: registry-credentials
  namespace: %[1]s
type: Opaque
stringData:
  user: %[2]s
  pass: %[3]s
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: registry
  namespace: %[1]s
spec:
  replicas: 1
  selector:
    matchLabels: {app: registry}
  template:
    metadata:
      labels: {app: registry}
    spec:
      automountServiceAccountToken: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        runAsGroup: 1000
        seccompProfile: {type: RuntimeDefault}
      initContainers:
        - name: htpasswd
          image: %[4]s
          command: [sh, -c, 'htpasswd -Bbc /auth/htpasswd "$REG_USER" "$REG_PASS"']
          env:
            - name: REG_USER
              valueFrom: {secretKeyRef: {name: registry-credentials, key: user}}
            - name: REG_PASS
              valueFrom: {secretKeyRef: {name: registry-credentials, key: pass}}
          resources:
            requests: {cpu: 10m, memory: 16Mi}
            limits: {memory: 64Mi}
          securityContext: &sc
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities: {drop: [ALL]}
          volumeMounts:
            - {name: auth, mountPath: /auth}
            - {name: tmp, mountPath: /tmp}
      containers:
        - name: registry
          image: %[5]s
          ports:
            - {containerPort: 5000}
          env:
            - {name: REGISTRY_AUTH, value: htpasswd}
            - {name: REGISTRY_AUTH_HTPASSWD_REALM, value: opm-e2e}
            - {name: REGISTRY_AUTH_HTPASSWD_PATH, value: /auth/htpasswd}
            - {name: REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY, value: /data}
          readinessProbe:
            tcpSocket: {port: 5000}
          resources:
            requests: {cpu: 50m, memory: 64Mi}
            limits: {memory: 256Mi}
          securityContext: *sc
          volumeMounts:
            - {name: auth, mountPath: /auth, readOnly: true}
            - {name: data, mountPath: /data}
      volumes:
        - {name: auth, emptyDir: {}}
        - {name: data, emptyDir: {}}
        - {name: tmp, emptyDir: {}}
---
apiVersion: v1
kind: Service
metadata:
  name: registry
  namespace: %[1]s
spec:
  selector: {app: registry}
  ports:
    - {port: 5000, targetPort: 5000}
`, ns, l.user, l.pass, htpasswdImage, registryImage)
}

// pullSecretManifest is the kubernetes.io/dockerconfigjson Secret in the
// controller namespace that registryAuth names.
func pullSecretManifest(name, host string, l registryLogin) []byte {
	b, err := json.Marshal(map[string]any{
		"apiVersion": "v1", "kind": "Secret",
		"metadata":   map[string]any{"name": name, "namespace": namespace},
		"type":       "kubernetes.io/dockerconfigjson",
		"stringData": map[string]string{".dockerconfigjson": string(l.dockerConfig(host))},
	})
	Expect(err).NotTo(HaveOccurred())
	return b
}

// applyStdin applies a manifest from memory, so no login is on a command line.
func applyStdin(manifest []byte) {
	GinkgoHelper()
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(manifest)
	_, err := utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred(), "kubectl apply failed")
}

func waitAvailable(ns, deployment string, timeout time.Duration) {
	GinkgoHelper()
	Eventually(func(g Gomega) {
		out, err := utils.Run(exec.Command("kubectl", "-n", ns, "get", "deploy", deployment,
			"-o", "jsonpath={.status.availableReplicas}"))
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(out).To(Equal("1"), deployment+" not Available yet")
	}, timeout, 3*time.Second).Should(Succeed())
}

// portForward starts kubectl port-forward on a free local port and waits for it.
func portForward(ns, target string, remote int) (int, *exec.Cmd) {
	GinkgoHelper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	Expect(err).NotTo(HaveOccurred())
	port := l.Addr().(*net.TCPAddr).Port
	Expect(l.Close()).To(Succeed())

	cmd := exec.Command("kubectl", "-n", ns, "port-forward", "--address", "127.0.0.1",
		target, fmt.Sprintf("%d:%d", port, remote))
	Expect(cmd.Start()).To(Succeed())
	Eventually(func() error {
		c, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), time.Second)
		if err == nil {
			_ = c.Close()
		}
		return err
	}, 30*time.Second, 500*time.Millisecond).Should(Succeed(), "port-forward did not open")
	return port, cmd
}

// publishFixture publishes a fixture module at its declared version to the
// registry behind the forwarded port, with a throwaway DOCKER_CONFIG and CUE
// cache. Only the test registry sees the login.
func publishFixture(name string, port int, l registryLogin) {
	GinkgoHelper()
	coord, err := fixtures.Load(name)
	Expect(err).NotTo(HaveOccurred())
	dir, err := fixtures.Dir()
	Expect(err).NotTo(HaveOccurred())

	// A CUE module cache is read-only, so t.TempDir cannot remove it: make it
	// writable first.
	work, err := os.MkdirTemp("", "registry-auth-e2e-")
	Expect(err).NotTo(HaveOccurred())
	DeferCleanup(func() {
		_ = filepath.WalkDir(work, func(p string, d os.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				_ = os.Chmod(p, 0o700)
			}
			return nil
		})
		_ = os.RemoveAll(work)
	})
	src := filepath.Join(work, name)
	Expect(os.CopyFS(src, os.DirFS(filepath.Join(dir, name)))).To(Succeed())
	cfgDir := filepath.Join(work, "docker")
	Expect(os.MkdirAll(cfgDir, 0o700)).To(Succeed())
	host := fmt.Sprintf("localhost:%d", port)
	Expect(os.WriteFile(filepath.Join(cfgDir, "config.json"), l.dockerConfig(host), 0o600)).To(Succeed())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cue", "mod", "publish", coord.Tag())
	cmd.Dir = src
	cmd.Env = append(os.Environ(),
		"CUE_REGISTRY=testing.opmodel.dev="+host+"+insecure,opmodel.dev=ghcr.io/open-platform-model,registry.cue.works",
		"DOCKER_CONFIG="+cfgDir,
		"CUE_CACHE_DIR="+filepath.Join(work, "cue-cache"),
	)
	out, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "cue mod publish failed: %s", out)
}

// renderRegistryAuthParts renders the controller module with registryAuth set
// and returns the strategic merge patch that carries its three parts: the
// DOCKER_CONFIG env, the registry-auth mount and the registry-auth volume of
// the manager Deployment.
func renderRegistryAuthParts(projectDir, secret string) string {
	GinkgoHelper()
	ctx := context.Background()
	moduleDir := filepath.Join(projectDir, "modules", "opm_controller")
	const catalogPath = "opmodel.dev/catalogs/opm@v4"
	const reg = "opmodel.dev=ghcr.io/open-platform-model,registry.cue.works"

	pinSrc, err := os.ReadFile(filepath.Join(moduleDir, "cue.mod", "module.cue"))
	Expect(err).NotTo(HaveOccurred())
	pin, err := cuecontext.New().CompileBytes(pinSrc).LookupPath(
		cue.MakePath(cue.Str("deps"), cue.Str(catalogPath), cue.Str("v"))).String()
	Expect(err).NotTo(HaveOccurred(), "the module cue.mod pins no catalog")

	src, err := platformmodule.NewRegistry(platformmodule.RegistryConfig{Registry: reg, Env: os.Environ()})
	Expect(err).NotTo(HaveOccurred())
	entries := []platformmodule.Entry{{Path: catalogPath, Version: pin, Enable: true}}
	deps, err := platformmodule.Closure(ctx, src, platformmodule.Roots(entries))
	Expect(err).NotTo(HaveOccurred(), "resolving the catalog closure")
	files, err := platformmodule.Generate(platformmodule.Input{
		Name: "cluster", Type: "kubernetes", ModulePath: opmcontroller.PlatformModulePath,
		Entries: entries, Deps: deps,
	})
	Expect(err).NotTo(HaveOccurred())
	platDir := GinkgoT().TempDir()
	Expect(files.WriteTo(platDir)).To(Succeed())

	k := kernel.New(kernel.WithRegistry(reg))
	plat, err := k.AcquirePlatformFromDir(ctx, platDir)
	Expect(err).NotTo(HaveOccurred(), "building the platform module")
	mod, err := k.AcquireModuleFromDir(ctx, moduleDir)
	Expect(err).NotTo(HaveOccurred(), "loading the controller module")
	vals, err := k.LoadSourceFromBytes("values.cue",
		fmt.Appendf(nil, "registryAuth: secretName: %q\n", secret))
	Expect(err).NotTo(HaveOccurred())
	inst, err := k.SynthesizeInstance(ctx, kernel.InstanceInput{
		Module: mod, Name: "opm-controller", Namespace: namespace, Values: []kernel.Source{vals},
	})
	Expect(err).NotTo(HaveOccurred(), "registryAuth.secretName was refused")
	res, err := k.Render(ctx, kernel.RenderInput{Instance: inst, Platform: plat, RuntimeName: "registry-auth-e2e"})
	Expect(err).NotTo(HaveOccurred(), "rendering the controller module")

	for _, c := range res.Compiled {
		b, err := c.Value.MarshalJSON()
		Expect(err).NotTo(HaveOccurred())
		var o struct {
			Kind     string `json:"kind"`
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Spec struct {
				Template struct {
					Spec struct {
						Containers []struct {
							Name         string           `json:"name"`
							Env          []map[string]any `json:"env"`
							VolumeMounts []map[string]any `json:"volumeMounts"`
						} `json:"containers"`
						Volumes []map[string]any `json:"volumes"`
					} `json:"spec"`
				} `json:"template"`
			} `json:"spec"`
		}
		Expect(json.Unmarshal(b, &o)).To(Succeed())
		if o.Kind != "Deployment" || o.Metadata.Name != "opm-controller-manager" {
			continue
		}
		pod := o.Spec.Template.Spec
		var env, mount, volume []map[string]any
		for _, c := range pod.Containers {
			if c.Name != "manager" {
				continue
			}
			for _, e := range c.Env {
				if e["name"] == "DOCKER_CONFIG" {
					env = append(env, e)
				}
			}
			for _, m := range c.VolumeMounts {
				if m["name"] == "registry-auth" {
					mount = append(mount, m)
				}
			}
		}
		for _, v := range pod.Volumes {
			if v["name"] == "registry-auth" {
				volume = append(volume, v)
			}
		}
		Expect(env).To(HaveLen(1), "the render must carry the DOCKER_CONFIG env")
		Expect(mount).To(HaveLen(1), "the render must carry the registry-auth mount")
		Expect(volume).To(HaveLen(1), "the render must carry the registry-auth volume")
		patch, err := json.Marshal(map[string]any{"spec": map[string]any{"template": map[string]any{"spec": map[string]any{
			"volumes": volume,
			"containers": []map[string]any{{
				"name": "manager", "env": env, "volumeMounts": mount,
			}},
		}}}})
		Expect(err).NotTo(HaveOccurred())
		return string(patch)
	}
	Fail("the render has no opm-controller-manager Deployment")
	return ""
}

// containerCreatingPod returns the name of a controller pod that waits in
// ContainerCreating, or "" when there is none.
func containerCreatingPod(g Gomega) string {
	out, err := utils.Run(exec.Command("kubectl", "-n", namespace, "get", "pods",
		"-l", "control-plane=manager", "-o", "json"))
	g.Expect(err).NotTo(HaveOccurred())
	var list struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Status struct {
				ContainerStatuses []struct {
					State struct {
						Waiting *struct {
							Reason string `json:"reason"`
						} `json:"waiting"`
					} `json:"state"`
				} `json:"containerStatuses"`
				Phase string `json:"phase"`
			} `json:"status"`
		} `json:"items"`
	}
	g.Expect(json.Unmarshal([]byte(out), &list)).To(Succeed())
	for _, p := range list.Items {
		if p.Status.Phase != "Pending" {
			continue
		}
		for _, cs := range p.Status.ContainerStatuses {
			if cs.State.Waiting != nil && cs.State.Waiting.Reason == "ContainerCreating" {
				return p.Metadata.Name
			}
		}
	}
	return ""
}
