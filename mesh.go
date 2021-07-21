package main

import (
	m "github.com/go-gl/mathgl/mgl64"
)

type Model struct {
	MaterialIndex int
	Position      m.Vec3

	Faces []Face
	//Vertices      []m.Vec3

	meshBvh MeshBVH
}

func (mod Model) MakeAABB() AABB {
	return mod.meshBvh.aabb
}

func (mod *Model) Setup() {

	//This only seems to work at 0(x) for some odd reason
	splitAxis := 2
	mod.meshBvh.axis = splitAxis
	mod.meshBvh.aabb.Min = mod.meshBvh.aabb.Min.Add(mod.Position)
	mod.meshBvh.aabb.Max = mod.meshBvh.aabb.Max.Add(mod.Position)

	mid := (mod.meshBvh.aabb.Max.Add(mod.meshBvh.aabb.Min)).Mul(.5)

	left := SubMesh{
		model: mod,
		Faces: make([]int, 0, len(mod.Faces)/2), //Just a wild approximation as to how many in each to try to preserve memory coherency
	}
	right := SubMesh{
		model: mod,
		Faces: make([]int, 0, len(mod.Faces)/2),
	}

	//For now just split along splitaxis
	for i := range mod.Faces {
		inLeft := 0
		for _, v := range mod.Faces[i].vs {
			if v[splitAxis] <= mid[splitAxis] {
				inLeft++
			}
		}
		//Put in left
		if inLeft == 3 {
			left.Faces = append(left.Faces, i)
		} else if inLeft == 0 { //Put in Right
			right.Faces = append(right.Faces, i)
		} else { //Put in both
			left.Faces = append(left.Faces, i)
			right.Faces = append(right.Faces, i)
		}
		if mod.Faces[i].MaterialIndex != 2 {
			println("Bad")
		}
	}

	mod.meshBvh.parts[0] = left
	mod.meshBvh.parts[1] = right

}

func (mod Model) Intersect(r Ray) (float64, m.Vec3, Intersector) {
	return mod.meshBvh.Intersect(r)
}
func (mod Model) Material() int {
	return mod.MaterialIndex
}

type SubMesh struct {
	model *Model
	Faces []int //Indices of the tris that are part of this submesh
}

//TODO
func (sub SubMesh) MakeAABB() AABB {
	return AABB{}
}

func (sub SubMesh) Intersect(r Ray) (float64, m.Vec3, Intersector) {
	//Intersect tris
	var intersectFace Intersector
	shortestDist := MAXDIST
	intersectNormal := m.Vec3{}
	for FaceIndexIndex := range sub.Faces {
		t, N, possibleIntersector := sub.model.Faces[sub.Faces[FaceIndexIndex]].Intersect(r)
		if t == -1 {
			continue
		}
		if t < shortestDist {
			shortestDist = t
			intersectNormal = N
			intersectFace = possibleIntersector
		}
	}
	//If it still hasnt hit anything
	if shortestDist == MAXDIST {
		return -1, m.Vec3{}, nil
	}
	//if intersectFace.Material() == 2 {
	//	fmt.Println("m", intersectFace.Material())
	//}
	return shortestDist, intersectNormal, intersectFace

}
