# Diskcern Architecture

Diskcern is built using Go, emphasizing performance, modularity, and cross-platform compatibility. The system is designed to handle rapid scanning of file systems and extracting insights using a pluggable provider model.

## High-Level Components

* **CLI (Command-Line Interface):** Built with [Cobra](https://github.com/spf13/cobra), handling initialization, command parsing, and directing flow to either the CLI commands (like `scan`, `diff`) or launching the UI. Located in `cmd/`.
* **UI (Graphical User Interface):** A fast, native GUI built with [Gio](https://gioui.org/). Gio allows Diskcern to compile to a native binary without relying on web views or heavy UI frameworks. Located in `internal/ui/`.
* **Scanner & Engine:** The core traversing logic. It rapidly scans directories, optionally running in parallel, and processes metadata. Located in `internal/engine/` and `internal/scanner/`.
* **Database (DB):** Diskcern uses [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (a pure Go SQLite implementation) to store snapshot data, enabling fast queries and offline analysis without CGO dependencies. Located in `internal/db/`.
* **Providers:** A modular system for detecting and analyzing specific types of data (e.g., Steam libraries, Docker images, node_modules). Providers implement a common interface to evaluate file paths during a scan. Located in `internal/providers/`.

## Data Flow

1. A scan is initiated via the CLI or UI.
2. The **Scanner** traverses the file system. For each file/directory, it passes the path to the **Provider Registry**.
3. The **Registry** iterates through all registered **Providers**, calling `Detect()`.
4. If a Provider matches the path, it can analyze the contents, dictate traversal directives (e.g., stop traversing this folder), and generate insights (like cleanup actions).
5. The results and file metadata are saved to the **Database** for caching, querying, and diffing.
