package buildid

import (
	"testing"
)

func TestFromELFPath(t *testing.T) {
	t.Run("valid ELF file", func(t *testing.T) {
		got, err := FromPath("./testdata/ld-musl-x86_64.so.1")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		want := "d3247ee608f25ff832421ce8c0d42e069e8c21ad"
		if got != want {
			t.Errorf("FromPath() = %v, want %v", got, want)
		}
	})

	t.Run("missing .text section", func(t *testing.T) {
		got, err := FromPath("./testdata/missing-text-sections")
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}
		if got != "" {
			t.Errorf("FromPath() = %v, want empty string", got)
		}
	})
}
