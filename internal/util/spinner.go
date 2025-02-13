package util

var CustomSpinner = []string{
	// Transition: empty → growing left fill (0 → 3 equals)
	"[            ]", // 0: empty
	"[=           ]", // 1: 1 equals
	"[==          ]", // 2: 2 equals
	"[===         ]", // 3: 3 equals

	// Transition: growing left fill (3 → 6 equals)
	"[====        ]", // 4: 4 equals
	"[=====       ]", // 5: 5 equals
	"[======      ]", // 6: 6 equals

	// Transition: growing left fill (6 → 9 equals)
	"[=======     ]", // 7: 7 equals
	"[========    ]", // 8: 8 equals
	"[=========   ]", // 9: 9 equals

	// Transition: slide block from left to right (keep 9 equals)
	"[ =========  ]", // 10: shift right by 1 space
	"[  ========= ]", // 11: shift right by 2 spaces
	"[   =========]", // 12: 9 equals flush right

	// Transition: shrinking right fill (9 → 6 equals)
	"[    ========]", // 13: 8 equals (right-aligned)
	"[     =======]", // 14: 7 equals
	"[      ======]", // 15: 6 equals

	// Transition: shrinking right fill (6 → 3 equals)
	"[       =====]", // 16: 5 equals
	"[        ====]", // 17: 4 equals
	"[         ===]", // 18: 3 equals

	// Transition: disappearing right fill (3 → 0 equals)
	"[          ==]", // 19: 2 equals
	"[           =]", // 20: 1 equals
	"[            ]", // 21: empty

	// Transition: empty → growing right fill (0 → 3 equals)
	"[           =]", // 22: 1 equals (right-aligned)
	"[          ==]", // 23: 2 equals
	"[         ===]", // 24: 3 equals

	// Transition: growing right fill (3 → 6 equals)
	"[        ====]", // 25: 4 equals
	"[       =====]", // 26: 5 equals
	"[      ======]", // 27: 6 equals

	// Transition: growing right fill (6 → 9 equals)
	"[     =======]", // 28: 7 equals
	"[    ========]", // 29: 8 equals
	"[   =========]", // 30: 9 equals

	// Transition: growing right fill (9 → full 12 equals)
	"[  ==========]", // 31: 10 equals
	"[ ===========]", // 32: 11 equals
	"[============]", // 33: 12 equals (full)

	// Transition: slide block from full to left fill (12 → 9 equals)
	"[=========== ]", // 34: 11 equals flush left
	"[==========  ]", // 35: 10 equals flush left
	"[=========   ]", // 36: 9 equals flush left

	// Transition: shrinking left fill (9 → 6 equals)
	"[========    ]", // 37: 8 equals
	"[=======     ]", // 38: 7 equals
	"[======      ]", // 39: 6 equals

	// Transition: shrinking left fill (6 → 3 equals)
	"[=====       ]", // 40: 5 equals
	"[====        ]", // 41: 4 equals
	"[===         ]", // 42: 3 equals

	// Transition: disappearing left fill (3 → 0 equals)
	"[==          ]", // 43: 2 equals
	"[=           ]", // 44: 1 equals
	"[            ]", // 45: empty
}
