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

// SearchEvent is an object representing the database table.
type SearchEvent struct {
	SessionID string      `boil:"session_id" json:"session_id" toml:"session_id" yaml:"session_id"`
	Location  null.String `boil:"location" json:"location,omitempty" toml:"location" yaml:"location,omitempty"`
	Tag       null.String `boil:"tag" json:"tag,omitempty" toml:"tag" yaml:"tag,omitempty"`
	Results   int         `boil:"results" json:"results" toml:"results" yaml:"results"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Type      null.String `boil:"type" json:"type,omitempty" toml:"type" yaml:"type,omitempty"`

	R *searchEventR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L searchEventL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SearchEventColumns = struct {
	SessionID string
	Location  string
	Tag       string
	Results   string
	CreatedAt string
	Type      string
}{
	SessionID: "session_id",
	Location:  "location",
	Tag:       "tag",
	Results:   "results",
	CreatedAt: "created_at",
	Type:      "type",
}

var SearchEventTableColumns = struct {
	SessionID string
	Location  string
	Tag       string
	Results   string
	CreatedAt string
	Type      string
}{
	SessionID: "search_event.session_id",
	Location:  "search_event.location",
	Tag:       "search_event.tag",
	Results:   "search_event.results",
	CreatedAt: "search_event.created_at",
	Type:      "search_event.type",
}

// Generated where

var SearchEventWhere = struct {
	SessionID whereHelperstring
	Location  whereHelpernull_String
	Tag       whereHelpernull_String
	Results   whereHelperint
	CreatedAt whereHelpertime_Time
	Type      whereHelpernull_String
}{
	SessionID: whereHelperstring{field: "\"search_event\".\"session_id\""},
	Location:  whereHelpernull_String{field: "\"search_event\".\"location\""},
	Tag:       whereHelpernull_String{field: "\"search_event\".\"tag\""},
	Results:   whereHelperint{field: "\"search_event\".\"results\""},
	CreatedAt: whereHelpertime_Time{field: "\"search_event\".\"created_at\""},
	Type:      whereHelpernull_String{field: "\"search_event\".\"type\""},
}

// SearchEventRels is where relationship names are stored.
var SearchEventRels = struct {
}{}

// searchEventR is where relationships are stored.
type searchEventR struct {
}

// NewStruct creates a new relationship struct
func (*searchEventR) NewStruct() *searchEventR {
	return &searchEventR{}
}

// searchEventL is where Load methods for each relationship are stored.
type searchEventL struct{}

var (
	searchEventAllColumns            = []string{"session_id", "location", "tag", "results", "created_at", "type"}
	searchEventColumnsWithoutDefault = []string{"session_id", "results", "created_at"}
	searchEventColumnsWithDefault    = []string{"location", "tag", "type"}
	searchEventPrimaryKeyColumns     = []string{"session_id"}
	searchEventGeneratedColumns      = []string{}
)

type (
	// SearchEventSlice is an alias for a slice of pointers to SearchEvent.
	// This should almost always be used instead of []SearchEvent.
	SearchEventSlice []*SearchEvent
	// SearchEventHook is the signature for custom SearchEvent hook methods
	SearchEventHook func(context.Context, boil.ContextExecutor, *SearchEvent) error

	searchEventQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	searchEventType                 = reflect.TypeOf(&SearchEvent{})
	searchEventMapping              = queries.MakeStructMapping(searchEventType)
	searchEventPrimaryKeyMapping, _ = queries.BindMapping(searchEventType, searchEventMapping, searchEventPrimaryKeyColumns)
	searchEventInsertCacheMut       sync.RWMutex
	searchEventInsertCache          = make(map[string]insertCache)
	searchEventUpdateCacheMut       sync.RWMutex
	searchEventUpdateCache          = make(map[string]updateCache)
	searchEventUpsertCacheMut       sync.RWMutex
	searchEventUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var searchEventAfterSelectHooks []SearchEventHook

var searchEventBeforeInsertHooks []SearchEventHook
var searchEventAfterInsertHooks []SearchEventHook

var searchEventBeforeUpdateHooks []SearchEventHook
var searchEventAfterUpdateHooks []SearchEventHook

var searchEventBeforeDeleteHooks []SearchEventHook
var searchEventAfterDeleteHooks []SearchEventHook

var searchEventBeforeUpsertHooks []SearchEventHook
var searchEventAfterUpsertHooks []SearchEventHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SearchEvent) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SearchEvent) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SearchEvent) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SearchEvent) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SearchEvent) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SearchEvent) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SearchEvent) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SearchEvent) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SearchEvent) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchEventAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSearchEventHook registers your hook function for all future operations.
func AddSearchEventHook(hookPoint boil.HookPoint, searchEventHook SearchEventHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		searchEventAfterSelectHooks = append(searchEventAfterSelectHooks, searchEventHook)
	case boil.BeforeInsertHook:
		searchEventBeforeInsertHooks = append(searchEventBeforeInsertHooks, searchEventHook)
	case boil.AfterInsertHook:
		searchEventAfterInsertHooks = append(searchEventAfterInsertHooks, searchEventHook)
	case boil.BeforeUpdateHook:
		searchEventBeforeUpdateHooks = append(searchEventBeforeUpdateHooks, searchEventHook)
	case boil.AfterUpdateHook:
		searchEventAfterUpdateHooks = append(searchEventAfterUpdateHooks, searchEventHook)
	case boil.BeforeDeleteHook:
		searchEventBeforeDeleteHooks = append(searchEventBeforeDeleteHooks, searchEventHook)
	case boil.AfterDeleteHook:
		searchEventAfterDeleteHooks = append(searchEventAfterDeleteHooks, searchEventHook)
	case boil.BeforeUpsertHook:
		searchEventBeforeUpsertHooks = append(searchEventBeforeUpsertHooks, searchEventHook)
	case boil.AfterUpsertHook:
		searchEventAfterUpsertHooks = append(searchEventAfterUpsertHooks, searchEventHook)
	}
}

// OneG returns a single searchEvent record from the query using the global executor.
func (q searchEventQuery) OneG(ctx context.Context) (*SearchEvent, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single searchEvent record from the query.
func (q searchEventQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SearchEvent, error) {
	o := &SearchEvent{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for search_event")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all SearchEvent records from the query using the global executor.
func (q searchEventQuery) AllG(ctx context.Context) (SearchEventSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all SearchEvent records from the query.
func (q searchEventQuery) All(ctx context.Context, exec boil.ContextExecutor) (SearchEventSlice, error) {
	var o []*SearchEvent

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to SearchEvent slice")
	}

	if len(searchEventAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all SearchEvent records in the query using the global executor
func (q searchEventQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all SearchEvent records in the query.
func (q searchEventQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count search_event rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q searchEventQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q searchEventQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if search_event exists")
	}

	return count > 0, nil
}

// SearchEvents retrieves all the records using an executor.
func SearchEvents(mods ...qm.QueryMod) searchEventQuery {
	mods = append(mods, qm.From("\"search_event\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"search_event\".*"})
	}

	return searchEventQuery{q}
}

// FindSearchEventG retrieves a single record by ID.
func FindSearchEventG(ctx context.Context, sessionID string, selectCols ...string) (*SearchEvent, error) {
	return FindSearchEvent(ctx, boil.GetContextDB(), sessionID, selectCols...)
}

// FindSearchEvent retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSearchEvent(ctx context.Context, exec boil.ContextExecutor, sessionID string, selectCols ...string) (*SearchEvent, error) {
	searchEventObj := &SearchEvent{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"search_event\" where \"session_id\"=$1", sel,
	)

	q := queries.Raw(query, sessionID)

	err := q.Bind(ctx, exec, searchEventObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from search_event")
	}

	if err = searchEventObj.doAfterSelectHooks(ctx, exec); err != nil {
		return searchEventObj, err
	}

	return searchEventObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *SearchEvent) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SearchEvent) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no search_event provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(searchEventColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	searchEventInsertCacheMut.RLock()
	cache, cached := searchEventInsertCache[key]
	searchEventInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			searchEventAllColumns,
			searchEventColumnsWithDefault,
			searchEventColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(searchEventType, searchEventMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(searchEventType, searchEventMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"search_event\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"search_event\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "database: unable to insert into search_event")
	}

	if !cached {
		searchEventInsertCacheMut.Lock()
		searchEventInsertCache[key] = cache
		searchEventInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single SearchEvent record using the global executor.
// See Update for more documentation.
func (o *SearchEvent) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the SearchEvent.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SearchEvent) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	searchEventUpdateCacheMut.RLock()
	cache, cached := searchEventUpdateCache[key]
	searchEventUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			searchEventAllColumns,
			searchEventPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update search_event, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"search_event\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, searchEventPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(searchEventType, searchEventMapping, append(wl, searchEventPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "database: unable to update search_event row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for search_event")
	}

	if !cached {
		searchEventUpdateCacheMut.Lock()
		searchEventUpdateCache[key] = cache
		searchEventUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q searchEventQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q searchEventQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for search_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for search_event")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SearchEventSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SearchEventSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), searchEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"search_event\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, searchEventPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in searchEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all searchEvent")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *SearchEvent) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SearchEvent) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no search_event provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(searchEventColumnsWithDefault, o)

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

	searchEventUpsertCacheMut.RLock()
	cache, cached := searchEventUpsertCache[key]
	searchEventUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			searchEventAllColumns,
			searchEventColumnsWithDefault,
			searchEventColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			searchEventAllColumns,
			searchEventPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert search_event, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(searchEventPrimaryKeyColumns))
			copy(conflict, searchEventPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"search_event\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(searchEventType, searchEventMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(searchEventType, searchEventMapping, ret)
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
		return errors.Wrap(err, "database: unable to upsert search_event")
	}

	if !cached {
		searchEventUpsertCacheMut.Lock()
		searchEventUpsertCache[key] = cache
		searchEventUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single SearchEvent record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *SearchEvent) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single SearchEvent record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SearchEvent) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no SearchEvent provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), searchEventPrimaryKeyMapping)
	sql := "DELETE FROM \"search_event\" WHERE \"session_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from search_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for search_event")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q searchEventQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q searchEventQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no searchEventQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from search_event")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for search_event")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o SearchEventSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SearchEventSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(searchEventBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), searchEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"search_event\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, searchEventPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from searchEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for search_event")
	}

	if len(searchEventAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *SearchEvent) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no SearchEvent provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SearchEvent) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSearchEvent(ctx, exec, o.SessionID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SearchEventSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty SearchEventSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SearchEventSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SearchEventSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), searchEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"search_event\".* FROM \"search_event\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, searchEventPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in SearchEventSlice")
	}

	*o = slice

	return nil
}

// SearchEventExistsG checks if the SearchEvent row exists.
func SearchEventExistsG(ctx context.Context, sessionID string) (bool, error) {
	return SearchEventExists(ctx, boil.GetContextDB(), sessionID)
}

// SearchEventExists checks if the SearchEvent row exists.
func SearchEventExists(ctx context.Context, exec boil.ContextExecutor, sessionID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"search_event\" where \"session_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, sessionID)
	}
	row := exec.QueryRowContext(ctx, sql, sessionID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if search_event exists")
	}

	return exists, nil
}

// Exists checks if the SearchEvent row exists.
func (o *SearchEvent) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SearchEventExists(ctx, exec, o.SessionID)
}
