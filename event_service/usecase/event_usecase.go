package usecase

import (
	"context"
	"log"

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
}

func NewEventUseCase(r repository.EventRepository, p rabbitmq.Publisher) EventUseCase {
	return &eventUseCase{repo: r, publisher: p}
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
	return uc.repo.Get(ctx, id)
}

func (uc *eventUseCase) UpdateEvent(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	return uc.repo.Update(ctx, e)
}

func (uc *eventUseCase) DeleteEvent(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *eventUseCase) ListEvents(ctx context.Context) ([]*domain.Event, error) {
	return uc.repo.List(ctx)
}
