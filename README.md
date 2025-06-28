# NLD - Next-Gen Layout Document

## Project Overview
NLD (Next-Gen Layout Document) is an open-source document standard and ecosystem designed to replace PDF for structured documents like legal contracts, receipts, invoices, and agreements. The project consists of a JSON-based schema specification and reference implementation tools built in Go.

## Problem Statement
- PDF parsing is unreliable and complex
- Modern documents are not LLM-friendly or machine-readable
- Legal and business documents lack semantic structure
- Document workflows require expensive, proprietary tools
- No open standard exists for structured business documents

## Project Goals
1. **Define an open document schema** that captures semantic structure, not just layout
2. **Build a Go-based validator** as the reference implementation
3. **Create a web-based reader** that demonstrates the format's capabilities
4. **Enable ecosystem adoption** through open standards and tooling

## Installation

### From Source
```bash
# Clone the repository
git clone https://github.com/colemalphrus/nld.git
cd nld

# Build the tool
go build ./cmd/nld

# Install the tool
go install ./cmd/nld
```

### Using Go Install
```bash
# Install the NLD tool directly
go install github.com/colemalphrus/nld/cmd/nld@latest
```

## Usage

### Validating Documents
Validate a single NLD document against its schema:
```bash
nld validate path/to/document.json
```

Validate multiple documents at once:
```bash
nld validate doc1.json doc2.json doc3.json
```

Use a specific schema file:
```bash
nld validate --schema path/to/schema.json document.json
```

Additional options:
- `--verbose` or `-v`: Show detailed validation information
- `--quiet` or `-q`: Suppress all output except errors
- `--output-format` or `-o`: Output format (text, json)
- `--force` or `-f`: Continue validation even if some files fail

### Creating New Documents
Create a new document using a template:
```bash
nld init --type contract --output my-contract.json
```

Available document types:
- `contract`: Basic contract template
- `receipt`: Receipt template
- `agreement`: General agreement template

Interactive mode for guided document creation:
```bash
nld init --interactive
```

Additional options:
- `--title`: Set the document title
- `--force` or `-f`: Overwrite existing files

### Version Information
Display version information:
```bash
nld version
```

## Document Schema
NLD documents follow a structured JSON schema with the following main components:

```json
{
  "metadata": {
    "version": "1.0.0",
    "type": "contract",
    "created": "2025-06-27T12:00:00Z",
    "title": "Service Agreement"
  },
  "content": {
    "sections": [
      {
        "id": "introduction",
        "title": "Introduction",
        "content": "This is the introduction section."
      }
    ]
  }
}
```

### Required Fields
- `metadata.version`: Document schema version (e.g., "1.0.0")
- `metadata.type`: Document type (contract, receipt, agreement)
- `metadata.created`: Creation date in ISO format
- `metadata.title`: Document title
- `content.sections`: Array of document sections

### Optional Fields
- `metadata.author`: Document author
- `metadata.entities`: Entities involved in the document
- `metadata.jurisdiction`: Legal jurisdiction
- `relationships`: Document relationships and dependencies

## Examples
Example documents can be found in the `examples/` directory:
- `valid-contract.json`: A complete contract example
- `valid-receipt.json`: A receipt example
- `invalid-missing-fields.json`: Example with missing required fields
- `invalid-type.json`: Example with invalid document type

## License
This project is licensed under [LICENSE] - see the LICENSE file for details.

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
