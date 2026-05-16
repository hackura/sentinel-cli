package auth

import (
	"testing"
)

func TestGenerateDeviceID(t *testing.T) {
	id1, err := GenerateDeviceID()
	if err != nil {
		t.Fatalf("failed to generate device ID: %v", err)
	}

	if len(id1) != 32 {
		t.Errorf("expected 32 characters, got %d", len(id1))
	}

	id2, _ := GenerateDeviceID()
	if id1 == id2 {
		t.Errorf("generated IDs are not unique")
	}
}

func TestMaskDeviceID(t *testing.T) {
	id := "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
	expected := "a1b2••••••o5p6"
	masked := MaskDeviceID(id)

	if masked != expected {
		t.Errorf("expected %s, got %s", expected, masked)
	}
}
