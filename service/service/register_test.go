package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sawitpro/UserService/common"
	commonErr "github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/mocks"
	"github.com/sawitpro/UserService/repository"
	"github.com/sawitpro/UserService/service"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		name           string
		phone          string
		expectedError  common.Error
		expectedUserID int64
	}{
		{
			name:           "Successful Registration",
			phone:          "+628232482440",
			expectedError:  nil,
			expectedUserID: 1,
		},
		{
			name:          "Error DB",
			phone:         "+628232482440",
			expectedError: commonErr.NewError("some error", commonErr.SystemErrorType),
		},
		{
			name:          "Phone Already Used",
			phone:         "+628232482440",
			expectedError: commonErr.NewError(commonErr.NewPhoneAlreadyUsedErrorMessage("+628232482440"), commonErr.BadRequestErrorType),
		},
		{
			name:          "Wrong Password",
			phone:         "+628232482440",
			expectedError: commonErr.NewError("some error", commonErr.SystemErrorType),
		},
		{
			name:          "Insert User Error",
			phone:         "+628232482440",
			expectedError: commonErr.NewError("insert user error", commonErr.SystemErrorType),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			mockHasher := mocks.NewMockPasswordHasher(ctrl)

			switch tc.name {
			case "Successful Registration":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phone).Return(nil, nil)
				mockHasher.EXPECT().HashPassword(gomock.Any()).Return("any hash password", nil)
				mockRepo.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&repository.User{ID: tc.expectedUserID}, nil).Times(1)
			case "Error DB":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phone).Return(nil, errors.New("some error"))
			case "Phone Already Used":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phone).Return(&repository.User{}, nil)
			case "Wrong Password":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phone).Return(nil, nil)
				mockHasher.EXPECT().HashPassword(gomock.Any()).Return("", errors.New("some error"))
			case "Insert User Error":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phone).Return(nil, nil)
				mockHasher.EXPECT().HashPassword(gomock.Any()).Return("any hash password", nil)
				mockRepo.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("insert user error")).Times(1)
			}

			svc := NewService(ServiceOpts{
				Repository: mockRepo,
				Hasher:     mockHasher,
			})

			response, err := svc.Register(context.Background(), service.RegisterParam{
				FullName:    "maulana aji satrio",
				PhoneNumber: tc.phone,
				Password:    "@Maulana",
			})

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.GetErrorMessage(), err.GetErrorMessage())
				assert.Equal(t, tc.expectedError.GetErrorType(), err.GetErrorType())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tc.expectedUserID, response.UserID)
			}
		})
	}
}
