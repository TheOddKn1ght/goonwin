package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	length := flag.Int("l", 20, "length of the password to generate")
	clip := flag.Bool("c", false, "copy to clipboard")
	silent := flag.Bool("s", false, "show no password")
	flag.Parse()
	if !(*clip) && *silent {
		fmt.Println("Kinda useless. Use -c with -s flag to copy to the clipboard without printing to stdout.")
		os.Exit(1)
	}
	if *length < 1 {
		fmt.Println("Incorrect length.")
		os.Exit(1)
	}
	const ascii = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890!@#$%^&*()-='"
	result := make([]byte, *length)
	for i := 0; i < *length; i++ {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(ascii))))
		if err != nil {
			panic(err)
		}
		n := nBig.Int64()
		result[i] = ascii[n]
	}
	if !(*silent) {
		fmt.Println(string(result))
	}
	if *clip {
		clipboard.WriteAll(string(result))
	}
}
