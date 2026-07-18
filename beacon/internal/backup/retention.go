package backup

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type RetentionPolicy struct {
	MaxBackups  int
	MaxAge      time.Duration
	KeepDaily   int
	KeepWeekly  int
	KeepMonthly int
}

func (p RetentionPolicy) Apply(ctx context.Context, store Store, serverID string) error {
	backups, err := store.List(ctx, serverID, 0)
	if err != nil {
		return err
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CompletedAt.After(backups[j].CompletedAt)
	})

	keep := make(map[string]bool, len(backups))
	for _, b := range backups {
		keep[b.ID] = true
	}

	cutoff := time.Now()
	if p.MaxAge > 0 {
		cutoff = cutoff.Add(-p.MaxAge)
		for _, b := range backups {
			if b.CompletedAt.Before(cutoff) {
				delete(keep, b.ID)
			}
		}
	}

	if p.KeepDaily > 0 {
		keepTime := func(t time.Time) string { return t.Format("2006-01-02") }
		applyKeepPerPeriod(backups, keep, p.KeepDaily, keepTime)
	}
	if p.KeepWeekly > 0 {
		keepTime := func(t time.Time) string {
			y, w := t.ISOWeek()
			return fmt.Sprintf("%d-%d", y, w)
		}
		applyKeepPerPeriod(backups, keep, p.KeepWeekly, keepTime)
	}
	if p.KeepMonthly > 0 {
		keepTime := func(t time.Time) string { return t.Format("2006-01") }
		applyKeepPerPeriod(backups, keep, p.KeepMonthly, keepTime)
	}

	if p.MaxBackups > 0 {
		kept := 0
		for _, b := range backups {
			if keep[b.ID] {
				kept++
			}
		}
		if kept > p.MaxBackups {
			for i := len(backups) - 1; i >= 0 && kept > p.MaxBackups; i-- {
				if keep[backups[i].ID] {
					delete(keep, backups[i].ID)
					kept--
				}
			}
		}
	}

	for _, b := range backups {
		if !keep[b.ID] {
			if err := store.Delete(ctx, b.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func applyKeepPerPeriod(backups []Backup, keep map[string]bool, maxPerPeriod int, period func(time.Time) string) {
	periodCount := make(map[string]int)
	for _, b := range backups {
		if !keep[b.ID] {
			continue
		}
		periodCount[period(b.CompletedAt)]++
	}
	for i := len(backups) - 1; i >= 0; i-- {
		if !keep[backups[i].ID] {
			continue
		}
		p := period(backups[i].CompletedAt)
		if periodCount[p] > maxPerPeriod {
			delete(keep, backups[i].ID)
			periodCount[p]--
		}
	}
}
