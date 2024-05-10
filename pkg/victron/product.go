package victron

type ProductID int

func (pid ProductID) String() string {
	switch pid {
	case PIDBMV700:
		return "BMV 700"
	// Add other cases similarly...
	case PIDBlueSolarMPPT15070Rev2:
		return "BlueSolar MPPT 150/70 Rev2"
	default:
		return "Unknown"
	}
}

func ProductIDFromString(s string) ProductID {
	switch s {
	case "BMV 700":
		return PIDBMV700
	// Add other cases similarly...
	case "BlueSolar MPPT 150/70 Rev2":
		return PIDBlueSolarMPPT15070Rev2
	default:
		return 0
	}
}

type Category struct {
	RequiredFields map[string]struct{}
}

func (c *Category) Validate(b *Blocks) bool {
	for field := range b.Fields {
		if _, ok := fields[field]; !ok {
			return false
		}
	}
	return true
}

var Categories = map[ProductID]*Category{
	PIDBMV700: {
		RequiredFields: map[string]struct{}{
			"V": {}, "I": {}, // Define all expected fields for BMV 700
		},
	},
	// Define other categories similarly...
	PIDBlueSolarMPPT15070Rev2: {
		RequiredFields: map[string]struct{}{
			"V": {}, "PPV": {}, // Define all expected fields for BlueSolar MPPT 150/70 Rev2
		},
	},
}
