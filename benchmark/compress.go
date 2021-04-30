package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dsnet/compress/bzip2"
)

type compressor interface {
	io.WriteCloser
}

type compressorFactory func(w io.Writer) compressor

func main() {
	dict, err := os.ReadFile("dict.dat")
	must(err)

	must(runTest("flatedict", func(w io.Writer) compressor {
		c, err := flate.NewWriterDict(w, flate.BestCompression, dict)
		must(err)
		return c
	}))
	must(runTest("flate", func(w io.Writer) compressor {
		c, err := flate.NewWriter(w, flate.BestCompression)
		must(err)
		return c
	}))
	must(runTest("gzip", func(w io.Writer) compressor {
		c, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		must(err)
		return c
	}))
	must(runTest("bzip", func(w io.Writer) compressor {
		c, err := bzip2.NewWriter(w, &bzip2.WriterConfig{
			Level: bzip2.BestCompression,
		})
		must(err)
		return c
	}))
}

func runTest(name string, factory compressorFactory) error {
	matches, err := filepath.Glob("data/*.xml")
	if err != nil {
		return err
	}

	for _, fn := range matches {
		if err := run(name+" "+fn, factory, fn); err != nil {
			return err
		}
	}

	return nil
}

func run(name string, factory compressorFactory, input string) error {
	f, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)

	c := factory(buf)

	start := time.Now()
	if _, err = c.Write(f); err != nil {
		return err
	}
	if err = c.Close(); err != nil {
		return err
	}
	duration := time.Since(start)

	orig := len(f)
	compressed := len(buf.Bytes())
	percent := float64(compressed) / float64(orig) * 100
	fmt.Printf("%s %d => %d (%f%%) in %s\n", name, orig, compressed, percent, duration.String())

	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
