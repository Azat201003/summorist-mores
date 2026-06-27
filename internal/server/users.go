package server

import (
	"context"

	"github.com/Azat201003/summorist-shared/gen/go/users"
)

func authorize(usersClient users.UsersClient, jwtToken string) (uint64, error) {
	response, err := usersClient.Authorize(context.Background(), &users.AuthRequest{JwtToken: jwtToken})
	if err != nil {
		return 0, err
	}
	if response.Code != 0 {
		return 0, ErrNotAuthorized
	}
	return response.UserId, nil
}
