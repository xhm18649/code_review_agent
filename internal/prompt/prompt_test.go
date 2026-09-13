package prompt

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadSkillsDiscoversAndSortsSkills(t *testing.T) {
	root := t.TempDir()
	for name, content := range map[string]string{
		"dotnet-security": "dotnet",
		"php-security":    "php",
		"evidence-ledger": "evidence",
	} {
		if err := os.Mkdir(filepath.Join(root, name), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, name, "SKILL.md"), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(root, "without-skill"), 0700); err != nil {
		t.Fatal(err)
	}

	got, err := loadSkills(root)
	if err != nil {
		t.Fatal(err)
	}
	if names := []string{got[0].Name, got[1].Name, got[2].Name}; !reflect.DeepEqual(names, []string{"dotnet-security", "evidence-ledger", "php-security"}) {
		t.Fatalf("skills were not sorted or discovered correctly: %v", names)
	}
}
