package grpc

import (
	"context"
	"log"
	pb "muchway/payment_service/pb"
	"muchway/payment_service/usecase"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUsecase
}

func NewPaymentServer(uc *usecase.PaymentUsecase) *PaymentServer {
	return &PaymentServer{uc: uc}
}

func (s *PaymentServer) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("Processing %s payment of %.2f for user %s", req.Type, req.Amount, req.UserId)

	p, err := s.uc.ProcessPayment(req.UserId, req.Amount, req.Type)
	if err != nil {
		log.Printf("Payment processing failed: %v", err)
		return nil, err
	}

	log.Printf("Payment processed successfully: ID=%s, Status=%s", p.ID, p.Status)
	return &pb.PaymentResponse{
		Id:        p.ID,
		UserId:    p.UserID,
		Type:      p.Type,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}, nil
}

func (s *PaymentServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	p, err := s.uc.GetPaymentByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResponse{
		Id:        p.ID,
		UserId:    p.UserID,
		Type:      p.Type,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}, nil
}

func (s *PaymentServer) GetAllPayments(ctx context.Context, req *pb.Empty) (*pb.PaymentsResponse, error) {
	list, err := s.uc.GetAllPayments()
	if err != nil {
		return nil, err
	}
	res := &pb.PaymentsResponse{}
	for _, p := range list {
		res.Payments = append(res.Payments, &pb.PaymentResponse{
			Id:        p.ID,
			UserId:    p.UserID,
			Type:      p.Type,
			Amount:    p.Amount,
			Status:    p.Status,
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		})
	}
	return res, nil
}

func (s *PaymentServer) DeletePayment(ctx context.Context, req *pb.DeletePaymentRequest) (*pb.Empty, error) {
	err := s.uc.DeletePaymentByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
