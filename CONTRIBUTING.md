# Contributing to BX

Thank you for your interest in contributing to [BX](https://github.com/pixel365/bx)! Please follow these guidelines to ensure your contributions are helpful and easy to integrate.

## General Guidelines

- Install the following tools for testing and verification before contributing:
    - [field alignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment)
    - [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
    - [golines](https://github.com/segmentio/golines)
    - [golangci-lint](https://github.com/golangci/golangci-lint)
- Read the [README.md](README.md) and documentation before starting.
- Ensure your code follows Go's style (`gofmt`, `golangci-lint` and `golines`).
- Keep the code readable and add comments where necessary.
- Open an issue before making significant changes.

## How to Contribute

1. **Fork the repository** and create a new branch:
   ```sh
   git checkout -b feature/my-feature
   ```
2. **Develop your changes**, ensuring tests and documentation are updated.
3. **Run checks** before committing:
   ```sh
   make
   ```
4. **Commit your changes** with a meaningful message:
   ```sh
   git commit -m "feat: add new functionality"
   ```
5. **Push your changes** to your fork and create a Pull Request (PR).

## Pull Request Requirements

- PRs should be small and focused on a single task.
- Provide a description of changes and link to the relevant issue (if applicable).
- PRs must pass all automated checks (CI/CD) before merging.
- Ensure your code is covered by tests where necessary.

## Working with Issues

- Before opening a new issue, check if it already exists.
- Provide a clear problem description with examples if possible.
- For bugs, include your Go version, system details, reproduction steps, and expected behavior.

Thank you for contributing to [BX](https://github.com/pixel365/bx)! 🚀
