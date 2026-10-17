package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// TrashSkill represents the database model
type TrashSkill struct {
	TrashId int32  `db:"trashId" json:"trashId"`
	SkillId string `db:"skillId" json:"skillId"`
	Trash   *Trash `json:"trash,omitempty"`
	Skill   *Skill `json:"skill,omitempty"`
}

// TrashSkillCreate contains model input fields for TrashSkill creation operations.
//
// Fields for TrashSkill:
//
//	trashId int32  required
//	skillId string required
type TrashSkillCreate struct {
	TrashId int32  `json:"trashId"`
	SkillId string `json:"skillId"`
}

// colMask returns a bit mask of columns that are set
func (s *TrashSkillCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// TrashSkillUpdate contains model input fields for TrashSkill update operations.
type TrashSkillUpdate struct {
	TrashId *int32  `json:"trashId"`
	SkillId *string `json:"skillId"`
}

func (u *TrashSkillUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.TrashId != nil {
		cols = append(cols, "trashId")
		vals = append(vals, u.TrashId)
	}
	if u.SkillId != nil {
		cols = append(cols, "skillId")
		vals = append(vals, u.SkillId)
	}
	return cols, vals
}

func assignmentsToTrashSkillUpdate(assignments []FieldAssignment) (TrashSkillUpdate, error) {
	var input TrashSkillUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "trashId":
			if v, ok := a.Val.(int32); ok {
				input.TrashId = &v
			} else if v, ok := a.Val.(*int32); ok {
				input.TrashId = v
			} else {
				errs.Add("trashId", a.Val, "type", "field trashId must be of type int32")
			}
		case "skillId":
			if v, ok := a.Val.(string); ok {
				input.SkillId = &v
				errs.ValidateString("skillId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.SkillId = v
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// TrashSkillSelect specifies which scalar and relation fields to select for TrashSkill.
//
// Selectable fields:
//
//	-- Scalars --
//	trashId (bool)
//	skillId (bool)
//	-- Relations --
//	trash   (Trash)
//	skill   (Skill)
type TrashSkillSelect struct {
	TrashId bool         `json:"trashId"`
	SkillId bool         `json:"skillId"`
	Trash   *TrashSelect `json:"trash,omitempty"`
	Skill   *SkillSelect `json:"skill,omitempty"`
}

var fullTrashSkillSelectVal = &TrashSkillSelect{
	TrashId: true,
	SkillId: true,
}

func fullTrashSkillSelect() *TrashSkillSelect {
	return fullTrashSkillSelectVal
}

func (s *TrashSkillSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.TrashId || s.SkillId
}

func (s *TrashSkillSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type TrashSkillOmit struct {
	TrashId bool `json:"trashId"`
	SkillId bool `json:"skillId"`
}

type TrashSkillSelectQuery interface {
	GetRelationParams() (*TrashSkillSelect, *TrashSkillOmit, QueryParams[TrashSkill])
}

func (s *TrashSkillSelect) GetRelationParams() (*TrashSkillSelect, *TrashSkillOmit, QueryParams[TrashSkill]) {
	return s, nil, QueryParams[TrashSkill]{}
}

type TrashSkillQueryBuilder struct {
	selects *TrashSkillSelect
	omits   *TrashSkillOmit
	where   []PredicateOf[TrashSkill]
	take    *int
	skip    *int
	orderBy []OrderBy[TrashSkill]
	cursor  UniquePredicate[TrashSkill]
}

func (b *TrashSkillQueryBuilder) Where(preds ...PredicateOf[TrashSkill]) *TrashSkillQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *TrashSkillQueryBuilder) Take(limit int) *TrashSkillQueryBuilder {
	b.take = &limit
	return b
}

func (b *TrashSkillQueryBuilder) Skip(offset int) *TrashSkillQueryBuilder {
	b.skip = &offset
	return b
}

func (b *TrashSkillQueryBuilder) OrderBy(orders ...OrderBy[TrashSkill]) *TrashSkillQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *TrashSkillQueryBuilder) Cursor(where UniquePredicate[TrashSkill]) *TrashSkillQueryBuilder {
	b.cursor = where
	return b
}

func (b *TrashSkillQueryBuilder) Select(s TrashSkillSelect) *TrashSkillQueryBuilder {
	b.selects = &s
	return b
}

func (b *TrashSkillQueryBuilder) Omit(o TrashSkillOmit) *TrashSkillQueryBuilder {
	b.omits = &o
	return b
}

func (b *TrashSkillQueryBuilder) GetRelationParams() (*TrashSkillSelect, *TrashSkillOmit, QueryParams[TrashSkill]) {
	if b == nil {
		return nil, nil, QueryParams[TrashSkill]{}
	}
	return b.selects, b.omits, QueryParams[TrashSkill]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// TrashSkillCreateArgs is the input argument passed to TrashSkill Create extension hooks.
//
// Fields for TrashSkill:
//
//	trashId int32  required
//	skillId string required
//
// Relations for TrashSkill:
//
//	trash   (Trash)
//	skill   (Skill)
type TrashSkillCreateArgs struct {
	// Data contains the model fields to insert.
	Data *TrashSkillCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *TrashSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// TrashSkillCreateManyArgs is the input argument passed to TrashSkill CreateMany extension hooks.
//
// Fields for TrashSkill:
//
//	trashId int32  required
//	skillId string required
type TrashSkillCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*TrashSkillCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *TrashSkillCreateManyArgs) AppendData(builders ...*TrashSkillCreateBuilder) *TrashSkillCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToTrashSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// TrashSkillCreateManyAndReturnArgs is the input argument passed to TrashSkill CreateManyAndReturn extension hooks.
//
// Fields for TrashSkill:
//
//	trashId int32  required
//	skillId string required
//
// Relations for TrashSkill:
//
//	trash   (Trash)
//	skill   (Skill)
type TrashSkillCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*TrashSkillCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *TrashSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *TrashSkillCreateManyAndReturnArgs) AppendData(builders ...*TrashSkillCreateBuilder) *TrashSkillCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToTrashSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// TrashSkillFindUniqueArgs is the input argument passed to TrashSkill FindUnique extension hooks.
type TrashSkillFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[TrashSkill]
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashSkillSelect
}

func (a *TrashSkillFindUniqueArgs) SetWhere(unique UniquePredicate[TrashSkill], additional ...PredicateOf[TrashSkill]) *TrashSkillFindUniqueArgs {
	a.Where = make([]PredicateOf[TrashSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashSkillFindFirstArgs is the input argument passed to TrashSkill FindFirst extension hooks.
type TrashSkillFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[TrashSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[TrashSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashSkillSelect
}

func (a *TrashSkillFindFirstArgs) SetWhere(preds ...PredicateOf[TrashSkill]) *TrashSkillFindFirstArgs {
	a.Where = preds
	return a
}

func (a *TrashSkillFindFirstArgs) SetOrderBy(orders ...OrderBy[TrashSkill]) *TrashSkillFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *TrashSkillFindFirstArgs) SetCursor(cursor UniquePredicate[TrashSkill]) *TrashSkillFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *TrashSkillFindFirstArgs) SetSkip(n int) *TrashSkillFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *TrashSkillFindFirstArgs) SetTake(n int) *TrashSkillFindFirstArgs {
	a.Take = &n
	return a
}

// TrashSkillFindManyArgs is the input argument passed to TrashSkill FindMany extension hooks.
type TrashSkillFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[TrashSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[TrashSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashSkillSelect
}

func (a *TrashSkillFindManyArgs) SetWhere(preds ...PredicateOf[TrashSkill]) *TrashSkillFindManyArgs {
	a.Where = preds
	return a
}

func (a *TrashSkillFindManyArgs) SetOrderBy(orders ...OrderBy[TrashSkill]) *TrashSkillFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *TrashSkillFindManyArgs) SetCursor(cursor UniquePredicate[TrashSkill]) *TrashSkillFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *TrashSkillFindManyArgs) SetSkip(n int) *TrashSkillFindManyArgs {
	a.Skip = &n
	return a
}

func (a *TrashSkillFindManyArgs) SetTake(n int) *TrashSkillFindManyArgs {
	a.Take = &n
	return a
}

// TrashSkillCountArgs is the input argument passed to TrashSkill Count extension hooks.
type TrashSkillCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *TrashSkillCountArgs) SetWhere(preds ...PredicateOf[TrashSkill]) *TrashSkillCountArgs {
	a.Where = preds
	return a
}

func (a *TrashSkillCountArgs) SetSkip(n int) *TrashSkillCountArgs {
	a.Skip = &n
	return a
}

func (a *TrashSkillCountArgs) SetTake(n int) *TrashSkillCountArgs {
	a.Take = &n
	return a
}

// TrashSkillDeleteArgs is the input argument passed to TrashSkill Delete extension hooks.
type TrashSkillDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[TrashSkill]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *TrashSkillSelect
}

func (a *TrashSkillDeleteArgs) SetWhere(unique UniquePredicate[TrashSkill], additional ...PredicateOf[TrashSkill]) *TrashSkillDeleteArgs {
	a.Where = make([]PredicateOf[TrashSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashSkillDeleteManyArgs is the input argument passed to TrashSkill DeleteMany extension hooks.
type TrashSkillDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashSkill]
}

func (a *TrashSkillDeleteManyArgs) SetWhere(preds ...PredicateOf[TrashSkill]) *TrashSkillDeleteManyArgs {
	a.Where = preds
	return a
}

// TrashSkillUpdateArgs is the input argument passed to TrashSkill Update extension hooks.
type TrashSkillUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[TrashSkill]
	// Data contains the model fields to update.
	Data *TrashSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *TrashSkillSelect
}

func (a *TrashSkillUpdateArgs) SetWhere(unique UniquePredicate[TrashSkill], additional ...PredicateOf[TrashSkill]) *TrashSkillUpdateArgs {
	a.Where = make([]PredicateOf[TrashSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashSkillUpdateManyArgs is the input argument passed to TrashSkill UpdateMany extension hooks.
type TrashSkillUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashSkill]
	// Data contains the model fields to update.
	Data *TrashSkillUpdate
}

func (a *TrashSkillUpdateManyArgs) SetWhere(preds ...PredicateOf[TrashSkill]) *TrashSkillUpdateManyArgs {
	a.Where = preds
	return a
}

// TrashSkillUpdateManyAndReturnArgs is the input argument passed to TrashSkill UpdateManyAndReturn extension hooks.
type TrashSkillUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashSkill]
	// Data contains the model fields to update.
	Data *TrashSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *TrashSkillSelect
}

func (a *TrashSkillUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[TrashSkill]) *TrashSkillUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type TrashSkillCreateQuery = func(ctx context.Context, args *TrashSkillCreateArgs) (*TrashSkill, error)
type TrashSkillCreateManyQuery = func(ctx context.Context, args *TrashSkillCreateManyArgs) (int64, error)
type TrashSkillCreateManyAndReturnQuery = func(ctx context.Context, args *TrashSkillCreateManyAndReturnArgs) ([]*TrashSkill, error)
type TrashSkillFindUniqueQuery = func(ctx context.Context, args *TrashSkillFindUniqueArgs) (*TrashSkill, error)
type TrashSkillFindFirstQuery = func(ctx context.Context, args *TrashSkillFindFirstArgs) (*TrashSkill, error)
type TrashSkillFindManyQuery = func(ctx context.Context, args *TrashSkillFindManyArgs) ([]*TrashSkill, error)
type TrashSkillDeleteQuery = func(ctx context.Context, args *TrashSkillDeleteArgs) (*TrashSkill, error)
type TrashSkillDeleteManyQuery = func(ctx context.Context, args *TrashSkillDeleteManyArgs) (int64, error)
type TrashSkillCountQuery = func(ctx context.Context, args *TrashSkillCountArgs) (int64, error)
type TrashSkillUpdateQuery = func(ctx context.Context, args *TrashSkillUpdateArgs) (*TrashSkill, error)
type TrashSkillUpdateManyQuery = func(ctx context.Context, args *TrashSkillUpdateManyArgs) (int64, error)
type TrashSkillUpdateManyAndReturnQuery = func(ctx context.Context, args *TrashSkillUpdateManyAndReturnArgs) ([]*TrashSkill, error)

type TrashSkillExtension struct {
	Create              func(ctx context.Context, args *TrashSkillCreateArgs, next TrashSkillCreateQuery) (*TrashSkill, error)
	CreateMany          func(ctx context.Context, args *TrashSkillCreateManyArgs, next TrashSkillCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *TrashSkillCreateManyAndReturnArgs, next TrashSkillCreateManyAndReturnQuery) ([]*TrashSkill, error)
	FindUnique          func(ctx context.Context, args *TrashSkillFindUniqueArgs, next TrashSkillFindUniqueQuery) (*TrashSkill, error)
	FindFirst           func(ctx context.Context, args *TrashSkillFindFirstArgs, next TrashSkillFindFirstQuery) (*TrashSkill, error)
	FindMany            func(ctx context.Context, args *TrashSkillFindManyArgs, next TrashSkillFindManyQuery) ([]*TrashSkill, error)
	Delete              func(ctx context.Context, args *TrashSkillDeleteArgs, next TrashSkillDeleteQuery) (*TrashSkill, error)
	DeleteMany          func(ctx context.Context, args *TrashSkillDeleteManyArgs, next TrashSkillDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *TrashSkillCountArgs, next TrashSkillCountQuery) (int64, error)
	Update              func(ctx context.Context, args *TrashSkillUpdateArgs, next TrashSkillUpdateQuery) (*TrashSkill, error)
	UpdateMany          func(ctx context.Context, args *TrashSkillUpdateManyArgs, next TrashSkillUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *TrashSkillUpdateManyAndReturnArgs, next TrashSkillUpdateManyAndReturnQuery) ([]*TrashSkill, error)
}

type TrashSkillDelegate struct {
	client     *Queries
	extensions []TrashSkillExtension
}

func (d *TrashSkillDelegate) Use(exts ...TrashSkillExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *TrashSkill) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "trashId":
			targets[i] = &m.TrashId
		case "skillId":
			targets[i] = &m.SkillId
		}
	}
	return targets
}

var trashSkillDefaultCols = []string{
	"trashId",
	"skillId",
}

var trashSkillPKCols = []string{
	"trashId",
	"skillId",
}

var trashSkillUniqueCols = []string{}

func selectTrashSkillCols(selects *TrashSkillSelect, omits *TrashSkillOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return trashSkillDefaultCols
	}

	anySelected := selects != nil && (selects.TrashId || selects.SkillId || selects.Trash != nil || selects.Skill != nil)

	specs := []colSpec{
		{"trashId", selects != nil && selects.TrashId, omits != nil && omits.TrashId, selects != nil && selects.Trash != nil},
		{"skillId", selects != nil && selects.SkillId, omits != nil && omits.SkillId, selects != nil && selects.Skill != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *TrashSkillSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Trash != nil || s.Skill != nil
}

type TrashSkillUpsertBuilder struct {
	*CreateBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]
}

func (b *TrashSkillUpsertBuilder) Select(s TrashSkillSelect) *TrashSkillUpsertBuilder {
	b.selects = &s
	return b
}

func (b *TrashSkillUpsertBuilder) Omit(o TrashSkillOmit) *TrashSkillUpsertBuilder {
	b.omits = &o
	return b
}

type TrashSkillCreateBuilder struct {
	*CreateBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]
}

func (b *TrashSkillCreateBuilder) Select(s TrashSkillSelect) *TrashSkillCreateBuilder {
	b.selects = &s
	return b
}

func (b *TrashSkillCreateBuilder) Omit(o TrashSkillOmit) *TrashSkillCreateBuilder {
	b.omits = &o
	return b
}

func (b *TrashSkillCreateBuilder) OnConflict(target UniqueConstraintTarget) *TrashSkillConflictBuilder[TrashSkillUpsertBuilder] {
	upsertBuilder := &TrashSkillUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &TrashSkillConflictBuilder[TrashSkillUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *TrashSkillCreateBuilder) SetTrashId(v int32) *TrashSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}
func (b *TrashSkillCreateBuilder) SetSkillId(v string) *TrashSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *TrashSkillCreateBuilder) Assignments(assignments ...FieldAssignmentOf[TrashSkill]) *TrashSkillCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *TrashSkillDelegate) Create() *TrashSkillCreateBuilder {
	return &TrashSkillCreateBuilder{
		CreateBuilder: &CreateBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedTrashSkillTrashId uint64 = 1 << 0
	providedTrashSkillSkillId uint64 = 1 << 1
)

func assignmentsToTrashSkillCreate(assignments []FieldAssignment) (TrashSkillCreate, error) {
	var input TrashSkillCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "trashId":
			provided |= providedTrashSkillTrashId
			if v, ok := a.Val.(int32); ok {
				input.TrashId = v
				errs.ValidateInt32("trashId", v, "")
			} else {
				errs.Add("trashId", a.Val, "type", "field trashId must be of type int32")
			}
		case "skillId":
			provided |= providedTrashSkillSkillId
			if v, ok := a.Val.(string); ok {
				input.SkillId = v
				errs.ValidateString("skillId", v, true, 0, false, false)
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		}
	}
	if provided&providedTrashSkillTrashId == 0 {
		errs.Add("trashId", nil, "required", "field TrashId is required")
	}
	if provided&providedTrashSkillSkillId == 0 {
		errs.Add("skillId", "", "required", "field SkillId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *TrashSkillCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "trashId")
	vals = append(vals, s.TrashId)
	cols = append(cols, "skillId")
	vals = append(vals, s.SkillId)
	return
}

func partitionTrashSkillInputs(dialect Dialect, inputs []*TrashSkillCreate) [][]*TrashSkillCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*TrashSkillCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*TrashSkillCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*TrashSkillCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*TrashSkillCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*TrashSkillCreate{inputs}
}

func (d *TrashSkillDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *TrashSkillSelect, omits *TrashSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*TrashSkill, error) {
	input, err := assignmentsToTrashSkillCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectTrashSkillCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashSkillCreateArgs) (*TrashSkill, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectTrashSkillCols(a.Select, omits)
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
			curr = func(c context.Context, a *TrashSkillCreateArgs) (*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type TrashSkillCreateManyBuilder struct {
	*CreateManyBuilder[TrashSkill]
}

func (b *TrashSkillCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *TrashSkillConflictBuilder[TrashSkillCreateManyBuilder] {
	return &TrashSkillConflictBuilder[TrashSkillCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type TrashSkillCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]
}

func (b *TrashSkillCreateManyAndReturnBuilder) Select(s TrashSkillSelect) *TrashSkillCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *TrashSkillCreateManyAndReturnBuilder) Omit(o TrashSkillOmit) *TrashSkillCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *TrashSkillCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *TrashSkillConflictBuilder[TrashSkillCreateManyAndReturnBuilder] {
	return &TrashSkillConflictBuilder[TrashSkillCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToTrashSkillRecordInputs(builders []*TrashSkillCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *TrashSkillDelegate) CreateMany(builders ...*TrashSkillCreateBuilder) *TrashSkillCreateManyBuilder {
	return &TrashSkillCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[TrashSkill]{
			records:  createBuildersToTrashSkillRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *TrashSkillDelegate) CreateManyAndReturn(builders ...*TrashSkillCreateBuilder) *TrashSkillCreateManyAndReturnBuilder {
	return &TrashSkillCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
			records:  createBuildersToTrashSkillRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToTrashSkillCreateInputs(records []RecordInput) ([]*TrashSkillCreate, error) {
	structs := make([]TrashSkillCreate, len(records))
	inputs := make([]*TrashSkillCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToTrashSkillCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *TrashSkillDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToTrashSkillCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &TrashSkillCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashSkillCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashSkillCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *TrashSkillSelect, omits *TrashSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*TrashSkill, error) {
	inputs, err := recordsToTrashSkillCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashSkillCreateManyAndReturnArgs) ([]*TrashSkill, error) {
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
			curr = func(c context.Context, a *TrashSkillCreateManyAndReturnArgs) ([]*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *TrashSkillSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*TrashSkill, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *TrashSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.TrashSkill.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.TrashSkill.loadRelations(ctx, []*TrashSkill{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "TrashSkill", cols, returningCols, trashSkillPKCols, conflictTarget, conflictAction, len(vals))
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

		var res TrashSkill
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, trashSkillPKCols)
}

func (d *TrashSkillDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*TrashSkill, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("TrashSkill") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "TrashSkill")
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

	var res TrashSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *TrashSkillDelegate) buildBulkInsertSQL(q *Queries, batch []*TrashSkillCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range trashSkillDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "TrashSkill")
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
			case "trashId":
				vals = append(vals, input.TrashId)
			case "skillId":
				vals = append(vals, input.SkillId)
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

func applyTrashSkillConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanTrashSkillRows(rows *sql.Rows, returningCols []string) ([]*TrashSkill, error) {
	var records []*TrashSkill
	for rows.Next() {
		var res TrashSkill
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

func (d *TrashSkillDelegate) runCreateMany(ctx context.Context, inputs []*TrashSkillCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionTrashSkillInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashSkillConflictClause(d.client.dialect, queryStr, vals, cols, trashSkillPKCols, conflictTarget, conflictAction)

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

func (d *TrashSkillDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*TrashSkillCreate,
	selects *TrashSkillSelect,
	omits *TrashSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*TrashSkill, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionTrashSkillInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*TrashSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.TrashSkill.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.TrashSkill.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectTrashSkillCols(selects, omits, trashSkillPKCols...)
	recordsOut := make([]*TrashSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashSkillConflictClause(d.client.dialect, queryStr, vals, cols, trashSkillPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanTrashSkillRows(rows, returningCols)
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

func (d *TrashSkillDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*TrashSkillCreate,
	selects *TrashSkillSelect,
	omits *TrashSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*TrashSkill, error) {
	batches := partitionTrashSkillInputs(d.client.dialect, inputs)
	returningCols := selectTrashSkillCols(selects, omits, trashSkillPKCols...)
	recordsOut := make([]*TrashSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashSkillConflictClause(d.client.dialect, queryStr, vals, cols, trashSkillPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("TrashSkill") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "TrashSkill")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, trashSkillPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, trashSkillPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanTrashSkillRows(rows, returningCols)
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

type TrashSkillConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *TrashSkillConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *TrashSkillConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *TrashSkillConflictBuilder[B]) Update(fn func(u *TrashSkillUpsert)) *B {
	cb.setAction(*TrashSkillConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// TrashSkillConflictUpdate creates a custom ConflictAction for TrashSkill upsert conflicts.
func TrashSkillConflictUpdate(fn func(u *TrashSkillUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newTrashSkillUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type TrashSkillUpsert struct {
	TrashId numericFieldUpsert[int32]
	SkillId fieldUpsert[string]
}

func newTrashSkillUpsert(up *ConflictUpdate) *TrashSkillUpsert {
	return &TrashSkillUpsert{
		TrashId: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "trashId", update: up},
			tableName:   "TrashSkill",
		},
		SkillId: fieldUpsert[string]{column: "skillId", update: up},
	}
}

type TrashSkillUpdateBuilder struct {
	*UpdateBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]
}

type TrashSkillUpdateManyBuilder struct {
	*UpdateManyBuilder[TrashSkill]
}

type TrashSkillUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]
}

func (b *TrashSkillUpdateBuilder) SetTrashId(v int32) *TrashSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}

func (b *TrashSkillUpdateManyBuilder) SetTrashId(v int32) *TrashSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}

func (b *TrashSkillUpdateManyAndReturnBuilder) SetTrashId(v int32) *TrashSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}
func (b *TrashSkillUpdateBuilder) SetSkillId(v string) *TrashSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *TrashSkillUpdateManyBuilder) SetSkillId(v string) *TrashSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *TrashSkillUpdateManyAndReturnBuilder) SetSkillId(v string) *TrashSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *TrashSkillUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[TrashSkill]) *TrashSkillUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *TrashSkillUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[TrashSkill]) *TrashSkillUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *TrashSkillUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[TrashSkill]) *TrashSkillUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *TrashSkillDelegate) Update(where UniquePredicate[TrashSkill], additional ...PredicateOf[TrashSkill]) *TrashSkillUpdateBuilder {
	return &TrashSkillUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *TrashSkillDelegate) UpdateMany(preds ...PredicateOf[TrashSkill]) *TrashSkillUpdateManyBuilder {
	return &TrashSkillUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[TrashSkill]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *TrashSkillDelegate) UpdateManyAndReturn(preds ...PredicateOf[TrashSkill]) *TrashSkillUpdateManyAndReturnBuilder {
	return &TrashSkillUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *TrashSkillDelegate) buildUpdateSQL(preds []PredicateOf[TrashSkill], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "TrashSkill")
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

func (d *TrashSkillDelegate) executeUpdate(ctx context.Context, where UniquePredicate[TrashSkill], additional []PredicateOf[TrashSkill], assignments []FieldAssignment, selects *TrashSkillSelect, omits *TrashSkillOmit) (*TrashSkill, error) {
	allWhere := make([]PredicateOf[TrashSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToTrashSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashSkillUpdateArgs) (*TrashSkill, error) {
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
			curr = func(c context.Context, a *TrashSkillUpdateArgs) (*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) runUpdate(ctx context.Context, preds []PredicateOf[TrashSkill], cols []string, vals []any, selects *TrashSkillSelect, omits *TrashSkillOmit) (*TrashSkill, error) {
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
		var res *TrashSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.TrashSkill.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.TrashSkill.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectTrashSkillCols(selects, omits, trashSkillPKCols...)
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
		return nil, &NotFoundError{Model: "TrashSkill"}
	}

	var res TrashSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*TrashSkill{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *TrashSkillDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[TrashSkill], cols []string, vals []any) (int64, error) {
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

func (d *TrashSkillDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[TrashSkill], cols []string, vals []any, selects *TrashSkillSelect, omits *TrashSkillOmit) (*TrashSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "TrashSkill"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *TrashSkillDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[TrashSkill], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToTrashSkillUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &TrashSkillUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *TrashSkillUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashSkillUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *TrashSkillDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[TrashSkill], assignments []FieldAssignment, selects *TrashSkillSelect, omits *TrashSkillOmit) ([]*TrashSkill, error) {
	input, err := assignmentsToTrashSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashSkillUpdateManyAndReturnArgs) ([]*TrashSkill, error) {
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
			curr = func(c context.Context, a *TrashSkillUpdateManyAndReturnArgs) ([]*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[TrashSkill], cols []string, vals []any, selects *TrashSkillSelect, omits *TrashSkillOmit) ([]*TrashSkill, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[TrashSkill]{Where: preds}, selects, omits)
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
		var res []*TrashSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.TrashSkill.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.TrashSkill.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectTrashSkillCols(selects, omits, trashSkillPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanTrashSkillRows(rows, returningCols)
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

func (d *TrashSkillDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[TrashSkill], cols []string, vals []any, selects *TrashSkillSelect, omits *TrashSkillOmit) ([]*TrashSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*TrashSkill{}, nil
	}
	return d.runFindMany(ctx, QueryParams[TrashSkill]{Where: preds}, selects, omits)
}
func (d *TrashSkillDelegate) FindUnique(where UniquePredicate[TrashSkill], additional ...PredicateOf[TrashSkill]) *FindUniqueBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit] {
	return &FindUniqueBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *TrashSkillDelegate) FindFirst(preds ...PredicateOf[TrashSkill]) *FindFirstBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit] {
	return &FindFirstBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *TrashSkillDelegate) FindMany(preds ...PredicateOf[TrashSkill]) *FindManyBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit] {
	return &FindManyBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *TrashSkillDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[TrashSkill], additional []PredicateOf[TrashSkill], selects *TrashSkillSelect, omits *TrashSkillOmit) (*TrashSkill, error) {
	allWhere := make([]PredicateOf[TrashSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashSkillFindUniqueArgs) (*TrashSkill, error) {
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
			curr = func(c context.Context, a *TrashSkillFindUniqueArgs) (*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[TrashSkill],
	selects *TrashSkillSelect,
	omits *TrashSkillOmit,
) (*TrashSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *TrashSkillFindFirstArgs) (*TrashSkill, error) {
		return d.runFindFirst(c, QueryParams[TrashSkill]{
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
			curr = func(c context.Context, a *TrashSkillFindFirstArgs) (*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[TrashSkill],
	selects *TrashSkillSelect,
	omits *TrashSkillOmit,
) ([]*TrashSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *TrashSkillFindManyArgs) ([]*TrashSkill, error) {
		return d.runFindMany(c, QueryParams[TrashSkill]{
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
			curr = func(c context.Context, a *TrashSkillFindManyArgs) ([]*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) runFindUnique(ctx context.Context, where []PredicateOf[TrashSkill], selects *TrashSkillSelect, omits *TrashSkillOmit) (*TrashSkill, error) {
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
	returningCols := selectTrashSkillCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*TrashSkill{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *TrashSkillDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[TrashSkill],
	selects *TrashSkillSelect,
	omits *TrashSkillOmit,
) (*TrashSkill, error) {
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

func (d *TrashSkillDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[TrashSkill],
	selects *TrashSkillSelect,
	omits *TrashSkillOmit,
) ([]*TrashSkill, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, trashSkillPKCols, trashSkillUniqueCols, "TrashSkill", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, trashSkillPKCols, trashSkillUniqueCols, isCursorQuery, params.Take)
	returningCols := selectTrashSkillCols(selects, omits)

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

func (d *TrashSkillDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*TrashSkill, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "TrashSkill", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res TrashSkill
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *TrashSkillDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*TrashSkill, error) {
	query := buildSelectSQL(d.client, "TrashSkill", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*TrashSkill, 0)
	for rows.Next() {
		var res TrashSkill
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
func (d *TrashSkillDelegate) DeleteMany(preds ...PredicateOf[TrashSkill]) *DeleteManyBuilder[TrashSkill] {
	return &DeleteManyBuilder[TrashSkill]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *TrashSkillDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[TrashSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &TrashSkillDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *TrashSkillDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashSkillDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[TrashSkill]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "TrashSkill")
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

func (d *TrashSkillDelegate) Delete(where UniquePredicate[TrashSkill], additional ...PredicateOf[TrashSkill]) *DeleteBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit] {
	return &DeleteBuilder[TrashSkill, TrashSkillSelect, TrashSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *TrashSkillDelegate) executeDelete(ctx context.Context, where UniquePredicate[TrashSkill], additional []PredicateOf[TrashSkill], selects *TrashSkillSelect, omits *TrashSkillOmit) (*TrashSkill, error) {
	allWhere := make([]PredicateOf[TrashSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSkillSelect()
	}

	args := &TrashSkillDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashSkillDeleteArgs) (*TrashSkill, error) {
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
			curr = func(c context.Context, a *TrashSkillDeleteArgs) (*TrashSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) runDelete(ctx context.Context, where []PredicateOf[TrashSkill], selects *TrashSkillSelect, omits *TrashSkillOmit) (*TrashSkill, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectTrashSkillCols(selects, omits, trashSkillPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *TrashSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.TrashSkill.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "TrashSkill"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "TrashSkill")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[TrashSkill]
			pkPreds = append(pkPreds, Predicate[TrashSkill]{
				Data: PredicateData{
					Column:   "trashId",
					Operator: "=",
					Value:    res.TrashId,
				},
			})
			pkPreds = append(pkPreds, Predicate[TrashSkill]{
				Data: PredicateData{
					Column:   "skillId",
					Operator: "=",
					Value:    res.SkillId,
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
	d.client.dialect.WriteQuotedIdent(&sb, "TrashSkill")

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
		return nil, &NotFoundError{Model: "TrashSkill"}
	}

	var row TrashSkill
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *TrashSkillDelegate) Count(preds ...PredicateOf[TrashSkill]) *CountBuilder[TrashSkill] {
	return &CountBuilder[TrashSkill]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *TrashSkillDelegate) executeCount(ctx context.Context, params QueryParams[TrashSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &TrashSkillCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *TrashSkillCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[TrashSkill]{
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
			curr = func(c context.Context, a *TrashSkillCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashSkillDelegate) runCount(ctx context.Context, params QueryParams[TrashSkill]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "TrashSkill")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "TrashSkill")
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
func (d *TrashSkillDelegate) loadRelations(ctx context.Context, records []*TrashSkill, selects *TrashSkillSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Trash != nil {
		relationSelects, relationOmits, relationParams := selects.Trash.GetRelationParams()
		returningCols := selectTrashCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: TrashSkill.trashId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *TrashSkill) int32 { return p.TrashId }),
			"Trash",
			"id",
			returningCols,
			scanInto(returningCols, (*Trash).ScanFields),
			directKey(func(c *Trash) int32 { return c.Id }),
			setOne(func(p *TrashSkill, c *Trash) { p.Trash = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading trash: %w", err)
		}
		if err := d.client.Trash.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Skill != nil {
		relationSelects, relationOmits, relationParams := selects.Skill.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: TrashSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *TrashSkill) string { return p.SkillId }),
			"Skill",
			"id",
			returningCols,
			scanInto(returningCols, (*Skill).ScanFields),
			directKey(func(c *Skill) string { return c.Id }),
			setOne(func(p *TrashSkill, c *Skill) { p.Skill = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skill: %w", err)
		}
		if err := d.client.Skill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
