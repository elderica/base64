package main

import (
	"flag"
	"os"
	"fmt"
)

var (
	decode bool
	igngarbage bool
	wrap int
	filename string
	input = os.Stdin
	output = os.Stdout
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
}
