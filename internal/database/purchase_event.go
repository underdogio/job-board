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

// PurchaseEvent is an object representing the database table.
type PurchaseEvent struct {
	StripeSessionID string    `boil:"stripe_session_id" json:"stripe_session_id" toml:"stripe_session_id" yaml:"stripe_session_id"`
	Email           string    `boil:"email" json:"email" toml:"email" yaml:"email"`
	PlanID          string    `boil:"plan_id" json:"plan_id" toml:"plan_id" yaml:"plan_id"`
	Description     string    `boil:"description" json:"description" toml:"description" yaml:"description"`
	JobID           string    `boil:"job_id" json:"job_id" toml:"job_id" yaml:"job_id"`
	CreatedAt       time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	CompletedAt     null.Time `boil:"completed_at" json:"completed_at,omitempty" toml:"completed_at" yaml:"completed_at,omitempty"`

	R *purchaseEventR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L purchaseEventL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PurchaseEventColumns = struct {
	StripeSessionID string
	Email           string
	PlanID          string
	Description     string
	JobID           string
	CreatedAt       string
	CompletedAt     string
}{
	StripeSessionID: "stripe_session_id",
	Email:           "email",
	PlanID:          "plan_id",
	Description:     "description",
	JobID:           "job_id",
	CreatedAt:       "created_at",
	CompletedAt:     "completed_at",
}

var PurchaseEventTableColumns = struct {
	StripeSessionID string
	Email           string
	PlanID          string
	Description     string
	JobID           string
	CreatedAt       string
	CompletedAt     string
}{
	StripeSessionID: "purchase_event.stripe_session_id",
	Email:           "purchase_event.email",
	PlanID:          "purchase_event.plan_id",
	Description:     "purchase_event.description",
	JobID:           "purchase_event.job_id",
	CreatedAt:       "purchase_event.created_at",
	CompletedAt:     "purchase_event.completed_at",
}

// Generated where

var PurchaseEventWhere = struct {
	StripeSessionID whereHelperstring
	Email           whereHelperstring
	PlanID          whereHelperstring
	Description     whereHelperstring
	JobID           whereHelperstring
	CreatedAt       whereHelpertime_Time
	CompletedAt     whereHelpernull_Time
}{
	StripeSessionID: whereHelperstring{field: "\"purchase_event\".\"stripe_session_id\""},
	Email:           whereHelperstring{field: "\"purchase_event\".\"email\""},
	PlanID:          whereHelperstring{field: "\"purchase_event\".\"plan_id\""},
	Description:     whereHelperstring{field: "\"purchase_event\".\"description\""},
	JobID:           whereHelperstring{field: "\"purchase_event\".\"job_id\""},
	CreatedAt:       whereHelpertime_Time{field: "\"purchase_event\".\"created_at\""},
	CompletedAt:     whereHelpernull_Time{field: "\"purchase_event\".\"completed_at\""},
}

// PurchaseEventRels is where relationship names are stored.
var PurchaseEventRels = struct {
	Job string
}{
	Job: "Job",
}

// purchaseEventR is where relationships are stored.
type purchaseEventR struct {
	Job *Job `boil:"Job" json:"Job" toml:"Job" yaml:"Job"`
}

// NewStruct creates a new relationship struct
func (*purchaseEventR) NewStruct() *purchaseEventR {
	return &purchaseEventR{}
}

func (r *purchaseEventR) GetJob() *Job {
	if r == nil {
		return nil
	}
	return r.Job
}

// purchaseEventL is where Load methods for each relationship are stored.
type purchaseEventL struct{}

var (
	purchaseEventAllColumns            = []string{"stripe_session_id", "email", "plan_id", "description", "job_id", "created_at", "completed_at"}
	purchaseEventColumnsWithoutDefault = []string{"stripe_session_id", "description", "job_id", "created_at"}
	purchaseEventColumnsWithDefault    = []string{"email", "plan_id", "completed_at"}
	purchaseEventPrimaryKeyColumns     = []string{"stripe_session_id"}
	purchaseEventGeneratedColumns      = []string{}
)

type (
	// PurchaseEventSlice is an alias for a slice of pointers to PurchaseEvent.
	// This should almost always be used instead of []PurchaseEvent.
	PurchaseEventSlice []*PurchaseEvent
	// PurchaseEventHook is the signature for custom PurchaseEvent hook methods
	PurchaseEventHook func(context.Context, boil.ContextExecutor, *PurchaseEvent) error

	purchaseEventQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	purchaseEventType                 = reflect.TypeOf(&PurchaseEvent{})
	purchaseEventMapping              = queries.MakeStructMapping(purchaseEventType)
	purchaseEventPrimaryKeyMapping, _ = queries.BindMapping(purchaseEventType, purchaseEventMapping, purchaseEventPrimaryKeyColumns)
	purchaseEventInsertCacheMut       sync.RWMutex
	purchaseEventInsertCache          = make(map[string]insertCache)
	purchaseEventUpdateCacheMut       sync.RWMutex
	purchaseEventUpdateCache          = make(map[string]updateCache)
	purchaseEventUpsertCacheMut       sync.RWMutex
	purchaseEventUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var purchaseEventAfterSelectHooks []PurchaseEventHook

var purchaseEventBeforeInsertHooks []PurchaseEventHook
var purchaseEventAfterInsertHooks []PurchaseEventHook

var purchaseEventBeforeUpdateHooks []PurchaseEventHook
var purchaseEventAfterUpdateHooks []PurchaseEventHook

var purchaseEventBeforeDeleteHooks []PurchaseEventHook
var purchaseEventAfterDeleteHooks []PurchaseEventHook

var purchaseEventBeforeUpsertHooks []PurchaseEventHook
var purchaseEventAfterUpsertHooks []PurchaseEventHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *PurchaseEvent) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *PurchaseEvent) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *PurchaseEvent) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *PurchaseEvent) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *PurchaseEvent) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *PurchaseEvent) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *PurchaseEvent) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *PurchaseEvent) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *PurchaseEvent) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range purchaseEventAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPurchaseEventHook registers your hook function for all future operations.
func AddPurchaseEventHook(hookPoint boil.HookPoint, purchaseEventHook PurchaseEventHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		purchaseEventAfterSelectHooks = append(purchaseEventAfterSelectHooks, purchaseEventHook)
	case boil.BeforeInsertHook:
		purchaseEventBeforeInsertHooks = append(purchaseEventBeforeInsertHooks, purchaseEventHook)
	case boil.AfterInsertHook:
		purchaseEventAfterInsertHooks = append(purchaseEventAfterInsertHooks, purchaseEventHook)
	case boil.BeforeUpdateHook:
		purchaseEventBeforeUpdateHooks = append(purchaseEventBeforeUpdateHooks, purchaseEventHook)
	case boil.AfterUpdateHook:
		purchaseEventAfterUpdateHooks = append(purchaseEventAfterUpdateHooks, purchaseEventHook)
	case boil.BeforeDeleteHook:
		purchaseEventBeforeDeleteHooks = append(purchaseEventBeforeDeleteHooks, purchaseEventHook)
	case boil.AfterDeleteHook:
		purchaseEventAfterDeleteHooks = append(purchaseEventAfterDeleteHooks, purchaseEventHook)
	case boil.BeforeUpsertHook:
		purchaseEventBeforeUpsertHooks = append(purchaseEventBeforeUpsertHooks, purchaseEventHook)
	case boil.AfterUpsertHook:
		purchaseEventAfterUpsertHooks = append(purchaseEventAfterUpsertHooks, purchaseEventHook)
	}
}

// OneG returns a single purchaseEvent record from the query using the global executor.
func (q purchaseEventQuery) OneG(ctx context.Context) (*PurchaseEvent, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single purchaseEvent record from the query.
func (q purchaseEventQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PurchaseEvent, error) {
	o := &PurchaseEvent{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for purchase_event")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all PurchaseEvent records from the query using the global executor.
func (q purchaseEventQuery) AllG(ctx context.Context) (PurchaseEventSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all PurchaseEvent records from the query.
func (q purchaseEventQuery) All(ctx context.Context, exec boil.ContextExecutor) (PurchaseEventSlice, error) {
	var o []*PurchaseEvent

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to PurchaseEvent slice")
	}

	if len(purchaseEventAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all PurchaseEvent records in the query using the global executor
func (q purchaseEventQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all PurchaseEvent records in the query.
func (q purchaseEventQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count purchase_event rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q purchaseEventQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q purchaseEventQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if purchase_event exists")
	}

	return count > 0, nil
}

// Job pointed to by the foreign key.
func (o *PurchaseEvent) Job(mods ...qm.QueryMod) jobQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.JobID),
	}

	queryMods = append(queryMods, mods...)

	return Jobs(queryMods...)
}

// LoadJob allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (purchaseEventL) LoadJob(ctx context.Context, e boil.ContextExecutor, singular bool, maybePurchaseEvent interface{}, mods queries.Applicator) error {
	var slice []*PurchaseEvent
	var object *PurchaseEvent

	if singular {
		var ok bool
		object, ok = maybePurchaseEvent.(*PurchaseEvent)
		if !ok {
			object = new(PurchaseEvent)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePurchaseEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePurchaseEvent))
			}
		}
	} else {
		s, ok := maybePurchaseEvent.(*[]*PurchaseEvent)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePurchaseEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePurchaseEvent))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &purchaseEventR{}
		}
		args = append(args, object.JobID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &purchaseEventR{}
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
		foreign.R.PurchaseEvents = append(foreign.R.PurchaseEvents, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.JobID == foreign.ID {
				local.R.Job = foreign
				if foreign.R == nil {
					foreign.R = &jobR{}
				}
				foreign.R.PurchaseEvents = append(foreign.R.PurchaseEvents, local)
				break
			}
		}
	}

	return nil
}

// SetJobG of the purchaseEvent to the related item.
// Sets o.R.Job to related.
// Adds o to related.R.PurchaseEvents.
// Uses the global database handle.
func (o *PurchaseEvent) SetJobG(ctx context.Context, insert bool, related *Job) error {
	return o.SetJob(ctx, boil.GetContextDB(), insert, related)
}

// SetJob of the purchaseEvent to the related item.
// Sets o.R.Job to related.
// Adds o to related.R.PurchaseEvents.
func (o *PurchaseEvent) SetJob(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Job) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"purchase_event\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"job_id"}),
		strmangle.WhereClause("\"", "\"", 2, purchaseEventPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.StripeSessionID}

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
		o.R = &purchaseEventR{
			Job: related,
		}
	} else {
		o.R.Job = related
	}

	if related.R == nil {
		related.R = &jobR{
			PurchaseEvents: PurchaseEventSlice{o},
		}
	} else {
		related.R.PurchaseEvents = append(related.R.PurchaseEvents, o)
	}

	return nil
}

// PurchaseEvents retrieves all the records using an executor.
func PurchaseEvents(mods ...qm.QueryMod) purchaseEventQuery {
	mods = append(mods, qm.From("\"purchase_event\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"purchase_event\".*"})
	}

	return purchaseEventQuery{q}
}

// FindPurchaseEventG retrieves a single record by ID.
func FindPurchaseEventG(ctx context.Context, stripeSessionID string, selectCols ...string) (*PurchaseEvent, error) {
	return FindPurchaseEvent(ctx, boil.GetContextDB(), stripeSessionID, selectCols...)
}

// FindPurchaseEvent retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPurchaseEvent(ctx context.Context, exec boil.ContextExecutor, stripeSessionID string, selectCols ...string) (*PurchaseEvent, error) {
	purchaseEventObj := &PurchaseEvent{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"purchase_event\" where \"stripe_session_id\"=$1", sel,
	)

	q := queries.Raw(query, stripeSessionID)

	err := q.Bind(ctx, exec, purchaseEventObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from purchase_event")
	}

	if err = purchaseEventObj.doAfterSelectHooks(ctx, exec); err != nil {
		return purchaseEventObj, err
	}

	return purchaseEventObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *PurchaseEvent) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PurchaseEvent) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no purchase_event provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(purchaseEventColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	purchaseEventInsertCacheMut.RLock()
	cache, cached := purchaseEventInsertCache[key]
	purchaseEventInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			purchaseEventAllColumns,
			purchaseEventColumnsWithDefault,
			purchaseEventColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(purchaseEventType, purchaseEventMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(purchaseEventType, purchaseEventMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"purchase_event\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"purchase_event\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "database: unable to insert into purchase_event")
	}

	if !cached {
		purchaseEventInsertCacheMut.Lock()
		purchaseEventInsertCache[key] = cache
		purchaseEventInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single PurchaseEvent record using the global executor.
// See Update for more documentation.
func (o *PurchaseEvent) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the PurchaseEvent.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PurchaseEvent) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	purchaseEventUpdateCacheMut.RLock()
	cache, cached := purchaseEventUpdateCache[key]
	purchaseEventUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			purchaseEventAllColumns,
			purchaseEventPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update purchase_event, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"purchase_event\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, purchaseEventPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(purchaseEventType, purchaseEventMapping, append(wl, purchaseEventPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "database: unable to update purchase_event row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for purchase_event")
	}

	if !cached {
		purchaseEventUpdateCacheMut.Lock()
		purchaseEventUpdateCache[key] = cache
		purchaseEventUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q purchaseEventQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q purchaseEventQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for purchase_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for purchase_event")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o PurchaseEventSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PurchaseEventSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), purchaseEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"purchase_event\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, purchaseEventPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in purchaseEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all purchaseEvent")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *PurchaseEvent) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PurchaseEvent) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no purchase_event provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(purchaseEventColumnsWithDefault, o)

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

	purchaseEventUpsertCacheMut.RLock()
	cache, cached := purchaseEventUpsertCache[key]
	purchaseEventUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			purchaseEventAllColumns,
			purchaseEventColumnsWithDefault,
			purchaseEventColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			purchaseEventAllColumns,
			purchaseEventPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert purchase_event, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(purchaseEventPrimaryKeyColumns))
			copy(conflict, purchaseEventPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"purchase_event\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(purchaseEventType, purchaseEventMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(purchaseEventType, purchaseEventMapping, ret)
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
		return errors.Wrap(err, "database: unable to upsert purchase_event")
	}

	if !cached {
		purchaseEventUpsertCacheMut.Lock()
		purchaseEventUpsertCache[key] = cache
		purchaseEventUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single PurchaseEvent record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *PurchaseEvent) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single PurchaseEvent record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PurchaseEvent) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no PurchaseEvent provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), purchaseEventPrimaryKeyMapping)
	sql := "DELETE FROM \"purchase_event\" WHERE \"stripe_session_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from purchase_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for purchase_event")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q purchaseEventQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q purchaseEventQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no purchaseEventQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from purchase_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for purchase_event")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o PurchaseEventSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PurchaseEventSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(purchaseEventBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), purchaseEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"purchase_event\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, purchaseEventPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from purchaseEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for purchase_event")
	}

	if len(purchaseEventAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *PurchaseEvent) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no PurchaseEvent provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PurchaseEvent) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPurchaseEvent(ctx, exec, o.StripeSessionID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PurchaseEventSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty PurchaseEventSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PurchaseEventSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PurchaseEventSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), purchaseEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"purchase_event\".* FROM \"purchase_event\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, purchaseEventPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in PurchaseEventSlice")
	}

	*o = slice

	return nil
}

// PurchaseEventExistsG checks if the PurchaseEvent row exists.
func PurchaseEventExistsG(ctx context.Context, stripeSessionID string) (bool, error) {
	return PurchaseEventExists(ctx, boil.GetContextDB(), stripeSessionID)
}

// PurchaseEventExists checks if the PurchaseEvent row exists.
func PurchaseEventExists(ctx context.Context, exec boil.ContextExecutor, stripeSessionID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"purchase_event\" where \"stripe_session_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, stripeSessionID)
	}
	row := exec.QueryRowContext(ctx, sql, stripeSessionID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if purchase_event exists")
	}

	return exists, nil
}

// Exists checks if the PurchaseEvent row exists.
func (o *PurchaseEvent) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return PurchaseEventExists(ctx, exec, o.StripeSessionID)
}
