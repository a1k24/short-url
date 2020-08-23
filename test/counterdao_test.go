package test

import (
	"github.com/a1k24/short-url/internal/app"
	"testing"
)

func TestCounterDao(t *testing.T) {
	name := "dummy_sequence"

	sequence, err := app.GenerateNextSequence(name)
	if nil != err {
		t.Error("Failed to generate sequence.")
	}
	if sequence <= 0 {
		t.Errorf("Invalid sequence: %d", sequence)
	}

	defer deleteSequence(t, name)

	nextSequence, err := app.GenerateNextSequence(name)
	if nil != err {
		t.Error("Failed to generate sequence.")
	}
	if g, w := nextSequence, sequence+1; g != w {
		t.Errorf("sequence = %v, want %v", g, w)
	}
}

func deleteSequence(t *testing.T, name string) {
	deleteResult, err := app.DropSequence(name)
	if nil != err {
		t.Error("Failed to drop sequence.")
	}
	if g, w := deleteResult.DeletedCount, int64(1); g != w {
		t.Errorf("linkhash: %s, delete_count = %v, want %v", info.LinkHash, g, w)
	}
}
