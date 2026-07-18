package evacuationplanner

import (
	"testing"

	"gamepanel/forge/internal/store"
)

func TestEvacuationMigrationOutcomeOnlyCompletesTerminalMigrations(t *testing.T) {
	tests := map[string]string{
		string(store.MigrationStatusPlanned):      "pending",
		string(store.MigrationStatusPreparing):    "pending",
		string(store.MigrationStatusTransferring): "pending",
		string(store.MigrationStatusRestoring):    "pending",
		string(store.MigrationStatusCompleted):    "completed",
		string(store.MigrationStatusFailed):       "failed",
		string(store.MigrationStatusCancelled):    "failed",
	}
	for status, want := range tests {
		t.Run(status, func(t *testing.T) {
			if got := evacuationMigrationOutcome(status); got != want {
				t.Fatalf("evacuationMigrationOutcome(%q) = %q, want %q", status, got, want)
			}
		})
	}
}

func TestEvacuationPlanFinishedRequiresEveryItemToBeTerminal(t *testing.T) {
	tests := []struct {
		name  string
		items []store.EvacuationItem
		want  bool
	}{
		{name: "empty", want: true},
		{name: "completed and failed", items: []store.EvacuationItem{{Status: "completed"}, {Status: "failed"}}, want: true},
		{name: "pending", items: []store.EvacuationItem{{Status: "pending"}}, want: false},
		{name: "preparing", items: []store.EvacuationItem{{Status: "preparing"}}, want: false},
		{name: "running", items: []store.EvacuationItem{{Status: "running"}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evacuationPlanFinished(tt.items); got != tt.want {
				t.Fatalf("evacuationPlanFinished() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestEvacuationPlanStatusReflectsPlanningOnly(t *testing.T) {
	tests := []struct {
		name  string
		items []PlanItem
		want  store.EvacuationPlanStatus
	}{
		{
			name:  "eligible targets remain pending execution",
			items: []PlanItem{{EvacuationItem: store.EvacuationItem{Eligible: true}}},
			want:  store.EvacuationPlanStatusPending,
		},
		{
			name: "empty plan remains pending execution",
			want: store.EvacuationPlanStatusPending,
		},
		{
			name:  "ineligible target fails planning",
			items: []PlanItem{{EvacuationItem: store.EvacuationItem{Eligible: false}}},
			want:  store.EvacuationPlanStatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evacuationPlanStatus(tt.items); got != tt.want {
				t.Fatalf("evacuationPlanStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}
