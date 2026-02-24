# Validation Call Tree (`cmd/devenv` -> `internal/config/validation.go`)

This document maps runtime call paths from CLI entrypoints in `cmd/devenv` to all functions and variables declared in `internal/config/validation.go`.

## 1) CLI entrypoints that reach config validation

### Path A: `devenv generate ...`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/main.go` -> `rootCmd.Execute()`
3. `cmd/devenv/root.go` -> `generateCmd`
4. `cmd/devenv/generate.go` -> `generateSingleDeveloper(...)` or `processSingleDeveloperForBatchWithError(...)`
5. `cmd/devenv/generate.go` -> `config.LoadDeveloperConfigWithBaseConfig(...)`
6. `internal/config/parser.go` -> `userConfig.Validate()`
7. `internal/config/types.go` -> `(*DevEnvConfig).Validate()`
8. `internal/config/validation.go` -> `ValidateDevEnvConfig(...)`

### Path B: `devenv validate ...`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/main.go` -> `rootCmd.Execute()`
3. `cmd/devenv/root.go` -> `validateCmd`
4. `cmd/devenv/validate.go` -> `validation.NewPortValidator(...)`
5. `cmd/devenv/validate.go` -> `(*PortValidator).ValidateAll()` or `(*PortValidator).ValidateSingle(...)`
6. `internal/validation/ports.go` -> `config.LoadDeveloperConfig(...)`
7. `internal/config/parser.go` -> `config.Validate()`
8. `internal/config/types.go` -> `(*DevEnvConfig).Validate()`
9. `internal/config/validation.go` -> `ValidateDevEnvConfig(...)`

## 2) Symbol-by-symbol map for `internal/config/validation.go`

### Variables

- `validate`:
  - Initialized in `init()`.
  - Used by `ValidateDevEnvConfig(...)` (`validate.Struct(config)`).
  - Used by `ValidateBaseConfig(...)` (`validate.Struct(config)`).

- `sshKeyRegex`:
  - Used in `validateSSHKeys(...)` to verify SSH public key format.

- `numberRe`:
  - Used in `validateKubernetesCPU(...)` for numeric CPU strings.

- `cpuMillicoresRe`:
  - Used in `validateKubernetesCPU(...)` for `"500m"` style values.

- `memoryRe`:
  - Used in `validateKubernetesMemory(...)` for memory quantity strings.

### Functions

- `init()`:
  - Runs automatically when package `config` is imported.
  - Registers custom validators:
    - tag `"ssh_keys"` -> `validateSSHKeys`
    - tag `"k8s_cpu"` -> `validateKubernetesCPU`
    - tag `"k8s_memory"` -> `validateKubernetesMemory`
  - Registers struct-level validator:
    - `GitRepo{}` -> `validateGitRepo`
  - Reached in both CLI paths (A and B) because both import/use `internal/config`.

- `validateSSHKeys(fl validator.FieldLevel) bool`:
  - Invoked by validator engine when a struct field has tag `ssh_keys`.
  - Uses `sshKeyRegex`.
  - Runtime reachable via `ValidateDevEnvConfig(...)` -> `validate.Struct(...)`.

- `validateGitRepo(sl validator.StructLevel)`:
  - Invoked by validator engine for `GitRepo` struct due to `RegisterStructValidation`.
  - Runtime reachable via `ValidateDevEnvConfig(...)` -> `validate.Struct(...)`.

- `validateKubernetesCPU(fl validator.FieldLevel) bool`:
  - Invoked by validator engine for fields tagged `k8s_cpu`.
  - Uses `numberRe` and `cpuMillicoresRe`.
  - Runtime reachable via `ValidateDevEnvConfig(...)` -> `validate.Struct(...)`.

- `validateKubernetesMemory(fl validator.FieldLevel) bool`:
  - Invoked by validator engine for fields tagged `k8s_memory`.
  - Uses `memoryRe`.
  - Runtime reachable via `ValidateDevEnvConfig(...)` -> `validate.Struct(...)`.

- `ValidateDevEnvConfig(config *DevEnvConfig) error`:
  - Main runtime entrypoint in this file.
  - Called from `(*DevEnvConfig).Validate()` in `internal/config/types.go`.
  - Does:
    - Tag/struct validation (`validate.Struct`).
    - SSH semantic checks (`config.GetSSHKeys()` and required >= 1).
    - Resource semantic checks (`getCanonicalCPU`, `getCanonicalMemory`, GPU >= 0).
    - Error formatting through `formatValidationError(...)` when tag validation fails.

- `ValidateBaseConfig(config *BaseConfig) error`:
  - Uses `validate.Struct(...)` + `formatValidationError(...)`.
  - Not currently called by `cmd/devenv` runtime paths.
  - Presently used in tests.

- `formatValidationError(err error) error`:
  - Called by:
    - `ValidateDevEnvConfig(...)`
    - `ValidateBaseConfig(...)`
  - Iterates validation errors and delegates per-field message formatting to `formatFieldError(...)`.

- `formatFieldError(fieldError validator.FieldError) string`:
  - Called only by `formatValidationError(...)`.
  - Produces user-facing text per validation tag.

## 3) Reachability summary from `cmd/devenv`

Runtime-reached from CLI:

- `init`
- `validateSSHKeys`
- `validateGitRepo`
- `validateKubernetesCPU`
- `validateKubernetesMemory`
- `ValidateDevEnvConfig`
- `formatValidationError`
- `formatFieldError`
- variables: `validate`, `sshKeyRegex`, `numberRe`, `cpuMillicoresRe`, `memoryRe`

Not reached from CLI runtime (currently):

- `ValidateBaseConfig` (test-only in current codebase)
