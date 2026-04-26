package service

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"golang/internal/config"
	"golang/internal/model"
)

type mockUserRepo struct {
	users   map[int]*model.User
	byEmail map[string]*model.User
	nextID  int
}

func (m *mockUserRepo) Create(user *model.User) error {
	if _, exists := m.byEmail[user.Email]; exists {
		return ErrEmailExists
	}
	m.nextID++
	user.ID = m.nextID
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(email string) (*model.User, error) {
	user, ok := m.byEmail[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepo) GetByID(id int) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[int]*model.User),
		byEmail: make(map[string]*model.User),
		nextID:  0,
	}
}

func TestRegister_Success(t *testing.T) {
	repo := newMockUserRepo()
	cfg := &config.Config{JWTSecret: "test", JWTExpirationHours: 1}
	svc := NewAuthService(repo, cfg)

	req := &model.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	user, err := svc.Register(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	// check password hashed
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		t.Error("password was not correctly hashed")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	cfg := &config.Config{}
	svc := NewAuthService(repo, cfg)

	req := &model.RegisterRequest{Email: "dup@example.com", Password: "pass"}
	_, err := svc.Register(req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Register(req)
	if err != ErrEmailExists {
		t.Errorf("expected ErrEmailExists, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpirationHours: 1}
	svc := NewAuthService(repo, cfg)

	// Создаём пользователя через сервис, чтобы пароль был хеширован
	req := &model.RegisterRequest{Email: "login@example.com", Password: "mypassword"}
	_, err := svc.Register(req)
	if err != nil {
		t.Fatal(err)
	}

	token, err := svc.Login(&model.LoginRequest{Email: "login@example.com", Password: "mypassword"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	// Проверяем токен
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("invalid token: %v", err)
	}
	if claims.Subject != "1" {
		t.Errorf("expected subject '1', got '%s'", claims.Subject)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := newMockUserRepo()
	cfg := &config.Config{JWTSecret: "test"}
	svc := NewAuthService(repo, cfg)

	svc.Register(&model.RegisterRequest{Email: "x@x.com", Password: "right"})
	_, err := svc.Login(&model.LoginRequest{Email: "x@x.com", Password: "wrong"})
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, &config.Config{})
	_, err := svc.Login(&model.LoginRequest{Email: "no@user.com", Password: "pass"})
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc := NewAuthService(nil, &config.Config{JWTSecret: "test"})
	_, err := svc.ValidateToken("garbage")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}
