package main

import (
	"image"
	"image/draw"
	"math"
	"os"

	_ "image/jpeg"
	_ "image/png"

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

type Environment interface {
	At(Dir m.Vec3) m.Vec3
}

type SimpleEnv struct {
	col1, col2 m.Vec3
}

func (s SimpleEnv) At(Dir m.Vec3) m.Vec3 {
	return lerpColor(s.col1, s.col2, Dir.Normalize().Y()/2+.5)
}

type HDRIEnv struct {
	Filename string
	Rotation float64 //0-1 rotation to rotate the environment
	image    *image.RGBA
}

func (h *HDRIEnv) LoadImg() {
	f, err := os.Open(h.Filename)
	check(err)
	src, _, err := image.Decode(f)
	check(err)
	b := src.Bounds()
	h.image = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(h.image, h.image.Bounds(), src, b.Min, draw.Src)

}

func (h *HDRIEnv) At(Dir m.Vec3) m.Vec3 {
	Dir = Dir.Normalize()
	u := (math.Atan2(Dir[0], Dir[2]) + math.Pi) / (math.Pi * 2)
	v := Dir[1]/2 + .5
	u += h.Rotation
	u *= float64(h.image.Bounds().Dx())
	u = float64(int(u) % h.image.Bounds().Dx())
	v *= float64(h.image.Bounds().Dy())
	c := h.image.RGBAAt(int(u), int(v))
	return RGBA2V3(c)
}

func IntersectScene(r Ray, scene *Scene, MinDist float64) (float64, m.Vec3, m.Vec3, *Intersector) {

	var intersector *Intersector = nil
	var intersectNormal m.Vec3
	var shortestDist float64 = MAXDIST
	for i := range scene.Geometry {
		//scene[i]
		t, N := scene.Geometry[i].Intersect(r)
		if t < shortestDist && t >= MinDist {
			intersector = &scene.Geometry[i]

			intersectNormal = N

			shortestDist = t
		}
	}
	return shortestDist, r.At(shortestDist), intersectNormal, intersector

}
