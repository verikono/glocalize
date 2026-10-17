package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// SkillDependency represents the database model
type SkillDependency struct {
	DependencyId string `db:"dependencyId" json:"dependencyId"`
	DependentId  string `db:"dependentId" json:"dependentId"`
	Dependency   *Skill `json:"dependency,omitempty"`
	Dependent    *Skill `json:"dependent,omitempty"`
}

// SkillDependencyCreate contains model input fields for SkillDependency creation operations.
//
// Fields for SkillDependency:
//
//	dependencyId string required
//	dependentId  string required
type SkillDependencyCreate struct {
	DependencyId string `json:"dependencyId"`
	DependentId  string `json:"dependentId"`
}

// colMask returns a bit mask of columns that are set
func (s *SkillDependencyCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// SkillDependencyUpdate contains model input fields for SkillDependency update operations.
type SkillDependencyUpdate struct {
	DependencyId *string `json:"dependencyId"`
	DependentId  *string `json:"dependentId"`
}

func (u *SkillDependencyUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.DependencyId != nil {
		cols = append(cols, "dependencyId")
		vals = append(vals, u.DependencyId)
	}
	if u.DependentId != nil {
		cols = append(cols, "dependentId")
		vals = append(vals, u.DependentId)
	}
	return cols, vals
}

func assignmentsToSkillDependencyUpdate(assignments []FieldAssignment) (SkillDependencyUpdate, error) {
	var input SkillDependencyUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "dependencyId":
			if v, ok := a.Val.(string); ok {
				input.DependencyId = &v
				errs.ValidateString("dependencyId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.DependencyId = v
			} else {
				errs.Add("dependencyId", a.Val, "type", "field dependencyId must be of type string")
			}
		case "dependentId":
			if v, ok := a.Val.(string); ok {
				input.DependentId = &v
				errs.ValidateString("dependentId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.DependentId = v
			} else {
				errs.Add("dependentId", a.Val, "type", "field dependentId must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// SkillDependencySelect specifies which scalar and relation fields to select for SkillDependency.
//
// Selectable fields:
//
//	-- Scalars --
//	dependencyId (bool)
//	dependentId  (bool)
//	-- Relations --
//	dependency   (Skill)
//	dependent    (Skill)
type SkillDependencySelect struct {
	DependencyId bool         `json:"dependencyId"`
	DependentId  bool         `json:"dependentId"`
	Dependency   *SkillSelect `json:"dependency,omitempty"`
	Dependent    *SkillSelect `json:"dependent,omitempty"`
}

var fullSkillDependencySelectVal = &SkillDependencySelect{
	DependencyId: true,
	DependentId:  true,
}

func fullSkillDependencySelect() *SkillDependencySelect {
	return fullSkillDependencySelectVal
}

func (s *SkillDependencySelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.DependencyId || s.DependentId
}

func (s *SkillDependencySelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type SkillDependencyOmit struct {
	DependencyId bool `json:"dependencyId"`
	DependentId  bool `json:"dependentId"`
}

type SkillDependencySelectQuery interface {
	GetRelationParams() (*SkillDependencySelect, *SkillDependencyOmit, QueryParams[SkillDependency])
}

func (s *SkillDependencySelect) GetRelationParams() (*SkillDependencySelect, *SkillDependencyOmit, QueryParams[SkillDependency]) {
	return s, nil, QueryParams[SkillDependency]{}
}

type SkillDependencyQueryBuilder struct {
	selects *SkillDependencySelect
	omits   *SkillDependencyOmit
	where   []PredicateOf[SkillDependency]
	take    *int
	skip    *int
	orderBy []OrderBy[SkillDependency]
	cursor  UniquePredicate[SkillDependency]
}

func (b *SkillDependencyQueryBuilder) Where(preds ...PredicateOf[SkillDependency]) *SkillDependencyQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *SkillDependencyQueryBuilder) Take(limit int) *SkillDependencyQueryBuilder {
	b.take = &limit
	return b
}

func (b *SkillDependencyQueryBuilder) Skip(offset int) *SkillDependencyQueryBuilder {
	b.skip = &offset
	return b
}

func (b *SkillDependencyQueryBuilder) OrderBy(orders ...OrderBy[SkillDependency]) *SkillDependencyQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *SkillDependencyQueryBuilder) Cursor(where UniquePredicate[SkillDependency]) *SkillDependencyQueryBuilder {
	b.cursor = where
	return b
}

func (b *SkillDependencyQueryBuilder) Select(s SkillDependencySelect) *SkillDependencyQueryBuilder {
	b.selects = &s
	return b
}

func (b *SkillDependencyQueryBuilder) Omit(o SkillDependencyOmit) *SkillDependencyQueryBuilder {
	b.omits = &o
	return b
}

func (b *SkillDependencyQueryBuilder) GetRelationParams() (*SkillDependencySelect, *SkillDependencyOmit, QueryParams[SkillDependency]) {
	if b == nil {
		return nil, nil, QueryParams[SkillDependency]{}
	}
	return b.selects, b.omits, QueryParams[SkillDependency]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// SkillDependencyCreateArgs is the input argument passed to SkillDependency Create extension hooks.
//
// Fields for SkillDependency:
//
//	dependencyId string required
//	dependentId  string required
//
// Relations for SkillDependency:
//
//	dependency   (Skill)
//	dependent    (Skill)
type SkillDependencyCreateArgs struct {
	// Data contains the model fields to insert.
	Data *SkillDependencyCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *SkillDependencySelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// SkillDependencyCreateManyArgs is the input argument passed to SkillDependency CreateMany extension hooks.
//
// Fields for SkillDependency:
//
//	dependencyId string required
//	dependentId  string required
type SkillDependencyCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*SkillDependencyCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *SkillDependencyCreateManyArgs) AppendData(builders ...*SkillDependencyCreateBuilder) *SkillDependencyCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToSkillDependencyCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// SkillDependencyCreateManyAndReturnArgs is the input argument passed to SkillDependency CreateManyAndReturn extension hooks.
//
// Fields for SkillDependency:
//
//	dependencyId string required
//	dependentId  string required
//
// Relations for SkillDependency:
//
//	dependency   (Skill)
//	dependent    (Skill)
type SkillDependencyCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*SkillDependencyCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *SkillDependencySelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *SkillDependencyCreateManyAndReturnArgs) AppendData(builders ...*SkillDependencyCreateBuilder) *SkillDependencyCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToSkillDependencyCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// SkillDependencyFindUniqueArgs is the input argument passed to SkillDependency FindUnique extension hooks.
type SkillDependencyFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[SkillDependency]
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillDependencySelect
}

func (a *SkillDependencyFindUniqueArgs) SetWhere(unique UniquePredicate[SkillDependency], additional ...PredicateOf[SkillDependency]) *SkillDependencyFindUniqueArgs {
	a.Where = make([]PredicateOf[SkillDependency], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillDependencyFindFirstArgs is the input argument passed to SkillDependency FindFirst extension hooks.
type SkillDependencyFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillDependency]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[SkillDependency]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[SkillDependency]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillDependencySelect
}

func (a *SkillDependencyFindFirstArgs) SetWhere(preds ...PredicateOf[SkillDependency]) *SkillDependencyFindFirstArgs {
	a.Where = preds
	return a
}

func (a *SkillDependencyFindFirstArgs) SetOrderBy(orders ...OrderBy[SkillDependency]) *SkillDependencyFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *SkillDependencyFindFirstArgs) SetCursor(cursor UniquePredicate[SkillDependency]) *SkillDependencyFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *SkillDependencyFindFirstArgs) SetSkip(n int) *SkillDependencyFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *SkillDependencyFindFirstArgs) SetTake(n int) *SkillDependencyFindFirstArgs {
	a.Take = &n
	return a
}

// SkillDependencyFindManyArgs is the input argument passed to SkillDependency FindMany extension hooks.
type SkillDependencyFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillDependency]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[SkillDependency]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[SkillDependency]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillDependencySelect
}

func (a *SkillDependencyFindManyArgs) SetWhere(preds ...PredicateOf[SkillDependency]) *SkillDependencyFindManyArgs {
	a.Where = preds
	return a
}

func (a *SkillDependencyFindManyArgs) SetOrderBy(orders ...OrderBy[SkillDependency]) *SkillDependencyFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *SkillDependencyFindManyArgs) SetCursor(cursor UniquePredicate[SkillDependency]) *SkillDependencyFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *SkillDependencyFindManyArgs) SetSkip(n int) *SkillDependencyFindManyArgs {
	a.Skip = &n
	return a
}

func (a *SkillDependencyFindManyArgs) SetTake(n int) *SkillDependencyFindManyArgs {
	a.Take = &n
	return a
}

// SkillDependencyCountArgs is the input argument passed to SkillDependency Count extension hooks.
type SkillDependencyCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillDependency]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *SkillDependencyCountArgs) SetWhere(preds ...PredicateOf[SkillDependency]) *SkillDependencyCountArgs {
	a.Where = preds
	return a
}

func (a *SkillDependencyCountArgs) SetSkip(n int) *SkillDependencyCountArgs {
	a.Skip = &n
	return a
}

func (a *SkillDependencyCountArgs) SetTake(n int) *SkillDependencyCountArgs {
	a.Take = &n
	return a
}

// SkillDependencyDeleteArgs is the input argument passed to SkillDependency Delete extension hooks.
type SkillDependencyDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[SkillDependency]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *SkillDependencySelect
}

func (a *SkillDependencyDeleteArgs) SetWhere(unique UniquePredicate[SkillDependency], additional ...PredicateOf[SkillDependency]) *SkillDependencyDeleteArgs {
	a.Where = make([]PredicateOf[SkillDependency], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillDependencyDeleteManyArgs is the input argument passed to SkillDependency DeleteMany extension hooks.
type SkillDependencyDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillDependency]
}

func (a *SkillDependencyDeleteManyArgs) SetWhere(preds ...PredicateOf[SkillDependency]) *SkillDependencyDeleteManyArgs {
	a.Where = preds
	return a
}

// SkillDependencyUpdateArgs is the input argument passed to SkillDependency Update extension hooks.
type SkillDependencyUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[SkillDependency]
	// Data contains the model fields to update.
	Data *SkillDependencyUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *SkillDependencySelect
}

func (a *SkillDependencyUpdateArgs) SetWhere(unique UniquePredicate[SkillDependency], additional ...PredicateOf[SkillDependency]) *SkillDependencyUpdateArgs {
	a.Where = make([]PredicateOf[SkillDependency], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillDependencyUpdateManyArgs is the input argument passed to SkillDependency UpdateMany extension hooks.
type SkillDependencyUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillDependency]
	// Data contains the model fields to update.
	Data *SkillDependencyUpdate
}

func (a *SkillDependencyUpdateManyArgs) SetWhere(preds ...PredicateOf[SkillDependency]) *SkillDependencyUpdateManyArgs {
	a.Where = preds
	return a
}

// SkillDependencyUpdateManyAndReturnArgs is the input argument passed to SkillDependency UpdateManyAndReturn extension hooks.
type SkillDependencyUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillDependency]
	// Data contains the model fields to update.
	Data *SkillDependencyUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *SkillDependencySelect
}

func (a *SkillDependencyUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[SkillDependency]) *SkillDependencyUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type SkillDependencyCreateQuery = func(ctx context.Context, args *SkillDependencyCreateArgs) (*SkillDependency, error)
type SkillDependencyCreateManyQuery = func(ctx context.Context, args *SkillDependencyCreateManyArgs) (int64, error)
type SkillDependencyCreateManyAndReturnQuery = func(ctx context.Context, args *SkillDependencyCreateManyAndReturnArgs) ([]*SkillDependency, error)
type SkillDependencyFindUniqueQuery = func(ctx context.Context, args *SkillDependencyFindUniqueArgs) (*SkillDependency, error)
type SkillDependencyFindFirstQuery = func(ctx context.Context, args *SkillDependencyFindFirstArgs) (*SkillDependency, error)
type SkillDependencyFindManyQuery = func(ctx context.Context, args *SkillDependencyFindManyArgs) ([]*SkillDependency, error)
type SkillDependencyDeleteQuery = func(ctx context.Context, args *SkillDependencyDeleteArgs) (*SkillDependency, error)
type SkillDependencyDeleteManyQuery = func(ctx context.Context, args *SkillDependencyDeleteManyArgs) (int64, error)
type SkillDependencyCountQuery = func(ctx context.Context, args *SkillDependencyCountArgs) (int64, error)
type SkillDependencyUpdateQuery = func(ctx context.Context, args *SkillDependencyUpdateArgs) (*SkillDependency, error)
type SkillDependencyUpdateManyQuery = func(ctx context.Context, args *SkillDependencyUpdateManyArgs) (int64, error)
type SkillDependencyUpdateManyAndReturnQuery = func(ctx context.Context, args *SkillDependencyUpdateManyAndReturnArgs) ([]*SkillDependency, error)

type SkillDependencyExtension struct {
	Create              func(ctx context.Context, args *SkillDependencyCreateArgs, next SkillDependencyCreateQuery) (*SkillDependency, error)
	CreateMany          func(ctx context.Context, args *SkillDependencyCreateManyArgs, next SkillDependencyCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *SkillDependencyCreateManyAndReturnArgs, next SkillDependencyCreateManyAndReturnQuery) ([]*SkillDependency, error)
	FindUnique          func(ctx context.Context, args *SkillDependencyFindUniqueArgs, next SkillDependencyFindUniqueQuery) (*SkillDependency, error)
	FindFirst           func(ctx context.Context, args *SkillDependencyFindFirstArgs, next SkillDependencyFindFirstQuery) (*SkillDependency, error)
	FindMany            func(ctx context.Context, args *SkillDependencyFindManyArgs, next SkillDependencyFindManyQuery) ([]*SkillDependency, error)
	Delete              func(ctx context.Context, args *SkillDependencyDeleteArgs, next SkillDependencyDeleteQuery) (*SkillDependency, error)
	DeleteMany          func(ctx context.Context, args *SkillDependencyDeleteManyArgs, next SkillDependencyDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *SkillDependencyCountArgs, next SkillDependencyCountQuery) (int64, error)
	Update              func(ctx context.Context, args *SkillDependencyUpdateArgs, next SkillDependencyUpdateQuery) (*SkillDependency, error)
	UpdateMany          func(ctx context.Context, args *SkillDependencyUpdateManyArgs, next SkillDependencyUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *SkillDependencyUpdateManyAndReturnArgs, next SkillDependencyUpdateManyAndReturnQuery) ([]*SkillDependency, error)
}

type SkillDependencyDelegate struct {
	client     *Queries
	extensions []SkillDependencyExtension
}

func (d *SkillDependencyDelegate) Use(exts ...SkillDependencyExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *SkillDependency) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "dependencyId":
			targets[i] = &m.DependencyId
		case "dependentId":
			targets[i] = &m.DependentId
		}
	}
	return targets
}

var skillDependencyDefaultCols = []string{
	"dependencyId",
	"dependentId",
}

var skillDependencyPKCols = []string{
	"dependencyId",
	"dependentId",
}

var skillDependencyUniqueCols = []string{}

func selectSkillDependencyCols(selects *SkillDependencySelect, omits *SkillDependencyOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return skillDependencyDefaultCols
	}

	anySelected := selects != nil && (selects.DependencyId || selects.DependentId || selects.Dependency != nil || selects.Dependent != nil)

	specs := []colSpec{
		{"dependencyId", selects != nil && selects.DependencyId, omits != nil && omits.DependencyId, selects != nil && selects.Dependency != nil},
		{"dependentId", selects != nil && selects.DependentId, omits != nil && omits.DependentId, selects != nil && selects.Dependent != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *SkillDependencySelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Dependency != nil || s.Dependent != nil
}

type SkillDependencyUpsertBuilder struct {
	*CreateBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]
}

func (b *SkillDependencyUpsertBuilder) Select(s SkillDependencySelect) *SkillDependencyUpsertBuilder {
	b.selects = &s
	return b
}

func (b *SkillDependencyUpsertBuilder) Omit(o SkillDependencyOmit) *SkillDependencyUpsertBuilder {
	b.omits = &o
	return b
}

type SkillDependencyCreateBuilder struct {
	*CreateBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]
}

func (b *SkillDependencyCreateBuilder) Select(s SkillDependencySelect) *SkillDependencyCreateBuilder {
	b.selects = &s
	return b
}

func (b *SkillDependencyCreateBuilder) Omit(o SkillDependencyOmit) *SkillDependencyCreateBuilder {
	b.omits = &o
	return b
}

func (b *SkillDependencyCreateBuilder) OnConflict(target UniqueConstraintTarget) *SkillDependencyConflictBuilder[SkillDependencyUpsertBuilder] {
	upsertBuilder := &SkillDependencyUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &SkillDependencyConflictBuilder[SkillDependencyUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *SkillDependencyCreateBuilder) SetDependencyId(v string) *SkillDependencyCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependencyId", Val: v})
	return b
}
func (b *SkillDependencyCreateBuilder) SetDependentId(v string) *SkillDependencyCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependentId", Val: v})
	return b
}

func (b *SkillDependencyCreateBuilder) Assignments(assignments ...FieldAssignmentOf[SkillDependency]) *SkillDependencyCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *SkillDependencyDelegate) Create() *SkillDependencyCreateBuilder {
	return &SkillDependencyCreateBuilder{
		CreateBuilder: &CreateBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedSkillDependencyDependencyId uint64 = 1 << 0
	providedSkillDependencyDependentId  uint64 = 1 << 1
)

func assignmentsToSkillDependencyCreate(assignments []FieldAssignment) (SkillDependencyCreate, error) {
	var input SkillDependencyCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "dependencyId":
			provided |= providedSkillDependencyDependencyId
			if v, ok := a.Val.(string); ok {
				input.DependencyId = v
				errs.ValidateString("dependencyId", v, true, 0, false, false)
			} else {
				errs.Add("dependencyId", a.Val, "type", "field dependencyId must be of type string")
			}
		case "dependentId":
			provided |= providedSkillDependencyDependentId
			if v, ok := a.Val.(string); ok {
				input.DependentId = v
				errs.ValidateString("dependentId", v, true, 0, false, false)
			} else {
				errs.Add("dependentId", a.Val, "type", "field dependentId must be of type string")
			}
		}
	}
	if provided&providedSkillDependencyDependencyId == 0 {
		errs.Add("dependencyId", "", "required", "field DependencyId is required")
	}
	if provided&providedSkillDependencyDependentId == 0 {
		errs.Add("dependentId", "", "required", "field DependentId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *SkillDependencyCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "dependencyId")
	vals = append(vals, s.DependencyId)
	cols = append(cols, "dependentId")
	vals = append(vals, s.DependentId)
	return
}

func partitionSkillDependencyInputs(dialect Dialect, inputs []*SkillDependencyCreate) [][]*SkillDependencyCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*SkillDependencyCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*SkillDependencyCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*SkillDependencyCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*SkillDependencyCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*SkillDependencyCreate{inputs}
}

func (d *SkillDependencyDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *SkillDependencySelect, omits *SkillDependencyOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*SkillDependency, error) {
	input, err := assignmentsToSkillDependencyCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectSkillDependencyCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillDependencyCreateArgs) (*SkillDependency, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectSkillDependencyCols(a.Select, omits)
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
			curr = func(c context.Context, a *SkillDependencyCreateArgs) (*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type SkillDependencyCreateManyBuilder struct {
	*CreateManyBuilder[SkillDependency]
}

func (b *SkillDependencyCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *SkillDependencyConflictBuilder[SkillDependencyCreateManyBuilder] {
	return &SkillDependencyConflictBuilder[SkillDependencyCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type SkillDependencyCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]
}

func (b *SkillDependencyCreateManyAndReturnBuilder) Select(s SkillDependencySelect) *SkillDependencyCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *SkillDependencyCreateManyAndReturnBuilder) Omit(o SkillDependencyOmit) *SkillDependencyCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *SkillDependencyCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *SkillDependencyConflictBuilder[SkillDependencyCreateManyAndReturnBuilder] {
	return &SkillDependencyConflictBuilder[SkillDependencyCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToSkillDependencyRecordInputs(builders []*SkillDependencyCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *SkillDependencyDelegate) CreateMany(builders ...*SkillDependencyCreateBuilder) *SkillDependencyCreateManyBuilder {
	return &SkillDependencyCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[SkillDependency]{
			records:  createBuildersToSkillDependencyRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *SkillDependencyDelegate) CreateManyAndReturn(builders ...*SkillDependencyCreateBuilder) *SkillDependencyCreateManyAndReturnBuilder {
	return &SkillDependencyCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
			records:  createBuildersToSkillDependencyRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToSkillDependencyCreateInputs(records []RecordInput) ([]*SkillDependencyCreate, error) {
	structs := make([]SkillDependencyCreate, len(records))
	inputs := make([]*SkillDependencyCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToSkillDependencyCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *SkillDependencyDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToSkillDependencyCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &SkillDependencyCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillDependencyCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillDependencyCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *SkillDependencySelect, omits *SkillDependencyOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*SkillDependency, error) {
	inputs, err := recordsToSkillDependencyCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillDependencyCreateManyAndReturnArgs) ([]*SkillDependency, error) {
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
			curr = func(c context.Context, a *SkillDependencyCreateManyAndReturnArgs) ([]*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *SkillDependencySelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*SkillDependency, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *SkillDependency
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.SkillDependency.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.SkillDependency.loadRelations(ctx, []*SkillDependency{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "SkillDependency", cols, returningCols, skillDependencyPKCols, conflictTarget, conflictAction, len(vals))
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

		var res SkillDependency
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, skillDependencyPKCols)
}

func (d *SkillDependencyDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*SkillDependency, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("SkillDependency") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "SkillDependency")
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

	var res SkillDependency
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *SkillDependencyDelegate) buildBulkInsertSQL(q *Queries, batch []*SkillDependencyCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range skillDependencyDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "SkillDependency")
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
			case "dependencyId":
				vals = append(vals, input.DependencyId)
			case "dependentId":
				vals = append(vals, input.DependentId)
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

func applySkillDependencyConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanSkillDependencyRows(rows *sql.Rows, returningCols []string) ([]*SkillDependency, error) {
	var records []*SkillDependency
	for rows.Next() {
		var res SkillDependency
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

func (d *SkillDependencyDelegate) runCreateMany(ctx context.Context, inputs []*SkillDependencyCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionSkillDependencyInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillDependencyConflictClause(d.client.dialect, queryStr, vals, cols, skillDependencyPKCols, conflictTarget, conflictAction)

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

func (d *SkillDependencyDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*SkillDependencyCreate,
	selects *SkillDependencySelect,
	omits *SkillDependencyOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*SkillDependency, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionSkillDependencyInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*SkillDependency
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.SkillDependency.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.SkillDependency.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectSkillDependencyCols(selects, omits, skillDependencyPKCols...)
	recordsOut := make([]*SkillDependency, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillDependencyConflictClause(d.client.dialect, queryStr, vals, cols, skillDependencyPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanSkillDependencyRows(rows, returningCols)
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

func (d *SkillDependencyDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*SkillDependencyCreate,
	selects *SkillDependencySelect,
	omits *SkillDependencyOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*SkillDependency, error) {
	batches := partitionSkillDependencyInputs(d.client.dialect, inputs)
	returningCols := selectSkillDependencyCols(selects, omits, skillDependencyPKCols...)
	recordsOut := make([]*SkillDependency, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillDependencyConflictClause(d.client.dialect, queryStr, vals, cols, skillDependencyPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("SkillDependency") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "SkillDependency")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, skillDependencyPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, skillDependencyPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanSkillDependencyRows(rows, returningCols)
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

type SkillDependencyConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *SkillDependencyConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *SkillDependencyConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *SkillDependencyConflictBuilder[B]) Update(fn func(u *SkillDependencyUpsert)) *B {
	cb.setAction(*SkillDependencyConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// SkillDependencyConflictUpdate creates a custom ConflictAction for SkillDependency upsert conflicts.
func SkillDependencyConflictUpdate(fn func(u *SkillDependencyUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newSkillDependencyUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type SkillDependencyUpsert struct {
	DependencyId fieldUpsert[string]
	DependentId  fieldUpsert[string]
}

func newSkillDependencyUpsert(up *ConflictUpdate) *SkillDependencyUpsert {
	return &SkillDependencyUpsert{
		DependencyId: fieldUpsert[string]{column: "dependencyId", update: up},
		DependentId:  fieldUpsert[string]{column: "dependentId", update: up},
	}
}

type SkillDependencyUpdateBuilder struct {
	*UpdateBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]
}

type SkillDependencyUpdateManyBuilder struct {
	*UpdateManyBuilder[SkillDependency]
}

type SkillDependencyUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]
}

func (b *SkillDependencyUpdateBuilder) SetDependencyId(v string) *SkillDependencyUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependencyId", Val: v})
	return b
}

func (b *SkillDependencyUpdateManyBuilder) SetDependencyId(v string) *SkillDependencyUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependencyId", Val: v})
	return b
}

func (b *SkillDependencyUpdateManyAndReturnBuilder) SetDependencyId(v string) *SkillDependencyUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependencyId", Val: v})
	return b
}
func (b *SkillDependencyUpdateBuilder) SetDependentId(v string) *SkillDependencyUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependentId", Val: v})
	return b
}

func (b *SkillDependencyUpdateManyBuilder) SetDependentId(v string) *SkillDependencyUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependentId", Val: v})
	return b
}

func (b *SkillDependencyUpdateManyAndReturnBuilder) SetDependentId(v string) *SkillDependencyUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dependentId", Val: v})
	return b
}

func (b *SkillDependencyUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[SkillDependency]) *SkillDependencyUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *SkillDependencyUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[SkillDependency]) *SkillDependencyUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *SkillDependencyUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[SkillDependency]) *SkillDependencyUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *SkillDependencyDelegate) Update(where UniquePredicate[SkillDependency], additional ...PredicateOf[SkillDependency]) *SkillDependencyUpdateBuilder {
	return &SkillDependencyUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *SkillDependencyDelegate) UpdateMany(preds ...PredicateOf[SkillDependency]) *SkillDependencyUpdateManyBuilder {
	return &SkillDependencyUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[SkillDependency]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *SkillDependencyDelegate) UpdateManyAndReturn(preds ...PredicateOf[SkillDependency]) *SkillDependencyUpdateManyAndReturnBuilder {
	return &SkillDependencyUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *SkillDependencyDelegate) buildUpdateSQL(preds []PredicateOf[SkillDependency], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "SkillDependency")
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

func (d *SkillDependencyDelegate) executeUpdate(ctx context.Context, where UniquePredicate[SkillDependency], additional []PredicateOf[SkillDependency], assignments []FieldAssignment, selects *SkillDependencySelect, omits *SkillDependencyOmit) (*SkillDependency, error) {
	allWhere := make([]PredicateOf[SkillDependency], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToSkillDependencyUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillDependencyUpdateArgs) (*SkillDependency, error) {
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
			curr = func(c context.Context, a *SkillDependencyUpdateArgs) (*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) runUpdate(ctx context.Context, preds []PredicateOf[SkillDependency], cols []string, vals []any, selects *SkillDependencySelect, omits *SkillDependencyOmit) (*SkillDependency, error) {
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
		var res *SkillDependency
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.SkillDependency.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.SkillDependency.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectSkillDependencyCols(selects, omits, skillDependencyPKCols...)
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
		return nil, &NotFoundError{Model: "SkillDependency"}
	}

	var res SkillDependency
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*SkillDependency{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *SkillDependencyDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[SkillDependency], cols []string, vals []any) (int64, error) {
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

func (d *SkillDependencyDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[SkillDependency], cols []string, vals []any, selects *SkillDependencySelect, omits *SkillDependencyOmit) (*SkillDependency, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "SkillDependency"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *SkillDependencyDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[SkillDependency], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToSkillDependencyUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &SkillDependencyUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *SkillDependencyUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillDependencyUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *SkillDependencyDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[SkillDependency], assignments []FieldAssignment, selects *SkillDependencySelect, omits *SkillDependencyOmit) ([]*SkillDependency, error) {
	input, err := assignmentsToSkillDependencyUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillDependencyUpdateManyAndReturnArgs) ([]*SkillDependency, error) {
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
			curr = func(c context.Context, a *SkillDependencyUpdateManyAndReturnArgs) ([]*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[SkillDependency], cols []string, vals []any, selects *SkillDependencySelect, omits *SkillDependencyOmit) ([]*SkillDependency, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[SkillDependency]{Where: preds}, selects, omits)
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
		var res []*SkillDependency
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.SkillDependency.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.SkillDependency.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectSkillDependencyCols(selects, omits, skillDependencyPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanSkillDependencyRows(rows, returningCols)
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

func (d *SkillDependencyDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[SkillDependency], cols []string, vals []any, selects *SkillDependencySelect, omits *SkillDependencyOmit) ([]*SkillDependency, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*SkillDependency{}, nil
	}
	return d.runFindMany(ctx, QueryParams[SkillDependency]{Where: preds}, selects, omits)
}
func (d *SkillDependencyDelegate) FindUnique(where UniquePredicate[SkillDependency], additional ...PredicateOf[SkillDependency]) *FindUniqueBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit] {
	return &FindUniqueBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *SkillDependencyDelegate) FindFirst(preds ...PredicateOf[SkillDependency]) *FindFirstBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit] {
	return &FindFirstBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *SkillDependencyDelegate) FindMany(preds ...PredicateOf[SkillDependency]) *FindManyBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit] {
	return &FindManyBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *SkillDependencyDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[SkillDependency], additional []PredicateOf[SkillDependency], selects *SkillDependencySelect, omits *SkillDependencyOmit) (*SkillDependency, error) {
	allWhere := make([]PredicateOf[SkillDependency], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillDependencyFindUniqueArgs) (*SkillDependency, error) {
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
			curr = func(c context.Context, a *SkillDependencyFindUniqueArgs) (*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[SkillDependency],
	selects *SkillDependencySelect,
	omits *SkillDependencyOmit,
) (*SkillDependency, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *SkillDependencyFindFirstArgs) (*SkillDependency, error) {
		return d.runFindFirst(c, QueryParams[SkillDependency]{
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
			curr = func(c context.Context, a *SkillDependencyFindFirstArgs) (*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[SkillDependency],
	selects *SkillDependencySelect,
	omits *SkillDependencyOmit,
) ([]*SkillDependency, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *SkillDependencyFindManyArgs) ([]*SkillDependency, error) {
		return d.runFindMany(c, QueryParams[SkillDependency]{
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
			curr = func(c context.Context, a *SkillDependencyFindManyArgs) ([]*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) runFindUnique(ctx context.Context, where []PredicateOf[SkillDependency], selects *SkillDependencySelect, omits *SkillDependencyOmit) (*SkillDependency, error) {
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
	returningCols := selectSkillDependencyCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*SkillDependency{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *SkillDependencyDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[SkillDependency],
	selects *SkillDependencySelect,
	omits *SkillDependencyOmit,
) (*SkillDependency, error) {
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

func (d *SkillDependencyDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[SkillDependency],
	selects *SkillDependencySelect,
	omits *SkillDependencyOmit,
) ([]*SkillDependency, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, skillDependencyPKCols, skillDependencyUniqueCols, "SkillDependency", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, skillDependencyPKCols, skillDependencyUniqueCols, isCursorQuery, params.Take)
	returningCols := selectSkillDependencyCols(selects, omits)

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

func (d *SkillDependencyDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*SkillDependency, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "SkillDependency", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res SkillDependency
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *SkillDependencyDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*SkillDependency, error) {
	query := buildSelectSQL(d.client, "SkillDependency", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*SkillDependency, 0)
	for rows.Next() {
		var res SkillDependency
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
func (d *SkillDependencyDelegate) DeleteMany(preds ...PredicateOf[SkillDependency]) *DeleteManyBuilder[SkillDependency] {
	return &DeleteManyBuilder[SkillDependency]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *SkillDependencyDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[SkillDependency]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &SkillDependencyDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *SkillDependencyDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillDependencyDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[SkillDependency]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "SkillDependency")
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

func (d *SkillDependencyDelegate) Delete(where UniquePredicate[SkillDependency], additional ...PredicateOf[SkillDependency]) *DeleteBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit] {
	return &DeleteBuilder[SkillDependency, SkillDependencySelect, SkillDependencyOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *SkillDependencyDelegate) executeDelete(ctx context.Context, where UniquePredicate[SkillDependency], additional []PredicateOf[SkillDependency], selects *SkillDependencySelect, omits *SkillDependencyOmit) (*SkillDependency, error) {
	allWhere := make([]PredicateOf[SkillDependency], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillDependencySelect()
	}

	args := &SkillDependencyDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillDependencyDeleteArgs) (*SkillDependency, error) {
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
			curr = func(c context.Context, a *SkillDependencyDeleteArgs) (*SkillDependency, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) runDelete(ctx context.Context, where []PredicateOf[SkillDependency], selects *SkillDependencySelect, omits *SkillDependencyOmit) (*SkillDependency, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectSkillDependencyCols(selects, omits, skillDependencyPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *SkillDependency
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.SkillDependency.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "SkillDependency"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "SkillDependency")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[SkillDependency]
			pkPreds = append(pkPreds, Predicate[SkillDependency]{
				Data: PredicateData{
					Column:   "dependencyId",
					Operator: "=",
					Value:    res.DependencyId,
				},
			})
			pkPreds = append(pkPreds, Predicate[SkillDependency]{
				Data: PredicateData{
					Column:   "dependentId",
					Operator: "=",
					Value:    res.DependentId,
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
	d.client.dialect.WriteQuotedIdent(&sb, "SkillDependency")

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
		return nil, &NotFoundError{Model: "SkillDependency"}
	}

	var row SkillDependency
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *SkillDependencyDelegate) Count(preds ...PredicateOf[SkillDependency]) *CountBuilder[SkillDependency] {
	return &CountBuilder[SkillDependency]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *SkillDependencyDelegate) executeCount(ctx context.Context, params QueryParams[SkillDependency]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &SkillDependencyCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *SkillDependencyCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[SkillDependency]{
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
			curr = func(c context.Context, a *SkillDependencyCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillDependencyDelegate) runCount(ctx context.Context, params QueryParams[SkillDependency]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "SkillDependency")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "SkillDependency")
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
func (d *SkillDependencyDelegate) loadRelations(ctx context.Context, records []*SkillDependency, selects *SkillDependencySelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Dependency != nil {
		relationSelects, relationOmits, relationParams := selects.Dependency.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: SkillDependency.dependencyId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *SkillDependency) string { return p.DependencyId }),
			"Skill",
			"id",
			returningCols,
			scanInto(returningCols, (*Skill).ScanFields),
			directKey(func(c *Skill) string { return c.Id }),
			setOne(func(p *SkillDependency, c *Skill) { p.Dependency = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading dependency: %w", err)
		}
		if err := d.client.Skill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Dependent != nil {
		relationSelects, relationOmits, relationParams := selects.Dependent.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: SkillDependency.dependentId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *SkillDependency) string { return p.DependentId }),
			"Skill",
			"id",
			returningCols,
			scanInto(returningCols, (*Skill).ScanFields),
			directKey(func(c *Skill) string { return c.Id }),
			setOne(func(p *SkillDependency, c *Skill) { p.Dependent = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading dependent: %w", err)
		}
		if err := d.client.Skill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
