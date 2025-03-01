# Bill - Invoice Generator

A modern, user-friendly invoice generator with Bitcoin payment support.

## Features

- Create professional invoices with a modern UI
- Add multiple items with automatic total calculation
- Support for Bitcoin payments
- Save and load default values
- Export to PDF
- Cross-platform support (macOS, Linux, Windows)

## Installation

### Prerequisites

- Go 1.21 or later
- Git

### From Source

1. Clone the repository:
```bash
git clone https://github.com/louisinger/bill.git
cd bill
```

2. Install the application:
```bash
go install ./cmd/bill
```

The `bill` binary will be installed in your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH.

### Using Go Install (directly from GitHub)

```bash
go install github.com/louisinger/bill/cmd/bill@latest
```

## Usage

### GUI Mode

Simply run:
```bash
bill
```

### CLI Mode

Generate an invoice using command line:
```bash
bill generate -o invoice.pdf
```

Use a template:
```bash
bill generate -t template.json -o invoice.pdf
```

Show version:
```bash
bill version
```

## Development

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/louisinger/bill.git
cd bill
```

2. Build the application:
```bash
go build -o bill ./cmd/bill
```

3. Run the application:
```bash
./bill
```

## License

MIT License - see LICENSE file for details