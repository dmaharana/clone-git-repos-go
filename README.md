# Git Repository Bulk Cloner

A Go-based CLI tool that reads repository URLs from a CSV file and clones them in parallel. This tool is useful for managing multiple Git repositories, especially when migrating repositories between different Git hosting services.

## Features

- Bulk clone Git repositories from a CSV file
- Support for both HTTPS and SSH repository URLs
- Authentication support for private repositories
- Parallel repository cloning
- Automatic retry mechanism for failed clones (up to 3 retries)
- Error handling for common Git operations
- Branch checkout for all available branches
- Configuration via INI file or command-line arguments

## Prerequisites

- Go 1.22 or higher
- Git installed on your system

## Installation

```bash
# Clone this repository
git clone https://github.com/dmaharana/clone-git-repo.git

# Change to the project directory
cd clone-git-repo

# Build the project
go build ./...

# Copy the example config file (optional)
cp config.ini.example config.ini
```

## Configuration

The tool supports two methods of configuration:

### 1. INI Configuration File (Recommended)

Copy the example configuration file and modify it according to your needs:

```bash
cp config.ini.example config.ini
```

Example `config.ini`:
```ini
[credentials]
username = your_username_here
token = your_token_here

[paths]
clone_dir = clonedir
csv_file = repositories.csv

[logging]
log_dir = logs
log_max_size = 10485760  # 10MB in bytes
```

### 2. Command Line Arguments

If no config file is found or if you prefer using command-line arguments:

- `-c`: Path to config file (default: "config.ini")
- `-f`: Path to the CSV file containing repository URLs
- `-d`: Directory where repositories will be cloned
- `-u`: Username for authentication (required for private repositories)
- `-t`: Token for authentication (required for private repositories)

## Usage

1. Create a CSV file with repository URLs (e.g., `repositories.csv`):
```csv
repo_url
https://github.com/user/repo1.git
git@github.com:user/repo2.git
```

2. Run the tool:

Using config file:
```bash
go run cmd/clone-git-repo/main.go
```

Using command line arguments:
```bash
go run cmd/clone-git-repo/main.go -f repositories.csv -d clonedir -u username -t token
```

## Error Handling

The tool includes robust error handling for common scenarios:
- Authentication errors
- Directory already exists
- Network issues
- Invalid repository URLs
- Automatic retries for failed clone attempts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Development

This project was developed using [Codeium](https://codeium.com/) and Windsurf, leveraging AI-powered development tools for enhanced productivity and code quality.
