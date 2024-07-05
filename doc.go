// Package buildid provides functionality to extract the build ID from ELF binaries on Linux systems.
//
// The build ID is a unique identifier embedded in an ELF binary that helps identify the specific build of the binary.
//
// Example usage:
//
//	package main
//
//	import (
//		"fmt"
//		"path/to/buildid"
//	)
//
//	func main() {
//		buildID, err := buildid.FromPath("path/to/binary")
//		if err != nil {
//			fmt.Println("Error:", err)
//			return
//		}
//		fmt.Println("Build ID:", buildID)
//	}
//
// Note: This package is intended for use on Linux systems only, as it relies on the ELF binary format.
//
// For more information and detailed documentation, please refer to the package's README file and source code.
package buildid
