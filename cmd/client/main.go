package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	boardSize    = 10
	boardX       = 120
	boardY       = 174
	cellSize     = 44
)

var (
	backgroundColor = color.RGBA{R: 22, G: 131, B: 196, A: 255}
	cellColor       = color.RGBA{R: 247, G: 232, B: 207, A: 255}
	markedColor     = color.RGBA{R: 215, G: 38, B: 61, A: 255}
	hoverColor      = color.RGBA{R: 255, G: 255, B: 0, A: 255}
)

type game struct {
	cells    [boardSize][boardSize]bool // cells[строка][столбец]
	hoverCol int
	hoverRow int
	hoverOk  bool
}

// логика, по умолчанию 60 раз в секунду
func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	x, y := ebiten.CursorPosition()
	g.hoverCol, g.hoverRow, g.hoverOk = cellAt(x, y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.hoverOk {
		g.cells[g.hoverRow][g.hoverCol] = !g.cells[g.hoverRow][g.hoverCol]
	}

	return nil
}

// отрисовка кадра
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	for col := range boardSize {
		x := float32(boardX + col*cellSize)

		for row := range boardSize {
			y := float32(boardY + row*cellSize)

			clr := cellColor

			if g.hoverOk && row == g.hoverRow && col == g.hoverCol {
				clr = hoverColor
			}

			if g.cells[row][col] {
				clr = markedColor
			}

			vector.FillRect(screen, x, y, cellSize-1, cellSize-1, clr, false)
		}

	}
}

// логическое разрешение
func (g *game) Layout(outW, outH int) (int, int) {
	return screenWidth, screenHeight
}

func cellAt(x, y int) (col, row int, ok bool) {
	if x < boardX || y < boardY {
		return 0, 0, false
	}

	col = (x - boardX) / cellSize
	row = (y - boardY) / cellSize

	if col >= boardSize || row >= boardSize {
		return 0, 0, false
	}

	return col, row, true
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sea Battle")

	g := game{}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
