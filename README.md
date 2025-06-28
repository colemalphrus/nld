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

## Key Components
- **JSON Schema specification** with versioning and extension system
- **Go validator CLI** (`nld validate`, `nld convert`, `nld init`)
- **Go library package** for programmatic document manipulation
- **Web reader application** demonstrating semantic document features

## Getting Started
```bash
# Install the NLD tool
go install github.com/colemalphrus/nld/cmd/nld@latest

# Validate an NLD document
nld validate path/to/document.nld

# Convert a document to NLD format
nld convert path/to/document.pdf
```

## License
This project is licensed under [LICENSE] - see the LICENSE file for details.

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
