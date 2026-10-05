module: "testing.opmodel.dev/releases/controller/hello@v0"
language: {
	version: "v0.17.0"
}
source: {
	kind: "self"
}
deps: {
	"opmodel.dev/catalogs/opm@v4": {
		v: "v4.6.0"
	}
	"opmodel.dev/core@v2": {
		v: "v2.0.0-beta.1"
	}
	"testing.opmodel.dev/modules/controller/hello@v0": {
		v: "v0.0.13"
	}
}
