# EzLang - PO File Translation CLI Tool

EzLang is a command-line tool that automates the translation of PO (Portable Object) files using AWS Translate. It helps developers and translators manage internationalization by automatically translating untranslated entries, handling fuzzy translations, and maintaining proper PO file formatting.

## Features

- 🌍 Automatic translation using AWS Translate
- 🔄 Handles fuzzy translations (removes fuzzy flags and retranslates)
- 📝 Preserves PO file structure and formatting
- 🔀 Supports plural forms
- 📄 Handles multi-line translations
- ⚙️ Configuration via YAML file
- 🚀 High performance with proper error handling

## Prerequisites

- Go 1.16 or higher
- AWS Account with Translate service access
- AWS Access Key ID and Secret Access Key
- AWS Region where the Translate service is available

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ezlang.git
cd ezlang
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o ezlang cmd/cli/main.go
```

## Configuration

### AWS Credentials

Set up your AWS credentials using environment variables:

```bash
export AWS_ACCESS_KEY_ID=your_access_key_here
export AWS_SECRET_ACCESS_KEY=your_secret_key_here
export AWS_REGION=ap-south-1 # or your preferred region
```

Alternatively, you can use AWS credentials file (`~/.aws/credentials`):

```ini
[default]
aws_access_key_id = your_access_key_here
aws_secret_access_key = your_secret_key_here
```

### EzLang Configuration

Create a `.ezlang.yml` file in your project root:

```yaml
files:
  - path: "locale/en/LC_MESSAGES/django.po"
    lang: "en"
  - path: "locale/sv/LC_MESSAGES/django.po"
    lang: "sv"
```

### Basic Usage

Run with default configuration (`.ezlang.yml`):
```bash
./ezlang
```

Run with custom configuration:
```bash
./ezlang --config custom-config.yml
```

### Command Line Options

```bash
./ezlang [options]

Options:
  --config string     Path to configuration file (default ".ezlang.yml")
  --rewrite string    Rewrite main file. Default is false
  --help              Show help message
```

### Example

```bash
# Translate PO files defined in .ezlang.yml
./ezlang

# Use custom configuration
./ezlang --config translations.yml

# Override output directory
./ezlang --rewrite
```

## How It Works

1. **Reads PO Files**: Parses PO files while preserving their structure
2. **Identifies Translations**: Finds empty msgstr entries and fuzzy translations
3. **Translates Content**: Uses AWS Translate to translate msgid to target language
4. **Preserves Formatting**: Maintains original file structure, comments, and metadata
5. **Handles Special Cases**: Properly processes plural forms and multi-line strings
6. **Removes Fuzzy Flags**: Cleans up fuzzy markers after translation

## Development

### Project Structure

```
ezlang/
├── .ezlang.yml
├── .gitignore
├── LICENSE
├── README.md
├── cmd
│   └── cli
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── application
│   │   └── services
│   │   └── file_service.go
│   └── config
│       └── config.go
└── pkg
    ├── po
    │   └── po.go
    └── translation
        └── translation.go
```

### Running Tests
Not implemented yet.
```bash
go test ./...
```

### Building from Source

```bash
# Build for current platform
go build -o ezlang cmd/cli/main.go

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o ezlang-linux cmd/cli/main.go
GOOS=darwin GOARCH=amd64 go build -o ezlang-mac cmd/cli/main.go
GOOS=windows GOARCH=amd64 go build -o ezlang.exe cmd/cli/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- AWS Translate for providing the translation service
- The Go community for excellent libraries and tools

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.