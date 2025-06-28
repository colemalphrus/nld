# Next-Gen Layout Document (.nld) Project Brief

## Project Overview
Create an open-source document standard and ecosystem to replace PDF for structured documents like legal contracts, receipts, invoices, and agreements. The project consists of a JSON-based schema specification and reference implementation tools built in Go.

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

## Technical Architecture

### Core Schema Design
```
Document {
  metadata: { type, version, created, entities, jurisdiction }
  structure: { sections[], items[], definitions[] }
  relationships: { dependencies, references, conditions }
  verification: { signatures, timestamps, attestations }
}
```

### Document Types (Initial Support)
- Legal contracts (NDAs, service agreements)
- Purchase agreements
- Receipts and invoices
- Insurance policies

### Go Validator Implementation
- **CLI tool** for document validation and conversion
- **Library package** for integration into other Go applications
- **HTTP API** for web service integration
- **JSON Schema validation** using established Go libraries

## Technical Stack

### Validator (Go)
- **Core libraries**: `encoding/json`, `github.com/santhosh-tekuri/jsonschema`
- **CLI framework**: `github.com/spf13/cobra`
- **Testing**: Standard Go testing + property-based testing
- **Distribution**: Single binary, cross-platform compilation

### Reader (Web)
- **Frontend**: React/TypeScript for rich document interaction
- **Backend**: Go HTTP server serving validation API
- **Rendering**: CSS Grid/Flexbox for responsive document layout
- **Features**: Multi-view rendering, semantic navigation, live validation

## Project Phases

### Phase 1: Foundation (Months 1-2)
- Define base document schema (JSON Schema)
- Build Go CLI validator with basic validation
- Create 3-5 example documents (NDA, receipt, invoice)
- Set up project repository and documentation

### Phase 2: Reference Implementation (Months 3-4)  
- Complete Go validator library with full schema support
- Build web-based document reader
- Add conversion tools from common formats
- Comprehensive test suite and documentation

### Phase 3: Ecosystem (Months 5-6)
- Community feedback and schema refinement
- Integration examples for popular languages
- Performance optimization and benchmarking
- Adoption outreach to legal tech and document management companies

## Success Metrics
- **Technical**: Validator handles 1000+ document validation/sec
- **Adoption**: 5+ companies experimenting with the format
- **Community**: 100+ GitHub stars, active issue discussions
- **Ecosystem**: 3+ third-party tools supporting the format

## Key Deliverables
1. **JSON Schema specification** with versioning and extension system
2. **Go validator CLI** (`nld validate`, `nld convert`, `nld init`)
3. **Go library package** for programmatic document manipulation
4. **Web reader application** demonstrating semantic document features
5. **Documentation** including specification, API docs, and integration guides
6. **Example documents** showcasing different document types and features

## Open Source Strategy
- **License**: Apache 2.0 or MIT for maximum adoption
- **Repository**: GitHub with clear contribution guidelines
- **Governance**: Benevolent dictator initially, evolving to consortium model
- **Community**: Discord/Slack for real-time discussion, GitHub issues for formal feedback

## Next Steps
1. Set up Go project structure with basic CLI skeleton
2. Define initial JSON schema for simple contracts
3. Implement core validation logic
4. Create first example document (NDA)
5. Build minimal web reader to demonstrate the concept

## Long-term Vision
Become the standard interchange format for structured business documents, enabling:
- Seamless document automation and workflow integration
- AI-powered document analysis and generation
- Vendor-neutral document ecosystems
- Dramatically improved document accessibility and processing