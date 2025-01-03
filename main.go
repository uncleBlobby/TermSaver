package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

type TermSaver struct {
	Width, Height int
	Saver         [][]string
	Interrupts    chan os.Signal
}

type Vector2 struct {
	X int
	Y int
}

type Ball struct {
	X        int
	Y        int
	Char     string
	Velocity Vector2
	History  []Vector2
}

func Init() *TermSaver {
	ts := &TermSaver{}
	ts.InitInterrupts()
	ts.InitDrawing()
	return ts
}

func (ts *TermSaver) SetSize() {
	var err error
	ts.Width, ts.Height, err = term.GetSize(0)
	if err != nil {
		log.Printf("[INIT] Could not get terminal size: %s", err)
		log.Printf("[ABORT]")
		os.Exit(1)
	}
}

func (ts *TermSaver) InitInterrupts() {
	ts.Interrupts = make(chan os.Signal, 1)
	signal.Notify(ts.Interrupts, syscall.SIGINT, syscall.SIGTERM)
}

func (ts *TermSaver) InitDrawing() {
	ts.Saver = [][]string{}
}

func (ts *TermSaver) Draw(b *Ball) {
	for {
		ts.SetSize()
		ts.InitDrawing()

		for i := 0; i < ts.Height-1; i++ {
			row := make([]string, ts.Width)
			if i == 0 || i == ts.Height-2 {
				for j := 0; j < len(row)-1; j++ {
					row[j] = "="
				}
			} else {
				for j := 0; j < len(row)-1; j++ {
					if j == 0 {
						row[j] = "|"
					} else if j == ts.Width-2 {
						row[j] = "|"
					} else {
						row[j] = " "
					}
				}
			}
			ts.Saver = append(ts.Saver, row)
		}

		ts.Saver[b.Y][b.X] = b.Char
		b.History = append(b.History, Vector2{X: b.X, Y: b.Y})

		b.X += b.Velocity.X
		b.Y += b.Velocity.Y

		if b.X >= ts.Width-3 || b.X <= 1 {
			b.Velocity.X *= -1
		}
		if b.Y >= ts.Height-3 || b.Y <= 1 {
			b.Velocity.Y *= -1
		}

		b.History = append(b.History, Vector2{X: b.X, Y: b.Y})

		if len(b.History) > 10 {
			b.History = b.History[len(b.History)-10:]
		}

		for ind, historyDraw := range b.History {
			if ind == len(b.History)-1 {
				ts.Saver[historyDraw.Y][historyDraw.X] = "O"
			} else if ind > 2 {
				ts.Saver[historyDraw.Y][historyDraw.X] = "o"
			} else {
				ts.Saver[historyDraw.Y][historyDraw.X] = "."
			}
		}

		for i := 0; i < len(ts.Saver); i++ {
			for j := 0; j < len(ts.Saver[i]); j++ {
				// if j == b.X && i == b.Y {
				// 	fmt.Printf("%s", b.Char)
				// } else {
				fmt.Printf("%s", ts.Saver[i][j])
				// }
			}
			fmt.Printf("\n")
		}
		time.Sleep(16 * time.Millisecond)
		ts.ClearScreen()
	}
}

func (ts *TermSaver) ClearScreen() {
	// fmt.Print("\033[H\033[2J")
	fmt.Print("\033[H")
	// cmd := exec.Command("tput reset")
	// cmd.Stdout = os.Stdout
	// cmd.Run()
}

func main() {

	ts := Init()

	// clear terminal before running termsaver
	fmt.Print("\033[H\033[2J")
	// hide cursor while termsaver is running
	fmt.Print("\x1b[?25l")

	// initialize channel for listening to terminal interrupt signal
	// TODO: listen for
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// initialize goroutine listener for signal event, which is discarded before redrawing cursor and exiting
	go func() {
		_ = <-ts.Interrupts
		// fmt.Println()
		// fmt.Println(sig)
		ts.ClearScreen()
		fmt.Print("\x1b[?25h")
		os.Exit(1)
	}()

	b := Ball{
		X: 5, Y: 5, Char: "O", Velocity: Vector2{X: 1, Y: 1},
	}
	// main application loop
	ts.Draw(&b)
}
