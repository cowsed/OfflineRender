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

	_ "net/http/pprof"
)

//Image and Renderring Settings
const (
	BlockSize      = 64
	ImageWidth     = 1920 //1920 / 20
	ImageHeight    = 1080 //1080 / 20
	NumFrames      = 1
	SamplePerPixel = 4
)

//Render Settings
const (
	MaxDepth = 48
	Gamma    = .5
)

func main() {
	//Consistency across runs for easier debugging
	rand.Seed(0)

	l, _ := os.Create("CPU2.pprof")
	pprof.StartCPUProfile(l)

	env := &HDRIEnv{
		Filename: "TestResources/round_platform_4k.hdr",
		Rotation: .2,
		image:    &image.RGBA{},
	}
	env.LoadImg()

	for Frame := 0; Frame < NumFrames; Frame++ {
		initTime := time.Now()
		fmt.Printf("Beginning Frame %d of %d...", Frame, NumFrames)

		//Create Image
		img := image.NewRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))

		//Setup rendering things
		var MainCam Camera = Camera{
			Pos:         m.Vec3{0, -1.5, 0},
			Rot:         m.Vec3{0, 0, 0},
			Aspect:      float64(ImageWidth) / float64(ImageHeight),
			FocalLength: 1.0,
		}
		MainCam.Init()

		//Scene Components
		//===========================
		//Initialize Materials
		//====
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
				Fuzziness:   0,
			},
			//Bronze
			Diffuse{
				Albedo:      m.Vec3{1, .5, .5},
				Attenuation: 1,
				//Fuzziness:   .7,
			},
			//Floor
			Diffuse{
				Albedo:      m.Vec3{.9, .9, .9},
				Attenuation: .8,
			},
			//Floor 2
			ShadowCatcher{
				Attenuation: 1,
			},
		}
		//Initialize Geometry
		//====
		model1 := CreateModelFromSTL("TestResources/cube.stl", m.Vec3{-0, -1, 2.2}, 2)
		model1.Setup()
		fmt.Println(model1.meshBvh.aabb)

		intersectors := []Intersector{
			Sphere{
				Center:        m.Vec3{-1.75, -1, 3},
				Radius:        1,
				MaterialIndex: 0,
			},
			//Sphere{
			//	Center:        m.Vec3{0, .5, 2.4},
			//	Radius:        .5,
			//	MaterialIndex: 1,
			//},
			Sphere{
				Center:        m.Vec3{1.75, -1, 3},
				Radius:        1,
				MaterialIndex: 0,
			},
			Sphere{
				Center:        m.Vec3{0, 20000, 0},
				Radius:        20000,
				MaterialIndex: 4,
			},
			model1,
		}

		//Create the actual scene
		MainScene := Scene{
			Env:       env,
			Cam:       MainCam,
			Geometry:  intersectors,
			Materials: materials,
		}
		//Create BVH
		//MainScene.CreateBVH()

		//Render the scene
		MakeImage(img, ImageWidth, ImageHeight, &MainScene)

		//Save the Image
		f, err := os.Create(fmt.Sprintf("Outputs/out%d.png", Frame))
		check(err)
		err = png.Encode(f, img)
		check(err)
		err = f.Close()
		check(err)

		//Timing
		fmt.Println("Finished in", time.Since(initTime))
	}
	//Profiling
	pprof.StopCPUProfile()
	l.Close()

}
