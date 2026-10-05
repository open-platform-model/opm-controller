// provider — the controller's test catalog: the smallest valid core #Catalog,
// one resource and one transformer. The registry-backed specs use it as a
// real, resolvable catalog other than opmodel.dev/catalogs/opm, for example
// the catalog a TransformerRegistration claim contributes to the generated
// platform. Every member FQN derives from this catalog's own path
// (identity/identity.cue), so it can sit in a platform beside the opm catalog
// without a key collision.
package provider

import (
	c "opmodel.dev/core@v2"
	id "testing.opmodel.dev/catalogs/controller/provider/identity"
	t "testing.opmodel.dev/catalogs/controller/provider/transformers"
	res "testing.opmodel.dev/catalogs/controller/provider/resources/v1alpha1"
)

c.#Catalog
metadata: {
	modulePath:  id.ModulePath
	version:     id.Version
	description: "Controller test catalog: one probe resource and its transformer"
}

#resources: {
	(res.#ProbeResource.metadata.fqn): res.#ProbeResource
}

#transformers: {
	(t.#ProbeTransformer.metadata.fqn): t.#ProbeTransformer
}
