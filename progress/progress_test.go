package progress

import (
	"bytes"
	"strings"
	"testing"
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
	for i, e := range ProgressBar(s, WithWriter(&buf)) {
		got = append(got, e)
		if i == 1 {
			break
		}
	}

	if want := []int{1, 2}; len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
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
