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

type orientation int

const (
	horizontal orientation = iota
	vertical
)

var (
	backgroundColor = color.RGBA{R: 22, G: 131, B: 196, A: 255}
	cellColor       = color.RGBA{R: 247, G: 232, B: 207, A: 255}
	shipColor       = color.RGBA{R: 126, G: 144, B: 168, A: 255}
	ghostOKColor    = color.RGBA{R: 80, G: 200, B: 120, A: 255}
	ghostBadColor   = color.RGBA{R: 220, G: 70, B: 70, A: 255}
	fleet           = []int{4, 3, 3, 2, 2, 2, 1, 1, 1, 1}
)

type coord struct {
	row, col int
}
type game struct {
	ships   []ship
	hover   coord
	hoverOK bool
	orient  orientation
}

// логика, по умолчанию 60 раз в секунду
func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	x, y := ebiten.CursorPosition()
	g.hover, g.hoverOK = cellAt(x, y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.hoverOK && g.isPlacing() {
		newShip := g.candidate()
		if newShip.fitsBoard() {
			g.ships = append(g.ships, newShip)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		switch g.orient {
		case vertical:
			g.orient = horizontal
		case horizontal:
			g.orient = vertical
		}
	}

	return nil
}

// отрисовка кадра
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	for col := range boardSize {
		for row := range boardSize {
			c := coord{row: row, col: col}
			drawCell(screen, c, cellColor)
		}
	}

	for _, s := range g.ships {
		for _, c := range s.cells() {
			drawCell(screen, c, shipColor)
		}
	}

	if g.hoverOK && g.isPlacing() {
		newShip := g.candidate()
		clr := ghostOKColor

		if !newShip.fitsBoard() {
			clr = ghostBadColor
		}

		for _, c := range newShip.cells() {
			drawCell(screen, c, clr)
		}
	}
}

// логическое разрешение
func (g *game) Layout(outW, outH int) (int, int) {
	return screenWidth, screenHeight
}

// candidate возвращает корабль, который встанет при клике. Длина берётся из fleet по числу уже поставленных кораблей.
func (g *game) candidate() ship {
	return ship{start: g.hover, length: fleet[len(g.ships)], orient: g.orient}
}

func (g *game) isPlacing() bool {
	return len(g.ships) < len(fleet)
}

// cellAt переводит пиксели экрана в клетку поля. ok = false, если точка вне поля.
func cellAt(x, y int) (coord, bool) {
	if x < boardX || y < boardY {
		return coord{}, false
	}

	row := (y - boardY) / cellSize
	col := (x - boardX) / cellSize

	if row >= boardSize || col >= boardSize {
		return coord{}, false
	}

	return coord{row: row, col: col}, true
}

func drawCell(screen *ebiten.Image, c coord, clr color.Color) {
	x := float32(boardX + c.col*cellSize)
	y := float32(boardY + c.row*cellSize)
	vector.FillRect(screen, x, y, cellSize-1, cellSize-1, clr, false)
}

type ship struct {
	start  coord
	length int
	orient orientation
}

func (s ship) cells() []coord {
	result := make([]coord, 0, s.length)

	for i := range s.length {
		c := s.start
		if s.orient == horizontal {
			c.col += i
		} else {
			c.row += i
		}
		result = append(result, c)
	}
	return result
}

func (s ship) fitsBoard() bool {
	for _, c := range s.cells() {
		if c.col >= boardSize || c.row >= boardSize {
			return false
		}
	}
	return true
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sea Battle")

	g := game{}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
