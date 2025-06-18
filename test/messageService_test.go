package test

import (
	"generalChat/internal/model"
	"generalChat/internal/service"

	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveMessage(message model.Message) {
	m.Called(message)
}

func (m *MockRepo) GetLastNMessages(n int64) []model.Message {
	args := m.Called(n)
	return args.Get(0).([]model.Message)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) AddMessage(msg model.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockCache) GetRecentMessages() ([]model.Message, error) {
	args := m.Called()
	return args.Get(0).([]model.Message), args.Error(1)
}

func TestGetLast20Messages_FromCache(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)

	cachedMsgs := make([]model.Message, 20)
	mockCache.On("GetRecentMessages").Return(cachedMsgs, nil)

	mService := service.MessageService{ChatRepo: mockRepo, ChatCache: mockCache}
	result := mService.GetLast20Messages()

	assert.Equal(t, cachedMsgs, result)
	mockRepo.AssertNotCalled(t, "GetLastNMessages", mock.Anything)
}

func TestGetLast20Messages_FromRepo(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)

	cachedMsgs := make([]model.Message, 5)
	dbMsgs := make([]model.Message, 20)

	mockCache.On("GetRecentMessages").Return(cachedMsgs, nil)
	mockRepo.On("GetLastNMessages", int64(20)).Return(dbMsgs)
	mockCache.On("AddMessage", mock.Anything).Return(nil)

	mService := service.MessageService{ChatRepo: mockRepo, ChatCache: mockCache}
	result := mService.GetLast20Messages()

	assert.Equal(t, dbMsgs, result)
	mockRepo.AssertCalled(t, "GetLastNMessages", int64(20))
	mockCache.AssertNumberOfCalls(t, "AddMessage", 20)
}

func TestGetLast20Messages_CacheError(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)

	dbMsgs := make([]model.Message, 20)

	mockCache.On("GetRecentMessages").Return([]model.Message{}, fmt.Errorf("cache error"))
	mockRepo.On("GetLastNMessages", int64(20)).Return(dbMsgs)
	mockCache.On("AddMessage", mock.Anything).Return(nil)

	mService := service.MessageService{ChatRepo: mockRepo, ChatCache: mockCache}
	result := mService.GetLast20Messages()

	assert.Equal(t, dbMsgs, result)
	mockRepo.AssertCalled(t, "GetLastNMessages", int64(20))
}
