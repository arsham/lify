package ui

// An Asset is a sprite on the UI.
type Asset uint8

const (
	// AssetHerb1 is the herb version 1.
	AssetHerb1 Asset = iota + 1
	// AssetFontInfo is a large font for showing information.
	AssetFontInfo
)

func (a Asset) String() string {
	if a == AssetHerb1 {
		return "herb/herb1"
	}
	return ""
}
