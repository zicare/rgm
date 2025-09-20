package ds

import "fmt"

var deferredInits []func() error

// Queue a startup task to run later.
func DeferInit(f func() error) { deferredInits = append(deferredInits, f) }

// Run all queued tasks once. Panics on first error.
func RunDeferredInits() {
	for i, f := range deferredInits {
		if err := f(); err != nil {
			panic(fmt.Errorf("deferred init %d: %w", i+1, err))
		}
	}
	// prevent re-run
	deferredInits = nil
}
