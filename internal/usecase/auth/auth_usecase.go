package auth

import (
	"fmt"
	"task-management-backend/config"
	"task-management-backend/internal/adapter/security"
	"task-management-backend/internal/domain/entity"
	"task-management-backend/internal/domain/ports"
	"time"
)

type AuthUseCase struct {
	jwtSecret string
	userRepo  ports.UserRepository
}

func NewAuthUseCase(jwtSecret string, userRepo ports.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		jwtSecret: jwtSecret,
		userRepo:  userRepo,
	}
}

func (uc *AuthUseCase) Login(username, password string) (*entity.LoginResponse, error) {
	cfg := config.GetConfig()
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	bcryptHasher := security.NewBcryptHasher()
	jwtTokenService := security.NewJWTTokenService()

	hashedPassword, err := bcryptHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Username: username,
		Password: string(hashedPassword),
	}

	existingUser, err := uc.userRepo.Upsert(user)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	token, err := jwtTokenService.Generate(uint(existingUser.ID), time.Duration(cfg.TokenDuration)*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entity.LoginResponse{
		Token:  token,
		UserID: existingUser.ID,
	}, nil
}
