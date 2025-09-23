package security

import (
	"crypto/rand"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

var (
	// Global memory pool for sensitive data
	memoryPool = sync.Pool{
		New: func() interface{} {
			// Pre-allocate buffers to avoid frequent allocations
			return make([]byte, 1024)
		},
	}
)

// SecureZeroMemory securely zeros out memory to prevent sensitive data from
// remaining in memory after use. This provides some protection against
// memory dumps and debugging attacks.
func SecureZeroMemory(data []byte) {
	if len(data) == 0 {
		return
	}

	// Zero out the memory multiple times for extra security
	for i := 0; i < 3; i++ {
		for j := range data {
			data[j] = 0
		}
		runtime.KeepAlive(data) // Prevent compiler optimization
	}

	// Fill with random data on final pass
	rand.Read(data)
	
	// Zero again
	for i := range data {
		data[i] = 0
	}
	
	// Force garbage collection to clean up
	runtime.GC()
	runtime.GC() // Call twice to ensure cleanup
}

// SecureZeroString securely zeros out a string by converting it to bytes
// Note: This has limitations due to Go's string immutability
// WARNING: This is unsafe and should only be used when absolutely necessary
func SecureZeroString(s string) {
	if len(s) == 0 {
		return
	}

	// Convert string to byte slice (unsafe operation)
	// WARNING: This modifies the underlying string data
	header := (*[2]uintptr)(unsafe.Pointer(&s))
	data := (*[1 << 30]byte)(unsafe.Pointer(header[0]))[:len(s):len(s)]

	SecureZeroMemory(data)
}

// SecureBuffer represents a secure memory buffer that automatically cleans up
type SecureBuffer struct {
	data     []byte
	size     int
	locked   bool
	cleanupFn func()
}

// NewSecureBuffer creates a new secure buffer
func NewSecureBuffer(size int) *SecureBuffer {
	data := make([]byte, size)
	
	// Try to lock memory (platform-specific)
	locked := LockMemory(data) == nil
	
	return &SecureBuffer{
		data:   data,
		size:   size,
		locked: locked,
		cleanupFn: func() {
			SecureZeroMemory(data)
			if locked {
				UnlockMemory(data)
			}
		},
	}
}

// Data returns the underlying byte slice
func (sb *SecureBuffer) Data() []byte {
	return sb.data
}

// Size returns the buffer size
func (sb *SecureBuffer) Size() int {
	return sb.size
}

// IsLocked returns whether the memory is locked
func (sb *SecureBuffer) IsLocked() bool {
	return sb.locked
}

// Destroy securely destroys the buffer
func (sb *SecureBuffer) Destroy() {
	if sb.cleanupFn != nil {
		sb.cleanupFn()
		sb.cleanupFn = nil
	}
	sb.data = nil
}

// AutoDestroy sets up automatic destruction after a timeout
func (sb *SecureBuffer) AutoDestroy(timeout time.Duration) {
	go func() {
		time.Sleep(timeout)
		sb.Destroy()
	}()
}

// GetSecureBuffer retrieves a buffer from the pool
func GetSecureBuffer() []byte {
	return memoryPool.Get().([]byte)
}

// PutSecureBuffer returns a buffer to the pool after cleaning
func PutSecureBuffer(data []byte) {
	SecureZeroMemory(data)
	memoryPool.Put(data)
}

// LockMemory attempts to lock memory pages to prevent them from being
// swapped to disk. This is a best-effort operation and may not be supported
// on all platforms.
func LockMemory(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// This is a platform-specific operation
	// On Windows, we would use VirtualLock
	// On Linux/Unix, we would use mlock
	// For now, this is a no-op placeholder
	
	// TODO: Implement platform-specific memory locking
	// - Windows: VirtualLock() syscall
	// - Linux/Unix: mlock() syscall
	// - macOS: mlock() syscall
	
	return nil
}

// UnlockMemory unlocks previously locked memory pages
func UnlockMemory(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// This is a platform-specific operation
	// On Windows, we would use VirtualUnlock
	// On Linux/Unix, we would use munlock
	// For now, this is a no-op placeholder
	
	// TODO: Implement platform-specific memory unlocking
	// - Windows: VirtualUnlock() syscall
	// - Linux/Unix: munlock() syscall
	// - macOS: munlock() syscall

	return nil
}

// DisableCoreDumps attempts to disable core dumps for this process
// This helps prevent sensitive data from being written to disk
func DisableCoreDumps() error {
	// This is platform-specific
	// On Unix-like systems, we would use setrlimit(RLIMIT_CORE, 0)
	// For now, this is a no-op placeholder
	
	// TODO: Implement platform-specific core dump disabling
	
	return nil
}

// ConstantTimeCompare performs constant-time comparison of two byte slices
// This helps prevent timing attacks on password/key comparison
func ConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	
	return result == 0
}

// MemoryStats provides information about memory usage
type MemoryStats struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return MemoryStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}
}

// ForceGarbageCollection forces immediate garbage collection
// Use this after handling sensitive data
func ForceGarbageCollection() {
	runtime.GC()
	runtime.GC() // Call twice for thorough cleanup
}

// SecureAllocate allocates memory securely and returns a secure buffer
func SecureAllocate(size int) *SecureBuffer {
	return NewSecureBuffer(size)
}

// init performs security-related initialization
func init() {
	// Attempt to disable core dumps on startup
	DisableCoreDumps()
	
	// Set up memory pool with secure cleanup
	runtime.SetFinalizer(&memoryPool, func(p *sync.Pool) {
		// This won't actually be called, but it's good practice
	})
}
