# ELF Build ID Extractor

This project provides a Go package for extracting the build ID from ELF binaries on Linux systems.

## Description

The `buildid` package allows you to retrieve the build ID from an ELF binary file. It supports extracting the build ID using the following methods:

1. Reading the `.note.go.buildid` section for Go binaries.
2. Reading the `.note.gnu.build-id` section for binaries with GNU build ID.
3. Hashing the `.text` section if no build ID is found.

## Usage

```go
import "github.com/offlinehacker/buildid"

buildID, err := buildid.FromPath("path/to/binary")
if err != nil {
    // Handle the error
}

fmt.Println("Build ID:", buildID)
```

## Platform Support

Please note that this package is intended for use on Linux systems only, as it relies on the ELF binary format.

## License

This project is licensed under the Apache License, Version 2.0. See the LICENSE file for more information.

## Contributing

Contributions to this project are welcome. If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the project's GitHub repository.

## Acknowledgments
 
This project is based on code from [parca-agent](https://github.com/parca-dev/parca-agent/tree/main/pkg/buildid)
