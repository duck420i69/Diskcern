# Contributing to Diskcern

First off, thank you for considering contributing to Diskcern! It's people like you that make Diskcern a great tool.

## Getting Started

1. **Fork the repository** on GitHub.
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/diskcern.git
   cd diskcern
   ```
3. **Ensure you have Go 1.24+** installed.

## Development Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/my-awesome-feature
   ```
2. Make your changes in the code.
3. **Test your changes**. Ensure all existing and new tests pass:
   ```bash
   go test ./...
   ```
4. If you added a new provider, please ensure you include corresponding tests and update the documentation if necessary.

## Submitting Changes

1. Commit your changes with a clear and descriptive commit message:
   ```bash
   git commit -m "Add feature XYZ"
   ```
2. Push your branch to your fork on GitHub:
   ```bash
   git push origin feature/my-awesome-feature
   ```
3. Open a **Pull Request** (PR) against the `main` branch of the original `diskcern` repository.
4. Provide a clear description of the problem you are solving and how you solved it.

## Code Style

* We adhere to standard Go formatting. Please run `go fmt ./...` before committing.
* Write unit tests for new features.
* Keep the CLI commands and the Gio UI in sync where applicable.

We look forward to reviewing your PRs!
