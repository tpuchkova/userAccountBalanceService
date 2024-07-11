package service

import (
	"testing"

	handlermodel "userAccountBalanceService/pkg/handler/model"
	repositorymodel "userAccountBalanceService/pkg/repository/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBalanceRepo struct {
	mock.Mock
}

func (m *MockBalanceRepo) GetBalance() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockBalanceRepo) SaveBalance(balance float64) error {
	args := m.Called(balance)
	return args.Error(0)
}

type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) GetLatestTransactions(count int) ([]repositorymodel.Transaction, error) {
	args := m.Called(count)
	return args.Get(0).([]repositorymodel.Transaction), args.Error(1)
}

func (m *MockTransactionRepo) CancelByIDs(ids []int) error {
	args := m.Called(ids)
	return args.Error(0)
}

func (m *MockTransactionRepo) AddTransaction(transaction repositorymodel.Transaction) (repositorymodel.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(repositorymodel.Transaction), args.Error(1)
}

func TestAddTransaction_InCaseOfWin(t *testing.T) {
	// Arrange
	balanceRepo := new(MockBalanceRepo)
	transactionRepo := new(MockTransactionRepo)

	service := &TransactionService{
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}

	transaction := handlermodel.Transaction{
		State:         "win",
		Amount:        "100.00",
		TransactionID: "txn_12345",
	}
	expectedTransaction := repositorymodel.Transaction{
		State:         "win",
		Amount:        100.00,
		TransactionID: "txn_12345",
	}

	balanceRepo.On("GetBalance").Return(500.00, nil)
	transactionRepo.On("AddTransaction", expectedTransaction).Return(repositorymodel.Transaction{
		ID:            1,
		State:         "win",
		Amount:        100.00,
		TransactionID: "txn_12345",
	}, nil)
	balanceRepo.On("SaveBalance", 600.00).Return(nil)

	// Act
	result, err := service.AddTransaction(transaction)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "win", result.State)
	assert.Equal(t, 100.00, result.Amount)
	assert.Equal(t, "txn_12345", result.TransactionID)
	assert.Equal(t, 600.00, result.UserBalance)

	balanceRepo.AssertExpectations(t)
	transactionRepo.AssertExpectations(t)
}

func TestAddTransaction_InCaseOfLost(t *testing.T) {
	// Arrange
	balanceRepo := new(MockBalanceRepo)
	transactionRepo := new(MockTransactionRepo)

	service := &TransactionService{
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}

	transaction := handlermodel.Transaction{
		State:         "lost",
		Amount:        "100.00",
		TransactionID: "txn_12345",
	}
	expectedTransaction := repositorymodel.Transaction{
		State:         "lost",
		Amount:        100.00,
		TransactionID: "txn_12345",
	}

	balanceRepo.On("GetBalance").Return(500.00, nil)
	transactionRepo.On("AddTransaction", expectedTransaction).Return(repositorymodel.Transaction{
		ID:            1,
		State:         "win",
		Amount:        100.00,
		TransactionID: "txn_12345",
	}, nil)
	balanceRepo.On("SaveBalance", 400.00).Return(nil)

	// Act
	result, err := service.AddTransaction(transaction)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "win", result.State)
	assert.Equal(t, 100.00, result.Amount)
	assert.Equal(t, "txn_12345", result.TransactionID)
	assert.Equal(t, 400.00, result.UserBalance)

	balanceRepo.AssertExpectations(t)
	transactionRepo.AssertExpectations(t)
}

func TestCancelLatestOddTransactions_InCaseOfEnoughBalance(t *testing.T) {
	// Arrange
	balanceRepo := new(MockBalanceRepo)
	transactionRepo := new(MockTransactionRepo)

	service := &TransactionService{
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}

	transactions := []repositorymodel.Transaction{
		{ID: 1, State: "win", Amount: 100},
		{ID: 2, State: "lost", Amount: 50},
		{ID: 3, State: "win", Amount: 150},
		{ID: 4, State: "lost", Amount: 200},
		{ID: 5, State: "lost", Amount: 100},
		{ID: 6, State: "win", Amount: 200},
		{ID: 7, State: "win", Amount: 100},
		{ID: 8, State: "win", Amount: 50},
		{ID: 9, State: "lost", Amount: 100},
		{ID: 10, State: "win", Amount: 100},
	}
	initialBalance := 500.0
	expectedCanceledIDs := []int{2, 4, 6, 8, 10}
	newBalance := initialBalance + 50 + 200 - 200 - 50 - 100

	transactionRepo.On("GetLatestTransactions", 19).Return(transactions, nil)
	balanceRepo.On("GetBalance").Return(initialBalance, nil)
	transactionRepo.On("CancelByIDs", expectedCanceledIDs).Return(nil)
	balanceRepo.On("SaveBalance", newBalance).Return(nil)

	// Act
	err := service.CancelLatestOddTransactions()

	// Assert
	assert.NoError(t, err)
	transactionRepo.AssertExpectations(t)
	balanceRepo.AssertExpectations(t)
}

func TestCancelLatestOddTransactions_InCaseOfTooSmallBalance(t *testing.T) {
	// Arrange
	balanceRepo := new(MockBalanceRepo)
	transactionRepo := new(MockTransactionRepo)

	service := &TransactionService{
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}

	transactions := []repositorymodel.Transaction{
		{ID: 1, State: "win", Amount: 100},
		{ID: 2, State: "win", Amount: 70},
		{ID: 3, State: "win", Amount: 150},
		{ID: 4, State: "win", Amount: 200},
		{ID: 5, State: "win", Amount: 100},
		{ID: 6, State: "win", Amount: 200},
		{ID: 7, State: "win", Amount: 100},
		{ID: 8, State: "win", Amount: 50},
		{ID: 9, State: "win", Amount: 100},
		{ID: 10, State: "win", Amount: 100},
	}
	initialBalance := 500.0
	expectedCanceledIDs := []int{2, 4, 6}
	newBalance := initialBalance - 70 - 200 - 200

	transactionRepo.On("GetLatestTransactions", 19).Return(transactions, nil)
	balanceRepo.On("GetBalance").Return(initialBalance, nil)
	transactionRepo.On("CancelByIDs", expectedCanceledIDs).Return(nil)
	balanceRepo.On("SaveBalance", newBalance).Return(nil)

	// Act
	err := service.CancelLatestOddTransactions()

	// Assert
	assert.NoError(t, err)
	transactionRepo.AssertExpectations(t)
	balanceRepo.AssertExpectations(t)
}

func TestCancelLatestOddTransactions_InCaseOfNoTransactionToCancel(t *testing.T) {
	// Arrange
	balanceRepo := new(MockBalanceRepo)
	transactionRepo := new(MockTransactionRepo)

	service := &TransactionService{
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}

	initialBalance := 500.0

	transactionRepo.On("GetLatestTransactions", 19).Return([]repositorymodel.Transaction{}, nil)
	balanceRepo.On("GetBalance").Return(initialBalance, nil)

	// Act
	err := service.CancelLatestOddTransactions()

	// Assert
	assert.NoError(t, err)
	transactionRepo.AssertExpectations(t)
	balanceRepo.AssertExpectations(t)
}
