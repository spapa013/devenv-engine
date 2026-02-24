# Resources Call Tree (`cmd/devenv` -> `internal/config/resources.go`)

This document maps runtime call paths from CLI entrypoints in `cmd/devenv` to all functions and variables declared in `internal/config/resources.go`.

## 1) CLI entrypoints that reach resource normalization

### Path A: `devenv generate ...`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/root.go` -> `generateCmd`
3. `cmd/devenv/generate.go` loads config via parser functions
4. Validation path during load:
   - `(*DevEnvConfig).Validate()` -> `ValidateDevEnvConfig(...)`
   - `ValidateDevEnvConfig(...)` calls:
     - `config.Resources.getCanonicalCPU()`
     - `config.Resources.getCanonicalMemory()`
5. Rendering/summary path:
   - `cfg.CPU()` -> `getCanonicalCPU()`
   - `cfg.Memory()` -> `getCanonicalMemory()`

### Path B: `devenv validate ...`

1. `cmd/devenv/validate.go` -> `PortValidator`
2. `internal/validation/ports.go` -> `config.LoadDeveloperConfig(...)`
3. During load, `config.Validate()` -> `ValidateDevEnvConfig(...)`
4. `ValidateDevEnvConfig(...)` calls:
   - `getCanonicalCPU()`
   - `getCanonicalMemory()`

## 2) Symbol-by-symbol map for `internal/config/resources.go`

### Variables

- `binToMi`:
  - Used by `memoryTextToMi(...)` for binary-unit conversion (Ki/Mi/Gi/...).

- `decBytesToMi`:
  - Used by `memoryTextToMi(...)` for decimal-byte conversion (k/M/G/...).

### Functions/methods

- `normalizeToCPUText(v any) (string, error)`:
  - Called by `(*ResourceConfig).getCanonicalCPU()`.

- `cpuTextToMillicores(s string) (int64, error)`:
  - Called by `(*ResourceConfig).getCanonicalCPU()`.

- `(r *ResourceConfig) getCanonicalCPU() (int64, error)`:
  - Main CPU entrypoint in this file.
  - Called by:
    - `(*DevEnvConfig).CPU()` in `types.go`
    - `ValidateDevEnvConfig(...)` in `validation.go`

- `normalizeToMemoryText(v any) (string, error)`:
  - Called by `(*ResourceConfig).getCanonicalMemory()`.
  - Uses `hasSuffixFold(...)`.

- `memoryTextToMi(s string) (int64, error)`:
  - Called by `(*ResourceConfig).getCanonicalMemory()`.
  - Uses `binToMi`, `decBytesToMi`, `bytesToMi(...)`, and `roundFloatToInt64(...)`.

- `(r *ResourceConfig) getCanonicalMemory() (int64, error)`:
  - Main memory entrypoint in this file.
  - Called by:
    - `(*DevEnvConfig).Memory()` in `types.go`
    - `ValidateDevEnvConfig(...)` in `validation.go`

- `bytesToMi(bytes float64) (int64, error)`:
  - Called by `memoryTextToMi(...)`.

- `roundFloatToInt64(v float64) (int64, error)`:
  - Called by:
    - `memoryTextToMi(...)`
    - `bytesToMi(...)`

- `hasSuffixFold(s, suf string) bool`:
  - Called by `normalizeToMemoryText(...)`.

## 3) Reachability summary from `cmd/devenv`

Runtime-reached from CLI:

- `normalizeToCPUText`
- `cpuTextToMillicores`
- `(*ResourceConfig).getCanonicalCPU`
- `normalizeToMemoryText`
- `memoryTextToMi`
- `(*ResourceConfig).getCanonicalMemory`
- `bytesToMi`
- `roundFloatToInt64`
- `hasSuffixFold`
- variables: `binToMi`, `decBytesToMi`
