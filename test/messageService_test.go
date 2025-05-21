package test

import (
	"fmt"
	"generalChat/entity"
	"generalChat/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveMessage(message entity.Message) {
	m.Called(message)
}

func (m *MockRepo) GetLastNMessages(n int64) []entity.Message {
	args := m.Called(n)
	return args.Get(0).([]entity.Message)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) AddMessage(msg entity.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockCache) GetRecentMessages() ([]entity.Message, error) {
	args := m.Called()
	return args.Get(0).([]entity.Message), args.Error(1)
}

func TestGetLast20Messages_FromCache(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)

	cachedMsgs := make([]entity.Message, 20)
	mockCache.On("GetRecentMessages").Return(cachedMsgs, nil)

	service := services.MessageService{ChatRepo: mockRepo, ChatCache: mockCache}
	result := service.GetLast20Messages()

	assert.Equal(t, cachedMsgs, result)
	mockRepo.AssertNotCalled(t, "GetLastNMessages", mock.Anything)
}

func TestGetLast20Messages_FromRepo(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)

	cachedMsgs := make([]entity.Message, 5)
	dbMsgs := make([]entity.Message, 20)

	mockCache.On("GetRecentMessages").Return(cachedMsgs, nil)
	mockRepo.On("GetLastNMessages", int64(20)).Return(dbMsgs)
	mockCache.On("AddMessage", mock.Anything).Return(nil)

	service := services.MessageService{ChatRepo: mockRepo, ChatCache: mockCache}
	result := service.GetLast20Messages()

	assert.Equal(t, dbMsgs, result)
	mockRepo.AssertCalled(t, "GetLastNMessages", int64(20))
	mockCache.AssertNumberOfCalls(t, "AddMessage", 20)
}

func TestGetLast20Messages_CacheError(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)

	dbMsgs := make([]entity.Message, 20)

	mockCache.On("GetRecentMessages").Return([]entity.Message{}, fmt.Errorf("cache error"))
	mockRepo.On("GetLastNMessages", int64(20)).Return(dbMsgs)
	mockCache.On("AddMessage", mock.Anything).Return(nil)

	service := services.MessageService{ChatRepo: mockRepo, ChatCache: mockCache}
	result := service.GetLast20Messages()

	assert.Equal(t, dbMsgs, result)
	mockRepo.AssertCalled(t, "GetLastNMessages", int64(20))
}
