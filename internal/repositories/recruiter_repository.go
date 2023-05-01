package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/underdogio/job-board/internal/database"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func RecruiterProfileByID(ctx context.Context, id string) (*database.RecruiterProfile, error) {
	return database.FindRecruiterProfileG(ctx, id)
}

func ActivateRecruiterProfile(ctx context.Context, email string) error {
	_, err := database.RecruiterProfiles(
		database.RecruiterProfileWhere.Email.EQ(email),
	).UpdateAllG(ctx, database.M{
		database.RecruiterProfileColumns.UpdatedAt: time.Now().UTC(),
	})
	return err
}

func RecruiterProfileByEmail(ctx context.Context, email string) (*database.RecruiterProfile, error) {
	return database.RecruiterProfiles(
		database.RecruiterProfileWhere.Email.EQ(email),
	).OneG(ctx)
}

func SaveRecruiterProfile(ctx context.Context, rec *database.RecruiterProfile) error {
	rec.Slug = slug.Make(fmt.Sprintf("%s %d", rec.Name.String, time.Now().UTC().Unix()))
	return rec.InsertG(ctx, boil.Infer())
}
