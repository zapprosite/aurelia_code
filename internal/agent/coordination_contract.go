package agent

import "strings"

type CoordinationMode string

const (
	CoordinationDelegation CoordinationMode = "delegation"
	CoordinationHandoff    CoordinationMode = "handoff"
	CoordinationAssist     CoordinationMode = "assist"
)

func DefaultCoordinationModes() []string {
	return []string{
		string(CoordinationDelegation),
		string(CoordinationHandoff),
		string(CoordinationAssist),
	}
}

func CoordinationLabel(modes []string) string {
	if len(modes) == 0 {
		modes = DefaultCoordinationModes()
	}
	return strings.Join(modes, " + ")
}
