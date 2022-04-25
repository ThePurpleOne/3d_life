package main

import (
	"math/rand"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type cell struct {
	mesh *graphic.Mesh
	pos math32.Vector3
	size int 
	active bool
}

type world struct{
	cells ([][]cell)
	w, h, l int
}


func create_world(shape *geometry.Geometry ,mat *material.Standard , w, h, l int) world{
	cells := make([][]cell, 0)
	for i := 0; i < w; i++ {
		line := make([]cell, 0)
		for j := 0; j < h; j++ {
			mesh := graphic.NewMesh(shape, mat)
			mesh.SetPosition(float32(i) + (float32(i) * 0.1), float32(j) + (float32(j) * 0.1) , 0)
			line = append(line, cell{mesh, math32.Vector3{float32(i), float32(j), 0.0}, 1, (rand.Intn(2) == 0)})
			//line = append(line, cell{mesh, math32.Vector3{float32(i), float32(j), 0.0}, 1, true}); //(i % 2 == 0)})
		}
		cells = append(cells, line)
	}
	return world{cells, w, h, l}
}

func (w world) add(scene* core.Node){
	for i := 0; i < w.w; i++ {
		for j := 0; j < w.h; j++ {
			w.cells[i][j].add(scene)
		}
	}
}

// Brief : ADD THE CELL TO THE SCENE TO BE DRAWN
// Args : [core.Node] scene to draw to
func (c cell) add(scene* core.Node){
	if c.is_active(){
		scene.Add(c.mesh)
	}
}

func (w world) show(){
	for i := 0; i < w.w; i++ {
		for j := 0; j < w.h; j++ {
			w.cells[i][j].show()
		}
	}
}

// Brief : ADD THE CELL TO THE SCENE TO BE DRAWN
// Args : [core.Node] scene to draw to
func (c* cell) show(){
	if c.is_active(){
		c.mesh.SetVisible(true)
	}else{
		c.mesh.SetVisible(false)	
	}
}

// Brief : Check wether a cell is active or not
// Returns : True if active
func (c cell) is_active() bool{
	return c.active
}

// Brief : Count the number of active cells in a world
// Returns : Number of active cells around a cell
func count_neighbors(w world, c* cell) int{
	neig_count := 0
	resolution := w.w / c.size

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			col := ((int(c.pos.X) + i  + resolution) % resolution)
			row := ((int(c.pos.Y) + j  + resolution) % resolution)
			neig_count += bool2int(w.cells[col][row].is_active());
		}
	}

	// Remove self if active
	neig_count -= bool2int(w.cells[int(c.pos.X)][int(c.pos.Y)].is_active())
	return neig_count
}

func (w* world) process_world(world_copy world){
	
	neigh_count := 0
	
	removed, added, stayed := 0, 0, 0
	for i := 0; i < w.w; i++ {
		for j := 0; j < w.h; j++ {
			current := &w.cells[i][j]

			neigh_count = count_neighbors(world_copy, current)

			// MAKE CHANGES ON THE CURRENT CELL
			if current.is_active() && neigh_count < 2{								// ! UNDERNPOPULATION
				current.active = false; removed++;
			}else if current.is_active() && neigh_count > 3{ 						// ! OVERPOPULATION
				current.active = false ; removed++;
			}else if !current.is_active() && neigh_count == 3{						// ! REPRODUCTION
				current.active = true ; added++;
			}else if current.is_active() && (neigh_count == 2 || neigh_count == 3){	// ! STAY ALIVE
				current.active = true ; stayed++;
			}
		}
	}
	//fmt.Printf("Removed: %d Added: %d Stayed : %d", removed, added, stayed)
}

// -----------------------------------------------------------------------------
// ----------------------------------- UTILS ----------------------------------- 
// -----------------------------------------------------------------------------
func bool2int(b bool) int {
	if b {return 1}
	return 0
} 
