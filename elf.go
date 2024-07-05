// This file contains code derived from Parca project.
// The original code is available at https://github.com/parca-dev/parca-agent/blob/045369547b2facda1f1dff96e917f45f63c74741/pkg/buildid/elf.go

package buildid

import (
	"debug/elf"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/cespare/xxhash/v2"
)

const goBuildIDSectionName = ".note.go.buildid"

// FromELF returns the build ID for an ELF binary.
//
// This method takes `elf.File` as an input and will attempt to extract buildid using different methods:
//
//   - Reading the ".note.go.buildid" section for Go binaries.
//   - Reading the ".note.gnu.build-id" section for binaries with GNU build ID.
//   - Hashing the ".text" section if no build ID is found.
func FromELF(ef *elf.File) (string, error) {
	// First, try fast methods.
	hasGoBuildIDSection := false
	for _, s := range ef.Sections {
		if s.Name == goBuildIDSectionName {
			hasGoBuildIDSection = true
			break
		}
	}

	if hasGoBuildIDSection {
		if id, err := fastGo(ef); err == nil && len(id) > 0 {
			return hex.EncodeToString(id), nil
		}
	}

	if id, err := fastGNU(ef); err == nil && len(id) > 0 {
		return hex.EncodeToString(id), nil
	}

	// If that fails, try the slow methods.
	return buildid(ef)
}

// buildid returns the build id for an ELF binary by:
// 1. First, looking for a GNU build-id note.
// 2. If fails, hashing the .text section.
func buildid(ef *elf.File) (string, error) {
	// Search through all the notes for a GNU build ID.
	b, err := slowGNU(ef)
	if err == nil {
		if len(b) > 0 {
			return hex.EncodeToString(b), nil
		}
	}

	// If we didn't find a GNU build ID, try hashing the .text section.
	text := ef.Section(".text")
	if text == nil {
		return "", errors.New("could not find .text section")
	}
	h := xxhash.New()
	if _, err := io.Copy(h, text.Open()); err != nil {
		return "", fmt.Errorf("hash elf .text section: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// fastGo returns the Go build-ID for an ELF binary by searching specific locations.
// (nil, nil) is returned if no build-ID is found.
func fastGo(ef *elf.File) ([]byte, error) {
	s := ef.Section(goBuildIDSectionName)
	if s == nil {
		return nil, fmt.Errorf("failed to find %s section", goBuildIDSectionName)
	}

	notes, err := parseNotes(s.Open(), int(s.Addralign), ef.ByteOrder)
	if err != nil {
		return nil, err
	}

	var buildID []byte
	for _, note := range notes {
		if note.Name == "Go" && note.Type == noteTypeGoBuildID {
			if len(buildID) == 0 {
				buildID = note.Desc
			} else {
				return nil, fmt.Errorf("multiple build ids found, don't know which to use")
			}
		}
	}
	if len(buildID) > 0 {
		return buildID, nil
	}
	return nil, nil
}

const gnuBuildIDSectionName = ".note.gnu.build-id"

// fastGNU returns the GNU build-ID for an ELF binary by searching specific locations.
// (nil, nil) is returned if no build-ID is found.
func fastGNU(ef *elf.File) ([]byte, error) {
	s := ef.Section(gnuBuildIDSectionName)
	if s == nil {
		return nil, fmt.Errorf("failed to find %s section", gnuBuildIDSectionName)
	}

	notes, err := parseNotes(s.Open(), int(s.Addralign), ef.ByteOrder)
	if err != nil {
		return nil, err
	}

	return findGNU(notes)
}

// findGNU returns the GNU build-ID for an ELF binary by searching through the given notes.
// (nil, nil) is returned if no build-ID is found.
func findGNU(notes []elfNote) ([]byte, error) {
	var buildID []byte
	for _, note := range notes {
		if note.Name == "GNU" && note.Type == noteTypeGNUBuildID {
			if len(buildID) == 0 {
				buildID = note.Desc
			} else {
				return nil, fmt.Errorf("multiple build ids found, don't know which to use")
			}
		}
	}
	if len(buildID) > 0 {
		return buildID, nil
	}
	return nil, nil
}

// slowGNU returns the GNU build-ID for an ELF binary by searching through all.
// (nil, nil) is returned if no build-ID is found.
func slowGNU(ef *elf.File) ([]byte, error) {
	for _, p := range ef.Progs {
		if p.Type != elf.PT_NOTE {
			continue
		}
		notes, err := parseNotes(p.Open(), int(p.Align), ef.ByteOrder)
		if err != nil {
			return nil, fmt.Errorf("parse notes: %w", err)
		}
		b, err := findGNU(notes)
		if err != nil {
			return nil, err
		}
		if len(b) > 0 {
			return b, nil
		}
	}
	for _, s := range ef.Sections {
		if s.Type != elf.SHT_NOTE {
			continue
		}
		notes, err := parseNotes(s.Open(), int(s.Addralign), ef.ByteOrder)
		if err != nil {
			return nil, fmt.Errorf("parse notes: %w", err)
		}
		b, err := findGNU(notes)
		if err != nil {
			return nil, err
		}
		if len(b) > 0 {
			return b, nil
		}
	}
	return nil, nil
}
