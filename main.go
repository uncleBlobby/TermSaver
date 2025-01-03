package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"
)

func main() {
	// if term.IsTerminal(0) {
	// 	fmt.Println("in a term")
	// } else {
	// 	fmt.Println("not a term")
	// }

	fmt.Print("\033[H\033[2J")
	for {
		width, height, err := term.GetSize(0)
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(1)
		}
		// fmt.Printf("width: %d, height: %d", width, height)

		// make 2d array of chars for screen drawing

		saver := [][]string{}

		for i := 0; i < height-1; i++ {
			row := make([]string, width)
			if i == 0 || i == height-2 {
				for j := 0; j < len(row)-1; j++ {
					row[j] = "="
				}
			} else {
				for j := 0; j < len(row)-1; j++ {
					if j == 0 {
						row[j] = "|"
					} else if j == width-2 {
						row[j] = "|"
					} else {
						row[j] = " "
					}
				}
				// row[0] = "="
				// row[width-1] = "^"
			}
			saver = append(saver, row)
		}

		for i := 0; i < len(saver); i++ {
			for j := 0; j < len(saver[i]); j++ {
				fmt.Printf("%s", saver[i][j])
			}
			fmt.Printf("\n")
		}
		time.Sleep(1 * time.Second)
	}
}
