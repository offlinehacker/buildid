package main

import (
	"fmt"
	"os"

	"github.com/offlinehacker/buildid"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path-to-binary>\n", os.Args[0])
		os.Exit(2)
	}

	id, err := buildid.FromPath(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "get build id: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(id)
}
