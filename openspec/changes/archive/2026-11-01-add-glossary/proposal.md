## Why

Agents and the owner name the same concept in several ways in this repo: instance and release, claim for a TransformerRegistration and for a PersistentVolumeClaim, platform module and platform package. A root `GLOSSARY.md` fixes one word per concept, so design and refactoring work reads the same language as the code and the docs.

## What Changes

- Add `GLOSSARY.md` at the repo root with the domain terms of this controller, each with a definition and the loose words to avoid.
- Route to it from `AGENTS.md` (layout) and `CONTEXT.md` (one router row).

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. The change has no spec delta (`skip_specs: true`).

## Impact

- Files: `GLOSSARY.md`, `AGENTS.md`, `CONTEXT.md`. `AGENTS.md` and `CONTEXT.md` are agent instruction files, so the diff is reviewed as code and scanned for hidden Unicode.
- No code, no API, no CRD, no generated file. SemVer class after GA: none; commit type `docs` releases nothing.
