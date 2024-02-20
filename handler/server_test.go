package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sawitpro/UserService/generated"
	mock_service "github.com/sawitpro/UserService/mocks"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockServiceInterface(ctrl)

	opts := ServerOpts{
		Service: mockService,
	}

	server := NewServer(opts)

	// Assert that the server implements the generated.ServerInterface
	_, ok := server.(generated.ServerInterface)
	assert.True(t, ok)

	// Assert that the server's service matches the mock service
	assert.Equal(t, mockService, server.(*Server).Service)
}
