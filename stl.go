package main

import (
	"github.com/hschendel/stl"

	m "github.com/go-gl/mathgl/mgl64"
)

func CreateModelFromSTL(fname string, pos m.Vec3, mat int) Model {
	mod := Model{}
	mod.bvh = MeshBVH{}
	mod.bvh.aabb = AABB{
		Min: [3]float64{MAXDIST, MAXDIST, MAXDIST},
		Max: [3]float64{-MAXDIST, -MAXDIST, -MAXDIST},
	}

	solid, err := stl.ReadFile(fname)
	check(err)
	faces := make([]Face, len(solid.Triangles))

	for i := range solid.Triangles {
		vs := solid.Triangles[i].Vertices
		f := Face{}

		for j := 0; j < 3; j++ {

			v := m.Vec3{
				float64(vs[j][0]),
				float64(vs[j][1]),
				float64(vs[j][2]),
			}
			v = v.Add(pos)
			//AABB Checking
			{
				if v.X() < mod.bvh.aabb.Min.X() {
					mod.bvh.aabb.Min[0] = v.X()
				}
				if v.Y() < mod.bvh.aabb.Min.Y() {
					mod.bvh.aabb.Min[1] = v.Y()
				}
				if v.Z() < mod.bvh.aabb.Min.Z() {
					mod.bvh.aabb.Min[2] = v.Z()
				}
				if v.X() > mod.bvh.aabb.Max.X() {
					mod.bvh.aabb.Max[0] = v.X()
				}
				if v.Y() > mod.bvh.aabb.Max.Y() {
					mod.bvh.aabb.Max[1] = v.Y()
				}
				if v.Z() > mod.bvh.aabb.Min.Z() {
					mod.bvh.aabb.Max[2] = v.Z()
				}
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
	mod.MaterialIndex = mat
	mod.Faces = faces
	mod.bvh.model = &mod
	//Vertices:      []mgl64.Vec3{},
	return mod

}
