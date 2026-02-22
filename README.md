# BulkCaller

Bulk HTTP request tool that substitutes values from CSV/XLS/XLSX files into request bodies.

## Features

- 📁 Supports CSV, XLS, and XLSX files
- 🔄 Replace placeholders in JSON body templates
- ⚡ Concurrent workers for high throughput
- 🔄 Automatic retries with exponential backoff
- 📊 Progress logging
- 💾 Optional output saving

## Installation

```bash
go install github.com/wellrcosta/bulkcaller@latest
```

Or download from releases.

## Usage

```bash
./bulkcaller \
  -file data.csv \
  -url "https://api.example.com/webhook" \
  -method POST \
  -body '{"name":"{{name}}","email":"{{email}}"}' \
  -headers "Authorization:Bearer token" \
  -concurrency 20 \
  -delay 100
```

## Parameters

| Flag | Description | Default |
|------|-------------|---------|
| `-file` | Path to data file (csv/xls/xlsx) | required |
| `-url` | Target URL | required |
| `-body` | JSON template with placeholders | required |
| `-method` | HTTP method | POST |
| `-headers` | Headers (key:value, comma-separated) | - |
| `-query` | Query params (key=value, comma-separated) | - |
| `-concurrency` | Number of workers | 10 |
| `-delay` | Delay between requests (ms) | 0 |
| `-retries` | Max retries on failure | 3 |
| `-output` | Output directory for responses | - |
| `-print` | Print responses to stdout | false |

## Input File Format

First row must contain headers that match placeholders in the body template.

**Example CSV:**
```csv
name,email,phone
John Doe,john@example.com,+5511999999999
Jane Smith,jane@example.com,+5511888888888
```

**Example Body Template:**
```json
{"customer_name":"{{name}}","contact_email":"{{email}}","contact_phone":"{{phone}}"}
```

## License

MIT