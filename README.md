# Tidbits

`tidbits` is a collection of functions, patterns, routines, and sample applications built for self-education. This repository serves as a personal laboratory for exploring Go language features, algorithms, and software design patterns.

## Purpose

The primary goal of this repository is to provide a compilation of code snippets that help in understanding:
- **Algorithms & Data Structures**: Implementation of common structures (trees, queues, lists) and algorithmic challenges.
- **Go Language Features**: Practical examples of Go-specific concepts like concurrency, embedding, and error handling.
- **Design Patterns**: Real-world applications of patterns like Factory, Outbox, and Event Bus.
- **Tooling & Integration**: Examples of CLI building (Cobra/Viper), Docker integration, and API consumption.

## Project Structure

The repository is organized into several modules, each focusing on a specific area:

- **`arrays/`**: Common array manipulation and sorting algorithms.
- **`cmd/`**: CLI command implementations for various sample apps.
- **`concurrency/`**: Exploring Go's concurrency primitives and potential pitfalls (e.g., leaking goroutines).
- **`dynamic/`**: Dynamic programming problems (e.g., coin change).
- **`embed/`**: Usage of Go's `embed` package for static assets and templates.
- **`errorh/`**: Strategies for robust error handling.
- **`event/`**: Implementation of an event bus and event store.
- **`factory/`**: The Factory design pattern in action.
- **`structures/`**: Core data structures including Binary Trees, Priority Queues, and Linked Lists.
- **`toy/`**: Miscellaneous small puzzles and utility functions (GCD, Oddword, etc.).
- **`weather/`**: A practical example of integrating with external APIs.

## Getting Started

### Prerequisites

- Go 1.26 or later

### Running Tidbits

The project uses a CLI interface powered by [Cobra](https://github.com/spf13/cobra). You can explore the various "tidbits" by running the main entry point.

```bash
# Build the application
go build -o tidbits main.go

# Run the help command to see available apps
./tidbits --help
```

### Examples

Many directories contain their own `README.md` or a suite of tests that demonstrate the functionality:

- **Data Structures**: See `structures/` for implementations of trees and lists.
- **Patterns**: Check `factory/` or `event/` for design pattern examples.
- **Tests**: Run `go test ./...` to see all routines in action.

## License

This project is licensed under the [LICENSE](LICENSE) file found in the root directory.
