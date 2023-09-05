package genome

// These constants determine the index of the trait in the traits array. These
// should not be accesses from outside of this package, unless for testing.
const (
	IndexNutritionCunsumption = iota
	IndexGrowth
	IndexMaxGrowth
)

var charValues = map[rune]int32{}

func init() {
	runes := "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i, r := range runes {
		charValues[r] = int32(i)
	}
}

// NutritionConsumption returns the nutrition consumption trait value for the
// given DNA.
func NutritionConsumption(dna *DNA) int32 {
	return charValues[dna.traits[IndexNutritionCunsumption]]
}

// Growth returns the growth value when eating food.
func Growth(dna *DNA) int32 {
	return charValues[dna.traits[IndexGrowth]]
}

// MaxGrowth returns the maximum an organism can grow.
func MaxGrowth(dna *DNA) int32 {
	return charValues[dna.traits[IndexMaxGrowth]]
}
