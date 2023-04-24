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

// SeoLocation is an object representing the database table.
type SeoLocation struct {
	Name       string       `boil:"name" json:"name" toml:"name" yaml:"name"`
	Currency   string       `boil:"currency" json:"currency" toml:"currency" yaml:"currency"`
	Country    null.String  `boil:"country" json:"country,omitempty" toml:"country" yaml:"country,omitempty"`
	Iso2       null.String  `boil:"iso2" json:"iso2,omitempty" toml:"iso2" yaml:"iso2,omitempty"`
	Region     null.String  `boil:"region" json:"region,omitempty" toml:"region" yaml:"region,omitempty"`
	Population null.Int     `boil:"population" json:"population,omitempty" toml:"population" yaml:"population,omitempty"`
	Lat        null.Float64 `boil:"lat" json:"lat,omitempty" toml:"lat" yaml:"lat,omitempty"`
	Long       null.Float64 `boil:"long" json:"long,omitempty" toml:"long" yaml:"long,omitempty"`
	Emoji      null.String  `boil:"emoji" json:"emoji,omitempty" toml:"emoji" yaml:"emoji,omitempty"`

	R *seoLocationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L seoLocationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SeoLocationColumns = struct {
	Name       string
	Currency   string
	Country    string
	Iso2       string
	Region     string
	Population string
	Lat        string
	Long       string
	Emoji      string
}{
	Name:       "name",
	Currency:   "currency",
	Country:    "country",
	Iso2:       "iso2",
	Region:     "region",
	Population: "population",
	Lat:        "lat",
	Long:       "long",
	Emoji:      "emoji",
}

var SeoLocationTableColumns = struct {
	Name       string
	Currency   string
	Country    string
	Iso2       string
	Region     string
	Population string
	Lat        string
	Long       string
	Emoji      string
}{
	Name:       "seo_location.name",
	Currency:   "seo_location.currency",
	Country:    "seo_location.country",
	Iso2:       "seo_location.iso2",
	Region:     "seo_location.region",
	Population: "seo_location.population",
	Lat:        "seo_location.lat",
	Long:       "seo_location.long",
	Emoji:      "seo_location.emoji",
}

// Generated where

type whereHelpernull_Int struct{ field string }

func (w whereHelpernull_Int) EQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int) NEQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int) LT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int) LTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int) GT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int) GTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_Int) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_Int) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_Int) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelpernull_Float64 struct{ field string }

func (w whereHelpernull_Float64) EQ(x null.Float64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Float64) NEQ(x null.Float64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Float64) LT(x null.Float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Float64) LTE(x null.Float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Float64) GT(x null.Float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Float64) GTE(x null.Float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_Float64) IN(slice []float64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_Float64) NIN(slice []float64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_Float64) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Float64) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var SeoLocationWhere = struct {
	Name       whereHelperstring
	Currency   whereHelperstring
	Country    whereHelpernull_String
	Iso2       whereHelpernull_String
	Region     whereHelpernull_String
	Population whereHelpernull_Int
	Lat        whereHelpernull_Float64
	Long       whereHelpernull_Float64
	Emoji      whereHelpernull_String
}{
	Name:       whereHelperstring{field: "\"seo_location\".\"name\""},
	Currency:   whereHelperstring{field: "\"seo_location\".\"currency\""},
	Country:    whereHelpernull_String{field: "\"seo_location\".\"country\""},
	Iso2:       whereHelpernull_String{field: "\"seo_location\".\"iso2\""},
	Region:     whereHelpernull_String{field: "\"seo_location\".\"region\""},
	Population: whereHelpernull_Int{field: "\"seo_location\".\"population\""},
	Lat:        whereHelpernull_Float64{field: "\"seo_location\".\"lat\""},
	Long:       whereHelpernull_Float64{field: "\"seo_location\".\"long\""},
	Emoji:      whereHelpernull_String{field: "\"seo_location\".\"emoji\""},
}

// SeoLocationRels is where relationship names are stored.
var SeoLocationRels = struct {
}{}

// seoLocationR is where relationships are stored.
type seoLocationR struct {
}

// NewStruct creates a new relationship struct
func (*seoLocationR) NewStruct() *seoLocationR {
	return &seoLocationR{}
}

// seoLocationL is where Load methods for each relationship are stored.
type seoLocationL struct{}

var (
	seoLocationAllColumns            = []string{"name", "currency", "country", "iso2", "region", "population", "lat", "long", "emoji"}
	seoLocationColumnsWithoutDefault = []string{"name"}
	seoLocationColumnsWithDefault    = []string{"currency", "country", "iso2", "region", "population", "lat", "long", "emoji"}
	seoLocationPrimaryKeyColumns     = []string{"name"}
	seoLocationGeneratedColumns      = []string{}
)

type (
	// SeoLocationSlice is an alias for a slice of pointers to SeoLocation.
	// This should almost always be used instead of []SeoLocation.
	SeoLocationSlice []*SeoLocation
	// SeoLocationHook is the signature for custom SeoLocation hook methods
	SeoLocationHook func(context.Context, boil.ContextExecutor, *SeoLocation) error

	seoLocationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	seoLocationType                 = reflect.TypeOf(&SeoLocation{})
	seoLocationMapping              = queries.MakeStructMapping(seoLocationType)
	seoLocationPrimaryKeyMapping, _ = queries.BindMapping(seoLocationType, seoLocationMapping, seoLocationPrimaryKeyColumns)
	seoLocationInsertCacheMut       sync.RWMutex
	seoLocationInsertCache          = make(map[string]insertCache)
	seoLocationUpdateCacheMut       sync.RWMutex
	seoLocationUpdateCache          = make(map[string]updateCache)
	seoLocationUpsertCacheMut       sync.RWMutex
	seoLocationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var seoLocationAfterSelectHooks []SeoLocationHook

var seoLocationBeforeInsertHooks []SeoLocationHook
var seoLocationAfterInsertHooks []SeoLocationHook

var seoLocationBeforeUpdateHooks []SeoLocationHook
var seoLocationAfterUpdateHooks []SeoLocationHook

var seoLocationBeforeDeleteHooks []SeoLocationHook
var seoLocationAfterDeleteHooks []SeoLocationHook

var seoLocationBeforeUpsertHooks []SeoLocationHook
var seoLocationAfterUpsertHooks []SeoLocationHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SeoLocation) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SeoLocation) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SeoLocation) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SeoLocation) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SeoLocation) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SeoLocation) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SeoLocation) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SeoLocation) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SeoLocation) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range seoLocationAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSeoLocationHook registers your hook function for all future operations.
func AddSeoLocationHook(hookPoint boil.HookPoint, seoLocationHook SeoLocationHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		seoLocationAfterSelectHooks = append(seoLocationAfterSelectHooks, seoLocationHook)
	case boil.BeforeInsertHook:
		seoLocationBeforeInsertHooks = append(seoLocationBeforeInsertHooks, seoLocationHook)
	case boil.AfterInsertHook:
		seoLocationAfterInsertHooks = append(seoLocationAfterInsertHooks, seoLocationHook)
	case boil.BeforeUpdateHook:
		seoLocationBeforeUpdateHooks = append(seoLocationBeforeUpdateHooks, seoLocationHook)
	case boil.AfterUpdateHook:
		seoLocationAfterUpdateHooks = append(seoLocationAfterUpdateHooks, seoLocationHook)
	case boil.BeforeDeleteHook:
		seoLocationBeforeDeleteHooks = append(seoLocationBeforeDeleteHooks, seoLocationHook)
	case boil.AfterDeleteHook:
		seoLocationAfterDeleteHooks = append(seoLocationAfterDeleteHooks, seoLocationHook)
	case boil.BeforeUpsertHook:
		seoLocationBeforeUpsertHooks = append(seoLocationBeforeUpsertHooks, seoLocationHook)
	case boil.AfterUpsertHook:
		seoLocationAfterUpsertHooks = append(seoLocationAfterUpsertHooks, seoLocationHook)
	}
}

// OneG returns a single seoLocation record from the query using the global executor.
func (q seoLocationQuery) OneG(ctx context.Context) (*SeoLocation, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single seoLocation record from the query.
func (q seoLocationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SeoLocation, error) {
	o := &SeoLocation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for seo_location")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all SeoLocation records from the query using the global executor.
func (q seoLocationQuery) AllG(ctx context.Context) (SeoLocationSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all SeoLocation records from the query.
func (q seoLocationQuery) All(ctx context.Context, exec boil.ContextExecutor) (SeoLocationSlice, error) {
	var o []*SeoLocation

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to SeoLocation slice")
	}

	if len(seoLocationAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all SeoLocation records in the query using the global executor
func (q seoLocationQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all SeoLocation records in the query.
func (q seoLocationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count seo_location rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q seoLocationQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q seoLocationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if seo_location exists")
	}

	return count > 0, nil
}

// SeoLocations retrieves all the records using an executor.
func SeoLocations(mods ...qm.QueryMod) seoLocationQuery {
	mods = append(mods, qm.From("\"seo_location\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"seo_location\".*"})
	}

	return seoLocationQuery{q}
}

// FindSeoLocationG retrieves a single record by ID.
func FindSeoLocationG(ctx context.Context, name string, selectCols ...string) (*SeoLocation, error) {
	return FindSeoLocation(ctx, boil.GetContextDB(), name, selectCols...)
}

// FindSeoLocation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSeoLocation(ctx context.Context, exec boil.ContextExecutor, name string, selectCols ...string) (*SeoLocation, error) {
	seoLocationObj := &SeoLocation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"seo_location\" where \"name\"=$1", sel,
	)

	q := queries.Raw(query, name)

	err := q.Bind(ctx, exec, seoLocationObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from seo_location")
	}

	if err = seoLocationObj.doAfterSelectHooks(ctx, exec); err != nil {
		return seoLocationObj, err
	}

	return seoLocationObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *SeoLocation) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SeoLocation) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no seo_location provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(seoLocationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	seoLocationInsertCacheMut.RLock()
	cache, cached := seoLocationInsertCache[key]
	seoLocationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			seoLocationAllColumns,
			seoLocationColumnsWithDefault,
			seoLocationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(seoLocationType, seoLocationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(seoLocationType, seoLocationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"seo_location\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"seo_location\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "database: unable to insert into seo_location")
	}

	if !cached {
		seoLocationInsertCacheMut.Lock()
		seoLocationInsertCache[key] = cache
		seoLocationInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single SeoLocation record using the global executor.
// See Update for more documentation.
func (o *SeoLocation) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the SeoLocation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SeoLocation) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	seoLocationUpdateCacheMut.RLock()
	cache, cached := seoLocationUpdateCache[key]
	seoLocationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			seoLocationAllColumns,
			seoLocationPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update seo_location, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"seo_location\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, seoLocationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(seoLocationType, seoLocationMapping, append(wl, seoLocationPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "database: unable to update seo_location row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for seo_location")
	}

	if !cached {
		seoLocationUpdateCacheMut.Lock()
		seoLocationUpdateCache[key] = cache
		seoLocationUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q seoLocationQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q seoLocationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for seo_location")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for seo_location")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SeoLocationSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SeoLocationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), seoLocationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"seo_location\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, seoLocationPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in seoLocation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all seoLocation")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *SeoLocation) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SeoLocation) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no seo_location provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(seoLocationColumnsWithDefault, o)

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

	seoLocationUpsertCacheMut.RLock()
	cache, cached := seoLocationUpsertCache[key]
	seoLocationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			seoLocationAllColumns,
			seoLocationColumnsWithDefault,
			seoLocationColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			seoLocationAllColumns,
			seoLocationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert seo_location, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(seoLocationPrimaryKeyColumns))
			copy(conflict, seoLocationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"seo_location\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(seoLocationType, seoLocationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(seoLocationType, seoLocationMapping, ret)
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
		return errors.Wrap(err, "database: unable to upsert seo_location")
	}

	if !cached {
		seoLocationUpsertCacheMut.Lock()
		seoLocationUpsertCache[key] = cache
		seoLocationUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single SeoLocation record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *SeoLocation) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single SeoLocation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SeoLocation) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no SeoLocation provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), seoLocationPrimaryKeyMapping)
	sql := "DELETE FROM \"seo_location\" WHERE \"name\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from seo_location")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for seo_location")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q seoLocationQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q seoLocationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no seoLocationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from seo_location")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for seo_location")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o SeoLocationSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SeoLocationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(seoLocationBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), seoLocationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"seo_location\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, seoLocationPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from seoLocation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for seo_location")
	}

	if len(seoLocationAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *SeoLocation) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no SeoLocation provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SeoLocation) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSeoLocation(ctx, exec, o.Name)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SeoLocationSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty SeoLocationSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SeoLocationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SeoLocationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), seoLocationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"seo_location\".* FROM \"seo_location\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, seoLocationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in SeoLocationSlice")
	}

	*o = slice

	return nil
}

// SeoLocationExistsG checks if the SeoLocation row exists.
func SeoLocationExistsG(ctx context.Context, name string) (bool, error) {
	return SeoLocationExists(ctx, boil.GetContextDB(), name)
}

// SeoLocationExists checks if the SeoLocation row exists.
func SeoLocationExists(ctx context.Context, exec boil.ContextExecutor, name string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"seo_location\" where \"name\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, name)
	}
	row := exec.QueryRowContext(ctx, sql, name)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if seo_location exists")
	}

	return exists, nil
}

// Exists checks if the SeoLocation row exists.
func (o *SeoLocation) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SeoLocationExists(ctx, exec, o.Name)
}