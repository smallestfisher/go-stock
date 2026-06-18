package data

// DefaultSponsorAESKeyHex is kept for compatibility with older builds and configs.
const DefaultSponsorAESKeyHex = ""

// SponsorDecryptKeyHex is kept for compatibility with older startup code.
var SponsorDecryptKeyHex string

// EffectiveSponsorVipLevel reports the free entitlement used by all local features.
func EffectiveSponsorVipLevel() (level int, active bool) {
	return 2, true
}
