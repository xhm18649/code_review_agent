package agent

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"code-review-agent/internal/config"
	"code-review-agent/internal/prompt"
	"code-review-agent/internal/tools"
)

func TestAuditStageAssignsAndLoadsLanguageSpecialists(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"index.php", "Api.cs"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	registry, err := tools.NewRegistry(root, 4096)
	if err != nil {
		t.Fatal(err)
	}
	prompts := prompt.Prompts{Skills: []prompt.Skill{
		{Name: "dotnet-security", Content: "dotnet skill"},
		{Name: "php-security", Content: "php skill"},
	}}
	team := NewTeam(config.Config{Agent: config.AgentConfig{ReconAgents: 1, AuditAgents: 2}}, prompts, nil, nil, registry)
	team.createStageLocked(phaseAudit)

	if got := team.securityAssignments(phaseAudit, 2); !reflect.DeepEqual(got, []string{"dotnet-security", "php-security"}) {
		t.Fatalf("unexpected specialist assignments: %v", got)
	}
	for i, worker := range team.workers {
		want := []string{"dotnet-security", "php-security"}[i]
		if got := worker.agent.prompts.LoadedSkillNames(); !reflect.DeepEqual(got, []string{want}) {
			t.Fatalf("worker %d loaded %v, want %v", i, got, want)
		}
		if !containsText(worker.agent.assignment, want) {
			t.Fatalf("worker %d assignment omitted %s: %q", i, want, worker.agent.assignment)
		}
		for _, action := range []string{"review_state", "search_content", "file_review_update", "todo_create", "verify_finding"} {
			if !containsText(worker.agent.assignment, action) {
				t.Fatalf("worker %d assignment omitted required action %s", i, action)
			}
		}
	}
}

func containsText(value, want string) bool {
	return strings.Contains(value, want)
}
