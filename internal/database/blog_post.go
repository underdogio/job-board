// Code generated by SQLBoiler 4.14.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// BlogPost is an object representing the database table.
type BlogPost struct {
	ID          string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Title       string    `boil:"title" json:"title" toml:"title" yaml:"title"`
	Description string    `boil:"description" json:"description" toml:"description" yaml:"description"`
	Slug        string    `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	Tags        string    `boil:"tags" json:"tags" toml:"tags" yaml:"tags"`
	Text        string    `boil:"text" json:"text" toml:"text" yaml:"text"`
	CreatedAt   time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	CreatedBy   string    `boil:"created_by" json:"created_by" toml:"created_by" yaml:"created_by"`
	PublishedAt null.Time `boil:"published_at" json:"published_at,omitempty" toml:"published_at" yaml:"published_at,omitempty"`

	R *blogPostR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L blogPostL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var BlogPostColumns = struct {
	ID          string
	Title       string
	Description string
	Slug        string
	Tags        string
	Text        string
	CreatedAt   string
	UpdatedAt   string
	CreatedBy   string
	PublishedAt string
}{
	ID:          "id",
	Title:       "title",
	Description: "description",
	Slug:        "slug",
	Tags:        "tags",
	Text:        "text",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	CreatedBy:   "created_by",
	PublishedAt: "published_at",
}

var BlogPostTableColumns = struct {
	ID          string
	Title       string
	Description string
	Slug        string
	Tags        string
	Text        string
	CreatedAt   string
	UpdatedAt   string
	CreatedBy   string
	PublishedAt string
}{
	ID:          "blog_post.id",
	Title:       "blog_post.title",
	Description: "blog_post.description",
	Slug:        "blog_post.slug",
	Tags:        "blog_post.tags",
	Text:        "blog_post.text",
	CreatedAt:   "blog_post.created_at",
	UpdatedAt:   "blog_post.updated_at",
	CreatedBy:   "blog_post.created_by",
	PublishedAt: "blog_post.published_at",
}

// Generated where

var BlogPostWhere = struct {
	ID          whereHelperstring
	Title       whereHelperstring
	Description whereHelperstring
	Slug        whereHelperstring
	Tags        whereHelperstring
	Text        whereHelperstring
	CreatedAt   whereHelpertime_Time
	UpdatedAt   whereHelpertime_Time
	CreatedBy   whereHelperstring
	PublishedAt whereHelpernull_Time
}{
	ID:          whereHelperstring{field: "\"blog_post\".\"id\""},
	Title:       whereHelperstring{field: "\"blog_post\".\"title\""},
	Description: whereHelperstring{field: "\"blog_post\".\"description\""},
	Slug:        whereHelperstring{field: "\"blog_post\".\"slug\""},
	Tags:        whereHelperstring{field: "\"blog_post\".\"tags\""},
	Text:        whereHelperstring{field: "\"blog_post\".\"text\""},
	CreatedAt:   whereHelpertime_Time{field: "\"blog_post\".\"created_at\""},
	UpdatedAt:   whereHelpertime_Time{field: "\"blog_post\".\"updated_at\""},
	CreatedBy:   whereHelperstring{field: "\"blog_post\".\"created_by\""},
	PublishedAt: whereHelpernull_Time{field: "\"blog_post\".\"published_at\""},
}

// BlogPostRels is where relationship names are stored.
var BlogPostRels = struct {
}{}

// blogPostR is where relationships are stored.
type blogPostR struct {
}

// NewStruct creates a new relationship struct
func (*blogPostR) NewStruct() *blogPostR {
	return &blogPostR{}
}

// blogPostL is where Load methods for each relationship are stored.
type blogPostL struct{}

var (
	blogPostAllColumns            = []string{"id", "title", "description", "slug", "tags", "text", "created_at", "updated_at", "created_by", "published_at"}
	blogPostColumnsWithoutDefault = []string{"id", "title", "description", "slug", "tags", "text", "created_at", "updated_at", "created_by"}
	blogPostColumnsWithDefault    = []string{"published_at"}
	blogPostPrimaryKeyColumns     = []string{"id"}
	blogPostGeneratedColumns      = []string{}
)

type (
	// BlogPostSlice is an alias for a slice of pointers to BlogPost.
	// This should almost always be used instead of []BlogPost.
	BlogPostSlice []*BlogPost
	// BlogPostHook is the signature for custom BlogPost hook methods
	BlogPostHook func(context.Context, boil.ContextExecutor, *BlogPost) error

	blogPostQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	blogPostType                 = reflect.TypeOf(&BlogPost{})
	blogPostMapping              = queries.MakeStructMapping(blogPostType)
	blogPostPrimaryKeyMapping, _ = queries.BindMapping(blogPostType, blogPostMapping, blogPostPrimaryKeyColumns)
	blogPostInsertCacheMut       sync.RWMutex
	blogPostInsertCache          = make(map[string]insertCache)
	blogPostUpdateCacheMut       sync.RWMutex
	blogPostUpdateCache          = make(map[string]updateCache)
	blogPostUpsertCacheMut       sync.RWMutex
	blogPostUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var blogPostAfterSelectHooks []BlogPostHook

var blogPostBeforeInsertHooks []BlogPostHook
var blogPostAfterInsertHooks []BlogPostHook

var blogPostBeforeUpdateHooks []BlogPostHook
var blogPostAfterUpdateHooks []BlogPostHook

var blogPostBeforeDeleteHooks []BlogPostHook
var blogPostAfterDeleteHooks []BlogPostHook

var blogPostBeforeUpsertHooks []BlogPostHook
var blogPostAfterUpsertHooks []BlogPostHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *BlogPost) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *BlogPost) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *BlogPost) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *BlogPost) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *BlogPost) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *BlogPost) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *BlogPost) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *BlogPost) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *BlogPost) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range blogPostAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddBlogPostHook registers your hook function for all future operations.
func AddBlogPostHook(hookPoint boil.HookPoint, blogPostHook BlogPostHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		blogPostAfterSelectHooks = append(blogPostAfterSelectHooks, blogPostHook)
	case boil.BeforeInsertHook:
		blogPostBeforeInsertHooks = append(blogPostBeforeInsertHooks, blogPostHook)
	case boil.AfterInsertHook:
		blogPostAfterInsertHooks = append(blogPostAfterInsertHooks, blogPostHook)
	case boil.BeforeUpdateHook:
		blogPostBeforeUpdateHooks = append(blogPostBeforeUpdateHooks, blogPostHook)
	case boil.AfterUpdateHook:
		blogPostAfterUpdateHooks = append(blogPostAfterUpdateHooks, blogPostHook)
	case boil.BeforeDeleteHook:
		blogPostBeforeDeleteHooks = append(blogPostBeforeDeleteHooks, blogPostHook)
	case boil.AfterDeleteHook:
		blogPostAfterDeleteHooks = append(blogPostAfterDeleteHooks, blogPostHook)
	case boil.BeforeUpsertHook:
		blogPostBeforeUpsertHooks = append(blogPostBeforeUpsertHooks, blogPostHook)
	case boil.AfterUpsertHook:
		blogPostAfterUpsertHooks = append(blogPostAfterUpsertHooks, blogPostHook)
	}
}

// OneG returns a single blogPost record from the query using the global executor.
func (q blogPostQuery) OneG(ctx context.Context) (*BlogPost, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single blogPost record from the query.
func (q blogPostQuery) One(ctx context.Context, exec boil.ContextExecutor) (*BlogPost, error) {
	o := &BlogPost{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for blog_post")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all BlogPost records from the query using the global executor.
func (q blogPostQuery) AllG(ctx context.Context) (BlogPostSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all BlogPost records from the query.
func (q blogPostQuery) All(ctx context.Context, exec boil.ContextExecutor) (BlogPostSlice, error) {
	var o []*BlogPost

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to BlogPost slice")
	}

	if len(blogPostAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all BlogPost records in the query using the global executor
func (q blogPostQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all BlogPost records in the query.
func (q blogPostQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count blog_post rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q blogPostQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q blogPostQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if blog_post exists")
	}

	return count > 0, nil
}

// BlogPosts retrieves all the records using an executor.
func BlogPosts(mods ...qm.QueryMod) blogPostQuery {
	mods = append(mods, qm.From("\"blog_post\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"blog_post\".*"})
	}

	return blogPostQuery{q}
}

// FindBlogPostG retrieves a single record by ID.
func FindBlogPostG(ctx context.Context, iD string, selectCols ...string) (*BlogPost, error) {
	return FindBlogPost(ctx, boil.GetContextDB(), iD, selectCols...)
}

// FindBlogPost retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBlogPost(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*BlogPost, error) {
	blogPostObj := &BlogPost{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"blog_post\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, blogPostObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from blog_post")
	}

	if err = blogPostObj.doAfterSelectHooks(ctx, exec); err != nil {
		return blogPostObj, err
	}

	return blogPostObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *BlogPost) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *BlogPost) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no blog_post provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(blogPostColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	blogPostInsertCacheMut.RLock()
	cache, cached := blogPostInsertCache[key]
	blogPostInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			blogPostAllColumns,
			blogPostColumnsWithDefault,
			blogPostColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(blogPostType, blogPostMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(blogPostType, blogPostMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"blog_post\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"blog_post\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "database: unable to insert into blog_post")
	}

	if !cached {
		blogPostInsertCacheMut.Lock()
		blogPostInsertCache[key] = cache
		blogPostInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single BlogPost record using the global executor.
// See Update for more documentation.
func (o *BlogPost) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the BlogPost.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *BlogPost) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	blogPostUpdateCacheMut.RLock()
	cache, cached := blogPostUpdateCache[key]
	blogPostUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			blogPostAllColumns,
			blogPostPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update blog_post, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"blog_post\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, blogPostPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(blogPostType, blogPostMapping, append(wl, blogPostPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update blog_post row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for blog_post")
	}

	if !cached {
		blogPostUpdateCacheMut.Lock()
		blogPostUpdateCache[key] = cache
		blogPostUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q blogPostQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q blogPostQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for blog_post")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for blog_post")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o BlogPostSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BlogPostSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("database: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blogPostPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"blog_post\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, blogPostPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in blogPost slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all blogPost")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *BlogPost) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *BlogPost) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no blog_post provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(blogPostColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	blogPostUpsertCacheMut.RLock()
	cache, cached := blogPostUpsertCache[key]
	blogPostUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			blogPostAllColumns,
			blogPostColumnsWithDefault,
			blogPostColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			blogPostAllColumns,
			blogPostPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert blog_post, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(blogPostPrimaryKeyColumns))
			copy(conflict, blogPostPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"blog_post\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(blogPostType, blogPostMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(blogPostType, blogPostMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "database: unable to upsert blog_post")
	}

	if !cached {
		blogPostUpsertCacheMut.Lock()
		blogPostUpsertCache[key] = cache
		blogPostUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single BlogPost record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *BlogPost) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single BlogPost record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *BlogPost) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no BlogPost provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), blogPostPrimaryKeyMapping)
	sql := "DELETE FROM \"blog_post\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from blog_post")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for blog_post")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q blogPostQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q blogPostQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no blogPostQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from blog_post")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for blog_post")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o BlogPostSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BlogPostSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(blogPostBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blogPostPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"blog_post\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, blogPostPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from blogPost slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for blog_post")
	}

	if len(blogPostAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *BlogPost) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no BlogPost provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *BlogPost) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindBlogPost(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BlogPostSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty BlogPostSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BlogPostSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := BlogPostSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blogPostPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"blog_post\".* FROM \"blog_post\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, blogPostPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in BlogPostSlice")
	}

	*o = slice

	return nil
}

// BlogPostExistsG checks if the BlogPost row exists.
func BlogPostExistsG(ctx context.Context, iD string) (bool, error) {
	return BlogPostExists(ctx, boil.GetContextDB(), iD)
}

// BlogPostExists checks if the BlogPost row exists.
func BlogPostExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"blog_post\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if blog_post exists")
	}

	return exists, nil
}

// Exists checks if the BlogPost row exists.
func (o *BlogPost) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return BlogPostExists(ctx, exec, o.ID)
}
