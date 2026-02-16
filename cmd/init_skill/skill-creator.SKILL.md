---
name: gh-skill-creator
description: "Create, structure, and publish AI agent skills as GitHub Gists using `gh skill`. Use when building a new skill, packaging an existing skill folder for distribution, or helping someone publish a skill to the gist-based skill registry."
---

# Skill Creator (gh-skill)

Create and publish AI agent skills backed by GitHub Gists using the `gh skill` CLI extension.

## What is a Skill?

A skill is a self-contained folder that gives an AI agent specialized knowledge for a domain or task. At minimum it contains a `SKILL.md` with YAML front matter. Optionally it includes scripts, references, and assets.

## Skill Structure

```
my-skill/
├── SKILL.md              # Required — instructions + front matter
├── scripts/              # Optional — deterministic code (shell, python, etc.)
├── references/           # Optional — docs loaded into context on demand
└── assets/               # Optional — templates, images, files used in output
```

## Writing SKILL.md

### Front Matter (required)

```yaml
---
name: my-skill
description: "What this skill does and WHEN to use it. Be specific — this is the trigger."
---
```

Only `name` and `description` are required. Optional fields: `version`, `tags`, `author`, `tools`.

**The description is critical.** It's how agents decide whether to load the skill. Include:
- What the skill does
- Specific triggers / contexts for activation
- Example phrases that should match

### Body (required)

Concise instructions. The agent is already smart — only include what it doesn't already know.

**Guidelines:**
- Imperative form ("Run X", "Check Y"), not conversational
- Prefer examples over explanations
- Keep under 500 lines; split into reference files if longer
- Don't duplicate info between SKILL.md and reference files

## Subdirectory Convention (Gists are flat)

Gists don't support folders. Use `--` as a path separator in filenames:

| Local path | Gist filename |
|---|---|
| `scripts/setup.sh` | `scripts--setup.sh` |
| `references/api.md` | `references--api.md` |

`gh skill publish` handles this automatically. `gh skill add` expands them back.

## Creating a Skill

### 1. Create the folder

```bash
mkdir -p my-skill
```

### 2. Write SKILL.md

Start with front matter + core instructions. Add reference files for anything large or conditional.

### 3. Add resources as needed

- **scripts/**: Code that would be rewritten every time without the skill
- **references/**: Domain docs, schemas, API specs (loaded on demand)
- **assets/**: Templates, images, boilerplate (used in output, not loaded into context)

### 4. Test locally

Point your agent tool at the skill folder and try real tasks. Iterate on SKILL.md based on where the agent struggles.

### 5. Publish

```bash
gh skill publish ./my-skill             # secret gist (default)
gh skill publish ./my-skill --public    # discoverable via search
```

This creates a gist with `[gh-skill]` prefix in the description. The output gives you the install command:

```
✓ Published: https://gist.github.com/you/abc123
  Install with: gh skill add abc123
```

## Installing Skills

```bash
gh skill add <gist-url-or-id>       # install + auto-link to detected tools
gh skill add abc123 --yes            # skip trust prompt
gh skill install abc123 -o ./skills  # download to a specific directory (no linking)
```

## Updating & Managing

```bash
gh skill update my-skill    # pull latest gist revision
gh skill update --all       # update everything
gh skill list               # show installed skills
gh skill info my-skill      # show metadata
gh skill remove my-skill    # uninstall + remove symlinks
gh skill link my-skill --target copilot  # link to specific tool
```

## Progressive Disclosure

Keep SKILL.md lean. Use reference files for variant-specific or deep content:

```markdown
# My Skill

## Quick Start
[core workflow here]

## Advanced
- **AWS specifics**: See references/aws.md
- **GCP specifics**: See references/gcp.md
```

The agent loads reference files only when relevant to the task.

## Anti-Patterns

- **Don't** include README.md, CHANGELOG.md, or meta-docs — the skill is for agents, not humans
- **Don't** auto-execute scripts on install
- **Don't** put "When to use this skill" in the body — put it in the `description` front matter
- **Don't** nest references more than one level deep from SKILL.md
- **Don't** duplicate information across SKILL.md and reference files

## Supported Tools

Skills auto-link to these tools when detected:

| Tool | Directory |
|------|-----------|
| Claude Code | `~/.claude/skills/` |
| OpenClaw | Per-agent (from config) |
| Copilot CLI | `~/.copilot/skills/` |
| Codex | `~/.codex/skills/` |
| OpenCode | `~/.opencode/skills/` |
| Cursor | `.cursor/skills/` (manual) |
