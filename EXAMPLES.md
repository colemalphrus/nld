# NLD Usage Examples

This document provides examples of common usage patterns for the NLD tool.

## Validation Examples

### Basic Validation

Validate a single document:

```bash
nld validate examples/valid-contract.json
```

Expected output:
```
✓ examples/valid-contract.json is valid
```

### Validating Multiple Files

Validate multiple documents at once:

```bash
nld validate examples/valid-contract.json examples/valid-receipt.json
```

Expected output:
```
✓ examples/valid-contract.json is valid
✓ examples/valid-receipt.json is valid

Validation summary: 2 valid, 0 invalid
```

### Handling Invalid Documents

When validating an invalid document:

```bash
nld validate examples/invalid-missing-fields.json
```

Expected output:
```
✗ examples/invalid-missing-fields.json has 2 errors:
  - Line 5: metadata.title is required
  - Line 7: content.sections is required
```

### Using a Custom Schema

Validate a document against a specific schema:

```bash
nld validate --schema schemas/nda.schema.json examples/nda.json
```

### Verbose Output

Get more detailed validation information:

```bash
nld validate --verbose examples/valid-contract.json
```

Expected output:
```
Validating file: examples/valid-contract.json
✓ examples/valid-contract.json is valid
```

### JSON Output Format

Get validation results in JSON format:

```bash
nld validate --output-format json examples/valid-contract.json
```

Expected output:
```json
{
  "Valid": true,
  "Errors": [],
  "Warnings": []
}
```

## Document Creation Examples

### Creating a Basic Contract

Create a new contract document:

```bash
nld init --type contract --output my-contract.json
```

Expected output:
```
Created new contract document: my-contract.json
```

### Creating a Receipt

Create a receipt document:

```bash
nld init --type receipt --output my-receipt.json
```

### Interactive Document Creation

Create a document with interactive prompts:

```bash
nld init --interactive --output my-document.json
```

Example interaction:
```
Enter document metadata (press Enter to use default):
Title [New contract]: Service Agreement
Author: John Doe
Jurisdiction: California, USA
Created new contract document: my-document.json
```

### Overwriting Existing Files

Force overwrite of an existing file:

```bash
nld init --force --output existing-file.json
```

### Custom Document Title

Create a document with a specific title:

```bash
nld init --title "Software Development Agreement" --output sda.json
```

## Combining Commands

Validate a document immediately after creation:

```bash
nld init --type contract --output new-contract.json && nld validate new-contract.json
```

Expected output:
```
Created new contract document: new-contract.json
✓ new-contract.json is valid
```

## Document Structure Examples

### Minimal Valid Document

```json
{
  "metadata": {
    "version": "1.0.0",
    "type": "contract",
    "created": "2025-06-27T12:00:00Z",
    "title": "Minimal Contract"
  },
  "content": {
    "sections": [
      {
        "id": "section1",
        "title": "Section 1",
        "content": "This is the content of section 1."
      }
    ]
  }
}
```

### Complete Contract Example

```json
{
  "metadata": {
    "version": "1.0.0",
    "type": "contract",
    "created": "2025-06-27T12:00:00Z",
    "title": "Service Agreement",
    "author": "Legal Department",
    "entities": [
      {
        "id": "provider",
        "name": "Acme Corporation",
        "role": "Service Provider"
      },
      {
        "id": "client",
        "name": "XYZ Ltd",
        "role": "Client"
      }
    ],
    "jurisdiction": "California, USA"
  },
  "content": {
    "sections": [
      {
        "id": "introduction",
        "title": "Introduction",
        "content": "This Service Agreement is entered into by and between Acme Corporation and XYZ Ltd."
      },
      {
        "id": "scope",
        "title": "Scope of Services",
        "content": "Acme Corporation will provide the following services..."
      },
      {
        "id": "terms",
        "title": "Terms and Conditions",
        "content": "The following terms and conditions apply to this agreement..."
      }
    ]
  },
  "relationships": {
    "references": [
      {
        "source": "terms",
        "target": "scope",
        "type": "references"
      }
    ]
  }
}