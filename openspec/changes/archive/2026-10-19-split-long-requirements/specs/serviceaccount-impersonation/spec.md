## REMOVED Requirements

### Requirement: Impersonated identity has the ServiceAccount's groups
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "An impersonated request reaches the apiserver as the ServiceAccount user with its groups", "The controller names no group on an impersonated request", "An RBAC binding on a ServiceAccount group authorises the impersonated request".

### Requirement: Deletion-cleanup SA-missing stalls with distinct reason
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing", "The DeletionSAMissing message names the ServiceAccount and the recovery options", "The DeletionSAMissing stall keeps the finalizer and warns once", "The DeletionSAMissing stall applies to both deletion paths", "Two cases of a missing deletion-cleanup ServiceAccount are not the DeletionSAMissing stall", "Both cases remove the finalizer and warn DeletionUnconfirmed".

### Requirement: Orphan-exit annotation removes finalizer on SA-missing
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The orphan annotation skips prune and removes the finalizer", "The orphan annotation emits OrphanedOnDeletion and clears the inventory", "The orphan annotation takes effect only for a missing ServiceAccount and only with the value true".

### Requirement: A created ServiceAccount wakes the instances that impersonate it
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A created ServiceAccount enqueues the instances stalled on it", "A created ServiceAccount enqueues an instance being deleted and one stalled on the apply path", "A created ServiceAccount enqueues no instance that is not stalled on it", "A created ServiceAccount enqueues no instance it cannot change", "The restore of ServiceAccount RBAC alone is not a trigger".


## ADDED Requirements

### Requirement: An impersonated request reaches the apiserver as the ServiceAccount user with its groups
An impersonated request MUST reach the apiserver as the user `system:serviceaccount:<namespace>:<name>` with exactly the groups a token of that ServiceAccount carries: `system:serviceaccounts`, `system:serviceaccounts:<namespace>` and `system:authenticated`.

#### Scenario: Identity seen by the apiserver
- **GIVEN** a ServiceAccount `deploy-sa` in namespace `team-a`
- **WHEN** the controller's impersonated client for it asks the apiserver who it is
- **THEN** the user is `system:serviceaccount:team-a:deploy-sa`
- **AND** the groups are exactly `system:serviceaccounts`, `system:serviceaccounts:team-a` and `system:authenticated`

### Requirement: The controller names no group on an impersonated request
An impersonated request MUST reach the apiserver as the user `system:serviceaccount:<namespace>:<name>` with exactly the groups a token of that ServiceAccount carries: `system:serviceaccounts`, `system:serviceaccounts:<namespace>` and `system:authenticated`. The controller MUST NOT name a group on an impersonated request; the apiserver derives these groups from the ServiceAccount user name.

#### Scenario: No group is sent
- **GIVEN** a ModuleRelease in namespace `team-a` with `spec.serviceAccountName=deploy-sa`
- **WHEN** the controller builds the impersonated client
- **THEN** the impersonation config names the user `system:serviceaccount:team-a:deploy-sa`
- **AND** it names no group

### Requirement: An RBAC binding on a ServiceAccount group authorises the impersonated request
An impersonated request MUST reach the apiserver as the user `system:serviceaccount:<namespace>:<name>` with exactly the groups a token of that ServiceAccount carries: `system:serviceaccounts`, `system:serviceaccounts:<namespace>` and `system:authenticated`. An RBAC binding whose subject is one of these groups MUST authorise the impersonated request as it authorises the same ServiceAccount with its own token.

#### Scenario: Group-subject RoleBinding authorizes apply
- **GIVEN** a ModuleRelease in namespace `team-a` with `spec.serviceAccountName=deploy-sa` and a RoleBinding in `team-a` whose subjects are `[{Kind: "Group", Name: "system:serviceaccounts:team-a"}]` granting permissions on the resources to be applied
- **WHEN** the controller runs Phase 5 (Apply)
- **THEN** the apply succeeds (the impersonated identity is recognized as a member of `system:serviceaccounts:team-a`)
- **AND** `Ready=True` is set on the ModuleRelease

#### Scenario: ServiceAccount without a binding is refused
- **GIVEN** a ServiceAccount in namespace `team-a` that no RoleBinding names, directly or through a group
- **WHEN** the controller applies a resource in `team-a` as that ServiceAccount
- **THEN** the apiserver refuses the apply as forbidden
- **AND** the resource is not created

### Requirement: A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing
When the deletion cleanup path attempts to build an impersonated client and the target ServiceAccount does not exist (apiserver returns NotFound for the SA lookup), the controller MUST:
- NOT silently fall back to the controller's own client.
- Stall the release with `Ready=False` and condition reason `DeletionSAMissing` (distinct from the generic `ImpersonationFailed` used on the apply path).
- Requeue on the stalled-recheck interval, not the tight transient backoff.

#### Scenario: SA deleted before finalizer can prune
- **GIVEN** a ModuleRelease with `spec.serviceAccountName=hello-applier`, `spec.prune=true`, and a non-empty inventory
- **AND** the ServiceAccount `hello-applier` has been deleted from the release's namespace
- **WHEN** the ModuleRelease is deleted and the finalizer runs deletion cleanup
- **THEN** the impersonated client build fails with SA-NotFound
- **AND** the release's Ready condition becomes False with reason `DeletionSAMissing`
- **AND** the condition message names `default/hello-applier` and lists recovery options
- **AND** a Warning event is emitted with reason `DeletionSAMissing`
- **AND** the finalizer is NOT removed
- **AND** the controller client is NOT used as a fallback for prune

#### Scenario: Other impersonation errors keep existing behavior
- **GIVEN** a ModuleRelease with `spec.serviceAccountName=deploy-sa` and the controller lacks `impersonate` RBAC
- **WHEN** the ModuleRelease is deleted and the finalizer runs deletion cleanup
- **THEN** the release stalls with reason `ImpersonationFailed` (the existing generic reason)
- **AND** the reason is NOT `DeletionSAMissing`
- **AND** the orphan annotation has no effect on this case

### Requirement: The DeletionSAMissing message names the ServiceAccount and the recovery options
When the deletion cleanup path attempts to build an impersonated client and the target ServiceAccount does not exist (apiserver returns NotFound for the SA lookup), the controller MUST:
- Populate the condition message with the SA's namespaced name and list the operator recovery options (restore SA, set `spec.prune=false`, or set the orphan annotation).

#### Scenario: The stall message names the ServiceAccount and lists the recovery options
- **WHEN** the deletion cleanup path attempts to build an impersonated client and the target ServiceAccount does not exist (apiserver returns NotFound for the SA lookup)
- **THEN** the condition message holds the SA's namespaced name and lists the operator recovery options (restore SA, set `spec.prune=false`, or set the orphan annotation)

### Requirement: The DeletionSAMissing stall keeps the finalizer and warns once
When the deletion cleanup path attempts to build an impersonated client and the target ServiceAccount does not exist (apiserver returns NotFound for the SA lookup), the controller MUST:
- Retain the finalizer.
- Emit a `Warning` event with reason `DeletionSAMissing` on the Ready transition (not on every requeue).

#### Scenario: The stalled release keeps its finalizer and warns once
- **WHEN** the deletion cleanup path attempts to build an impersonated client and the target ServiceAccount does not exist (apiserver returns NotFound for the SA lookup)
- **THEN** the finalizer is retained
- **AND** a `Warning` event with reason `DeletionSAMissing` is emitted on the Ready transition (not on every requeue)

### Requirement: The DeletionSAMissing stall applies to both deletion paths
As in requirements "A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing", "The DeletionSAMissing message names the ServiceAccount and the recovery options" and "The DeletionSAMissing stall keeps the finalizer and warns once": This behavior MUST apply symmetrically to both `ModuleRelease` and `Release` deletion paths.

#### Scenario: Both deletion paths stall alike
- **WHEN** the deletion cleanup path of a `ModuleRelease` or of a `Release` attempts to build an impersonated client and the target ServiceAccount does not exist
- **THEN** this behavior applies symmetrically to both `ModuleRelease` and `Release` deletion paths

### Requirement: Two cases of a missing deletion-cleanup ServiceAccount are not the DeletionSAMissing stall
As in requirement "A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing": Two cases SHALL NOT be this stall, because the missing ServiceAccount has nothing left to do (capability `finalizer-and-deletion`): an inventory that holds nothing the cleanup would delete (it is empty, or holds only PersistentVolumeClaims that `spec.dataPolicy` keeps), and a cleanup that already sent every delete and only waits for the objects to be gone.

#### Scenario: ServiceAccount removed while the cleanup waits
- **GIVEN** a ModuleInstance being deleted whose `Ready` condition carries the reason `DeletionInProgress`
- **WHEN** its ServiceAccount is deleted and the controller reconciles the instance
- **THEN** the finalizer is removed and the reason `DeletionSAMissing` is never set

### Requirement: Both cases remove the finalizer and warn DeletionUnconfirmed
As in requirement "Two cases of a missing deletion-cleanup ServiceAccount are not the DeletionSAMissing stall": In both the finalizer SHALL be removed, the controller's own client SHALL still never be used, and a `Warning` event with reason `DeletionUnconfirmed` SHALL say what could not be read.

#### Scenario: Only kept claims and no ServiceAccount
- **GIVEN** a ModuleInstance with `spec.prune=true` and no `spec.dataPolicy` being deleted, whose inventory holds only PersistentVolumeClaims, and whose ServiceAccount does not exist
- **WHEN** the controller reconciles it
- **THEN** no object is read or deleted and the finalizer is removed
- **AND** the reason `DeletionSAMissing` is never set
- **AND** a `Warning` event with reason `DeletionUnconfirmed` says that the claims were left in place without a read

### Requirement: The orphan annotation skips prune and removes the finalizer
When a release is in the `DeletionSAMissing` stall state AND the annotation `opm.dev/force-delete-orphan` is set to `"true"` on the release, the controller MUST:
- Skip prune entirely.
- Remove the finalizer.

#### Scenario: Orphan annotation releases a stuck deletion
- **GIVEN** a ModuleRelease stalled with reason `DeletionSAMissing` and `status.inventory` listing 3 entries
- **WHEN** an operator patches the release with annotation `opm.dev/force-delete-orphan=true`
- **AND** the next reconcile runs
- **THEN** prune is not attempted
- **AND** a Warning event is emitted with reason `OrphanedOnDeletion` and message referencing 3 orphaned entries
- **AND** the finalizer is removed
- **AND** the release object is garbage-collected by the apiserver
- **AND** the 3 previously-managed resources remain in the cluster

### Requirement: The orphan annotation emits OrphanedOnDeletion and clears the inventory
When a release is in the `DeletionSAMissing` stall state AND the annotation `opm.dev/force-delete-orphan` is set to `"true"` on the release, the controller MUST:
- Emit a `Warning` event with reason `OrphanedOnDeletion` naming the number of inventory entries left behind.
- Clear `status.inventory` in the final status patch so the last-observed state does not claim ownership of resources the controller is abandoning.

#### Scenario: The orphan exit warns and clears the inventory
- **WHEN** a release is in the `DeletionSAMissing` stall state AND the annotation `opm.dev/force-delete-orphan` is set to `"true"` on the release
- **THEN** a `Warning` event with reason `OrphanedOnDeletion` names the number of inventory entries left behind
- **AND** `status.inventory` is cleared in the final status patch

### Requirement: The orphan annotation takes effect only for a missing ServiceAccount and only with the value true
As in requirement "The orphan annotation skips prune and removes the finalizer": The annotation MUST only take effect when the deletion cleanup impersonation failure is SA-NotFound. For any other impersonation error (transient, RBAC-denied, etc.), the annotation MUST be ignored and the existing stall behavior preserved.
Any annotation value other than the literal string `"true"` MUST be treated as absent.

#### Scenario: Orphan annotation ignored for non-SA-missing stall
- **GIVEN** a ModuleRelease stalled with reason `ImpersonationFailed` (RBAC-denied impersonate verb, not SA-missing)
- **AND** the annotation `opm.dev/force-delete-orphan=true` is set
- **WHEN** the next reconcile runs
- **THEN** the annotation has no effect
- **AND** the release remains stalled with `ImpersonationFailed`
- **AND** the finalizer is NOT removed

#### Scenario: Annotation value other than "true" is treated as absent
- **GIVEN** a ModuleRelease stalled with reason `DeletionSAMissing`
- **AND** the annotation `opm.dev/force-delete-orphan=yes` is set (value is not the literal `"true"`)
- **WHEN** the next reconcile runs
- **THEN** the annotation is ignored
- **AND** the release remains stalled with `DeletionSAMissing`
- **AND** the finalizer is NOT removed

### Requirement: A created ServiceAccount enqueues the instances stalled on it
When a ServiceAccount is created, the controller MUST enqueue every `ModuleInstance` in that ServiceAccount's namespace whose effective ServiceAccount (`spec.serviceAccountName`, else the manager's `--default-service-account`) has that name and that is stalled on it: its `Ready` condition is `False` with reason `DeletionSAMissing` or `ImpersonationFailed`. Such an instance then recovers without waiting for its stalled recheck.

#### Scenario: A created ServiceAccount enqueues the instances that name it

- **GIVEN** instance `a` with `spec.serviceAccountName: deploy-sa` and instance `b` with `spec.serviceAccountName: other-sa`, both in namespace `team-a` and both stalled with reason `ImpersonationFailed`
- **WHEN** the ServiceAccount `deploy-sa` is created in `team-a`
- **THEN** instance `a` is enqueued and instance `b` is not

#### Scenario: The flag default counts as the effective ServiceAccount

- **GIVEN** a manager started with `--default-service-account=opm-deployer` and an instance in `team-a` with no `spec.serviceAccountName`, stalled with reason `ImpersonationFailed`
- **WHEN** the ServiceAccount `opm-deployer` is created in `team-a`
- **THEN** that instance is enqueued

### Requirement: A created ServiceAccount enqueues an instance being deleted and one stalled on the apply path
When a ServiceAccount is created, the controller MUST enqueue every `ModuleInstance` in that ServiceAccount's namespace whose effective ServiceAccount (`spec.serviceAccountName`, else the manager's `--default-service-account`) has that name and that is stalled on it: its `Ready` condition is `False` with reason `DeletionSAMissing` or `ImpersonationFailed`. This holds for an instance being deleted and for one stalled on the apply path.

#### Scenario: An instance being deleted and one stalled on the apply path are both enqueued
- **WHEN** a ServiceAccount is created
- **THEN** the instances in its namespace whose effective ServiceAccount has that name are enqueued, for an instance being deleted and for one stalled on the apply path

### Requirement: A created ServiceAccount enqueues no instance that is not stalled on it
The controller MUST NOT enqueue an instance that is not stalled on its ServiceAccount, whatever ServiceAccount is created. A user who may create ServiceAccounts in a namespace, but no `ModuleInstance`, MUST NOT be able to make a healthy instance reconcile.

#### Scenario: A healthy instance is not enqueued

- **GIVEN** an instance in `team-a` with `spec.serviceAccountName: deploy-sa` whose `Ready` condition is `True`, and one that was never reconciled
- **WHEN** the ServiceAccount `deploy-sa` is created in `team-a`, once or many times
- **THEN** neither instance is enqueued

### Requirement: A created ServiceAccount enqueues no instance it cannot change
As in requirement "A created ServiceAccount enqueues the instances stalled on it": The controller MUST NOT enqueue an instance the ServiceAccount cannot change: a CLI-owned instance, an instance that names another ServiceAccount, an instance in another namespace, and a suspended instance that is not being deleted. An update or a deletion of a ServiceAccount MUST NOT enqueue anything.

#### Scenario: Instances the ServiceAccount cannot change are not enqueued

- **GIVEN** a CLI-owned instance, a suspended instance that is not being deleted, and an instance in namespace `team-b`, all with `spec.serviceAccountName: deploy-sa` and all stalled with reason `ImpersonationFailed`
- **WHEN** the ServiceAccount `deploy-sa` is created in `team-a`
- **THEN** none of them is enqueued

#### Scenario: An update or a deletion of a ServiceAccount enqueues nothing
- **WHEN** a ServiceAccount is updated or deleted
- **THEN** nothing is enqueued

### Requirement: The restore of ServiceAccount RBAC alone is not a trigger
As in requirement "A created ServiceAccount enqueues the instances stalled on it": The restore of a ServiceAccount's RBAC alone (a Role or a binding) SHALL NOT be a trigger; an instance stalled with `ImpersonationFailed` on a forbidden request still waits for its stalled recheck.

#### Scenario: A restored Role or binding enqueues nothing
- **WHEN** a Role or a binding of a ServiceAccount is restored, and not the ServiceAccount itself
- **THEN** it is not a trigger
- **AND** an instance stalled with `ImpersonationFailed` on a forbidden request still waits for its stalled recheck
