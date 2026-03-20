package agent

import "testing"

func TestBrowserLoginPolicyValidateAllowedHost(t *testing.T) {
	policy := BrowserLoginPolicy{
		AllowedHosts: []string{"example.com"},
		MaxSteps:     5,
	}
	if err := policy.Validate(BrowserLoginRequest{
		TargetURL: "https://app.example.com/login",
		Stage:     BrowserLoginStageUsername,
		StepCount: 1,
	}); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestBrowserLoginPolicyRejectsHostOutsideAllowlist(t *testing.T) {
	policy := BrowserLoginPolicy{
		AllowedHosts: []string{"example.com"},
		MaxSteps:     5,
	}
	if err := policy.Validate(BrowserLoginRequest{
		TargetURL: "https://evil.com/login",
		Stage:     BrowserLoginStageUsername,
		StepCount: 1,
	}); err == nil {
		t.Fatal("expected error for disallowed host")
	}
}

func TestBrowserLoginPolicyNeedsHumanGate(t *testing.T) {
	policy := DefaultBrowserLoginPolicy()
	if !policy.NeedsHumanGate(BrowserLoginRequest{Stage: BrowserLoginStagePassword}) {
		t.Fatal("expected password stage to require human gate")
	}
	if !policy.NeedsHumanGate(BrowserLoginRequest{Stage: BrowserLoginStageTwoFA}) {
		t.Fatal("expected two-factor stage to require human gate")
	}
	if policy.NeedsHumanGate(BrowserLoginRequest{Stage: BrowserLoginStageUsername}) {
		t.Fatal("did not expect username stage to require human gate")
	}
}

func TestBrowserLoginPolicyStepBudget(t *testing.T) {
	policy := BrowserLoginPolicy{MaxSteps: 2}
	if err := policy.Validate(BrowserLoginRequest{
		TargetURL: "https://example.com/login",
		Stage:     BrowserLoginStageStart,
		StepCount: 3,
	}); err == nil {
		t.Fatal("expected step budget error")
	}
}
