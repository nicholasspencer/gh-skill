# gh-skill

A [GitHub CLI](https://cli.github.com/) extension that turns GitHub Gists into a universal skill registry for AI coding agents.

**The idea:** A GitHub Gist already *is* a skill folder — multiple files, versioning, forks, stars, raw URLs, API access. No new infrastructure needed. `gh skill` adds a thin convention on top.

## Why?

AI agent skills (SKILL.md + scripts + docs) are becoming a standard across tools — Claude Code, Copilot CLI, OpenClaw, Codex, Cursor, OpenCode. But sharing them is fragmented. You end up copy-pasting files between repos or publishing to tool-specific registries.

`gh skill` gives you one command to publish a skill and one command to install it — everywhere.

## Install

```bash
gh extension install nicholasspencer/gh-skill
```

Requires the [GitHub CLI](https://cli.github.com/) (`gh`).

## Quick Start

```bash
# Install a skill from a gist
gh skill add https://gist.github.com/user/abc123
gh skill add abc123  # or just the ID

# Publish your own skill
gh skill publish ./my-skill          # secret (unlisted) by default
gh skill publish ./my-skill --public # discoverable via search

# That's it. The output gives you a shareable install command:
# ✓ Published: https://gist.github.com/you/abc123
#   Install with: gh skill add abc123
```

## How It Works

1. **You write a skill** — a folder with a `SKILL.md` (YAML front matter + instructions) and optionally scripts, references, or assets.
2. **You publish it** — `gh skill publish` creates a gist. Subdirectories are flattened using `--` separators (`scripts/setup.sh` → `scripts--setup.sh`).
3. **Anyone installs it** — `gh skill add <gist-id>` downloads, expands paths, and symlinks the skill into every detected AI tool on the machine.

### Trust & Security

Skills from unknown authors go through a trust gate — you see the files, any scripts are flagged, and you choose whether to install, trust the author for future installs, or abort. Your own gists and trusted authors skip the prompt.

```bash
gh skill trust list              # see who you trust
gh skill trust add <username>    # trust an author
gh skill trust remove <username> # revoke trust
```

## Commands

| Command | What it does |
|---------|-------------|
| `gh skill add <gist>` | Install a skill + auto-link to tools |
| `gh skill install <gist>` | Download skill files to a directory (no linking) |
| `gh skill publish <path>` | Publish a local skill folder as a gist |
| `gh skill list` | List installed skills |
| `gh skill info <name>` | Show skill metadata |
| `gh skill update <name>` | Pull latest gist revision |
| `gh skill update --all` | Update all installed skills |
| `gh skill remove <name>` | Uninstall + remove symlinks |
| `gh skill search <query>` | Search public gists tagged `[gh-skill]` |
| `gh skill link <name> --target <tool>` | Symlink to a specific tool |
| `gh skill fork <gist>` | Fork a skill gist for customization |

## Supported Tools

Installed skills auto-link to every detected tool:

| Tool | Skill Directory |
|------|----------------|
| [Claude Code](https://docs.anthropic.com/en/docs/claude-code) | `~/.claude/skills/` |
| [OpenClaw](https://openclaw.ai) | `~/.chad/skills/` |
| [Copilot CLI](https://githubnext.com/projects/copilot-cli) | `~/.copilot/skills/` |
| [Codex](https://openai.com/index/introducing-codex/) | `~/.codex/skills/` |
| [OpenCode](https://opencode.ai) | `~/.opencode/skills/` |
| [Cursor](https://cursor.sh) | `.cursor/skills/` (project-level, manual) |

## Skill Format

A skill is a folder with a `SKILL.md`:

```yaml
---
name: my-skill
description: What this skill does and when to use it
version: 1.0.0
tags: [automation, git]
author: username
---

# My Skill

Instructions for the AI agent...
```

Add `scripts/`, `references/`, or `assets/` directories as needed. See the [skill-creator skill](skills/skill-creator/SKILL.md) for the full authoring guide.

## For AI Agents

See [AGENTS.md](AGENTS.md) for development instructions if you're an AI agent contributing to this project.

## License

MIT
