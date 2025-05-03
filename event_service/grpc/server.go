package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"muchway/event_service/domain"
	pb "muchway/event_service/proto"
	"muchway/event_service/usecase"
)

type Server struct {
	uc usecase.EventUseCase
	pb.UnimplementedEventServiceServer
}

func NewGRPCServer(uc usecase.EventUseCase) *Server {
	return &Server{uc: uc}
}

func toProto(e *domain.Event) *pb.Event {
	var wid string
	if e.WinnerID != nil {
		wid = *e.WinnerID
	}
	return &pb.Event{
		Id:        e.ID,
		Name:      e.Name,
		StartTime: e.StartTime.Format(time.RFC3339),
		Status:    e.Status,
		WinnerId:  wid,
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
		UpdatedAt: e.CreatedAt.Format(time.RFC3339),
	}
}

func fromProto(p *pb.Event) (*domain.Event, error) {
	if p == nil {
		return nil, status.Errorf(codes.InvalidArgument, "event must be provided")
	}
	st, err := time.Parse(time.RFC3339, p.StartTime)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid start_time: %v", err)
	}
	var wid *string
	if p.WinnerId != "" {
		wid = &p.WinnerId
	}
	return &domain.Event{
		ID:        p.Id,
		Name:      p.Name,
		StartTime: st,
		Status:    p.Status,
		WinnerID:  wid,
	}, nil
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	if req.Event == nil {
		return nil, status.Errorf(codes.InvalidArgument, "event must be provided")
	}
	e, err := fromProto(req.Event)
	if err != nil {
		return nil, err
	}
	saved, err := s.uc.CreateEvent(ctx, e)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create event: %v", err)
	}
	return &pb.CreateEventResponse{Event: toProto(saved)}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	e, err := s.uc.GetEvent(ctx, req.Id)
	return &pb.GetEventResponse{Event: toProto(e)}, err
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	if req.Event == nil {
		return nil, status.Errorf(codes.InvalidArgument, "event must be provided")
	}
	e, err := fromProto(req.Event)
	if err != nil {
		return nil, err
	}

	updated, err := s.uc.UpdateEvent(ctx, e)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not update event: %v", err)
	}
	return &pb.UpdateEventResponse{Event: toProto(updated)}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := s.uc.DeleteEvent(ctx, req.Id)
	return &pb.DeleteEventResponse{}, err
}

func (s *Server) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	list, err := s.uc.ListEvents(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListEventsResponse{}
	for _, e := range list {
		resp.Events = append(resp.Events, toProto(e))
	}
	return resp, nil
}
