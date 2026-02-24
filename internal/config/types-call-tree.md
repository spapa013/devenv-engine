# Types Call Tree (`cmd/devenv` -> `internal/config/types.go`)

This document maps runtime call paths from CLI entrypoints in `cmd/devenv` to type constructors/methods declared in `internal/config/types.go`.

## 1) CLI entrypoints that reach `types.go`

### Path A: `devenv generate ...`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/main.go` -> `rootCmd.Execute()`
3. `cmd/devenv/root.go` -> `generateCmd`
4. `cmd/devenv/generate.go`:
   - `config.LoadGlobalConfig(...)` -> `NewBaseConfigWithDefaults()`
   - `config.LoadDeveloperConfigWithBaseConfig(...)` -> `(*DevEnvConfig).Validate()`
   - Optional summary path calls methods like `GetSSHKeys()`, `CPU()`, `Memory()`, `GetDeveloperDir()`
5. Template rendering path (`templates.RenderAll`) evaluates template methods from `*DevEnvConfig`:
   - `GPU()`, `CPU()`, `CPURequest()`, `Memory()`, `MemoryRequest()`, `NodePort()`, `GetUserID()`, `GetSSHKeysString()`

### Path B: `devenv validate ...`

1. `cmd/devenv/validate.go` -> `validation.PortValidator`
2. `internal/validation/ports.go` -> `config.LoadDeveloperConfig(...)`
3. `internal/config/parser.go` -> `config.Validate()`
4. `internal/config/types.go` -> `(*DevEnvConfig).Validate()`

## 2) Symbol-by-symbol map for `internal/config/types.go`

### Type declarations

- `BaseConfig`
- `DevEnvConfig`
- `GitConfig`
- `PackageConfig`
- `GitRepo`
- `ResourceConfig`
- `VolumeMount`
- `RefreshConfig`

All of the above are runtime data models used by YAML unmarshal, validation, and template rendering paths.

### Functions and methods

- `NewBaseConfigWithDefaults() BaseConfig`:
  - Called by `LoadGlobalConfig(...)` in `internal/config/parser.go`.

- `(c *BaseConfig) GetSSHKeys() ([]string, error)`:
  - Calls `normalizeSSHKeys(...)` from `parser.go`.
  - Used by parser merge logic, validation checks, and CLI summary output.

- `(c *DevEnvConfig) GetDeveloperDir() string`:
  - Used by `cmd/devenv/generate.go` summary output.

- `(c *DevEnvConfig) GetUserID() string`:
  - Used by templates:
    - `internal/templates/template_files/dev/manifests/env-vars.tmpl`
    - `internal/templates/template_files/dev/scripts/templated/startup.sh`

- `(c *DevEnvConfig) GPU() int`:
  - Used by `statefulset.tmpl`.

- `(c *DevEnvConfig) CPU() string`:
  - Used by `cmd/devenv/generate.go` summary output.
  - Used by `statefulset.tmpl`.
  - Calls `c.Resources.getCanonicalCPU()` from `resources.go`.

- `(c *DevEnvConfig) Memory() string`:
  - Used by `cmd/devenv/generate.go` summary output.
  - Used by `statefulset.tmpl`.
  - Calls `c.Resources.getCanonicalMemory()` from `resources.go`.

- `(c *DevEnvConfig) CPURequest() string`:
  - Used by `statefulset.tmpl`.
  - Delegates to `CPU()`.

- `(c *DevEnvConfig) MemoryRequest() string`:
  - Used by `statefulset.tmpl`.
  - Delegates to `Memory()`.

- `(c *DevEnvConfig) NodePort() int`:
  - Used by `internal/templates/template_files/dev/manifests/service.tmpl`.

- `(c *DevEnvConfig) VolumeMounts() []VolumeMount`:
  - No current runtime call sites from `cmd/devenv`; presently test/doc-facing.

- `(c *DevEnvConfig) GetSSHKeysSlice() []string`:
  - Called by `GetSSHKeysString()`.
  - No independent template/runtime call site currently.

- `(c *DevEnvConfig) GetSSHKeysString() string`:
  - Used by `internal/templates/template_files/dev/scripts/templated/startup.sh`.

- `(c *DevEnvConfig) Validate() error`:
  - Called by parser loading functions:
    - `LoadDeveloperConfig(...)`
    - `LoadDeveloperConfigWithBaseConfig(...)`
  - Delegates to `ValidateDevEnvConfig(...)` in `validation.go`.

## 3) Reachability summary from `cmd/devenv`

Runtime-reached from CLI:

- `NewBaseConfigWithDefaults`
- `(*BaseConfig).GetSSHKeys`
- `(*DevEnvConfig).GetDeveloperDir`
- `(*DevEnvConfig).GetUserID`
- `(*DevEnvConfig).GPU`
- `(*DevEnvConfig).CPU`
- `(*DevEnvConfig).Memory`
- `(*DevEnvConfig).CPURequest`
- `(*DevEnvConfig).MemoryRequest`
- `(*DevEnvConfig).NodePort`
- `(*DevEnvConfig).GetSSHKeysSlice` (indirect through `GetSSHKeysString`)
- `(*DevEnvConfig).GetSSHKeysString`
- `(*DevEnvConfig).Validate`

Currently not runtime-reached from `cmd/devenv`:

- `(*DevEnvConfig).VolumeMounts`
