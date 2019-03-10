package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/nsf/termbox-go"
)

// perlint noise
const seed = 123
const alpha = 2  // noise-yness (more noisey as approaches 1)
const beta = 1.1 // spacing
const iterations = 3

func debug(thing interface{}) {
	f, _ := os.OpenFile("./debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	txt := fmt.Sprintf("%v\n", thing)
	f.WriteString(txt)
	f.Close()
}

// Converts a cord to the spaced cord grid
func grid(i int) int {
	if i%2 == 0 {
		return i
	}
	return i + 1
}

func noise(x int, y int) float64 {
	p := perlin.NewPerlin(alpha, beta, iterations, seed)
	return p.Noise2D(float64(x), float64(y))
}

func renderPebble(x int, y int) {
	termbox.SetCell(x, y, '.', termbox.ColorMagenta, termbox.ColorDefault)
}

func renderLeaf(x int, y int) {
	termbox.SetCell(x, y, ',', termbox.ColorMagenta, termbox.ColorDefault)
}

func renderFlat(x int, y int) {
	termbox.SetCell(x, y, '_', termbox.ColorMagenta, termbox.ColorDefault)
}

func renderText() {
	w, h := termbox.Size()
	termbox.SetCell(w/2-6, h/2, 'R', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(w/2-5, h/2, ' ', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(w/2-4, h/2, 'a', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(w/2-3, h/2, ' ', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(w/2-2, h/2, 'i', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(w/2-1, h/2, ' ', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	termbox.SetCell(w/2+0, h/2, 'n', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
}

func renderGround() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	for y := 1; y < h-2; y++ {
		for x := 1; x < w-3; x++ {
			n := noise(x, y)
			switch {
			case n < 0:
				renderPebble(grid(x), y)
				break
			case n < 0.2:
				renderFlat(grid(x), y)
				break
			case n < 0.4:
				renderLeaf(grid(x), y)
				break
			}
		}
	}
	renderText()
	termbox.Flush()
}

type drop struct {
	x   int
	y   int
	age int
}

func renderDrop(d drop) {

	color := termbox.ColorBlue
	x := grid(d.x)
	y := d.y

	switch d.age {
	case 4:
		termbox.SetCell(x, y, '.', color, termbox.ColorDefault)
		break
	case 5:
		termbox.SetCell(x, y, '-', color, termbox.ColorDefault)
		break
	case 6:
		termbox.SetCell(x, y, 'o', color, termbox.ColorDefault)
		break
	case 7:
		termbox.SetCell(x, y, 'O', color, termbox.ColorDefault)
		break
	case 8:
		termbox.SetCell(x, y, '(', color, termbox.ColorDefault)
		termbox.SetCell(x+2, y, ')', color, termbox.ColorDefault)
		break
	case 9:
		termbox.SetCell(x-2, y, '(', color, termbox.ColorDefault)
		termbox.SetCell(x+2, y, ')', color, termbox.ColorDefault)
		break
	}

}

func renderRain(drops []drop) {
	for i, d := range drops {
		renderDrop(d)
		drops[i].age++
	}

	termbox.Flush()
}

func newDrop() drop {
	w, h := termbox.Size()
	return drop{
		x:   rand.Intn(w-3) + 2,
		y:   rand.Intn(h-2) + 2,
		age: 0,
	}
}

// This is pretty slow/bad go code; don't copy this
// esp for hot code paths
func updateRain(drops []drop) []drop {
	newDrops := []drop{}
	for _, d := range drops {
		if rand.Intn(3) == 0 {
			d.age++
		}

		if d.age < 30 {
			newDrops = append(newDrops, d)
		}
	}

	for i := 0; i < rand.Intn(7); i++ {
		newDrops = append(newDrops, newDrop())
	}

	return newDrops
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	drops := []drop{}

	renderGround()
	for {
		// Exit on any key
		select {
		case ev := <-events:
			if ev.Type == termbox.EventKey {
				termbox.Close()
				os.Exit(0)
			}
		default:
			debug(drops)
			renderGround()
			renderRain(drops)
			drops = updateRain(drops)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
