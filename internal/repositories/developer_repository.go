package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/underdogio/job-board/internal/database"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	developerProfileEventPageView    = "developer_profile_page_view"
	developerProfileEventMessageSent = "developer_profile_message_sent"
	SearchTypeDeveloper              = "developer"
)

func DeveloperProfileBySlug(ctx context.Context, slug string) (*database.DeveloperProfile, error) {
	return database.DeveloperProfiles(
		database.DeveloperProfileWhere.Slug.EQ(slug),
	).OneG(ctx)
}

func DeveloperProfileByEmail(ctx context.Context, email string) (*database.DeveloperProfile, error) {
	return database.DeveloperProfiles(
		database.DeveloperProfileWhere.Email.EQ(email),
	).OneG(ctx)
}

func DeveloperProfileByID(ctx context.Context, id string) (*database.DeveloperProfile, error) {
	return database.FindDeveloperProfileG(ctx, id)
}

func SendMessageDeveloperProfile(ctx context.Context, message *database.DeveloperProfileMessage) error {
	return message.InsertG(ctx, boil.Infer())
}

func MessageForDeliveryByID(ctx context.Context, id string) (*database.DeveloperProfileMessage, error) {
	return database.DeveloperProfileMessages(
		qm.Load(database.DeveloperProfileMessageRels.Profile),
		database.DeveloperProfileMessageWhere.ID.EQ(id),
		database.DeveloperProfileMessageWhere.SentAt.IsNull(),
	).OneG(ctx)
}

func MarkDeveloperMessageAsSent(ctx context.Context, id string) error {
	_, err := database.DeveloperProfileMessages(
		database.DeveloperProfileMessageWhere.ID.EQ(id),
	).UpdateAllG(ctx, database.M{
		database.DeveloperProfileMessageColumns.SentAt: time.Now().UTC(),
	})
	return err
}

func DevelopersByLocationAndTag(ctx context.Context, loc, tag string, pageID, pageSize int) ([]*database.DeveloperProfile, int, error) {
	offset := pageID*pageSize - pageSize
	conditions := []qm.QueryMod{}
	if tag != "" {
		conditions = append(conditions, qm.Where(database.DeveloperProfileColumns.Skills+" ILIKE ?", "%"+tag+"%"))
	}
	if loc != "" {
		conditions = append(conditions, qm.Where(database.DeveloperProfileColumns.Location+" ILIKE ?", "%"+loc+"%"))
	}
	conditions = append(
		conditions,
		qm.Where(database.DeveloperProfileColumns.CreatedAt+" != "+database.DeveloperProfileColumns.UpdatedAt),
		qm.OrderBy(database.DeveloperProfileColumns.UpdatedAt+" DESC"),
		qm.Limit(pageSize),
		qm.Offset(offset),
	)
	query := database.DeveloperProfiles(conditions...)
	count, err := query.CountG(ctx)
	if err != nil {
		return nil, 0, err
	}
	developers, err := query.AllG(ctx)
	return developers, int(count), err

}

func UpdateDeveloperProfile(ctx context.Context, dev *database.DeveloperProfile) error {
	_, err := dev.UpdateG(ctx, boil.Infer())
	return err
}

func DeleteDeveloperProfile(ctx context.Context, id, email string) error {
	_, err := database.DeveloperProfiles(
		database.DeveloperProfileWhere.ID.EQ(id),
		database.DeveloperProfileWhere.Email.EQ(email),
	).DeleteAllG(ctx)
	return err
}

func ActivateDeveloperProfile(ctx context.Context, email string) error {
	_, err := database.DeveloperProfiles(
		database.DeveloperProfileWhere.Email.EQ(email),
	).UpdateAllG(ctx, database.M{
		database.DeveloperProfileColumns.UpdatedAt: time.Now().UTC(),
	})
	return err
}

func SaveDeveloperProfile(ctx context.Context, dev *database.DeveloperProfile) error {
	dev.Slug = slug.Make(fmt.Sprintf("%s %d", dev.Name, time.Now().UTC().Unix()))
	return dev.InsertG(ctx, boil.Infer())
}

func GetTopDevelopers(ctx context.Context, limit int) ([]*database.DeveloperProfile, error) {
	return database.DeveloperProfiles(
		qm.Where(database.DeveloperProfileColumns.UpdatedAt+" != "+database.DeveloperProfileColumns.CreatedAt),
		qm.OrderBy(database.DeveloperProfileColumns.UpdatedAt+" DESC"),
		qm.Limit(limit),
	).AllG(ctx)
}

func GetTopDeveloperSkills(ctx context.Context, limit int) ([]string, error) {
	skills := make([]string, 0, limit)
	err := database.DeveloperProfiles(
		qm.Select("trim(both from unnest(regexp_split_to_array(skills, ','))) as skill"),
		qm.Where(database.DeveloperProfileColumns.UpdatedAt+" != "+database.DeveloperProfileColumns.CreatedAt),
		qm.GroupBy("skill"),
		qm.OrderBy("count(*) desc"),
		qm.Limit(limit),
	).BindG(ctx, &skills)
	return skills, err

}

func GetDeveloperSkills(ctx context.Context) ([]string, error) {
	var skills []string
	err := database.DeveloperProfiles(
		qm.Select("distinct trim(both from unnest(regexp_split_to_array(skills, ','))) as skill"),
		qm.Where(database.DeveloperProfileColumns.UpdatedAt+" != "+database.DeveloperProfileColumns.CreatedAt),
	).BindG(ctx, &skills)
	return skills, err
}

func GetDeveloperSlugs(ctx context.Context) ([]string, error) {
	slugs := make([]string, 0)
	err := database.DeveloperProfiles(
		qm.Where(database.DeveloperProfileColumns.UpdatedAt+" != "+database.DeveloperProfileColumns.CreatedAt),
	).BindG(ctx, &slugs)
	if err != nil {
		return slugs, err
	}
	return slugs, nil
}

func GetLastDevUpdatedAt(ctx context.Context) (time.Time, error) {
	dev, err := database.DeveloperProfiles(
		qm.Where(database.DeveloperProfileColumns.UpdatedAt+" != "+database.DeveloperProfileColumns.CreatedAt),
		qm.OrderBy(database.DeveloperProfileColumns.UpdatedAt+" DESC"),
		qm.Limit(1),
	).OneG(ctx)
	if err != nil {
		return time.Time{}, err
	}

	return dev.UpdatedAt.Time, nil
}

func GetDevelopersRegisteredLastMonth(ctx context.Context) (int, error) {
	count, err := database.DeveloperProfiles(
		qm.Where(database.DeveloperProfileColumns.CreatedAt + " > NOW() - INTERVAL '30 days'"),
	).CountG(ctx)
	return int(count), err
}

func GetDeveloperMessagesSentLastMonth(ctx context.Context) (int, error) {
	count, err := database.DeveloperProfileMessages(
		qm.Where(database.DeveloperProfileMessageColumns.CreatedAt + " > NOW() - INTERVAL '30 days'"),
	).CountG(ctx)

	return int(count), err
}

func GetDeveloperProfilePageViewsLastMonth(ctx context.Context) (int, error) {
	count, err := database.DeveloperProfileEvents(
		database.DeveloperProfileEventWhere.EventType.EQ(developerProfileEventPageView),
		qm.Where(database.DeveloperProfileEventColumns.CreatedAt+" > NOW() - INTERVAL '30 days'"),
	).CountG(ctx)

	return int(count), err
}

func TrackDeveloperProfileView(ctx context.Context, dev *database.DeveloperProfile) error {
	event := &database.DeveloperProfileEvent{
		DeveloperProfileID: dev.ID,
		EventType:          developerProfileEventPageView,
		CreatedAt:          time.Now().UTC(),
	}
	return event.InsertG(ctx, boil.Infer())
}

func TrackDeveloperProfileMessageSent(ctx context.Context, dev *database.DeveloperProfile) error {
	event := &database.DeveloperProfileEvent{
		DeveloperProfileID: dev.ID,
		EventType:          developerProfileEventMessageSent,
		CreatedAt:          time.Now().UTC(),
	}
	return event.InsertG(ctx, boil.Infer())
}
