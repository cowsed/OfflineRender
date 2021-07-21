package main

import (
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/mdouchement/hdr/codec/rgbe"

	m "github.com/go-gl/mathgl/mgl64"
)

const MAXDIST float64 = 10000000000
const MINDIST float64 = 1e-8

type Scene struct {
	Env       Environment
	Cam       Camera
	Geometry  []Intersector
	Materials []Material
}

func (scene *Scene) Intersect(r Ray, MinDist float64) (float64, m.Vec3, m.Vec3, Intersector) {
	//Traverse BVH to find Geometry intersections

	var intersector Intersector = nil
	var intersectNormal m.Vec3
	var shortestDist float64 = MAXDIST
	for i := range scene.Geometry {
		t, N, tempIntersector := scene.Geometry[i].Intersect(r)
		if t < shortestDist && t >= MinDist {
			intersector = tempIntersector //scene.Geometry[i] //possibleIntersector

			intersectNormal = N

			shortestDist = t
		}
	}
	return shortestDist, r.At(shortestDist), intersectNormal, intersector

}
