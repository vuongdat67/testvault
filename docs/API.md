# FileVault API Documentation

**Version:** 1.0  
**Last Updated:** November 2024

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Basic Usage](#basic-usage)
4. [Command Reference](#command-reference)
5. [Configuration](#configuration)
6. [Examples](#examples)
7. [Error Handling](#error-handling)
8. [Security Best Practices](#security-best-practices)

---

## Overview

FileVault is a command-line tool for secure file encryption using AES-256-GCM with PBKDF2 key derivation. It provides authenticated encryption with integrity verification for protecting sensitive files.

### Key Features
- **AES-256-GCM** authenticated encryption
- **PBKDF2-SHA256** key derivation (100,000+ iterations)
- **Unique salt** per file (32 bytes)
- **Progress tracking** for large files
- **Batch processing** support
- **Cross-platform** compatibility (Windows, Linux, macOS)

### System Requirements
- **Operating System**: Windows 10+, Linux, macOS 10.15+
- **Architecture**: x86_64, ARM64
- **Memory**: 256MB RAM minimum
- **Disk Space**: 50MB available space

---

## Installation

### Binary Installation

#### Windows
```powershell
# Download from GitHub releases
Invoke-WebRequest -Uri "https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/releases/download/v1.0.0/filevault-windows-amd64.exe" -OutFile "filevault.exe"

# Add to PATH (optional)
$env:PATH += ";C:\path\to\filevault"
```

#### Linux
```bash
# Download and install
wget https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/releases/download/v1.0.0/filevault-linux-amd64.tar.gz
tar -xzf filevault-linux-amd64.tar.gz
sudo mv filevault /usr/local/bin/
```

#### macOS
```bash
# Download and install
curl -L https://github.com/vuongdat67/NT140.Q11.ANTT-Group15/releases/download/v1.0.0/filevault-darwin-amd64.tar.gz -o filevault-darwin-amd64.tar.gz
tar -xzf filevault-darwin-amd64.tar.gz
sudo mv filevault /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/vuongdat67/NT140.Q11.ANTT-Group15.git
cd NT140.Q11.ANTT-Group15
make build
```

---

## Basic Usage

### Quick Start

```bash
# Encrypt a file
filevault encrypt document.pdf

# Decrypt a file
filevault decrypt document.pdf.enc

# View file information
filevault info document.pdf.enc

# Verify file integrity
filevault verify document.pdf.enc
```

### Global Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--quiet` | `-q` | Suppress non-error output | `false` |
| `--help` | `-h` | Show help message | - |
| `--version` | - | Show version information | - |

---

## Command Reference

### `filevault encrypt`

Encrypt files using AES-256-GCM with PBKDF2 key derivation.

#### Syntax
```bash
filevault encrypt [file...] [flags]
```

#### Arguments
- `file...` - One or more files to encrypt

#### Flags
| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--output` | `-o` | string | Output file or directory | `{input}.enc` |
| `--force` | `-f` | bool | Overwrite existing files | `false` |
| `--keep` | `-k` | bool | Keep original files | `false` |
| `--iterations` | - | int | PBKDF2 iterations | `100000` |

#### Examples
```bash
# Basic encryption
filevault encrypt document.pdf

# Specify output file
filevault encrypt document.pdf -o secure.enc

# Encrypt multiple files
filevault encrypt *.txt *.pdf

# Encrypt to directory
filevault encrypt file1.txt file2.pdf -o encrypted/

# Keep original files
filevault encrypt important.doc --keep

# High security (more iterations)
filevault encrypt secret.txt --iterations 200000

# Force overwrite
filevault encrypt data.xlsx -o backup.enc --force
```

#### Output Format
```
✅ Encrypted: document.pdf → document.pdf.enc
   File size: 2.4 MB → 2.4 MB
   Encryption completed in 1.2s
```

---

### `filevault decrypt`

Decrypt FileVault encrypted files.

#### Syntax
```bash
filevault decrypt [file...] [flags]
```

#### Arguments
- `file...` - One or more encrypted files to decrypt

#### Flags
| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--output` | `-o` | string | Output file or directory | Auto-detected |
| `--force` | `-f` | bool | Overwrite existing files | `false` |

#### Examples
```bash
# Basic decryption
filevault decrypt document.pdf.enc

# Specify output file
filevault decrypt encrypted.enc -o recovered.pdf

# Decrypt multiple files
filevault decrypt *.enc

# Decrypt to directory
filevault decrypt file1.enc file2.enc -o decrypted/

# Force overwrite
filevault decrypt backup.enc -o original.txt --force
```

#### Output Format
```
✅ Decrypted: document.pdf.enc → document.pdf
   Original size: 2.4 MB
   Decryption completed in 0.8s
```

---

### `filevault info`

Display comprehensive information about encrypted files.

#### Syntax
```bash
filevault info [file...] [flags]
```

#### Arguments
- `file...` - One or more encrypted files to analyze

#### Flags
| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--hex` | - | bool | Show crypto parameters in hex | `false` |

#### Examples
```bash
# Basic file information
filevault info document.pdf.enc

# Show cryptographic parameters
filevault info --hex secret.enc

# Analyze multiple files
filevault info backup1.enc backup2.enc backup3.enc

# Verbose batch analysis
filevault info -v *.enc
```

#### Output Format
```
═══════════════════════════════════════════════════════════════
                     FILE ANALYSIS REPORT                      
═══════════════════════════════════════════════════════════════

Basic Information:
  File: document.pdf.enc
  File Size: 2.4 MB
  Modified: 2024-11-15T14:30:22Z
  Permissions: -rw-r--r--

Format Information:
  Status: ✅ Valid FileVault File
  Format: FileVault v1
  Algorithm: AES-256-GCM
  Key Derivation: PBKDF2-SHA256 (100,000 iterations)

Original File Information:
  Original Filename: document.pdf
  Original Size: 2.4 MB
  Compression Ratio: 100.1%
```

---

### `filevault verify`

Verify file integrity and format validity.

#### Syntax
```bash
filevault verify [file...] [flags]
```

#### Arguments
- `file...` - One or more files to verify

#### Flags
| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--deep` | - | bool | Deep integrity check (requires password) | `false` |

#### Examples
```bash
# Basic verification
filevault verify document.pdf.enc

# Verify multiple files
filevault verify *.enc

# Batch verification with summary
filevault verify backups/*.enc

# Deep verification (with password)
filevault verify --deep sensitive.enc
```

#### Output Format
```bash
# Single file
✅ File verification successful
   Format: FileVault v1
   Algorithm: AES-256-GCM
   Original file: document.pdf (2.4 MB)

# Batch verification
Verification Summary:
====================
Total files: 10
✅ Valid: 9
❌ Invalid: 1
```

---

### `filevault version`

Display version and build information.

#### Syntax
```bash
filevault version
```

#### Output Format
```
FileVault 1.0.0
Commit: abc123def
Built: 2024-11-15T10:30:00Z
Go version: go1.21.0
```

---

### `filevault help`

Get help for any command.

#### Syntax
```bash
filevault help [command]
```

#### Examples
```bash
# General help
filevault help

# Command-specific help
filevault help encrypt
filevault help decrypt
```

---

## Configuration

### Configuration File

FileVault supports a JSON configuration file for default settings:

**Location**: `~/.filevault/config.json`

#### Example Configuration
```json
{
  "default_iterations": 150000,
  "default_algorithm": "AES-256-GCM",
  "buffer_size": 65536,
  "password_min_length": 12,
  "require_strong_password": true,
  "secure_memory": true,
  "use_colors": true,
  "show_progress": true,
  "verbose_output": false,
  "max_file_size": 10737418240,
  "streaming_threshold": 1048576
}
```

#### Configuration Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `default_iterations` | int | PBKDF2 iterations | `100000` |
| `default_algorithm` | string | Encryption algorithm | `"AES-256-GCM"` |
| `buffer_size` | int | I/O buffer size (bytes) | `65536` |
| `password_min_length` | int | Minimum password length | `8` |
| `require_strong_password` | bool | Enforce strong passwords | `false` |
| `secure_memory` | bool | Enable secure memory handling | `true` |
| `use_colors` | bool | Enable colored output | `true` |
| `show_progress` | bool | Show progress bars | `true` |
| `verbose_output` | bool | Default verbose mode | `false` |
| `max_file_size` | int | Maximum file size (bytes) | `10737418240` |
| `streaming_threshold` | int | Streaming threshold (bytes) | `1048576` |

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FILEVAULT_CONFIG_DIR` | Configuration directory | `~/.filevault` |
| `FILEVAULT_TEMP_DIR` | Temporary file directory | System temp |
| `FILEVAULT_NO_COLOR` | Disable colored output | `false` |

---

## Examples

### Basic File Encryption

```bash
# Encrypt a single document
filevault encrypt contract.pdf
# Output: contract.pdf.enc

# Provide password when prompted:
# Enter password for encryption: ********
# Confirm password: ********
# ✅ Encrypted: contract.pdf → contract.pdf.enc
```

### Batch Processing

```bash
# Encrypt all text files in current directory
filevault encrypt *.txt

# Encrypt multiple files to specific directory
filevault encrypt file1.pdf file2.docx file3.xlsx -o encrypted/

# Decrypt all .enc files in a directory
filevault decrypt encrypted/*.enc -o decrypted/
```

### Advanced Usage

```bash
# High-security encryption (200k iterations)
filevault encrypt top-secret.pdf --iterations 200000 --keep

# Encrypt with custom output name
filevault encrypt data.csv -o backup-2024.enc

# Force overwrite existing encrypted file
filevault encrypt updated-data.xlsx -o data.xlsx.enc --force

# Verify integrity of backup files
filevault verify backups/*.enc
```

### Automation and Scripting

#### Windows PowerShell
```powershell
# Batch encrypt all PDFs
Get-ChildItem "*.pdf" | ForEach-Object {
    filevault encrypt $_.Name
}

# Verify all encrypted files
Get-ChildItem "*.enc" | ForEach-Object {
    filevault verify $_.Name
}
```

#### Bash Script
```bash
#!/bin/bash
# Automated backup encryption

BACKUP_DIR="/home/user/backups"
ENCRYPTED_DIR="/home/user/encrypted-backups"

# Create encrypted backup directory
mkdir -p "$ENCRYPTED_DIR"

# Encrypt all files in backup directory
for file in "$BACKUP_DIR"/*; do
    if [[ -f "$file" ]]; then
        echo "Encrypting: $(basename "$file")"
        filevault encrypt "$file" -o "$ENCRYPTED_DIR/$(basename "$file").enc"
    fi
done

# Verify all encrypted files
filevault verify "$ENCRYPTED_DIR"/*.enc
```

### Integration Examples

#### Python Integration
```python
import subprocess
import sys

def encrypt_file(file_path, password=None):
    """Encrypt a file using FileVault."""
    cmd = ['filevault', 'encrypt', file_path]
    
    try:
        if password:
            # For automated scripts (security risk!)
            proc = subprocess.Popen(cmd, stdin=subprocess.PIPE, 
                                  stdout=subprocess.PIPE, 
                                  stderr=subprocess.PIPE, text=True)
            stdout, stderr = proc.communicate(input=f"{password}\n{password}\n")
        else:
            # Interactive mode
            result = subprocess.run(cmd, check=True)
            
        return True
    except subprocess.CalledProcessError as e:
        print(f"Encryption failed: {e}")
        return False

# Usage
if encrypt_file("document.pdf"):
    print("File encrypted successfully")
```

#### Node.js Integration
```javascript
const { exec } = require('child_process');
const path = require('path');

function encryptFile(filePath) {
    return new Promise((resolve, reject) => {
        const cmd = `filevault encrypt "${filePath}"`;
        
        exec(cmd, (error, stdout, stderr) => {
            if (error) {
                reject(new Error(`Encryption failed: ${error.message}`));
                return;
            }
            
            if (stderr) {
                console.warn('Warning:', stderr);
            }
            
            console.log('Success:', stdout);
            resolve(stdout);
        });
    });
}

// Usage
encryptFile('./document.pdf')
    .then(result => console.log('Encrypted:', result))
    .catch(error => console.error('Error:', error));
```

---

## Error Handling

### Exit Codes

| Code | Description | Common Causes |
|------|-------------|---------------|
| `0` | Success | Operation completed successfully |
| `1` | General error | Unspecified error occurred |
| `2` | File not found | Input file doesn't exist |
| `3` | Permission denied | Insufficient file permissions |
| `4` | Authentication failed | Wrong password or corrupted file |
| `5` | Corrupted file | Invalid file format or data corruption |
| `6` | Insufficient resources | Out of memory or disk space |
| `7` | Invalid arguments | Malformed command line arguments |

### Common Error Messages

#### File Operations
```bash
❌ Error: File not found: document.pdf
Suggestions:
  • Verify the file path is correct
  • Check if the file exists in the current directory
  • Use absolute path if relative path doesn't work

❌ Error: Permission denied: /protected/file.txt
Suggestions:
  • Check file permissions
  • Run with appropriate user privileges
  • Ensure directory is writable
```

#### Encryption/Decryption
```bash
❌ Error: Authentication failed - wrong password or corrupted file
Suggestions:
  • Make sure you're using the correct password
  • Check for typos in the password
  • Verify the file hasn't been corrupted

❌ Error: File appears to be corrupted or damaged
Suggestions:
  • Verify file integrity with: filevault verify file.enc
  • Check if the file was completely downloaded/copied
  • Restore from backup if available
```

#### Password Issues
```bash
❌ Error: Password is too weak
Suggestions:
  • Use at least 12 characters
  • Include uppercase and lowercase letters
  • Add numbers and special characters
  • Use --force to proceed with weak password (not recommended)
```

### Debugging

#### Verbose Mode
```bash
# Enable detailed logging
filevault encrypt document.pdf --verbose

# Output includes:
# - File size and processing information
# - Cryptographic parameters used
# - Performance metrics
# - Step-by-step operation details
```

#### Verification and Testing
```bash
# Test file integrity
filevault verify suspicious-file.enc

# Check if file is encrypted
filevault info unknown-file.dat

# Verify multiple files
filevault verify *.enc --verbose
```

---

## Security Best Practices

### Password Security

#### Strong Password Guidelines
```bash
✅ Minimum 12 characters (recommended 16+)
✅ Mix of uppercase, lowercase, numbers, symbols
✅ Avoid dictionary words and personal information
✅ Use unique passwords for different purposes
✅ Consider using a password manager
```

#### Examples
```bash
# Weak passwords (avoid)
❌ "password123"
❌ "JohnSmith1990"
❌ "qwerty"

# Strong passwords
✅ "Tr0ub4dor&3"
✅ "correct horse battery staple"
✅ "My$ecur3P@ssw0rd2024!"
```

### File Security

#### Secure File Handling
```bash
# Always verify files after encryption
filevault encrypt important.pdf
filevault verify important.pdf.enc

# Keep backups of encrypted files
cp important.pdf.enc backup/

# Securely delete original files (Linux/macOS)
shred -vfz -n 3 important.pdf  # Linux
rm -P important.pdf            # macOS

# Or use FileVault's automatic cleanup
filevault encrypt important.pdf  # Original deleted by default
```

#### Storage Recommendations
```bash
✅ Store encrypted files in multiple locations
✅ Use cloud storage for encrypted files (safe)
✅ Regular backup verification
✅ Document password recovery procedures
❌ Don't store passwords with encrypted files
❌ Don't rely on single backup location
```

### Operational Security

#### Environment Security
```bash
# Check for memory dumps/core dumps
ulimit -c 0  # Disable core dumps

# Clear command history after use
history -c
export HISTSIZE=0

# Use secure temporary directories
export TMPDIR=/secure/temp
```

#### Process Security
```bash
# Monitor for FileVault processes
ps aux | grep filevault

# Check memory usage during operations
top -p $(pgrep filevault)

# Verify no sensitive data in logs
journalctl -u filevault  # Linux
Console.app             # macOS
```

### Network Security

#### File Transfer Security
```bash
# Safe: Transfer encrypted files over insecure channels
scp document.pdf.enc user@remote:/backup/

# Unsafe: Transfer unencrypted files
❌ scp document.pdf user@remote:/backup/

# Email encrypted files (safe)
# Email unencrypted files (unsafe)
```

### Compliance Considerations

#### Data Protection Regulations
- **GDPR**: Encryption satisfies "appropriate technical measures"
- **HIPAA**: AES-256 encryption meets security requirements
- **PCI DSS**: Strong cryptography for cardholder data protection
- **SOX**: File encryption supports data integrity requirements

#### Audit Requirements
```bash
# Document encryption procedures
filevault info *.enc > encryption-audit.txt

# Verify file integrity for audits
filevault verify audit-files/*.enc

# Maintain encryption logs
filevault encrypt --verbose file.pdf 2>&1 | tee encryption.log
```

---

This API documentation provides comprehensive guidance for using FileVault securely and effectively. For additional support or questions, please refer to the project's GitHub repository or security documentation.
