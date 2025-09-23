package fileops

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "sync"
)

// Default buffer sizes for different operations
const (
    DefaultBufferSize = 64 * 1024  // 64KB default
    SmallFileBuffer   = 8 * 1024   // 8KB for small files
    LargeFileBuffer   = 256 * 1024 // 256KB for large files
)

// StreamReader handles streaming file reading for encryption/decryption
type StreamReader struct {
    file      *os.File
    reader    *bufio.Reader
    fileSize  int64
    bytesRead int64
    bufferSize int
    mutex     sync.RWMutex
}

// NewStreamReader creates a new stream reader for the given file
func NewStreamReader(filepath string) (*StreamReader, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
    }
    
    // Get file size
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to get file stats: %w", err)
    }
    
    // Choose buffer size based on file size
    bufferSize := DefaultBufferSize
    if stat.Size() < 1024*1024 { // < 1MB
        bufferSize = SmallFileBuffer
    } else if stat.Size() > 100*1024*1024 { // > 100MB
        bufferSize = LargeFileBuffer
    }
    
    return &StreamReader{
        file:       file,
        reader:     bufio.NewReaderSize(file, bufferSize),
        fileSize:   stat.Size(),
        bytesRead:  0,
        bufferSize: bufferSize,
    }, nil
}

// NewStreamReaderWithBuffer creates a reader with custom buffer size
func NewStreamReaderWithBuffer(filepath string, bufferSize int) (*StreamReader, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
    }
    
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to get file stats: %w", err)
    }
    
    return &StreamReader{
        file:       file,
        reader:     bufio.NewReaderSize(file, bufferSize),
        fileSize:   stat.Size(),
        bytesRead:  0,
        bufferSize: bufferSize,
    }, nil
}

// Read reads data into the provided buffer
func (sr *StreamReader) Read(buffer []byte) (int, error) {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()
    
    n, err := sr.reader.Read(buffer)
    sr.bytesRead += int64(n)
    return n, err
}

// ReadChunk reads a specific chunk size
func (sr *StreamReader) ReadChunk(chunkSize int) ([]byte, error) {
    buffer := make([]byte, chunkSize)
    n, err := sr.Read(buffer)
    if err != nil && err != io.EOF {
        return nil, err
    }
    
    // Return only the bytes actually read
    return buffer[:n], err
}

// ReadChunkOptimized reads a chunk with optimized buffer reuse
func (sr *StreamReader) ReadChunkOptimized(buffer []byte) (int, error) {
    return sr.Read(buffer)
}

// GetProgress returns reading progress (0.0 to 1.0)
func (sr *StreamReader) GetProgress() float64 {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()
    
    if sr.fileSize == 0 {
        return 1.0
    }
    return float64(sr.bytesRead) / float64(sr.fileSize)
}

// GetBytesRead returns total bytes read so far
func (sr *StreamReader) GetBytesRead() int64 {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()
    return sr.bytesRead
}

// GetFileSize returns total file size
func (sr *StreamReader) GetFileSize() int64 {
    return sr.fileSize
}

// GetBufferSize returns current buffer size
func (sr *StreamReader) GetBufferSize() int {
    return sr.bufferSize
}

// IsEOF checks if we've reached end of file
func (sr *StreamReader) IsEOF() bool {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()
    return sr.bytesRead >= sr.fileSize
}

// Reset resets the reader to the beginning of the file
func (sr *StreamReader) Reset() error {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()
    
    _, err := sr.file.Seek(0, io.SeekStart)
    if err != nil {
        return fmt.Errorf("failed to reset file position: %w", err)
    }
    
    sr.reader = bufio.NewReaderSize(sr.file, sr.bufferSize)
    sr.bytesRead = 0
    return nil
}

// Seek moves to a specific position in the file
func (sr *StreamReader) Seek(offset int64, whence int) (int64, error) {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()
    
    pos, err := sr.file.Seek(offset, whence)
    if err != nil {
        return pos, fmt.Errorf("failed to seek to position: %w", err)
    }
    
    // Reset buffer after seeking
    sr.reader = bufio.NewReaderSize(sr.file, sr.bufferSize)
    sr.bytesRead = pos
    
    return pos, nil
}

// ReadAt reads data at a specific offset without changing file position
func (sr *StreamReader) ReadAt(buffer []byte, offset int64) (int, error) {
    return sr.file.ReadAt(buffer, offset)
}

// Prefetch attempts to prefetch data for better performance
func (sr *StreamReader) Prefetch() error {
    // This is a hint to the OS to read ahead
    // Implementation would be platform-specific
    return nil
}

// GetReadSpeed calculates current read speed in bytes per second
func (sr *StreamReader) GetReadSpeed(elapsedTime float64) float64 {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()
    
    if elapsedTime <= 0 {
        return 0
    }
    return float64(sr.bytesRead) / elapsedTime
}

// Close closes the underlying file
func (sr *StreamReader) Close() error {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()
    
    if sr.file != nil {
        return sr.file.Close()
    }
    return nil
}

// StreamingReadCallback is called during streaming operations
type StreamingReadCallback func(bytesRead, totalBytes int64, speed float64)

// ReadWithCallback reads with progress callback
func (sr *StreamReader) ReadWithCallback(buffer []byte, callback StreamingReadCallback, startTime float64) (int, error) {
    n, err := sr.Read(buffer)
    
    if callback != nil {
        elapsed := startTime // This should be calculated by caller
        speed := sr.GetReadSpeed(elapsed)
        callback(sr.GetBytesRead(), sr.GetFileSize(), speed)
    }
    
    return n, err
}