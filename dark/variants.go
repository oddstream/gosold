package dark

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

var variants = map[string]scripter{
	"Alhambra": &Alhambra{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Alhambra_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"Bisley": &Bisley{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Bisley_(card_game)",
			cardColors: 4,
		},
	},
	"Blockade": &Blockade{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Blockade_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"Colorado": &Colorado{
		scriptBase: scriptBase{
			packs: 2,
		},
	},
	"Eagle Wing": &EagleWing{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Eagle_Wing",
		},
	},
	"Klondike": &Klondike{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
		draw:     1,
		recycles: 2,
	},
	"Klondike Draw Three": &Klondike{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
		draw:     3,
		recycles: 2,
	},
	// "Thoughtful": &Klondike{
	// 	scriptBase: scriptBase{
	// 		wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
	// 	},
	// 	draw:       1,
	// 	recycles:   2,
	// 	thoughtful: true,
	// },
	"Gargantua": &Klondike{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Gargantua_(card_game)",
			packs:     2,
		},
		draw:     1,
		recycles: 2,
		founds:   []int{3, 4, 5, 6, 7, 8, 9, 10},    // 8
		tabs:     []int{2, 3, 4, 5, 6, 7, 8, 9, 10}, // 9
	},
	"Triple Klondike": &Klondike{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Gargantua_(card_game)",
			packs:     3,
		},
		draw:     1,
		recycles: 2,
		founds:   []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},             // 12
		tabs:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, // 16
	},
	"Freecell": &Freecell{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/FreeCell",
		},
	},
	"Forty Thieves": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
	},
	"Josephine": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
		moveType:    MOVE_ANY,
	},
	"Rank and File": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		proneRows:      []int{0, 1, 2},
		tabCompareFunc: dyad.compare_DownAltColor,
		moveType:       MOVE_ANY,
	},
	"Indian": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    3,
		proneRows:      []int{0},
		tabCompareFunc: dyad.compare_DownOtherSuit,
	},
	"Streets": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: dyad.compare_DownAltColor,
	},
	"Number Ten": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		proneRows:      []int{0, 1},
		tabCompareFunc: dyad.compare_DownAltColor,
		moveType:       MOVE_ANY,
	},
	"Limited": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 3,
	},
	"Forty and Eight": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 5,
		recycles:    1,
	},
	"Red and Black": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: dyad.compare_DownAltColor,
	},
	"Light and Shadow": &LightAndShadow{
		scriptBase: scriptBase{
			packs: 2,
		},
	},
	"Lucas": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{5, 6, 7, 8, 9, 10, 11, 12},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		cardsPerTab: 3,
		dealAces:    true,
	},
	"Busy Aces": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 1,
	},
	"Maria": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: dyad.compare_DownAltColor,
	},
	"Sixty Thieves": &FortyThieves{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      3,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		cardsPerTab: 5,
	},
	"Penguin": &Penguin{
		scriptBase: scriptBase{
			wikipedia:  "https://www.parlettgames.uk/patience/penguin.html",
			cardColors: 4,
		},
	},
	"Scorpion": &Scorpion{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Scorpion_(solitaire)",
			cardColors: 4,
		},
	},
	"Spider One Suit": &Spider{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 1,
			packs:      8,
			suits:      1,
		},
	},
	"Spider Two Suits": &Spider{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 2,
			packs:      4,
			suits:      2,
		},
	},
	"Spider Four Suits": &Spider{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 4,
			packs:      2,
			suits:      4,
		},
	},
	"Spiderette": &Spiderette{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)#Variants",
			cardColors: 4,
		},
	},
	"The Rainbow": &TheRainbow{
		scriptBase: scriptBase{},
	},
	"Classic Westcliff": &Westcliff{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Westcliff_(card_game)",
		},
		variant: "Classic",
	},
	"American Westcliff": &Westcliff{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Westcliff_(card_game)",
		},
		variant: "American",
	},
	// "Easthaven": &Westcliff{
	// 	scriptBase: scriptBase{
	// 		wikipedia: "https://en.wikipedia.org/wiki/Westcliff_(card_game)",
	// 	},
	// 	variant: "Easthaven",
	// },
	// "Whitehead": &Whitehead{
	// 	scriptBase: scriptBase{
	// 		wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
	// 	},
	// },
	"Uncle Sam": &UncleSam{
		scriptBase: scriptBase{
			cardColors: 2,
			packs:      2,
		},
	},
	"Yukon": &Yukon{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
		},
	},
	"Yukon Cells": &Yukon{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
		},
		extraCells: 2,
	},
}

var variantGroups = map[string][]string{
	// "> All" added dynamically by func init()
	// don't have any group that comes alphabetically before "> All"
	"> Canfields":     {"Eagle Wing"},
	"> Easier":        {"American Toad", "American Westcliff", "Blockade", "Classic Westcliff", "Lucas", "Spider One Suit", "Usk Relaxed"},
	"> Hapgood":       {"Light and Shadow", "The Rainbow", "Uncle Sam"},
	"> Harder":        {"Forty Thieves", "Spider Four Suits"},
	"> Forty Thieves": {"Forty Thieves", "Number Ten", "Red and Black", "Indian", "Rank and File", "Sixty Thieves", "Josephine", "Limited", "Forty and Eight", "Lucas", "Busy Aces", "Maria", "Streets"},
	"> Freecells":     {"Freecell"},
	"> Klondikes":     {"Gargantua", "Triple Klondike", "Klondike", "Klondike Draw Three"},
	"> People":        {"Duchess", "Josephine", "Maria", "Baker's Game"},
	"> Places":        {"Bisley", "Colorado", "Yukon", "Klondike"},
	"> Puzzlers":      {"Bisley", "Penguin"},
	"> Spiders":       {"Spider One Suit", "Spider Two Suits", "Spider Four Suits", "Scorpion", "Spiderette"},
	"> Yukons":        {"Yukon", "Yukon Cells"},
}

// init is used to assemble the "> All" alpha-sorted group of variants
func init() {
	// look in the scripts folder tree (depth one only) for *.lua files
	// turn subfolder names as group names
	{
		type scriptInfo struct {
			path, name, group string
		}

		var files []scriptInfo

		err := filepath.Walk("./scripts", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if !info.IsDir() && filepath.Ext(path) == ".lua" {
				var sinfo scriptInfo = scriptInfo{path: path, name: strings.TrimSuffix(filepath.Base(path), ".lua")}
				var splits []string
				if runtime.GOOS == "windows" {
					splits = strings.Split(path, "\\")
				} else {
					splits = strings.Split(path, "/")
				}
				// [scripts Duchess.lua]
				// [scripts Canfields Duchess.lua]
				if len(splits) == 3 {
					sinfo.group = "> " + splits[1]
				}
				// fmt.Println(splits)
				files = append(files, sinfo)
			}
			return nil
		})

		if err != nil {
			fmt.Println(err)
		} else {
			// fmt.Println(files)
			for _, sinfo := range files {
				variants[sinfo.name] = &MoonGame{scriptBase: scriptBase{fname: sinfo.path}}
				if sinfo.group != "" {
					if _, ok := variantGroups[sinfo.group]; !ok {
						variantGroups[sinfo.group] = []string{sinfo.name}
					} else {
						// You can't change values associated with keys in a map, you can only reassign values
						gameNames := variantGroups[sinfo.group]
						gameNames = append(gameNames, sinfo.name)
						variantGroups[sinfo.group] = gameNames
					}
				}
			}
		}

		// alternatives to filepath.Walk
		//
		// files, err := filepath.Glob("./scripts/*.lua")
		// if err != nil {
		// 	log.Println(err) // "open scripts: no such file or directory"
		// } else {
		// 	fmt.Println(files) // [scripts/FreecellScript.lua]
		// }
		// entries, err := os.ReadDir("scripts")
		// if err != nil {
		// 	log.Println(err) // "open scripts: no such file or directory"
		// } else {
		// 	for _, e := range entries {
		// 		println(e.Name())
		// 	}
		// }
	}

	var vnames []string = make([]string, 0, len(variants))
	for k := range variants {
		vnames = append(vnames, k)
	}

	// no need to sort here, sort gets done on-demand by func VariantNames()
	variantGroups["> All"] = vnames
	variantGroups["> All by Played"] = vnames
}

func (d *dark) ListVariantGroups() []string {
	var vnames []string = make([]string, 0, len(variantGroups))
	for k := range variantGroups {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

func (d *dark) ListVariants(group string) []string {
	var vnames []string = []string{}
	vnames = append(vnames, variantGroups[group]...)
	if group == "> All by Played" {
		sort.Slice(vnames, func(i, j int) bool {
			return d.stats.played(vnames[i]) > d.stats.played(vnames[j])
		})
	} else {
		sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	}
	return vnames
}
