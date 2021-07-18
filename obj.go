package main

import (
	"fmt"

	"github.com/hschendel/stl"

	m "github.com/go-gl/mathgl/mgl64"
)

func CreateModelFromSTL(fname string) Model {
	solid, err := stl.ReadFile(fname)
	check(err)
	faces := make([]Face, len(solid.Triangles))
	aabb := AABB{
		Min: [3]float64{MAXDIST, MAXDIST, MAXDIST},
		Max: [3]float64{-MAXDIST, -MAXDIST, -MAXDIST},
	}
	for i := range solid.Triangles {
		vs := solid.Triangles[i].Vertices
		f := Face{}

		for j := 0; j < 3; j++ {

			v := m.Vec3{
				float64(vs[j][0]),
				float64(vs[j][1]),
				float64(vs[j][2]),
			}
			if v.X() < aabb.Min.X() {
				aabb.Min[0] = v.X()
			}
			if v.Y() < aabb.Min.Y() {
				aabb.Min[1] = v.Y()
			}
			if v.Z() < aabb.Min.Z() {
				aabb.Min[2] = v.Z()
			}
			if v.X() > aabb.Max.X() {
				aabb.Max[0] = v.X()
			}
			if v.Y() > aabb.Max.Y() {
				aabb.Max[1] = v.Y()
			}
			if v.Z() > aabb.Min.Z() {
				aabb.Max[2] = v.Z()
			}
			f.vs[j] = v
		}
		f.n = m.Vec3{
			float64(solid.Triangles[i].Normal[0]) * 1,
			float64(solid.Triangles[i].Normal[1]) * 1,
			float64(solid.Triangles[i].Normal[2]) * 1,
		}
		faces[i] = f
	}
	fmt.Println(aabb)
	return Model{
		MaterialIndex: 0,
		aabb:          aabb,
		Faces:         faces,
		//Vertices:      []mgl64.Vec3{},
	}

}
