package rl

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// AdamOptimizer implements the Adam optimization algorithm
type AdamOptimizer struct {
	LearningRate float64
	Beta1        float64
	Beta2        float64
	M            *mat.Dense
	V            *mat.Dense
	T            int
	Epsilon      float64
}

// NewAdamOptimizer creates a new AdamOptimizer with the given parameters
func NewAdamOptimizer(learningRate, beta1, beta2 float64) *AdamOptimizer {
	return &AdamOptimizer{
		LearningRate: learningRate,
		Beta1:        beta1,
		Beta2:        beta2,
		Epsilon:      1e-8,
		T:            0,
	}
}

// Update updates the weights using the Adam optimization algorithm
func (a *AdamOptimizer) Update(weights *mat.Dense, gradient float64) {
	a.T++

	rows, cols := weights.Dims()

	// Initialize momentum and velocity if not already done
	if a.M == nil {
		a.M = mat.NewDense(rows, cols, nil)
		a.V = mat.NewDense(rows, cols, nil)
	}

	// Update momentum and velocity
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			// Simplified implementation using a single gradient value
			// In a real implementation, you would have a gradient matrix
			g := gradient * weights.At(i, j)

			// Update biased first moment estimate
			m := a.M.At(i, j)
			m = a.Beta1*m + (1-a.Beta1)*g
			a.M.Set(i, j, m)

			// Update biased second raw moment estimate
			v := a.V.At(i, j)
			v = a.Beta2*v + (1-a.Beta2)*g*g
			a.V.Set(i, j, v)

			// Compute bias-corrected first moment estimate
			mCorrected := m / (1 - math.Pow(a.Beta1, float64(a.T)))

			// Compute bias-corrected second raw moment estimate
			vCorrected := v / (1 - math.Pow(a.Beta2, float64(a.T)))

			// Update parameters
			weights.Set(i, j, weights.At(i, j)-a.LearningRate*mCorrected/(math.Sqrt(vCorrected)+a.Epsilon))
		}
	}
}
