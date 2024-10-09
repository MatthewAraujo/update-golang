# Update-Golang CLI

`update-golang` is a simple CLI tool that automates the process of updating your Go installation on Linux. It fetches the latest version (or a specified version) of Go, downloads it, removes the old version, and installs the new one.

## Features

- Automatically fetches and installs the latest Go version.
- Optionally install a specific Go version.
- Removes the existing Go installation before installing the new one.
- Simple, easy-to-use CLI interface.

## Installation

To build and install the `update-golang` tool, follow these steps:

1. Clone this repository:
    ```bash
    git clone https://github.com/yourusername/update-golang
    cd update-golang
    ```

2. Install Go if you haven't already. Instructions are available [here](https://golang.org/doc/install).

3. Build the CLI tool:
    ```bash
    go build -o update-golang
    ```

4. (Optional) Move the binary to a directory in your `$PATH`:
    ```bash
    sudo mv update-golang /usr/local/bin/
    ```

## Usage

### Fetch and install the latest Go version

By default, the tool fetches the latest version of Go from the official Go website and installs it:

```bash
update-golang
```

### Install a specific Go version

You can specify a specific Go version to install using the `--version` or `-v` flag:

```bash
update-golang --version go1.18
```

### Example

To install the latest Go version:

```bash
./update-golang
```

To install a specific version:

```bash
./update-golang --version go1.18
```

## How It Works

1. **Fetch Version**: The tool scrapes the Go downloads page to find the latest Go version if none is provided.
2. **Download**: It downloads the tarball file for the specified Go version.
3. **Remove Old Go**: The current Go installation (typically located in `/usr/local/go`) is removed.
4. **Install New Go**: The new Go version is extracted and installed in `/usr/local/go`.

## Requirements

- **Linux**: This tool is designed to work on Linux distributions.
- **Go**: You need Go installed to build the tool.

## Development

### Running the CLI locally

To run the tool without building a binary:

```bash
go run main.go
```

### Build

To build the tool as a binary:

```bash
go build -o update-golang
```

### Testing

You can pass different Go versions to test the installation:

```bash
./update-golang --version go1.16.5
```

## Contributing

Pull requests and issues are welcome! If you have any suggestions or improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

This `README.md` covers:

- Introduction and features.
- Installation and usage instructions.
- How the tool works.
- Contribution guidelines.