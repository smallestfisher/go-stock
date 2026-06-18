package data

import "testing"

func TestEffectiveSponsorVipLevelIsFreeAndActive(t *testing.T) {
	SponsorDecryptKeyHex = ""

	level, active := EffectiveSponsorVipLevel()

	if level != 2 {
		t.Fatalf("level = %d, want 2", level)
	}
	if !active {
		t.Fatal("active = false, want true")
	}
}
