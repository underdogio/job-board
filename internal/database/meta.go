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

// Metum is an object representing the database table.
type Metum struct {
	Key   string `boil:"key" json:"key" toml:"key" yaml:"key"`
	Value string `boil:"value" json:"value" toml:"value" yaml:"value"`

	R *metumR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L metumL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MetumColumns = struct {
	Key   string
	Value string
}{
	Key:   "key",
	Value: "value",
}

var MetumTableColumns = struct {
	Key   string
	Value string
}{
	Key:   "meta.key",
	Value: "meta.value",
}

// Generated where

var MetumWhere = struct {
	Key   whereHelperstring
	Value whereHelperstring
}{
	Key:   whereHelperstring{field: "\"meta\".\"key\""},
	Value: whereHelperstring{field: "\"meta\".\"value\""},
}

// MetumRels is where relationship names are stored.
var MetumRels = struct {
}{}

// metumR is where relationships are stored.
type metumR struct {
}

// NewStruct creates a new relationship struct
func (*metumR) NewStruct() *metumR {
	return &metumR{}
}

// metumL is where Load methods for each relationship are stored.
type metumL struct{}

var (
	metumAllColumns            = []string{"key", "value"}
	metumColumnsWithoutDefault = []string{"key", "value"}
	metumColumnsWithDefault    = []string{}
	metumPrimaryKeyColumns     = []string{"key"}
	metumGeneratedColumns      = []string{}
)

type (
	// MetumSlice is an alias for a slice of pointers to Metum.
	// This should almost always be used instead of []Metum.
	MetumSlice []*Metum
	// MetumHook is the signature for custom Metum hook methods
	MetumHook func(context.Context, boil.ContextExecutor, *Metum) error

	metumQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	metumType                 = reflect.TypeOf(&Metum{})
	metumMapping              = queries.MakeStructMapping(metumType)
	metumPrimaryKeyMapping, _ = queries.BindMapping(metumType, metumMapping, metumPrimaryKeyColumns)
	metumInsertCacheMut       sync.RWMutex
	metumInsertCache          = make(map[string]insertCache)
	metumUpdateCacheMut       sync.RWMutex
	metumUpdateCache          = make(map[string]updateCache)
	metumUpsertCacheMut       sync.RWMutex
	metumUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var metumAfterSelectHooks []MetumHook

var metumBeforeInsertHooks []MetumHook
var metumAfterInsertHooks []MetumHook

var metumBeforeUpdateHooks []MetumHook
var metumAfterUpdateHooks []MetumHook

var metumBeforeDeleteHooks []MetumHook
var metumAfterDeleteHooks []MetumHook

var metumBeforeUpsertHooks []MetumHook
var metumAfterUpsertHooks []MetumHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Metum) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Metum) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Metum) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Metum) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Metum) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Metum) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Metum) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Metum) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Metum) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range metumAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMetumHook registers your hook function for all future operations.
func AddMetumHook(hookPoint boil.HookPoint, metumHook MetumHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		metumAfterSelectHooks = append(metumAfterSelectHooks, metumHook)
	case boil.BeforeInsertHook:
		metumBeforeInsertHooks = append(metumBeforeInsertHooks, metumHook)
	case boil.AfterInsertHook:
		metumAfterInsertHooks = append(metumAfterInsertHooks, metumHook)
	case boil.BeforeUpdateHook:
		metumBeforeUpdateHooks = append(metumBeforeUpdateHooks, metumHook)
	case boil.AfterUpdateHook:
		metumAfterUpdateHooks = append(metumAfterUpdateHooks, metumHook)
	case boil.BeforeDeleteHook:
		metumBeforeDeleteHooks = append(metumBeforeDeleteHooks, metumHook)
	case boil.AfterDeleteHook:
		metumAfterDeleteHooks = append(metumAfterDeleteHooks, metumHook)
	case boil.BeforeUpsertHook:
		metumBeforeUpsertHooks = append(metumBeforeUpsertHooks, metumHook)
	case boil.AfterUpsertHook:
		metumAfterUpsertHooks = append(metumAfterUpsertHooks, metumHook)
	}
}

// OneG returns a single metum record from the query using the global executor.
func (q metumQuery) OneG(ctx context.Context) (*Metum, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single metum record from the query.
func (q metumQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Metum, error) {
	o := &Metum{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for meta")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all Metum records from the query using the global executor.
func (q metumQuery) AllG(ctx context.Context) (MetumSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all Metum records from the query.
func (q metumQuery) All(ctx context.Context, exec boil.ContextExecutor) (MetumSlice, error) {
	var o []*Metum

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to Metum slice")
	}

	if len(metumAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all Metum records in the query using the global executor
func (q metumQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all Metum records in the query.
func (q metumQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count meta rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q metumQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q metumQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if meta exists")
	}

	return count > 0, nil
}

// Meta retrieves all the records using an executor.
func Meta(mods ...qm.QueryMod) metumQuery {
	mods = append(mods, qm.From("\"meta\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"meta\".*"})
	}

	return metumQuery{q}
}

// FindMetumG retrieves a single record by ID.
func FindMetumG(ctx context.Context, key string, selectCols ...string) (*Metum, error) {
	return FindMetum(ctx, boil.GetContextDB(), key, selectCols...)
}

// FindMetum retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMetum(ctx context.Context, exec boil.ContextExecutor, key string, selectCols ...string) (*Metum, error) {
	metumObj := &Metum{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"meta\" where \"key\"=$1", sel,
	)

	q := queries.Raw(query, key)

	err := q.Bind(ctx, exec, metumObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from meta")
	}

	if err = metumObj.doAfterSelectHooks(ctx, exec); err != nil {
		return metumObj, err
	}

	return metumObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Metum) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Metum) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no meta provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(metumColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	metumInsertCacheMut.RLock()
	cache, cached := metumInsertCache[key]
	metumInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			metumAllColumns,
			metumColumnsWithDefault,
			metumColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(metumType, metumMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(metumType, metumMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"meta\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"meta\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "database: unable to insert into meta")
	}

	if !cached {
		metumInsertCacheMut.Lock()
		metumInsertCache[key] = cache
		metumInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single Metum record using the global executor.
// See Update for more documentation.
func (o *Metum) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the Metum.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Metum) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	metumUpdateCacheMut.RLock()
	cache, cached := metumUpdateCache[key]
	metumUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			metumAllColumns,
			metumPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update meta, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"meta\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, metumPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(metumType, metumMapping, append(wl, metumPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "database: unable to update meta row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for meta")
	}

	if !cached {
		metumUpdateCacheMut.Lock()
		metumUpdateCache[key] = cache
		metumUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q metumQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q metumQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for meta")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for meta")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o MetumSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MetumSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), metumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"meta\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, metumPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in metum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all metum")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Metum) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Metum) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no meta provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(metumColumnsWithDefault, o)

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

	metumUpsertCacheMut.RLock()
	cache, cached := metumUpsertCache[key]
	metumUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			metumAllColumns,
			metumColumnsWithDefault,
			metumColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			metumAllColumns,
			metumPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert meta, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(metumPrimaryKeyColumns))
			copy(conflict, metumPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"meta\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(metumType, metumMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(metumType, metumMapping, ret)
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
		return errors.Wrap(err, "database: unable to upsert meta")
	}

	if !cached {
		metumUpsertCacheMut.Lock()
		metumUpsertCache[key] = cache
		metumUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single Metum record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Metum) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single Metum record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Metum) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no Metum provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), metumPrimaryKeyMapping)
	sql := "DELETE FROM \"meta\" WHERE \"key\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from meta")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for meta")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q metumQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q metumQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no metumQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from meta")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for meta")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o MetumSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MetumSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(metumBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), metumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"meta\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, metumPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from metum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for meta")
	}

	if len(metumAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Metum) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no Metum provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Metum) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMetum(ctx, exec, o.Key)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MetumSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty MetumSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MetumSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MetumSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), metumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"meta\".* FROM \"meta\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, metumPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in MetumSlice")
	}

	*o = slice

	return nil
}

// MetumExistsG checks if the Metum row exists.
func MetumExistsG(ctx context.Context, key string) (bool, error) {
	return MetumExists(ctx, boil.GetContextDB(), key)
}

// MetumExists checks if the Metum row exists.
func MetumExists(ctx context.Context, exec boil.ContextExecutor, key string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"meta\" where \"key\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, key)
	}
	row := exec.QueryRowContext(ctx, sql, key)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if meta exists")
	}

	return exists, nil
}

// Exists checks if the Metum row exists.
func (o *Metum) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return MetumExists(ctx, exec, o.Key)
}
