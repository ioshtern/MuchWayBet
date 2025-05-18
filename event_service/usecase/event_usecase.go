package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"muchway/event_service/domain"
	"muchway/event_service/rabbitmq"
	"muchway/event_service/repository"
)

type EventUseCase interface {
	CreateEvent(ctx context.Context, e *domain.Event) (*domain.Event, error)
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, e *domain.Event) (*domain.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context) ([]*domain.Event, error)
}

type eventUseCase struct {
	repo      repository.EventRepository
	publisher rabbitmq.Publisher
	rdb       *redis.Client
}

func NewEventUseCase(r repository.EventRepository, p rabbitmq.Publisher, rdb *redis.Client) EventUseCase {
	return &eventUseCase{repo: r, publisher: p, rdb: rdb}
}

func (uc *eventUseCase) CreateEvent(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	saved, err := uc.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	if err := uc.publisher.Publish("events", "created", saved); err != nil {
		log.Printf("publish error: %v", err)
	}
	return saved, nil
}

func (uc *eventUseCase) GetEvent(ctx context.Context, id string) (*domain.Event, error) {
	key := "event:" + id

	if raw, err := uc.rdb.Get(ctx, key).Result(); err == nil {
		var ev domain.Event
		if err := json.Unmarshal([]byte(raw), &ev); err == nil {
			log.Printf("Cache hit for event %s", id)
			return &ev, nil
		}
	}
	log.Printf("Cache miss for event %s, loading from DB", id)

	// Cache miss -> DB
	evPtr, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Populate cache
	if buf, err := json.Marshal(evPtr); err == nil {
		uc.rdb.Set(ctx, key, buf, 5*time.Minute)
		log.Printf("Cached event %s", id)
	}
	return evPtr, nil
}

func (uc *eventUseCase) UpdateEvent(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	updated, err := uc.repo.Update(ctx, e)
	if err != nil {
		return nil, err
	}
	uc.rdb.Del(ctx, "event:"+e.ID)
	log.Printf("Invalidated cache for event %s", e.ID)
	return updated, nil
}

func (uc *eventUseCase) DeleteEvent(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	uc.rdb.Del(ctx, "event:"+id)
	log.Printf("Invalidated cache for event %s", id)
	return nil
}

func (uc *eventUseCase) ListEvents(ctx context.Context) ([]*domain.Event, error) {
	return uc.repo.List(ctx)
}
