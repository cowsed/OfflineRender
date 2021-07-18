package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	m "github.com/go-gl/mathgl/mgl64"
	//	"image/color"

	_ "net/http/pprof"
)

//Image and Renderring Settings
const (
	BlockSize      = 128
	ImageWidth     = 1920 / 2
	ImageHeight    = 1080 / 2
	NumFrames      = 1
	SamplePerPixel = 1400
)

//Render Settings
const (
	MaxDepth = 48
	Gamma    = .5
)

//uv is [-1,1] and  [-1,1]
func main() {
	//Consistency across runs for easier debugging
	rand.Seed(0)

	l, _ := os.Create("CPU2.pprof")
	pprof.StartCPUProfile(l)

	for Frame := 0; Frame < NumFrames; Frame++ {
		initTime := time.Now()
		fmt.Printf("Beginning Frame %d...", Frame)
		//Time := float64(Frame) / float64(NumFrames)
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
			//Silver
			Metal{
				Albedo:      m.Vec3{.5, 0.5, .5},
				Reflectance: 1,
				Fuzziness:   0,
			},
			//Gold
			Metal{
				Albedo:      m.Vec3{0.72, 0.53, 0.04},
				Reflectance: .6,
				Fuzziness:   .2,
			},
			//Bronze
			Diffuse{
				Albedo:      m.Vec3{1, 1, 1},
				Attenuation: 1,
				//Fuzziness:   .7,
			},
			//Floor
			Diffuse{
				Albedo:      m.Vec3{.9, .9, .9},
				Attenuation: .8,
			},
		}
		intersectors := []Intersector{
			Sphere{
				Center:        m.Vec3{-.3, 0, 3.7},
				Radius:        1,
				MaterialIndex: 0,
			},
			//Sphere{
			//	Center:        m.Vec3{0, .5, 2.4},
			//	Radius:        .5,
			//	MaterialIndex: 1,
			//},
			//Sphere{
			//	Center:        m.Vec3{1.55, 0, 3},
			//	Radius:        1,
			//	MaterialIndex: 2,
			//},
			Sphere{
				Center:        m.Vec3{0, 20000, 0},
				Radius:        19999,
				MaterialIndex: 3,
			},
		}

		rander := rand.New(rand.NewSource(1))
		for i := 0; i < 10; i++ {
			p := RandomVec3InUnitSphere(rander)
			p[1] = 0
			if p.Len() < .5 {
				i--
				continue
			}

			p = p.Mul(-2.2)
			p[1] = .8

			p = p.Add(m.Vec3{0, 0, 3.8})
			fmt.Println("NewSphere", p)
			ns := Sphere{
				Center:        p,
				Radius:        .2,
				MaterialIndex: len(materials),
			}
			materials = append(materials, Diffuse{
				Albedo:      RandomVec3InUnitSphere(rander),
				Attenuation: .9,
			})
			intersectors = append(intersectors, ns)
		}

		MainScene := Scene{
			Env:       SimpleEnv{m.Vec3{.2, .5, 1}, m.Vec3{1, 1, 1}},
			Cam:       MainCam,
			Geometry:  intersectors,
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

	/*
		f, err := os.Create("mem.pprof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	*/
}

//Assorted helper functions
func check(err error) {
	if err != nil {
		panic(err)
	}
}
