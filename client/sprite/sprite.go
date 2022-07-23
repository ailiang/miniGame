package sprite

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	offsetX = 256
	offsetY = 282

	animationFrames = 5
)

type Sprite struct {
	x, y float64
	imgs []*ebiten.Image
	cur  int

	running        bool
	animationCount int
}

func New(imgs []image.Image) *Sprite {
	s := &Sprite{}
	for _, img := range imgs {
		s.imgs = append(s.imgs, ebiten.NewImageFromImage(img))
	}
	return s
}

func (s *Sprite) Update() {
	if s.running {
		s.animationCount++
		if s.animationCount >= animationFrames {
			s.cur++
			s.animationCount = 0
			if s.cur >= len(s.imgs) {
				s.cur = 0
				s.running = false
			}
		}
	}
}

func (s *Sprite) Position() (x, y float64) { return s.x, s.y }

func (s *Sprite) Move(dx, dy float64) {
	s.x += dx
	s.y += dy

	if !s.running {
		s.running = true
		s.animationCount = 0
	}
}

func (s *Sprite) Draw(screen *ebiten.Image, view image.Rectangle) {
	dx := s.x - float64(view.Min.X)
	dy := s.y - float64(view.Min.Y)
	if dx < 0 || dy < 0 || dx >= float64(view.Max.X) || dy >= float64(view.Max.Y) {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-offsetX, -offsetY)
	op.GeoM.Translate(dx, dy)
	screen.DrawImage(s.imgs[s.cur], op)
}
