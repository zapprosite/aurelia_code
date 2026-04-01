package config

import (
	"os"
	"strings"
)

type Features struct {
	VoiceEnabled       bool
	DreamEnabled       bool
	KAIROSEnabled      bool
	ComputerUseEnabled bool
}

func LoadFeatures() Features {
	return Features{
		VoiceEnabled:       isEnvTrue("FEAT_VOICE"),
		DreamEnabled:       isEnvTrue("FEAT_DREAM"),
		KAIROSEnabled:      isEnvTrue("FEAT_KAIROS"),
		ComputerUseEnabled: isEnvTrue("FEAT_COMPUTER_USE"),
	}
}

func isEnvTrue(key string) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return val == "true" || val == "1" || val == "yes" || val == "on"
}
