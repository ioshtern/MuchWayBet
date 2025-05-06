package grpc

import (
	"bet_service/domain"
	"bet_service/muchway/bet_service/proto/betpb"
	"bet_service/transport/rabbitmq"
	"bet_service/usecase"
	"context"

	"github.com/google/uuid"
)

type BetServer struct {
	betpb.UnimplementedBetServiceServer
	usecase   *usecase.BetUsecase
	publisher *rabbitmq.Publisher
}

func NewBetServer(u *usecase.BetUsecase, p *rabbitmq.Publisher) *BetServer {
	return &BetServer{usecase: u, publisher: p}
}


func (s *BetServer) CreateBet(ctx context.Context, req *betpb.CreateBetRequest) (*betpb.CreateBetResponse, error) {
	bet := &domain.Bet{
		ID:      uuid.New().String(),
		UserID:  req.Bet.UserId,
		EventID: req.Bet.EventId,
		Amount:  req.Bet.Amount,
		Odds:    req.Bet.Odds,
		Status:  "pending",
	}

	if err := s.usecase.CreateBet(bet); err != nil {
		return nil, err
	}

	_ = s.publisher.PublishBetCreated(bet)

	return &betpb.CreateBetResponse{
		Bet: &betpb.Bet{
			Id:      bet.ID,
			UserId:  bet.UserID,
			EventId: bet.EventID,
			Amount:  bet.Amount,
			Odds:    bet.Odds,
			Status:  bet.Status,
			Payout:  bet.Payout,
		},
	}, nil
}


func (s *BetServer) GetBetByID(ctx context.Context, req *betpb.GetBetByIDRequest) (*betpb.GetBetByIDResponse, error) {
	bet, err := s.usecase.GetBetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &betpb.GetBetByIDResponse{
		Bet: &betpb.Bet{
			Id:      bet.ID,
			UserId:  bet.UserID,
			EventId: bet.EventID,
			Amount:  bet.Amount,
			Odds:    bet.Odds,
			Status:  bet.Status,
			Payout:  bet.Payout,
		},
	}, nil
}


func (s *BetServer) GetBetsByUserID(ctx context.Context, req *betpb.GetBetsByUserIDRequest) (*betpb.GetBetsByUserIDResponse, error) {
	bets, err := s.usecase.GetBetsByUserID(req.UserId)
	if err != nil {
		return nil, err
	}

	var pbBets []*betpb.Bet
	for _, bet := range bets {
		pbBets = append(pbBets, &betpb.Bet{
			Id:      bet.ID,
			UserId:  bet.UserID,
			EventId: bet.EventID,
			Amount:  bet.Amount,
			Odds:    bet.Odds,
			Status:  bet.Status,
			Payout:  bet.Payout,
		})
	}

	return &betpb.GetBetsByUserIDResponse{Bets: pbBets}, nil
}


func (s *BetServer) UpdateBet(ctx context.Context, req *betpb.UpdateBetRequest) (*betpb.UpdateBetResponse, error) {
	bet := &domain.Bet{
		ID:      req.Bet.Id,
		UserID:  req.Bet.UserId,
		EventID: req.Bet.EventId,
		Amount:  req.Bet.Amount,
		Odds:    req.Bet.Odds,
		Status:  req.Bet.Status,
		Payout:  req.Bet.Payout,
	}

	if err := s.usecase.UpdateBet(bet); err != nil {
		return nil, err
	}

	return &betpb.UpdateBetResponse{
		Bet: &betpb.Bet{
			Id:      bet.ID,
			UserId:  bet.UserID,
			EventId: bet.EventID,
			Amount:  bet.Amount,
			Odds:    bet.Odds,
			Status:  bet.Status,
			Payout:  bet.Payout,
		},
	}, nil
}


func (s *BetServer) DeleteBet(ctx context.Context, req *betpb.DeleteBetRequest) (*betpb.DeleteBetResponse, error) {
	err := s.usecase.DeleteBet(req.Id)
	if err != nil {
		return nil, err
	}

	return &betpb.DeleteBetResponse{
		Success: true,
	}, nil
}
