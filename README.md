# gh-prs

A GitHub CLI extension for interactively listing and opening Pull Requests based on review status.

## Features

- 🔍 **Smart PR Discovery**: Automatically detects your organization or specify with `--org`
- 📝 **Review Status Filtering**: View PRs you need to review or have already reviewed
- 🎯 **Interactive Selection**: Beautiful terminal UI for browsing and selecting PRs
- 🌐 **Browser Integration**: Opens selected PRs directly in your default browser
- ⚡ **Fast & Lightweight**: Built in Go for optimal performance

## Installation

### Prerequisites
- [GitHub CLI](https://cli.github.com/) must be installed and authenticated

### Install from GitHub Releases

```bash
gh extension install your-username/gh-prs
```

### Install from Source

```bash
git clone https://github.com/your-username/gh-prs.git
cd gh-prs
gh extension install .
```

## Usage

### Basic Commands

```bash
# List PRs waiting for your review
gh prs to-review

# List PRs you have already reviewed
gh prs reviewed

# Specify an organization
gh prs to-review --org myorg
gh prs reviewed --org myorg
```

### Options

- `--org, -o`: Specify GitHub organization (defaults to current repo's organization)
- `--help, -h`: Show help information
- `--version`: Show version information

### Interactive Controls

- `↑/↓` or `k/j`: Navigate through PR list
- `Enter`: Open selected PR in browser
- `q` or `Esc` or `Ctrl+C`: Quit

## Examples

```bash
# Review PRs in your current organization
gh prs to-review

# Check PRs you've reviewed in a specific org
gh prs reviewed --org acme-corp

# Get help
gh prs --help
```

## Development

### Prerequisites
- Go 1.21 or later
- Make (optional, for build scripts)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally for testing
make local-install

# Run tests
make test
```

### Project Structure

```
gh-prs/
├── main.go              # Main application code
├── go.mod               # Go module definition
├── Makefile             # Build scripts
├── README.md            # This file
└── .github/
    └── workflows/       # GitHub Actions for CI/CD
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [Bubble Tea](https://github.com/charmbracelet/bubbletea) for interactive UI
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss)