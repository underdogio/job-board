package repositories

import (
	"context"
	"time"

	"github.com/underdogio/job-board/internal/database"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func GetBlogPostByIDAndAuthor(ctx context.Context, id, authorID string) (*database.BlogPost, error) {
	return database.BlogPosts(
		database.BlogPostWhere.ID.EQ(id),
		database.BlogPostWhere.CreatedBy.EQ(authorID),
	).OneG(ctx)
}

func GetBlogPostBySlug(ctx context.Context, slug string) (*database.BlogPost, error) {
	return database.BlogPosts(
		database.BlogPostWhere.Slug.EQ(slug),
		database.BlogPostWhere.PublishedAt.IsNotNull(),
	).OneG(ctx)
}

func GetBlogPostByCreatedBy(ctx context.Context, userID string) ([]*database.BlogPost, error) {
	return database.BlogPosts(
		database.BlogPostWhere.CreatedBy.EQ(userID),
	).AllG(ctx)
}

func GetAllPublishedBlogPost(ctx context.Context) ([]*database.BlogPost, error) {
	return database.BlogPosts(
		database.BlogPostWhere.PublishedAt.IsNotNull(),
	).AllG(ctx)
}

func CreateBlogPost(ctx context.Context, bp *database.BlogPost) error {
	return bp.InsertG(ctx, boil.Infer())
}

func UpdateBlogPost(ctx context.Context, bp *database.BlogPost) error {
	_, err := bp.UpdateG(ctx, boil.Infer())
	return err
}

func PublishBlogPost(ctx context.Context, bp *database.BlogPost) error {
	bp.PublishedAt = null.TimeFrom(time.Now())
	return UpdateBlogPost(ctx, bp)
}

func UnpublishBlogPost(ctx context.Context, bp *database.BlogPost) error {
	bp.PublishedAt = null.TimeFromPtr(nil)
	return UpdateBlogPost(ctx, bp)
}
