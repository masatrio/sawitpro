package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sawitpro/UserService/common"
	commonErr "github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/mocks"
	"github.com/sawitpro/UserService/repository"
	"github.com/sawitpro/UserService/service"
)

func TestGetProfile(t *testing.T) {
	testCases := []struct {
		name           string
		userID         int64
		user           *repository.User
		expectedError  common.Error
		expectedResult *service.UserInfoResponse
	}{
		{
			name:           "Successful GetProfile",
			userID:         1,
			user:           &repository.User{ID: 1, FullName: "maulana aji", Phone: "+628232482488"},
			expectedError:  nil,
			expectedResult: &service.UserInfoResponse{FullName: "maulana aji", PhoneNumber: "+628232482488"},
		},
		{
			name:           "Error DB",
			userID:         1,
			expectedError:  commonErr.NewError("some error", commonErr.SystemErrorType),
			expectedResult: nil,
		},
		{
			name:           "User Not Found",
			userID:         1,
			expectedError:  commonErr.NewError(errors.New(commonErr.UserDataNotFoundErrorMessage).Error(), commonErr.BadRequestErrorType),
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)

			switch tc.name {
			case "Successful GetProfile":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.userID).Return(tc.user, nil)
			case "Error DB":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.userID).Return(nil, errors.New("some error"))
			case "User Not Found":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.userID).Return(nil, nil)
			}

			svc := NewService(ServiceOpts{
				Repository: mockRepo,
			})

			response, err := svc.GetProfile(context.Background(), tc.userID)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.GetErrorMessage(), err.GetErrorMessage())
				assert.Equal(t, tc.expectedError.GetErrorType(), err.GetErrorType())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tc.expectedResult.FullName, response.FullName)
				assert.Equal(t, tc.expectedResult.PhoneNumber, response.PhoneNumber)
			}
		})
	}
}
