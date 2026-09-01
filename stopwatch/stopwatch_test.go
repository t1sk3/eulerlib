package stopwatch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestStartPrintsLabelAndReturnsElapsed(t *testing.T) {
	var buf bytes.Buffer

	stop := Start("brute force", WithWriter(&buf))
	time.Sleep(5 * time.Millisecond)
	elapsed := stop()

	if elapsed < 5*time.Millisecond {
		t.Errorf("elapsed = %s, want >= 5ms", elapsed)
	}
	out := buf.String()
	if !strings.Contains(out, "brute force:") {
		t.Errorf("expected output to contain the label, got %q", out)
	}
}

func TestTimeRunsFnAndReturnsElapsed(t *testing.T) {
	var buf bytes.Buffer

	ran := false
	elapsed := Time("part 1", func() {
		ran = true
		time.Sleep(5 * time.Millisecond)
	}, WithWriter(&buf))

	if !ran {
		t.Error("expected fn to be called")
	}
	if elapsed < 5*time.Millisecond {
		t.Errorf("elapsed = %s, want >= 5ms", elapsed)
	}
	if !strings.Contains(buf.String(), "part 1:") {
		t.Errorf("expected output to contain the label, got %q", buf.String())
	}
}
