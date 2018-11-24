package math

// UNorm performs unity-based normalization.
//
// Scales v to within [0 .. 1]
func UNorm(v, min, max float64) float64 {
	return (v - min) / (max - min)
}
