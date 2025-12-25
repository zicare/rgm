package ds

import (
	"context"
	"io"
)

// FileID is the stable identifier exposed to the app.
type FileID string

// StoragePath is the backend-relative path/key where the blob lives.
// Example: "files/ab/cd/uuid.bin" or "uuid".
type StoragePath string

// FileInfo is IO-level information returned after upload.
// This is NOT a database record.
// The full path to the file in the filesystem would be BaseDir + Path
type FileInfo struct {
	ID           FileID
	Path         StoragePath
	OriginalName string
	ContentType  string
	SizeBytes    int64
}

// UploadInput is produced by the controller after parsing the HTTP request.
// It is transport-agnostic.
type UploadInput struct {
	OriginalName string
	ContentType  string
	SizeBytes    int64
	Body         io.Reader
}

// IFileDataSource provides blob storage operations only.
type IFileDataSource interface {

	// Upload stores the blob and returns its identifier and storage path.
	Upload(ctx context.Context, in UploadInput, scope StoragePath) (FileInfo, error)

	// Open returns a stream to read the blob at the given path.
	// Caller MUST close the returned reader.
	Open(ctx context.Context, path StoragePath) (io.ReadCloser, error)

	// Delete removes the blob at the given path.
	Delete(ctx context.Context, path StoragePath) error
}
