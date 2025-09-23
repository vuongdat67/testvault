package benchmarks

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/crypto"
)

// Benchmark configurations
const (
	TestPassword = "BenchmarkTestPassword123!"
)

// File sizes for benchmarking
var benchmarkFileSizes = []struct {
	name string
	size int64
}{
	{"1KB", 1024},
	{"10KB", 10 * 1024},
	{"100KB", 100 * 1024},
	{"1MB", 1024 * 1024},
	{"10MB", 10 * 1024 * 1024},
	{"100MB", 100 * 1024 * 1024},
}

// BenchmarkEncryptionSpeed benchmarks file encryption speed
func BenchmarkEncryptionSpeed(b *testing.B) {
	for _, size := range benchmarkFileSizes {
		b.Run(size.name, func(b *testing.B) {
			benchmarkEncryptFile(b, size.size)
		})
	}
}

// BenchmarkDecryptionSpeed benchmarks file decryption speed
func BenchmarkDecryptionSpeed(b *testing.B) {
	for _, size := range benchmarkFileSizes {
		b.Run(size.name, func(b *testing.B) {
			benchmarkDecryptFile(b, size.size)
		})
	}
}

// BenchmarkCryptoOperations benchmarks core crypto operations
func BenchmarkCryptoOperations(b *testing.B) {
	dataSizes := []int{1024, 64 * 1024, 1024 * 1024} // 1KB, 64KB, 1MB

	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("AESEncrypt_%dB", size), func(b *testing.B) {
			benchmarkAESEncryption(b, size)
		})

		b.Run(fmt.Sprintf("AESDecrypt_%dB", size), func(b *testing.B) {
			benchmarkAESDecryption(b, size)
		})

		b.Run(fmt.Sprintf("PBKDF2_%dB", size), func(b *testing.B) {
			benchmarkPBKDF2(b, size)
		})
	}
}

// BenchmarkMemoryUsage benchmarks memory usage patterns
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("SmallFiles", func(b *testing.B) {
		benchmarkMemoryUsageForSize(b, 10*1024) // 10KB
	})

	b.Run("MediumFiles", func(b *testing.B) {
		benchmarkMemoryUsageForSize(b, 1024*1024) // 1MB
	})

	b.Run("LargeFiles", func(b *testing.B) {
		benchmarkMemoryUsageForSize(b, 50*1024*1024) // 50MB
	})
}

// Helper functions

func benchmarkEncryptFile(b *testing.B, fileSize int64) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "filevault_benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file with random data
	testFile := filepath.Join(tempDir, "test_input.dat")
	if err := createRandomFile(testFile, fileSize); err != nil {
		b.Fatal(err)
	}

	outputFile := filepath.Join(tempDir, "test_output.enc")

	b.ResetTimer()
	b.SetBytes(fileSize)

	for i := 0; i < b.N; i++ {
		// Remove output file if exists
		os.Remove(outputFile)

		// Measure encryption time
		start := time.Now()
		err := core.EncryptFile(testFile, outputFile, TestPassword)
		if err != nil {
			b.Fatal(err)
		}
		elapsed := time.Since(start)

		// Report custom metrics
		mbPerSec := float64(fileSize) / (1024 * 1024) / elapsed.Seconds()
		b.ReportMetric(mbPerSec, "MB/sec")
	}
}

func benchmarkDecryptFile(b *testing.B, fileSize int64) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "filevault_benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file and encrypt it
	testFile := filepath.Join(tempDir, "test_input.dat")
	encryptedFile := filepath.Join(tempDir, "test_encrypted.enc")
	decryptedFile := filepath.Join(tempDir, "test_decrypted.dat")

	if err := createRandomFile(testFile, fileSize); err != nil {
		b.Fatal(err)
	}

	if err := core.EncryptFile(testFile, encryptedFile, TestPassword); err != nil {
		b.Fatal(err)
	}

	// Get encrypted file size for accurate measurement
	encryptedInfo, err := os.Stat(encryptedFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.SetBytes(encryptedInfo.Size())

	for i := 0; i < b.N; i++ {
		// Remove output file if exists
		os.Remove(decryptedFile)

		// Measure decryption time
		start := time.Now()
		err := core.DecryptFile(encryptedFile, decryptedFile, TestPassword)
		if err != nil {
			b.Fatal(err)
		}
		elapsed := time.Since(start)

		// Report custom metrics
		mbPerSec := float64(fileSize) / (1024 * 1024) / elapsed.Seconds()
		b.ReportMetric(mbPerSec, "MB/sec")
	}
}

func benchmarkAESEncryption(b *testing.B, dataSize int) {
	// Generate test data
	data := make([]byte, dataSize)
	rand.Read(data)

	// Generate key from password
	salt, _ := crypto.GenerateSalt32()
	cipher, _ := crypto.NewAESCipherFromPassword(TestPassword, salt)

	b.ResetTimer()
	b.SetBytes(int64(dataSize))

	for i := 0; i < b.N; i++ {
		_, err := cipher.Encrypt(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkAESDecryption(b *testing.B, dataSize int) {
	// Generate test data and encrypt it
	data := make([]byte, dataSize)
	rand.Read(data)

	salt, _ := crypto.GenerateSalt32()
	cipher, _ := crypto.NewAESCipherFromPassword(TestPassword, salt)
	encryptedData, _ := cipher.Encrypt(data)

	b.ResetTimer()
	b.SetBytes(int64(dataSize))

	for i := 0; i < b.N; i++ {
		_, err := cipher.Decrypt(encryptedData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkPBKDF2(b *testing.B, saltSize int) {
	salt := make([]byte, saltSize)
	rand.Read(salt)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = crypto.DeriveKey(TestPassword, salt, 100000)
	}
}

func benchmarkMemoryUsageForSize(b *testing.B, fileSize int64) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "filevault_memory_benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "memory_test.dat")
	if err := createRandomFile(testFile, fileSize); err != nil {
		b.Fatal(err)
	}

	outputFile := filepath.Join(tempDir, "memory_test.enc")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Remove output file
		os.Remove(outputFile)

		// Track memory before operation
		var m1 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		// Perform encryption
		err := core.EncryptFile(testFile, outputFile, TestPassword)
		if err != nil {
			b.Fatal(err)
		}

		// Track memory after operation
		var m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m2)

		// Report memory metrics
		memUsed := float64(m2.Alloc-m1.Alloc) / (1024 * 1024) // MB
		b.ReportMetric(memUsed, "MB_memory")
	}
}

// createRandomFile creates a file with random data of specified size
func createRandomFile(filename string, size int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write random data in chunks to avoid memory issues
	const chunkSize = 64 * 1024 // 64KB chunks
	buffer := make([]byte, chunkSize)

	written := int64(0)
	for written < size {
		remaining := size - written
		currentChunkSize := chunkSize
		if remaining < chunkSize {
			currentChunkSize = int(remaining)
			buffer = buffer[:currentChunkSize]
		}

		if _, err := rand.Read(buffer); err != nil {
			return err
		}

		if _, err := file.Write(buffer); err != nil {
			return err
		}

		written += int64(currentChunkSize)
	}

	return nil
}

// BenchmarkSuite runs comprehensive benchmark suite
func BenchmarkSuite(b *testing.B) {
	// This function can be called to run comprehensive benchmarks
	// and report performance characteristics
	b.Run("EncryptionSpeed", BenchmarkEncryptionSpeed)
	b.Run("DecryptionSpeed", BenchmarkDecryptionSpeed)
	b.Run("CryptoOperations", BenchmarkCryptoOperations)
	b.Run("MemoryUsage", BenchmarkMemoryUsage)
}

// Throughput measurement helpers
func MeasureThroughput(dataSize int64, duration time.Duration) float64 {
	return float64(dataSize) / (1024 * 1024) / duration.Seconds() // MB/s
}

// Memory efficiency measurement
func MeasureMemoryEfficiency(fileSize, memoryUsed int64) float64 {
	return float64(fileSize) / float64(memoryUsed) // bytes per byte of memory
}

// CPU efficiency measurement
func MeasureCPUEfficiency(dataSize int64, cpuTime time.Duration) float64 {
	return float64(dataSize) / (1024 * 1024) / cpuTime.Seconds() // MB/s CPU time
}
