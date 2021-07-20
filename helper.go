package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

func Mul3x3(a, b mgl64.Vec3) mgl64.Vec3 {
	return mgl64.Vec3{
		a[0] * b[0],
		a[1] * b[1],
		a[2] * b[2],
	}
}

func RandomVec3(rander *rand.Rand) mgl64.Vec3 {
	return mgl64.Vec3{
		rander.Float64()*2 - 1,
		rander.Float64()*2 - 1,
		rander.Float64()*2 - 1,
	}
}

func RandomVec3InUnitSphere(rander *rand.Rand) mgl64.Vec3 {
	for {
		var p = RandomVec3(rander)
		if p.LenSqr() >= 1 {
			continue
		}
		return p
	}
}
func RandomUnitVec3(rander *rand.Rand) mgl64.Vec3 {
	return RandomVec3InUnitSphere(rander).Normalize()
}

func minI8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return a
}
func minI(a, b int) int {
	if a < b {
		return a
	}
	return a
}

//Vector to color functions
func V22RGBA(v mgl64.Vec2) color.RGBA {

	return color.RGBA{
		minI8(uint8(255*math.Abs(v.X())), 255),
		minI8(uint8(255*math.Abs(v.Y())), 255),
		0,
		255}
}
func V32RGBA(v mgl64.Vec3) color.RGBA {

	return color.RGBA{
		minI8(uint8(float64(255)*math.Abs(v.X())), 255),
		minI8(uint8(float64(255)*math.Abs(v.Y())), 255),
		minI8(uint8(float64(255)*math.Abs(v.Z())), 255),
		255}
}

func RGBA2V3(c color.RGBA) mgl64.Vec3 {

	return mgl64.Vec3{
		float64(c.R) / 255,
		float64(c.G) / 255,
		float64(c.B) / 255,
	}
}

func V42RGBA(v mgl64.Vec4) color.RGBA {

	return color.RGBA{
		uint8(float64(255) * math.Abs(v.X())),
		uint8(float64(255) * math.Abs(v.Y())),
		uint8(float64(255) * math.Abs(v.Z())),
		uint8(float64(255) * math.Abs(v.W())),
	}
}

func Abs2(v mgl64.Vec2) (o mgl64.Vec2) {
	o[0] = mgl64.Abs(v[0])
	o[1] = mgl64.Abs(v[1])
	return o
}
