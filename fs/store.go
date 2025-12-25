package fs

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
)

// Store implements ds.IFileDataSource using the local filesystem.
type Store struct {
	baseDir string
	pathFn  PathFunc
}

// Opts configures the filesystem store.
type Opts struct {
	// BaseDir is the root directory where files are stored (required).
	BaseDir string

	// PathFn controls how a FileID maps to a relative StoragePath.
	// If nil, DefaultPathFn is used (flat: "<id>").
	PathFn PathFunc
}

// PathFunc maps a file id to a relative storage path (no leading slash).
// Examples: "uuid", "aa/bb/uuid", "tenant/123/uuid".
type PathFunc func(scope ds.StoragePath, id ds.FileID) ds.StoragePath

// FlatPathFn stores files flat under BaseDir as "<uuid>".
func FlatPathFn(scope ds.StoragePath, id ds.FileID) ds.StoragePath {
	return ds.StoragePath(string(id))
}

// ScopedPathFn stores files flat under BaseDir as "<uid>/<uuid>".
func ScopedPathFn(scope ds.StoragePath, id ds.FileID) ds.StoragePath {

	p := strings.TrimSpace(string(scope))
	if p == "" {
		return ds.StoragePath(string(id))
	}

	return ds.StoragePath(filepath.ToSlash(
		filepath.Join(p, string(id)),
	))
}

// Shard2PathFn shards by first 2+2 hex chars of the id.
// Example: "a1b2c3..." -> "a1/b2/a1b2c3..."
// If scope is provided, it becomes: "<scope>/a1/b2/a1b2c3..."
func Shard2PathFn(scope ds.StoragePath, id ds.FileID) ds.StoragePath {

	s := strings.ToLower(string(id))

	// No shard possible, keep flat (optionally under scope)
	if len(s) < 4 {
		if string(scope) == "" {
			return ds.StoragePath(s)
		}
		return ds.StoragePath(filepath.ToSlash(filepath.Join(string(scope), s)))
	}

	rel := filepath.ToSlash(filepath.Join(s[0:2], s[2:4], s))

	if string(scope) == "" {
		return ds.StoragePath(rel)
	}

	return ds.StoragePath(filepath.ToSlash(filepath.Join(string(scope), rel)))
}

// New creates a new filesystem-backed store.
func New(opts Opts) (*Store, error) {
	base := strings.TrimSpace(opts.BaseDir)
	if base == "" {
		return nil, errors.New("filestore: BaseDir is required")
	}
	base = filepath.Clean(base)

	fn := opts.PathFn
	if fn == nil {
		fn = FlatPathFn
	}

	return &Store{
		baseDir: base,
		pathFn:  fn,
	}, nil
}

// Upload writes the blob to disk and returns FileInfo + StoragePath.
// It is the caller's job (controller) to populate OriginalName/ContentType/SizeBytes in UploadInput.
func (s *Store) Upload(ctx context.Context, in ds.UploadInput, scope ds.StoragePath) (ds.FileInfo, error) {
	_ = ctx // reserved for future cancellation-aware writes if needed

	if in.Body == nil {
		return ds.FileInfo{}, &ds.InvalidFileError{Message: msg.Get("invalid_file")}
	}

	id := ds.FileID(uuid.New().String())
	rel := s.pathFn(scope, id)
	if err := validateRelPath(rel); err != nil {
		return ds.FileInfo{}, &ds.InvalidFileError{Message: msg.Get("invalid_file_path").SetArgs(err.Error())}
	}

	abs := filepath.Join(s.baseDir, filepath.FromSlash(string(rel)))
	dir := filepath.Dir(abs)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ds.FileInfo{}, err
	}

	tmp := abs + ".tmp"

	// Ensure any stale tmp is removed.
	_ = os.Remove(tmp)

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return ds.FileInfo{}, err
	}

	_, copyErr := io.Copy(f, in.Body)
	closeErr := f.Close()

	if copyErr != nil {
		_ = os.Remove(tmp)
		return ds.FileInfo{}, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return ds.FileInfo{}, closeErr
	}

	// Atomic finalize.
	if err := os.Rename(tmp, abs); err != nil {
		_ = os.Remove(tmp)
		return ds.FileInfo{}, err
	}

	return ds.FileInfo{
		ID:           id,
		Path:         rel,
		OriginalName: in.OriginalName,
		ContentType:  in.ContentType,
		SizeBytes:    in.SizeBytes,
	}, nil
}

// Open opens the stored blob for streaming. Caller must close.
func (s *Store) Open(ctx context.Context, path ds.StoragePath) (io.ReadCloser, error) {
	_ = ctx

	if err := validateRelPath(path); err != nil {
		return nil, &ds.InvalidFileError{Message: msg.Get("invalid_file_path").SetArgs(err.Error())}
	}

	abs := filepath.Join(s.baseDir, filepath.FromSlash(string(path)))
	f, err := os.Open(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &ds.NotFoundFileError{Message: msg.Get("file_not_found")}
		}
		return nil, err
	}
	return f, nil
}

// Delete removes the stored blob.
func (s *Store) Delete(ctx context.Context, path ds.StoragePath) error {
	_ = ctx

	if err := validateRelPath(path); err != nil {
		return &ds.InvalidFileError{Message: msg.Get("invalid_file_path").SetArgs(err.Error())}
	}

	abs := filepath.Join(s.baseDir, filepath.FromSlash(string(path)))
	err := os.Remove(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return &ds.NotFoundFileError{Message: msg.Get("file_not_found")}
		}
		return err
	}
	return nil
}

// validateRelPath rejects absolute paths and traversal attempts.
// Minimal safety: filestore must never allow writes outside BaseDir.
func validateRelPath(p ds.StoragePath) error {
	s := strings.TrimSpace(string(p))
	if s == "" {
		return errors.New("empty path")
	}
	// Disallow absolute paths (both Unix and Windows-ish).
	if strings.HasPrefix(s, "/") || strings.HasPrefix(s, `\`) {
		return errors.New("absolute path not allowed")
	}
	// Disallow traversal.
	clean := filepath.Clean(filepath.FromSlash(s))
	if clean == "." || clean == string(filepath.Separator) {
		return errors.New("invalid path")
	}
	parts := strings.Split(clean, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." {
			return errors.New("path traversal not allowed")
		}
	}
	return nil
}
