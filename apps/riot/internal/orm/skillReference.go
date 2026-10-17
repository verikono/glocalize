package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// SkillReference represents the database model
type SkillReference struct {
	SkillId       string     `db:"skillId" json:"skillId"`
	ReferenceName string     `db:"referenceName" json:"referenceName"`
	Skill         *Skill     `json:"skill,omitempty"`
	Reference     *Reference `json:"reference,omitempty"`
}

// SkillReferenceCreate contains model input fields for SkillReference creation operations.
//
// Fields for SkillReference:
//
//	skillId       string required
//	referenceName string required
type SkillReferenceCreate struct {
	SkillId       string `json:"skillId"`
	ReferenceName string `json:"referenceName"`
}

// colMask returns a bit mask of columns that are set
func (s *SkillReferenceCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// SkillReferenceUpdate contains model input fields for SkillReference update operations.
type SkillReferenceUpdate struct {
	SkillId       *string `json:"skillId"`
	ReferenceName *string `json:"referenceName"`
}

func (u *SkillReferenceUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.SkillId != nil {
		cols = append(cols, "skillId")
		vals = append(vals, u.SkillId)
	}
	if u.ReferenceName != nil {
		cols = append(cols, "referenceName")
		vals = append(vals, u.ReferenceName)
	}
	return cols, vals
}

func assignmentsToSkillReferenceUpdate(assignments []FieldAssignment) (SkillReferenceUpdate, error) {
	var input SkillReferenceUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "skillId":
			if v, ok := a.Val.(string); ok {
				input.SkillId = &v
				errs.ValidateString("skillId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.SkillId = v
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		case "referenceName":
			if v, ok := a.Val.(string); ok {
				input.ReferenceName = &v
				errs.ValidateString("referenceName", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.ReferenceName = v
			} else {
				errs.Add("referenceName", a.Val, "type", "field referenceName must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// SkillReferenceSelect specifies which scalar and relation fields to select for SkillReference.
//
// Selectable fields:
//
//	-- Scalars --
//	skillId       (bool)
//	referenceName (bool)
//	-- Relations --
//	skill         (Skill)
//	reference     (Reference)
type SkillReferenceSelect struct {
	SkillId       bool             `json:"skillId"`
	ReferenceName bool             `json:"referenceName"`
	Skill         *SkillSelect     `json:"skill,omitempty"`
	Reference     *ReferenceSelect `json:"reference,omitempty"`
}

var fullSkillReferenceSelectVal = &SkillReferenceSelect{
	SkillId:       true,
	ReferenceName: true,
}

func fullSkillReferenceSelect() *SkillReferenceSelect {
	return fullSkillReferenceSelectVal
}

func (s *SkillReferenceSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.SkillId || s.ReferenceName
}

func (s *SkillReferenceSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type SkillReferenceOmit struct {
	SkillId       bool `json:"skillId"`
	ReferenceName bool `json:"referenceName"`
}

type SkillReferenceSelectQuery interface {
	GetRelationParams() (*SkillReferenceSelect, *SkillReferenceOmit, QueryParams[SkillReference])
}

func (s *SkillReferenceSelect) GetRelationParams() (*SkillReferenceSelect, *SkillReferenceOmit, QueryParams[SkillReference]) {
	return s, nil, QueryParams[SkillReference]{}
}

type SkillReferenceQueryBuilder struct {
	selects *SkillReferenceSelect
	omits   *SkillReferenceOmit
	where   []PredicateOf[SkillReference]
	take    *int
	skip    *int
	orderBy []OrderBy[SkillReference]
	cursor  UniquePredicate[SkillReference]
}

func (b *SkillReferenceQueryBuilder) Where(preds ...PredicateOf[SkillReference]) *SkillReferenceQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *SkillReferenceQueryBuilder) Take(limit int) *SkillReferenceQueryBuilder {
	b.take = &limit
	return b
}

func (b *SkillReferenceQueryBuilder) Skip(offset int) *SkillReferenceQueryBuilder {
	b.skip = &offset
	return b
}

func (b *SkillReferenceQueryBuilder) OrderBy(orders ...OrderBy[SkillReference]) *SkillReferenceQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *SkillReferenceQueryBuilder) Cursor(where UniquePredicate[SkillReference]) *SkillReferenceQueryBuilder {
	b.cursor = where
	return b
}

func (b *SkillReferenceQueryBuilder) Select(s SkillReferenceSelect) *SkillReferenceQueryBuilder {
	b.selects = &s
	return b
}

func (b *SkillReferenceQueryBuilder) Omit(o SkillReferenceOmit) *SkillReferenceQueryBuilder {
	b.omits = &o
	return b
}

func (b *SkillReferenceQueryBuilder) GetRelationParams() (*SkillReferenceSelect, *SkillReferenceOmit, QueryParams[SkillReference]) {
	if b == nil {
		return nil, nil, QueryParams[SkillReference]{}
	}
	return b.selects, b.omits, QueryParams[SkillReference]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// SkillReferenceCreateArgs is the input argument passed to SkillReference Create extension hooks.
//
// Fields for SkillReference:
//
//	skillId       string required
//	referenceName string required
//
// Relations for SkillReference:
//
//	skill         (Skill)
//	reference     (Reference)
type SkillReferenceCreateArgs struct {
	// Data contains the model fields to insert.
	Data *SkillReferenceCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *SkillReferenceSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// SkillReferenceCreateManyArgs is the input argument passed to SkillReference CreateMany extension hooks.
//
// Fields for SkillReference:
//
//	skillId       string required
//	referenceName string required
type SkillReferenceCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*SkillReferenceCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *SkillReferenceCreateManyArgs) AppendData(builders ...*SkillReferenceCreateBuilder) *SkillReferenceCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToSkillReferenceCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// SkillReferenceCreateManyAndReturnArgs is the input argument passed to SkillReference CreateManyAndReturn extension hooks.
//
// Fields for SkillReference:
//
//	skillId       string required
//	referenceName string required
//
// Relations for SkillReference:
//
//	skill         (Skill)
//	reference     (Reference)
type SkillReferenceCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*SkillReferenceCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *SkillReferenceSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *SkillReferenceCreateManyAndReturnArgs) AppendData(builders ...*SkillReferenceCreateBuilder) *SkillReferenceCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToSkillReferenceCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// SkillReferenceFindUniqueArgs is the input argument passed to SkillReference FindUnique extension hooks.
type SkillReferenceFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[SkillReference]
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillReferenceSelect
}

func (a *SkillReferenceFindUniqueArgs) SetWhere(unique UniquePredicate[SkillReference], additional ...PredicateOf[SkillReference]) *SkillReferenceFindUniqueArgs {
	a.Where = make([]PredicateOf[SkillReference], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillReferenceFindFirstArgs is the input argument passed to SkillReference FindFirst extension hooks.
type SkillReferenceFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillReference]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[SkillReference]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[SkillReference]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillReferenceSelect
}

func (a *SkillReferenceFindFirstArgs) SetWhere(preds ...PredicateOf[SkillReference]) *SkillReferenceFindFirstArgs {
	a.Where = preds
	return a
}

func (a *SkillReferenceFindFirstArgs) SetOrderBy(orders ...OrderBy[SkillReference]) *SkillReferenceFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *SkillReferenceFindFirstArgs) SetCursor(cursor UniquePredicate[SkillReference]) *SkillReferenceFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *SkillReferenceFindFirstArgs) SetSkip(n int) *SkillReferenceFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *SkillReferenceFindFirstArgs) SetTake(n int) *SkillReferenceFindFirstArgs {
	a.Take = &n
	return a
}

// SkillReferenceFindManyArgs is the input argument passed to SkillReference FindMany extension hooks.
type SkillReferenceFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillReference]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[SkillReference]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[SkillReference]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *SkillReferenceSelect
}

func (a *SkillReferenceFindManyArgs) SetWhere(preds ...PredicateOf[SkillReference]) *SkillReferenceFindManyArgs {
	a.Where = preds
	return a
}

func (a *SkillReferenceFindManyArgs) SetOrderBy(orders ...OrderBy[SkillReference]) *SkillReferenceFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *SkillReferenceFindManyArgs) SetCursor(cursor UniquePredicate[SkillReference]) *SkillReferenceFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *SkillReferenceFindManyArgs) SetSkip(n int) *SkillReferenceFindManyArgs {
	a.Skip = &n
	return a
}

func (a *SkillReferenceFindManyArgs) SetTake(n int) *SkillReferenceFindManyArgs {
	a.Take = &n
	return a
}

// SkillReferenceCountArgs is the input argument passed to SkillReference Count extension hooks.
type SkillReferenceCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillReference]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *SkillReferenceCountArgs) SetWhere(preds ...PredicateOf[SkillReference]) *SkillReferenceCountArgs {
	a.Where = preds
	return a
}

func (a *SkillReferenceCountArgs) SetSkip(n int) *SkillReferenceCountArgs {
	a.Skip = &n
	return a
}

func (a *SkillReferenceCountArgs) SetTake(n int) *SkillReferenceCountArgs {
	a.Take = &n
	return a
}

// SkillReferenceDeleteArgs is the input argument passed to SkillReference Delete extension hooks.
type SkillReferenceDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[SkillReference]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *SkillReferenceSelect
}

func (a *SkillReferenceDeleteArgs) SetWhere(unique UniquePredicate[SkillReference], additional ...PredicateOf[SkillReference]) *SkillReferenceDeleteArgs {
	a.Where = make([]PredicateOf[SkillReference], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillReferenceDeleteManyArgs is the input argument passed to SkillReference DeleteMany extension hooks.
type SkillReferenceDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillReference]
}

func (a *SkillReferenceDeleteManyArgs) SetWhere(preds ...PredicateOf[SkillReference]) *SkillReferenceDeleteManyArgs {
	a.Where = preds
	return a
}

// SkillReferenceUpdateArgs is the input argument passed to SkillReference Update extension hooks.
type SkillReferenceUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[SkillReference]
	// Data contains the model fields to update.
	Data *SkillReferenceUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *SkillReferenceSelect
}

func (a *SkillReferenceUpdateArgs) SetWhere(unique UniquePredicate[SkillReference], additional ...PredicateOf[SkillReference]) *SkillReferenceUpdateArgs {
	a.Where = make([]PredicateOf[SkillReference], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// SkillReferenceUpdateManyArgs is the input argument passed to SkillReference UpdateMany extension hooks.
type SkillReferenceUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillReference]
	// Data contains the model fields to update.
	Data *SkillReferenceUpdate
}

func (a *SkillReferenceUpdateManyArgs) SetWhere(preds ...PredicateOf[SkillReference]) *SkillReferenceUpdateManyArgs {
	a.Where = preds
	return a
}

// SkillReferenceUpdateManyAndReturnArgs is the input argument passed to SkillReference UpdateManyAndReturn extension hooks.
type SkillReferenceUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[SkillReference]
	// Data contains the model fields to update.
	Data *SkillReferenceUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *SkillReferenceSelect
}

func (a *SkillReferenceUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[SkillReference]) *SkillReferenceUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type SkillReferenceCreateQuery = func(ctx context.Context, args *SkillReferenceCreateArgs) (*SkillReference, error)
type SkillReferenceCreateManyQuery = func(ctx context.Context, args *SkillReferenceCreateManyArgs) (int64, error)
type SkillReferenceCreateManyAndReturnQuery = func(ctx context.Context, args *SkillReferenceCreateManyAndReturnArgs) ([]*SkillReference, error)
type SkillReferenceFindUniqueQuery = func(ctx context.Context, args *SkillReferenceFindUniqueArgs) (*SkillReference, error)
type SkillReferenceFindFirstQuery = func(ctx context.Context, args *SkillReferenceFindFirstArgs) (*SkillReference, error)
type SkillReferenceFindManyQuery = func(ctx context.Context, args *SkillReferenceFindManyArgs) ([]*SkillReference, error)
type SkillReferenceDeleteQuery = func(ctx context.Context, args *SkillReferenceDeleteArgs) (*SkillReference, error)
type SkillReferenceDeleteManyQuery = func(ctx context.Context, args *SkillReferenceDeleteManyArgs) (int64, error)
type SkillReferenceCountQuery = func(ctx context.Context, args *SkillReferenceCountArgs) (int64, error)
type SkillReferenceUpdateQuery = func(ctx context.Context, args *SkillReferenceUpdateArgs) (*SkillReference, error)
type SkillReferenceUpdateManyQuery = func(ctx context.Context, args *SkillReferenceUpdateManyArgs) (int64, error)
type SkillReferenceUpdateManyAndReturnQuery = func(ctx context.Context, args *SkillReferenceUpdateManyAndReturnArgs) ([]*SkillReference, error)

type SkillReferenceExtension struct {
	Create              func(ctx context.Context, args *SkillReferenceCreateArgs, next SkillReferenceCreateQuery) (*SkillReference, error)
	CreateMany          func(ctx context.Context, args *SkillReferenceCreateManyArgs, next SkillReferenceCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *SkillReferenceCreateManyAndReturnArgs, next SkillReferenceCreateManyAndReturnQuery) ([]*SkillReference, error)
	FindUnique          func(ctx context.Context, args *SkillReferenceFindUniqueArgs, next SkillReferenceFindUniqueQuery) (*SkillReference, error)
	FindFirst           func(ctx context.Context, args *SkillReferenceFindFirstArgs, next SkillReferenceFindFirstQuery) (*SkillReference, error)
	FindMany            func(ctx context.Context, args *SkillReferenceFindManyArgs, next SkillReferenceFindManyQuery) ([]*SkillReference, error)
	Delete              func(ctx context.Context, args *SkillReferenceDeleteArgs, next SkillReferenceDeleteQuery) (*SkillReference, error)
	DeleteMany          func(ctx context.Context, args *SkillReferenceDeleteManyArgs, next SkillReferenceDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *SkillReferenceCountArgs, next SkillReferenceCountQuery) (int64, error)
	Update              func(ctx context.Context, args *SkillReferenceUpdateArgs, next SkillReferenceUpdateQuery) (*SkillReference, error)
	UpdateMany          func(ctx context.Context, args *SkillReferenceUpdateManyArgs, next SkillReferenceUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *SkillReferenceUpdateManyAndReturnArgs, next SkillReferenceUpdateManyAndReturnQuery) ([]*SkillReference, error)
}

type SkillReferenceDelegate struct {
	client     *Queries
	extensions []SkillReferenceExtension
}

func (d *SkillReferenceDelegate) Use(exts ...SkillReferenceExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *SkillReference) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "skillId":
			targets[i] = &m.SkillId
		case "referenceName":
			targets[i] = &m.ReferenceName
		}
	}
	return targets
}

var skillReferenceDefaultCols = []string{
	"skillId",
	"referenceName",
}

var skillReferencePKCols = []string{
	"skillId",
	"referenceName",
}

var skillReferenceUniqueCols = []string{}

func selectSkillReferenceCols(selects *SkillReferenceSelect, omits *SkillReferenceOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return skillReferenceDefaultCols
	}

	anySelected := selects != nil && (selects.SkillId || selects.ReferenceName || selects.Skill != nil || selects.Reference != nil)

	specs := []colSpec{
		{"skillId", selects != nil && selects.SkillId, omits != nil && omits.SkillId, selects != nil && selects.Skill != nil},
		{"referenceName", selects != nil && selects.ReferenceName, omits != nil && omits.ReferenceName, selects != nil && selects.Reference != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *SkillReferenceSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Skill != nil || s.Reference != nil
}

type SkillReferenceUpsertBuilder struct {
	*CreateBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]
}

func (b *SkillReferenceUpsertBuilder) Select(s SkillReferenceSelect) *SkillReferenceUpsertBuilder {
	b.selects = &s
	return b
}

func (b *SkillReferenceUpsertBuilder) Omit(o SkillReferenceOmit) *SkillReferenceUpsertBuilder {
	b.omits = &o
	return b
}

type SkillReferenceCreateBuilder struct {
	*CreateBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]
}

func (b *SkillReferenceCreateBuilder) Select(s SkillReferenceSelect) *SkillReferenceCreateBuilder {
	b.selects = &s
	return b
}

func (b *SkillReferenceCreateBuilder) Omit(o SkillReferenceOmit) *SkillReferenceCreateBuilder {
	b.omits = &o
	return b
}

func (b *SkillReferenceCreateBuilder) OnConflict(target UniqueConstraintTarget) *SkillReferenceConflictBuilder[SkillReferenceUpsertBuilder] {
	upsertBuilder := &SkillReferenceUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &SkillReferenceConflictBuilder[SkillReferenceUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *SkillReferenceCreateBuilder) SetSkillId(v string) *SkillReferenceCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}
func (b *SkillReferenceCreateBuilder) SetReferenceName(v string) *SkillReferenceCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referenceName", Val: v})
	return b
}

func (b *SkillReferenceCreateBuilder) Assignments(assignments ...FieldAssignmentOf[SkillReference]) *SkillReferenceCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *SkillReferenceDelegate) Create() *SkillReferenceCreateBuilder {
	return &SkillReferenceCreateBuilder{
		CreateBuilder: &CreateBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedSkillReferenceSkillId       uint64 = 1 << 0
	providedSkillReferenceReferenceName uint64 = 1 << 1
)

func assignmentsToSkillReferenceCreate(assignments []FieldAssignment) (SkillReferenceCreate, error) {
	var input SkillReferenceCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "skillId":
			provided |= providedSkillReferenceSkillId
			if v, ok := a.Val.(string); ok {
				input.SkillId = v
				errs.ValidateString("skillId", v, true, 0, false, false)
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		case "referenceName":
			provided |= providedSkillReferenceReferenceName
			if v, ok := a.Val.(string); ok {
				input.ReferenceName = v
				errs.ValidateString("referenceName", v, true, 0, false, false)
			} else {
				errs.Add("referenceName", a.Val, "type", "field referenceName must be of type string")
			}
		}
	}
	if provided&providedSkillReferenceSkillId == 0 {
		errs.Add("skillId", "", "required", "field SkillId is required")
	}
	if provided&providedSkillReferenceReferenceName == 0 {
		errs.Add("referenceName", "", "required", "field ReferenceName is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *SkillReferenceCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "skillId")
	vals = append(vals, s.SkillId)
	cols = append(cols, "referenceName")
	vals = append(vals, s.ReferenceName)
	return
}

func partitionSkillReferenceInputs(dialect Dialect, inputs []*SkillReferenceCreate) [][]*SkillReferenceCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*SkillReferenceCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*SkillReferenceCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*SkillReferenceCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*SkillReferenceCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*SkillReferenceCreate{inputs}
}

func (d *SkillReferenceDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *SkillReferenceSelect, omits *SkillReferenceOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*SkillReference, error) {
	input, err := assignmentsToSkillReferenceCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectSkillReferenceCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillReferenceCreateArgs) (*SkillReference, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectSkillReferenceCols(a.Select, omits)
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
			curr = func(c context.Context, a *SkillReferenceCreateArgs) (*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type SkillReferenceCreateManyBuilder struct {
	*CreateManyBuilder[SkillReference]
}

func (b *SkillReferenceCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *SkillReferenceConflictBuilder[SkillReferenceCreateManyBuilder] {
	return &SkillReferenceConflictBuilder[SkillReferenceCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type SkillReferenceCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]
}

func (b *SkillReferenceCreateManyAndReturnBuilder) Select(s SkillReferenceSelect) *SkillReferenceCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *SkillReferenceCreateManyAndReturnBuilder) Omit(o SkillReferenceOmit) *SkillReferenceCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *SkillReferenceCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *SkillReferenceConflictBuilder[SkillReferenceCreateManyAndReturnBuilder] {
	return &SkillReferenceConflictBuilder[SkillReferenceCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToSkillReferenceRecordInputs(builders []*SkillReferenceCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *SkillReferenceDelegate) CreateMany(builders ...*SkillReferenceCreateBuilder) *SkillReferenceCreateManyBuilder {
	return &SkillReferenceCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[SkillReference]{
			records:  createBuildersToSkillReferenceRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *SkillReferenceDelegate) CreateManyAndReturn(builders ...*SkillReferenceCreateBuilder) *SkillReferenceCreateManyAndReturnBuilder {
	return &SkillReferenceCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
			records:  createBuildersToSkillReferenceRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToSkillReferenceCreateInputs(records []RecordInput) ([]*SkillReferenceCreate, error) {
	structs := make([]SkillReferenceCreate, len(records))
	inputs := make([]*SkillReferenceCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToSkillReferenceCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *SkillReferenceDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToSkillReferenceCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &SkillReferenceCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillReferenceCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillReferenceCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *SkillReferenceSelect, omits *SkillReferenceOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*SkillReference, error) {
	inputs, err := recordsToSkillReferenceCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *SkillReferenceCreateManyAndReturnArgs) ([]*SkillReference, error) {
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
			curr = func(c context.Context, a *SkillReferenceCreateManyAndReturnArgs) ([]*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *SkillReferenceSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*SkillReference, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *SkillReference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.SkillReference.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.SkillReference.loadRelations(ctx, []*SkillReference{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "SkillReference", cols, returningCols, skillReferencePKCols, conflictTarget, conflictAction, len(vals))
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

		var res SkillReference
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, skillReferencePKCols)
}

func (d *SkillReferenceDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*SkillReference, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("SkillReference") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "SkillReference")
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

	var res SkillReference
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *SkillReferenceDelegate) buildBulkInsertSQL(q *Queries, batch []*SkillReferenceCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range skillReferenceDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "SkillReference")
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
			case "skillId":
				vals = append(vals, input.SkillId)
			case "referenceName":
				vals = append(vals, input.ReferenceName)
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

func applySkillReferenceConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanSkillReferenceRows(rows *sql.Rows, returningCols []string) ([]*SkillReference, error) {
	var records []*SkillReference
	for rows.Next() {
		var res SkillReference
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

func (d *SkillReferenceDelegate) runCreateMany(ctx context.Context, inputs []*SkillReferenceCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionSkillReferenceInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillReferenceConflictClause(d.client.dialect, queryStr, vals, cols, skillReferencePKCols, conflictTarget, conflictAction)

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

func (d *SkillReferenceDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*SkillReferenceCreate,
	selects *SkillReferenceSelect,
	omits *SkillReferenceOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*SkillReference, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionSkillReferenceInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*SkillReference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.SkillReference.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.SkillReference.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectSkillReferenceCols(selects, omits, skillReferencePKCols...)
	recordsOut := make([]*SkillReference, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillReferenceConflictClause(d.client.dialect, queryStr, vals, cols, skillReferencePKCols, conflictTarget, conflictAction)

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

		scanned, err := scanSkillReferenceRows(rows, returningCols)
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

func (d *SkillReferenceDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*SkillReferenceCreate,
	selects *SkillReferenceSelect,
	omits *SkillReferenceOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*SkillReference, error) {
	batches := partitionSkillReferenceInputs(d.client.dialect, inputs)
	returningCols := selectSkillReferenceCols(selects, omits, skillReferencePKCols...)
	recordsOut := make([]*SkillReference, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applySkillReferenceConflictClause(d.client.dialect, queryStr, vals, cols, skillReferencePKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("SkillReference") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "SkillReference")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, skillReferencePKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, skillReferencePKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanSkillReferenceRows(rows, returningCols)
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

type SkillReferenceConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *SkillReferenceConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *SkillReferenceConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *SkillReferenceConflictBuilder[B]) Update(fn func(u *SkillReferenceUpsert)) *B {
	cb.setAction(*SkillReferenceConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// SkillReferenceConflictUpdate creates a custom ConflictAction for SkillReference upsert conflicts.
func SkillReferenceConflictUpdate(fn func(u *SkillReferenceUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newSkillReferenceUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type SkillReferenceUpsert struct {
	SkillId       fieldUpsert[string]
	ReferenceName fieldUpsert[string]
}

func newSkillReferenceUpsert(up *ConflictUpdate) *SkillReferenceUpsert {
	return &SkillReferenceUpsert{
		SkillId:       fieldUpsert[string]{column: "skillId", update: up},
		ReferenceName: fieldUpsert[string]{column: "referenceName", update: up},
	}
}

type SkillReferenceUpdateBuilder struct {
	*UpdateBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]
}

type SkillReferenceUpdateManyBuilder struct {
	*UpdateManyBuilder[SkillReference]
}

type SkillReferenceUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]
}

func (b *SkillReferenceUpdateBuilder) SetSkillId(v string) *SkillReferenceUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *SkillReferenceUpdateManyBuilder) SetSkillId(v string) *SkillReferenceUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *SkillReferenceUpdateManyAndReturnBuilder) SetSkillId(v string) *SkillReferenceUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}
func (b *SkillReferenceUpdateBuilder) SetReferenceName(v string) *SkillReferenceUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referenceName", Val: v})
	return b
}

func (b *SkillReferenceUpdateManyBuilder) SetReferenceName(v string) *SkillReferenceUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referenceName", Val: v})
	return b
}

func (b *SkillReferenceUpdateManyAndReturnBuilder) SetReferenceName(v string) *SkillReferenceUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referenceName", Val: v})
	return b
}

func (b *SkillReferenceUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[SkillReference]) *SkillReferenceUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *SkillReferenceUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[SkillReference]) *SkillReferenceUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *SkillReferenceUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[SkillReference]) *SkillReferenceUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *SkillReferenceDelegate) Update(where UniquePredicate[SkillReference], additional ...PredicateOf[SkillReference]) *SkillReferenceUpdateBuilder {
	return &SkillReferenceUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *SkillReferenceDelegate) UpdateMany(preds ...PredicateOf[SkillReference]) *SkillReferenceUpdateManyBuilder {
	return &SkillReferenceUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[SkillReference]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *SkillReferenceDelegate) UpdateManyAndReturn(preds ...PredicateOf[SkillReference]) *SkillReferenceUpdateManyAndReturnBuilder {
	return &SkillReferenceUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *SkillReferenceDelegate) buildUpdateSQL(preds []PredicateOf[SkillReference], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "SkillReference")
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

func (d *SkillReferenceDelegate) executeUpdate(ctx context.Context, where UniquePredicate[SkillReference], additional []PredicateOf[SkillReference], assignments []FieldAssignment, selects *SkillReferenceSelect, omits *SkillReferenceOmit) (*SkillReference, error) {
	allWhere := make([]PredicateOf[SkillReference], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToSkillReferenceUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillReferenceUpdateArgs) (*SkillReference, error) {
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
			curr = func(c context.Context, a *SkillReferenceUpdateArgs) (*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) runUpdate(ctx context.Context, preds []PredicateOf[SkillReference], cols []string, vals []any, selects *SkillReferenceSelect, omits *SkillReferenceOmit) (*SkillReference, error) {
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
		var res *SkillReference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.SkillReference.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.SkillReference.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectSkillReferenceCols(selects, omits, skillReferencePKCols...)
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
		return nil, &NotFoundError{Model: "SkillReference"}
	}

	var res SkillReference
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*SkillReference{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *SkillReferenceDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[SkillReference], cols []string, vals []any) (int64, error) {
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

func (d *SkillReferenceDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[SkillReference], cols []string, vals []any, selects *SkillReferenceSelect, omits *SkillReferenceOmit) (*SkillReference, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "SkillReference"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *SkillReferenceDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[SkillReference], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToSkillReferenceUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &SkillReferenceUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *SkillReferenceUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillReferenceUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *SkillReferenceDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[SkillReference], assignments []FieldAssignment, selects *SkillReferenceSelect, omits *SkillReferenceOmit) ([]*SkillReference, error) {
	input, err := assignmentsToSkillReferenceUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillReferenceUpdateManyAndReturnArgs) ([]*SkillReference, error) {
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
			curr = func(c context.Context, a *SkillReferenceUpdateManyAndReturnArgs) ([]*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[SkillReference], cols []string, vals []any, selects *SkillReferenceSelect, omits *SkillReferenceOmit) ([]*SkillReference, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[SkillReference]{Where: preds}, selects, omits)
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
		var res []*SkillReference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.SkillReference.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.SkillReference.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectSkillReferenceCols(selects, omits, skillReferencePKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanSkillReferenceRows(rows, returningCols)
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

func (d *SkillReferenceDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[SkillReference], cols []string, vals []any, selects *SkillReferenceSelect, omits *SkillReferenceOmit) ([]*SkillReference, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*SkillReference{}, nil
	}
	return d.runFindMany(ctx, QueryParams[SkillReference]{Where: preds}, selects, omits)
}
func (d *SkillReferenceDelegate) FindUnique(where UniquePredicate[SkillReference], additional ...PredicateOf[SkillReference]) *FindUniqueBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit] {
	return &FindUniqueBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *SkillReferenceDelegate) FindFirst(preds ...PredicateOf[SkillReference]) *FindFirstBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit] {
	return &FindFirstBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *SkillReferenceDelegate) FindMany(preds ...PredicateOf[SkillReference]) *FindManyBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit] {
	return &FindManyBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *SkillReferenceDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[SkillReference], additional []PredicateOf[SkillReference], selects *SkillReferenceSelect, omits *SkillReferenceOmit) (*SkillReference, error) {
	allWhere := make([]PredicateOf[SkillReference], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillReferenceFindUniqueArgs) (*SkillReference, error) {
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
			curr = func(c context.Context, a *SkillReferenceFindUniqueArgs) (*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[SkillReference],
	selects *SkillReferenceSelect,
	omits *SkillReferenceOmit,
) (*SkillReference, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *SkillReferenceFindFirstArgs) (*SkillReference, error) {
		return d.runFindFirst(c, QueryParams[SkillReference]{
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
			curr = func(c context.Context, a *SkillReferenceFindFirstArgs) (*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[SkillReference],
	selects *SkillReferenceSelect,
	omits *SkillReferenceOmit,
) ([]*SkillReference, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *SkillReferenceFindManyArgs) ([]*SkillReference, error) {
		return d.runFindMany(c, QueryParams[SkillReference]{
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
			curr = func(c context.Context, a *SkillReferenceFindManyArgs) ([]*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) runFindUnique(ctx context.Context, where []PredicateOf[SkillReference], selects *SkillReferenceSelect, omits *SkillReferenceOmit) (*SkillReference, error) {
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
	returningCols := selectSkillReferenceCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*SkillReference{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *SkillReferenceDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[SkillReference],
	selects *SkillReferenceSelect,
	omits *SkillReferenceOmit,
) (*SkillReference, error) {
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

func (d *SkillReferenceDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[SkillReference],
	selects *SkillReferenceSelect,
	omits *SkillReferenceOmit,
) ([]*SkillReference, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, skillReferencePKCols, skillReferenceUniqueCols, "SkillReference", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, skillReferencePKCols, skillReferenceUniqueCols, isCursorQuery, params.Take)
	returningCols := selectSkillReferenceCols(selects, omits)

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

func (d *SkillReferenceDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*SkillReference, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "SkillReference", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res SkillReference
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *SkillReferenceDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*SkillReference, error) {
	query := buildSelectSQL(d.client, "SkillReference", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*SkillReference, 0)
	for rows.Next() {
		var res SkillReference
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
func (d *SkillReferenceDelegate) DeleteMany(preds ...PredicateOf[SkillReference]) *DeleteManyBuilder[SkillReference] {
	return &DeleteManyBuilder[SkillReference]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *SkillReferenceDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[SkillReference]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &SkillReferenceDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *SkillReferenceDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *SkillReferenceDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[SkillReference]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "SkillReference")
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

func (d *SkillReferenceDelegate) Delete(where UniquePredicate[SkillReference], additional ...PredicateOf[SkillReference]) *DeleteBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit] {
	return &DeleteBuilder[SkillReference, SkillReferenceSelect, SkillReferenceOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *SkillReferenceDelegate) executeDelete(ctx context.Context, where UniquePredicate[SkillReference], additional []PredicateOf[SkillReference], selects *SkillReferenceSelect, omits *SkillReferenceOmit) (*SkillReference, error) {
	allWhere := make([]PredicateOf[SkillReference], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullSkillReferenceSelect()
	}

	args := &SkillReferenceDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *SkillReferenceDeleteArgs) (*SkillReference, error) {
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
			curr = func(c context.Context, a *SkillReferenceDeleteArgs) (*SkillReference, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) runDelete(ctx context.Context, where []PredicateOf[SkillReference], selects *SkillReferenceSelect, omits *SkillReferenceOmit) (*SkillReference, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectSkillReferenceCols(selects, omits, skillReferencePKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *SkillReference
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.SkillReference.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "SkillReference"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "SkillReference")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[SkillReference]
			pkPreds = append(pkPreds, Predicate[SkillReference]{
				Data: PredicateData{
					Column:   "skillId",
					Operator: "=",
					Value:    res.SkillId,
				},
			})
			pkPreds = append(pkPreds, Predicate[SkillReference]{
				Data: PredicateData{
					Column:   "referenceName",
					Operator: "=",
					Value:    res.ReferenceName,
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
	d.client.dialect.WriteQuotedIdent(&sb, "SkillReference")

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
		return nil, &NotFoundError{Model: "SkillReference"}
	}

	var row SkillReference
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *SkillReferenceDelegate) Count(preds ...PredicateOf[SkillReference]) *CountBuilder[SkillReference] {
	return &CountBuilder[SkillReference]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *SkillReferenceDelegate) executeCount(ctx context.Context, params QueryParams[SkillReference]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &SkillReferenceCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *SkillReferenceCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[SkillReference]{
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
			curr = func(c context.Context, a *SkillReferenceCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *SkillReferenceDelegate) runCount(ctx context.Context, params QueryParams[SkillReference]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "SkillReference")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "SkillReference")
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
func (d *SkillReferenceDelegate) loadRelations(ctx context.Context, records []*SkillReference, selects *SkillReferenceSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Skill != nil {
		relationSelects, relationOmits, relationParams := selects.Skill.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: SkillReference.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *SkillReference) string { return p.SkillId }),
			"Skill",
			"id",
			returningCols,
			scanInto(returningCols, (*Skill).ScanFields),
			directKey(func(c *Skill) string { return c.Id }),
			setOne(func(p *SkillReference, c *Skill) { p.Skill = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skill: %w", err)
		}
		if err := d.client.Skill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Reference != nil {
		relationSelects, relationOmits, relationParams := selects.Reference.GetRelationParams()
		returningCols := selectReferenceCols(relationSelects, relationOmits, "name")
		// Current model holds the FK: SkillReference.referenceName
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *SkillReference) string { return p.ReferenceName }),
			"Reference",
			"name",
			returningCols,
			scanInto(returningCols, (*Reference).ScanFields),
			directKey(func(c *Reference) string { return c.Name }),
			setOne(func(p *SkillReference, c *Reference) { p.Reference = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading reference: %w", err)
		}
		if err := d.client.Reference.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
