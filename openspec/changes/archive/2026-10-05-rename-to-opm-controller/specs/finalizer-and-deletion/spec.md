## MODIFIED Requirements

### Requirement: A wait survives the loss of the ServiceAccount
A deletion cleanup has sent every delete when one reconcile of it ended with a release verdict for every plan: each plan step ended as deleted (the API server accepted its DELETE) or as skipped, and no step failed. That reconcile records the fact by setting the `Ready` reason `DeletionInProgress`, or it removes the finalizer at once when nothing is left. The reason `DeletionBlocked` carries the same fact. No other field records it. Source: owner decisions of 2026-10-09 (the `Ready` reason is the record, with no new CRD field; during the wait the finalizer is released "only when the ServiceAccount is gone").

The controller MUST read the record only from the `Ready` condition the object carried when the reconcile started, and only on an object that has a `deletionTimestamp`. It MUST use the record for one decision: to remove the finalizer when the ServiceAccount the cleanup deletes as is gone. While the ServiceAccount exists and every request succeeds, the record MUST NOT change any outcome.

The ServiceAccount is gone in exactly one case: the controller's own read of it, not an impersonated one, by the name the object uses and in the object's namespace, is answered NotFound for that ServiceAccount. The controller MUST decide this from the typed API status (a NotFound whose details name the core resource `serviceaccounts` and that name) and MUST NOT read message text. A NotFound that names another object or none is not this case. A ServiceAccount that was deleted and created again under the same name exists.

With the record present:

- When the ServiceAccount is gone, the controller MUST remove the finalizer without a read or a delete, and MUST emit one `Warning` event with reason `DeletionUnconfirmed` and action `Delete` that states how many inventory objects it could not confirm as gone.
- Every other answer MUST hold the finalizer, on any read or delete of the cleanup and whoever sent it: Forbidden, Unauthorized, a server error, a throttle, a timeout, a connection error, a cancelled context, a ServiceAccount that could not be looked up. This includes a Forbidden answer for one object while other objects are readable, and a ServiceAccount that exists and has lost its rights. The controller MUST retry with its normal backoff, MUST NOT replace the wait reason with a stall reason (the wait reason is the record), and MUST say in the `Ready` message of that reconcile what could not be checked and why. The message MUST name at most three failed requests, MUST count the rest, MUST cut each cause, and MUST stay far below the length a condition message may have. When the status cannot be written, the reconcile MUST fail so that it is retried.
- Once the object has been deleting for 10 minutes by the controller's clock, the controller MUST report a recheck that could not be done with reason `DeletionBlocked` and `Stalled=True`, with one `Warning` event and the ways out, as it reports a deletion whose objects do not go.

Known limit, by the owner's decision: when the rights of the ServiceAccount are removed during the wait and the ServiceAccount stays (its RoleBinding is deleted), the deletion does not release by itself. It reports `DeletionBlocked` after 10 minutes and names the cause; `spec.prune=false` releases it, and so does deleting the ServiceAccount.

Without the record, a missing or failed identity and a Forbidden step MUST hold the finalizer with `DeletionSAMissing` or `ImpersonationFailed`, however many deletes an earlier, failed reconcile already sent.

The controller MUST NOT set the reasons `DeletionInProgress` and `DeletionBlocked` anywhere else, and MUST NOT set them on an object that is not being deleted. The roles the controller ships MUST NOT grant a user a write verb on `moduleinstances/status` or `modulepackages/status`.

#### Scenario: ServiceAccount deleted during the wait
- **GIVEN** a ModuleInstance being deleted whose `Ready` condition carries the reason `DeletionInProgress` and whose Deployment is still terminating
- **WHEN** its ServiceAccount is deleted and the controller reconciles the instance
- **THEN** the finalizer is removed without a read or a delete
- **AND** one `Warning` event with reason `DeletionUnconfirmed` is emitted
- **AND** `Ready` never carries the reason `DeletionSAMissing`

#### Scenario: A blocked deletion is released the same way
- **GIVEN** a ModuleInstance being deleted whose `Ready` condition carries the reason `DeletionBlocked`
- **WHEN** its ServiceAccount is deleted and the controller reconciles the instance
- **THEN** the finalizer is removed

#### Scenario: RBAC removed during the wait holds
- **GIVEN** a ModuleInstance being deleted with reason `DeletionInProgress`, whose ServiceAccount exists and has lost every right because its RoleBinding was deleted
- **WHEN** the controller reconciles the instance and every read is refused as Forbidden
- **THEN** the finalizer stays, `Ready` keeps the reason `DeletionInProgress` and says that the deleted objects could not be checked, and no event with reason `DeletionUnconfirmed` is emitted
- **AND** once the instance has been deleting for 10 minutes, `Ready` carries the reason `DeletionBlocked` and names `spec.prune=false` as a way out

#### Scenario: One refused object next to readable ones holds
- **GIVEN** a ModuleInstance being deleted with reason `DeletionInProgress`, whose read of one inventory object is refused as Forbidden while the others are read
- **WHEN** the readable objects are gone and the controller reconciles the instance
- **THEN** the finalizer stays

#### Scenario: A NotFound that is not the ServiceAccount's does not release
- **GIVEN** a ModuleInstance being deleted with reason `DeletionInProgress`
- **WHEN** the lookup of its ServiceAccount is answered NotFound for another object, or with a NotFound that names none
- **THEN** the finalizer stays and the reconcile is retried

#### Scenario: A recreated ServiceAccount does not release
- **GIVEN** a ModuleInstance being deleted with reason `DeletionInProgress`, whose ServiceAccount was deleted and created again under the same name
- **WHEN** the controller reconciles the instance
- **THEN** the recheck runs as that ServiceAccount, and the finalizer stays while a deleted object still exists

#### Scenario: A failed lookup of the ServiceAccount during the wait keeps the finalizer
- **GIVEN** a ModuleInstance being deleted with reason `DeletionBlocked`, whose ServiceAccount exists
- **WHEN** the controller's read of the ServiceAccount fails with a server error or a timeout
- **THEN** the finalizer stays, `Ready` keeps the reason `DeletionBlocked` and says that the check failed, and the reconcile is retried
- **AND** no event with reason `DeletionUnconfirmed` is emitted

#### Scenario: An Unauthorized answer does not release
- **GIVEN** a ModuleInstance being deleted with reason `DeletionBlocked`
- **WHEN** a read of the cleanup is answered Unauthorized
- **THEN** the finalizer stays, `Ready` keeps the reason `DeletionBlocked`, and the reconcile is retried

#### Scenario: A wait whose rechecks keep failing is reported as blocked
- **GIVEN** a ModuleInstance being deleted with reason `DeletionInProgress`, whose rechecks fail with a server error
- **WHEN** the controller reconciles it before and after 10 minutes of deletion, by the controller's clock
- **THEN** the first reconcile keeps `DeletionInProgress` and its message says that the deleted objects could not be checked, with the cause
- **AND** the later one sets `DeletionBlocked` with `Stalled=True`, names `spec.prune=false` as a way out and emits one `Warning` event

#### Scenario: The message of a failed recheck is bounded
- **GIVEN** a ModuleInstance being deleted with reason `DeletionInProgress` and 400 inventory objects, whose every read fails
- **WHEN** the controller reconciles it
- **THEN** the `Ready` message names three failed requests, says that 397 more failed, and has at most 2048 characters

#### Scenario: Identity lost before every delete was sent
- **GIVEN** a ModuleInstance being deleted whose first cleanup reconcile deleted one ConfigMap and failed on a second, so that `Ready` carries no deletion wait reason
- **WHEN** its ServiceAccount is deleted and the controller reconciles the instance
- **THEN** the finalizer stays and `Ready` is False with reason `DeletionSAMissing`

#### Scenario: A readable terminating object still holds
- **GIVEN** a ModuleInstance with reason `DeletionInProgress` whose ServiceAccount may read Deployments and may no longer delete them, and a Deployment that is still terminating
- **WHEN** the controller reconciles the instance
- **THEN** the finalizer stays and `Ready` keeps the reason `DeletionInProgress`

#### Scenario: The record is ignored on a live object
- **GIVEN** a ModuleInstance that is not being deleted and whose `Ready` condition was written by another client with the reason `DeletionInProgress`
- **WHEN** the controller reconciles it
- **THEN** the reconcile renders and applies as usual and replaces the condition

#### Scenario: The shipped roles cannot write the record
- **GIVEN** every role under `config/rbac` that the controller ships for users (the editor, admin and viewer roles)
- **WHEN** their rules are read
- **THEN** no rule grants `update` or `patch` on a `status` subresource
