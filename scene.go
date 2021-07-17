package main

import (
	"math"

	m "github.com/go-gl/mathgl/mgl64"
)

const MAXDIST float64 = 10000000000

type Scene struct {
	Cam       Camera
	Geometry  []Sphere
	Materials []Material
}

type Sphere struct {
	Center        m.Vec3
	Radius        float64
	MaterialIndex int
}

func IntersectScene(r Ray, scene *Scene, MinDist float64) (float64, m.Vec3, m.Vec3, *Sphere) {
	
	var intersector *Sphere = nil
	var intersectNormal m.Vec3
	var shortestDist float64 = MAXDIST
	for i := range scene.Geometry {
		//scene[i]
		t := IntersectSphere(scene.Geometry[i].Center, scene.Geometry[i].Radius, r)
		if t < shortestDist && t >= MinDist {
			intersector = &scene.Geometry[i]

			intersectNormal = (r.At(t).Sub(scene.Geometry[i].Center)).Normalize()

			shortestDist = t
		}
	}
	return shortestDist, r.At(shortestDist), intersectNormal, intersector

}

func IntersectSphere(center m.Vec3, radius float64, r Ray) float64 {
	var oc = r.Origin.Sub(center)
	a := r.Dir.LenSqr()
	half_b := oc.Dot(r.Dir)
	c := oc.LenSqr() - radius*radius
	discriminant := half_b*half_b - a*c

	if discriminant < 0 {
		return -1.0
	} else {
		return (-half_b - math.Sqrt(discriminant)) / a
	}

}
