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

func (ts *TermSaver) Draw() {
	for {
		ts.SetSize()

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

		for i := 0; i < len(ts.Saver); i++ {
			for j := 0; j < len(ts.Saver[i]); j++ {
				fmt.Printf("%s", ts.Saver[i][j])
			}
			fmt.Printf("\n")
		}
		time.Sleep(1 * time.Second)
	}
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
		fmt.Print("\x1b[?25h")
		os.Exit(1)
	}()

	// main application loop
	ts.Draw()
	// 	for {
	// 		// get dimensions of terminal window
	// 		// runs each update for hot - resizing ie: terminal can be resized and screensaver will update to the new dimensions
	// 		width, height, err := term.GetSize(0)
	// 		if err != nil {
	// 			log.Printf("Error: %s", err)
	// 			os.Exit(1)
	// 		}

	// 		// make 2d array of chars for screen drawing
	// 		saver := [][]string{}

	// 		for i := 0; i < height-1; i++ {
	// 			row := make([]string, width)
	// 			if i == 0 || i == height-2 {
	// 				for j := 0; j < len(row)-1; j++ {
	// 					row[j] = "="
	// 				}
	// 			} else {
	// 				for j := 0; j < len(row)-1; j++ {
	// 					if j == 0 {
	// 						row[j] = "|"
	// 					} else if j == width-2 {
	// 						row[j] = "|"
	// 					} else {
	// 						row[j] = " "
	// 					}
	// 				}
	// 				// row[0] = "="
	// 				// row[width-1] = "^"
	// 			}
	// 			saver = append(saver, row)
	// 		}

	// 		for i := 0; i < len(saver); i++ {
	// 			for j := 0; j < len(saver[i]); j++ {
	// 				fmt.Printf("%s", saver[i][j])
	// 			}
	// 			fmt.Printf("\n")
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 	}
}
