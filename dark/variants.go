package dark

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var variants = map[string]scripter{
	"Agnes Bernauer": &Agnes{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Agnes_(solitaire)",
			cardColors: 2,
		},
	},
	"Alhambra": &Alhambra{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Alhambra_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"American Toad": &Toad{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/American_Toad_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"Australian": &Australian{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Australian_Patience",
			cardColors: 4,
		},
	},
	"Baker's Dozen": &BakersDozen{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Baker%27s_Dozen_(solitaire)",
			cardColors: 1,
		},
	},
	"Baker's Game": &Freecell{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Baker%27s_Game",
			cardColors: 4,
		},
		tabCompareFunc: dyad.compare_DownSuit,
	},
	"Bisley": &Bisley{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Bisley_(card_game)",
			cardColors: 4,
		},
	},
	"Blind Freecell": &Freecell{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/FreeCell",
			cardColors: 2,
		},
		tabCompareFunc: dyad.compare_DownAltColor,
		blind:          true,
	},
	"Blockade": &Blockade{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Blockade_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"Canfield": &Canfield{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		},
		draw:     3,
		recycles: 32767,
	},
	"Rainbow Canfield": &Canfield{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		},
		draw:     1,
		recycles: 2,
	},
	"Selective Canfield": &Canfield{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		},
		draw:     3,
		recycles: 32767,
		variant:  "selective",
	},
	"Storehouse": &Canfield{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
			cardColors: 4,
		},
		draw:     1,
		recycles: 2,
		variant:  "storehouse",
	},
	"Chameleon": &Chameleon{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Chameleon_(solitaire)",
		},
	},
	"Colorado": &Colorado{
		scriptBase: scriptBase{
			packs: 2,
		},
	},
	"Duchess": &Duchess{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Duchess_(solitaire)",
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
	"Thoughtful": &Klondike{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
		draw:       1,
		recycles:   2,
		thoughtful: true,
	},
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
	"Eight Off": &EightOff{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Eight_Off",
			cardColors: 4,
		},
	},
	"Freecell": &Freecell{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/FreeCell",
		},
		tabCompareFunc: dyad.compare_DownAltColor,
	},
	"Freecell Easy": &Freecell{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/FreeCell",
		},
		tabCompareFunc: dyad.compare_DownAltColor,
		easy:           true,
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
	"Mrs Mop": &MrsMop{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Mrs._Mop",
			cardColors: 4,
			packs:      2,
		},
	},
	"Mrs Mop Easy": &MrsMop{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Mrs._Mop",
			cardColors: 4,
			packs:      2,
		},
		easy: true,
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
	"Seahaven Towers": &Seahaven{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Seahaven_Towers",
			cardColors: 4,
		},
	},
	"Simple Simon": &SimpleSimon{
		scriptBase: scriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Simple_Simon_(solitaire)",
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
	"Easthaven": &Westcliff{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Westcliff_(card_game)",
		},
		variant: "Easthaven",
	},
	"Whitehead": &Whitehead{
		scriptBase: scriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
	},
	"Uncle Sam": &UncleSam{
		scriptBase: scriptBase{
			cardColors: 2,
			packs:      2,
		},
	},
	"Usk": &Usk{
		scriptBase: scriptBase{
			wikipedia: "https://politaire.com/help/usk",
		},
		tableauLabel: "K",
	},
	"Usk Relaxed": &Usk{
		scriptBase: scriptBase{
			wikipedia: "https://politaire.com/help/usk",
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
	"> Canfields":     {"Canfield", "Chameleon", "Eagle Wing", "Rainbow Canfield", "Selective Canfield", "Storehouse", "Duchess", "American Toad"},
	"> Easier":        {"American Toad", "American Westcliff", "Blockade", "Classic Westcliff", "Lucas", "Spider One Suit", "Usk Relaxed"},
	"> Hapgood":       {"Light and Shadow", "The Rainbow", "Uncle Sam"},
	"> Harder":        {"Baker's Dozen", "Easthaven", "Forty Thieves", "Spider Four Suits", "Usk"},
	"> Forty Thieves": {"Forty Thieves", "Number Ten", "Red and Black", "Indian", "Rank and File", "Sixty Thieves", "Josephine", "Limited", "Forty and Eight", "Lucas", "Busy Aces", "Maria", "Streets"},
	"> Freecells":     {"Baker's Game", "Blind Freecell", "Freecell", "Freecell Easy", "Eight Off", "Seahaven Towers"},
	"> Klondikes":     {"Gargantua", "Triple Klondike", "Klondike", "Klondike Draw Three", "Thoughtful", "Whitehead"},
	"> People":        {"Agnes Bernauer", "Duchess", "Josephine", "Maria", "Simple Simon", "Baker's Game"},
	"> Places":        {"Australian", "Bisley", "Colorado", "Yukon", "Klondike", "Usk", "Usk Relaxed"},
	"> Puzzlers":      {"Bisley", "Usk", "Mrs Mop", "Penguin", "Simple Simon", "Baker's Dozen"},
	"> Spiders":       {"Spider One Suit", "Spider Two Suits", "Spider Four Suits", "Scorpion", "Spiderette"},
	"> Yukons":        {"Yukon", "Yukon Cells"},
}

// init is used to assemble the "> All" alpha-sorted group of variants
func init() {
	// look in the scripts folder for *.lua files
	{
		var files []string

		err := filepath.Walk("./scripts", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}
			// fmt.Println("filepath.Walk ./scripts", path, "dir", filepath.Dir(path), strings.Split(filepath.Dir(path), "/"))
			// if info.IsDir() {
			// 	splits := strings.Split(path, "/")
			// 	cat := splits[len(splits)-1]
			// 	fmt.Println(splits, cat)
			// }
			if !info.IsDir() && filepath.Ext(path) == ".lua" {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			fmt.Println(err)
		} else {
			for _, file := range files {
				fname := strings.TrimSuffix(filepath.Base(file), ".lua")
				// fmt.Println("found variant", fname)
				variants[fname] = &MoonGame{scriptBase: scriptBase{fname: file}}
			}
		}

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
