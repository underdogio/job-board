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

// UserSignOnToken is an object representing the database table.
type UserSignOnToken struct {
	Token    string `boil:"token" json:"token" toml:"token" yaml:"token"`
	Email    string `boil:"email" json:"email" toml:"email" yaml:"email"`
	UserType string `boil:"user_type" json:"user_type" toml:"user_type" yaml:"user_type"`

	R *userSignOnTokenR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userSignOnTokenL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserSignOnTokenColumns = struct {
	Token    string
	Email    string
	UserType string
}{
	Token:    "token",
	Email:    "email",
	UserType: "user_type",
}

var UserSignOnTokenTableColumns = struct {
	Token    string
	Email    string
	UserType string
}{
	Token:    "user_sign_on_token.token",
	Email:    "user_sign_on_token.email",
	UserType: "user_sign_on_token.user_type",
}

// Generated where

var UserSignOnTokenWhere = struct {
	Token    whereHelperstring
	Email    whereHelperstring
	UserType whereHelperstring
}{
	Token:    whereHelperstring{field: "\"user_sign_on_token\".\"token\""},
	Email:    whereHelperstring{field: "\"user_sign_on_token\".\"email\""},
	UserType: whereHelperstring{field: "\"user_sign_on_token\".\"user_type\""},
}

// UserSignOnTokenRels is where relationship names are stored.
var UserSignOnTokenRels = struct {
}{}

// userSignOnTokenR is where relationships are stored.
type userSignOnTokenR struct {
}

// NewStruct creates a new relationship struct
func (*userSignOnTokenR) NewStruct() *userSignOnTokenR {
	return &userSignOnTokenR{}
}

// userSignOnTokenL is where Load methods for each relationship are stored.
type userSignOnTokenL struct{}

var (
	userSignOnTokenAllColumns            = []string{"token", "email", "user_type"}
	userSignOnTokenColumnsWithoutDefault = []string{"token", "email"}
	userSignOnTokenColumnsWithDefault    = []string{"user_type"}
	userSignOnTokenPrimaryKeyColumns     = []string{"token"}
	userSignOnTokenGeneratedColumns      = []string{}
)

type (
	// UserSignOnTokenSlice is an alias for a slice of pointers to UserSignOnToken.
	// This should almost always be used instead of []UserSignOnToken.
	UserSignOnTokenSlice []*UserSignOnToken
	// UserSignOnTokenHook is the signature for custom UserSignOnToken hook methods
	UserSignOnTokenHook func(context.Context, boil.ContextExecutor, *UserSignOnToken) error

	userSignOnTokenQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userSignOnTokenType                 = reflect.TypeOf(&UserSignOnToken{})
	userSignOnTokenMapping              = queries.MakeStructMapping(userSignOnTokenType)
	userSignOnTokenPrimaryKeyMapping, _ = queries.BindMapping(userSignOnTokenType, userSignOnTokenMapping, userSignOnTokenPrimaryKeyColumns)
	userSignOnTokenInsertCacheMut       sync.RWMutex
	userSignOnTokenInsertCache          = make(map[string]insertCache)
	userSignOnTokenUpdateCacheMut       sync.RWMutex
	userSignOnTokenUpdateCache          = make(map[string]updateCache)
	userSignOnTokenUpsertCacheMut       sync.RWMutex
	userSignOnTokenUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userSignOnTokenAfterSelectHooks []UserSignOnTokenHook

var userSignOnTokenBeforeInsertHooks []UserSignOnTokenHook
var userSignOnTokenAfterInsertHooks []UserSignOnTokenHook

var userSignOnTokenBeforeUpdateHooks []UserSignOnTokenHook
var userSignOnTokenAfterUpdateHooks []UserSignOnTokenHook

var userSignOnTokenBeforeDeleteHooks []UserSignOnTokenHook
var userSignOnTokenAfterDeleteHooks []UserSignOnTokenHook

var userSignOnTokenBeforeUpsertHooks []UserSignOnTokenHook
var userSignOnTokenAfterUpsertHooks []UserSignOnTokenHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserSignOnToken) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserSignOnToken) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserSignOnToken) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserSignOnToken) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserSignOnToken) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserSignOnToken) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserSignOnToken) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserSignOnToken) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserSignOnToken) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSignOnTokenAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserSignOnTokenHook registers your hook function for all future operations.
func AddUserSignOnTokenHook(hookPoint boil.HookPoint, userSignOnTokenHook UserSignOnTokenHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		userSignOnTokenAfterSelectHooks = append(userSignOnTokenAfterSelectHooks, userSignOnTokenHook)
	case boil.BeforeInsertHook:
		userSignOnTokenBeforeInsertHooks = append(userSignOnTokenBeforeInsertHooks, userSignOnTokenHook)
	case boil.AfterInsertHook:
		userSignOnTokenAfterInsertHooks = append(userSignOnTokenAfterInsertHooks, userSignOnTokenHook)
	case boil.BeforeUpdateHook:
		userSignOnTokenBeforeUpdateHooks = append(userSignOnTokenBeforeUpdateHooks, userSignOnTokenHook)
	case boil.AfterUpdateHook:
		userSignOnTokenAfterUpdateHooks = append(userSignOnTokenAfterUpdateHooks, userSignOnTokenHook)
	case boil.BeforeDeleteHook:
		userSignOnTokenBeforeDeleteHooks = append(userSignOnTokenBeforeDeleteHooks, userSignOnTokenHook)
	case boil.AfterDeleteHook:
		userSignOnTokenAfterDeleteHooks = append(userSignOnTokenAfterDeleteHooks, userSignOnTokenHook)
	case boil.BeforeUpsertHook:
		userSignOnTokenBeforeUpsertHooks = append(userSignOnTokenBeforeUpsertHooks, userSignOnTokenHook)
	case boil.AfterUpsertHook:
		userSignOnTokenAfterUpsertHooks = append(userSignOnTokenAfterUpsertHooks, userSignOnTokenHook)
	}
}

// OneG returns a single userSignOnToken record from the query using the global executor.
func (q userSignOnTokenQuery) OneG(ctx context.Context) (*UserSignOnToken, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single userSignOnToken record from the query.
func (q userSignOnTokenQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UserSignOnToken, error) {
	o := &UserSignOnToken{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for user_sign_on_token")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all UserSignOnToken records from the query using the global executor.
func (q userSignOnTokenQuery) AllG(ctx context.Context) (UserSignOnTokenSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all UserSignOnToken records from the query.
func (q userSignOnTokenQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserSignOnTokenSlice, error) {
	var o []*UserSignOnToken

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to UserSignOnToken slice")
	}

	if len(userSignOnTokenAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all UserSignOnToken records in the query using the global executor
func (q userSignOnTokenQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all UserSignOnToken records in the query.
func (q userSignOnTokenQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count user_sign_on_token rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q userSignOnTokenQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q userSignOnTokenQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if user_sign_on_token exists")
	}

	return count > 0, nil
}

// UserSignOnTokens retrieves all the records using an executor.
func UserSignOnTokens(mods ...qm.QueryMod) userSignOnTokenQuery {
	mods = append(mods, qm.From("\"user_sign_on_token\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"user_sign_on_token\".*"})
	}

	return userSignOnTokenQuery{q}
}

// FindUserSignOnTokenG retrieves a single record by ID.
func FindUserSignOnTokenG(ctx context.Context, token string, selectCols ...string) (*UserSignOnToken, error) {
	return FindUserSignOnToken(ctx, boil.GetContextDB(), token, selectCols...)
}

// FindUserSignOnToken retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserSignOnToken(ctx context.Context, exec boil.ContextExecutor, token string, selectCols ...string) (*UserSignOnToken, error) {
	userSignOnTokenObj := &UserSignOnToken{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_sign_on_token\" where \"token\"=$1", sel,
	)

	q := queries.Raw(query, token)

	err := q.Bind(ctx, exec, userSignOnTokenObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from user_sign_on_token")
	}

	if err = userSignOnTokenObj.doAfterSelectHooks(ctx, exec); err != nil {
		return userSignOnTokenObj, err
	}

	return userSignOnTokenObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *UserSignOnToken) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserSignOnToken) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no user_sign_on_token provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userSignOnTokenColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userSignOnTokenInsertCacheMut.RLock()
	cache, cached := userSignOnTokenInsertCache[key]
	userSignOnTokenInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userSignOnTokenAllColumns,
			userSignOnTokenColumnsWithDefault,
			userSignOnTokenColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userSignOnTokenType, userSignOnTokenMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userSignOnTokenType, userSignOnTokenMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_sign_on_token\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_sign_on_token\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "database: unable to insert into user_sign_on_token")
	}

	if !cached {
		userSignOnTokenInsertCacheMut.Lock()
		userSignOnTokenInsertCache[key] = cache
		userSignOnTokenInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single UserSignOnToken record using the global executor.
// See Update for more documentation.
func (o *UserSignOnToken) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the UserSignOnToken.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserSignOnToken) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userSignOnTokenUpdateCacheMut.RLock()
	cache, cached := userSignOnTokenUpdateCache[key]
	userSignOnTokenUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userSignOnTokenAllColumns,
			userSignOnTokenPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update user_sign_on_token, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_sign_on_token\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userSignOnTokenPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userSignOnTokenType, userSignOnTokenMapping, append(wl, userSignOnTokenPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "database: unable to update user_sign_on_token row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for user_sign_on_token")
	}

	if !cached {
		userSignOnTokenUpdateCacheMut.Lock()
		userSignOnTokenUpdateCache[key] = cache
		userSignOnTokenUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q userSignOnTokenQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q userSignOnTokenQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for user_sign_on_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for user_sign_on_token")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o UserSignOnTokenSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserSignOnTokenSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSignOnTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_sign_on_token\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userSignOnTokenPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in userSignOnToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all userSignOnToken")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *UserSignOnToken) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserSignOnToken) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no user_sign_on_token provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userSignOnTokenColumnsWithDefault, o)

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

	userSignOnTokenUpsertCacheMut.RLock()
	cache, cached := userSignOnTokenUpsertCache[key]
	userSignOnTokenUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userSignOnTokenAllColumns,
			userSignOnTokenColumnsWithDefault,
			userSignOnTokenColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			userSignOnTokenAllColumns,
			userSignOnTokenPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert user_sign_on_token, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userSignOnTokenPrimaryKeyColumns))
			copy(conflict, userSignOnTokenPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_sign_on_token\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userSignOnTokenType, userSignOnTokenMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userSignOnTokenType, userSignOnTokenMapping, ret)
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
		return errors.Wrap(err, "database: unable to upsert user_sign_on_token")
	}

	if !cached {
		userSignOnTokenUpsertCacheMut.Lock()
		userSignOnTokenUpsertCache[key] = cache
		userSignOnTokenUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single UserSignOnToken record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *UserSignOnToken) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single UserSignOnToken record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserSignOnToken) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no UserSignOnToken provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userSignOnTokenPrimaryKeyMapping)
	sql := "DELETE FROM \"user_sign_on_token\" WHERE \"token\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from user_sign_on_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for user_sign_on_token")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q userSignOnTokenQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q userSignOnTokenQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no userSignOnTokenQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from user_sign_on_token")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for user_sign_on_token")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o UserSignOnTokenSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserSignOnTokenSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userSignOnTokenBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSignOnTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_sign_on_token\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userSignOnTokenPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from userSignOnToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for user_sign_on_token")
	}

	if len(userSignOnTokenAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *UserSignOnToken) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no UserSignOnToken provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserSignOnToken) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUserSignOnToken(ctx, exec, o.Token)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSignOnTokenSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty UserSignOnTokenSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSignOnTokenSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserSignOnTokenSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSignOnTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_sign_on_token\".* FROM \"user_sign_on_token\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userSignOnTokenPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in UserSignOnTokenSlice")
	}

	*o = slice

	return nil
}

// UserSignOnTokenExistsG checks if the UserSignOnToken row exists.
func UserSignOnTokenExistsG(ctx context.Context, token string) (bool, error) {
	return UserSignOnTokenExists(ctx, boil.GetContextDB(), token)
}

// UserSignOnTokenExists checks if the UserSignOnToken row exists.
func UserSignOnTokenExists(ctx context.Context, exec boil.ContextExecutor, token string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_sign_on_token\" where \"token\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, token)
	}
	row := exec.QueryRowContext(ctx, sql, token)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if user_sign_on_token exists")
	}

	return exists, nil
}

// Exists checks if the UserSignOnToken row exists.
func (o *UserSignOnToken) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return UserSignOnTokenExists(ctx, exec, o.Token)
}