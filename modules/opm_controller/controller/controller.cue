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
Version: "1.0.0-beta.9"

Image: {
	repository: "ghcr.io/open-platform-model/opm-controller"
	tag:        "v\(Version)"
	// The content digest of `tag` on GHCR, read when Version is bumped.
	digest: "sha256:99749b5dfb51fa96d16dc40068f73d56aa0d9be13e3ae71aad56b2bded14190e"
}
