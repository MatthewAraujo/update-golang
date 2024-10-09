# Update-Golang CLI

`update-golang` is a simple CLI tool that automates the process of updating your Go installation on Linux. It fetches the latest version (or a specified version) of Go, removes the old version, and installs the new one.

## Demo Video

[Watch the demo video](https://drive.google.com/file/d/1RfJe9IazW80lA8Os4hovXiKU3-yyzZmr/view?usp=drive_link)

## Features

- Automatically fetches and installs the latest Go version.
- Optionally install a specific Go version.
- Removes the existing Go installation before installing the new one.
- Simple, easy-to-use CLI interface.

## Installation

To install the `update-golang` tool, use the following command:

```bash
go install github.com/MatthewAraujo/update-golang@latest
```

This command will place the `update-golang` binary in your `$GOPATH/bin` directory, which should be included in your system's `PATH`.

You need to run the following command with `sudo` because the tool installs the Go version into the `/usr/local/bin` directory, which requires elevated permissions to access:

```bash
sudo mv ~/go/bin/update-golang /usr/local/bin/
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
update-golang
```

To install a specific version:

```bash
update-golang --version go1.18
```

## How It Works

1. **Fetch Version**: The tool scrapes the Go downloads page to find the latest Go version if none is provided.
2. **Download**: It downloads the tarball file for the specified Go version.
3. **Remove Old Go**: The current Go installation (typically located in `/usr/local/go`) is removed.
4. **Install New Go**: The new Go version is extracted and installed in `/usr/local/go`.

## Requirements

- **Linux**: This tool is designed to work on Linux distributions.
- **Go**: You need Go installed to use the tool.

## Contributing

Pull requests and issues are welcome! If you have any suggestions or improvements, feel free to open an issue or submit a pull request.