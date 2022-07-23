package world

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	width     int
	height    int
	tileImage *ebiten.Image
}

func New(width, height int, tile image.Image) *World {
	return &World{
		width:     width,
		height:    height,
		tileImage: ebiten.NewImageFromImage(tile),
	}
}

func (wld *World) DrawTerrian(screen *ebiten.Image, view image.Rectangle) {
	var (
		w, h int
		dw   int
		dh   = view.Dy()
	)
	for y := 0; y < view.Dy(); y += h {
		dw = view.Dx()
		for x := 0; x < view.Dx(); x += w {
			w, h = wld.drawTile(screen, view.Min.X, view.Min.Y, x, y, dw, dh)
			dw -= w
		}
		dh -= h
		//log.Printf("-")
	}
	//log.Printf("--")
}

func (w *World) drawTile(screen *ebiten.Image, vx, vy, dx, dy, dw, dh int) (sw, sh int) {
	iw := w.tileImage.Bounds().Dx()
	ih := w.tileImage.Bounds().Dy()

	sx := (vx + dx) % iw
	sy := (vy + dy) % ih

	sw = iw
	if sw > dw {
		sw = dw
	}
	if sx+sw > iw {
		sw = iw - sx
	}
	sh = ih
	if sh > dh {
		sh = dh
	}
	if sy+sh > ih {
		sh = ih - sy
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	//log.Printf("%d %d %d %d -> %d %d", sx, sy, sw, sh, dx, dy)
	screen.DrawImage(w.tileImage.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image), op)
	return
}
