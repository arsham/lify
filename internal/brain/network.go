package brain

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"

	"github.com/arsham/neuragene/internal/config"
)

// Layer is a layer of the neural network.
type Layer struct {
	Weights []float64
	Biases  []float64
}

// Config is the configuration for the neural network. On production you should
// not set the TestCheck to true, otherwise it will check the input values and
// it will slow down the process.
type Config struct {
	HiddenLayer   Layer
	OutputLayer   Layer
	InputNeurons  int
	OutputNeurons int
	TestCheck     bool // for testing
}

type layer struct {
	weights *mat.Dense
	biases  *mat.Dense
}

// Network is a fully connected multi-layer network with one hidden layer.
type Network struct {
	hidden       layer
	output       layer
	inputNeurons int
}

// New creates a neural network with the input weights, hidden weights and with
// the amount of given output nodes.
func New(c *Config) (*Network, error) {
	if c.TestCheck {
		// Only applied in testing to reduce the amount of checks in
		// production.
		if c.InputNeurons < 1 {
			return nil, errors.New("zero input neurons")
		}
		if c.OutputNeurons < 1 {
			return nil, errors.New("zero output neurons")
		}
		if len(c.HiddenLayer.Weights) != c.InputNeurons*len(c.HiddenLayer.Biases) {
			return nil, fmt.Errorf("hidden layer weights size %d does not match %d", len(c.HiddenLayer.Weights), c.InputNeurons*len(c.HiddenLayer.Biases))
		}
		if len(c.OutputLayer.Weights) != len(c.HiddenLayer.Biases)*c.OutputNeurons {
			return nil, fmt.Errorf("output layer weights size %d does not match %d", len(c.OutputLayer.Weights), len(c.HiddenLayer.Biases)*c.OutputNeurons)
		}
	}

	hiddenNeurons := len(c.HiddenLayer.Weights) / c.InputNeurons
	n := &Network{
		inputNeurons: c.InputNeurons,
		hidden: layer{
			weights: mat.NewDense(c.InputNeurons, hiddenNeurons, c.HiddenLayer.Weights),
			biases:  mat.NewDense(1, hiddenNeurons, c.HiddenLayer.Biases),
		},
		output: layer{
			weights: mat.NewDense(hiddenNeurons, c.OutputNeurons, c.OutputLayer.Weights),
			biases:  mat.NewDense(1, c.OutputNeurons, c.OutputLayer.Biases),
		},
	}
	return n, nil
}

func sigmoid(_, _ int, z float64) float64 {
	return 1.0 / (1 + math.Exp(-z))
}

// Predict makes a prediction based on a trained neural network.
func (n *Network) Predict(input []float64) ([]float64, error) {
	if len(input) != n.inputNeurons {
		return nil, fmt.Errorf("wrong input size: %d, want %d", len(input), n.inputNeurons)
	}
	hiddenIn := &mat.Dense{}
	inputMat := mat.NewDense(1, len(input), input)
	hiddenIn.Mul(inputMat, n.hidden.weights)
	hiddenIn.Apply(func(_, j int, v float64) float64 {
		return v + n.hidden.biases.At(0, j)
	}, hiddenIn)

	hiddenActivation := &mat.Dense{}
	hiddenActivation.Apply(sigmoid, hiddenIn)

	outputIn := &mat.Dense{}
	outputIn.Mul(hiddenActivation, n.output.weights)
	outputIn.Apply(func(_, col int, v float64) float64 {
		return v + n.output.biases.At(0, col)
	}, outputIn)

	output := &mat.Dense{}
	output.Apply(sigmoid, outputIn)
	return output.RawMatrix().Data, nil
}

const (
	hiddenWeightsBlob = "hidden_weights.blob"
	hiddenBiasesBlob  = "hidden_biases.blob"
	outputWeightsBlob = "output_weights.blob"
	outputBiasesBlob  = "output_biases.blob"
)

// Save creates a tar file and saves the neural network to it.
func (n Network) Save(filename string) error {
	f, err := os.Create(filename) // nolint:gosec // user will provide the file.
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			config.Logger().Error("closing file", "error", err)
		}
	}()

	tw := tar.NewWriter(f)
	defer func() {
		if err := tw.Close(); err != nil {
			config.Logger().Error("closing tar writer", "error", err)
		}
	}()

	blobs := []struct {
		matric *mat.Dense
		name   string
	}{
		{n.hidden.weights, hiddenWeightsBlob},
		{n.hidden.biases, hiddenBiasesBlob},
		{n.output.weights, outputWeightsBlob},
		{n.output.biases, outputBiasesBlob},
	}

	for _, b := range blobs {
		blobData, err := b.matric.MarshalBinary()
		if err != nil {
			return fmt.Errorf("marshalling %s: %w", b.name, err)
		}
		blobHdr := &tar.Header{
			Name: b.name,
			Size: int64(len(blobData)),
		}
		if err := tw.WriteHeader(blobHdr); err != nil {
			return fmt.Errorf("writing %s header: %w", b.name, err)
		}
		if _, err := tw.Write(blobData); err != nil {
			return fmt.Errorf("writing %s: %w", b.name, err)
		}
	}

	if err := tw.Flush(); err != nil {
		return fmt.Errorf("flushing tar writer: %w", err)
	}

	return nil
}

// Load reads a tar file and loads the neural network from it.
func Load(filename string) (*Network, error) {
	n := &Network{}
	f, err := os.Open(filename) // nolint:gosec // user will provide the file.
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			config.Logger().Error("closing file", "error", err)
		}
	}()

	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar header: %w", err)
		}

		switch hdr.Name {
		case hiddenWeightsBlob:
			hidWeights := &mat.Dense{}
			if err := hidWeights.UnmarshalBinary(make([]byte, hdr.Size)); err != nil {
				return nil, fmt.Errorf("unmarshalling hidden weights: %w", err)
			}
			n.hidden.weights = hidWeights
		case hiddenBiasesBlob:
			hidBiases := &mat.Dense{}
			if err := hidBiases.UnmarshalBinary(make([]byte, hdr.Size)); err != nil {
				return nil, fmt.Errorf("unmarshalling hidden biases: %w", err)
			}
			n.hidden.biases = hidBiases
		case outputWeightsBlob:
			outWeights := &mat.Dense{}
			if err := outWeights.UnmarshalBinary(make([]byte, hdr.Size)); err != nil {
				return nil, fmt.Errorf("unmarshalling output weights: %w", err)
			}
			n.output.weights = outWeights
		case outputBiasesBlob:
			outBiases := &mat.Dense{}
			if err := outBiases.UnmarshalBinary(make([]byte, hdr.Size)); err != nil {
				return nil, fmt.Errorf("unmarshalling output biases: %w", err)
			}
			n.output.biases = outBiases
		default:
			return nil, fmt.Errorf("unknown file: %s", hdr.Name)
		}
	}

	return n, nil
}
