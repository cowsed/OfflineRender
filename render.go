package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"sync"

	m "github.com/go-gl/mathgl/mgl64"
)

func MakeImage(img *image.RGBA, width, height int, scene *Scene) {

	//Do the rendering
	var wg sync.WaitGroup
	NumRoutines := 0
	for r := 0; r < height; r += BlockSize {
		for c := 0; c < width; c += BlockSize {
			wg.Add(1)
			NumRoutines++
			go func(r1, c1 int) {
				FillBlock(scene, c1, r1, width, height, img)
				wg.Done()
				//fmt.Printf("Finished Block %d,%d\n", r1, c1)
			}(r, c)

		}
	}
	fmt.Printf("Using %d goroutines \n", NumRoutines)
	wg.Wait()
}

func FillBlock(scene *Scene, startx, starty, width, height int, img *image.RGBA) {
	//New Source of random numbers to avoid locking times
	rander := rand.New(rand.NewSource(int64(startx*width + starty)))
	for x := startx; x < minI(startx+BlockSize, width); x++ {
		for y := starty; y < minI(starty+BlockSize, height); y++ {
			c := MakePixel(scene, x, y, width, height, rander)

			img.SetRGBA(x, y, V32RGBA(c))

		}
	}
}

func MakePixel(scene *Scene, x, y, width, height int, rander *rand.Rand) m.Vec3 {

	fullColor := m.Vec3{}

	//for X := -.5; X < .5; X += 1 / float64(MSAA) {
	//	for Y := -.5; Y < .5; Y += 1 / float64(MSAA) {
	for sample := 0; sample < SamplePerPixel; sample++ {

		uv := m.Vec2{}
		uv[0] = (float64(x) + rander.Float64()) / float64(width)
		uv[1] = (float64(y) + rander.Float64()) / float64(height)

		c := lerpColor(m.Vec3{.5, .7, 1.0}, m.Vec3{1.0, 1.0, 1.0}, uv.Y())

		var r = scene.Cam.GetRay(uv)

		c = colorRay(r, scene, rander, MaxDepth)

		fullColor = fullColor.Add(c)
	}
	//		}
	//}
	//fullColor = fullColor.Mul(1.0 / float64(MSAA*MSAA))
	fullColor = fullColor.Mul(1.0 / SamplePerPixel)
	fullColor = m.Vec3{
		math.Pow(fullColor.X(), Gamma),
		math.Pow(fullColor.Y(), Gamma),
		math.Pow(fullColor.Z(), Gamma),
	}
	return fullColor
}

func colorRay(r Ray, scene *Scene, rander *rand.Rand, depth int) m.Vec3 {
	if depth <= 0 {
		return m.Vec3{}
	}
	var c m.Vec3

	t, p, N, intersector := IntersectScene(r, scene, 0.001)
	//Hit Nothing, sky color
	if intersector == nil {
		
		c = scene.Env.At(r.Dir)
		return c
	}
	//Color by material
	if t >= 0 {

		//newDir := reflect(r.Dir, N)
		SInfo := scene.Materials[(*intersector).Material()].Scatter(r, N, rander)

		newR := Ray{
			Origin: p,
			Dir:    SInfo.NewRayDir,
		}

		c = colorRay(newR, scene, rander, depth-1).Mul(SInfo.Attenuation)
		c = Mul3x3(SInfo.Color, c)
		//c = c.Mul(.5).Add(colorRay(newR, scene, rander, depth-1).Mul(.5)).Mul(.5)
	}
	return c

}

func lerpColor(a, b m.Vec3, amt float64) m.Vec3 {
	c := m.Vec3{}
	c[0] = b.X()*amt + (1-amt)*a.X()
	c[1] = b.Y()*amt + (1-amt)*a.Y()
	c[2] = b.Z()*amt + (1-amt)*a.Z()

	return c
}

func reflect(I, N m.Vec3) m.Vec3 {

	return I.Sub(N.Mul(2.0 * N.Dot(I)))
}
