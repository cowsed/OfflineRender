package main

import (
	"math/rand"

	m "github.com/go-gl/mathgl/mgl64"
)

type ScatterInfo struct {
	NewRay           bool
	Attenuation      float64
	NewRayDir, Color m.Vec3
}

type Material interface {
	Scatter(inRay Ray, Normal m.Vec3, env Environment, rander *rand.Rand) ScatterInfo
}
type Diffuse struct {
	Albedo      m.Vec3
	Attenuation float64
}

func (d Diffuse) Scatter(inRay Ray, Normal m.Vec3, env Environment, rander *rand.Rand) ScatterInfo {
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
	Fuzziness   float64
}

func (m Metal) Scatter(inRay Ray, Normal m.Vec3, env Environment, rander *rand.Rand) ScatterInfo {
	return ScatterInfo{
		NewRay:      true,
		NewRayDir:   reflect(inRay.Dir, Normal).Add(RandomVec3InUnitSphere(rander).Mul(m.Fuzziness)),
		Color:       m.Albedo,
		Attenuation: m.Reflectance,
	}

}

//Attenuation needs to change to fit specific environments
type ShadowCatcher struct {
	Attenuation float64
}

func (s ShadowCatcher) Scatter(inray Ray, Normal m.Vec3, env Environment, rander *rand.Rand) ScatterInfo {
	//The attenuation is kinda weird and hacky
	return ScatterInfo{
		NewRay:      true,
		NewRayDir:   Normal.Add(RandomUnitVec3(rander)),
		Color:       env.At(inray.Dir),
		Attenuation: s.Attenuation,
	}
}

type Emmisive struct {
	Albedo m.Vec3
	//Strength float64
}

func (e Emmisive) Scatter(inray Ray, Normal m.Vec3, env Environment, rander *rand.Rand) ScatterInfo {

	return ScatterInfo{
		NewRay:      false,
		Attenuation: 0,
		NewRayDir:   [3]float64{},
		Color:       e.Albedo,
	}
}
