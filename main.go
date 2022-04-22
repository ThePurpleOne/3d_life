package main

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

func main() {

	var CUBE_SIZE float32
	CUBE_SIZE = 1
	SIZE_WORLD := 10

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam_pos := float32(SIZE_WORLD * int(CUBE_SIZE) / 2)
	cam.SetPosition( cam_pos, cam_pos, 15)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create a blue torus and add it to the scene
	//geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	geom := geometry.NewCube(CUBE_SIZE)
	mat := material.NewStandard(math32.NewColor("red"))

	world1 := create_world(geom, mat, 10, 10, 0)

	world1.show(scene)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{0.5, 0.5, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 10.0)
	light_pos := float32(SIZE_WORLD * int(CUBE_SIZE) / 2)

	pointLight.SetPosition(light_pos, light_pos, 5)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	//scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(1.0, 1.0, 1.0, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}