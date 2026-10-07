package users_transport_http

import "github.com/apollobing/todo-app/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id"           example:"10"`
	Version     int     `json:"version"      example:"3"`
	FullName    string  `json:"full_name"    example:"John Doe"`
	PhoneNumber *string `json:"phone_number" example:"+491234567890"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
