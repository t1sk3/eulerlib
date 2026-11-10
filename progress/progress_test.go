package progress

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/schollz/progressbar/v3"
)

func TestProgressBarYieldsIndexAndElement(t *testing.T) {
	s := []string{"a", "b", "c"}
	var buf bytes.Buffer

	var gotIdx []int
	var gotVal []string
	for i, e := range ProgressBar(s, WithWriter(&buf)) {
		gotIdx = append(gotIdx, i)
		gotVal = append(gotVal, e)
	}

	if len(gotIdx) != len(s) {
		t.Fatalf("yielded %d items, want %d", len(gotIdx), len(s))
	}
	for i, e := range s {
		if gotIdx[i] != i || gotVal[i] != e {
			t.Errorf("item %d = (%d, %q), want (%d, %q)", i, gotIdx[i], gotVal[i], i, e)
		}
	}
	if buf.Len() == 0 {
		t.Error("expected some output to be written to the writer")
	}
}

func TestProgressBarStopsOnBreak(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	var buf bytes.Buffer

	var got []int
	for i, e := range ProgressBar(s, WithWriter(&buf), progressbar.OptionThrottle(0)) {
		got = append(got, e)
		if i == 1 {
			break
		}
	}

	if want := []int{1, 2}; len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if !strings.Contains(buf.String(), "2/5") {
		t.Fatalf("expected final count 2/5 after break, got %q", buf.String())
	}
}

func TestProgressBarN(t *testing.T) {
	var buf bytes.Buffer

	var got []int
	for i := range ProgressBarN(4, WithWriter(&buf)) {
		got = append(got, i)
	}

	want := []int{0, 1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestProgressBarNEmpty(t *testing.T) {
	var buf bytes.Buffer

	count := 0
	for range ProgressBarN(0, WithWriter(&buf)) {
		count++
	}

	if count != 0 {
		t.Errorf("expected 0 iterations, got %d", count)
	}
	if buf.Len() == 0 {
		t.Error("expected the bar to still finalize output for an empty range")
	}
}

func TestSpinnerStopsOnBreak(t *testing.T) {
	var buf bytes.Buffer

	var got []int
	for i := range Spinner(WithWriter(&buf)) {
		got = append(got, i)
		if i == 4 {
			break
		}
	}

	want := []int{0, 1, 2, 3, 4}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
	if buf.Len() == 0 {
		t.Error("expected some output to be written to the writer")
	}
}

func TestWithLabelAndWidth(t *testing.T) {
	var buf bytes.Buffer

	for range ProgressBarN(1, WithWriter(&buf), WithLabel("euler"), WithWidth(10)) {
	}

	out := buf.String()
	if !strings.Contains(out, "euler ") {
		t.Errorf("expected output to contain label, got %q", out)
	}
	if !strings.Contains(out, "1/1") {
		t.Errorf("expected output to contain final count, got %q", out)
	}
}

func TestWithWriterDoesNotWriteToStderr(t *testing.T) {
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	var buf bytes.Buffer
	for range ProgressBarN(1, WithWriter(&buf)) {
	}

	if err := w.Close(); err != nil {
		t.Fatalf("stderr pipe close error = %v", err)
	}

	stderrOut, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stderr read error = %v", err)
	}
	if len(stderrOut) != 0 {
		t.Fatalf("expected no stderr output with WithWriter, got %q", string(stderrOut))
	}
}

func TestBarsAreNotDrawnBeforeIteration(t *testing.T) {
	for name, mk := range map[string]func(w io.Writer){
		"ProgressBar":  func(w io.Writer) { ProgressBar([]int{1, 2, 3}, WithWriter(w)) },
		"ProgressBarN": func(w io.Writer) { ProgressBarN(3, WithWriter(w)) },
		"Spinner":      func(w io.Writer) { Spinner(WithWriter(w)) },
	} {
		var buf bytes.Buffer
		mk(&buf)
		if buf.Len() != 0 {
			t.Errorf("%s wrote %q before iteration started", name, buf.String())
		}
	}
}

func TestRangingTwiceDrawsTwoCompleteBars(t *testing.T) {
	var buf bytes.Buffer
	seq := ProgressBarN(3, WithWriter(&buf), progressbar.OptionThrottle(0))
	for _, pass := range []int{1, 2} {
		buf.Reset()
		n := 0
		for range seq {
			n++
		}
		if n != 3 {
			t.Fatalf("pass %d: got %d iterations, want 3", pass, n)
		}
		if out := buf.String(); !strings.Contains(out, "3/3") || !strings.Contains(out, "0/3") {
			t.Errorf("pass %d: expected a full bar from 0/3 to 3/3, got %q", pass, out)
		}
	}

	s := ProgressBar([]string{"a", "b"}, WithWriter(&buf), progressbar.OptionThrottle(0))
	for _, pass := range []int{1, 2} {
		buf.Reset()
		for range s {
		}
		if out := buf.String(); !strings.Contains(out, "2/2") || !strings.Contains(out, "0/2") {
			t.Errorf("ProgressBar pass %d: expected a full bar from 0/2 to 2/2, got %q", pass, out)
		}
	}
}
