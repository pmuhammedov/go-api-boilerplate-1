package user

import (
	"app/pkg/domain"
	"context"

	"github.com/google/uuid"
)

type eventSourcedRepository struct {
	streamName string
	eventStore domain.EventStore
	eventBus   domain.EventBus
}

func (r *eventSourcedRepository) Save(ctx context.Context, u *User) error {
	r.eventStore.Store(u.Changes())

	for _, event := range u.Changes() {
		r.eventBus.Publish(event.Metadata.Type, ctx, event)
	}

	return nil
}

func (r *eventSourcedRepository) Get(id uuid.UUID) *User {
	events := r.eventStore.GetStream(id, r.streamName)

	aggregateRoot := New()
	aggregateRoot.FromHistory(events)

	return aggregateRoot
}

func newEventSourcedRepository(streamName string, store domain.EventStore, bus domain.EventBus) *eventSourcedRepository {
	return &eventSourcedRepository{streamName, store, bus}
}
