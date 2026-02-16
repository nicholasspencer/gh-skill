# gh-skill

A [GitHub CLI](https://cli.github.com/) extension for managing AI agent skills stored as GitHub Gists.

## Install

```bash
gh extension install nicholasspencer/gh-skill
```

## Usage

### Install a skill from a gist

```bash
gh skill add https://gist.github.com/user/abc123
gh skill add abc123  # by ID
```

Skills are installed to `~/.gistskills/<name>/` and automatically symlinked to detected AI tool skill directories (Claude Code, OpenClaw, Copilot CLI, Codex, OpenCode).

### List installed skills

```bash
gh skill list
```

### Get skill info

```bash
gh skill info my-skill
```

### Update skills

```bash
gh skill update my-skill
gh skill update --all
```

### Remove a skill

```bash
gh skill remove my-skill
```

### Publish a local skill as a gist

```bash
gh skill publish ./my-skill-folder
```

Creates a secret (unlisted) gist from a folder containing a `SKILL.md`. Gist files are flat — all files in the skill folder root are published directly.

```bash
gh skill publish ./my-skill --public   # make it discoverable
gh skill publish ./my-skill --secret   # explicit secret (default)
```

### Search for skills

```bash
gh skill search "git automation"
```

Searches public gists tagged with `#gistskill`.

### Link to a specific tool

```bash
gh skill link my-skill --target claude-code
gh skill link my-skill --target openclaw
```

## Skill Format

A skill is a GitHub Gist with a `SKILL.md` file containing YAML front matter:

```yaml
---
name: my-skill
description: What this skill does
version: 1.0.0
tags: [automation, git]
author: username
---

# My Skill

Instructions for the AI agent...
```

## Supported Tools

| Tool | Skill Directory | Auto-linked |
|------|----------------|-------------|
| Claude Code | `~/.claude/skills/` | ✓ |
| OpenClaw | `~/.chad/skills/` | ✓ |
| Copilot CLI | `~/.copilot/skills/` | ✓ |
| Codex | `~/.codex/skills/` | ✓ |
| OpenCode | `~/.opencode/skills/` | ✓ |
| Cursor | `.cursor/skills/` | ✗ (project-level) |

## License

MIT
