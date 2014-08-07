package main

import (
	"flag"
	"os"
	"fmt"
	"io"
	"bytes"
	"encoding/base64"
)

const base64alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

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

type Base64Cleaner struct {
	r io.Reader
}

func NewBase64Cleaner(r io.Reader) io.Reader {
	return &Base64Cleaner{r}
}

func (c *Base64Cleaner) Read(p []byte) (int, error) {
	buf := make([]byte, len(p))
	nr, err := c.r.Read(buf)
	if err != nil {
		return 0, err
	}
	var j int
	for i := 0; i < nr; i++ {
		if bytes.IndexByte([]byte(base64alpha), buf[i]) >= 0 || buf[i] == '=' {
			p[j] = buf[i]
			j++
		}
	}
	return j, nil
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
		if igngarbage {
			input = NewBase64Cleaner(input)
		}
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
