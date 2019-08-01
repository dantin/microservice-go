package dbclient

import (
	"github.com/dantin/microservice-go/account/model"
	"github.com/stretchr/testify/mock"
)

// MockBoltClient is a mock implementation of a datastore client for testing purposes.
// Instead of the bolt.DB pointer, we're just putting a generic mock object from
// strechr/testify
type MockBoltClient struct {
	mock.Mock
}

// QueryAccount return a mock account.
func (m *MockBoltClient) QueryAccount(accountID string) (model.Account, error) {
	args := m.Mock.Called(accountID)
	return args.Get(0).(model.Account), args.Error(1)
}

// OpenBoltDB is a mock method.
func (m *MockBoltClient) OpenBoltDB() {
	// Does nothing
}

// Seed is a mock method.
func (m *MockBoltClient) Seed() {
	// Does nothing
}
