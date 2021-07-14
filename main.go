package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	m "github.com/go-gl/mathgl/mgl64"
	//	"image/color"

	_ "net/http/pprof"
)

//Renderring Settings
const (
	BlockSize   = 256
	ImageWidth  = 1920
	ImageHeight = 1080
	NumFrames   = 1
	//MSAA           = 1
	SamplePerPixel = 60
)

//Render Settings
const (
	MaxDepth = 60
	Gamma    = 1
)

//uv is [-1,1] and  [-1,1]
func main() {
	//Consistency
	rand.Seed(0)

	l, _ := os.Create("CPU2.pprof")
	pprof.StartCPUProfile(l)

	for Frame := 0; Frame < NumFrames; Frame++ {
		initTime := time.Now()
		fmt.Printf("Beginning Frame %d...", Frame)
		Time := float64(Frame) / float64(NumFrames)
		img := image.NewRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))

		//Setup rendering things
		var MainCam Camera = Camera{
			Pos:         m.Vec3{0, -.8, 0},
			Rot:         m.Vec3{0, 0, 0},
			Aspect:      float64(ImageWidth) / float64(ImageHeight),
			FocalLength: 1.0,
		}
		MainCam.Init()

		materials := []Material{
			{
				Color: m.Vec3{.5, 0, 1},
			},
			{
				Color: m.Vec3{.6, 0, .2},
			},
			{
				Color: m.Vec3{.1, .6, .2},
			},
		}
		spheres := []Sphere{
			{
				Center:        m.Vec3{-1, 0, 3 - math.Sin(2*math.Pi*Time)/2},
				Radius:        1,
				MaterialIndex: 0,
			},
			{
				Center:        m.Vec3{math.Sin(Time * 2 * math.Pi), -2, 3},
				Radius:        1,
				MaterialIndex: 1,
			},
			{
				Center:        m.Vec3{1, 0, 3 + math.Sin(2*math.Pi*Time)/2},
				Radius:        1,
				MaterialIndex: 2,
			},
			{
				Center:        m.Vec3{0, 2000, 0},
				Radius:        1999,
				MaterialIndex: 1,
			},
		}
		MainScene := Scene{
			Cam:       MainCam,
			Geometry:  spheres,
			Materials: materials,
		}
		MakeImage(img, ImageWidth, ImageHeight, &MainScene)

		f, err := os.Create(fmt.Sprintf("Outputs/out%d.png", Frame))
		check(err)

		err = png.Encode(f, img)
		check(err)

		err = f.Close()
		check(err)
		fmt.Println("Finished in", time.Since(initTime))
	}
	pprof.StopCPUProfile()
	l.Close()
}

//Assorted helper functions
func check(err error) {
	if err != nil {
		panic(err)
	}
}
func minI(a, b int) int {
	if a < b {
		return a
	}
	return a
}
