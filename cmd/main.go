package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tamasd/xsddict"
)

var (
	extra = flag.String("extra", "", "Extra file to include in the dictionary. Useful to include XML headers.")
	usage = `Usage xsddict [INPUT XSD] [OUTPUT FILE]`
)

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		printUsage()
	}

	xsdFile, err := os.Open(flag.Arg(0))
	handleErr(err)

	outputFile, err := os.OpenFile(flag.Arg(1), os.O_CREATE|os.O_RDWR, 0644)
	handleErr(err)

	if *extra != "" {
		extraFile, err := os.Open(*extra)
		handleErr(err)

		_, err = io.Copy(outputFile, extraFile)
		handleErr(err)

		handleErr(extraFile.Close())
		handleErr(outputFile.Sync())
	}

	handleErr(xsddict.WhiteSpaces(outputFile))
	handleErr(outputFile.Sync())

	handleErr(xsddict.GenerateDict(outputFile, xsdFile))
	handleErr(outputFile.Sync())

	handleErr(xsdFile.Close())
	handleErr(outputFile.Close())
}

func printUsage() {
	fmt.Println(usage)
	os.Exit(1)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
