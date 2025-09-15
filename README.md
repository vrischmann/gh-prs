# gh-prs

A GitHub CLI extension for interactively listing and opening open Pull Requests based on review status.

## Features

- Automatically detects your organization or specify with `--org`
- Filter open PRs by review status (need to review vs already reviewed)
- Interactive terminal UI for browsing and selecting PRs
- Opens selected PRs directly in your default browser

## Installation

### Prerequisites
- [GitHub CLI](https://cli.github.com/) must be installed and authenticated

### Install from GitHub Releases

```bash
gh extension install vrischmann/gh-prs
```

### Install from Source

```bash
git clone https://github.com/vrischmann/gh-prs.git
cd gh-prs
gh extension install .
```

## Usage

### Basic Commands

```bash
# List open PRs waiting for your review
gh prs to-review

# List open PRs you have already reviewed
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


## Development

### Prerequisites
- Go 1.25 or later
- [just](https://github.com/casey/just) (optional, for build scripts)

### Building

```bash
# Build for current platform
just build

# Build for all platforms
just build-all

# Install locally for testing
just local-install

# Run tests
just test
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
