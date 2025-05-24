package usecase

import (
	"bet_service/domain"
	"bet_service/repository"
	"os"
	"testing"
)

// --- Моки ---

type mockBetRepo struct {
	createCalled bool
	updateCalled bool
	deleteCalled bool
	getByIDFunc  func(id string) (*domain.Bet, error)
}

func (m *mockBetRepo) Create(bet *domain.Bet) error {
	m.createCalled = true
	return nil
}

func (m *mockBetRepo) Update(bet *domain.Bet) error {
	m.updateCalled = true
	return nil
}

func (m *mockBetRepo) Delete(id string) error {
	m.deleteCalled = true
	return nil
}

func (m *mockBetRepo) GetByID(id string) (*domain.Bet, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return &domain.Bet{ID: id}, nil
}

func (m *mockBetRepo) GetByUserID(userID string) ([]*domain.Bet, error) {
	return []*domain.Bet{
		{ID: "bet123", UserID: userID},
	}, nil
}

type mockPublisher struct {
	created bool
	updated bool
	deleted bool
}

func (m *mockPublisher) PublishBetCreated(bet *domain.Bet) error {
	m.created = true
	return nil
}

func (m *mockPublisher) PublishBetUpdated(bet *domain.Bet) error {
	m.updated = true
	return nil
}

func (m *mockPublisher) PublishBetDeleted(bet *domain.Bet) error {
	m.deleted = true
	return nil
}

// --- Тесты ---

func TestCreateBet(t *testing.T) {
	mockRepo := &mockBetRepo{}
	mockPub := &mockPublisher{}
	uc := NewBetUsecase(mockRepo, mockPub)

	bet := &domain.Bet{ID: "bet123", UserID: "user1"}

	err := uc.CreateBet(bet)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !mockRepo.createCalled {
		t.Error("expected Create to be called")
	}
	if !mockPub.created {
		t.Error("expected PublishBetCreated to be called")
	}
}

func TestUpdateBet(t *testing.T) {
	mockRepo := &mockBetRepo{}
	mockPub := &mockPublisher{}
	uc := NewBetUsecase(mockRepo, mockPub)

	bet := &domain.Bet{ID: "bet123", UserID: "user1"}

	err := uc.UpdateBet(bet)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !mockRepo.updateCalled {
		t.Error("expected Update to be called")
	}
	if !mockPub.updated {
		t.Error("expected PublishBetUpdated to be called")
	}
}

func TestDeleteBet(t *testing.T) {
	mockRepo := &mockBetRepo{}
	mockPub := &mockPublisher{}
	mockRepo.getByIDFunc = func(id string) (*domain.Bet, error) {
		return &domain.Bet{ID: id}, nil
	}

	uc := NewBetUsecase(mockRepo, mockPub)

	err := uc.DeleteBet("bet123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !mockRepo.deleteCalled {
		t.Error("expected Delete to be called")
	}
	if !mockPub.deleted {
		t.Error("expected PublishBetDeleted to be called")
	}
}
func TestMain(m *testing.M) {
	repository.InitRedisClient("localhost:6379", "", 0)
	os.Exit(m.Run())
}
