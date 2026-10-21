// A fixture of the controller module's deployed controller for test-release-check.sh.
package controller

Version: "1.0.0-beta.2"

Image: {
	repository: "ghcr.io/open-platform-model/opm-controller"
	tag:        "v\(Version)"
	digest:     "sha256:1111111111111111111111111111111111111111111111111111111111111111"
}
