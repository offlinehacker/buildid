// This file contains code derived from Parca project.
// The original code is available at https://github.com/parca-dev/parca-agent/blob/045369547b2facda1f1dff96e917f45f63c74741/pkg/buildid/buildid_test.go

package buildid

import (
	"debug/elf"
	"encoding/hex"
	"os"
	"testing"
)

func TestFromELF(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "go binary",
			args: args{
				path: "./testdata/readelf-sections",
			},
			want: "38485a695f33313366465a4977783952383553352f7061675079616d5137476a525276786b447243682f564636356c4b554450384b684e71766d5133314a2f49765f39585a33486b576a684f57306661525158",
		},
		{
			name: "rust binary",
			args: args{
				path: "./testdata/rust",
			},
			want: "ea8a38018312ad155fa70e471d4e0039ff9971c6",
		},
		{
			name: "rust binary build with bazel",
			args: args{
				path: "./testdata/bazel-rust",
			},
			want: "983bd888c60ead8e",
		},
		{
			name: "missing .text section",
			args: args{
				path: "./testdata/missing-text-section",
			},
			wantErr: true,
		},
		{
			name: "shared object",
			args: args{
				path: "./testdata/ld-musl-x86_64.so.1",
			},
			want: "d3247ee608f25ff832421ce8c0d42e069e8c21ad",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.path)
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}

			ef, err := elf.NewFile(f)
			if err != nil {
				t.Fatalf("Failed to create ELF file: %v", err)
			}

			got, err := FromELF(ef)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
			if got != tt.want {
				t.Errorf("FromELF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fastGNU(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "rust binary",
			args: args{
				path: "./testdata/rust",
			},
			want: "ea8a38018312ad155fa70e471d4e0039ff9971c6",
		},
		{
			name: "rust binary build with bazel",
			args: args{
				path: "./testdata/bazel-rust",
			},
			want: "983bd888c60ead8e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := elf.Open(tt.args.path)
			if err != nil {
				t.Fatalf("Failed to open ELF file: %v", err)
			}

			got, err := fastGNU(file)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
			if hex.EncodeToString(got) != tt.want {
				t.Errorf("fastGNU() = %v, want %v", hex.EncodeToString(got), tt.want)
			}
		})
	}
}

func Test_buildid(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "go binary",
			args: args{
				path: "./testdata/readelf-sections",
			},
			want: "bd1ca7c3af25af95", // fallbacks to hash of .text
		},
		{
			name: "rust binary",
			args: args{
				path: "./testdata/rust",
			},
			want: "ea8a38018312ad155fa70e471d4e0039ff9971c6",
		},
		{
			name: "rust binary build with bazel",
			args: args{
				path: "./testdata/bazel-rust",
			},
			want: "983bd888c60ead8e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.path)
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}

			ef, err := elf.NewFile(f)
			if err != nil {
				t.Fatalf("Failed to create ELF file: %v", err)
			}

			got, err := buildid(ef)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
			if got != tt.want {
				t.Errorf("buildid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fastGo(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "go binary",
			args: args{
				path: "./testdata/readelf-sections",
			},
			want: "8HZi_313fFZIwx9R85S5/pagPyamQ7GjRRvxkDrCh/VF65lKUDP8KhNqvmQ31J/Iv_9XZ3HkWjhOW0faRQX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := elf.Open(tt.args.path)

			if err != nil {
				t.Fatalf("Failed to open ELF file: %v", err)
			}
			got, err := fastGo(file)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
			if string(got) != tt.want {
				t.Errorf("fastGo() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
