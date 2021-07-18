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

type Model struct {
	MaterialIndex int
	Position      m.Vec3
	aabb          AABB
	Faces         []Face
	//Vertices      []m.Vec3
}
type Face struct {
	vs [3]m.Vec3
	n  m.Vec3
}

func (face Face) Intersects(r Ray) (float64, m.Vec3) {

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

func (mod *Model) Setup() {
	for i := range mod.Faces {
		mod.Faces[i].vs[0] = mod.Faces[i].vs[0].Add(mod.Position)
		mod.Faces[i].vs[1] = mod.Faces[i].vs[1].Add(mod.Position)
		mod.Faces[i].vs[2] = mod.Faces[i].vs[2].Add(mod.Position)

	}
	mod.aabb.Min = mod.aabb.Min.Add(mod.Position)
	mod.aabb.Max = mod.aabb.Max.Add(mod.Position)

}

func (mod Model) Intersect(r Ray) (float64, m.Vec3) {
	if !mod.aabb.Intersects(r) {
		return -1, m.Vec3{}
	}
	//Intersect tris
	shortestDist := MAXDIST
	intersectNormal := m.Vec3{}
	for i := range mod.Faces {
		t, N := mod.Faces[i].Intersects(r)
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
