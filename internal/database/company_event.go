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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// CompanyEvent is an object representing the database table.
type CompanyEvent struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	EventType string    `boil:"event_type" json:"event_type" toml:"event_type" yaml:"event_type"`
	CompanyID string    `boil:"company_id" json:"company_id" toml:"company_id" yaml:"company_id"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *companyEventR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L companyEventL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CompanyEventColumns = struct {
	ID        string
	EventType string
	CompanyID string
	CreatedAt string
}{
	ID:        "id",
	EventType: "event_type",
	CompanyID: "company_id",
	CreatedAt: "created_at",
}

var CompanyEventTableColumns = struct {
	ID        string
	EventType string
	CompanyID string
	CreatedAt string
}{
	ID:        "company_event.id",
	EventType: "company_event.event_type",
	CompanyID: "company_event.company_id",
	CreatedAt: "company_event.created_at",
}

// Generated where

var CompanyEventWhere = struct {
	ID        whereHelperstring
	EventType whereHelperstring
	CompanyID whereHelperstring
	CreatedAt whereHelpertime_Time
}{
	ID:        whereHelperstring{field: "\"company_event\".\"id\""},
	EventType: whereHelperstring{field: "\"company_event\".\"event_type\""},
	CompanyID: whereHelperstring{field: "\"company_event\".\"company_id\""},
	CreatedAt: whereHelpertime_Time{field: "\"company_event\".\"created_at\""},
}

// CompanyEventRels is where relationship names are stored.
var CompanyEventRels = struct {
	Company string
}{
	Company: "Company",
}

// companyEventR is where relationships are stored.
type companyEventR struct {
	Company *Company `boil:"Company" json:"Company" toml:"Company" yaml:"Company"`
}

// NewStruct creates a new relationship struct
func (*companyEventR) NewStruct() *companyEventR {
	return &companyEventR{}
}

func (r *companyEventR) GetCompany() *Company {
	if r == nil {
		return nil
	}
	return r.Company
}

// companyEventL is where Load methods for each relationship are stored.
type companyEventL struct{}

var (
	companyEventAllColumns            = []string{"id", "event_type", "company_id", "created_at"}
	companyEventColumnsWithoutDefault = []string{"id", "event_type", "company_id", "created_at"}
	companyEventColumnsWithDefault    = []string{}
	companyEventPrimaryKeyColumns     = []string{"id"}
	companyEventGeneratedColumns      = []string{}
)

type (
	// CompanyEventSlice is an alias for a slice of pointers to CompanyEvent.
	// This should almost always be used instead of []CompanyEvent.
	CompanyEventSlice []*CompanyEvent
	// CompanyEventHook is the signature for custom CompanyEvent hook methods
	CompanyEventHook func(context.Context, boil.ContextExecutor, *CompanyEvent) error

	companyEventQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	companyEventType                 = reflect.TypeOf(&CompanyEvent{})
	companyEventMapping              = queries.MakeStructMapping(companyEventType)
	companyEventPrimaryKeyMapping, _ = queries.BindMapping(companyEventType, companyEventMapping, companyEventPrimaryKeyColumns)
	companyEventInsertCacheMut       sync.RWMutex
	companyEventInsertCache          = make(map[string]insertCache)
	companyEventUpdateCacheMut       sync.RWMutex
	companyEventUpdateCache          = make(map[string]updateCache)
	companyEventUpsertCacheMut       sync.RWMutex
	companyEventUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var companyEventAfterSelectHooks []CompanyEventHook

var companyEventBeforeInsertHooks []CompanyEventHook
var companyEventAfterInsertHooks []CompanyEventHook

var companyEventBeforeUpdateHooks []CompanyEventHook
var companyEventAfterUpdateHooks []CompanyEventHook

var companyEventBeforeDeleteHooks []CompanyEventHook
var companyEventAfterDeleteHooks []CompanyEventHook

var companyEventBeforeUpsertHooks []CompanyEventHook
var companyEventAfterUpsertHooks []CompanyEventHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *CompanyEvent) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *CompanyEvent) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *CompanyEvent) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *CompanyEvent) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *CompanyEvent) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *CompanyEvent) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *CompanyEvent) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *CompanyEvent) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *CompanyEvent) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range companyEventAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddCompanyEventHook registers your hook function for all future operations.
func AddCompanyEventHook(hookPoint boil.HookPoint, companyEventHook CompanyEventHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		companyEventAfterSelectHooks = append(companyEventAfterSelectHooks, companyEventHook)
	case boil.BeforeInsertHook:
		companyEventBeforeInsertHooks = append(companyEventBeforeInsertHooks, companyEventHook)
	case boil.AfterInsertHook:
		companyEventAfterInsertHooks = append(companyEventAfterInsertHooks, companyEventHook)
	case boil.BeforeUpdateHook:
		companyEventBeforeUpdateHooks = append(companyEventBeforeUpdateHooks, companyEventHook)
	case boil.AfterUpdateHook:
		companyEventAfterUpdateHooks = append(companyEventAfterUpdateHooks, companyEventHook)
	case boil.BeforeDeleteHook:
		companyEventBeforeDeleteHooks = append(companyEventBeforeDeleteHooks, companyEventHook)
	case boil.AfterDeleteHook:
		companyEventAfterDeleteHooks = append(companyEventAfterDeleteHooks, companyEventHook)
	case boil.BeforeUpsertHook:
		companyEventBeforeUpsertHooks = append(companyEventBeforeUpsertHooks, companyEventHook)
	case boil.AfterUpsertHook:
		companyEventAfterUpsertHooks = append(companyEventAfterUpsertHooks, companyEventHook)
	}
}

// OneG returns a single companyEvent record from the query using the global executor.
func (q companyEventQuery) OneG(ctx context.Context) (*CompanyEvent, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single companyEvent record from the query.
func (q companyEventQuery) One(ctx context.Context, exec boil.ContextExecutor) (*CompanyEvent, error) {
	o := &CompanyEvent{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for company_event")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all CompanyEvent records from the query using the global executor.
func (q companyEventQuery) AllG(ctx context.Context) (CompanyEventSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all CompanyEvent records from the query.
func (q companyEventQuery) All(ctx context.Context, exec boil.ContextExecutor) (CompanyEventSlice, error) {
	var o []*CompanyEvent

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to CompanyEvent slice")
	}

	if len(companyEventAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all CompanyEvent records in the query using the global executor
func (q companyEventQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all CompanyEvent records in the query.
func (q companyEventQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count company_event rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q companyEventQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q companyEventQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if company_event exists")
	}

	return count > 0, nil
}

// Company pointed to by the foreign key.
func (o *CompanyEvent) Company(mods ...qm.QueryMod) companyQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.CompanyID),
	}

	queryMods = append(queryMods, mods...)

	return Companies(queryMods...)
}

// LoadCompany allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (companyEventL) LoadCompany(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCompanyEvent interface{}, mods queries.Applicator) error {
	var slice []*CompanyEvent
	var object *CompanyEvent

	if singular {
		var ok bool
		object, ok = maybeCompanyEvent.(*CompanyEvent)
		if !ok {
			object = new(CompanyEvent)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCompanyEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCompanyEvent))
			}
		}
	} else {
		s, ok := maybeCompanyEvent.(*[]*CompanyEvent)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCompanyEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCompanyEvent))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &companyEventR{}
		}
		args = append(args, object.CompanyID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &companyEventR{}
			}

			for _, a := range args {
				if a == obj.CompanyID {
					continue Outer
				}
			}

			args = append(args, obj.CompanyID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`company`),
		qm.WhereIn(`company.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Company")
	}

	var resultSlice []*Company
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Company")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for company")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for company")
	}

	if len(companyAfterSelectHooks) != 0 {
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
		object.R.Company = foreign
		if foreign.R == nil {
			foreign.R = &companyR{}
		}
		foreign.R.CompanyEvents = append(foreign.R.CompanyEvents, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CompanyID == foreign.ID {
				local.R.Company = foreign
				if foreign.R == nil {
					foreign.R = &companyR{}
				}
				foreign.R.CompanyEvents = append(foreign.R.CompanyEvents, local)
				break
			}
		}
	}

	return nil
}

// SetCompanyG of the companyEvent to the related item.
// Sets o.R.Company to related.
// Adds o to related.R.CompanyEvents.
// Uses the global database handle.
func (o *CompanyEvent) SetCompanyG(ctx context.Context, insert bool, related *Company) error {
	return o.SetCompany(ctx, boil.GetContextDB(), insert, related)
}

// SetCompany of the companyEvent to the related item.
// Sets o.R.Company to related.
// Adds o to related.R.CompanyEvents.
func (o *CompanyEvent) SetCompany(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Company) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"company_event\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"company_id"}),
		strmangle.WhereClause("\"", "\"", 2, companyEventPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CompanyID = related.ID
	if o.R == nil {
		o.R = &companyEventR{
			Company: related,
		}
	} else {
		o.R.Company = related
	}

	if related.R == nil {
		related.R = &companyR{
			CompanyEvents: CompanyEventSlice{o},
		}
	} else {
		related.R.CompanyEvents = append(related.R.CompanyEvents, o)
	}

	return nil
}

// CompanyEvents retrieves all the records using an executor.
func CompanyEvents(mods ...qm.QueryMod) companyEventQuery {
	mods = append(mods, qm.From("\"company_event\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"company_event\".*"})
	}

	return companyEventQuery{q}
}

// FindCompanyEventG retrieves a single record by ID.
func FindCompanyEventG(ctx context.Context, iD string, selectCols ...string) (*CompanyEvent, error) {
	return FindCompanyEvent(ctx, boil.GetContextDB(), iD, selectCols...)
}

// FindCompanyEvent retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCompanyEvent(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*CompanyEvent, error) {
	companyEventObj := &CompanyEvent{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"company_event\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, companyEventObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from company_event")
	}

	if err = companyEventObj.doAfterSelectHooks(ctx, exec); err != nil {
		return companyEventObj, err
	}

	return companyEventObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *CompanyEvent) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *CompanyEvent) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no company_event provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(companyEventColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	companyEventInsertCacheMut.RLock()
	cache, cached := companyEventInsertCache[key]
	companyEventInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			companyEventAllColumns,
			companyEventColumnsWithDefault,
			companyEventColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(companyEventType, companyEventMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(companyEventType, companyEventMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"company_event\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"company_event\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "database: unable to insert into company_event")
	}

	if !cached {
		companyEventInsertCacheMut.Lock()
		companyEventInsertCache[key] = cache
		companyEventInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single CompanyEvent record using the global executor.
// See Update for more documentation.
func (o *CompanyEvent) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the CompanyEvent.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *CompanyEvent) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	companyEventUpdateCacheMut.RLock()
	cache, cached := companyEventUpdateCache[key]
	companyEventUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			companyEventAllColumns,
			companyEventPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update company_event, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"company_event\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, companyEventPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(companyEventType, companyEventMapping, append(wl, companyEventPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "database: unable to update company_event row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for company_event")
	}

	if !cached {
		companyEventUpdateCacheMut.Lock()
		companyEventUpdateCache[key] = cache
		companyEventUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q companyEventQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q companyEventQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for company_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for company_event")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o CompanyEventSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CompanyEventSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), companyEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"company_event\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, companyEventPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in companyEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all companyEvent")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *CompanyEvent) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *CompanyEvent) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no company_event provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(companyEventColumnsWithDefault, o)

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

	companyEventUpsertCacheMut.RLock()
	cache, cached := companyEventUpsertCache[key]
	companyEventUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			companyEventAllColumns,
			companyEventColumnsWithDefault,
			companyEventColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			companyEventAllColumns,
			companyEventPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert company_event, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(companyEventPrimaryKeyColumns))
			copy(conflict, companyEventPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"company_event\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(companyEventType, companyEventMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(companyEventType, companyEventMapping, ret)
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
		return errors.Wrap(err, "database: unable to upsert company_event")
	}

	if !cached {
		companyEventUpsertCacheMut.Lock()
		companyEventUpsertCache[key] = cache
		companyEventUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single CompanyEvent record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *CompanyEvent) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single CompanyEvent record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CompanyEvent) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no CompanyEvent provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), companyEventPrimaryKeyMapping)
	sql := "DELETE FROM \"company_event\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from company_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for company_event")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q companyEventQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q companyEventQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no companyEventQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from company_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for company_event")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o CompanyEventSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CompanyEventSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(companyEventBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), companyEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"company_event\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, companyEventPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from companyEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for company_event")
	}

	if len(companyEventAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *CompanyEvent) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no CompanyEvent provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *CompanyEvent) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCompanyEvent(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CompanyEventSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty CompanyEventSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CompanyEventSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CompanyEventSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), companyEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"company_event\".* FROM \"company_event\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, companyEventPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in CompanyEventSlice")
	}

	*o = slice

	return nil
}

// CompanyEventExistsG checks if the CompanyEvent row exists.
func CompanyEventExistsG(ctx context.Context, iD string) (bool, error) {
	return CompanyEventExists(ctx, boil.GetContextDB(), iD)
}

// CompanyEventExists checks if the CompanyEvent row exists.
func CompanyEventExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"company_event\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if company_event exists")
	}

	return exists, nil
}

// Exists checks if the CompanyEvent row exists.
func (o *CompanyEvent) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return CompanyEventExists(ctx, exec, o.ID)
}
