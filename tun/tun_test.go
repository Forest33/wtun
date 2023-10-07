package tun

import (
	"testing"
)

func TestCreate(t *testing.T) {
	d, err := CreateTUN("", 1500)
	if err != nil {
		t.Errorf("failed create device: %v", err)
	}
	_ = d
}
