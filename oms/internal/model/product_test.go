package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProductJSONSerialization(t *testing.T) {
	now := time.Now()
	product := Product{
		ID:        1,
		Name:      "Test Product",
		Price:     99.99,
		Stock:     100,
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("failed to marshal product: %v", err)
	}

	var unmarshaled Product
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal product: %v", err)
	}

	if unmarshaled.ID != product.ID {
		t.Errorf("expected ID %d, got %d", product.ID, unmarshaled.ID)
	}
	if unmarshaled.Name != product.Name {
		t.Errorf("expected name %s, got %s", product.Name, unmarshaled.Name)
	}
	if unmarshaled.Price != product.Price {
		t.Errorf("expected price %f, got %f", product.Price, unmarshaled.Price)
	}
}

func TestProductTableName(t *testing.T) {
	product := Product{}
	if product.TableName() != "products" {
		t.Errorf("expected table name 'products', got '%s'", product.TableName())
	}
}

func TestCreateProductRequestValidation(t *testing.T) {
	req := CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var unmarshaled CreateProductRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if unmarshaled.Name != req.Name {
		t.Errorf("expected name %s, got %s", req.Name, unmarshaled.Name)
	}
}
