package fileops

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

// StreamReader handles streaming file reading for encryption/decryption
type StreamReader struct {
    file     *os.File
    reader   *bufio.Reader
    fileSize int64
    bytesRead int64
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
    
    return &StreamReader{
        file:     file,
        reader:   bufio.NewReader(file),
        fileSize: stat.Size(),
        bytesRead: 0,
    }, nil
}

// Read reads data into the provided buffer
func (sr *StreamReader) Read(buffer []byte) (int, error) {
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

// GetProgress returns reading progress (0.0 to 1.0)
func (sr *StreamReader) GetProgress() float64 {
    if sr.fileSize == 0 {
        return 1.0
    }
    return float64(sr.bytesRead) / float64(sr.fileSize)
}

// GetBytesRead returns total bytes read so far
func (sr *StreamReader) GetBytesRead() int64 {
    return sr.bytesRead
}

// GetFileSize returns total file size
func (sr *StreamReader) GetFileSize() int64 {
    return sr.fileSize
}

// IsEOF checks if we've reached end of file
func (sr *StreamReader) IsEOF() bool {
    return sr.bytesRead >= sr.fileSize
}

// Reset resets the reader to the beginning of the file
func (sr *StreamReader) Reset() error {
    _, err := sr.file.Seek(0, io.SeekStart)
    if err != nil {
        return fmt.Errorf("failed to reset file position: %w", err)
    }
    
    sr.reader = bufio.NewReader(sr.file)
    sr.bytesRead = 0
    return nil
}

// Close closes the underlying file
func (sr *StreamReader) Close() error {
    if sr.file != nil {
        return sr.file.Close()
    }
    return nil
}