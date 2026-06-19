# Contributing to Polaris

First off, thank you for considering contributing to Polaris! It's people like you that make Polaris such a great tool.

## How Can I Contribute?

### Reporting Bugs

If you find a bug in the source code, you can help us by [submitting an issue](#) to our GitHub Repository. Even better, you can submit a Pull Request with a fix.

Before submitting a bug report, please check the existing issues to avoid duplicates.

### Suggesting Enhancements

If you want to suggest an enhancement, please submit an issue with:
* A clear and descriptive title
* A detailed description of the proposed enhancement
* Any relevant examples or screenshots

### Pull Requests

1. Fork the repository and create your branch from `main`.
2. Ensure your code follows the standard Go formatting guidelines by running `go fmt ./...`.
3. If you've changed APIs, update the documentation.
4. Make sure your code passes standard linters (e.g., `golangci-lint`).
5. Issue that pull request!

## Local Development Setup

Ensure you have Go installed on your system.

```bash
# Clone the repository
git clone https://github.com/yourusername/polaris.git
cd polaris

# Install dependencies
go mod download

# Build the project
go build -o polaris cmd/polaris/main.go
```

## Coding Style

* Follow the standard Go conventions.
* Format your code with `gofmt` or `goimports` before committing.
* Add GoDoc comments to exported functions, types, and packages.

## Commit Messages

* Use the present tense ("Add feature" not "Added feature").
* Limit the first line to 72 characters or less.
* Reference issues and pull requests liberally after the first line.

We look forward to your contributions!
