package meander

import "strings"

// Cost is enum
type Cost int8

const (
	_ Cost = iota
	// Cost1 is
	Cost1
	// Cost2 is
	Cost2
	// Cost3 is
	Cost3
	// Cost4 is
	Cost4
	// Cost5 is
	Cost5
)

var costStrings = map[string]Cost{
	"$":     Cost1,
	"$$":    Cost2,
	"$$$":   Cost3,
	"$$$$":  Cost4,
	"$$$$$": Cost5,
}

func (l Cost) String() string {
	for s, v := range costStrings {
		if l == v {
			return s
		}
	}
	return "invalid value"
}

// ParseCost return Cost type
func ParseCost(s string) Cost {
	return costStrings[s]
}

// CostRange is cost range
type CostRange struct {
	From Cost
	To   Cost
}

func (r CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

// ParseCostRange parse string to CostRange type.
func ParseCostRange(s string) *CostRange {
	segs := strings.Split(s, "...")
	return &CostRange{
		From: ParseCost(segs[0]),
		To:   ParseCost(segs[1]),
	}
}
