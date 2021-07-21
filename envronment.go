package main

import (
	"image"
	"image/draw"
	"math"
	"os"

	m "github.com/go-gl/mathgl/mgl64"
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/tmo"
)

type Environment interface {
	At(Dir m.Vec3) m.Vec3
}

type SimpleEnv struct {
	TopCol, BottomCol m.Vec3
}

func (s SimpleEnv) At(Dir m.Vec3) m.Vec3 {
	return lerpColor(s.TopCol, s.BottomCol, Dir.Normalize().Y()/2+.5)
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

	if hdrm, ok := src.(hdr.Image); ok {

		//t := tmo.NewLinear(hdrm)
		//t := tmo.NewLogarithmic(hdrm)
		//t := tmo.NewDefaultDrago03(hdrm)
		//t := tmo.NewDefaultDurand(hdrm)
		//t := tmo.NewDefaultCustomReinhard05(hdrm)
		t := tmo.NewDefaultReinhard05(hdrm)
		//t := tmo.NewDefaultICam06(hdrm)
		src = t.Perform()

	}
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
	//c.R = minI8(c.R, 255)
	//c.G = minI8(c.G, 255)
	//c.B = minI8(c.B, 255)
	//c.A = minI8(c.A, 255)
	return RGBA2V3(c)
}
