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

// ApplyToken is an object representing the database table.
type ApplyToken struct {
	Token       string    `boil:"token" json:"token" toml:"token" yaml:"token"`
	JobID       string    `boil:"job_id" json:"job_id" toml:"job_id" yaml:"job_id"`
	CreatedAt   time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	ConfirmedAt null.Time `boil:"confirmed_at" json:"confirmed_at,omitempty" toml:"confirmed_at" yaml:"confirmed_at,omitempty"`
	Email       string    `boil:"email" json:"email" toml:"email" yaml:"email"`
	CV          []byte    `boil:"cv" json:"cv" toml:"cv" yaml:"cv"`

	R *applyTokenR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L applyTokenL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ApplyTokenColumns = struct {
	Token       string
	JobID       string
	CreatedAt   string
	ConfirmedAt string
	Email       string
	CV          string
}{
	Token:       "token",
	JobID:       "job_id",
	CreatedAt:   "created_at",
	ConfirmedAt: "confirmed_at",
	Email:       "email",
	CV:          "cv",
}

var ApplyTokenTableColumns = struct {
	Token       string
	JobID       string
	CreatedAt   string
	ConfirmedAt string
	Email       string
	CV          string
}{
	Token:       "apply_token.token",
	JobID:       "apply_token.job_id",
	CreatedAt:   "apply_token.created_at",
	ConfirmedAt: "apply_token.confirmed_at",
	Email:       "apply_token.email",
	CV:          "apply_token.cv",
}

// Generated where

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelpernull_Time struct{ field string }

func (w whereHelpernull_Time) EQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Time) NEQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Time) LT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Time) LTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Time) GT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Time) GTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Time) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Time) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelper__byte struct{ field string }

func (w whereHelper__byte) EQ(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelper__byte) NEQ(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelper__byte) LT(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelper__byte) LTE(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelper__byte) GT(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelper__byte) GTE(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

var ApplyTokenWhere = struct {
	Token       whereHelperstring
	JobID       whereHelperstring
	CreatedAt   whereHelpertime_Time
	ConfirmedAt whereHelpernull_Time
	Email       whereHelperstring
	CV          whereHelper__byte
}{
	Token:       whereHelperstring{field: "\"apply_token\".\"token\""},
	JobID:       whereHelperstring{field: "\"apply_token\".\"job_id\""},
	CreatedAt:   whereHelpertime_Time{field: "\"apply_token\".\"created_at\""},
	ConfirmedAt: whereHelpernull_Time{field: "\"apply_token\".\"confirmed_at\""},
	Email:       whereHelperstring{field: "\"apply_token\".\"email\""},
	CV:          whereHelper__byte{field: "\"apply_token\".\"cv\""},
}

// ApplyTokenRels is where relationship names are stored.
var ApplyTokenRels = struct {
	Job string
}{
	Job: "Job",
}

// applyTokenR is where relationships are stored.
type applyTokenR struct {
	Job *Job `boil:"Job" json:"Job" toml:"Job" yaml:"Job"`
}

// NewStruct creates a new relationship struct
func (*applyTokenR) NewStruct() *applyTokenR {
	return &applyTokenR{}
}

func (r *applyTokenR) GetJob() *Job {
	if r == nil {
		return nil
	}
	return r.Job
}

// applyTokenL is where Load methods for each relationship are stored.
type applyTokenL struct{}

var (
	applyTokenAllColumns            = []string{"token", "job_id", "created_at", "confirmed_at", "email", "cv"}
	applyTokenColumnsWithoutDefault = []string{"token", "job_id", "created_at", "email", "cv"}
	applyTokenColumnsWithDefault    = []string{"confirmed_at"}
	applyTokenPrimaryKeyColumns     = []string{"token"}
	applyTokenGeneratedColumns      = []string{}
)

type (
	// ApplyTokenSlice is an alias for a slice of pointers to ApplyToken.
	// This should almost always be used instead of []ApplyToken.
	ApplyTokenSlice []*ApplyToken
	// ApplyTokenHook is the signature for custom ApplyToken hook methods
	ApplyTokenHook func(context.Context, boil.ContextExecutor, *ApplyToken) error

	applyTokenQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	applyTokenType                 = reflect.TypeOf(&ApplyToken{})
	applyTokenMapping              = queries.MakeStructMapping(applyTokenType)
	applyTokenPrimaryKeyMapping, _ = queries.BindMapping(applyTokenType, applyTokenMapping, applyTokenPrimaryKeyColumns)
	applyTokenInsertCacheMut       sync.RWMutex
	applyTokenInsertCache          = make(map[string]insertCache)
	applyTokenUpdateCacheMut       sync.RWMutex
	applyTokenUpdateCache          = make(map[string]updateCache)
	applyTokenUpsertCacheMut       sync.RWMutex
	applyTokenUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var applyTokenAfterSelectHooks []ApplyTokenHook

var applyTokenBeforeInsertHooks []ApplyTokenHook
var applyTokenAfterInsertHooks []ApplyTokenHook

var applyTokenBeforeUpdateHooks []ApplyTokenHook
var applyTokenAfterUpdateHooks []ApplyTokenHook

var applyTokenBeforeDeleteHooks []ApplyTokenHook
var applyTokenAfterDeleteHooks []ApplyTokenHook

var applyTokenBeforeUpsertHooks []ApplyTokenHook
var applyTokenAfterUpsertHooks []ApplyTokenHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ApplyToken) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ApplyToken) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ApplyToken) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ApplyToken) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ApplyToken) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ApplyToken) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ApplyToken) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ApplyToken) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ApplyToken) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applyTokenAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddApplyTokenHook registers your hook function for all future operations.
func AddApplyTokenHook(hookPoint boil.HookPoint, applyTokenHook ApplyTokenHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		applyTokenAfterSelectHooks = append(applyTokenAfterSelectHooks, applyTokenHook)
	case boil.BeforeInsertHook:
		applyTokenBeforeInsertHooks = append(applyTokenBeforeInsertHooks, applyTokenHook)
	case boil.AfterInsertHook:
		applyTokenAfterInsertHooks = append(applyTokenAfterInsertHooks, applyTokenHook)
	case boil.BeforeUpdateHook:
		applyTokenBeforeUpdateHooks = append(applyTokenBeforeUpdateHooks, applyTokenHook)
	case boil.AfterUpdateHook:
		applyTokenAfterUpdateHooks = append(applyTokenAfterUpdateHooks, applyTokenHook)
	case boil.BeforeDeleteHook:
		applyTokenBeforeDeleteHooks = append(applyTokenBeforeDeleteHooks, applyTokenHook)
	case boil.AfterDeleteHook:
		applyTokenAfterDeleteHooks = append(applyTokenAfterDeleteHooks, applyTokenHook)
	case boil.BeforeUpsertHook:
		applyTokenBeforeUpsertHooks = append(applyTokenBeforeUpsertHooks, applyTokenHook)
	case boil.AfterUpsertHook:
		applyTokenAfterUpsertHooks = append(applyTokenAfterUpsertHooks, applyTokenHook)
	}
}

// OneG returns a single applyToken record from the query using the global executor.
func (q applyTokenQuery) OneG(ctx context.Context) (*ApplyToken, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single applyToken record from the query.
func (q applyTokenQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ApplyToken, error) {
	o := &ApplyToken{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for apply_token")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all ApplyToken records from the query using the global executor.
func (q applyTokenQuery) AllG(ctx context.Context) (ApplyTokenSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all ApplyToken records from the query.
func (q applyTokenQuery) All(ctx context.Context, exec boil.ContextExecutor) (ApplyTokenSlice, error) {
	var o []*ApplyToken

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to ApplyToken slice")
	}

	if len(applyTokenAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all ApplyToken records in the query using the global executor
func (q applyTokenQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all ApplyToken records in the query.
func (q applyTokenQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count apply_token rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q applyTokenQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q applyTokenQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if apply_token exists")
	}

	return count > 0, nil
}

// Job pointed to by the foreign key.
func (o *ApplyToken) Job(mods ...qm.QueryMod) jobQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.JobID),
	}

	queryMods = append(queryMods, mods...)

	return Jobs(queryMods...)
}

// LoadJob allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (applyTokenL) LoadJob(ctx context.Context, e boil.ContextExecutor, singular bool, maybeApplyToken interface{}, mods queries.Applicator) error {
	var slice []*ApplyToken
	var object *ApplyToken

	if singular {
		var ok bool
		object, ok = maybeApplyToken.(*ApplyToken)
		if !ok {
			object = new(ApplyToken)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeApplyToken)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeApplyToken))
			}
		}
	} else {
		s, ok := maybeApplyToken.(*[]*ApplyToken)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeApplyToken)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeApplyToken))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &applyTokenR{}
		}
		args = append(args, object.JobID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &applyTokenR{}
			}

			for _, a := range args {
				if a == obj.JobID {
					continue Outer
				}
			}

			args = append(args, obj.JobID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`job`),
		qm.WhereIn(`job.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Job")
	}

	var resultSlice []*Job
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Job")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for job")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for job")
	}

	if len(jobAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Job = foreign
		if foreign.R == nil {
			foreign.R = &jobR{}
		}
		foreign.R.ApplyTokens = append(foreign.R.ApplyTokens, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.JobID == foreign.ID {
				local.R.Job = foreign
				if foreign.R == nil {
					foreign.R = &jobR{}
				}
				foreign.R.ApplyTokens = append(foreign.R.ApplyTokens, local)
				break
			}
		}
	}

	return nil
}

// SetJobG of the applyToken to the related item.
// Sets o.R.Job to related.
// Adds o to related.R.ApplyTokens.
// Uses the global database handle.
func (o *ApplyToken) SetJobG(ctx context.Context, insert bool, related *Job) error {
	return o.SetJob(ctx, boil.GetContextDB(), insert, related)
}

// SetJob of the applyToken to the related item.
// Sets o.R.Job to related.
// Adds o to related.R.ApplyTokens.
func (o *ApplyToken) SetJob(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Job) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"apply_token\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"job_id"}),
		strmangle.WhereClause("\"", "\"", 2, applyTokenPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.Token}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.JobID = related.ID
	if o.R == nil {
		o.R = &applyTokenR{
			Job: related,
		}
	} else {
		o.R.Job = related
	}

	if related.R == nil {
		related.R = &jobR{
			ApplyTokens: ApplyTokenSlice{o},
		}
	} else {
		related.R.ApplyTokens = append(related.R.ApplyTokens, o)
	}

	return nil
}

// ApplyTokens retrieves all the records using an executor.
func ApplyTokens(mods ...qm.QueryMod) applyTokenQuery {
	mods = append(mods, qm.From("\"apply_token\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"apply_token\".*"})
	}

	return applyTokenQuery{q}
}

// FindApplyTokenG retrieves a single record by ID.
func FindApplyTokenG(ctx context.Context, token string, selectCols ...string) (*ApplyToken, error) {
	return FindApplyToken(ctx, boil.GetContextDB(), token, selectCols...)
}

// FindApplyToken retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindApplyToken(ctx context.Context, exec boil.ContextExecutor, token string, selectCols ...string) (*ApplyToken, error) {
	applyTokenObj := &ApplyToken{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"apply_token\" where \"token\"=$1", sel,
	)

	q := queries.Raw(query, token)

	err := q.Bind(ctx, exec, applyTokenObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from apply_token")
	}

	if err = applyTokenObj.doAfterSelectHooks(ctx, exec); err != nil {
		return applyTokenObj, err
	}

	return applyTokenObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *ApplyToken) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ApplyToken) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no apply_token provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(applyTokenColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	applyTokenInsertCacheMut.RLock()
	cache, cached := applyTokenInsertCache[key]
	applyTokenInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			applyTokenAllColumns,
			applyTokenColumnsWithDefault,
			applyTokenColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(applyTokenType, applyTokenMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(applyTokenType, applyTokenMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"apply_token\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"apply_token\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "database: unable to insert into apply_token")
	}

	if !cached {
		applyTokenInsertCacheMut.Lock()
		applyTokenInsertCache[key] = cache
		applyTokenInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single ApplyToken record using the global executor.
// See Update for more documentation.
func (o *ApplyToken) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the ApplyToken.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ApplyToken) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	applyTokenUpdateCacheMut.RLock()
	cache, cached := applyTokenUpdateCache[key]
	applyTokenUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			applyTokenAllColumns,
			applyTokenPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update apply_token, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"apply_token\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, applyTokenPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(applyTokenType, applyTokenMapping, append(wl, applyTokenPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "database: unable to update apply_token row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for apply_token")
	}

	if !cached {
		applyTokenUpdateCacheMut.Lock()
		applyTokenUpdateCache[key] = cache
		applyTokenUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q applyTokenQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q applyTokenQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for apply_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for apply_token")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ApplyTokenSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ApplyTokenSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), applyTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"apply_token\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, applyTokenPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in applyToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all applyToken")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *ApplyToken) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ApplyToken) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no apply_token provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(applyTokenColumnsWithDefault, o)

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

	applyTokenUpsertCacheMut.RLock()
	cache, cached := applyTokenUpsertCache[key]
	applyTokenUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			applyTokenAllColumns,
			applyTokenColumnsWithDefault,
			applyTokenColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			applyTokenAllColumns,
			applyTokenPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert apply_token, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(applyTokenPrimaryKeyColumns))
			copy(conflict, applyTokenPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"apply_token\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(applyTokenType, applyTokenMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(applyTokenType, applyTokenMapping, ret)
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
		return errors.Wrap(err, "database: unable to upsert apply_token")
	}

	if !cached {
		applyTokenUpsertCacheMut.Lock()
		applyTokenUpsertCache[key] = cache
		applyTokenUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single ApplyToken record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *ApplyToken) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single ApplyToken record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ApplyToken) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no ApplyToken provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), applyTokenPrimaryKeyMapping)
	sql := "DELETE FROM \"apply_token\" WHERE \"token\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from apply_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for apply_token")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q applyTokenQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q applyTokenQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no applyTokenQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from apply_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for apply_token")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o ApplyTokenSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ApplyTokenSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(applyTokenBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), applyTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"apply_token\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, applyTokenPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from applyToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for apply_token")
	}

	if len(applyTokenAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *ApplyToken) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no ApplyToken provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ApplyToken) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindApplyToken(ctx, exec, o.Token)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ApplyTokenSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty ApplyTokenSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ApplyTokenSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ApplyTokenSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), applyTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"apply_token\".* FROM \"apply_token\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, applyTokenPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in ApplyTokenSlice")
	}

	*o = slice

	return nil
}

// ApplyTokenExistsG checks if the ApplyToken row exists.
func ApplyTokenExistsG(ctx context.Context, token string) (bool, error) {
	return ApplyTokenExists(ctx, boil.GetContextDB(), token)
}

// ApplyTokenExists checks if the ApplyToken row exists.
func ApplyTokenExists(ctx context.Context, exec boil.ContextExecutor, token string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"apply_token\" where \"token\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, token)
	}
	row := exec.QueryRowContext(ctx, sql, token)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if apply_token exists")
	}

	return exists, nil
}

// Exists checks if the ApplyToken row exists.
func (o *ApplyToken) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ApplyTokenExists(ctx, exec, o.Token)
}
