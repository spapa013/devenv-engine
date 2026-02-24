# Git Info Call Tree (`cmd/devenv` -> `internal/git/info.go`)

This document maps runtime call paths from CLI entrypoints in `cmd/devenv` to all functions and variables declared in `internal/git/info.go`.

## 1) CLI entrypoints that reach `internal/git/info.go`

Current state: there is no runtime import/use of `internal/git` from `cmd/devenv`.

- `cmd/devenv/generate.go` does not import `internal/git`.
- `cmd/devenv/validate.go` does not import `internal/git`.
- `cmd/devenv/version.go` does not import `internal/git`.

So there is no call path from CLI entrypoints to this package at present.

## 2) Symbol-by-symbol map for `internal/git/info.go`

### Type declarations

- `GitInfo`:
  - Holds:
    - `Repository *git.Repository`
    - `CommitHash string`
    - `Branch string`
    - `Tag []string`
    - `IsDirty bool`
  - Constructed and returned by `GetGitInfo(...)`.

### Functions

- `GetGitInfo(repoPath string) (*GitInfo, error)`:
  - Standalone function that:
    - opens a repository (`git.PlainOpenWithOptions`)
    - resolves `HEAD` commit and branch
    - scans tags pointing to current commit
    - checks worktree dirty status
    - returns a populated `*GitInfo`
  - Current usage in repo: tests only (`internal/git/info_test.go`).

## 3) Reachability summary from `cmd/devenv`

Runtime-reached from CLI:

- none

Used in repository:

- `GetGitInfo(...)` and `GitInfo` are used in `internal/git/info_test.go` (test-only).
