package fileops

import (
    "bufio"
    "fmt"
    "os"
    "sync"
    "time"
)

// StreamWriter handles streaming file writing for encryption/decryption
type StreamWriter struct {
    file         *os.File
    writer       *bufio.Writer
    bytesWritten int64
    filepath     string
    bufferSize   int
    mutex        sync.RWMutex
    startTime    time.Time
}

// NewStreamWriter creates a new stream writer for the given file
func NewStreamWriter(filepath string) (*StreamWriter, error) {
    file, err := os.Create(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to create file %s: %w", filepath, err)
    }
    
    // Choose buffer size for optimal performance
    bufferSize := DefaultBufferSize
    
    return &StreamWriter{
        file:         file,
        writer:       bufio.NewWriterSize(file, bufferSize),
        bytesWritten: 0,
        filepath:     filepath,
        bufferSize:   bufferSize,
        startTime:    time.Now(),
    }, nil
}

// NewStreamWriterWithBuffer creates a writer with custom buffer size
func NewStreamWriterWithBuffer(filepath string, bufferSize int) (*StreamWriter, error) {
    file, err := os.Create(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to create file %s: %w", filepath, err)
    }
    
    return &StreamWriter{
        file:         file,
        writer:       bufio.NewWriterSize(file, bufferSize),
        bytesWritten: 0,
        filepath:     filepath,
        bufferSize:   bufferSize,
        startTime:    time.Now(),
    }, nil
}

// NewStreamWriterAppend creates a stream writer in append mode
func NewStreamWriterAppend(filepath string) (*StreamWriter, error) {
    file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open file %s for append: %w", filepath, err)
    }
    
    // Get current file size
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to get file stats: %w", err)
    }
    
    return &StreamWriter{
        file:         file,
        writer:       bufio.NewWriterSize(file, DefaultBufferSize),
        bytesWritten: stat.Size(),
        filepath:     filepath,
        bufferSize:   DefaultBufferSize,
        startTime:    time.Now(),
    }, nil
}

// Write writes data from the provided buffer
func (sw *StreamWriter) Write(data []byte) (int, error) {
    sw.mutex.Lock()
    defer sw.mutex.Unlock()
    
    n, err := sw.writer.Write(data)
    sw.bytesWritten += int64(n)
    return n, err
}

// WriteChunk writes a chunk of data
func (sw *StreamWriter) WriteChunk(chunk []byte) error {
    _, err := sw.Write(chunk)
    return err
}

// WriteBatch writes multiple chunks efficiently
func (sw *StreamWriter) WriteBatch(chunks [][]byte) error {
    sw.mutex.Lock()
    defer sw.mutex.Unlock()
    
    for _, chunk := range chunks {
        n, err := sw.writer.Write(chunk)
        if err != nil {
            return err
        }
        sw.bytesWritten += int64(n)
    }
    return nil
}

// GetBytesWritten returns total bytes written so far
func (sw *StreamWriter) GetBytesWritten() int64 {
    sw.mutex.RLock()
    defer sw.mutex.RUnlock()
    return sw.bytesWritten
}

// GetFilepath returns the file path
func (sw *StreamWriter) GetFilepath() string {
    return sw.filepath
}

// GetBufferSize returns current buffer size
func (sw *StreamWriter) GetBufferSize() int {
    return sw.bufferSize
}

// GetWriteSpeed calculates current write speed in bytes per second
func (sw *StreamWriter) GetWriteSpeed() float64 {
    sw.mutex.RLock()
    defer sw.mutex.RUnlock()
    
    elapsed := time.Since(sw.startTime).Seconds()
    if elapsed <= 0 {
        return 0
    }
    return float64(sw.bytesWritten) / elapsed
}

// Flush flushes the buffered data to disk
func (sw *StreamWriter) Flush() error {
    sw.mutex.Lock()
    defer sw.mutex.Unlock()
    
    return sw.writer.Flush()
}

// Sync syncs the file to disk (calls fsync)
func (sw *StreamWriter) Sync() error {
    if err := sw.Flush(); err != nil {
        return err
    }
    return sw.file.Sync()
}

// Close flushes and closes the underlying file
func (sw *StreamWriter) Close() error {
    sw.mutex.Lock()
    defer sw.mutex.Unlock()
    
    if sw.writer != nil {
        if err := sw.writer.Flush(); err != nil {
            return fmt.Errorf("failed to flush buffer: %w", err)
        }
    }
    
    if sw.file != nil {
        return sw.file.Close()
    }
    return nil
}

// CleanupOnError removes the file if there was an error during writing
func (sw *StreamWriter) CleanupOnError() {
    sw.Close()
    os.Remove(sw.filepath)
}

// WriteAt writes data at a specific offset
func (sw *StreamWriter) WriteAt(data []byte, offset int64) (int, error) {
    sw.mutex.Lock()
    defer sw.mutex.Unlock()
    
    // Flush buffer first to ensure consistency
    if err := sw.writer.Flush(); err != nil {
        return 0, err
    }
    
    return sw.file.WriteAt(data, offset)
}

// Seek moves to a specific position in the file
func (sw *StreamWriter) Seek(offset int64, whence int) (int64, error) {
    sw.mutex.Lock()
    defer sw.mutex.Unlock()
    
    // Flush buffer first
    if err := sw.writer.Flush(); err != nil {
        return 0, err
    }
    
    pos, err := sw.file.Seek(offset, whence)
    if err != nil {
        return pos, fmt.Errorf("failed to seek to position: %w", err)
    }
    
    // Reset buffer after seeking
    sw.writer = bufio.NewWriterSize(sw.file, sw.bufferSize)
    
    return pos, nil
}

// StreamingWriteCallback is called during streaming operations
type StreamingWriteCallback func(bytesWritten int64, speed float64)

// WriteWithCallback writes with progress callback
func (sw *StreamWriter) WriteWithCallback(data []byte, callback StreamingWriteCallback) (int, error) {
    n, err := sw.Write(data)
    
    if callback != nil {
        speed := sw.GetWriteSpeed()
        callback(sw.GetBytesWritten(), speed)
    }
    
    return n, err
}

// SetBufferSize changes the buffer size (flushes current buffer first)
func (sw *StreamWriter) SetBufferSize(size int) error {
    sw.mutex.Lock()
    defer sw.mutex.Unlock()
    
    // Flush current buffer
    if err := sw.writer.Flush(); err != nil {
        return err
    }
    
    // Create new buffer with different size
    sw.writer = bufio.NewWriterSize(sw.file, size)
    sw.bufferSize = size
    
    return nil
}

// GetStats returns writing statistics
func (sw *StreamWriter) GetStats() (bytesWritten int64, speed float64, elapsed time.Duration) {
    sw.mutex.RLock()
    defer sw.mutex.RUnlock()
    
    elapsed = time.Since(sw.startTime)
    bytesWritten = sw.bytesWritten
    
    if elapsed.Seconds() > 0 {
        speed = float64(bytesWritten) / elapsed.Seconds()
    }
    
    return
}