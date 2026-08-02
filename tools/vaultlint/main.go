package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	root := flag.String("root", "../../Obsidian Vault", "path to the Obsidian Vault")
	flag.Parse()

	issues, err := Lint(*root)
	if err != nil {
		log.Fatal(err)
	}

	if len(issues) == 0 {
		fmt.Println("vault lint passed")
		return
	}

	for _, issue := range issues {
		fmt.Fprintln(os.Stderr, issue.String())
	}
	fmt.Fprintf(os.Stderr, "\nvault lint found %d issue(s)\n", len(issues))
	os.Exit(1)
}
