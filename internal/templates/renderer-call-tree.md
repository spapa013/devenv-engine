# Renderer Call Tree (`cmd/devenv` -> `internal/templates/renderer.go`)

This document maps runtime call paths from CLI entrypoints in `cmd/devenv` to all functions and variables declared in `internal/templates/renderer.go`.

## 1) CLI entrypoints that reach renderer

### Path A: `devenv generate ...`

1. `cmd/devenv/main.go` -> `main()`
2. `cmd/devenv/root.go` -> `generateCmd`
3. `cmd/devenv/generate.go`:
   - system manifests:
     - `generateSystemManifests(...)` -> `templates.NewSystemRenderer(...)` -> `renderer.RenderAll(...)`
   - developer manifests:
     - `generateDeveloperManifests(...)` -> `templates.NewDevRenderer(...)` -> `renderer.RenderAll(...)`

### Path B: `devenv validate ...`

- No usage of `internal/templates` in `devenv validate`.

## 2) Symbol-by-symbol map for `internal/templates/renderer.go`

### Variables

- `devTemplatesToRender`:
  - Used by `NewDevRenderer(...)` as target template list.

- `systemTemplatesToRender`:
  - Used by `NewSystemRenderer(...)` as target template list.

- `templates` (`embed.FS`):
  - Embedded filesystem for all template/script files.
  - Used by:
    - `templateFuncs(...)` (`ReadFile` for scripts)
    - `(*Renderer[T]).RenderTemplate(...)` (`ReadFile` for manifests)

### Types

- `Renderer[T config.BaseConfig | config.DevEnvConfig]`:
  - Generic renderer carrying:
    - `outputDir`
    - `templateRoot`
    - `targetTemplates`

### Functions/methods

- `NewDevRenderer(outputDir string) *Renderer[config.DevEnvConfig]`:
  - Called by `cmd/devenv/generate.go` in developer manifest generation.
  - Delegates to `NewRendererWithFS(...)` with `template_files/dev`.

- `NewSystemRenderer(outputDir string) *Renderer[config.BaseConfig]`:
  - Called by `cmd/devenv/generate.go` in system manifest generation.
  - Delegates to `NewRendererWithFS(...)` with `template_files/system`.

- `NewRendererWithFS[T ...](outputDir, templateRoot string, targetTemplates []string) *Renderer[T]`:
  - Called by `NewDevRenderer(...)` and `NewSystemRenderer(...)`.
  - Initializes a generic `Renderer[T]`.

- `templateFuncs(templateRoot string) template.FuncMap`:
  - Called by `RenderTemplate(...)` to register template helper functions.
  - Provides:
    - `b64enc`
    - `indent`
    - `getTemplatedScript` (recursive template execution for script templates)
    - `getStaticScript`

- `(r *Renderer[T]) RenderTemplate(templateName string, config *T) error`:
  - Called by `RenderAll(...)`.
  - Reads embedded template file, parses with `templateFuncs`, creates output file, executes template.

- `(r *Renderer[T]) RenderAll(config *T) error`:
  - Called by `cmd/devenv/generate.go`.
  - Iterates `targetTemplates` and calls `RenderTemplate(...)` for each.

## 3) Reachability summary from `cmd/devenv`

Runtime-reached from CLI:

- `devTemplatesToRender`
- `systemTemplatesToRender`
- `templates` (`embed.FS`)
- `Renderer[T ...]`
- `NewDevRenderer`
- `NewSystemRenderer`
- `NewRendererWithFS`
- `templateFuncs`
- `(*Renderer[T]).RenderTemplate`
- `(*Renderer[T]).RenderAll`
