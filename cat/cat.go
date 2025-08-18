package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) >= 2 {

		for i := 1; i < len(os.Args); i++ {

			data, err := os.ReadFile(os.Args[i])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(string(data))
		}
	}
}
