package grpc

import (
	"context"
	"muchway/payment_service/domain"
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
	p := &domain.Payment{UserID: req.UserId, Type: req.Type, Amount: req.Amount, Status: "pending"}
	err := s.uc.CreatePayment(p)
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResponse{Id: p.ID, UserId: p.UserID, Type: p.Type, Amount: p.Amount, Status: p.Status, CreatedAt: p.CreatedAt.Unix(), UpdatedAt: p.UpdatedAt.Unix()}, nil
}

func (s *PaymentServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	p, err := s.uc.GetPaymentByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResponse{Id: p.ID, UserId: p.UserID, Type: p.Type, Amount: p.Amount, Status: p.Status, CreatedAt: p.CreatedAt.Unix(), UpdatedAt: p.UpdatedAt.Unix()}, nil
}
func (s *PaymentServer) GetAllPayments(ctx context.Context, req *pb.Empty) (*pb.PaymentsResponse, error) {
	list, err := s.uc.GetAllPayments()
	if err != nil {
		return nil, err
	}
	res := &pb.PaymentsResponse{}
	for _, p := range list {
		res.Payments = append(res.Payments, &pb.PaymentResponse{Id: p.ID, UserId: p.UserID, Type: p.Type, Amount: p.Amount, Status: p.Status, CreatedAt: p.CreatedAt.Unix(), UpdatedAt: p.UpdatedAt.Unix()})
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
