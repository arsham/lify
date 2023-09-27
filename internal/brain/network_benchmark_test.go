package brain_test

import (
	"testing"

	"github.com/arsham/neuragene/internal/brain"
)

func BenchmarkNetwork(b *testing.B) {
	nn, err := brain.New(&brain.Config{
		InputNeurons: 8,
		HiddenLayer: brain.Layer{
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
		OutputLayer: brain.Layer{
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
	if err != nil {
		b.Fatal(err)
	}
	input := []float64{0.1, 0.2, 0.3, -0.1, 0.15, 1, 0, 0.3}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := nn.Predict(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
