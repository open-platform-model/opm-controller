## MODIFIED Requirements

### Requirement: Registry fetch failures are transient wherever they occur
As in requirement "A render failure is classified by its type, never by its message text": A failure SHALL be transient when its error chain holds the library's typed registry fetch failure (`*oerrors.FetchError`, of any `Kind`: unreachable, not found, unauthorized or other) and holds none of the typed terminal causes (an identity mismatch, a wrong artifact kind, a structurally invalid package, a missing required field, unresolved platform demands, unmatched components).

#### Scenario: Registry outage during ModuleInstance acquisition retries on the backoff

- **GIVEN** a ModuleInstance with `failureCounters.reconcile=0` and a generated platform
- **WHEN** fetching its module from the registry fails because the registry cannot be reached
- **THEN** the ModuleInstance reports `Ready=False` with reason `ResolutionFailed`, no `Stalled` condition, a Warning event carrying the acquisition error, and the controller returns `RequeueAfter: 5s`

#### Scenario: Registry outage while loading a ModulePackage retries on the backoff

- **GIVEN** a ModulePackage whose artifact fetched and whose path resolved
- **WHEN** loading the package fails because a CUE dependency in its `cue.mod/module.cue` cannot be fetched from the registry
- **THEN** the ModulePackage reports `Ready=False` with reason `ResolutionFailed`, no `Stalled` condition, and requeues on the exponential backoff rather than the 30-minute recheck

#### Scenario: Registry failure during synthesis or render retries on the backoff

- **GIVEN** a ModuleInstance whose module was acquired
- **WHEN** instance synthesis or the render build fails with a typed registry fetch failure
- **THEN** the ModuleInstance reports `Ready=False` with reason `ResolutionFailed`, no `Stalled` condition, and requeues on the exponential backoff

#### Scenario: Not-found and refused credentials retry on the backoff

- **WHEN** a fetch fails because the registry does not hold the module or version, or refuses the credentials
- **THEN** the object reports `Ready=False` with reason `ResolutionFailed`, no `Stalled` condition, and requeues on the exponential backoff capped at 5 minutes

#### Scenario: Repeated fetch failures cap at five minutes

- **GIVEN** an object with `failureCounters.reconcile=10` whose fetch keeps failing
- **WHEN** the next reconcile fails the same way
- **THEN** the controller requeues after 5 minutes, not 30

#### Scenario: Recovery once the registry answers

- **GIVEN** an object reporting a transient `ResolutionFailed` after a registry fetch failure
- **WHEN** the next retry acquires the module and renders it
- **THEN** the object reports `Ready=True`, `status.nextRetryAt` is nil and `failureCounters.reconcile` resets

#### Scenario: Typed terminal causes stay stalled

- **WHEN** a render fails with an identity mismatch, a wrong artifact kind, a structurally invalid package or a missing required field, even when a registry fetch failure is joined to it
- **THEN** the object reports `Ready=False` and `Stalled=True` and requeues on the 30-minute recheck; the reason is `ResolutionFailed`, except a ModulePackage of the wrong kind, which keeps `UnsupportedKind`

#### Scenario: Author defect in a package stalls

- **GIVEN** a ModulePackage whose package has a CUE syntax error, values that conflict with `#config`, or non-concrete values
- **WHEN** loading the package fails
- **THEN** the ModulePackage reports `Ready=False` with reason `ResolutionFailed` and `Stalled=True`, and requeues on the 30-minute recheck, not the backoff

#### Scenario: Unclassified acquisition failure stalls

- **WHEN** a ModuleInstance's module acquisition fails with an error that is not a typed registry fetch failure, for example a version that does not parse
- **THEN** the ModuleInstance reports `Ready=False` with reason `ResolutionFailed` and `Stalled=True`, and requeues on the 30-minute recheck

#### Scenario: Message text does not classify

- **WHEN** a render fails with an untyped error whose message contains `loading package`, `resolving`, `fetching`, `connection refused` or `not found`
- **THEN** the failure is classified by its typed cause alone: with none, the object reports `Stalled=True` on the 30-minute recheck

### Requirement: A render timeout retries on the backoff
As in requirement "A render timeout is classified before any other render failure": A timeout SHALL set `Ready=False` with reason `RenderTimedOut` and `Reconciling=True`, SHALL NOT set `Stalled=True`, SHALL emit a Warning event with reason `RenderTimedOut` carrying the message, and SHALL produce the outcome `FailedTransient`, so the object requeues on the exponential backoff capped at 5 minutes and `status.nextRetryAt`, `failureCounters.reconcile` and the history follow the transient path.

#### Scenario: A ModuleInstance render times out

- **GIVEN** a ModuleInstance with `failureCounters.reconcile=0` and a renderer that blocks until its context is done
- **WHEN** the render deadline passes
- **THEN** the ModuleInstance reports `Ready=False` with reason `RenderTimedOut`, no `Stalled` condition, a Warning event with reason `RenderTimedOut`, `failureCounters.reconcile=1`, a failure history entry, `status.nextRetryAt` set, and the controller returns `RequeueAfter: 5s`

#### Scenario: A ModulePackage render times out

- **GIVEN** a ModulePackage whose artifact fetched and whose renderer blocks until its context is done
- **WHEN** the render deadline passes
- **THEN** the ModulePackage reports `Ready=False` with reason `RenderTimedOut`, no `Stalled` condition, and requeues on the exponential backoff rather than the 30-minute recheck

#### Scenario: Repeated timeouts cap at five minutes

- **GIVEN** an object with `failureCounters.reconcile=10` whose render keeps timing out
- **WHEN** the next render times out the same way
- **THEN** the controller requeues after 5 minutes, not 30

#### Scenario: A render that fails at its deadline is still a timeout

- **WHEN** a render returns a library stage error that wraps the expired render context, or returns a success, after its deadline has passed
- **THEN** the object reports `RenderTimedOut`, not `ResolutionFailed` or `RenderFailed`, and the result is not applied

#### Scenario: Recovery after a timeout

- **GIVEN** an object reporting `RenderTimedOut`
- **WHEN** a later render finishes within the timeout
- **THEN** the object reports `Ready=True`, `status.nextRetryAt` is nil and `failureCounters.reconcile` resets


## ADDED Requirements

### Requirement: A render failure is classified by its type, never by its message text
When a render fails, the ModuleInstance and ModulePackage reconcilers SHALL classify the failure by its type, never by its message text, with one shared classification so the two cannot drift. This holds for every kernel call the render makes: module acquisition, values compile, instance synthesis, the render build and the package load.

#### Scenario: A render failure is classified by its type
- **WHEN** a render fails
- **THEN** the ModuleInstance and ModulePackage reconcilers classify the failure by its type, never by its message text, with one shared classification so the two cannot drift
- **AND** this holds for every kernel call the render makes: module acquisition, values compile, instance synthesis, the render build and the package load

### Requirement: A transient failure sets ResolutionFailed and produces FailedTransient
As in requirement "Registry fetch failures are transient wherever they occur": A transient failure SHALL set `Ready=False` with reason `ResolutionFailed`, SHALL NOT set `Stalled=True`, SHALL emit a Warning event carrying the error message, and SHALL produce the outcome `FailedTransient`, so the object requeues on the exponential backoff capped at 5 minutes and `status.nextRetryAt` and `failureCounters.reconcile` follow the transient path.

#### Scenario: A transient failure retries on the backoff
- **WHEN** a failure is transient
- **THEN** `Ready=False` is set with reason `ResolutionFailed`, `Stalled=True` is not set, a Warning event carries the error message, and the outcome is `FailedTransient`, so the object requeues on the exponential backoff capped at 5 minutes

### Requirement: Every other failure stalls and requeues on the 30-minute recheck
As in requirement "Registry fetch failures are transient wherever they occur": Every other failure SHALL set `Ready=False` and `Stalled=True` and produce `FailedStalled`, requeued on the 30-minute recheck. Its reason follows the typed cause; a failure during module acquisition or package load that carries no other typed cause is `ResolutionFailed`. The status message SHALL be the error's message unchanged.

#### Scenario: A failure that is not transient stalls
- **WHEN** a failure is not transient
- **THEN** `Ready=False` and `Stalled=True` are set, the outcome is `FailedStalled` and the object is requeued on the 30-minute recheck

### Requirement: A render timeout is classified before any other render failure
When a ModuleInstance or ModulePackage render does not finish within `--render-timeout`, or a reconcile finds that an earlier timed-out render of the same object is still running, the reconciler SHALL classify the attempt as a timeout before any other render classification, whatever error the render returned.

#### Scenario: A render that does not finish is a timeout
- **WHEN** a ModuleInstance or ModulePackage render does not finish within `--render-timeout`
- **THEN** the reconciler classifies the attempt as a timeout before any other render classification, whatever error the render returned

### Requirement: The timeout message names the timeout and the render slot
As in requirement "A render timeout retries on the backoff": The message SHALL name the timeout and say that the render's slot stays held until the render returns, or, for a render still running, that no new render was started.

#### Scenario: The timeout message names the timeout
- **WHEN** a render times out
- **THEN** the message names the timeout and says that the render's slot stays held until the render returns, or, for a render still running, that no new render was started

### Requirement: A timeout keeps inventory and digests and applies and prunes nothing
As in requirement "A render timeout retries on the backoff": Inventory and the last-applied digests SHALL keep the last success, and nothing SHALL be applied or pruned.

#### Scenario: A timeout applies and prunes nothing
- **WHEN** a render times out
- **THEN** inventory and the last-applied digests keep the last success, and nothing is applied or pruned

### Requirement: A deadline error that is not the render deadline keeps the registry fetch classification
A render error that carries `context.DeadlineExceeded` from anything other than the render deadline, for example a registry client's own timeout in a render that returned in time, SHALL keep the classification of "Registry fetch failures are transient wherever they occur".

#### Scenario: A deadline error of another source keeps its classification
- **WHEN** a render error carries `context.DeadlineExceeded` from anything other than the render deadline, for example a registry client's own timeout in a render that returned in time
- **THEN** it keeps the classification of "Registry fetch failures are transient wherever they occur"
