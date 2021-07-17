package main

import m "github.com/go-gl/mathgl/mgl64"
import "math/rand"
type Material interface {
	Scatter(Normal m.Vec3, rander *rand.Rand) ScatterInfo
}
type Diffuse struct{
	Albedo m.Vec3
	Reflectance float64
}



func (d Diffuse) Scatter(Normal m.Vec3, rander *rand.Rand) ScatterInfo{
	return ScatterInfo{
		NewRay: true,
		NewRayDir: Normal.Add(RandomUnitVec3(rander)), 
		Color: d.Albedo,
		Attenuation: d.Reflectance,
	}
	
}

type ScatterInfo struct{
	NewRay bool
	Attenuation float64
	NewRayDir, Color m.Vec3
}