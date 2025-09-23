# 🔐 FileVault - Secure File Encryption Tool

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Security](https://img.shields.io/badge/security-AES--256--GCM-red.svg)]()
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)]()

**FileVault** là một công cụ mã hóa file command-line đơn giản, an toàn và hiệu quả được phát triển bởi **NT140.Q11.ANTT Group 15**. Sử dụng mã hóa AES-256-GCM và key derivation PBKDF2, FileVault đảm bảo bảo mật cao cho các file quan trọng của bạn.

## 📋 Mục Lục

- [🚀 Tính Năng Chính](#-tính-năng-chính)
- [🔧 Cài Đặt](#-cài-đặt)
- [📖 Sử Dụng Cơ Bản](#-sử-dụng-cơ-bản)
- [🛡️ Bảo Mật](#️-bảo-mật)
- [🏗️ Kiến Trúc](#️-kiến-trúc)
- [🧪 Testing](#-testing)
- [🤝 Đóng Góp](#-đóng-góp)
- [📄 License](#-license)

## 🚀 Tính Năng Chính

### ✅ **Core Features**
- **🔐 Mã Hóa Mạnh**: AES-256-GCM với authenticated encryption
- **🔑 Key Derivation An Toàn**: PBKDF2-SHA256 với 100,000 iterations
- **🎲 Random Salt**: Mỗi file được mã hóa với salt unique 32-byte
- **📊 Progress Tracking**: Progress bar cho file lớn
- **🔍 File Verification**: Kiểm tra tính toàn vẹn file
- **🔄 Batch Processing**: Mã hóa nhiều file cùng lúc

### 🛠️ **Advanced Features**
- **📱 Cross-Platform**: Hỗ trợ Windows, Linux, macOS
- **🧠 Memory Security**: Secure memory cleanup sau xử lý
- **⚡ Streaming Encryption**: Xử lý file lớn hiệu quả
- **📄 Custom File Format**: Header metadata với version control
- **🎯 CLI Intuitive**: Command-line interface dễ sử dụng

## 🔧 Cài Đặt

### **Option 1: Download Binary (Recommended)**

Tải pre-built binary từ [GitHub Releases](https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/releases):

```bash
# Windows
wget https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/releases/latest/download/filevault-windows-amd64.exe

# Linux
wget https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/releases/latest/download/filevault-linux-amd64
chmod +x filevault-linux-amd64

# macOS  
wget https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/releases/latest/download/filevault-darwin-amd64
chmod +x filevault-darwin-amd64
```

### **Option 2: Build từ Source**

**Yêu cầu**: Go 1.25 hoặc cao hơn

```bash
# Clone repository
git clone https://github.com/vuongdat67/NT140.Q11.ANTT-Group15.git
cd NT140.Q11.ANTT-Group15

# Build
make build

# Hoặc build manual
go build -o filevault cmd/filevault/main.go
```

### **Option 3: Install Script**

```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/vuongdat67/NT140.Q11.ANTT-Group15/main/scripts/install.sh | bash

# Windows PowerShell
iwr -useb https://raw.githubusercontent.com/vuongdat67/NT140.Q11.ANTT-Group15/main/scripts/install.ps1 | iex
```

## 📖 Sử Dụng Cơ Bản

### **🔐 Mã Hóa File**

```bash
# Mã hóa file đơn lẻ
filevault encrypt document.pdf
# Output: document.pdf.enc

# Chỉ định file output
filevault encrypt document.pdf secure.enc

# Mã hóa với output directory
filevault encrypt document.pdf -o encrypted/
```

### **🔓 Giải Mã File**

```bash  
# Giải mã file
filevault decrypt document.pdf.enc
# Output: document.pdf

# Chỉ định file output
filevault decrypt secure.enc original.pdf

# Giải mã với output directory
filevault decrypt secure.enc -o decrypted/
```

### **📊 Thông Tin File**

```bash
# Xem thông tin file đã mã hóa
filevault info document.pdf.enc

# Kiểm tra tính toàn vẹn
filevault verify document.pdf.enc

# Kiểm tra file có bị mã hóa không
filevault check document.pdf
```

## 🚀 Sử Dụng Nâng Cao

### **📁 Batch Operations**

```bash
# Mã hóa nhiều file
filevault encrypt *.txt
filevault encrypt file1.pdf file2.docx file3.xlsx

# Mã hóa với pattern
filevault encrypt reports/*.pdf -o encrypted/
filevault encrypt "*.{txt,doc,pdf}" -o secure/
```

### **⚙️ Tùy Chọn Nâng Cao**

```bash
# Verbose output
filevault encrypt document.pdf --verbose

# Force overwrite
filevault encrypt document.pdf -f

# Keep original file
filevault encrypt document.pdf --keep

# Custom iterations
filevault encrypt document.pdf --iterations 200000

# Secure delete original
filevault encrypt document.pdf --secure-delete
```

### **📋 Examples với Output**

```bash
$ filevault encrypt financial_report.pdf --verbose
🔐 FileVault v1.0 - Encrypting: financial_report.pdf
Enter password: ********
Confirm password: ********
🔑 Generating salt and deriving key...
🔄 Using AES-256-GCM with PBKDF2 (100,000 iterations)
📊 [████████████████████████████████] 100% | 2.4MB/2.4MB | 15MB/s
✅ Successfully encrypted: financial_report.pdf → financial_report.pdf.enc
⏱️  Time elapsed: 1.2s
```

```bash
$ filevault info report.pdf.enc
📄 File Information
==================
File: report.pdf.enc
Original Name: financial_report.pdf
Format: FileVault v1.0
Algorithm: AES-256-GCM  
Key Derivation: PBKDF2-SHA256 (100,000 iterations)
Original Size: 2.4 MB (2,458,624 bytes)
Encrypted Size: 2.4 MB (2,458,756 bytes)
Overhead: 132 bytes (0.005%)
Created: 2024-09-23 14:30:22 UTC
Status: ✅ Valid & Intact
```

## 🛡️ Bảo Mật

### **🔐 Cryptographic Specifications**

| Component | Algorithm | Parameters |
|-----------|-----------|------------|
| **Encryption** | AES-256-GCM | 256-bit key, 128-bit auth tag |
| **Key Derivation** | PBKDF2-SHA256 | 100,000 iterations, 32-byte salt |
| **Random Generation** | crypto/rand | Cryptographically secure |
| **Authentication** | GCM Mode | Built-in authenticated encryption |

### **🔒 Security Features**

- **🧂 Unique Salt**: Mỗi file có 32-byte random salt riêng biệt
- **🎯 Authenticated Encryption**: GCM mode chống tampering
- **🧠 Memory Security**: Zero-ize sensitive data sau sử dụng
- **⚡ Constant-Time**: Password comparison chống timing attacks
- **🔄 Secure Random**: Sử dụng OS entropy pool
- **📋 Input Validation**: Comprehensive validation cho tất cả inputs

### **🔍 Security Analysis**

FileVault đã được phân tích bảo mật toàn diện. Xem chi tiết tại:
- [Security Analysis Document](docs/Security-analysis.md)
- [Threat Model](docs/Security-analysis.md#threat-model)  
- [Security Testing Results](docs/Security-analysis.md#security-testing)

### **🚨 Security Best Practices**

```bash
# ✅ SỬ DỤNG MẬT KHẨU MẠNH
- Tối thiểu 12 ký tự
- Bao gồm chữ hoa, thường, số, ký tự đặc biệt
- Không sử dụng từ điển hoặc thông tin cá nhân

# ✅ BẢO QUẢN AN TOÀN  
- Lưu file encrypted ở nơi an toàn
- Backup mật khẩu một cách bảo mật
- Không share mật khẩu qua kênh không an toàn

# ✅ VERIFY INTEGRITY
- Chạy 'filevault verify' định kỳ
- Kiểm tra file size và timestamp
- Test decrypt trước khi xóa file gốc
```

## 🏗️ Kiến Trúc

### **📦 Project Structure**

```
filevault/
├── cmd/filevault/           # CLI entry point
├── internal/                # Private packages
│   ├── crypto/              # Encryption algorithms  
│   ├── fileops/             # File I/O operations
│   ├── security/            # Security utilities
│   ├── cli/                 # CLI handlers
│   └── core/                # Business logic
├── pkg/filevault/           # Public API
├── test/                    # Test suites
├── docs/                    # Documentation
└── scripts/                 # Build scripts
```

### **🔧 Technical Stack**

- **Language**: Go 1.25+
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Progress Bar**: [pb](https://github.com/cheggaaa/pb)
- **Cryptography**: Go standard crypto libraries
- **Testing**: Go testing + testify
- **Build**: Make + GitHub Actions

### **📋 File Format**

FileVault sử dụng custom binary format:

```
┌─────────────────────────────────┐
│ Header (108+ bytes)             │
│ ├─ Magic: "FVLT" (4 bytes)      │
│ ├─ Version: uint32 (4 bytes)    │  
│ ├─ Algorithm: uint32 (4 bytes)  │
│ ├─ Salt: [32]byte (32 bytes)    │
│ ├─ IV: [16]byte (16 bytes)      │
│ ├─ Original Size: uint64 (8)    │
│ ├─ Filename Length: uint32 (4)  │
│ ├─ Original Name: variable      │
│ ├─ Reserved: [32]byte (32)      │
│ └─ Checksum: [16]byte (16)      │
├─────────────────────────────────┤
│ Encrypted Data (variable)       │
├─────────────────────────────────┤  
│ Auth Tag: [16]byte (16 bytes)   │
└─────────────────────────────────┘
```

## 🧪 Testing

### **🔬 Test Suite**

FileVault có comprehensive test suite:

```bash
# Chạy tất cả tests
make test

# Unit tests
go test ./internal/...

# Integration tests  
go test ./test/integration/...

# Benchmark tests
go test -bench=. ./test/benchmarks/...

# Coverage report
make test-coverage
```

### **📊 Test Coverage**

- **Unit Tests**: 95%+ coverage
- **Integration Tests**: End-to-end workflows
- **Security Tests**: Cryptographic validation
- **Performance Tests**: Speed và memory benchmarks
- **Cross-Platform Tests**: Windows, Linux, macOS

### **🎯 Test Categories**

```bash
# Crypto tests - Validate encryption/decryption
go test ./internal/crypto/... -v

# File operations - Test I/O operations  
go test ./internal/fileops/... -v

# Security tests - Validate security measures
go test ./internal/security/... -v

# CLI tests - Test command-line interface
go test ./internal/cli/... -v
```

## 🔧 Development

### **🚀 Getting Started**

```bash
# Setup development environment
git clone https://github.com/vuongdat67/NT140.Q11.ANTT-Group15.git
cd NT140.Q11.ANTT-Group15
make setup

# Install dependencies
go mod download

# Run in development mode
go run cmd/filevault/main.go --help
```

### **🔨 Build Commands**

```bash
make build          # Build cho current platform
make build-all      # Build cho tất cả platforms  
make clean          # Clean build artifacts
make lint           # Run code linting
make fmt            # Format code
make test           # Run test suite
make release        # Create release build
```

### **📋 API Usage**

FileVault cung cấp public Go API:

```go
package main

import (
    "log"
    "github.com/vuongdat67/NT140.Q11.ANTT-Group15/pkg/filevault"
)

func main() {
    // Create client
    client := filevault.NewClient(
        filevault.WithVerbose(true),
    )
    
    // Encrypt file
    err := client.EncryptFile("document.pdf", "mypassword")
    if err != nil {
        log.Fatal(err)
    }
    
    // Decrypt file  
    err = client.DecryptFile("document.pdf.enc", "mypassword")
    if err != nil {
        log.Fatal(err)
    }
    
    // Verify file
    result, err := client.VerifyFile("document.pdf.enc")
    if err != nil {
        log.Fatal(err) 
    }
    
    if result.IsValid() {
        log.Println("File is valid!")
    }
}
```

## 🤝 Đóng Góp

### **👥 Team Members**

| Thành viên | Vai trò | Chuyên môn |
|------------|---------|------------|
| **Toàn** | Tech Lead | Architecture, Build System |
| **Đạt** | Security Expert | Cryptography, Security Analysis |
| **Quân** | CLI Developer | User Interface, File Operations |
| **Trung** | QA Engineer | Testing, Validation |

### **🔄 Development Process**

1. **Fork** repository
2. **Create** feature branch (`git checkout -b feature/AmazingFeature`)
3. **Commit** changes (`git commit -m 'Add AmazingFeature'`)
4. **Push** branch (`git push origin feature/AmazingFeature`)  
5. **Open** Pull Request

### **📋 Contribution Guidelines**

- Tuân thủ Go coding standards
- Viết tests cho new features
- Update documentation
- Maintain security standards
- Follow semantic versioning

## 📚 Documentation

### **📖 Additional Documentation**

- [🏗️ Architecture Guide](docs/Architecture.md) - System design và patterns
- [🔐 Security Analysis](docs/Security-analysis.md) - Comprehensive security assessment
- [🚀 API Reference](docs/API.md) - CLI commands và Go API
- [⚙️ Development Guide](docs/Development.md) - Setup và development workflow

### **💡 Examples**

- [📝 Basic Usage](examples/basic_usage.md) - Getting started examples
- [📁 Batch Processing](examples/batch_processing.md) - Bulk operations
- [🛡️ Security Guide](examples/security_guide.md) - Security best practices

## 🐛 Troubleshooting

### **❌ Common Issues**

**Q: "Authentication failed" error**
```bash
A: Kiểm tra mật khẩu và file integrity
   filevault verify file.enc
```

**Q: "File not found" error**  
```bash
A: Kiểm tra đường dẫn file
   ls -la /path/to/file
```

**Q: Slow performance trên file lớn**
```bash  
A: Sử dụng SSD và đủ RAM
   Kiểm tra disk space
```

### **🔧 Debug Mode**

```bash
# Enable verbose output
filevault encrypt file.txt --verbose

# Check system info
filevault version --system

# Verify installation
filevault --help
```

## 📊 Performance

### **⚡ Benchmarks**

| File Size | Encrypt Time | Decrypt Time | Memory Usage |
|-----------|--------------|--------------|--------------|
| 1 MB      | 0.1s         | 0.1s         | 8 MB         |
| 10 MB     | 0.5s         | 0.4s         | 12 MB        |
| 100 MB    | 4.2s         | 3.8s         | 24 MB        |
| 1 GB      | 42s          | 38s          | 64 MB        |

*Benchmark trên: Intel i7-8700K, 16GB RAM, NVMe SSD*

## 📄 License

FileVault được phát hành dưới [MIT License](LICENSE).

```
MIT License

Copyright (c) 2024 NT140.Q11.ANTT Group 15

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🙏 Acknowledgments

- **UIT University** - Cung cấp môi trường học tập
- **NT140.Q11 Course** - Network Security coursework  
- **Go Community** - Excellent cryptography libraries
- **Security Researchers** - Inspiration và best practices

---

**📞 Liên Hệ**

- 📧 Email: nt140.group15@uit.edu.vn
- 🐙 GitHub: [NT140.Q11.ANTT-Group15](https://github.com/vuongdat67/NT140.Q11.ANTT-Group15)
- 📝 Issues: [GitHub Issues](https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/issues)

---

<div align="center">

**🔐 FileVault - Bảo mật file đơn giản, hiệu quả và an toàn 🔐**

Made with ❤️ by NT140.Q11.ANTT Group 15

</div>
