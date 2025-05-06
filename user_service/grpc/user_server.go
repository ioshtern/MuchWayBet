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

func (s *UserServer) GetUserByUsername(ctx context.Context, req *userpb.GetUserByUsernameRequest) (*userpb.GetUserByUsernameResponse, error) {
	user, err := s.usecase.GetUserByUsername(req.GetUsername())
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByUsernameResponse{
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

func (s *UserServer) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error) {
	user, err := s.usecase.GetUserByEmail(req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByEmailResponse{
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

func (s *UserServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	u := req.GetUser()
	user := &domain.User{
		ID:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Balance:  u.Balance,
		Role:     u.Role,
	}

	err := s.usecase.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &userpb.UpdateUserResponse{User: u}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	err := s.usecase.DeleteUser(req.GetUsername())
	if err != nil {
		return &userpb.DeleteUserResponse{Success: false}, err
	}

	return &userpb.DeleteUserResponse{Success: true}, nil
}
