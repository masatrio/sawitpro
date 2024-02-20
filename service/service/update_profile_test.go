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

func TestUpdateProfile(t *testing.T) {
	testCases := []struct {
		name              string
		params            service.UpdateProfileParam
		user              *repository.User
		getUserByIDErr    error
		getUserByPhoneErr error
		updateUserErr     error
		expectedError     common.Error
		expectedResult    *service.UpdateProfileResponse
	}{
		{
			name: "Successful Update",
			params: service.UpdateProfileParam{
				UserID:      1,
				FullName:    "",
				PhoneNumber: "+62823248244",
			},
			user:           &repository.User{ID: 1, FullName: "Old Name", Phone: "+62823248244"},
			expectedError:  nil,
			expectedResult: &service.UpdateProfileResponse{FullName: "Old Name", PhoneNumber: "+62823248244"},
		},
		{
			name: "Successful Update No Params",
			params: service.UpdateProfileParam{
				UserID:      1,
				FullName:    "",
				PhoneNumber: "",
			},
			user:           &repository.User{ID: 1, FullName: "Old Name", Phone: "+62823248244"},
			expectedError:  nil,
			expectedResult: &service.UpdateProfileResponse{FullName: "Old Name", PhoneNumber: "+62823248244"},
		},
		{
			name: "Error DB - Get User By ID",
			params: service.UpdateProfileParam{
				UserID:      1,
				FullName:    "maulana aji satrio",
				PhoneNumber: "+62823248244",
			},
			getUserByIDErr: errors.New("some error"),
			expectedError:  commonErr.NewError("some error", commonErr.SystemErrorType),
		},
		{
			name: "User Not Found - Get User By ID",
			params: service.UpdateProfileParam{
				UserID:      1,
				FullName:    "maulana aji satrio",
				PhoneNumber: "+62823248244",
			},
			expectedError: commonErr.NewError(commonErr.UserDataNotFoundErrorMessage, commonErr.BadRequestErrorType),
		},
		{
			name: "Error DB - Get User By Phone",
			params: service.UpdateProfileParam{
				UserID:      1,
				FullName:    "maulana aji satrio",
				PhoneNumber: "+62823248244",
			},
			user:              &repository.User{ID: 1, FullName: "Old Name", Phone: "+62823248244"},
			getUserByPhoneErr: errors.New("some error"),
			expectedError:     commonErr.NewError("some error", commonErr.SystemErrorType),
		},
		{
			name: "Phone Already Used",
			params: service.UpdateProfileParam{
				UserID:      1,
				FullName:    "maulana aji satrio",
				PhoneNumber: "+62823248244",
			},
			user:          &repository.User{ID: 1, FullName: "Old Name", Phone: "+62823248244"},
			expectedError: commonErr.NewError(commonErr.NewPhoneAlreadyUsedErrorMessage("+62823248244"), commonErr.ConflictErrorType),
		},
		{
			name: "Error DB - Update User",
			params: service.UpdateProfileParam{
				UserID:      1,
				FullName:    "maulana aji satrio",
				PhoneNumber: "+62823248244",
			},
			user:          &repository.User{ID: 1, FullName: "Old Name", Phone: "+62823248244"},
			updateUserErr: errors.New("some error"),
			expectedError: commonErr.NewError("some error", commonErr.SystemErrorType),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)

			switch tc.name {
			case "Successful Update":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.user.ID).Return(tc.user, nil)
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.user.Phone).Return(nil, nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any(), tc.user.ID, tc.user.FullName, tc.user.Phone).Return(nil)
			case "Successful Update No Params":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.user.ID).Return(tc.user, nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any(), tc.user.ID, tc.user.FullName, tc.user.Phone).Return(nil)
			case "Error DB - Get User By ID":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.params.UserID).Return(nil, errors.New("some error"))
			case "User Not Found - Get User By ID":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.params.UserID).Return(nil, nil)
			case "Error DB - Get User By Phone":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.params.UserID).Return(tc.user, nil)
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.params.PhoneNumber).Return(nil, errors.New("some error"))
			case "Phone Already Used":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.params.UserID).Return(tc.user, nil)
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.params.PhoneNumber).Return(&repository.User{ID: 10}, nil)
			case "Error DB - Update User":
				mockRepo.EXPECT().GetUserByID(gomock.Any(), tc.params.UserID).Return(tc.user, nil)
				mockRepo.EXPECT().GetUserByPhone(gomock.Any(), tc.params.PhoneNumber).Return(nil, nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any(), tc.params.UserID, tc.params.FullName, tc.params.PhoneNumber).Return(errors.New("some error"))
			}

			svc := NewService(ServiceOpts{
				Repository: mockRepo,
			})

			response, err := svc.UpdateProfile(context.Background(), tc.params)

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
