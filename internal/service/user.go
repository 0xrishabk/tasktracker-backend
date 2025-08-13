package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ryszhio/tasktracker/internal/model"
	"github.com/ryszhio/tasktracker/internal/repository"
	"github.com/ryszhio/tasktracker/internal/util"
)

type UserService struct {
	userRepo *repository.UserRepository
	timeout  time.Duration
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		timeout:  time.Duration(2) * time.Second,
	}
}

func (s *UserService) CreateUser(c context.Context, req model.RequestCreateUser) (*model.ResponseLoginUser, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("UserService.CreateUser - Starting user creation for: %s", req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Print("UserService.CreateUser - Validation Failed: missing required fields.")
		return nil, fmt.Errorf("username, email & password fields are required")
	}

	if len(req.Password) < 6 {
		log.Print("UserService.CreateUser - Validation failed: password too short.")
		return nil, fmt.Errorf("password must be atleast 6 characters long")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Printf("UserService.CreateUser - Passwording hasing failed: %v", err)
		return nil, fmt.Errorf("failed to process password")
	}

	u := &repository.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &hashedPassword,
	}

	user, err := s.userRepo.CreateUser(c, u)
	if err != nil {
		log.Printf("UserService.CreateUser - Database Error: %v", err)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("username or email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	log.Printf("UserService.CreateUser - User created successfully in database: %s", user.ID.String())

	return &model.ResponseLoginUser{
		AccessToken: "something",
		Username:    user.Username,
		ID:          user.ID.String(),
	}, nil
}

func (s *UserService) Login(c context.Context, req model.RequestLoginUser) (*model.ResponseLoginUser, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("UserService.Login - Starting login attempt for email: %s", req.Email)

	user, err := s.userRepo.GetUserByEmail(c, req.Email)
	if err != nil {
		log.Printf("UserService.Login - Database error: %v", err)
		return nil, fmt.Errorf("failed to authenticate user")
	}

	if user == nil {
		log.Printf("UserService.Login - User not found for email: %s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	if user.PasswordHash == nil {
		log.Printf("UserService.Login - User has no password hash: %s", req.Email)
		return nil, fmt.Errorf("invalid user account")
	}

	err = util.CheckPassword(req.Password, *user.PasswordHash)
	if err != nil {
		log.Printf("UserService.Login - Password check failed for the user: %s", user.ID.String())
		return nil, fmt.Errorf("invalid email or password")
	}

	log.Printf("UserService.Login - Password verification successful for the user: %s", user.ID.String())

	return &model.ResponseLoginUser{AccessToken: "", Username: user.Username, ID: user.ID.String()}, nil
}

func (s *UserService) GetUserByID(c context.Context, userID string) (*repository.User, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("UserService.GetUserByID - Starting attempt to fetch user by ID: %s", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("UserService.GetUserByID - UUID parsing error: %v", err)
		return nil, err
	}

	u, err := s.userRepo.GetUserByID(c, uid)
	if err != nil {
		log.Printf("UserService.GetUserByID - Database error: %v", err)
		return nil, err
	}

	log.Printf("UserService.GetUserByID - Sucessfuly fetched user by ID: %s", userID)

	return u, nil
}

func (s *UserService) DeleteUser(c context.Context, userID string) error {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("UserService.DeleteUser - Starting attempt to delete user by ID: %s", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("UserService.GetUserByID - UUID parsing error: %v", err)
		return err
	}

	err = s.userRepo.DeleteUser(c, uid)
	if err != nil {
		log.Printf("UserService.DeleteUser - Database Error: %v", err)
		return err
	}

	log.Printf("UserService.DeleteUser - Successfully deleted user with ID: %s", userID)
	return nil
}

func (s *UserService) UpdateUsername(c context.Context, userID string, newUsername string) (*model.ResponseLoginUser, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("UserService.UpdateUsername - Starting attempt to update username by ID: %s", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("UserService.UpdateUsername - UUID parsing error: %v", err)
		return nil, err
	}

	user, err := s.userRepo.UpdateUsername(c, uid, newUsername)
	if err != nil {
		log.Printf("UserService.UpdateUsername - Database error: %v", err)
		return nil, err
	}

	log.Printf("UserService.UpdateUsername - Successfully updated username to %s", user.Username)

	return &model.ResponseLoginUser{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}
