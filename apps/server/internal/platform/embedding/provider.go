package embedding

import "context"

// Provider defines the text embedding service interface
type Provider interface {
	// Embed returns vector representations of the input texts
	Embed(ctx context.Context, texts []string) ([][]float32, error)

	// Dimension returns the dimension of the vectors produced by this provider
	Dimension() int

	// HealthCheck performs a health check on the provider
	HealthCheck(ctx context.Context) error
}

// Config is the common configuration for embedding providers
type Config struct {
	Provider string `yaml:"provider" json:"provider"`
}
