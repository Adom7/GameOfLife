package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
)

const scale int = 8
const width = 400
const height = 400

var grid = [width][height]uint8{}
var buffer = [width][height]uint8{}
var count int = 0
var gameRunning = false

func update() error {
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			buffer[x][y] = 0
			n := grid[x-1][y-1] + grid[x-1][y+0] + grid[x-1][y+1] + grid[x+0][y-1] + grid[x+0][y+1] + grid[x+1][y-1] + grid[x+1][y+0] + grid[x+1][y+1]

			if grid[x][y] == 0 && n == 3 {
				buffer[x][y] = 1
			} else if n < 2 || n > 3 {
				buffer[x][y] = 0
			} else {
				buffer[x][y] = grid[x][y]
			}
		}
	}

	temp := buffer
	buffer = grid
	grid = temp
	return nil
}

func display(window *ebiten.Image) {
	window.Fill(color.Black)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			for i := 0; i < scale; i++ {
				for j := 0; j < scale; j++ {
					if grid[x][y] == 1 {
						window.Set(x*scale+i, y*scale+j, color.RGBA{R: 255})
					}
				}
			}
		}
	}
}

func frame(window *ebiten.Image) error {
	handleMouseInput()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		toggleGameState()
	}

	count++
	var err error = nil
	if count == 50 && gameRunning {
		err = update()
		count = 0
	}
	if !ebiten.IsDrawingSkipped() {
		display(window)
	}

	return err
}

func handleMouseInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		gridX, gridY := x/scale, y/scale
		if gridX >= 0 && gridX < width && gridY >= 0 && gridY < height {
			grid[gridX][gridY] = 1
		}
	}
}

func toggleGameState() {
	gameRunning = !gameRunning
	if gameRunning {
		count = 49
	}
	log.Printf("Game Running: %v\n", gameRunning)
}

func main() {
	if err := ebiten.Run(frame, width, height, 2, "Game of Life"); err != nil {
		log.Fatal(err)
	}
}
