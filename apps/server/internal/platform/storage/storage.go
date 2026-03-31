// Package storage defines the file storage abstraction used across the application.
//
// The Provider interface decouples upload handlers from the concrete storage
// backend (local filesystem, S3, MinIO, etc).  The default LocalProvider
// stores files on disk and is sufficient for single-node deployments.
//
// Runtime boundary: LocalProvider does NOT support cluster-wide consistency.
// For multi-node deployments, swap in an S3/MinIO provider.
package storage

import "io"

// ObjectInfo contains metadata about a stored object.
type ObjectInfo struct {
	Key  string // storage-relative key (e.g. "2024/01/15/1234_report.pdf")
	Size int64
	URL  string // publicly accessible URL (path or full URL depending on provider)
}

// Provider is the storage backend abstraction.
type Provider interface {
	// Save stores the content read from r under the given key.
	// Returns metadata about the stored object.
	Save(key string, r io.Reader, size int64) (*ObjectInfo, error)

	// Open returns a reader for the given key.
	// Caller is responsible for closing the reader.
	Open(key string) (io.ReadCloser, error)

	// Delete removes the object at the given key.
	Delete(key string) error

	// PresignedURL returns a time-limited URL for direct access.
	// For local storage this returns the relative path; for cloud storage
	// it returns a signed URL.
	PresignedURL(key string, expiresSeconds int) (string, error)
}
