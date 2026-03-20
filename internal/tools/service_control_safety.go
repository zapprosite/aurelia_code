package tools

// blocksSelfServiceMutation prevents the bot from mutating its own
// service (restart/stop/start) which would kill the running process.
// Read-only operations (status, logs) are always allowed.
func blocksSelfServiceMutation(action serviceAction, service string) string {
	if service != "aurelia.service" {
		return ""
	}
	switch action {
	case serviceActionStatus, serviceActionLogs, serviceActionList:
		return ""
	default:
		return "self-mutation blocked: cannot " + string(action) + " own service"
	}
}
