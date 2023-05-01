package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/underdogio/job-board/internal/database"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// GetOrCreateUserFromToken creates or get existing user given a token
// returns the user struct, whether the user existed already and an error
func GetOrCreateUserFromToken(ctx context.Context, token string) (*database.User, bool, error) {
	usoToken, err := database.UserSignOnTokens(
		database.UserSignOnTokenWhere.Token.EQ(token),
	).OneG(ctx)
	if err != nil {
		return nil, false, err
	}

	user, err := database.Users(
		database.UserWhere.Email.EQ(usoToken.Email),
	).OneG(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	} else if errors.Is(err, sql.ErrNoRows) {
		// user not found create new one
		userID, err := ksuid.NewRandom()
		if err != nil {
			return nil, false, err
		}
		user = &database.User{
			ID:        userID.String(),
			Email:     usoToken.Email,
			CreatedAt: null.TimeFrom(time.Now()),
			UserType:  usoToken.UserType,
		}
		if err := user.InsertG(ctx, boil.Infer()); err != nil {
			return nil, false, err
		}
	}
	return user, true, nil
}

func GetUser(ctx context.Context, email string) (*database.User, error) {
	return database.Users(
		database.UserWhere.Email.EQ(email),
	).OneG(ctx)
}

func SaveTokenSignOn(ctx context.Context, email, token, userType string) error {
	usoToken := &database.UserSignOnToken{
		Token:    token,
		Email:    email,
		UserType: userType,
	}
	return usoToken.InsertG(ctx, boil.Infer())
}
