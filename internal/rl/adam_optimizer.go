// Package rl provides reinforcement learning functionality
package rl

import (
	"math"
)

// Adam implements the Adam optimization algorithm
type Adam struct {
	LearningRate float64
	Beta1        float64
	Beta2        float64
	Epsilon      float64
	M            []float64
	V            []float64
	T            int
}

// NewAdam creates a new Adam optimizer
func NewAdam(learningRate, beta1, beta2 float64) *Adam {
	return &Adam{
		LearningRate: learningRate,
		Beta1:        beta1,
		Beta2:        beta2,
		Epsilon:      1e-8,
		M:            []float64{},
		V:            []float64{},
		T:            0,
	}
}

// Init initializes the optimizer for the given parameter size
func (o *Adam) Init(paramSize int) {
	o.M = make([]float64, paramSize)
	o.V = make([]float64, paramSize)
	o.T = 0
}

// Step performs one optimization step
func (o *Adam) Step(params []float64, grads []float64) []float64 {
	if len(o.M) == 0 {
		o.Init(len(params))
	}

	o.T++

	for i := 0; i < len(params); i++ {
		if i >= len(grads) {
			continue
		}

		// Update biased first moment estimate
		o.M[i] = o.Beta1*o.M[i] + (1-o.Beta1)*grads[i]

		// Update biased second raw moment estimate
		o.V[i] = o.Beta2*o.V[i] + (1-o.Beta2)*grads[i]*grads[i]

		// Compute bias-corrected first moment estimate
		mCorrected := o.M[i] / (1 - math.Pow(o.Beta1, float64(o.T)))

		// Compute bias-corrected second raw moment estimate
		vCorrected := o.V[i] / (1 - math.Pow(o.Beta2, float64(o.T)))

		// Update parameters
		params[i] -= o.LearningRate * mCorrected / (math.Sqrt(vCorrected) + o.Epsilon)
	}

	return params
}
