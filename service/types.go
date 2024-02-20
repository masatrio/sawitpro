package service

type RegisterParam struct {
	FullName    string
	PhoneNumber string
	Password    string
}

type RegisterResponse struct {
	UserID int64
}

type LoginParam struct {
	PhoneNumber string
	Password    string
}

type LoginResponse struct {
	UserID int64
	Token  string
}

type UserInfoResponse struct {
	FullName    string
	PhoneNumber string
}

type UpdateProfileParam struct {
	FullName    string
	PhoneNumber string
	UserID      int64
}

type UpdateProfileResponse struct {
	FullName    string
	PhoneNumber string
}
