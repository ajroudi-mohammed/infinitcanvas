
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
)

var (
	screenWidth  = 800
	screenHeight = 600
)

type Point struct {
	X, Y float64
}

type Drawing struct {
	Points []Point
}

type Game struct {
	offsetX     float64
	offsetY     float64
	mouseX      int
	mouseY      int
	dragging    bool
	lastMouseX  int
	lastMouseY  int
	isDrawing   bool
	drawings    []Drawing
	currentPath []Point
}

func (g *Game) Update() error {
	// Handle keyboard input
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.offsetX += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.offsetX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.offsetY += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.offsetY -= 5
	}

	// Handle mouse input
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.handleLeftMouse()
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.handleRightMouseDown()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		g.handleRightMouseUp()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the background with yellow color
	screen.Fill(color.NRGBA{R: 255, G: 255, B: 0, A: 255})

	// Draw a border around the screen
	ebitenutil.DrawLine(screen, 0, 0, float64(screenWidth-1), 0, color.Black)
	ebitenutil.DrawLine(screen, 0, 0, 0, float64(screenHeight-1), color.Black)
	ebitenutil.DrawLine(screen, float64(screenWidth-1), 0, float64(screenWidth-1), float64(screenHeight-1), color.Black)
	ebitenutil.DrawLine(screen, 0, float64(screenHeight-1), float64(screenWidth-1), float64(screenHeight-1), color.Black)

	// Draw a grid
	for x := -int(g.offsetX)%50 - 50; x < screenWidth; x += 50 {
		for y := -int(g.offsetY)%50 - 50; y < screenHeight; y += 50 {
			ebitenutil.DrawRect(screen, float64(x), float64(y), 48, 48, color.Black)
		}
	}

	// Draw the drawings
	for _, drawing := range g.drawings {
		g.drawSVGPath(screen, drawing.Points)
	}

	// Draw the current drawing path
	if g.isDrawing && len(g.currentPath) > 1 {
		g.drawSVGPath(screen, g.currentPath)
	}

	// Draw an "X" to indicate movement
	ebitenutil.DebugPrintAt(screen, "X", int(100-g.offsetX), int(100-g.offsetY))

	msg := "Use arrow keys to move, right-click and drag to draw"
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) handleLeftMouse() {
	if !g.dragging {
		g.mouseX, g.mouseY = ebiten.CursorPosition()
		g.dragging = true
	} else {
		curX, curY := ebiten.CursorPosition()
		g.offsetX -= float64(curX - g.mouseX)
		g.offsetY -= float64(curY - g.mouseY)
		g.mouseX = curX
		g.mouseY = curY
	}
}

func (g *Game) handleRightMouseDown() {
	if !g.isDrawing {
		g.currentPath = nil
	}
	curX, curY := ebiten.CursorPosition()
	g.currentPath = append(g.currentPath, Point{X: float64(curX), Y: float64(curY)})
	g.isDrawing = true
}

func (g *Game) handleRightMouseUp() {
	if g.isDrawing && len(g.currentPath) > 1 {
		g.drawings = append(g.drawings, Drawing{Points: g.currentPath})
	}
	g.isDrawing = false
}

func (g *Game) drawSVGPath(screen *ebiten.Image, points []Point) {
	if len(points) < 2 {
		return
	}
	for i := 0; i < len(points)-1; i++ {
		x1, y1 := points[i].X-g.offsetX, points[i].Y-g.offsetY
		x2, y2 := points[i+1].X-g.offsetX, points[i+1].Y-g.offsetY
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.White)
	}
}

func main() {
	game := &Game{}
	ebiten.SetWindowTitle("Infinite Canvas Example")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
