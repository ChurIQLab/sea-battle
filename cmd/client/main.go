package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	squareSize   = 100
	speed        = 4
)

type Direction int

const (
	Up   Direction = iota // iota - счетчик: Up = 0
	Down                  // Down = 1
)

type Game struct {
	ticks int
	red   uint8
	move  Direction
	x     float32
	y     float32
}

// логика, по умолчанию 60 раз в секунду
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.x = min(g.x+speed, screenWidth-squareSize)
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.x = max(g.x-speed, 0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.y = min(g.y+speed, screenHeight-squareSize)
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.y = max(g.y-speed, 0)
	}

	switch g.move {
	case Up:
		g.red++
	case Down:
		g.red--
	}

	switch g.red {
	case 255:
		g.move = Down
	case 0:
		g.move = Up
	}

	g.ticks++
	if g.ticks >= ebiten.TPS() {
		g.ticks = 0
		log.Printf("FPS: %.1f TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS())
	}

	return nil
}

// отрисовка кадра
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: g.red, G: 0, B: 0, A: 255})
	vector.FillRect(screen, g.x, g.y, squareSize, squareSize, color.RGBA{R: 255, G: 212, B: 71, A: 255}, false)
}

// логическое разрешение
func (g *Game) Layout(outW, outH int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sea Battle")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
