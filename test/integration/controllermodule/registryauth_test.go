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

package controllermodule_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	registryAuthSecret = "opm-registry-auth"
	registryAuthDir    = "/etc/opm/docker"
	registryAuthVolume = "registry-auth"
)

// podSpec returns the rendered Deployment's pod spec.
func podSpec(objs []obj) map[string]any {
	GinkgoHelper()
	d := find(objs, "Deployment", "opm-controller-manager")
	return dict(d.Object, "spec", "template", "spec")
}

// registryAuthItem returns the registry-auth entry of a list of maps with a
// "name" key (volumes), or nil.
func registryAuthItem(items []any) map[string]any {
	for _, i := range items {
		if m, ok := i.(map[string]any); ok && m["name"] == registryAuthVolume {
			return m
		}
	}
	return nil
}

var _ = Describe("The controller module's registryAuth value", func() {
	Context("unset", func() {
		It("renders no DOCKER_CONFIG, no registry-auth mount and no registry-auth volume", func() {
			objs := mustRender(moduleDir)
			c := managerContainer(objs)
			Expect(envValue(c, "DOCKER_CONFIG")).To(BeNil())
			Expect(registryAuthItem(list(c, "volumeMounts"))).To(BeNil())
			Expect(registryAuthItem(list(podSpec(objs), "volumes"))).To(BeNil())
		})
	})

	Context("set", Ordered, func() {
		var objs []obj

		BeforeAll(func() {
			objs = mustRenderWith(moduleDir, `{registryAuth: secretName: "`+registryAuthSecret+`"}`)
		})

		It("renders DOCKER_CONFIG, one read-only mount and one volume with the item", func() {
			c := managerContainer(objs)
			Expect(envValue(c, "DOCKER_CONFIG")).To(Equal(registryAuthDir),
				"the environment carries the directory, never a login")

			mounts := 0
			for _, m := range list(c, "volumeMounts") {
				if mm := m.(map[string]any); mm["name"] == registryAuthVolume {
					mounts++
					Expect(mm["mountPath"]).To(Equal(registryAuthDir))
					Expect(mm["readOnly"]).To(BeTrue())
					Expect(mm).NotTo(HaveKey("subPath"), "a subPath mount never receives an update")
				}
			}
			Expect(mounts).To(Equal(1))

			vol := registryAuthItem(list(podSpec(objs), "volumes"))
			Expect(vol).NotTo(BeNil())
			Expect(normalize(vol["secret"])).To(Equal(normalize(map[string]any{
				"secretName": registryAuthSecret,
				"items":      []any{map[string]any{"key": ".dockerconfigjson", "path": "config.json"}},
			})))
		})

		It("leaves the volume required and readable by the non-root pod user", func() {
			vol := registryAuthItem(list(podSpec(objs), "volumes"))
			secret := dict(vol, "secret")
			if opt, ok := secret["optional"]; ok {
				Expect(opt).To(BeFalse(), "a missing Secret must block the pod")
			}
			if mode, ok := secret["defaultMode"]; ok {
				Expect(mode).To(BeNumerically(">=", 0o444),
					"the file is root-owned and the pod has no fsGroup; a tighter mode makes it unreadable")
			}
		})

		It("carries no login in the arguments or the other environment values", func() {
			c := managerContainer(objs)
			Expect(c["args"]).To(Equal(moduleArgs))
			for _, e := range list(c, "env") {
				em := e.(map[string]any)
				if em["name"] != "DOCKER_CONFIG" && em["name"] != "GOMEMLIMIT" {
					Fail("unexpected environment value " + em["name"].(string))
				}
			}
		})

		It("renders no Secret and gives no rendered role a rule on secrets", func() {
			for _, o := range objs {
				Expect(o.kind()).NotTo(Equal("Secret"), "the module creates no Secret (%s)", o.name())
				if o.kind() != kindRole && o.kind() != kindClusterRole {
					continue
				}
				for _, r := range list(o.Object, "rules") {
					Expect(list(r.(map[string]any), "resources")).NotTo(ContainElement("secrets"),
						"%s %s must not read Secrets", o.kind(), o.name())
				}
			}
		})
	})
})
