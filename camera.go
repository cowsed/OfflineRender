package main

import (
	m "github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	Pos                           m.Vec3
	Rot                           m.Vec3
	Aspect                        float64
	FocalLength                   float64
	viewportHeight, viewportWidth float64
	horizontal, vertical          m.Vec3
	lowerLeftCorner               m.Vec3

	rotmtx m.Mat3
}

func (c *Camera) Init() {
	c.viewportHeight = 2.0
	c.viewportWidth = c.Aspect * c.viewportHeight
	c.horizontal = m.Vec3{c.viewportWidth, 0, 0}
	c.vertical = m.Vec3{0, c.viewportHeight, 0}

	c.lowerLeftCorner = c.Pos.Sub(c.horizontal.Mul(.5).Add(c.vertical.Mul(.5).Sub(m.Vec3{0, 0, c.FocalLength})))
	c.rotmtx = m.Rotate3DZ(c.Rot.X())

}

//x,y are -1,1 normalized coords
func (c Camera) GetRay(uv m.Vec2) Ray {

	return Ray{
		Origin: c.Pos,
		Dir:    c.rotmtx.Mul3x1(c.lowerLeftCorner.Add(c.horizontal.Mul(uv.X())).Add(c.vertical.Mul(uv.Y())).Sub(c.Pos)),
	}
}

type Ray struct {
	Origin m.Vec3
	Dir    m.Vec3
}

func (r Ray) At(t float64) m.Vec3 {
	return r.Origin.Add(r.Dir.Mul(t))
}
