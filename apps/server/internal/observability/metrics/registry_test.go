package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Fatal("expected non-nil registry")
	}
	if reg.Gatherer() == nil {
		t.Fatal("expected non-nil gatherer")
	}
}

func TestRegistry_MustRegister(t *testing.T) {
	reg := NewRegistry()
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A test counter",
	})
	reg.MustRegister(counter)

	// Duplicate registration should panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate registration")
		}
	}()
	reg.MustRegister(counter)
}

func TestRegistry_RegisterGoCollector(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterGoCollector()

	mfs, err := reg.Gatherer().Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}
	found := false
	for _, mf := range mfs {
		if mf.GetName() == "go_goroutines" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected go_goroutines metric after RegisterGoCollector")
	}
}
