package service

import (
	"errors"
	"time"

	"oms/internal/model"
	"oms/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// LoginResponse is the response for successful login
type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(c *gin.Context, req model.LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(c *gin.Context, req model.RegisterRequest) (*model.User, error) {
	// Check if username exists
	existingUser, err := s.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil && existingUser.ID > 0 {
		return nil, errors.New("用户名已存在")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour * 7).Unix(), // 7 days
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名方法不正确")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDVal, hasUserID := claims["user_id"]
		usernameVal, hasUsername := claims["username"]

		if !hasUserID || !hasUsername {
			return 0, "", errors.New("令牌缺少必要字段")
		}

		userID, ok := userIDVal.(float64)
		if !ok {
			return 0, "", errors.New("令牌用户ID格式错误")
		}

		username, ok := usernameVal.(string)
		if !ok {
			return 0, "", errors.New("令牌用户名格式错误")
		}

		return int64(userID), username, nil
	}

	return 0, "", errors.New("令牌无效")
}

// GenerateTokenForUser generates a token for an existing user (for testing)
func (s *AuthService) GenerateTokenForUser(userID int64, username string) (string, error) {
	return s.generateToken(userID, username)
}

var (
	ErrUnauthorized = errors.New("未授权")
	ErrInvalidToken = errors.New("无效的令牌")
)
