package fileops

import (
    "bufio"
    "fmt"
    "os"
)

// StreamWriter handles streaming file writing for encryption/decryption
type StreamWriter struct {
    file        *os.File
    writer      *bufio.Writer
    bytesWritten int64
    filepath    string
}

// NewStreamWriter creates a new stream writer for the given file
func NewStreamWriter(filepath string) (*StreamWriter, error) {
    file, err := os.Create(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to create file %s: %w", filepath, err)
    }
    
    return &StreamWriter{
        file:         file,
        writer:       bufio.NewWriter(file),
        bytesWritten: 0,
        filepath:     filepath,
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
        writer:       bufio.NewWriter(file),
        bytesWritten: stat.Size(),
        filepath:     filepath,
    }, nil
}

// Write writes data from the provided buffer
func (sw *StreamWriter) Write(data []byte) (int, error) {
    n, err := sw.writer.Write(data)
    sw.bytesWritten += int64(n)
    return n, err
}

// WriteChunk writes a chunk of data
func (sw *StreamWriter) WriteChunk(chunk []byte) error {
    _, err := sw.Write(chunk)
    return err
}

// GetBytesWritten returns total bytes written so far
func (sw *StreamWriter) GetBytesWritten() int64 {
    return sw.bytesWritten
}

// GetFilepath returns the file path
func (sw *StreamWriter) GetFilepath() string {
    return sw.filepath
}

// Flush flushes the buffered data to disk
func (sw *StreamWriter) Flush() error {
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