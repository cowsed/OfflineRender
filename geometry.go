package main

import (
	"math"

	m "github.com/go-gl/mathgl/mgl64"
)

type Intersector interface {
	Intersect(r Ray) (float64, m.Vec3)
	Material() int
}

type Sphere struct {
	Center        m.Vec3
	Radius        float64
	MaterialIndex int
}

func (s Sphere) Intersect(r Ray) (float64, m.Vec3) {
	var oc = r.Origin.Sub(s.Center)
	a := r.Dir.LenSqr()
	half_b := oc.Dot(r.Dir)
	c := oc.LenSqr() - s.Radius*s.Radius
	discriminant := half_b*half_b - a*c

	if discriminant < 0 {
		return -1.0, m.Vec3{}
	} else {
		t := (-half_b - math.Sqrt(discriminant)) / a
		pos := r.At(t)
		nor := pos.Sub(s.Center).Normalize()

		return t, nor
	}

}

func (s Sphere) Material() int {
	return s.MaterialIndex
}

type Face struct {
	vs            [3]m.Vec3
	n             m.Vec3
	MaterialIndex int
}

func (face Face) Intersect(r Ray) (float64, m.Vec3) {

	var EPSILON float64 = 0.0000001
	vertex0 := face.vs[0]
	vertex1 := face.vs[1]
	vertex2 := face.vs[2]
	var edge1, edge2, h, s, q m.Vec3
	var a, f, u, v float64
	edge1 = vertex1.Sub(vertex0)
	edge2 = vertex2.Sub(vertex0)
	h = r.Dir.Cross(edge2)
	a = edge1.Dot(h)
	if a > -EPSILON && a < EPSILON {
		return -1, m.Vec3{} // This ray is parallel to this triangle.
	}
	f = 1.0 / a
	s = r.Origin.Sub(vertex0)
	u = f * s.Dot(h)
	if u < 0.0 || u > 1.0 {
		return -1, m.Vec3{}
	}
	q = s.Cross(edge1)
	v = f * r.Dir.Dot(q)
	if v < 0.0 || u+v > 1.0 {
		return -1, m.Vec3{}
	}
	// At this stage we can compute t to find out where the intersection point is on the line.
	t := f * edge2.Dot(q)
	if t > EPSILON { // ray intersection

		return t, face.n
	} else { // This means that there is a line intersection but not a ray intersection.
		return -1, m.Vec3{}
	}

}
func (face Face) Material() int {
	return face.MaterialIndex
}

type Model struct {
	MaterialIndex int
	Position      m.Vec3

	Faces []Face
	//Vertices      []m.Vec3

	bvh MeshBVH
}

func (mod *Model) Setup() {

	//This only seems to work at 0(x) for some odd reason
	splitAxis := 2
	mod.bvh.axis = splitAxis
	mod.bvh.aabb.Min = mod.bvh.aabb.Min.Add(mod.Position)
	mod.bvh.aabb.Max = mod.bvh.aabb.Max.Add(mod.Position)

	mid := (mod.bvh.aabb.Max.Add(mod.bvh.aabb.Min)).Mul(.5)

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
	}

	mod.bvh.parts[0] = left
	mod.bvh.parts[1] = right

}

type SubMesh struct {
	model *Model
	Faces []int //Indices of the tris that are part of this submesh
}

func (sub SubMesh) Material() int {
	return sub.model.MaterialIndex
}
func (sub SubMesh) Intersect(r Ray) (float64, m.Vec3) {
	//Intersect tris
	shortestDist := MAXDIST
	intersectNormal := m.Vec3{}
	for _, FaceIndex := range sub.Faces {
		t, N := sub.model.Faces[FaceIndex].Intersect(r)
		if t == -1 {
			continue
		}
		if t < shortestDist {
			shortestDist = t
			intersectNormal = N
		}
	}
	//If it still hasnt hit anything
	if shortestDist == MAXDIST {
		return -1, m.Vec3{}
	}
	return shortestDist, intersectNormal

}

type MeshBVH struct {
	model  *Model
	parts  [2]Intersector
	aabb   AABB    //AABB for both halves
	middle float64 //Center axis to split aabb over
	axis   int     //0-2
}

func (b MeshBVH) Intersect(r Ray) (float64, m.Vec3) {

	var tA, tB float64

	var nA, nB m.Vec3

	leftMax := b.aabb.Max
	leftMax[b.axis] = b.middle
	rightMin := b.aabb.Min
	rightMin[b.axis] = b.middle

	leftAABB := AABB{
		Min: b.aabb.Min,
		Max: leftMax,
	}
	rightAABB := AABB{
		Min: rightMin,
		Max: b.aabb.Max,
	}

	leftHit := leftAABB.Intersects(r)
	rightHit := rightAABB.Intersects(r)
	//Both missed
	if !leftHit && !rightHit {
		//return -1, m.Vec3{}
	}

	//Left Hit missed
	if !leftHit {
		//Check Right Mesh
		tB, nB = b.parts[1].Intersect(r)
		return tB, nB
	} else if !rightHit { //Right Hit missed
		//Check Left Mesh
		tA, nA = b.parts[0].Intersect(r)
		return tA, nA
	}

	//Both AABBS hit. must check mesh
	tA, nA = b.parts[0].Intersect(r)
	tB, nB = b.parts[1].Intersect(r)

	if tA == -1 && tB == -1 {
		//Both Miss
		return -1, m.Vec3{}
	} else if tA == -1 {
		//Case left misses inside
		return tB, nB
	} else if tB == -1 { //case Right Misses inside
		return tA, nA
	}

	//Case both hit something and a was closer
	if tA < tB {
		return tA, nA
	}
	//Case both hit something and b was closer

	return tB, nB
}

func (mod Model) Intersect(r Ray) (float64, m.Vec3) {
	/*
		if !mod.aabb.Intersects(r) {
			return -1, m.Vec3{}
		}
		//Intersect tris
		shortestDist := MAXDIST
		intersectNormal := m.Vec3{}
		for i := range mod.Faces {
			t, N := mod.Faces[i].Intersect(r)
			if t == -1 {
				continue
			}
			if t < shortestDist {
				shortestDist = t
				intersectNormal = N
			}
		}
		//If it still hasnt hit anything
		if shortestDist == MAXDIST {
			return -1, m.Vec3{}
		}
		return shortestDist, intersectNormal
	*/
	return mod.bvh.Intersect(r)
}
func (mod Model) Material() int {
	return mod.MaterialIndex
}

type AABB struct {
	Min m.Vec3
	Max m.Vec3
}

func (m AABB) Intersects(r Ray) bool {
	for a := 0; a < 3; a++ {
		t0 := fmin((m.Min[a]-r.Origin[a])/r.Dir[a],
			(m.Max[a]-r.Origin[a])/r.Dir[a])
		t1 := fmax((m.Min[a]-r.Origin[a])/r.Dir[a],
			(m.Max[a]-r.Origin[a])/r.Dir[a])
		t_min := fmax(t0, MINDIST)
		t_max := fmin(t1, MAXDIST)
		if t_max <= t_min {
			return false
		}
	}
	return true

}

func fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
