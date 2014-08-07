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

type WrapWriter struct {
	w io.Writer
	wn int
	c int
}

func NewWrapWriter(w io.Writer, lim int) *WrapWriter {
	return &WrapWriter{w, lim, 0}
}

func (w *WrapWriter) Write(p []byte) (int, error) {
	r := 0
	for i, _ := range p {
		nw, err := w.w.Write(p[i:i+1])
		if err != nil {
			return r + nw, err
		}
		w.c += nw
		r += nw
		if w.c % w.wn == 0 {
			_, err := w.w.Write([]byte("\n"))
			if err != nil {
				return r, err
			}
		}
	}
	return r, nil
}

func main() {
	flag.Parse()

	var (
		fi *os.File
		err error
	)
	filename := flag.Arg(0)
	if filename != "" {
		if filename == "-" {
			goto NOTFILE
		}
		fi, err = os.Open(filename)
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
		// When wrap is 0, disable line wrapping.
		if wrap < 0 {
			fmt.Fprintf(os.Stderr, "invalid wrap size: %d", wrap)
			fi.Close()
			os.Exit(-2)
		} else if 0 < wrap {
			output = NewWrapWriter(output, wrap)
		}
		output = base64.NewEncoder(base64.StdEncoding, output)
	}
	io.Copy(output, input)
}
