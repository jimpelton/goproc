package trfunc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TFOpacityFromString(t *testing.T) {
	tfString := `0.0 1.0
		0.1 	0.5
        0.5 0.5
        1.0 0.0`
	tf, err := OpenTFOpacityString(tfString)
	if err != nil {
		t.Error("failed opening tf from string:", err.Error())
	}

	expected := []TFOpacityKnot{
		{0.0, 1.0},
		{0.1, 0.5},
		{0.5, 0.5},
		{1.0, 0.0},
	}

	// check length is equal for both
	if len(tf.knots) != len(expected) {
		t.Error("lengths not same")
	}
	// compare each knot
	for i, k := range tf.knots {
		if k != expected[i] {
			t.Error("knots at index", i, "were not equal")
		}
	}
}

func TestTFOpacity_Interpolate(t *testing.T) {
	tfString := `0.00  0.00	
               0.16  0.00
                0.33  1.00
                0.50  1.00
                0.66  1.00
                0.83  0.50
                1.00  0.00`

	tf, err := OpenTFOpacityString(tfString)
	if err != nil {
		t.Error("failed opening tf from string:", err.Error())
	}

	t.Run("otf returns correct values on exact match", func(t *testing.T) {
		assert.Equal(t, 0.0, tf.Interpolate(0.0), "exact value not as expected")
		assert.Equal(t, 0.0, tf.Interpolate(0.16), "exact value not as expected")
		assert.Equal(t, 1.0, tf.Interpolate(0.5), "exact value not as expected")
		assert.Equal(t, 1.0, tf.Interpolate(0.66), "exact value not as expected")
		assert.Equal(t, 0.50, tf.Interpolate(0.83), "exact value not as expected")
		assert.Equal(t, 0.0, tf.Interpolate(1.0), "exact value not as expected")
	})

	t.Run("otf returns interpolated values", func (t *testing.T) {

		assert.Equal(t, 0.7205882352, tf.Interpolate(0.755), "interpolated value not as expected")
		assert.Equal(t, 1.0, tf.Interpolate(0.55), "interpolated value not as expected")
	})
}
