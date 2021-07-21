package main

import (
	"fmt"
	"sort"

	m "github.com/go-gl/mathgl/mgl64"
)

type BVH struct {
	Children [2]*BVH
	aabb     AABB        //Union of Children AABBs
	Geometry Intersector //Is nil unless this bvh is a leaf node
}

func (scene *Scene) MakeBVH() {
	//just one axis for now
	for axis := 0; axis <= 0; axis++ {
		//Sort objects by position on that axis or of bounding boxes
		sort.SliceStable(scene.Geometry,
			func(i, j int) bool {
				return scene.Geometry[i].GetAABB().Min[axis] < scene.Geometry[j].GetAABB().Min[axis]
			},
		)
		//Make into a tree
		scene.bvh = ConstructBVH(scene.Geometry, 0, len(scene.Geometry))
	}

}

func PrintBVH(b *BVH) {
	printBVH("", b)
}

func printBVH(prefix string, b *BVH) {
	if b.Geometry == nil {
		//Is a junction
		fmt.Println(prefix + "└+")
		printBVH(prefix+"  ", b.Children[1])
		printBVH(prefix+"  ", b.Children[0])

	} else {
		fmt.Println(prefix+"└─", b.Geometry.Name())
	}
}

func ConstructBVH(elements []Intersector, min, max int) BVH {
	//More than a leaf
	if max-min > 2 {
		newMid := min + ((max - min) / 2)
		A := ConstructBVH(elements, min, newMid)
		B := ConstructBVH(elements, newMid, max)

		return BVH{
			aabb:     A.aabb.Union(B.aabb),
			Children: [2]*BVH{&A, &B},
		}
	}
	//only 2 leaves
	if max-min == 2 {
		A := BVH{Geometry: elements[min]}
		B := BVH{Geometry: elements[min+1]}
		return BVH{
			aabb:     A.aabb.Union(B.aabb),
			Children: [2]*BVH{&A, &B},
		}
	}

	//branch with 1 leave
	if max-min == 1 {
		return BVH{
			Geometry: elements[min],
			aabb:     elements[min].GetAABB(),
		}
	}

	fmt.Println("Should never get here")
	return BVH{}

}

func (b BVH) Intersect(r Ray) (float64, m.Vec3) {

	return -1, m.Vec3{}
}

type MeshBVH struct {
	model  *Model
	parts  [2]SubMesh
	aabb   AABB    //AABB for both halves
	middle float64 //Center axis to split aabb over
	axis   int     //0-2
}

///The returning intersector here may need changing to return the face, rather than the model for texturing purposes
func (b MeshBVH) Intersect(r Ray) (float64, m.Vec3, Intersector) {

	var intersectA, intersectB Intersector
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
	//This is problematic. it makes some strange artifacts when it misses the inside but hits the box
	//if !leftHit && !rightHit {
	//	return -1, m.Vec3{}, nil
	//}

	//Left Hit missed
	if !leftHit {
		//Check Right Mesh
		tB, nB, intersectB = b.parts[1].Intersect(r)
		return tB, nB, intersectB
	} else if !rightHit { //Right Hit missed
		//Check Left Mesh
		tA, nA, intersectA = b.parts[0].Intersect(r)
		return tA, nA, intersectA
	}

	//Both AABBS hit. must check mesh
	tA, nA, intersectA = b.parts[0].Intersect(r)
	tB, nB, intersectB = b.parts[1].Intersect(r)

	if tA == -1 && tB == -1 {
		//Both Miss
		return -1, m.Vec3{}, nil
	} else if tA == -1 {
		//Case left misses inside
		return tB, nB, intersectB
	} else if tB == -1 { //case Right Misses inside
		return tA, nA, intersectA
	}

	//Case both hit something and a was closer

	if tA < tB {
		return tA, nA, intersectA
	}
	//Case both hit something and b was closer

	return tB, nB, intersectB
}
