package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

// SquadStatusReporter provides team agent status for the /status command.
type SquadStatusReporter interface {
	GetSquadStatus() []AgentStatus
}

// AgentStatus is a snapshot of one team member's status.
type AgentStatus struct {
	Name   string
	Icon   string
	Role   string
	Status string // "online", "offline", "busy"
	Load   int
}

// CronNextJobReporter provides upcoming cron job names for the /status command.
type CronNextJobReporter interface {
	GetNextJobs(ctx context.Context, limit int) []NextJob
}

// NextJob represents an upcoming scheduled job.
type NextJob struct {
	Name   string
	NextIn string
}

// startedAt is the process start time, used for uptime calculation.
var startedAt = time.Now()

// SetSquadReporter wires the team reporter.
func (bc *BotController) SetSquadReporter(r SquadStatusReporter) {
	bc.squadReporter = r
}

// SetCronJobReporter wires the cron job reporter.
func (bc *BotController) SetCronJobReporter(r CronNextJobReporter) {
	bc.cronJobReporter = r
}

// handleSquadStatus builds and sends a formatted team + cron status message.
func (bc *BotController) handleSquadStatus(c telebot.Context) error {
	var sb strings.Builder

	// Team
	if bc.squadReporter != nil {
		agents := bc.squadReporter.GetSquadStatus()
		online := 0
		for _, a := range agents {
			if a.Status == "online" || a.Status == "busy" {
				online++
			}
		}
		sb.WriteString(fmt.Sprintf("🟢 Team Online (%d/%d)\n", online, len(agents)))
		for i, a := range agents {
			prefix := "├─"
			if i == len(agents)-1 {
				prefix = "└─"
			}
			statusIcon := "🔴"
			if a.Status == "online" {
				statusIcon = "🟢"
			} else if a.Status == "busy" {
				statusIcon = "🟡"
			}
			sb.WriteString(fmt.Sprintf("%s %s %s — %s %d%%\n", prefix, statusIcon, a.Name, a.Status, a.Load))
		}
	} else {
		sb.WriteString("Team: indisponível\n")
	}

	// Crons
	if bc.cronJobReporter != nil {
		jobs := bc.cronJobReporter.GetNextJobs(context.Background(), 3)
		if len(jobs) > 0 {
			sb.WriteString("\n⏰ Próximos Crons\n")
			for i, j := range jobs {
				prefix := "├─"
				if i == len(jobs)-1 {
					prefix = "└─"
				}
				sb.WriteString(fmt.Sprintf("%s %s → em %s\n", prefix, j.Name, j.NextIn))
			}
		}
	}

	// Uptime
	uptime := time.Since(startedAt).Round(time.Second)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	sb.WriteString(fmt.Sprintf("\n⚙️ Uptime: %dh%02dmin\n", hours, minutes))

	return SendContextText(c, sb.String())
}
