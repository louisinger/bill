# bill

A command-line tool to generate professional PDF invoices with Bitcoin payment support.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- Interactive CLI for invoice generation
- Template support for recurring information
- Professional PDF output with:
  - Company and client details
  - Itemized billing
  - Bitcoin payment section with QR code
- EUR currency support
- Clean and modern design

## Installation

### From Source

Requirements:
- Go 1.21 or later
- make
- git

### Build from Source

```bash
make build
```

### Install

```bash
make install
```

## Usage

```bash
bill --help
```

### Template

The template is a JSON file that contains recurring information of the invoice.

```json
{
  "company_name": "My Company",
  "address": "123 Business Street\nCity, 12345\nCountry",
  "vat_number": "GB123456789",
  "to_company_name": "Client Company Ltd",
  "to_address": "456 Client Street\nClient City, 54321\nClient Country",
  "to_vat_number": "FR987654321",
  "items": [
      {
          "description": "ask chatGPT",
          "quantity": 1,
          "unit_price": 100000.00
      }
  ]
}
```

Then, use the template to avoid repeating yourself.

```bash
bill --template template.json
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.