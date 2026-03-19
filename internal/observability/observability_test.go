package observability

import (
	"reflect"
	"testing"
)

func TestRedact(t *testing.T) {
	if got := Redact(""); got != "" {
		t.Fatalf("Redact(empty) = %q, want empty string", got)
	}
	if got := Redact("secret"); got != "[redacted]" {
		t.Fatalf("Redact(secret) = %q", got)
	}
}

func TestBasename(t *testing.T) {
	if got := Basename("/tmp/aurelia/audio.wav"); got != "audio.wav" {
		t.Fatalf("Basename() = %q", got)
	}
}

func TestMapKeys(t *testing.T) {
	got := MapKeys(map[string]any{
		"beta":  2,
		"alpha": 1,
	})
	want := []string{"alpha", "beta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MapKeys() = %#v, want %#v", got, want)
	}
}
