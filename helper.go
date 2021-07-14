package main

import (
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

//Vector to color functions
func V22RGBA(v mgl64.Vec2) color.RGBA {

	return color.RGBA{
		uint8(255 * math.Abs(v.X())),
		uint8(255 * math.Abs(v.Y())),
		0,
		255}
}
func V32RGBA(v mgl64.Vec3) color.RGBA {

	return color.RGBA{
		uint8(float64(255) * math.Abs(v.X())),
		uint8(float64(255) * math.Abs(v.Y())),
		uint8(float64(255) * math.Abs(v.Z())),
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
