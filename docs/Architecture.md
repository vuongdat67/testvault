# FileVault Architecture Design

## System Architecture Overview

FileVault follows a layered architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    FileVault Architecture                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                CLI Layer (Cobra)                      │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────────┐ │  │
│  │  │ encrypt │ │ decrypt │ │  info   │ │ help/version │ │  │
│  │  │ command │ │ command │ │ command │ │   commands   │ │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └──────────────┘ │  │
│  └───────────────────────┬───────────────────────────────┘  │
│                          │                                  │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │              Business Logic Layer                     │  │
│  │  ┌─────────────┐ ┌──────────────┐ ┌─────────────────┐ │  │
│  │  │   Encrypt   │ │   Decrypt    │ │  File Analysis  │ │  │
│  │  │ Orchestrator│ │ Orchestrator │ │   Operations    │ │  │
│  │  └─────────────┘ └──────────────┘ └─────────────────┘ │  │
│  └───────────────────────┬───────────────────────────────┘  │
│                          │                                  │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │                Core Services Layer                    │  │
│  │  ┌─────────────┐ ┌──────────────┐ ┌─────────────────┐ │  │
│  │  │   Crypto    │ │   FileOps    │ │    Security     │ │  │
│  │  │   Service   │ │   Service    │ │    Service      │ │  │
│  │  │             │ │              │ │                 │ │  │
│  │  │ • AES-GCM   │ │ • Stream I/O │ │ • Password      │ │  │
│  │  │ • PBKDF2    │ │ • Format     │ │ • Memory Mgmt   │ │  │
│  │  │ • Random    │ │ • Validation │ │ • Input Valid   │ │  │
│  │  └─────────────┘ └──────────────┘ └─────────────────┘ │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Module Breakdown

### CLI Layer (`cmd/filevault/`, `internal/cli/`)
- **Responsibility**: User interface and command parsing
- **Components**: 
  - Cobra-based command structure
  - Input validation and user prompts
  - Progress indicators and messaging
- **Key Files**: `main.go`, `commands/*.go`, `messages.go`, `progress.go`

### Core Business Logic (`internal/core/`)
- **Responsibility**: Orchestrate encryption/decryption workflows
- **Components**:
  - File encryption orchestration
  - File decryption orchestration  
  - File verification operations
- **Key Files**: `encrypt.go`, `decrypt.go`, `verify.go`

### Cryptographic Services (`internal/crypto/`)
- **Responsibility**: All cryptographic operations
- **Components**:
  - AES-256-GCM symmetric encryption
  - PBKDF2 key derivation function
  - Secure random number generation
  - Cryptographic type definitions
- **Key Files**: `aes.go`, `kdf.go`, `random.go`, `types.go`

### File Operations (`internal/fileops/`)
- **Responsibility**: File I/O and format handling
- **Components**:
  - Streaming file reader/writer
  - Binary format specification
  - File format validation
- **Key Files**: `reader.go`, `writer.go`, `format.go`

### Security Services (`internal/security/`)
- **Responsibility**: Security controls and validations
- **Components**:
  - Password strength validation
  - Input sanitization and validation
  - Secure memory management
- **Key Files**: `password.go`, `validation.go`, `memory.go`

### Public API (`pkg/filevault/`)
- **Responsibility**: Clean external interface
- **Components**: Simple client API for programmatic usage
- **Key Files**: `client.go`

### Error Handling (`internal/errors/`)
- **Responsibility**: Centralized error management
- **Components**: Structured error types with user-friendly messages

## Data Flow

### Encryption Flow
1. **CLI Input** → Password prompt and file validation
2. **Security Layer** → Input validation and password strength check
3. **Crypto Layer** → Salt generation and key derivation
4. **Core Logic** → Orchestrate encryption process
5. **FileOps Layer** → Stream processing and binary format creation

### Decryption Flow  
1. **CLI Input** → Password prompt and file format check
2. **FileOps Layer** → Header parsing and format validation
3. **Security Layer** → Authentication tag verification
4. **Crypto Layer** → Key derivation and decryption
5. **Core Logic** → Stream decryption to output file

## Security Architecture

### Defense in Depth
1. **Input Validation**: All user inputs are validated
2. **Authentication**: GCM provides built-in authentication
3. **Key Management**: PBKDF2 with high iteration count
4. **Memory Protection**: Secure memory clearing
5. **File Protection**: Proper file permissions

### Security Boundaries
- User input validation at CLI layer
- Cryptographic operations isolated in crypto module
- File operations controlled through validation layer
- Memory management centralized in security module

## Design Patterns

### Strategy Pattern
Used for cryptographic algorithms to allow future extensibility.

### Factory Pattern  
Used for file format handling to support multiple versions.

### Observer Pattern
Used for progress reporting during long operations.

### Command Pattern
Used for CLI command structure and validation.

## Performance Considerations

### Streaming Processing
- Files processed in 64KB chunks to handle large files
- Memory usage remains constant regardless of file size
- Progress reporting for user feedback

### Optimization Points
- PBKDF2 iterations balanced for security vs performance
- Buffer sizes optimized for typical file operations
- Lazy loading of cryptographic contexts

## Future Extensibility

### Planned Extension Points
1. **Additional Algorithms**: ChaCha20-Poly1305 support
2. **Compression**: Pre-encryption compression options
3. **Key Management**: Hardware security module integration
4. **Batch Operations**: Enhanced parallel processing

### Modularity Benefits
- Each layer can be tested independently
- Cryptographic algorithms can be swapped easily
- CLI commands can be added without core changes
- File formats can be versioned and upgraded
