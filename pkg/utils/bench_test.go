package utils

import (
	"testing"
)

func TestCaptureMemSnapshot_NotZero(t *testing.T) {
	snap := CaptureMemSnapshot()
	if snap.SysMB == 0 {
		t.Error("expected non-zero SysMB")
	}
	if snap.Goroutines == 0 {
		t.Error("expected non-zero goroutines")
	}
}

func TestPrintMemDiff_NoPanic(t *testing.T) {
	before := CaptureMemSnapshot()
	after := CaptureMemSnapshot()
	PrintMemDiff(before, after)
}
