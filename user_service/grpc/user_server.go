package server

import (
	"context"
	"muchway/user_service/domain"
	"muchway/user_service/proto/userpb"
	"muchway/user_service/usecase"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserServer(usecase usecase.UserUsecase) *UserServer {
	return &UserServer{
		usecase: usecase,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	u := req.GetUser()
	user := &domain.User{
		ID:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Balance:  u.Balance,
		Role:     u.Role,
	}

	err := s.usecase.Register(user)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Balance:  user.Balance,
			Role:     user.Role,
		},
	}, nil
}
func (s *UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	user, err := s.usecase.Login(email, password)
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		User: &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Balance:  user.Balance,
			Role:     user.Role,
		},
	}, nil
}

func (s *UserServer) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	user, err := s.usecase.GetUserByID(req.GetId())
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByIDResponse{
		User: &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Balance:  user.Balance,
			Role:     user.Role,
		},
	}, nil
}

func (s *UserServer) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
	users, err := s.usecase.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var pbUsers []*userpb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Balance:  user.Balance,
			Role:     user.Role,
		})
	}

	return &userpb.GetAllUsersResponse{Users: pbUsers}, nil
}
