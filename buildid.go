package buildid

import (
	"debug/elf"
	"os"
)

// FromPath attempts to extract builid from provided path to the binary.
//
// See `FromELF` for more information.
func FromPath(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	ef, err := elf.NewFile(f)
	if err != nil {
		return "", err
	}
	defer ef.Close()

	return FromELF(ef)
}
