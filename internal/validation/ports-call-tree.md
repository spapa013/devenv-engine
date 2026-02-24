# Ports Validation Call Tree (`cmd/devenv` -> `internal/validation/ports.go`)

This document maps runtime call paths from CLI entrypoints in `cmd/devenv` to all functions, constants, and types declared in `internal/validation/ports.go`.

## 1) CLI entrypoints that reach `ports.go`

### Path A: `devenv validate`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/root.go` -> `validateCmd`
3. `cmd/devenv/validate.go`:
   - `validation.NewPortValidator(validateConfigDir)`
   - `validator.ValidateAll()` or `validator.ValidateSingle(developerName)`
4. `internal/validation/ports.go` methods execute:
   - directory discovery
   - per-developer config load
   - range + conflict checks

### Path B: `devenv generate`

- No usage of `internal/validation` from generate flow.

## 2) Symbol-by-symbol map for `internal/validation/ports.go`

### Constants

- `NodePortMin`:
  - Used by:
    - `validateSingleDeveloper(...)` range checks
    - `cmd/devenv/validate.go` help/suggestion output

- `NodePortMax`:
  - Used by:
    - `validateSingleDeveloper(...)` range checks
    - `cmd/devenv/validate.go` help/suggestion output

### Types

- `PortValidator`:
  - Created by `NewPortValidator(...)`.
  - Holds `configDir` scanned for developer configs.

- `ValidationResult`:
  - Returned by `ValidateAll(...)` and `ValidateSingle(...)`.
  - Consumed by `cmd/devenv/validate.go` for user-facing output and exit behavior.

- `ValidationError`:
  - Produced by:
    - `validateSingleDeveloper(...)` (`invalid`, `out_of_range`)
    - `ValidateAll(...)` (`conflict`)
  - Consumed by `cmd/devenv/validate.go`.

- `ValidationWarning`:
  - Produced by:
    - `ValidateAll(...)` (`no_configs`)
    - `validateSingleDeveloper(...)` (`no_ssh_port`)
  - Consumed by `cmd/devenv/validate.go`.

### Functions/methods

- `NewPortValidator(configDir string) *PortValidator`:
  - Called by `cmd/devenv/validate.go`.

- `(pv *PortValidator) ValidateAll() (*ValidationResult, error)`:
  - Called by `cmd/devenv/validate.go` (all-dev mode).
  - Calls:
    - `pv.findDeveloperDirs()`
    - `pv.validateSingleDeveloper(...)` for each developer
  - Performs conflict detection across collected ports.

- `(pv *PortValidator) validateSingleDeveloper(developerName string) (int, *ValidationError, *ValidationWarning)`:
  - Called by `ValidateAll()`.
  - Calls `config.LoadDeveloperConfig(...)` (which triggers config validation).
  - Applies:
    - missing SSH port warning
    - NodePort range error (`NodePortMin`/`NodePortMax`)

- `(pv *PortValidator) ValidateSingle(developerName string) (*ValidationResult, error)`:
  - Called by `cmd/devenv/validate.go` (single-dev mode).
  - Calls `pv.ValidateAll()` first, then filters results for one user.
  - Calls `pv.errorInvolvesUser(...)`.

- `(pv *PortValidator) errorInvolvesUser(err ValidationError, targetUser string) bool`:
  - Called by `ValidateSingle(...)` during filtering.

- `(pv *PortValidator) findDeveloperDirs() ([]string, error)`:
  - Called by `ValidateAll(...)`.
  - Scans `configDir` and keeps directories containing `devenv-config.yaml`.

## 3) Reachability summary from `cmd/devenv`

Runtime-reached from CLI:

- `NodePortMin`
- `NodePortMax`
- `PortValidator`
- `ValidationResult`
- `ValidationError`
- `ValidationWarning`
- `NewPortValidator`
- `(*PortValidator).ValidateAll`
- `(*PortValidator).validateSingleDeveloper`
- `(*PortValidator).ValidateSingle`
- `(*PortValidator).errorInvolvesUser`
- `(*PortValidator).findDeveloperDirs`
