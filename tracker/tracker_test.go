package tracker

import (
	"testing"
)

func TestTrack(t *testing.T) {
	input := "RV170984370CN"
	t.Logf("input=%s", input)

	result, err := NewTracker().Track(input)
	if err != nil {
		t.Errorf("Error returned %s", err)
	}

	t.Run("Testing return length", func(t *testing.T) {
		if len(result) != 11 {
			t.Errorf("Should return 11 lines, actual: %d", len(result))
		}
	})

	t.Run("Testing tracking number column", func(t *testing.T) {
		for _, elm := range result {
			if elm.TrackingNumber != input {
				t.Errorf("Should return %s in tracking number column, actual: %s", input, elm.TrackingNumber)
			}
		}
	})

}
