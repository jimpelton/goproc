package math

type Vec3UI64 [3]uint64

// CompProduct computes the component wise product of the vector.
func (v Vec3UI64) CompProduct() uint64 {
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

func (v Vec3UI64) CompDiv(v2 Vec3UI64) Vec3UI64 {
	return Vec3UI64{
		v.X() / v2.X(),
		v.Y() / v2.Y(),
		v.Z() / v2.Z(),
	}
}

func (v Vec3UI64) CompMul(v2 Vec3UI64) Vec3UI64 {
	return Vec3UI64{
		v.X() * v2.X(),
		v.Y() * v2.Y(),
		v.Z() * v2.Z(),
	}
}
