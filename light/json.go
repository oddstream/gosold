//go:build linux || windows || android

package light

import (
	"encoding/json"
	"log"

	"oddstream.games/gosold/util"
)

// Load an already existing Settings object from file
func (s *Settings) load() {
	// defer util.Duration(time.Now(), "Settings.Load")

	bytes, count, err := util.LoadBytesFromFile("settings.json", false)
	if err != nil || count == 0 || bytes == nil {
		return
	}

	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], s)
	if err != nil {
		log.Panic("Settings.Load Unmarshal", err)
	}
}

// Save writes the Settings object to file
func (s *Settings) save() {
	// defer util.Duration(time.Now(), "Settings.Save")

	s.LastVersionMajor = GosoldVersionMajor
	s.LastVersionMinor = GosoldVersionMinor
	// warning - calling ebiten function ouside RunGame loop will cause fatal panic
	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	util.SaveBytesToFile(bytes, "settings.json")
}
