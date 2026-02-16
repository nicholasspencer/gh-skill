package internal

import "testing"

func TestEffectiveProvider(t *testing.T) {
	tests := []struct {
		provider string
		want     string
	}{
		{"", "github"},
		{"github", "github"},
		{"gitlab", "gitlab"},
	}
	for _, tt := range tests {
		m := SkillMeta{Provider: tt.provider}
		if got := m.EffectiveProvider(); got != tt.want {
			t.Errorf("EffectiveProvider() with %q = %q, want %q", tt.provider, got, tt.want)
		}
	}
}

func TestExpandFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"SKILL.md", "SKILL.md"},
		{"scripts--setup.sh", "scripts/setup.sh"},
		{"references--api--docs.md", "references/api/docs.md"},
	}
	for _, tt := range tests {
		if got := ExpandFilename(tt.input); got != tt.want {
			t.Errorf("ExpandFilename(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFlattenFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"SKILL.md", "SKILL.md"},
		{"scripts/setup.sh", "scripts--setup.sh"},
	}
	for _, tt := range tests {
		if got := FlattenFilename(tt.input); got != tt.want {
			t.Errorf("FlattenFilename(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsSkillFile(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"weather.skill.md", true},
		{"WEATHER.SKILL.MD", true},
		{"README.md", false},
		{"SKILL.md", false},
	}
	for _, tt := range tests {
		if got := IsSkillFile(tt.input); got != tt.want {
			t.Errorf("IsSkillFile(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseFrontMatter(t *testing.T) {
	content := `---
name: my-skill
description: A test skill
version: 1.0.0
tags: [test, demo]
author: nico
---

# My Skill
Does things.
`
	fm, err := ParseFrontMatter(content)
	if err != nil {
		t.Fatalf("ParseFrontMatter() error: %v", err)
	}
	if fm.Name != "my-skill" {
		t.Errorf("Name = %q, want my-skill", fm.Name)
	}
	if fm.Version != "1.0.0" {
		t.Errorf("Version = %q, want 1.0.0", fm.Version)
	}
	if fm.Author != "nico" {
		t.Errorf("Author = %q, want nico", fm.Author)
	}
	if len(fm.Tags) != 2 {
		t.Errorf("Tags = %v, want [test demo]", fm.Tags)
	}
}

func TestParseFrontMatter_None(t *testing.T) {
	fm, err := ParseFrontMatter("# Just markdown\nNo front matter here.")
	if err != nil {
		t.Fatalf("ParseFrontMatter() error: %v", err)
	}
	if fm.Name != "" {
		t.Errorf("Expected empty name, got %q", fm.Name)
	}
}

func TestFindSkillFile(t *testing.T) {
	// *.skill.md preferred over SKILL.md
	files := map[string]GistFile{
		"weather.skill.md": {Filename: "weather.skill.md", Content: "test"},
		"SKILL.md":         {Filename: "SKILL.md", Content: "legacy"},
	}
	name, f, ok := FindSkillFile(files)
	if !ok {
		t.Fatal("FindSkillFile returned false")
	}
	if name != "weather.skill.md" {
		t.Errorf("name = %q, want weather.skill.md", name)
	}
	if f.Content != "test" {
		t.Errorf("content = %q, want test", f.Content)
	}

	// Fallback to SKILL.md
	files2 := map[string]GistFile{
		"SKILL.md": {Filename: "SKILL.md", Content: "legacy"},
	}
	name2, _, ok2 := FindSkillFile(files2)
	if !ok2 {
		t.Fatal("FindSkillFile fallback returned false")
	}
	if name2 != "SKILL.md" {
		t.Errorf("fallback name = %q, want SKILL.md", name2)
	}
}
