package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ai-chat-service-go/internal/middleware"
	"ai-chat-service-go/internal/models"
)

// MockDB is a mock database for testing
type MockDB struct {
	mock.Mock
}

// Begin mocks a transaction begin
func (m *MockDB) Begin() (*MockTx, error) {
	args := m.Called()
	return args.Get(0).(*MockTx), args.Error(1)
}

// Query mocks a query execution
func (m *MockDB) Query(query string, args ...interface{}) (*MockRows, error) {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*MockRows), callArgs.Error(1)
}

// QueryRow mocks a query row execution
func (m *MockDB) QueryRow(query string, args ...interface{}) *MockRow {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*MockRow)
}

// Close mocks database close
func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockTx is a mock transaction
type MockTx struct {
	mock.Mock
}

// Exec mocks an exec call
func (m *MockTx) Exec(query string, args ...interface{}) (interface{}, error) {
	callArgs := m.Called(query, args)
	return callArgs.Get(0), callArgs.Error(1)
}

// Commit mocks a transaction commit
func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

// Rollback mocks a transaction rollback
func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// MockRows is a mock result rows
type MockRows struct {
	mock.Mock
}

// Next mocks the next method
func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

// Scan mocks the scan method
func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

// Close mocks the close method
func (m *MockRows) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Err mocks the err method
func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

// MockRow is a mock result row
type MockRow struct {
	mock.Mock
}

// Scan mocks the scan method
func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

// TestCreateChat tests the CreateChat handler
func TestCreateChat(t *testing.T) {
	// Setup
	app := fiber.New()
	mockDB := new(MockDB)
	mockTx := new(MockTx)

	// Setup expectations
	mockDB.On("Begin").Return(mockTx, nil)
	mockTx.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Times(2)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Create handler
	handler := NewChatHandler(mockDB)

	// Set up the router
	app.Post("/chats", func(c *fiber.Ctx) error {
		// Mock authenticated user
		userInfo := &middleware.UserInfo{
			Email: "test@example.com",
			Name:  "Test User",
		}
		c.Locals(string(middleware.UserKey), userInfo)

		return handler.CreateChat(c)
	})

	// Create test request
	reqBody := CreateChatRequest{
		Content: "How do I configure my device?",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/chats", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var respBody CreateChatResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, respBody.ID)
	assert.Equal(t, reqBody.Content, respBody.Title)
	assert.NotEmpty(t, respBody.LastActiveDate)
	assert.NotEmpty(t, respBody.InitialMessage.ID)
	assert.Equal(t, reqBody.Content, respBody.InitialMessage.Content)
	assert.Equal(t, models.SenderTypeUser, respBody.InitialMessage.SenderType)

	// Verify expectations
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
