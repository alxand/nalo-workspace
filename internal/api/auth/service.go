package auth

import (
	"errors"
	"time"

	"github.com/alxand/nalo-workspace/internal/config"
	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// Service handles authentication operations
type Service struct {
	jwtConfig config.JWTConfig
	userRepo  interfaces.UserInterface
	logger    *zap.Logger
}

// NewService creates a new auth service
func NewService(jwtConfig config.JWTConfig, userRepo interfaces.UserInterface, logger *zap.Logger) *Service {
	return &Service{
		jwtConfig: jwtConfig,
		userRepo:  userRepo,
		logger:    logger,
	}
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email     string          `json:"email" validate:"required,email"`
	Username  string          `json:"username" validate:"required,min=3,max=50"`
	Password  string          `json:"password" validate:"required,min=8"`
	FirstName string          `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string          `json:"last_name" validate:"required,min=2,max=50"`
	Role      models.UserRole `json:"role" validate:"required,oneof=admin user manager"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token     string       `json:"token"`
	Type      string       `json:"type"`
	User      *models.User `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// Authenticate validates user credentials and returns user if valid
func (s *Service) Authenticate(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !user.CheckPassword(password) {
		s.logger.Error("Invalid password for user", zap.String("email", email))
		return nil, errors.New("invalid credentials")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		s.logger.Error("Failed to update last login", zap.Int64("user_id", user.ID), zap.Error(err))
		// Don't return error, just log it
	}

	return user, nil
}

// Register creates a new user account
func (s *Service) Register(req *RegisterRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	exists, err = s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Create new user
	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password, // Will be hashed by GORM hook
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("Failed to create user", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	s.logger.Info("User registered successfully", zap.String("email", req.Email), zap.Int64("user_id", user.ID))
	return user, nil
}

// Login authenticates user and returns JWT token
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.Authenticate(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	token, expiresAt, err := s.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		Type:      "Bearer",
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateJWT generates a new JWT token for a user
func (s *Service) GenerateJWT(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.jwtConfig.Expiration)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "nalo-workspace",
		"aud":      "nalo-workspace-users",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateJWT validates a JWT token and returns the claims
func (s *Service) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractUserID extracts user ID from JWT claims
func (s *Service) ExtractUserID(claims jwt.MapClaims) (int64, error) {
	if userID, ok := claims["user_id"].(float64); ok {
		return int64(userID), nil
	}
	return 0, errors.New("user_id not found in token claims")
}

// ExtractUserRole extracts user role from JWT claims
func (s *Service) ExtractUserRole(claims jwt.MapClaims) (models.UserRole, error) {
	if role, ok := claims["role"].(string); ok {
		return models.UserRole(role), nil
	}
	return "", errors.New("role not found in token claims")
}

// GetUserFromToken validates token and returns the user
func (s *Service) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	userID, err := s.ExtractUserID(claims)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return user, nil
}
