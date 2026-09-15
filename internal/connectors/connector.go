package connectors

import (
	"context"
)

type Config map[string]string

type Resource struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Connector interface {
	Name() string
	Connect(ctx context.Context, config Config) error
	Disconnect(ctx context.Context) error
	ListResources(ctx context.Context) ([]Resource, error)
}

type Registry struct {
	connectors map[string]Connector
}

func NewRegistry() *Registry {
	return &Registry{
		connectors: make(map[string]Connector),
	}
}

func (r *Registry) Register(connector Connector) {
	r.connectors[connector.Name()] = connector
}

func (r *Registry) Get(name string) (Connector, bool) {
	connector, ok := r.connectors[name]
	return connector, ok
}

func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.connectors))

	for name := range r.connectors {
		names = append(names, name)
	}

	return names
}
