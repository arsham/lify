package brain

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestPredictTableDriven(t *testing.T) {
	t.Parallel()
	t.Run("testPredictTableDriven4x3x2", testPredictTableDriven4x3x2)
	t.Run("testPredictTableDriven8x6x10", testPredictTableDriven8x6x10)
}

func testPredictTableDriven4x3x2(t *testing.T) {
	nn, err := New(&Config{
		InputNeurons: 4,
		HiddenLayer: Layer{
			Weights: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.3, 0.2, 0.1, 0.2, 0.3, 0.4},
			Biases:  []float64{0.1, 0.2, 0.5},
		},
		OutputLayer: Layer{
			Weights: []float64{0.7, 0.8, 0.9, 1.0, 1.1, 1.2},
			Biases:  []float64{0.3, 0.4},
		},
		OutputNeurons: 2,
		TestCheck:     true,
	})
	assert.NoError(t, err)

	testCases := []struct {
		input []float64
		want  []float64
	}{
		{input: []float64{0.1, 0.2, 0.3, 0.1}, want: []float64{0.878995, 0.906172}},
		{input: []float64{0.5, 0.6, 0.7, 0.3}, want: []float64{0.903224, 0.927303}},
		{input: []float64{0.9, 0.8, 0.7, 0.5}, want: []float64{0.913622, 0.936113}},
		{input: []float64{0.3, 0.2, 0.1, 0.7}, want: []float64{0.890796, 0.916501}},
		{input: []float64{0.0, 0.0, 0.0, 0.0}, want: []float64{0.863779, 0.892522}},
		{input: []float64{1.0, 1.0, 1.0, 1.0}, want: []float64{0.924425, 0.945125}},
		{input: []float64{0.4, 0.3, 0.2, 0.8}, want: []float64{0.897452, 0.922285}},
		{input: []float64{0.6, 0.5, 0.4, 0.4}, want: []float64{0.900440, 0.924882}},
		{input: []float64{0.8, 0.7, 0.6, 0.2}, want: []float64{0.906109, 0.929751}},
		{input: []float64{0.2, 0.1, 0.0, 0.9}, want: []float64{0.889097, 0.915004}},
		{input: []float64{0.2, 0.1, 0.3, 0.1}, want: []float64{0.876966, 0.904366}},
	}

	for i, tc := range testCases {
		tc := tc
		name := fmt.Sprintf("Test case %d", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			output, err := nn.Predict(tc.input)
			assert.NoError(t, err)
			if !almostEqual(output, tc.want, 0.0001) {
				t.Errorf("Predicted output %v does not match expected output %v", output, tc.want)
			}
		})
	}
}

// almostEqual checks if two float64 slices are equal within a specified
// epsilon.
func almostEqual(a, b []float64, epsilon float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > epsilon {
			return false
		}
	}
	return true
}

func testPredictTableDriven8x6x10(t *testing.T) {
	t.Parallel()
	nn, err := New(&Config{
		InputNeurons: 8,
		HiddenLayer: Layer{
			Weights: []float64{
				0.1, 0.2, 0.3, 0.4, 0.5, 0.6, -0.3, -0.2, // 1st neuron
				0.1, 0.2, 0.3, 0.4, 0.5, 0.6, -0.3, -0.2, // 2nd neuron
				0.06, -0.05, 0.04, -0.03, 0.02, -0.01, 0.0, 0.1, // 3rd neuron
				-0.01, 0.02, -0.05, 0.1, -0.02, 0.05, 0.1, 0.2, // 4th neuron
				0.03, -0.04, 0.05, -0.06, 0.07, -0.08, 0.09, 0.1, // 5th neuron
				0.1, -0.1, 0.1, -0.1, 0.1, -0.1, 0.1, 0.1, // 6th neuron
			},
			Biases: []float64{0.01, 0.02, 0.05, 0.01, -0.02, 0.05},
		},
		OutputLayer: Layer{
			Weights: []float64{
				0.7, 0.8, 0.9, 1.0, 1.1, 1.2, // 1st neuron
				6.7, -6.8, 6.9, 7.0, -7.1, 7.2, // 2nd neuron
				0.3, 0.4, 0.5, 0.6, 0.7, 0.8, // 3rd neuron
				0.9, 1.0, 1.1, -1.2, 1.3, 1.4, // 4th neuron
				1.5, 1.6, 1.7, 1.8, 1.9, 2.0, // 5th neuron
				2.1, 2.2, -2.3, 2.4, 2.5, 2.6, // 6th neuron
				2.7, 2.8, 2.9, 3.0, 3.1, 3.2, // 7th neuron
				3.3, 3.4, 3.5, 3.6, 3.7, 3.8, // 8th neuron
				3.9, 4.0, 4.1, 4.2, 4.3, 4.4, // 9th neuron
				4.5, 4.6, 4.7, 4.8, 4.9, -5.0, // 10th neuron
			},
			Biases: []float64{0.03, 0.04, 0.05, 0.06, 0.07, -0.8, 0.09, 0, 0.01, 0.02},
		},
		OutputNeurons: 10,
		TestCheck:     true,
	})
	assert.NoError(t, err)

	tcs := []struct {
		input []float64
		want  int
	}{
		{input: []float64{0.1, 0.2, 0.3, -0.1, 0.15, 1, 0, 0.3}, want: 8},
		{input: []float64{6, -10, -10, -10, 0, 6, -10, -10}, want: 1},
		{input: []float64{-20, 0, 0, 0, 0, 0, 0, 0}, want: 6},
	}
	for i, tc := range tcs {
		tc := tc
		name := fmt.Sprintf("Test case %d", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			output, err := nn.Predict(tc.input)
			assert.NoError(t, err)
			got := slices.Index(output, slices.Max(output))
			assert.Equal(t, tc.want, got)
		})
	}
}
