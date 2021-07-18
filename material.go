package main

import (
	"math/rand"

	m "github.com/go-gl/mathgl/mgl64"
)

type Material interface {
	Scatter(inRay Ray, Normal m.Vec3, rander *rand.Rand) ScatterInfo
}
type Diffuse struct {
	Albedo      m.Vec3
	Attenuation float64
}

func (d Diffuse) Scatter(inRay Ray, Normal m.Vec3, rander *rand.Rand) ScatterInfo {
	return ScatterInfo{
		NewRay:      true,
		NewRayDir:   Normal.Add(RandomUnitVec3(rander)),
		Color:       d.Albedo,
		Attenuation: d.Attenuation,
	}
}

type Metal struct {
	Albedo      m.Vec3
	Reflectance float64
	Fuzziness float64
}

func (m Metal) Scatter(inRay Ray, Normal m.Vec3, rander *rand.Rand) ScatterInfo {
	return ScatterInfo{
		NewRay:      true,
		NewRayDir:   reflect(inRay.Dir, Normal).Add(RandomVec3InUnitSphere(rander).Mul(m.Fuzziness)),
		Color:       m.Albedo,
		Attenuation: m.Reflectance,
	}

}

type ScatterInfo struct {
	NewRay           bool
	Attenuation      float64
	NewRayDir, Color m.Vec3
}
