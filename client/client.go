package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"miniGame/client/sprite"
	"miniGame/client/world"
	"net"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed resource/*
var resfs embed.FS

const (
	worldWidth   = 10000
	worldHeight  = 10000
	windowWidth  = 800
	windowHeight = 600

	moveStep     = 5
	viewMoveStep = 3
)

type Game struct {
	keys []ebiten.Key

	wld   *world.World
	viewX int
	viewY int

	mainCharacter *sprite.Sprite
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, k := range g.keys {
		switch k {
		case ebiten.KeyA:
			x, _ := g.mainCharacter.Position()
			step := math.Min(moveStep, x)
			g.mainCharacter.Move(-step, 0)
		case ebiten.KeyD:
			x, _ := g.mainCharacter.Position()
			step := math.Min(moveStep, worldWidth-x)
			g.mainCharacter.Move(step, 0)
		case ebiten.KeyW:
			_, y := g.mainCharacter.Position()
			step := math.Min(moveStep, y)
			g.mainCharacter.Move(0, -step)
		case ebiten.KeyS:
			_, y := g.mainCharacter.Position()
			step := math.Min(moveStep, worldHeight-y)
			g.mainCharacter.Move(0, step)

		case ebiten.KeyLeft:
			g.viewX -= viewMoveStep
			if g.viewX < 0 {
				g.viewX = 0
			}
		case ebiten.KeyRight:
			g.viewX += viewMoveStep
			if g.viewX > worldWidth-windowWidth {
				g.viewX = worldWidth - windowWidth
			}
		case ebiten.KeyUp:
			g.viewY -= viewMoveStep
			if g.viewY < 0 {
				g.viewY = 0
			}
		case ebiten.KeyDown:
			g.viewY += viewMoveStep
			if g.viewY > worldHeight-windowHeight {
				g.viewY = worldHeight - windowHeight
			}
		}
	}

	g.mainCharacter.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	view := image.Rect(g.viewX, g.viewY, g.viewX+windowWidth, g.viewY+windowHeight)
	g.wld.DrawTerrian(screen, view)
	g.mainCharacter.Draw(screen, view)
	x, y := g.mainCharacter.Position()
	t := fmt.Sprintf("%f,%f", x, y)
	ebiten.SetWindowTitle(t)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func loadImage(filename string) image.Image {
	b, err := resfs.ReadFile(filename)
	if err != nil {
		log.Fatalf("open tile: %v", err)
	}

	img, _, err := image.Decode(bytes.NewBuffer(b))
	if err != nil {
		log.Fatalf("open tile image: %v", err)
	}
	return img
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		msg := "hello server \n"
		_, e := conn.Write([]byte(msg))
		if e != nil {
			log.Fatal(err)
		}
	}

	tileImg := loadImage("resource/tile.png")
	runImgs := []image.Image{
		loadImage("resource/run0.png"),
		loadImage("resource/run1.png"),
	}
	g := &Game{
		wld:           world.New(worldWidth, worldHeight, tileImg),
		mainCharacter: sprite.New(runImgs),
	}
	g.mainCharacter.Move(400, 300)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Game World Demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
