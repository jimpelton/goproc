package math

type Vec3UI64 [3]uint64

// CProduct computes the component wise product of the vector.
func (v Vec3UI64) CProduct() uint64 {
	return v[0] * v[1] * v[2]
}

func (v Vec3UI64) X() uint64 {
	return v[0]
}

func (v Vec3UI64) Y() uint64 {
	return v[1]
}

func (v Vec3UI64) Z() uint64 {
	return v[2]
}

