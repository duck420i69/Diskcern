# Diskcern

**Diskcern** is a generic, cross-platform storage analyzer designed for finding large files, analyzing disk usage, and creating snapshots. It is built in Go and provides both a Command-Line Interface (CLI) and a Graphical User Interface (GUI).

## Core Features
* **Cross-platform Graphical User Interface:** A fast, native GUI built with [Gio](https://gioui.org/).
* **CLI for Automation:** Command-line tools for scanning directories and comparing snapshots.
* **Storage Analysis:** Quickly find large files and visualize disk usage.

## Target Audience
This documentation is intended for both **end-users** looking to install and use the application, as well as **developers** who want to build the project from source or contribute to its development.

## Basic Usage

### Building from Source

To build Diskcern from source, ensure you have Go 1.24+ installed on your system.

```bash
git clone https://github.com/diskcern/diskcern.git
cd diskcern
go build -o diskcern ./cmd/diskcern
```

### Running the CLI

You can use the CLI to scan directories or compare snapshots directly from your terminal:

* **Scan a directory:**
  ```bash
  ./diskcern scan /path/to/directory
  ```
* **Compare two snapshots:**
  ```bash
  ./diskcern diff snapshot1 snapshot2
  ```

### Launching the UI

If you run Diskcern without any specific CLI commands, it will automatically launch the Graphical User Interface:

```bash
./diskcern
```

## Documentation

More detailed documentation will be available in the `docs` directory soon.
