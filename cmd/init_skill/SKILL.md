---
name: gh-skill
description: "Search, install, and manage AI agent skills from the gh-skill registry. Use when the user asks for a capability you don't have, wants to find skills, install new skills, or manage installed skills."
---

# gh-skill — Skill Registry

Search, install, and manage AI agent skills using the `gh skill` CLI extension.

## When to Use

- User asks for a capability you don't have ("I need a weather skill", "find me a skill for X")
- User wants to browse available skills
- User asks to install, update, or remove skills
- You need a specialized skill for a task and suspect one exists

## Searching for Skills

```bash
gh skill search <query>          # search public skills by keyword
gh skill search weather          # example: find weather skills
gh skill search --tag <tag>      # filter by tag
```

Review results and suggest the best match. Check description, author, and star count.

## Installing Skills

```bash
gh skill add <gist-url-or-id>            # install + auto-link to all detected tools
gh skill add <gist-url-or-id> --yes      # skip trust prompt
```

After install, the skill is immediately available in your skill directory. No restart needed.

## Managing Skills

```bash
gh skill list                    # show all installed skills
gh skill info <name>             # show metadata for a skill
gh skill update <name>           # pull latest version
gh skill update --all            # update all installed skills
gh skill remove <name>           # uninstall + remove symlinks
```

## Workflow

1. **Search** for what you need: `gh skill search <query>`
2. **Review** the results — check description and author
3. **Install** the best match: `gh skill add <id>`
4. **Confirm** it's linked: `gh skill list`
5. **Use it** — the skill is now in your available skills

## Notes

- Skills are backed by GitHub Gists and discovered via the `[gh-skill]` prefix convention
- Installed skills live in `~/.gistskills/` and are symlinked to your tool's skill directory
- Trust gate: you'll be prompted before installing from untrusted authors
- Always tell the user what you found and let them decide before installing
