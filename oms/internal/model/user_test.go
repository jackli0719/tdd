package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUserJSONSerialization(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Phone:     "1234567890",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Serialize to JSON
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	// Deserialize from JSON
	var unmarshaled User
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal user: %v", err)
	}

	if unmarshaled.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, unmarshaled.ID)
	}
	if unmarshaled.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, unmarshaled.Username)
	}
	if unmarshaled.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, unmarshaled.Email)
	}
}

func TestUserTableName(t *testing.T) {
	user := User{}
	if user.TableName() != "users" {
		t.Errorf("expected table name 'users', got '%s'", user.TableName())
	}
}

func TestCreateUserRequestValidation(t *testing.T) {
	req := CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var unmarshaled CreateUserRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if unmarshaled.Username != req.Username {
		t.Errorf("expected username %s, got %s", req.Username, unmarshaled.Username)
	}
}
