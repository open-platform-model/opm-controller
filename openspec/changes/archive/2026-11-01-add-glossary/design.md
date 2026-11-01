## Context

The repo has ADRs, a constitution and site pages, but no list of its domain words. ADR-003 and ADR-008 still use the retired names `ModuleRelease` and `Release`, and the CRD type docs mix "claim" for two different objects.

## Goals / Non-Goals

**Goals:** one definition per repo-specific term, each backed by a file in this repo; a list of avoided words for terms that blur.

**Non-Goals:** rewriting the ADRs or any type comment to match; defining general Kubernetes or Go terms.

## Decisions

### Form

**Decision**: the form of the glossary skills: a context heading, a `## Language` section, terms grouped under subheadings, a bold term, a one or two sentence definition and an `_Avoid_` line.
**Rationale**: the architecture and domain-model skills read and extend a file in this form.

### Routing

**Decision**: `AGENTS.md` lists the file in its layout block; `CONTEXT.md` gets one router row. `GLOSSARY.md` is a root name the workspace docs-form rule allows.
**Rationale**: one home per rule; `CONTEXT.md` stays a router.

### Sources

**Decision**: every term is backed by the API types, `internal/`, `docs/`, `docs/site/`, `adr/` or `CONSTITUTION.md`. A term with no source here is left out.
**Rationale**: a glossary that invents terms misleads agents more than none.
