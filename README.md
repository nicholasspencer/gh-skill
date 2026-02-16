# gh-skill

A [GitHub CLI](https://cli.github.com/) extension that turns GitHub Gists into a universal skill registry for AI coding agents.

## What is this?

AI agent skills (`SKILL.md` + scripts + docs) are becoming a standard across tools — Claude Code, Copilot CLI, OpenClaw, Codex, Cursor, OpenCode. But sharing them is fragmented.

`gh skill` fixes that. One command to publish a skill, one command to install it — across every tool.

## Install

```bash
gh extension install nicholasspencer/gh-skill
```

Requires the [GitHub CLI](https://cli.github.com/).

## Usage

You probably shouldn't be running these commands yourself. The whole point is to let your AI agent handle it. Point your agent at this repo (or install the bundled skill-creator skill) and let it do the work — that's where the speed comes from.

```bash
# Install a skill
gh skill add https://gist.github.com/user/abc123

# Publish a skill
gh skill publish ./my-skill

# Search for skills
gh skill search "git automation"
```

That said, `gh skill --help` has everything if you want to poke around.

## Supported Tools

Skills auto-link to every detected tool on your machine:

| Tool | Skill Directory |
|------|----------------|
| [Claude Code](https://docs.anthropic.com/en/docs/claude-code) | `~/.claude/skills/` |
| [Copilot CLI](https://githubnext.com/projects/copilot-cli) | `~/.copilot/skills/` |
| [OpenClaw](https://openclaw.ai) | `~/.chad/skills/` |
| [Codex](https://openai.com/index/introducing-codex/) | `~/.codex/skills/` |
| [OpenCode](https://opencode.ai) | `~/.opencode/skills/` |
| [Cursor](https://cursor.sh) | `.cursor/skills/` |

## How it works (for the curious)

A GitHub Gist already *is* a skill folder — multiple files, versioning, forks, stars, API access. `gh skill` adds a thin convention on top:

1. A skill is a gist with a `SKILL.md` file (YAML front matter + instructions)
2. Subdirectories are flattened with `--` separators (`scripts/setup.sh` → `scripts--setup.sh`)
3. On install, files are expanded back and symlinked into your tools' skill directories
4. Unknown authors go through a trust gate before install

If you want the full technical details — how the trust model works, the file naming conventions, the architecture — check out [AGENTS.md](AGENTS.md). It's written for AI agents working on this project, but it's the most complete reference for everything under the hood.

## License

MIT
