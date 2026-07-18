package placement

import (
	"context"
	"errors"
	"sort"
)

type Logger interface {
	Log(ctx context.Context, msg string, args ...any)
}

type Engine struct {
	scorer  Scorer
	checker *ConstraintChecker
	logger  Logger
}

func NewEngine(scorer Scorer, checker *ConstraintChecker) *Engine {
	return &Engine{
		scorer:  scorer,
		checker: checker,
	}
}

func (e *Engine) Scorer() Scorer {
	return e.scorer
}

func (e *Engine) Place(ctx context.Context, candidates []Candidate, req WorkloadRequest) (ScoreResult, error) {
	filtered, _ := e.checker.FilterByConstraints(candidates, req.Constraints, ConstraintContext{})
	if len(filtered) == 0 {
		return ScoreResult{}, errors.New("no viable candidates after constraint filtering")
	}

	var results []ScoreResult
	for _, c := range filtered {
		score, reasons, err := e.scorer.Score(ctx, c, req)
		if err != nil {
			continue
		}
		results = append(results, ScoreResult{NodeID: c.NodeID, Score: score, Reasons: reasons})
	}

	if len(results) == 0 {
		return ScoreResult{}, errors.New("no viable candidates after scoring")
	}

	best := results[0]
	for _, r := range results[1:] {
		if r.Score > best.Score {
			best = r
		}
	}

	if e.logger != nil {
		e.logger.Log(ctx, "placement selected", "node", best.NodeID, "score", best.Score)
	}

	return best, nil
}

func (e *Engine) PlaceAll(ctx context.Context, candidates []Candidate, req WorkloadRequest) ([]ScoreResult, error) {
	filtered, _ := e.checker.FilterByConstraints(candidates, req.Constraints, ConstraintContext{})
	if len(filtered) == 0 {
		return nil, errors.New("no viable candidates after constraint filtering")
	}

	var results []ScoreResult
	for _, c := range filtered {
		score, reasons, err := e.scorer.Score(ctx, c, req)
		if err != nil {
			continue
		}
		results = append(results, ScoreResult{NodeID: c.NodeID, Score: score, Reasons: reasons})
	}

	if len(results) == 0 {
		return nil, errors.New("no viable candidates after scoring")
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}
