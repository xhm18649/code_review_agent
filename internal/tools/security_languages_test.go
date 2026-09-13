package tools

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSecurityLanguagesDetectsPHPAndDotNetFiles(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"index.php", "Admin.cshtml", "Api.cs", "worker.csproj", "README.md"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	registry, err := NewRegistry(root, 4096)
	if err != nil {
		t.Fatal(err)
	}
	if got := registry.SecurityLanguages(); !reflect.DeepEqual(got, []string{"dotnet-security", "php-security"}) {
		t.Fatalf("unexpected security languages: %v", got)
	}
}

func TestSecurityLanguagesIgnoresUnrelatedFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	registry, err := NewRegistry(root, 4096)
	if err != nil {
		t.Fatal(err)
	}
	if got := registry.SecurityLanguages(); len(got) != 0 {
		t.Fatalf("unrelated files selected security skills: %v", got)
	}
}

func TestSecurityLanguagesDetectsProjectAndDeploymentMarkers(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"composer.json", "web.config", "App.sln", "settings.json"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	registry, err := NewRegistry(root, 4096)
	if err != nil {
		t.Fatal(err)
	}
	if got := registry.SecurityLanguages(); !reflect.DeepEqual(got, []string{"dotnet-security", "php-security"}) {
		t.Fatalf("project/deployment markers were not detected: %v", got)
	}
}
