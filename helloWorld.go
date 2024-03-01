package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"time"
)

var scale int = 8

const width = 400
const height = 400

var grid = [width][height]uint8{}
var buffer = [width][height]uint8{}
var gameRunning = false
var generation int = 0
var livingCell int = 0
var lastUpdateTime time.Time
var timer time.Duration = 2000

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
	livingCell = 0

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			for i := 0; i < scale; i++ {
				for j := 0; j < scale; j++ {
					if grid[x][y] == 1 {
						window.Set(x*scale+i, y*scale+j, color.RGBA{R: 255})
						livingCell++
					}
				}
			}
		}
	}
}

func frame(window *ebiten.Image) error {
	if lastUpdateTime.IsZero() {
		lastUpdateTime = time.Now()
	}

	timeElapsed := time.Since(lastUpdateTime)
	desiredInterval := timer * time.Millisecond

	handleMouseInput()
	handleVerticalArrowsPress()
	handleHorizontalArrowsPress()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		toggleGameState()
	}

	var err error = nil

	if timeElapsed >= desiredInterval && gameRunning {
		err = update()
		log.Printf("Generation: %v\n", generation)
		log.Printf("Living cells: %v\n", livingCell/64)
		log.Printf("Interval: %v\n", desiredInterval)
		generation++
		lastUpdateTime = time.Now()
	}

	if !ebiten.IsDrawingSkipped() {
		display(window)
	}

	return err
}

func handleHorizontalArrowsPress() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		scale++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		scale--
	}
}

func handleVerticalArrowsPress() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && timer <= 4900 {
		timer += 100
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && timer >= 200 {
		timer -= 100
	}
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
}

func main() {
	if err := ebiten.Run(frame, width, height, 2, "Game of Life"); err != nil {
		log.Fatal(err)
	}
}
