package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/telegram"
)

// squadStatusAdapter implements telegram.SquadStatusReporter using agent.GetFixedSquad().
type squadStatusAdapter struct{}

func (squadStatusAdapter) GetSquadStatus() []telegram.AgentStatus {
	members := agent.GetFixedSquad()
	out := make([]telegram.AgentStatus, 0, len(members))
	for _, m := range members {
		out = append(out, telegram.AgentStatus{
			Name:   m.Name,
			Icon:   m.IconName,
			Role:   m.Role,
			Status: m.Status,
			Load:   m.Load,
		})
	}
	return out
}

// cronNextJobAdapter implements telegram.CronNextJobReporter using the cron store.
type cronNextJobAdapter struct {
	store *cron.SQLiteCronStore
}

func (a *cronNextJobAdapter) GetNextJobs(ctx context.Context, limit int) []telegram.NextJob {
	if a.store == nil {
		return nil
	}
	now := time.Now().UTC()
	// List due jobs in the next 24 hours
	jobs, err := a.store.ListDueJobs(ctx, now.Add(24*time.Hour), limit*5)
	if err != nil || len(jobs) == 0 {
		return nil
	}

	var result []telegram.NextJob
	seen := 0
	for _, j := range jobs {
		if !j.Active || j.NextRunAt == nil {
			continue
		}
		dur := j.NextRunAt.Sub(now)
		if dur < 0 {
			dur = 0
		}
		name := j.ID
		if len(name) > 16 {
			name = name[:16]
		}
		result = append(result, telegram.NextJob{
			Name:   name,
			NextIn: formatDuration(dur),
		})
		seen++
		if seen >= limit {
			break
		}
	}
	return result
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dmin", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh%02dmin", int(d.Hours()), int(d.Minutes())%60)
}
