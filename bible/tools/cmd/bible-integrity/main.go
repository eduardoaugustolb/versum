package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eduardoaugustolb/versum/bible/tools/internal/biblecorpus"
)

func main() {
	root := flag.String("root", "../corpus/v1", "path to the canonical Bible corpus")
	flag.Parse()

	if err := biblecorpus.Verify(*root); err != nil {
		log.Fatal(err)
	}

	fmt.Println("bible corpus integrity check passed")
}
