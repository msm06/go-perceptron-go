// Neural provides struct to represents most common neural networks model and algorithms to train / test them.
package neural

import (

	// sys import
	"math/rand"
	"os"

	// third part import
	log "github.com/sirupsen/logrus"

	// this repo internal import
	mu "github.com/made2591/go-perceptron-go/util"
)

const (

	SCALING_FACTOR = 0.0000000000001

)

// Neuron struct represents a simple Neuron network with a slice of n weights.
type Neuron struct {

	// Weights represents Neuron vector representation
	Weights []float64
	// Bias represents Neuron natural propensity to spread signal
	Bias float64
	// Lrate represents learning rate of neuron
	Lrate float64

	// Value represents desired value when loading input into network in Multi Layer Perceptron
	Value float64
	// Delta represents delta error for unit
	Delta float64

}

// #######################################################################################

func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// RandomNeuronInit initialize neuron weight, bias and learning rate using NormFloat64 random value.
func RandomNeuronInit(neuron *Neuron, dim int) {

	neuron.Weights = make([]float64, dim)

	// init random weights
	for index, _ := range neuron.Weights {
		// init random threshold weight
		neuron.Weights[index] = rand.NormFloat64() * SCALING_FACTOR
	}

	// init random bias and lrate
	neuron.Bias  = rand.NormFloat64() * SCALING_FACTOR
	neuron.Lrate = rand.NormFloat64() * SCALING_FACTOR
	neuron.Value = rand.NormFloat64() * SCALING_FACTOR
	neuron.Delta = rand.NormFloat64() * SCALING_FACTOR

	log.WithFields(log.Fields{
		"level":   "debug",
		"place":   "neuron",
		"func":    "RandomNeuronInit",
		"msg":     "random neuron weights init",
		"weights": neuron.Weights,
	}).Debug()

}

// UpdateWeights performs update in neuron weights with respect to passed stimulus.
// It returns error of prediction before and after updating weights.
func UpdateWeights(neuron *Neuron, stimulus *Stimulus) (float64, float64) {

	// compute prediction value and error for stimulus given neuron BEFORE update (actual state)
	var predictedValue, prevError, postError float64 = Predict(neuron, stimulus), 0.0, 0.0
	prevError = stimulus.Expected - predictedValue

	// performs weights update for neuron
	neuron.Bias = neuron.Bias + neuron.Lrate*prevError

	// performs weights update for neuron
	for index, _ := range neuron.Weights {
		neuron.Weights[index] = neuron.Weights[index] + neuron.Lrate*prevError*stimulus.Dimensions[index]
	}

	// compute prediction value and error for stimulus given neuron AFTER update (actual state)
	predictedValue = Predict(neuron, stimulus)
	postError = stimulus.Expected - predictedValue

	log.WithFields(log.Fields{
		"level":   "debug",
		"place":   "neuron",
		"func":    "UpdateWeights",
		"msg":     "updating weights of neuron",
		"weights": neuron.Weights,
	}).Debug()

	// return errors
	return prevError, postError

}

// TrainNeuron trains a passed neuron with stimuli passed, for specified number of epoch.
// If init is 0, leaves weights unchanged before training.
// If init is 1, reset weights and bias of neuron before training.
func TrainNeuron(neuron *Neuron, stimuli []Stimulus, epochs int, init int) {

	// init weights if specified
	if init == 1 {
		neuron.Weights = make([]float64, len(stimuli[0].Dimensions))
		neuron.Bias = 0.0
	}

	// init counter
	var epoch int = 0

	// accumulator errors prev and post weights updates
	var squaredPrevError, squaredPostError float64 = 0.0, 0.0

	// in each epoch
	for epoch < epochs {

		// update weight using each stimulus in training set
		for _, stimulus := range stimuli {
			prevError, postError := UpdateWeights(neuron, &stimulus)
			// NOTE: in each step, use weights already updated by previous
			squaredPrevError = squaredPrevError + (prevError * prevError)
			squaredPostError = squaredPostError + (postError * postError)
		}

		log.WithFields(log.Fields{
			"level":            "debug",
			"place":            "error evolution in epoch",
			"method":           "TrainNeuron",
			"msg":              "epoch and squared errors reached before and after updating weights",
			"epochReached":     epoch + 1,
			"squaredErrorPrev": squaredPrevError,
			"squaredErrorPost": squaredPostError,
		}).Debug()

		// increment epoch counter
		epoch++

	}

}

// Predict performs a neuron prediction to passed stimulus.
// It returns a float64 binary predicted value.
func Predict(neuron *Neuron, stimulus *Stimulus) float64 {

	if mu.ScalarProduct(neuron.Weights, stimulus.Dimensions)+neuron.Bias < 0.0 {
		return 0.0
	}
	return 1.0

}

// Accuracy calculate percentage of equal values between two float64 based slices.
// It returns int number and a float64 percentage value of corrected values.
func Accuracy(actual []float64, predicted []float64) (int, float64) {

	// if slices have different number of elements
	if len(actual) != len(predicted) {
		log.WithFields(log.Fields{
			"level":        "error",
			"place":        "neuron",
			"method":       "Accuracy",
			"msg":          "accuracy between actual and predicted slices of values",
			"actualLen":    len(actual),
			"predictedLen": len(predicted),
		}).Error("Failed to compute accuracy between actual values and predictions: different length.")
		return -1, -1.0
	}

	// init result
	var correct int = 0

	for index, value := range actual {
		if value == predicted[index] {
			correct++
		}
	}

	// return correct
	return correct, float64(correct) / float64(len(actual)) * 100.0

}
