package genome_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/arsham/lify/internal/genome"
)

func TestDNA(t *testing.T) {
	t.Parallel()
	t.Run("SetAndTraitAt", testDNASetAndTraitAt)
	t.Run("CalculateDifference", testDNACalculateDifference)
	t.Run("TraitStrength", testDNATraitStrength)
	t.Run("IsCompatibleWith", testDNAIsCompatibleWith)
	t.Run("CreateOffspring", testDNACreateOffspring)
}

func testDNASetAndTraitAt(t *testing.T) {
	t.Parallel()
	dna := genome.NewDNA(24)
	dna.SetTrait(0, 'a')
	dna.SetTrait(12, '6')

	got := dna.TraitAt(0)
	assert.Equal(t, got, 'a')

	got = dna.TraitAt(12)
	assert.Equal(t, got, '6')

	got = dna.TraitAt(100)
	assert.Equal(t, got, -1)

	got = dna.TraitAt(-10)
	assert.Equal(t, got, -1)
}

func testDNACalculateDifference(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		dna1, dna2 string
		around     float64
	}{
		"not compatible": {
			dna1:   "aaaaaaaaaaaaaaaaaaaaaaa",
			dna2:   "aaaaaaaaaaaaaaaaaaaaaaaa",
			around: 100,
		},
		"same": {
			dna1:   "aaaaaaaaaaaaaaaaaaaaaaaa",
			dna2:   "aaaaaaaaaaaaaaaaaaaaaaaa",
			around: 0,
		},
		"slightly different": {
			dna1:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			dna2:   "aaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			around: 0.033,
		},
		"a bit more different": {
			dna1:   "aaaaaaaaaaaaaaaaaaaaaagaaaaaaaaaaaaaaaaaaaaaaaaaa",
			dna2:   "aaaaaaaaaaaabaaaaaaaaaiaaaaaaaaaaaaaaaaaaaaaaaaaa",
			around: 0.098,
		},
		"very different": {
			dna1:   "aaaabaaaaaaaaaaaaaaaaaaaaa0aaaaaaaaaaaaaaaaaaaaaa",
			dna2:   "aaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			around: 1.61,
		},
		"different in different places": {
			dna1:   "aaaaZaaaaaaBaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			dna2:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaZaaBaaaaaaa",
			around: 100,
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			dna1 := genome.NewDNAFromString(tc.dna1)
			dna2 := genome.NewDNAFromString(tc.dna2)
			diff := dna1.CalculateDifference(dna2)
			assert.True(t, diff >= 0, "got %f, want a positive number", diff)
			// The absolute difference between diff and tc.around should be at
			// most 0.01.
			assert.True(t, math.Abs(diff-tc.around) <= 0.01, fmt.Sprintf("got %f, want %f", diff, tc.around))
		})
	}
}

func testDNATraitStrength(t *testing.T) {
	t.Parallel()
	dna := genome.NewDNAFromString("abcdefxyz123456789ABCDEFXYZ")
	tcs := []struct {
		index int
		want  int
	}{
		{index: 0, want: 97},
		{index: 1, want: 98},
		{index: 2, want: 99},
		{index: 3, want: 100},
		{index: 4, want: 101},
		{index: 5, want: 102},
		{index: 6, want: 120},
		{index: 7, want: 121},
		{index: 8, want: 122},
		{index: 9, want: 49},
		{index: 10, want: 50},
		{index: 11, want: 51},
		{index: 12, want: 52},
		{index: 13, want: 53},
		{index: 14, want: 54},
		{index: 15, want: 55},
		{index: 16, want: 56},
		{index: 17, want: 57},
		{index: 18, want: 65},
		{index: 19, want: 66},
		{index: 20, want: 67},
		{index: 21, want: 68},
		{index: 22, want: 69},
		{index: 23, want: 70},
		{index: 24, want: 88},
		{index: 25, want: 89},
		{index: 26, want: 90},
		{index: 27, want: -1},
	}
	for i, tc := range tcs {
		tc := tc
		name := fmt.Sprintf("case %d", i)
		t.Run(name, func(t *testing.T) {
			got := dna.TraitStrength(tc.index)
			assert.Equal(t, tc.want, got)
		})
	}
}

func testDNAIsCompatibleWith(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		dna1, dna2 string
		compatible bool
	}{
		"not compatible": {
			dna1:       "aaaaaaaaaaaaaaaaaaaaaaa",
			dna2:       "aaaaaaaaaaaaaaaaaaaaaaaa",
			compatible: false,
		},
		"same": {
			dna1:       "aaaaaaaaaaaaaaaaaaaaaaaa",
			dna2:       "aaaaaaaaaaaaaaaaaaaaaaaa",
			compatible: true,
		},
		"slightly different": {
			dna1:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			dna2:       "aaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			compatible: true,
		},
		"a bit more different": {
			dna1:       "aaaaaaaaaaaaaaaaaaaaaagaaaaaaaaaaaaaaaaaaaaaaaaaa",
			dna2:       "aaaaaaaaaaaabaaaaaaaaaiaaaaaaaaaaaaaaaaaaaaaaaaaa",
			compatible: false,
		},
		"different in different places": {
			dna1:       "aaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaa",
			dna2:       "aaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaabaaa",
			compatible: false,
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			dna1 := genome.NewDNAFromString(tc.dna1)
			dna2 := genome.NewDNAFromString(tc.dna2)
			ok1 := dna1.IsCompatibleWith(dna2)
			ok2 := dna2.IsCompatibleWith(dna1)
			assert.Equal(t, tc.compatible, ok1)
			assert.Equal(t, tc.compatible, ok2)
		})
	}
}

func testDNACreateOffspring(t *testing.T) {
	t.Parallel()
	// giving the mutation a high chance to happen.
	for i := 0; i < 10000; i++ {
		tcs := []struct {
			dna1, dna2 string
			want       string
		}{{
			dna1: "aaaaaaaaaaaaaaaaaaaaaaa",
			dna2: "aaaaaaaaaaaaaaaaaaaaaaa",
			want: "aaaaaaaaaaaaaaaaaaaaaaa",
		}, {
			dna1: "aaaaaasaaaafayaaaBaaaaa",
			dna2: "aaaaaaTaaaagaxaaabaaaaa",
			want: "aaaaaasaaaagayaaabaaaaa",
		}, {
			dna1: "aaaaaaBaaaaaaagaaTaacaa",
			dna2: "aaaaaaCaaaaaaaaaataaCaa",
			want: "aaaaaaCaaaaaaagaataacaa",
		}, {
			dna1: "aaaaa113451111111aaaaaa",
			dna2: "aaaaa111345111111aaaaaa",
			want: "aaaaa113455111111aaaaaa",
		}}
		for i, tc := range tcs {
			tc := tc
			name := fmt.Sprintf("case %d", i)
			t.Run(name, func(t *testing.T) {
				dna1 := genome.NewDNAFromString(tc.dna1)
				dna2 := genome.NewDNAFromString(tc.dna2)
				child := genome.CreateOffspring(dna1, dna2)
				want := genome.NewDNAFromString(tc.want)
				diff := child.CalculateDifference(want)
				assert.True(t, diff <= 0.1, "\nwant %s\ngot  %s (%f)", want, child, diff)
			})
		}
	}
}
