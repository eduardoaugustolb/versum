package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eduardoaugustolb/versum/bible/tools/internal/biblecorpus"
)

func main() {
	rawRoot := flag.String("raw", "../raw", "path to the raw Bible corpus")
	outputRoot := flag.String("output", "../corpus/v1", "path for the canonical Bible corpus")
	flag.Parse()

	manifest, err := biblecorpus.Generate(*rawRoot, *outputRoot)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("generated %d books, %d chapters and %d verses\n", manifest.BookCount, manifest.ChapterCount, manifest.VerseCount)
}
