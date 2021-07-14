package main

import (
	"image"
	"math"
	"sync"

	m "github.com/go-gl/mathgl/mgl64"
)

func MakeImage(img *image.RGBA, width, height int, scene *Scene) {

	//Do the rendering
	var wg sync.WaitGroup
	for r := 0; r < height; r += BlockSize {
		for c := 0; c < width; c += BlockSize {
			wg.Add(1)
			go func(r1, c1 int) {
				FillBlock(scene, c1, r1, width, height, img)
				wg.Done()
				//fmt.Printf("Finished Block %d,%d\n", r1, c1)
			}(r, c)

		}

	}

	wg.Wait()
}

func FillBlock(scene *Scene, startx, starty, width, height int, img *image.RGBA) {
	for x := startx; x < minI(startx+BlockSize, width); x++ {
		for y := starty; y < minI(starty+BlockSize, height); y++ {
			c := MakePixel(scene, x, y, width, height)

			img.SetRGBA(x, y, V32RGBA(c))

		}
	}
}

func MakePixel(scene *Scene, x, y, width, height int) m.Vec3 {

	uv := m.Vec2{}
	uv[0] = (float64(x) / float64(width))
	uv[1] = float64(y) / float64(height)

	c := lerpColor(m.Vec3{.5, .7, 1.0}, m.Vec3{1.0, 1.0, 1.0}, uv.Y())
	//uv = uv.Mul(2).Sub(mgl64.Vec2{1, 1})

	var r = scene.Cam.GetRay(uv)

	c = colorRay(r, scene, MaxDepth)
	c = m.Vec3{
		math.Pow(c.X(), Gamma),
		math.Pow(c.Y(), Gamma),
		math.Pow(c.Z(), Gamma),
	}

	return c
}

func colorRay(r Ray, scene *Scene, depth int) m.Vec3 {
	if depth == 0 {
		return m.Vec3{}
	}
	var c m.Vec3

	t, p, N, intersector := IntersectScene(r, scene)
	//Hit Nothing, sky color
	if intersector == nil {
		c = lerpColor(m.Vec3{.5, .7, 1.0}, m.Vec3{1.0, 1.0, 1.0}, r.Dir.Normalize().Y())
		return c
	}
	if t >= 0 {

		c = scene.Materials[intersector.MaterialIndex].Color

		newDir := reflect(r.Dir, N)
		newR := Ray{
			Origin: p,
			Dir:    newDir,
		}
		c2 := colorRay(newR, scene, depth-1)
		c = c.Mul(.5).Add(c2.Mul(.5))
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
