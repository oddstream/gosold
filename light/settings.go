package light

// Settings holds user preferences.
// Colors are named from the web extended colors at https://en.wikipedia.org/wiki/Web_colors
type Settings struct {
	// Capitals to emit to json
	Variant                            string
	BaizeColor                         string
	CardFaceColor                      string
	CardBackColor                      string
	MovableCardBackColor               string
	BlackColor                         string
	RedColor                           string
	ClubColor                          string
	DiamondColor                       string
	HeartColor                         string
	SpadeColor                         string
	ColorfulCards                      bool
	PowerMoves                         bool
	SafeCollect, AutoCollect           bool
	Mute                               bool
	Volume                             float64
	MirrorBaize                        bool
	ShowMovableCards                   bool
	AlwaysShowMovableCards             bool
	Timer                              bool
	CardRatio                          float64
	AniSpeed                           float64
	LastVersionMajor, LastVersionMinor int
	// FixedCards                         bool
	// FixedCardWidth, FixedCardHeight    int
}

func NewSettings() *Settings {
	s := &Settings{
		Variant:                "Klondike",
		BaizeColor:             "BaizeGreen",
		PowerMoves:             true,
		SafeCollect:            false,
		AutoCollect:            false,
		CardFaceColor:          "Ivory",
		CardBackColor:          "CornflowerBlue",
		MovableCardBackColor:   "Gold",
		ColorfulCards:          true,
		RedColor:               "Crimson",
		BlackColor:             "Black",
		ClubColor:              "DarkGreen",
		DiamondColor:           "DarkBlue",
		HeartColor:             "Crimson",
		SpadeColor:             "Black",
		Mute:                   false,
		Volume:                 0.75,
		ShowMovableCards:       true,
		AlwaysShowMovableCards: true,
		Timer:                  false,
		// FixedCards:       false,
		// FixedCardWidth:   90,
		// FixedCardHeight:  122,
		CardRatio:        1.39, // official poker size
		AniSpeed:         0.6,  // Normal
		LastVersionMajor: 0,
		LastVersionMinor: 0,
	}
	return s
}
