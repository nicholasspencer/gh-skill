# Agent Instructions

> **Want to build and publish skills?** This project ships a skill for that.
> Read [skills/skill-creator/SKILL.md](skills/skill-creator/SKILL.md) — it teaches you the full skill authoring workflow using `gh skill`.

---

## Architecture

```
gh-skill (Go, cobra CLI)
├── main.go              → entry point
├── cmd/                 → cobra commands (one file per subcommand)
│   ├── root.go          → root command + subcommand registration
│   ├── add.go           → install skill from gist (with trust gate)
│   ├── install.go       → download skill files to cwd (no linking)
│   ├── publish.go       → publish local folder as gist
│   ├── list.go          → list installed skills
│   ├── info.go          → show skill metadata
│   ├── update.go        → pull latest gist revision
│   ├── remove.go        → uninstall skill + remove symlinks
│   ├── search.go        → search public gists tagged [gh-skill]
│   ├── link.go          → symlink skill to a tool's skill dir
│   ├── fork.go          → fork a gist skill
│   └── trust.go         → manage trusted authors
└── internal/            → shared logic
    ├── skill.go         → front matter parsing, install/list/remove
    ├── gist.go          → GitHub Gist API (fetch, create, search via `gh api`)
    ├── linking.go       → tool detection, symlink management
    └── trust.go         → trust store, trust prompt UI
```

## Key Concepts

- **Skills** are GitHub Gists containing a `*.skill.md` (or legacy `SKILL.md`) file with YAML front matter
- **Gist file naming**: subdirectories flattened with `--` separator (e.g., `scripts--setup.sh` → `scripts/setup.sh`)
- **On publish**: `SKILL.md` → `<name>.skill.md`; **on install**: `<name>.skill.md` → `SKILL.md`
- **Trust gate**: untrusted authors prompt before install; own gists + trusted authors skip
- **Auto-linking**: `gh skill add` symlinks to all detected tool skill directories
- **State**: all in `~/.gistskills/` (skill folders + `.gistskill.json` metadata + `trusted-authors.json`)

## Supported Tool Targets

| Name | Directory | Auto-linked |
|------|-----------|-------------|
| claude-code | `~/.claude/skills/` | ✓ |
| openclaw | `~/.chad/skills/` (or config) | ✓ |
| copilot | `~/.copilot/skills/` | ✓ |
| codex | `~/.codex/skills/` | ✓ |
| opencode | `~/.opencode/skills/` | ✓ |
| cursor | `.cursor/skills/` (project) | manual only |

To add a new tool target, update `DetectToolDirs()` and `KnownTools()` in `internal/linking.go`, plus the error string in `ToolDirByName()` and help text in `cmd/link.go`.

## Development

```bash
# Build
go build -o gh-skill .

# Install as gh extension (local dev)
gh extension install .

# Test
go test ./...
```

### Adding a New Subcommand

1. Create `cmd/<name>.go` with a `cobra.Command`
2. Register it in `cmd/root.go` → `init()` → `rootCmd.AddCommand()`
3. Put reusable logic in `internal/`

### Conventions

- **Commits**: [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, etc.)
- **Dependencies**: keep minimal — `cobra` for CLI, `yaml.v3` for front matter, stdlib for everything else
- **All GitHub API calls** go through `gh api` (inherits auth, no token management)
- **No auto-execution** of scripts on install — security boundary

## Discoverability

### The Extension Itself

`gh-skill` is discoverable as a `gh` extension via:

```bash
gh extension search gh-skill
gh search repos --topic gh-extension gh-skill
```

Ensure the repo has the `gh-extension` topic on GitHub. The repo description shows in search results.

### Skills (Gist-based)

Published skills use `[gh-skill]` prefix in the gist description. `gh skill search` filters public gists by this prefix + presence of a `*.skill.md` file.

### Skills (Repo-based, future)

Repos tagged with the `gh-skill` GitHub topic could be searched via:

```bash
gh search repos --topic gh-skill --json fullName,description,stargazersCount,url,license,updatedAt
```

Not yet implemented in the CLI but the convention is documented for forward compatibility.

### `gh extension search` JSON Fields

Available via `--json`: `fullName`, `description`, `stargazersCount`, `forksCount`, `license`, `owner`, `url`, `createdAt`, `updatedAt`, `language`, `defaultBranch`, `isArchived`, `visibility`, `watchersCount`, `openIssuesCount`, `size`.

## Testing Changes

After modifying code, always run `go build -o gh-skill .` to verify compilation. For behavioral changes, test the relevant command end-to-end with a real gist when possible.
