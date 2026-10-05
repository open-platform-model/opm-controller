## MODIFIED Requirements

### Requirement: Registry configuration

The controller MUST accept a `--registry` flag for CUE registry configuration.

The controller MUST read the `OPM_REGISTRY` environment variable as a fallback
when the `--registry` flag value is empty.

The controller MUST set `CUE_REGISTRY` and `OPM_REGISTRY` to the resolved
registry value before any CUE module evaluation.

The controller ships a built-in default `--registry` value routing
`opmodel.dev/*` and `testing.opmodel.dev/*` to `ghcr.io/open-platform-model`
with `registry.cue.works` as a fallback mirror. Operators override by passing
an explicit `--registry=<mapping>`, or disable the built-in default by passing
`--registry=""` (which then falls through to `OPM_REGISTRY`).

Precedence (highest first):

1. `--registry` flag value (including the built-in default when the controller
   does not pass the flag).
2. `OPM_REGISTRY` environment variable (reached only when `--registry` is
   explicitly empty).
3. CUE's built-in default resolution (reached only when both `--registry` and
   `OPM_REGISTRY` are empty).

#### Scenario: Flag value wins
- **WHEN** the controller is started with an explicit `--registry=<mapping>`
- **THEN** `CUE_REGISTRY` and `OPM_REGISTRY` are set to that mapping before CUE evaluation

#### Scenario: Env fallback when flag empty
- **WHEN** the controller is started with `--registry=""` and `OPM_REGISTRY` is set
- **THEN** `CUE_REGISTRY` and `OPM_REGISTRY` are set to the `OPM_REGISTRY` value

#### Scenario: Built-in default applied
- **WHEN** the controller is started without passing the `--registry` flag
- **THEN** the built-in default mapping is used, routing `opmodel.dev/*` and `testing.opmodel.dev/*` to `ghcr.io/open-platform-model` with `registry.cue.works` as a fallback

#### Scenario: CUE default resolution
- **WHEN** both `--registry` is explicitly empty and `OPM_REGISTRY` is unset
- **THEN** CUE's built-in default resolution is used
