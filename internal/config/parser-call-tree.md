# Parser Call Tree (`cmd/devenv` -> `internal/config/parser.go`)

This document maps runtime call paths from CLI entrypoints in `cmd/devenv` to all functions declared in `internal/config/parser.go`.

## 1) CLI entrypoints that reach parser functions

### Path A: `devenv generate ...`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/main.go` -> `rootCmd.Execute()`
3. `cmd/devenv/root.go` -> `generateCmd`
4. `cmd/devenv/generate.go` -> `generateSingleDeveloper(...)` or `processSingleDeveloperForBatchWithError(...)`
5. `cmd/devenv/generate.go` -> `config.LoadGlobalConfig(...)`
6. `cmd/devenv/generate.go` -> `config.LoadDeveloperConfigWithBaseConfig(...)`

### Path B: `devenv validate ...`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/main.go` -> `rootCmd.Execute()`
3. `cmd/devenv/root.go` -> `validateCmd`
4. `cmd/devenv/validate.go` -> `validation.NewPortValidator(...)`
5. `internal/validation/ports.go` -> `config.LoadDeveloperConfig(...)`

## 2) Symbol-by-symbol map for `internal/config/parser.go`

- `LoadGlobalConfig(configDir string) (*BaseConfig, error)`:
  - Called directly from `cmd/devenv/generate.go`.
  - Calls `NewBaseConfigWithDefaults()` from `internal/config/types.go`.

- `LoadDeveloperConfig(configDir, developerName string) (*DevEnvConfig, error)`:
  - Called by `internal/validation/ports.go` (used by `devenv validate`).
  - Calls `config.Validate()` (`(*DevEnvConfig).Validate`) after YAML unmarshal.

- `LoadDeveloperConfigWithBaseConfig(configDir, developerName string, baseConfig *BaseConfig) (*DevEnvConfig, error)`:
  - Called directly from `cmd/devenv/generate.go`.
  - Calls `userConfig.mergeListFields(baseConfig)`.
  - Calls `userConfig.Validate()` (`(*DevEnvConfig).Validate`) after merge.

- `(config *DevEnvConfig) mergeListFields(globalConfig *BaseConfig)`:
  - Called only by `LoadDeveloperConfigWithBaseConfig(...)`.
  - Calls:
    - `mergeStringSlices(...)` (packages and SSH keys)
    - `mergeVolumes(...)`
    - `globalConfig.GetSSHKeys()` and `config.GetSSHKeys()`

- `mergeStringSlices(global, user []string) []string`:
  - Called by `mergeListFields(...)`.
  - Not called directly outside parser internals in runtime.

- `mergeVolumes(global, user []VolumeMount) []VolumeMount`:
  - Called by `mergeListFields(...)`.
  - Not called directly outside parser internals in runtime.

- `normalizeSSHKeys(sshKeyField any) ([]string, error)`:
  - Indirect runtime usage via `(*BaseConfig).GetSSHKeys()` in `internal/config/types.go`.
  - Reached from:
    - parser merging (`mergeListFields`)
    - validation path (`validateSSHKeys` and `ValidateDevEnvConfig`)
    - `cmd/devenv/generate.go` config summary (`cfg.GetSSHKeys()`).

## 3) Reachability summary from `cmd/devenv`

Runtime-reached from CLI:

- `LoadGlobalConfig`
- `LoadDeveloperConfig`
- `LoadDeveloperConfigWithBaseConfig`
- `(*DevEnvConfig).mergeListFields`
- `mergeStringSlices`
- `mergeVolumes`
- `normalizeSSHKeys` (indirect)

No package-level variables are declared in `internal/config/parser.go`.
