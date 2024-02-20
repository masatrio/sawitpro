package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sawitpro/UserService/common"
	commonErr "github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/helper"
	"github.com/sawitpro/UserService/mocks"
	"github.com/sawitpro/UserService/repository"
	"github.com/sawitpro/UserService/service"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name          string
		phoneNumber   string
		password      string
		user          *repository.User
		expectedError common.Error
		expectedToken string
	}{
		{
			name:          "Successful Login",
			phoneNumber:   "+628232482440",
			password:      "@Maulana",
			user:          &repository.User{ID: 1},
			expectedError: nil,
			expectedToken: "mocked_token",
		},
		{
			name:          "Error DB",
			phoneNumber:   "+628232482440",
			password:      "@Maulana",
			expectedError: commonErr.NewError("some error", commonErr.SystemErrorType),
		},
		{
			name:          "Wrong Phone or Password",
			phoneNumber:   "+628232482440",
			password:      "@Maulana",
			expectedError: commonErr.NewError(errors.New(commonErr.WrongPhonePasswordErrorMessage).Error(), commonErr.BadRequestErrorType),
		},
		{
			name:          "Wrong Password",
			phoneNumber:   "+628232482440",
			password:      "@Maulana",
			user:          &repository.User{ID: 1},
			expectedError: commonErr.NewError("password or phone number is incorrect.", commonErr.BadRequestErrorType),
		},
		{
			name:          "Insert User Activity Log Error",
			phoneNumber:   "+628232482440",
			password:      "@Maulana",
			user:          &repository.User{ID: 1},
			expectedError: commonErr.NewError("Insert User Activity Log Error", commonErr.SystemErrorType),
		},
		{
			name:          "Increment Login Count Error",
			phoneNumber:   "+628232482440",
			password:      "@Maulana",
			user:          &repository.User{ID: 1},
			expectedError: commonErr.NewError("Increment Login Count Error", commonErr.SystemErrorType),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)

			switch tc.name {
			case "Successful Login":
				hashedPassword, _ := helper.HashPassword(tc.password)
				tc.user.HashedPassword = hashedPassword
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(tc.user, nil)
				mockRepo.EXPECT().ExecTransaction(gomock.Any(), gomock.Any()).Return(nil)
			case "Error DB":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			case "Wrong Phone or Password":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phoneNumber).Return(nil, nil)
			case "Wrong Password":
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phoneNumber).Return(tc.user, nil)
			case "Insert User Activity Log Error":
				hashedPassword, _ := helper.HashPassword(tc.password)
				tc.user.HashedPassword = hashedPassword
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phoneNumber).Return(tc.user, nil)
				mockRepo.EXPECT().ExecTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					if err := fn(ctx); err != nil {
						return err
					}

					if err := mockRepo.InsertUserActivityLog(ctx, tc.user.ID, common.LOGIN_ACTIVITY); err != nil {
						return errors.New("InsertUserActivityLog error")
					}

					return nil
				})
				mockRepo.EXPECT().InsertUserActivityLog(gomock.Any(), tc.user.ID, common.LOGIN_ACTIVITY).Return(errors.New("Insert User Activity Log Error"))
			case "Increment Login Count Error":
				hashedPassword, _ := helper.HashPassword(tc.password)
				tc.user.HashedPassword = hashedPassword
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phoneNumber).Return(tc.user, nil)
				mockRepo.EXPECT().ExecTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					if err := fn(ctx); err != nil {
						return err
					}
					if err := mockRepo.InsertUserActivityLog(ctx, tc.user.ID, common.LOGIN_ACTIVITY); err != nil {
						return errors.New("InsertUserActivityLog error")
					}
					if err := mockRepo.IncrementLoginCount(ctx, tc.user.ID); err != nil {
						return errors.New("IncrementLoginCount error")
					}
					return nil
				})
				mockRepo.EXPECT().InsertUserActivityLog(gomock.Any(), tc.user.ID, common.LOGIN_ACTIVITY).Return(nil)
				mockRepo.EXPECT().IncrementLoginCount(gomock.Any(), tc.user.ID).Return(errors.New("Increment Login Count Error"))
			case "Token Generation Error":
				hashedPassword, _ := helper.HashPassword(tc.password)
				tc.user.HashedPassword = hashedPassword
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.phoneNumber).Return(tc.user, nil)
				mockRepo.EXPECT().ExecTransaction(gomock.Any(), gomock.Any()).Return(nil)
			}

			svc := NewService(ServiceOpts{
				Repository: mockRepo,
			})

			response, err := svc.Login(context.Background(), service.LoginParam{
				PhoneNumber: tc.phoneNumber,
				Password:    tc.password,
			})

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.GetErrorMessage(), err.GetErrorMessage())
				assert.Equal(t, tc.expectedError.GetErrorType(), err.GetErrorType())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, response)
			}
		})
	}
}
