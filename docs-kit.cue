// The opm-controller docs bundle, built by docs-kit's opm-docs (docs-kit
// docs/contracts.md C6, C15, C18) and published by docs.yml and release.yml.
bundles: "opm-controller": {
	placement: {kind: "docs", root: "/docs/", owns: ["reference/controller/"]}
	version: {from: "tag", prefix: "v"}
	sources: [
		// The authored pages ship in the same bundle (docs-kit DESIGN decision
		// 20). docs/site/reference/controller/_index.md holds only front matter
		// and an intro, which complete the crd source's generated section
		// index (C15); the kind pages beside it are generated only.
		{kind: "markdown", dir: "docs/site"},
		{
			kind:        "crd"
			dir:         "./config/crd/bases"
			samples:     "./config/samples"
			section:     "reference/controller/" // the index and one page per kind
			title:       "Controller Reference"
			description: "One generated page per resource kind the controller defines: ModuleInstance, ModulePackage, Platform and TransformerRegistration."
			weight:      7 // the committed index's weight
			order: ["ModuleInstance", "ModulePackage", "Platform", "TransformerRegistration"]
			// The hello fixture is a test module, not an example to copy.
			hideSamplesMatching: ["testing.opmodel.dev"]
			// kubebuilder's scaffold labels on every sample; removed only when the
			// value matches.
			stripLabels: {
				"app.kubernetes.io/name":       "opm-controller"
				"app.kubernetes.io/managed-by": "kustomize"
			}
			// The Named(...) of the reconciler whose builder calls For(&<Kind>{})
			// in internal/controller. Rename a reconciler, edit this map:
			// internal/controller/docskit_reconciledby_test.go fails until they agree.
			reconciledBy: {
				ModuleInstance:          "moduleinstance"
				ModulePackage:           "modulepackage"
				Platform:                "platform"
				TransformerRegistration: "transformerregistration"
			}
		},
	]
}
