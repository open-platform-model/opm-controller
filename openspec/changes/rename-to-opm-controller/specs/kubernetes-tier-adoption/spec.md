## REMOVED Requirements

### Requirement: The operator uses the library's object and label packages and keeps no copy

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller uses the library's object and label packages and keeps no copy" in this capability.

### Requirement: The operator uses the library's inventory package and keeps no copy

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller uses the library's inventory package and keeps no copy" in this capability.

## ADDED Requirements

### Requirement: The controller uses the library's object and label packages and keeps no copy

The controller SHALL take the Kubernetes resource wrapper over the kernel's compiled output, its JSON and unstructured conversion, the duplicate-identity check, the OPM label keys and values, and the recognition of an OPM manager label value from the library's Kubernetes tier (`opm/k8s/object` and `opm/k8s/labels`). It SHALL declare no type, constant or function of its own for any of them and no alias to them. Its lint SHALL refuse an import of the package path that held its earlier copy, `github.com/open-platform-model/opm-controller/pkg/core`. Source: 0012:D3:R6.

#### Scenario: No local copy

- **WHEN** the controller's Go source is searched for the package `pkg/core`, a `Resource` type wrapping a CUE value with instance, component and transformer provenance, or a declaration of the `app.kubernetes.io/managed-by` or `module-instance.opmodel.dev/uuid` label key
- **THEN** none is found; every reader imports `opm/k8s/object` or `opm/k8s/labels`

#### Scenario: Lint refuses the old path

- **WHEN** a package exists at `github.com/open-platform-model/opm-controller/pkg/core` and a Go file in the repository imports it
- **THEN** `task dev:lint` fails on the rule `no-local-kubernetes-tier-copy`, with a message naming the library packages to use instead

### Requirement: The controller uses the library's inventory package and keeps no copy

The controller SHALL compute the inventory stale set, the inventory digest and the render digest only with the library's `opm/k8s/inventory` (`StaleSet`, `Digest`, `RenderDigest`). It SHALL build inventory entries with that package's `NewEntry`. It SHALL declare no entry constructor, identity relation, stale-set function, inventory digest or render digest of its own. Its only inventory code is the conversion between the API's `InventoryEntry` and `inventory.Entry`, and a test of `internal/inventory` SHALL refuse any further export from that package, so a copy cannot come back there unnoticed. Source: 0012:D3:R6.

#### Scenario: No local copy

- **WHEN** the controller's Go source is searched for `ComputeStaleSet`, `ComputeDigest`, `NewEntryFromResource`, `IdentityEqual`, `K8sIdentityEqual` or a render digest function in `internal/status`
- **THEN** none is found, and every reader calls `opm/k8s/inventory`

#### Scenario: A copy in the conversion package fails the tests

- **WHEN** a function is exported from `internal/inventory` beyond `Current`, `FromEntries`, `FromEntry`, `ToEntries` and `ToEntry`
- **THEN** the package's closed-surface test fails
