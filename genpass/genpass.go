package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	
	"github.com/atotto/clipboard"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"
const specials = "!@#$%^&*()-='"

func main() {
	length := flag.Int("l", 20, "length of the password to generate")
	clip := flag.Bool("c", false, "copy to clipboard")
	silent := flag.Bool("s", false, "show no password")
	low := flag.Bool("low", false, "use only letters (low strength)")
	medium := flag.Bool("medium", false, "use letters and numbers (medium strength)")
	flag.Parse()

	if *low && *medium {
		fmt.Println("-low and -medium flags cannot be used simultaneously.")
		os.Exit(1)
	}

	if !(*clip) && *silent {
		fmt.Println("Kinda useless. Use -c with -s flag to copy to the clipboard without printing to stdout.")
		os.Exit(1)
	}

	if *length < 1 {
		fmt.Println("Incorrect length.")
		os.Exit(1)
	}

	var charset string
	if *low {
		charset = letters
	} else if *medium {
		charset = letters + digits
	} else {
		// high is default
		charset = letters + digits + specials
	}

	result := make([]byte, *length)
	for i := 0; i < *length; i++ {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		n := nBig.Int64()
		result[i] = charset[n]
	}

	if !(*silent) {
		fmt.Println(string(result))
	}

	if *clip {
		clipboard.WriteAll(string(result))
	}
}
