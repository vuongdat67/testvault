# FileVault Development Guide

## Getting Started

### Prerequisites
- **Go 1.21+**: Latest stable version recommended
- **Git**: For version control
- **Make**: For build automation (optional)
- **VS Code**: Recommended editor with Go extension

### Installation

```bash
# Clone the repository
git clone https://github.com/vuongdat67/NT140.Q11.ANTT-Group15.git
cd FileVault

# Install dependencies
go mod tidy

# Build the project
go build -o filevault ./cmd/filevault

# Run tests
go test ./...
```

## Development Environment Setup

### VS Code Configuration
Recommended extensions:
- Go (golang.org)
- Go Test Explorer
- GitLens
- Better Comments

### Code Style
We follow standard Go conventions:
- Use `gofmt` for formatting
- Follow `golint` recommendations  
- Use meaningful variable names
- Keep functions small and focused
- Add comments for exported functions

## Project Structure

```
filevault/
├── cmd/filevault/           # CLI entry point
├── internal/                # Private packages
│   ├── crypto/              # Encryption algorithms
│   ├── fileops/             # File I/O operations
│   ├── security/            # Security utilities
│   ├── cli/                 # CLI handlers
│   ├── core/                # Business logic
│   └── errors/              # Error handling
├── pkg/filevault/           # Public API
├── test/                    # Test suites
│   ├── unit/                # Unit tests
│   ├── integration/         # Integration tests
│   └── benchmarks/          # Performance tests
├── docs/                    # Documentation
├── examples/                # Usage examples
└── scripts/                 # Build/deployment scripts
```

## Development Workflow

### 1. Feature Development
```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make your changes
# Write tests first (TDD approach)
go test ./test/unit/...

# Implement the feature
# Run tests frequently
go test ./...

# Format code
go fmt ./...

# Commit changes
git add .
git commit -m "feat: add your feature description"
```

### 2. Testing Strategy

#### Unit Tests
Located in `test/unit/`, these test individual modules:

```bash
# Test crypto module
go test -v ./test/unit/crypto_test.go

# Test file operations
go test -v ./test/unit/fileops_test.go

# Test security functions
go test -v ./test/unit/security_test.go
```

#### Integration Tests
Located in `test/integration/`, these test full workflows:

```bash
# Test end-to-end encryption/decryption
go test -v ./test/integration/...
```

#### Benchmarks
Located in `test/benchmarks/`:

```bash
# Run performance benchmarks
go test -bench=. ./test/benchmarks/...
```

### 3. Build and Release

#### Development Build
```bash
# Build for current platform
go build -o filevault ./cmd/filevault

# Build with debug info
go build -ldflags "-X main.version=dev" -o filevault-dev ./cmd/filevault
```

#### Cross-platform Build
```bash
# Use provided script
./scripts/build.sh

# Or manual build
GOOS=windows GOARCH=amd64 go build -o filevault-windows.exe ./cmd/filevault
GOOS=linux GOARCH=amd64 go build -o filevault-linux ./cmd/filevault
GOOS=darwin GOARCH=amd64 go build -o filevault-darwin ./cmd/filevault
```

## Code Guidelines

### 1. Error Handling
Always use the custom error types from `internal/errors/`:

```go
import "github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/errors"

// Good
return errors.NewFileNotFoundError(filename)

// Bad
return fmt.Errorf("file not found: %s", filename)
```

### 2. Security Best Practices
- Never log sensitive data (passwords, keys, file contents)
- Always clear sensitive data from memory using `security.ClearMemory()`
- Validate all user inputs using `security.Validate*()` functions
- Use constant-time comparisons for authentication data

### 3. Logging
Use structured logging with appropriate levels:

```go
// Information for users
cli.PrintInfo("Encrypting file: %s", filename)

// Warnings for non-critical issues  
cli.PrintWarning("File already exists, use --force to overwrite")

// Errors for failures
cli.PrintError("Failed to encrypt file: %v", err)
```

### 4. Testing Best Practices
- Write tests before implementing features (TDD)
- Use table-driven tests for multiple scenarios
- Test both success and failure cases
- Mock external dependencies
- Use meaningful test names

```go
func TestPasswordValidation(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
    }{
        {"valid strong password", "MyStr0ng!Pass", false},
        {"too short", "short", true},
        {"no numbers", "NoNumbers!", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := security.ValidatePasswordBasic(tt.password)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidatePasswordBasic() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Module Development Guide

### Adding New Commands
1. Create command file in `internal/cli/commands/`
2. Implement command logic
3. Add to main.go
4. Write tests
5. Update documentation

### Adding New Crypto Algorithms
1. Implement in `internal/crypto/`
2. Follow the `CryptoStrategy` interface
3. Add algorithm constants to `types.go`
4. Update file format if needed
5. Add comprehensive tests

### Adding New File Formats
1. Update `internal/fileops/format.go`
2. Implement version handling
3. Add backward compatibility
4. Update verification logic
5. Test with existing encrypted files

## Debugging

### Common Issues
1. **Import path errors**: Use the full GitHub path in imports
2. **Test failures**: Check file permissions and temp directory access
3. **Build failures**: Ensure Go version compatibility

### Debugging Tools
```bash
# Run with race detection
go run -race ./cmd/filevault

# Run with verbose output
./filevault encrypt test.txt --verbose

# Profile memory usage
go tool pprof ./filevault mem.prof
```

## Contributing

### Commit Message Format
Follow conventional commits:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions/changes
- `refactor:` Code restructuring

### Pull Request Process
1. Create feature branch from main
2. Implement feature with tests
3. Update documentation if needed
4. Create pull request with clear description
5. Address review comments
6. Merge after approval

## Performance Guidelines

### Memory Management
- Use streaming for large files (>1MB)
- Clear sensitive data immediately after use
- Monitor memory usage in benchmarks
- Use buffer pools for frequent allocations

### Optimization Tips
- Profile before optimizing
- Focus on hot paths (encryption/decryption)
- Consider parallel processing for batch operations
- Cache expensive computations when safe

## Security Development

### Threat Modeling
When adding features, consider:
- Input validation requirements
- Attack surface changes
- Authentication/authorization impacts
- Data exposure risks

### Security Review Checklist
- [ ] All inputs validated
- [ ] No sensitive data in logs
- [ ] Proper error handling
- [ ] Memory cleared after use
- [ ] Timing attack prevention
- [ ] File permissions checked

## Release Process

### Version Management
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Tag releases in git
- Update CHANGELOG.md
- Build release binaries

### Release Checklist
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Performance benchmarks run
- [ ] Security review completed
- [ ] Cross-platform builds tested
- [ ] Release notes prepared

## Troubleshooting

### Build Issues
```bash
# Clean module cache
go clean -modcache

# Verify dependencies
go mod verify

# Update dependencies
go get -u ./...
```

### Test Issues
```bash
# Run specific test
go test -run TestEncryption ./test/unit/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Runtime Issues
```bash
# Enable debug logging
export FILEVAULT_DEBUG=1
./filevault encrypt test.txt

# Check file permissions
ls -la test.txt

# Verify file format
hexdump -C test.txt.enc | head
```

## Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Go Crypto Packages](https://pkg.go.dev/golang.org/x/crypto)

### Security Resources
- [OWASP Cryptographic Standards](https://owasp.org/www-project-cryptographic-storage-cheat-sheet/)
- [Go Security Guidelines](https://golang.org/doc/security)
- [Cryptographic Right Answers](https://latacora.micro.blog/2018/04/03/cryptographic-right-answers.html)
