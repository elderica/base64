package main

import (
	"flag"
	"os"
	"fmt"
	"io"
	"encoding/base64"
)

var (
	decode bool
	igngarbage bool
	wrap int
	filename string
	input io.Reader = os.Stdin
	output io.Writer = os.Stdout
)

func init() {
	flag.BoolVar(&decode, "d", false, "")
	flag.BoolVar(&decode, "decode", false, "Decode mode. Default is encode mode.")
	flag.BoolVar(&igngarbage,  "i", false, "")
	flag.BoolVar(&igngarbage,  "ignore-garbage", false, "Ignore unrecognized bytes.")
	flag.IntVar(&wrap, "w", 76, "")
	flag.IntVar(&wrap, "wrap", 76, "Wrap lines during encoding.")
}

func main() {
	flag.Parse()

	filename := flag.Arg(0)
	if filename != "" {
		if filename == "-" {
			goto NOTFILE
		}
		fi, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(-1)
		}
		defer fi.Close()
		input = fi
	}
	NOTFILE:

	if decode {
		input = base64.NewDecoder(base64.StdEncoding, input)
	} else {
		output = base64.NewEncoder(base64.StdEncoding, output)
	}
	io.Copy(output, input)
}
