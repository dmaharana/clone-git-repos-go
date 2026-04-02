# Git Repository Bulk Cloner

A Go-based CLI tool that reads repository URLs from a CSV file and clones them in parallel. This tool is useful for managing multiple Git repositories, especially when migrating repositories between different Git hosting services.

## Features

- Bulk clone Git repositories from a CSV file
- **Smart Protocol Selection**: Automatically handles HTTPS and SSH repository URLs
- **Automatic HTTPS Conversion**: If an authentication token is provided, the tool automatically converts SSH URLs to HTTPS for seamless authentication
- **Authentication**: Native support for private repositories using username and token
- **Security & Privacy**: Automatic masking of sensitive credentials (usernames and tokens) in all logs and console output
- Parallel repository cloning
- Automatic retry mechanism for failed clones (up to 3 retries)
- Error handling for common Git operations
- Branch and Tag checkout: Clones all available branches and tags
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
- **Authentication**: Identifies credential issues and provides secure feedback
- **Conflict Management**: Automatically handles cases where the repository directory already exists
- **Resilience**: Includes an automatic retry mechanism for transient network issues
- **Invalid URLs**: Validates repository URLs before attempting operations

## Security

Security is a priority for this tool:
- **Credential Masking**: Usernames and tokens are never written to logs or displayed in the console in plain text. They are replaced with `****` automatically.
- **Secure Transport**: When a token is provided, the tool defaults to HTTPS to ensure encrypted communication with the Git provider.
- **No In-URL Credentials**: Credentials are passed via secure authentication headers rather than being embedded directly in URLs.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Development

This project was developed using [Codeium](https://codeium.com/) and Windsurf, leveraging AI-powered development tools for enhanced productivity and code quality. Fixes implemented with Google Gemini CLI.
