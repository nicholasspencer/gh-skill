# GistSkills — PRD

**Status:** Draft  
**Author:** Nico + Chad  
**Date:** 2026-02-15  

---

## Problem

AI agent skills (SKILL.md + supporting files) are becoming a universal standard across tools (Claude Code, OpenClaw, Codex, Cursor, OpenCode, etc.). But sharing them is fragmented:

- ClawHub is OpenClaw-specific
- skillshare requires a Git repo + CLI
- Most skills are just 1-5 markdown/text files in a folder

There's no lightweight, universal way to share a skill. You end up copy-pasting files between repos or publishing to tool-specific registries.

## Insight

**A GitHub Gist _is_ a skill folder.** Gists already support:
- Multiple files (SKILL.md + scripts/ + references/)
- Versioning (every edit creates a revision)
- Raw URLs (for programmatic fetching)
- Forks (for customization)
- Stars (for discovery)
- API access (no auth needed for public gists)
- Embeddable (for docs/blogs)

No new infrastructure needed. GitHub already built the registry.

## Solution

**GistSkills** — a convention + thin CLI for using GitHub Gists as universal AI agent skills.

### The Convention

A GistSkill is a public GitHub Gist where:
1. **`SKILL.md`** is the primary file (required) — contains instructions + front matter
2. Additional gist files map to the skill folder structure (scripts, references, etc.)
3. The gist description serves as the skill's one-liner summary

That's it. No build step, no manifest, no registry submission.

### File Mapping

```
Gist files:                    →  Installed as:
─────────────────────────────────────────────────
SKILL.md                       →  my-skill/SKILL.md
setup.sh                       →  my-skill/scripts/setup.sh
api-reference.md               →  my-skill/references/api-reference.md
```

**Convention for subdirectories** (gists are flat): prefix filenames with the directory path using `--` as separator:
```
scripts--setup.sh              →  my-skill/scripts/setup.sh
references--api-docs.md        →  my-skill/references/api-docs.md
```

### Front Matter (in SKILL.md)

Standard YAML front matter, compatible with existing skill specs:

```yaml
---
name: my-skill
description: One-line description
version: 1.0.0
tags: [automation, git, productivity]
tools: [claude-code, openclaw, codex, cursor]  # optional compatibility hints
author: nicholasspencer
---
```

### CLI — `gh skill` (gh extension)

Distributed as a `gh` extension (`gh-skill`). Zero new tooling — piggybacks on GitHub CLI's auth, update system, and extension ecosystem.

```bash
# Install the extension
gh extension install nicholasspencer/gh-skill

# Install a skill from a gist URL
gh skill add https://gist.github.com/nico/abc123

# Install by gist ID
gh skill add abc123

# List installed skills
gh skill list

# Update a skill (pulls latest gist revision)
gh skill update my-skill

# Update all
gh skill update --all

# Remove
gh skill remove my-skill

# Publish a local skill folder as a gist
gh skill publish ./my-skill

# Search (uses GitHub Gist search + tag convention)
gh skill search "git automation"

# Info about an installed skill
gh skill info my-skill

# Link into a specific tool's skill directory
gh skill link my-skill --target claude-code
gh skill link my-skill --target openclaw
gh skill link my-skill --target cursor
```

### Installation Targets

Skills install to `~/.gistskills/` by default, then get symlinked into tool-specific directories:

```
~/.gistskills/
  my-skill/
    SKILL.md
    scripts/
      setup.sh
    .gistskill.json          # metadata (gist ID, version, installed date)
```

Auto-detection of installed tools:
- **Claude Code:** `~/.claude/skills/`  
- **OpenClaw:** `~/.chad/skills/` (or configured path)
- **Cursor:** `.cursor/skills/` (project-level)
- **Codex:** `.codex/skills/`
- **OpenCode:** `.opencode/skills/`

`gistskill add` installs + links to all detected tools by default. `--target` to be specific.

### Discovery

No centralized registry needed. Discovery happens through:

1. **GitHub Gist search** — search by description/filename
2. **Tag convention** — include `#gistskill` in gist description for discoverability
3. **Curated lists** — a gist can itself be an index (list of gist URLs)
4. **Stars** — GitHub's existing starring system
5. **Social sharing** — just share the URL

Optional: A simple static site (GitHub Pages) that indexes public gists tagged `#gistskill`. No auth, no API keys, just scraping the GitHub API periodically.

### Collections

A "collection" is just a gist containing a `gistskills.json`:

```json
{
  "name": "Productivity Pack",
  "description": "Essential skills for daily workflows",
  "skills": [
    "https://gist.github.com/user/abc123",
    "https://gist.github.com/user/def456",
    "https://gist.github.com/user/ghi789"
  ]
}
```

```bash
gistskill add-collection https://gist.github.com/nico/mycollection
```

## Technical Details

### Implementation

- **Distribution:** `gh` extension (`gh-skill` / `nicholasspencer/gh-skill`)
- **Language:** Go (matches `gh` ecosystem, compiles to single binary via `gh extension create --precompiled`)
- **Auth:** Inherited from `gh auth` — no separate token management
- **API:** Uses `gh api` under the hood for gist CRUD
- **Config:** `~/.gistskills/config.json`
- **State:** `~/.gistskills/` directory is the entire state (no database)

### Security Considerations

- Skills are text files (markdown, shell scripts) — same trust model as any skill system
- Scripts are NOT auto-executed on install (user must explicitly run them)
- `gistskill audit <skill>` shows all files + any executable content
- Private gists supported with GitHub token auth
- Checksum verification: `.gistskill.json` stores the commit SHA for integrity

### Rate Limits

- GitHub API: 60 req/hr unauthenticated, 5000/hr authenticated
- Sufficient for individual use; for the index site, use authenticated + caching

## Non-Goals (v1)

- No accounts or auth system
- No comments/reviews (use GitHub's gist comments)
- No paid skills or monetization
- No auto-execution of scripts on install
- No dependency resolution between skills

## Success Metrics

- Can install a skill in <5 seconds from a URL
- Can publish a skill in <10 seconds from a folder
- Works with 0 configuration for the common case
- No signup, no API keys, no accounts needed

## Decisions

1. ~~CLI name~~ → `gh skill` (gh extension)
2. ~~GitLab support~~ → GitHub-only via `gh` (v1)
3. ~~Subdirectory convention~~ → `--` separator for portability
4. `gh skill publish` creates the gist via `gh api` (auth already available)

## Open Questions

1. Should `gh skill publish` default to public or private gists?
2. Extension language: Go (precompiled) or shell script (simpler but less portable)?
3. Should we reserve a GitHub topic (`gh-skill`?) for discoverability alongside `#gistskill` in descriptions?

---

*"The best package manager is a URL."*
