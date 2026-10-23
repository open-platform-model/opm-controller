// Package controller names the one controller release this module version
// deploys. It is a package of its own so a tool reads the deployed controller
// version from the module's source without rendering it:
//
//	cue eval ./controller -e Version
//
// A module release that follows a controller release moves all three values
// together. The version never falls below hack/controller-module/min-controller-version
// (the first controller that refuses to reconcile its own instance); the render
// test enforces that.
package controller

// The controller release this module version deploys, bare SemVer.
Version: "1.0.0-beta.10"

Image: {
	repository: "ghcr.io/open-platform-model/opm-controller"
	tag:        "v\(Version)"
	// The content digest of `tag` on GHCR, read when Version is bumped.
	digest: "sha256:fb00c6414044055ed1ba09d41967c8366a950fbea3684d82b83258c1dc87adaa"
}
