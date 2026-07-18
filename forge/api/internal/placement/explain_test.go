package placement

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExplainPlacement_ProducesReport(t *testing.T) {
	candidates := []Candidate{
		{NodeID: "node-1", RegionID: "us-east", AvailableMemory: 100},
		{NodeID: "node-2", RegionID: "us-west", AvailableMemory: 200},
		{NodeID: "node-3", RegionID: "us-east", AvailableMemory: 300},
	}

	constraints := []Constraint{
		{Type: ConstraintRegion, Required: true, Values: []string{"us-east"}},
	}

	engine := NewEngine(&LeastLoadedScorer{}, NewConstraintChecker())
	report, err := ExplainPlacement(context.Background(), engine, candidates, WorkloadRequest{Constraints: constraints})
	require.NoError(t, err)

	assert.Equal(t, 3, report.TotalCandidates)
	assert.Len(t, report.FilteredOut, 1)
	assert.Equal(t, "node-2", report.FilteredOut[0].NodeID)
	assert.Len(t, report.ScoredCandidates, 2)
	require.NotNil(t, report.Selected)
	assert.Equal(t, "node-3", report.Selected.NodeID)
}

func TestExplainScores_ReturnsBreakdowns(t *testing.T) {
	candidates := []Candidate{
		{NodeID: "node-1", AvailableMemory: 100, AvailableCPU: 10, AvailableDisk: 1000},
		{NodeID: "node-2", AvailableMemory: 200, AvailableCPU: 20, AvailableDisk: 2000},
	}

	breakdowns := ExplainScores(context.Background(), &LeastLoadedScorer{}, candidates, WorkloadRequest{})
	require.Len(t, breakdowns, 2)
	assert.Equal(t, "node-1", breakdowns[0].NodeID)
	assert.True(t, breakdowns[0].BaseScore > 0)
	assert.Equal(t, "node-2", breakdowns[1].NodeID)
	assert.True(t, breakdowns[1].BaseScore > breakdowns[0].BaseScore)
}

func TestExplainPlacement_NoViableCandidates(t *testing.T) {
	candidates := []Candidate{
		{NodeID: "node-1", RegionID: "us-east"},
		{NodeID: "node-2", RegionID: "us-west"},
	}

	constraints := []Constraint{
		{Type: ConstraintRegion, Required: true, Values: []string{"eu-west"}},
	}

	engine := NewEngine(&LeastLoadedScorer{}, NewConstraintChecker())
	report, err := ExplainPlacement(context.Background(), engine, candidates, WorkloadRequest{Constraints: constraints})
	require.NoError(t, err)

	assert.Equal(t, 2, report.TotalCandidates)
	assert.Len(t, report.FilteredOut, 2)
	assert.Len(t, report.ScoredCandidates, 0)
	assert.Nil(t, report.Selected)
}
