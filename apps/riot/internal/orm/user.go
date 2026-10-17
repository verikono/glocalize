package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// User represents the database model
type User struct {
	Id        string   `db:"id" json:"id"`
	Name      string   `db:"name" json:"name"`
	Preferred []*Skill `json:"preferred,omitempty"`
	Blacklist []*Skill `json:"blacklist,omitempty"`
}

// UserCreate contains model input fields for User creation operations.
//
// Fields for User:
//
//	id   string required
//	name string required
type UserCreate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// colMask returns a bit mask of columns that are set
func (s *UserCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// UserUpdate contains model input fields for User update operations.
type UserUpdate struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

func (u *UserUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, u.Id)
	}
	if u.Name != nil {
		cols = append(cols, "name")
		vals = append(vals, u.Name)
	}
	return cols, vals
}

func assignmentsToUserUpdate(assignments []FieldAssignment) (UserUpdate, error) {
	var input UserUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				errs.ValidateString("id", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Id = v
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "name":
			if v, ok := a.Val.(string); ok {
				input.Name = &v
				errs.ValidateString("name", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Name = v
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// UserSelect specifies which scalar and relation fields to select for User.
//
// Selectable fields:
//
//	-- Scalars --
//	id   (bool)
//	name (bool)
//	-- Relations --
//	preferred ([]Skill)
//	blacklist ([]Skill)
type UserSelect struct {
	Id        bool             `json:"id"`
	Name      bool             `json:"name"`
	Preferred SkillSelectQuery `json:"preferred,omitempty"`
	Blacklist SkillSelectQuery `json:"blacklist,omitempty"`
}

var fullUserSelectVal = &UserSelect{
	Id:   true,
	Name: true,
}

func fullUserSelect() *UserSelect {
	return fullUserSelectVal
}

func (s *UserSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Id || s.Name
}

func (s *UserSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type UserOmit struct {
	Id   bool `json:"id"`
	Name bool `json:"name"`
}

type UserSelectQuery interface {
	GetRelationParams() (*UserSelect, *UserOmit, QueryParams[User])
}

func (s *UserSelect) GetRelationParams() (*UserSelect, *UserOmit, QueryParams[User]) {
	return s, nil, QueryParams[User]{}
}

type UserQueryBuilder struct {
	selects *UserSelect
	omits   *UserOmit
	where   []PredicateOf[User]
	take    *int
	skip    *int
	orderBy []OrderBy[User]
	cursor  UniquePredicate[User]
}

func (b *UserQueryBuilder) Where(preds ...PredicateOf[User]) *UserQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *UserQueryBuilder) Take(limit int) *UserQueryBuilder {
	b.take = &limit
	return b
}

func (b *UserQueryBuilder) Skip(offset int) *UserQueryBuilder {
	b.skip = &offset
	return b
}

func (b *UserQueryBuilder) OrderBy(orders ...OrderBy[User]) *UserQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *UserQueryBuilder) Cursor(where UniquePredicate[User]) *UserQueryBuilder {
	b.cursor = where
	return b
}

func (b *UserQueryBuilder) Select(s UserSelect) *UserQueryBuilder {
	b.selects = &s
	return b
}

func (b *UserQueryBuilder) Omit(o UserOmit) *UserQueryBuilder {
	b.omits = &o
	return b
}

func (b *UserQueryBuilder) GetRelationParams() (*UserSelect, *UserOmit, QueryParams[User]) {
	if b == nil {
		return nil, nil, QueryParams[User]{}
	}
	return b.selects, b.omits, QueryParams[User]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// UserCreateArgs is the input argument passed to User Create extension hooks.
//
// Fields for User:
//
//	id   string required
//	name string required
//
// Relations for User:
//
//	preferred ([]Skill)
//	blacklist ([]Skill)
type UserCreateArgs struct {
	// Data contains the model fields to insert.
	Data *UserCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *UserSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// UserCreateManyArgs is the input argument passed to User CreateMany extension hooks.
//
// Fields for User:
//
//	id   string required
//	name string required
type UserCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*UserCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *UserCreateManyArgs) AppendData(builders ...*UserCreateBuilder) *UserCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToUserCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// UserCreateManyAndReturnArgs is the input argument passed to User CreateManyAndReturn extension hooks.
//
// Fields for User:
//
//	id   string required
//	name string required
//
// Relations for User:
//
//	preferred ([]Skill)
//	blacklist ([]Skill)
type UserCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*UserCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *UserSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *UserCreateManyAndReturnArgs) AppendData(builders ...*UserCreateBuilder) *UserCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToUserCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// UserFindUniqueArgs is the input argument passed to User FindUnique extension hooks.
type UserFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[User]
	// Select specifies which scalar and relation fields to select and return.
	Select *UserSelect
}

func (a *UserFindUniqueArgs) SetWhere(unique UniquePredicate[User], additional ...PredicateOf[User]) *UserFindUniqueArgs {
	a.Where = make([]PredicateOf[User], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// UserFindFirstArgs is the input argument passed to User FindFirst extension hooks.
type UserFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[User]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[User]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[User]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *UserSelect
}

func (a *UserFindFirstArgs) SetWhere(preds ...PredicateOf[User]) *UserFindFirstArgs {
	a.Where = preds
	return a
}

func (a *UserFindFirstArgs) SetOrderBy(orders ...OrderBy[User]) *UserFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *UserFindFirstArgs) SetCursor(cursor UniquePredicate[User]) *UserFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *UserFindFirstArgs) SetSkip(n int) *UserFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *UserFindFirstArgs) SetTake(n int) *UserFindFirstArgs {
	a.Take = &n
	return a
}

// UserFindManyArgs is the input argument passed to User FindMany extension hooks.
type UserFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[User]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[User]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[User]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *UserSelect
}

func (a *UserFindManyArgs) SetWhere(preds ...PredicateOf[User]) *UserFindManyArgs {
	a.Where = preds
	return a
}

func (a *UserFindManyArgs) SetOrderBy(orders ...OrderBy[User]) *UserFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *UserFindManyArgs) SetCursor(cursor UniquePredicate[User]) *UserFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *UserFindManyArgs) SetSkip(n int) *UserFindManyArgs {
	a.Skip = &n
	return a
}

func (a *UserFindManyArgs) SetTake(n int) *UserFindManyArgs {
	a.Take = &n
	return a
}

// UserCountArgs is the input argument passed to User Count extension hooks.
type UserCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[User]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *UserCountArgs) SetWhere(preds ...PredicateOf[User]) *UserCountArgs {
	a.Where = preds
	return a
}

func (a *UserCountArgs) SetSkip(n int) *UserCountArgs {
	a.Skip = &n
	return a
}

func (a *UserCountArgs) SetTake(n int) *UserCountArgs {
	a.Take = &n
	return a
}

// UserDeleteArgs is the input argument passed to User Delete extension hooks.
type UserDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[User]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *UserSelect
}

func (a *UserDeleteArgs) SetWhere(unique UniquePredicate[User], additional ...PredicateOf[User]) *UserDeleteArgs {
	a.Where = make([]PredicateOf[User], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// UserDeleteManyArgs is the input argument passed to User DeleteMany extension hooks.
type UserDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[User]
}

func (a *UserDeleteManyArgs) SetWhere(preds ...PredicateOf[User]) *UserDeleteManyArgs {
	a.Where = preds
	return a
}

// UserUpdateArgs is the input argument passed to User Update extension hooks.
type UserUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[User]
	// Data contains the model fields to update.
	Data *UserUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *UserSelect
}

func (a *UserUpdateArgs) SetWhere(unique UniquePredicate[User], additional ...PredicateOf[User]) *UserUpdateArgs {
	a.Where = make([]PredicateOf[User], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// UserUpdateManyArgs is the input argument passed to User UpdateMany extension hooks.
type UserUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[User]
	// Data contains the model fields to update.
	Data *UserUpdate
}

func (a *UserUpdateManyArgs) SetWhere(preds ...PredicateOf[User]) *UserUpdateManyArgs {
	a.Where = preds
	return a
}

// UserUpdateManyAndReturnArgs is the input argument passed to User UpdateManyAndReturn extension hooks.
type UserUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[User]
	// Data contains the model fields to update.
	Data *UserUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *UserSelect
}

func (a *UserUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[User]) *UserUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type UserCreateQuery = func(ctx context.Context, args *UserCreateArgs) (*User, error)
type UserCreateManyQuery = func(ctx context.Context, args *UserCreateManyArgs) (int64, error)
type UserCreateManyAndReturnQuery = func(ctx context.Context, args *UserCreateManyAndReturnArgs) ([]*User, error)
type UserFindUniqueQuery = func(ctx context.Context, args *UserFindUniqueArgs) (*User, error)
type UserFindFirstQuery = func(ctx context.Context, args *UserFindFirstArgs) (*User, error)
type UserFindManyQuery = func(ctx context.Context, args *UserFindManyArgs) ([]*User, error)
type UserDeleteQuery = func(ctx context.Context, args *UserDeleteArgs) (*User, error)
type UserDeleteManyQuery = func(ctx context.Context, args *UserDeleteManyArgs) (int64, error)
type UserCountQuery = func(ctx context.Context, args *UserCountArgs) (int64, error)
type UserUpdateQuery = func(ctx context.Context, args *UserUpdateArgs) (*User, error)
type UserUpdateManyQuery = func(ctx context.Context, args *UserUpdateManyArgs) (int64, error)
type UserUpdateManyAndReturnQuery = func(ctx context.Context, args *UserUpdateManyAndReturnArgs) ([]*User, error)

type UserExtension struct {
	Create              func(ctx context.Context, args *UserCreateArgs, next UserCreateQuery) (*User, error)
	CreateMany          func(ctx context.Context, args *UserCreateManyArgs, next UserCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *UserCreateManyAndReturnArgs, next UserCreateManyAndReturnQuery) ([]*User, error)
	FindUnique          func(ctx context.Context, args *UserFindUniqueArgs, next UserFindUniqueQuery) (*User, error)
	FindFirst           func(ctx context.Context, args *UserFindFirstArgs, next UserFindFirstQuery) (*User, error)
	FindMany            func(ctx context.Context, args *UserFindManyArgs, next UserFindManyQuery) ([]*User, error)
	Delete              func(ctx context.Context, args *UserDeleteArgs, next UserDeleteQuery) (*User, error)
	DeleteMany          func(ctx context.Context, args *UserDeleteManyArgs, next UserDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *UserCountArgs, next UserCountQuery) (int64, error)
	Update              func(ctx context.Context, args *UserUpdateArgs, next UserUpdateQuery) (*User, error)
	UpdateMany          func(ctx context.Context, args *UserUpdateManyArgs, next UserUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *UserUpdateManyAndReturnArgs, next UserUpdateManyAndReturnQuery) ([]*User, error)
}

type UserDelegate struct {
	client     *Queries
	extensions []UserExtension
}

func (d *UserDelegate) Use(exts ...UserExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *User) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "name":
			targets[i] = &m.Name
		}
	}
	return targets
}

var userDefaultCols = []string{
	"id",
	"name",
}

var userPKCols = []string{
	"id",
}

var userUniqueCols = []string{
	"id",
}

func selectUserCols(selects *UserSelect, omits *UserOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return userDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Name || selects.Preferred != nil || selects.Blacklist != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *UserSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Preferred != nil || s.Blacklist != nil
}

type UserUpsertBuilder struct {
	*CreateBuilder[User, UserSelect, UserOmit]
}

func (b *UserUpsertBuilder) Select(s UserSelect) *UserUpsertBuilder {
	b.selects = &s
	return b
}

func (b *UserUpsertBuilder) Omit(o UserOmit) *UserUpsertBuilder {
	b.omits = &o
	return b
}

type UserCreateBuilder struct {
	*CreateBuilder[User, UserSelect, UserOmit]
}

func (b *UserCreateBuilder) Select(s UserSelect) *UserCreateBuilder {
	b.selects = &s
	return b
}

func (b *UserCreateBuilder) Omit(o UserOmit) *UserCreateBuilder {
	b.omits = &o
	return b
}

func (b *UserCreateBuilder) OnConflict(target UniqueConstraintTarget) *UserConflictBuilder[UserUpsertBuilder] {
	upsertBuilder := &UserUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &UserConflictBuilder[UserUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *UserCreateBuilder) SetId(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *UserCreateBuilder) SetName(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *UserCreateBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *UserDelegate) Create() *UserCreateBuilder {
	return &UserCreateBuilder{
		CreateBuilder: &CreateBuilder[User, UserSelect, UserOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedUserId   uint64 = 1 << 0
	providedUserName uint64 = 1 << 1
)

func assignmentsToUserCreate(assignments []FieldAssignment) (UserCreate, error) {
	var input UserCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedUserId
			if v, ok := a.Val.(string); ok {
				input.Id = v
				errs.ValidateString("id", v, true, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "name":
			provided |= providedUserName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				errs.ValidateString("name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		}
	}
	if provided&providedUserId == 0 {
		errs.Add("id", "", "required", "field Id is required")
	}
	if provided&providedUserName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *UserCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "id")
	if s.Id != "" {
		vals = append(vals, s.Id)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "name")
	vals = append(vals, s.Name)
	return
}

func partitionUserInputs(dialect Dialect, inputs []*UserCreate) [][]*UserCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*UserCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*UserCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*UserCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*UserCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*UserCreate{inputs}
}

func (d *UserDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *UserSelect, omits *UserOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*User, error) {
	input, err := assignmentsToUserCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectUserCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *UserCreateArgs) (*User, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectUserCols(a.Select, omits)
		return d.runCreate(c, cCols, cVals, cReturningCols, a.Select, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Create != nil {
			return ext.Create(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, a *UserCreateArgs) (*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type UserCreateManyBuilder struct {
	*CreateManyBuilder[User]
}

func (b *UserCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *UserConflictBuilder[UserCreateManyBuilder] {
	return &UserConflictBuilder[UserCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type UserCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[User, UserSelect, UserOmit]
}

func (b *UserCreateManyAndReturnBuilder) Select(s UserSelect) *UserCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *UserCreateManyAndReturnBuilder) Omit(o UserOmit) *UserCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *UserCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *UserConflictBuilder[UserCreateManyAndReturnBuilder] {
	return &UserConflictBuilder[UserCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToUserRecordInputs(builders []*UserCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *UserDelegate) CreateMany(builders ...*UserCreateBuilder) *UserCreateManyBuilder {
	return &UserCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[User]{
			records:  createBuildersToUserRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *UserDelegate) CreateManyAndReturn(builders ...*UserCreateBuilder) *UserCreateManyAndReturnBuilder {
	return &UserCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[User, UserSelect, UserOmit]{
			records:  createBuildersToUserRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToUserCreateInputs(records []RecordInput) ([]*UserCreate, error) {
	structs := make([]UserCreate, len(records))
	inputs := make([]*UserCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToUserCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *UserDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToUserCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &UserCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *UserCreateManyArgs) (int64, error) {
		return d.runCreateMany(c, a.Data, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.CreateMany != nil {
			return ext.CreateMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, a *UserCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *UserSelect, omits *UserOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*User, error) {
	inputs, err := recordsToUserCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *UserCreateManyAndReturnArgs) ([]*User, error) {
		return d.runCreateManyAndReturn(c, a.Data, a.Select, omits, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.CreateManyAndReturn != nil {
			return ext.CreateManyAndReturn(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, a *UserCreateManyAndReturnArgs) ([]*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *UserSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*User, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.User.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.User.loadRelations(ctx, []*User{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "User", cols, returningCols, userPKCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	if d.client.dialect.SupportsInsertReturning {
		rows, err := d.client.query(ctx, query, vals...)
		if err != nil {
			return nil, err
		}

		if !rows.Next() {
			err := rows.Err()
			rows.Close()
			if err != nil {
				return nil, TranslateDBError(err)
			}
			return nil, nil
		}

		var res User
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, userPKCols)
}

func (d *UserDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*User, error) {
	result, err := d.client.exec(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	var pkVals []any
	for _, pkCol := range pkCols {
		var val any
		for i, c := range cols {
			if c == pkCol {
				val = vals[i]
				break
			}
		}
		if val == nil && len(pkCols) == 1 {
			lastID, err := result.LastInsertId()
			if err != nil {
				return nil, TranslateDBError(err)
			}
			val = lastID
		}
		pkVals = append(pkVals, val)
	}

	var selectSb strings.Builder
	selectSb.Grow(64 + len(returningCols)*15 + len("User") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "User")
	selectSb.WriteString(" WHERE ")
	for i, pkCol := range pkCols {
		if i > 0 {
			selectSb.WriteString(" AND ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, pkCol)
		selectSb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&selectSb, i+1)
	}

	rows, err := d.client.query(ctx, selectSb.String(), pkVals...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err := rows.Err()
		rows.Close()
		if err != nil {
			return nil, TranslateDBError(err)
		}
		return nil, nil
	}

	var res User
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *UserDelegate) buildBulkInsertSQL(q *Queries, batch []*UserCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range userDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "User")
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WriteQuotedIdent(&sb, col)
	}
	sb.WriteString(") VALUES ")

	paramIdx := paramStartIdx
	for ri, input := range batch {
		if ri > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j, col := range cols {
			if j > 0 {
				sb.WriteString(", ")
			}
			writeDefault := false
			switch col {
			case "id":
				if input.Id != "" {
					vals = append(vals, input.Id)
				} else {
					vals = append(vals, generateCUID())
				}
			case "name":
				vals = append(vals, input.Name)
			}
			if writeDefault {
				sb.WriteString("DEFAULT")
			} else {
				q.dialect.WritePlaceholder(&sb, paramIdx)
				paramIdx++
			}
		}
		sb.WriteString(")")
	}
	queryStr = sb.String()
	return cols, vals, queryStr
}

func applyUserConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
	var conflictCols []string
	if conflictTarget != nil {
		conflictCols = conflictTarget.UniqueColumns()
	}
	var nonConflictCols []string
	if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
		nonConflictCols = computeNonConflictCols(cols, conflictCols, pkCols)
	}
	clause, clauseArgs := dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
	queryStr += clause
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}
	return queryStr, vals
}

func scanUserRows(rows *sql.Rows, returningCols []string) ([]*User, error) {
	var records []*User
	for rows.Next() {
		var res User
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			rows.Close()
			return nil, TranslateDBError(err)
		}
		records = append(records, &res)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, TranslateDBError(err)
	}
	rows.Close()
	return records, nil
}

func (d *UserDelegate) runCreateMany(ctx context.Context, inputs []*UserCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionUserInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyUserConflictClause(d.client.dialect, queryStr, vals, cols, userPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		count += affected
	}
	return count, nil
}

func (d *UserDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*UserCreate,
	selects *UserSelect,
	omits *UserOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*User, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionUserInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.User.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.User.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectUserCols(selects, omits, userPKCols...)
	recordsOut := make([]*User, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyUserConflictClause(d.client.dialect, queryStr, vals, cols, userPKCols, conflictTarget, conflictAction)

		if len(returningCols) > 0 {
			var retSb strings.Builder
			retSb.Grow(12 + len(returningCols)*15)
			retSb.WriteString(" RETURNING ")
			for i, col := range returningCols {
				if i > 0 {
					retSb.WriteString(", ")
				}
				d.client.dialect.WriteQuotedIdent(&retSb, col)
			}
			queryStr += retSb.String()
		}

		rows, err := d.client.query(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		scanned, err := scanUserRows(rows, returningCols)
		if err != nil {
			return nil, err
		}
		recordsOut = append(recordsOut, scanned...)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, recordsOut, selects); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
}

func (d *UserDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*UserCreate,
	selects *UserSelect,
	omits *UserOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*User, error) {
	batches := partitionUserInputs(d.client.dialect, inputs)
	returningCols := selectUserCols(selects, omits, userPKCols...)
	recordsOut := make([]*User, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyUserConflictClause(d.client.dialect, queryStr, vals, cols, userPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("User") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "User")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, userPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, userPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanUserRows(rows, returningCols)
		if err != nil {
			return nil, err
		}
		recordsOut = append(recordsOut, scanned...)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, recordsOut, selects); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
}

type UserConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *UserConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *UserConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *UserConflictBuilder[B]) Update(fn func(u *UserUpsert)) *B {
	cb.setAction(*UserConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// UserConflictUpdate creates a custom ConflictAction for User upsert conflicts.
func UserConflictUpdate(fn func(u *UserUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newUserUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type UserUpsert struct {
	Id   fieldUpsert[string]
	Name fieldUpsert[string]
}

func newUserUpsert(up *ConflictUpdate) *UserUpsert {
	return &UserUpsert{
		Id:   fieldUpsert[string]{column: "id", update: up},
		Name: fieldUpsert[string]{column: "name", update: up},
	}
}

type UserUpdateBuilder struct {
	*UpdateBuilder[User, UserSelect, UserOmit]
}

type UserUpdateManyBuilder struct {
	*UpdateManyBuilder[User]
}

type UserUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[User, UserSelect, UserOmit]
}

func (b *UserUpdateBuilder) SetId(v string) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetId(v string) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetId(v string) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetName(v string) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetName(v string) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetName(v string) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *UserUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *UserUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *UserUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *UserDelegate) Update(where UniquePredicate[User], additional ...PredicateOf[User]) *UserUpdateBuilder {
	return &UserUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[User, UserSelect, UserOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *UserDelegate) UpdateMany(preds ...PredicateOf[User]) *UserUpdateManyBuilder {
	return &UserUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[User]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *UserDelegate) UpdateManyAndReturn(preds ...PredicateOf[User]) *UserUpdateManyAndReturnBuilder {
	return &UserUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[User, UserSelect, UserOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *UserDelegate) buildUpdateSQL(preds []PredicateOf[User], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "User")
	sb.WriteString(" SET ")

	setVals := make([]any, 0, len(cols)+len(predVals))
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
		sb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&sb, i+1)
		setVals = append(setVals, vals[i])
	}

	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
		setVals = append(setVals, predVals...)
	}

	if len(returningCols) > 0 && d.client.dialect.SupportsUpdateReturning {
		sb.WriteString(" RETURNING ")
		for i, col := range returningCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&sb, col)
		}
	}

	return sb.String(), setVals
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

func (d *UserDelegate) executeUpdate(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) (*User, error) {
	allWhere := make([]PredicateOf[User], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToUserUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *UserUpdateArgs) (*User, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.runUpdate(c, a.Where, extCols, extVals, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Update != nil {
			return ext.Update(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Update != nil {
			next, hook := curr, ext.Update
			curr = func(c context.Context, a *UserUpdateArgs) (*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) runUpdate(ctx context.Context, preds []PredicateOf[User], cols []string, vals []any, selects *UserSelect, omits *UserOmit) (*User, error) {
	if len(cols) == 0 {
		return d.runFindUnique(ctx, preds, selects, omits)
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (!d.client.dialect.SupportsUpdateReturning || hasRelations) && !d.client.inTx()

	if useTx {
		var res *User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.User.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.User.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectUserCols(selects, omits, userPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err := rows.Err()
		rows.Close()
		if err != nil {
			return nil, TranslateDBError(err)
		}
		return nil, &NotFoundError{Model: "User"}
	}

	var res User
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*User{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *UserDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[User], cols []string, vals []any) (int64, error) {
	if len(cols) == 0 {
		return 0, nil
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	query, setVals := d.buildUpdateSQL(preds, cols, vals, nil)
	result, err := d.client.exec(ctx, query, setVals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *UserDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[User], cols []string, vals []any, selects *UserSelect, omits *UserOmit) (*User, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "User"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *UserDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToUserUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &UserUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *UserUpdateManyArgs) (int64, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.execUpdateStmt(c, a.Where, extCols, extVals)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.UpdateMany != nil {
			return ext.UpdateMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.UpdateMany != nil {
			next, hook := curr, ext.UpdateMany
			curr = func(c context.Context, a *UserUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *UserDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	input, err := assignmentsToUserUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *UserUpdateManyAndReturnArgs) ([]*User, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.runUpdateManyAndReturn(c, a.Where, extCols, extVals, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.UpdateManyAndReturn != nil {
			return ext.UpdateManyAndReturn(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.UpdateManyAndReturn != nil {
			next, hook := curr, ext.UpdateManyAndReturn
			curr = func(c context.Context, a *UserUpdateManyAndReturnArgs) ([]*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[User], cols []string, vals []any, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[User]{Where: preds}, selects, omits)
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (!d.client.dialect.SupportsUpdateReturning || hasRelations) && !d.client.inTx()

	if useTx {
		var res []*User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.User.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.User.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectUserCols(selects, omits, userPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanUserRows(rows, returningCols)
	if err != nil {
		return nil, err
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, scanned, selects); err != nil {
			return nil, err
		}
	}

	return scanned, nil
}

func (d *UserDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[User], cols []string, vals []any, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*User{}, nil
	}
	return d.runFindMany(ctx, QueryParams[User]{Where: preds}, selects, omits)
}
func (d *UserDelegate) FindUnique(where UniquePredicate[User], additional ...PredicateOf[User]) *FindUniqueBuilder[User, UserSelect, UserOmit] {
	return &FindUniqueBuilder[User, UserSelect, UserOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *UserDelegate) FindFirst(preds ...PredicateOf[User]) *FindFirstBuilder[User, UserSelect, UserOmit] {
	return &FindFirstBuilder[User, UserSelect, UserOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *UserDelegate) FindMany(preds ...PredicateOf[User]) *FindManyBuilder[User, UserSelect, UserOmit] {
	return &FindManyBuilder[User, UserSelect, UserOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *UserDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	allWhere := make([]PredicateOf[User], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *UserFindUniqueArgs) (*User, error) {
		return d.runFindUnique(c, a.Where, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindUnique != nil {
			return ext.FindUnique(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, a *UserFindUniqueArgs) (*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) (*User, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *UserFindFirstArgs) (*User, error) {
		return d.runFindFirst(c, QueryParams[User]{
			Where:   a.Where,
			OrderBy: a.OrderBy,
			Cursor:  a.Cursor,
			Skip:    a.Skip,
			Take:    a.Take,
		}, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindFirst != nil {
			return ext.FindFirst(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, a *UserFindFirstArgs) (*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) ([]*User, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *UserFindManyArgs) ([]*User, error) {
		return d.runFindMany(c, QueryParams[User]{
			Where:   a.Where,
			OrderBy: a.OrderBy,
			Cursor:  a.Cursor,
			Skip:    a.Skip,
			Take:    a.Take,
		}, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindMany != nil {
			return ext.FindMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, a *UserFindManyArgs) ([]*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) runFindUnique(ctx context.Context, where []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, _ := CompilePredicates(d.client.dialect, where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectUserCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*User{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *UserDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) (*User, error) {
	one := 1
	negOne := -1
	if params.Take != nil && *params.Take < 0 {
		params.Take = &negOne
	} else {
		params.Take = &one
	}

	results, err := d.runFindMany(ctx, params, selects, omits)
	if err != nil || len(results) == 0 {
		return nil, err
	}
	return results[0], nil
}

func (d *UserDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) ([]*User, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, nextIdx := CompilePredicates(d.client.dialect, params.Where)
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if isCursorQuery {
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, userPKCols, userUniqueCols, "User", nextIdx, params.Take)
		if err != nil {
			return nil, err
		}
		if cClause != "" {
			if whereClause == "" {
				whereClause = cClause
			} else {
				whereClause = "(" + whereClause + ") AND " + cClause
			}
			vals = append(vals, cVals...)
		}
	}
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, userPKCols, userUniqueCols, isCursorQuery, params.Take)
	returningCols := selectUserCols(selects, omits)

	results, err := d.queryMany(ctx, whereClause, orderByClause, vals, returningCols, params.Take, params.Skip)
	if err != nil || len(results) == 0 {
		return results, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, results, selects); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (d *UserDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*User, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "User", returningCols, whereClause, orderByClause, &limitOne, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
				return nil, nil
			}
			return nil, TranslateDBError(err)
		}
		return nil, nil
	}

	var res User
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *UserDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*User, error) {
	query := buildSelectSQL(d.client, "User", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*User, 0)
	for rows.Next() {
		var res User
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, TranslateDBError(err)
		}
		results = append(results, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, TranslateDBError(err)
	}
	if take != nil && *take < 0 {
		reverseSlice(results)
	}
	return results, nil
}
func (d *UserDelegate) DeleteMany(preds ...PredicateOf[User]) *DeleteManyBuilder[User] {
	return &DeleteManyBuilder[User]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *UserDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[User]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &UserDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *UserDeleteManyArgs) (int64, error) {
		return d.runDeleteMany(c, a.Where)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.DeleteMany != nil {
			return ext.DeleteMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, a *UserDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[User]) (int64, error) {
	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, preds)

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "User")
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	result, err := d.client.exec(ctx, sb.String(), vals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *UserDelegate) Delete(where UniquePredicate[User], additional ...PredicateOf[User]) *DeleteBuilder[User, UserSelect, UserOmit] {
	return &DeleteBuilder[User, UserSelect, UserOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *UserDelegate) executeDelete(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	allWhere := make([]PredicateOf[User], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullUserSelect()
	}

	args := &UserDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *UserDeleteArgs) (*User, error) {
		return d.runDelete(c, a.Where, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Delete != nil {
			return ext.Delete(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, a *UserDeleteArgs) (*User, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) runDelete(ctx context.Context, where []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectUserCols(selects, omits, userPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.User.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "User"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "User")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[User]
			pkPreds = append(pkPreds, Predicate[User]{
				Data: PredicateData{
					Column:   "id",
					Operator: "=",
					Value:    res.Id,
				},
			})

			whereClause, vals, _ := CompilePredicates(txQ.dialect, pkPreds)
			deleteSb.WriteString(whereClause)

			_, err = txQ.exec(ctx, deleteSb.String(), vals...)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// Dialect supports RETURNING, and no relations need loading: run direct DELETE ... RETURNING
	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "User")

	whereClause, vals, _ := CompilePredicates(d.client.dialect, where)
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	sb.WriteString(" RETURNING ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
	}

	rows, err := d.client.query(ctx, sb.String(), vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, TranslateDBError(err)
		}
		return nil, &NotFoundError{Model: "User"}
	}

	var row User
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *UserDelegate) Count(preds ...PredicateOf[User]) *CountBuilder[User] {
	return &CountBuilder[User]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *UserDelegate) executeCount(ctx context.Context, params QueryParams[User]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &UserCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *UserCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[User]{
			Where: a.Where,
			Skip:  a.Skip,
			Take:  a.Take,
		})
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Count != nil {
			return ext.Count(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, a *UserCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *UserDelegate) runCount(ctx context.Context, params QueryParams[User]) (int64, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	var query string
	if params.Take != nil || params.Skip != nil {
		var subQuery strings.Builder
		subQuery.WriteString("SELECT 1 FROM ")
		d.client.dialect.WriteQuotedIdent(&subQuery, "User")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "User")
		if whereClause != "" {
			sb.WriteString(whereClause)
		}
		query = sb.String()
	}

	rows, err := d.client.query(ctx, query, vals...)
	if err != nil {
		return 0, TranslateDBError(err)
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, TranslateDBError(err)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, TranslateDBError(err)
	}
	return count, nil
}
func (d *UserDelegate) loadRelations(ctx context.Context, records []*User, selects *UserSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Preferred != nil {
		// Many-to-many through join table _UserSkills
		relationSelects, relationOmits, relationParams := selects.Preferred.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits)
		allChildren, err := loadManyToMany(
			ctx, d.client, records,
			directKey(func(p *User) string { return p.Id }),
			"Skill",
			"id",
			"_UserSkills",
			"userId",
			"skillId",
			returningCols,
			(*Skill).ScanFields,
			appendMany(func(p *User) *[]*Skill { return &p.Preferred }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading preferred: %w", err)
		}
		if err := d.client.Skill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Blacklist != nil {
		// Many-to-many through join table _UserSkillBlacklist
		relationSelects, relationOmits, relationParams := selects.Blacklist.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits)
		allChildren, err := loadManyToMany(
			ctx, d.client, records,
			directKey(func(p *User) string { return p.Id }),
			"Skill",
			"id",
			"_UserSkillBlacklist",
			"userId",
			"skillId",
			returningCols,
			(*Skill).ScanFields,
			appendMany(func(p *User) *[]*Skill { return &p.Blacklist }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading blacklist: %w", err)
		}
		if err := d.client.Skill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
