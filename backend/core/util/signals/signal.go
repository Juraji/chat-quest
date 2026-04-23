package signals

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"

	"juraji.nl/chat-quest/core/util"
)

// SignalListener defines a callback function type for handling signal events.
// Use this to process payloads when a signal is emitted.
type SignalListener[T any] func(ctx context.Context, payload T) error

// SignalResult encapsulates the result of emitting signals to multiple listeners.
// It provides synchronization via a WaitGroup and aggregates errors from all listener executions.
// wg is a WaitGroup used to synchronize waiting for all listener goroutines to complete.
// errGroup accumulates errors from each listener execution, allowing retrieval of combined results.
type SignalResult struct {
	wg       *sync.WaitGroup
	errGroup *multierror.Error
}

func (sr SignalResult) Wait() error {
	sr.wg.Wait()
	return sr.errGroup.ErrorOrNil()
}

// subscriber represents an internal structure used to manage subscribers.
// You typically won't interact with this directly.
type subscriber[T any] = struct {
	key      string
	mut      sync.Mutex
	listener SignalListener[T]
}

// Signal implements a thread-safe publish-subscribe pattern for event handling.
// Use this to create signal channels that allow components to communicate asynchronously.
type Signal[T any] struct {
	subscribers    []*subscriber[T]
	subscribersMap *util.Set[string]
	subscriberMut  sync.Mutex
}

// New returns a new instance of a Signal.
func New[T any]() *Signal[T] {
	return &Signal[T]{
		subscribersMap: util.NewSet[string](0),
		subscriberMut:  sync.Mutex{},
	}
}

// MapSignal creates a new signal that emits transformed versions of the original signal's payload.
// The transformation is applied using the provided function, which takes each emitted value
// from the source signal and returns a new value of type R. Returns a new Signal[R] instance.
func MapSignal[T any, R any](source *Signal[T], name string, transform func(T) R) *Signal[R] {
	mappedSignal := New[R]()

	source.AddListener(name, func(ctx context.Context, payload T) error {
		return mappedSignal.
			Emit(ctx, transform(payload)).
			Wait()
	})

	return mappedSignal
}

// AddListener registers your callback function to be notified when the signal is emitted.
// Provide a unique key for each listener to avoid duplicates and ensure proper removal later.
func (s *Signal[T]) AddListener(key string, listener SignalListener[T]) {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()

	if exists := s.subscribersMap.Contains(key); exists {
		panic(fmt.Sprintf("key %s already exists", key))
	}

	s.subscribersMap.Add(key)
	s.subscribers = append(s.subscribers, &subscriber[T]{
		key:      key,
		mut:      sync.Mutex{},
		listener: listener,
	})
}

// RemoveListener unregisters a previously added signal listener.
// Use this when you no longer need to be notified about this signal.
func (s *Signal[T]) RemoveListener(key string) {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()

	if removed := s.subscribersMap.Remove(key); !removed {
		return
	}

	for i, sub := range s.subscribers {
		if sub.key == key {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			return
		}
	}
}

// Emit sends a signal to all registered listeners and returns a WaitGroup.
// Use this when you need to wait for all listener functions to complete execution.
func (s *Signal[T]) Emit(ctx context.Context, payload T) SignalResult {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()
	var wg sync.WaitGroup
	var errGroup *multierror.Error

	for _, sub := range s.subscribers {
		wg.Go(func() {
			sub.mut.Lock()
			defer sub.mut.Unlock()
			err := sub.listener(ctx, payload)
			errGroup = multierror.Append(errGroup, err)
		})
	}

	return SignalResult{
		wg:       &wg,
		errGroup: errGroup,
	}
}

// EmitBG is a convenience method for emitting signals with a background context.
// Use this when you don't need request-scoped data in your listeners.
func (s *Signal[T]) EmitBG(payload T) SignalResult {
	return s.Emit(context.Background(), payload)
}

// EmitAll sends multiple related signals to all registered listeners in sequence.
// Use this when you need to emitOnGroup several related signals and want to wait for all of them to complete.
func (s *Signal[T]) EmitAll(ctx context.Context, payloads []T) SignalResult {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()
	var wg sync.WaitGroup
	var errGroup *multierror.Error

	for _, payload := range payloads {
		for _, sub := range s.subscribers {
			wg.Go(func() {
				sub.mut.Lock()
				defer sub.mut.Unlock()
				err := sub.listener(ctx, payload)
				errGroup = multierror.Append(errGroup, err)
			})
		}
	}

	return SignalResult{
		wg:       &wg,
		errGroup: errGroup,
	}
}

// EmitAllBG sends multiple related signals using a background context.
// Use this when you need to emitOnGroup several related signals without request-scoped data.
func (s *Signal[T]) EmitAllBG(payloads []T) SignalResult {
	return s.EmitAll(context.Background(), payloads)
}
