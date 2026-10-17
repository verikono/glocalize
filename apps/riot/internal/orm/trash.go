package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Trash represents the database model
type Trash struct {
	Id        int32         `db:"id" json:"id"`
	CreatedAt time.Time     `db:"createdAt" json:"createdAt"`
	Content   string        `db:"content" json:"content"`
	Agents    []*TrashAgent `json:"agents,omitempty"`
	Skills    []*TrashSkill `json:"skills,omitempty"`
}

// TrashCreate contains model input fields for Trash creation operations.
//
// Fields for Trash:
//
//	id        int32     default: autoincrement()
//	createdAt time.Time default: now()
//	content   string    required
type TrashCreate struct {
	Id        *int32     `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	Content   string     `json:"content"`
}

// colMask returns a bit mask of columns that are set
func (s *TrashCreate) colMask() uint64 {
	var mask uint64
	if s.Id != nil {
		mask |= 1 << 0
	}
	mask |= 1 << 1
	mask |= 1 << 2
	return mask
}

// TrashUpdate contains model input fields for Trash update operations.
type TrashUpdate struct {
	Id        *int32     `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	Content   *string    `json:"content"`
}

func (u *TrashUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, u.Id)
	}
	if u.CreatedAt != nil {
		cols = append(cols, "createdAt")
		vals = append(vals, u.CreatedAt)
	}
	if u.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, u.Content)
	}
	return cols, vals
}

func assignmentsToTrashUpdate(assignments []FieldAssignment) (TrashUpdate, error) {
	var input TrashUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
			} else if v, ok := a.Val.(*int32); ok {
				input.Id = v
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "createdAt":
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			} else if v, ok := a.Val.(*time.Time); ok {
				input.CreatedAt = v
			} else {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		case "content":
			if v, ok := a.Val.(string); ok {
				input.Content = &v
				errs.ValidateString("content", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Content = v
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// TrashSelect specifies which scalar and relation fields to select for Trash.
//
// Selectable fields:
//
//	-- Scalars --
//	id        (bool)
//	createdAt (bool)
//	content   (bool)
//	-- Relations --
//	agents    ([]TrashAgent)
//	skills    ([]TrashSkill)
type TrashSelect struct {
	Id        bool                  `json:"id"`
	CreatedAt bool                  `json:"createdAt"`
	Content   bool                  `json:"content"`
	Agents    TrashAgentSelectQuery `json:"agents,omitempty"`
	Skills    TrashSkillSelectQuery `json:"skills,omitempty"`
}

var fullTrashSelectVal = &TrashSelect{
	Id:        true,
	CreatedAt: true,
	Content:   true,
}

func fullTrashSelect() *TrashSelect {
	return fullTrashSelectVal
}

func (s *TrashSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Id || s.CreatedAt || s.Content
}

func (s *TrashSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type TrashOmit struct {
	Id        bool `json:"id"`
	CreatedAt bool `json:"createdAt"`
	Content   bool `json:"content"`
}

type TrashSelectQuery interface {
	GetRelationParams() (*TrashSelect, *TrashOmit, QueryParams[Trash])
}

func (s *TrashSelect) GetRelationParams() (*TrashSelect, *TrashOmit, QueryParams[Trash]) {
	return s, nil, QueryParams[Trash]{}
}

type TrashQueryBuilder struct {
	selects *TrashSelect
	omits   *TrashOmit
	where   []PredicateOf[Trash]
	take    *int
	skip    *int
	orderBy []OrderBy[Trash]
	cursor  UniquePredicate[Trash]
}

func (b *TrashQueryBuilder) Where(preds ...PredicateOf[Trash]) *TrashQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *TrashQueryBuilder) Take(limit int) *TrashQueryBuilder {
	b.take = &limit
	return b
}

func (b *TrashQueryBuilder) Skip(offset int) *TrashQueryBuilder {
	b.skip = &offset
	return b
}

func (b *TrashQueryBuilder) OrderBy(orders ...OrderBy[Trash]) *TrashQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *TrashQueryBuilder) Cursor(where UniquePredicate[Trash]) *TrashQueryBuilder {
	b.cursor = where
	return b
}

func (b *TrashQueryBuilder) Select(s TrashSelect) *TrashQueryBuilder {
	b.selects = &s
	return b
}

func (b *TrashQueryBuilder) Omit(o TrashOmit) *TrashQueryBuilder {
	b.omits = &o
	return b
}

func (b *TrashQueryBuilder) GetRelationParams() (*TrashSelect, *TrashOmit, QueryParams[Trash]) {
	if b == nil {
		return nil, nil, QueryParams[Trash]{}
	}
	return b.selects, b.omits, QueryParams[Trash]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// TrashCreateArgs is the input argument passed to Trash Create extension hooks.
//
// Fields for Trash:
//
//	id        int32     default: autoincrement()
//	createdAt time.Time default: now()
//	content   string    required
//
// Relations for Trash:
//
//	agents    ([]TrashAgent)
//	skills    ([]TrashSkill)
type TrashCreateArgs struct {
	// Data contains the model fields to insert.
	Data *TrashCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *TrashSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// TrashCreateManyArgs is the input argument passed to Trash CreateMany extension hooks.
//
// Fields for Trash:
//
//	id        int32     default: autoincrement()
//	createdAt time.Time default: now()
//	content   string    required
type TrashCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*TrashCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *TrashCreateManyArgs) AppendData(builders ...*TrashCreateBuilder) *TrashCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToTrashCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// TrashCreateManyAndReturnArgs is the input argument passed to Trash CreateManyAndReturn extension hooks.
//
// Fields for Trash:
//
//	id        int32     default: autoincrement()
//	createdAt time.Time default: now()
//	content   string    required
//
// Relations for Trash:
//
//	agents    ([]TrashAgent)
//	skills    ([]TrashSkill)
type TrashCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*TrashCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *TrashSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *TrashCreateManyAndReturnArgs) AppendData(builders ...*TrashCreateBuilder) *TrashCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToTrashCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// TrashFindUniqueArgs is the input argument passed to Trash FindUnique extension hooks.
type TrashFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Trash]
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashSelect
}

func (a *TrashFindUniqueArgs) SetWhere(unique UniquePredicate[Trash], additional ...PredicateOf[Trash]) *TrashFindUniqueArgs {
	a.Where = make([]PredicateOf[Trash], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashFindFirstArgs is the input argument passed to Trash FindFirst extension hooks.
type TrashFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Trash]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Trash]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Trash]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashSelect
}

func (a *TrashFindFirstArgs) SetWhere(preds ...PredicateOf[Trash]) *TrashFindFirstArgs {
	a.Where = preds
	return a
}

func (a *TrashFindFirstArgs) SetOrderBy(orders ...OrderBy[Trash]) *TrashFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *TrashFindFirstArgs) SetCursor(cursor UniquePredicate[Trash]) *TrashFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *TrashFindFirstArgs) SetSkip(n int) *TrashFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *TrashFindFirstArgs) SetTake(n int) *TrashFindFirstArgs {
	a.Take = &n
	return a
}

// TrashFindManyArgs is the input argument passed to Trash FindMany extension hooks.
type TrashFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Trash]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Trash]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Trash]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashSelect
}

func (a *TrashFindManyArgs) SetWhere(preds ...PredicateOf[Trash]) *TrashFindManyArgs {
	a.Where = preds
	return a
}

func (a *TrashFindManyArgs) SetOrderBy(orders ...OrderBy[Trash]) *TrashFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *TrashFindManyArgs) SetCursor(cursor UniquePredicate[Trash]) *TrashFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *TrashFindManyArgs) SetSkip(n int) *TrashFindManyArgs {
	a.Skip = &n
	return a
}

func (a *TrashFindManyArgs) SetTake(n int) *TrashFindManyArgs {
	a.Take = &n
	return a
}

// TrashCountArgs is the input argument passed to Trash Count extension hooks.
type TrashCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Trash]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *TrashCountArgs) SetWhere(preds ...PredicateOf[Trash]) *TrashCountArgs {
	a.Where = preds
	return a
}

func (a *TrashCountArgs) SetSkip(n int) *TrashCountArgs {
	a.Skip = &n
	return a
}

func (a *TrashCountArgs) SetTake(n int) *TrashCountArgs {
	a.Take = &n
	return a
}

// TrashDeleteArgs is the input argument passed to Trash Delete extension hooks.
type TrashDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Trash]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *TrashSelect
}

func (a *TrashDeleteArgs) SetWhere(unique UniquePredicate[Trash], additional ...PredicateOf[Trash]) *TrashDeleteArgs {
	a.Where = make([]PredicateOf[Trash], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashDeleteManyArgs is the input argument passed to Trash DeleteMany extension hooks.
type TrashDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Trash]
}

func (a *TrashDeleteManyArgs) SetWhere(preds ...PredicateOf[Trash]) *TrashDeleteManyArgs {
	a.Where = preds
	return a
}

// TrashUpdateArgs is the input argument passed to Trash Update extension hooks.
type TrashUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Trash]
	// Data contains the model fields to update.
	Data *TrashUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *TrashSelect
}

func (a *TrashUpdateArgs) SetWhere(unique UniquePredicate[Trash], additional ...PredicateOf[Trash]) *TrashUpdateArgs {
	a.Where = make([]PredicateOf[Trash], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashUpdateManyArgs is the input argument passed to Trash UpdateMany extension hooks.
type TrashUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Trash]
	// Data contains the model fields to update.
	Data *TrashUpdate
}

func (a *TrashUpdateManyArgs) SetWhere(preds ...PredicateOf[Trash]) *TrashUpdateManyArgs {
	a.Where = preds
	return a
}

// TrashUpdateManyAndReturnArgs is the input argument passed to Trash UpdateManyAndReturn extension hooks.
type TrashUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Trash]
	// Data contains the model fields to update.
	Data *TrashUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *TrashSelect
}

func (a *TrashUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[Trash]) *TrashUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type TrashCreateQuery = func(ctx context.Context, args *TrashCreateArgs) (*Trash, error)
type TrashCreateManyQuery = func(ctx context.Context, args *TrashCreateManyArgs) (int64, error)
type TrashCreateManyAndReturnQuery = func(ctx context.Context, args *TrashCreateManyAndReturnArgs) ([]*Trash, error)
type TrashFindUniqueQuery = func(ctx context.Context, args *TrashFindUniqueArgs) (*Trash, error)
type TrashFindFirstQuery = func(ctx context.Context, args *TrashFindFirstArgs) (*Trash, error)
type TrashFindManyQuery = func(ctx context.Context, args *TrashFindManyArgs) ([]*Trash, error)
type TrashDeleteQuery = func(ctx context.Context, args *TrashDeleteArgs) (*Trash, error)
type TrashDeleteManyQuery = func(ctx context.Context, args *TrashDeleteManyArgs) (int64, error)
type TrashCountQuery = func(ctx context.Context, args *TrashCountArgs) (int64, error)
type TrashUpdateQuery = func(ctx context.Context, args *TrashUpdateArgs) (*Trash, error)
type TrashUpdateManyQuery = func(ctx context.Context, args *TrashUpdateManyArgs) (int64, error)
type TrashUpdateManyAndReturnQuery = func(ctx context.Context, args *TrashUpdateManyAndReturnArgs) ([]*Trash, error)

type TrashExtension struct {
	Create              func(ctx context.Context, args *TrashCreateArgs, next TrashCreateQuery) (*Trash, error)
	CreateMany          func(ctx context.Context, args *TrashCreateManyArgs, next TrashCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *TrashCreateManyAndReturnArgs, next TrashCreateManyAndReturnQuery) ([]*Trash, error)
	FindUnique          func(ctx context.Context, args *TrashFindUniqueArgs, next TrashFindUniqueQuery) (*Trash, error)
	FindFirst           func(ctx context.Context, args *TrashFindFirstArgs, next TrashFindFirstQuery) (*Trash, error)
	FindMany            func(ctx context.Context, args *TrashFindManyArgs, next TrashFindManyQuery) ([]*Trash, error)
	Delete              func(ctx context.Context, args *TrashDeleteArgs, next TrashDeleteQuery) (*Trash, error)
	DeleteMany          func(ctx context.Context, args *TrashDeleteManyArgs, next TrashDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *TrashCountArgs, next TrashCountQuery) (int64, error)
	Update              func(ctx context.Context, args *TrashUpdateArgs, next TrashUpdateQuery) (*Trash, error)
	UpdateMany          func(ctx context.Context, args *TrashUpdateManyArgs, next TrashUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *TrashUpdateManyAndReturnArgs, next TrashUpdateManyAndReturnQuery) ([]*Trash, error)
}

type TrashDelegate struct {
	client     *Queries
	extensions []TrashExtension
}

func (d *TrashDelegate) Use(exts ...TrashExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Trash) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "createdAt":
			targets[i] = &m.CreatedAt
		case "content":
			targets[i] = &m.Content
		}
	}
	return targets
}

var trashDefaultCols = []string{
	"id",
	"createdAt",
	"content",
}

var trashPKCols = []string{
	"id",
}

var trashUniqueCols = []string{
	"id",
}

func selectTrashCols(selects *TrashSelect, omits *TrashOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return trashDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.CreatedAt || selects.Content || selects.Agents != nil || selects.Skills != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"createdAt", selects != nil && selects.CreatedAt, omits != nil && omits.CreatedAt, false},
		{"content", selects != nil && selects.Content, omits != nil && omits.Content, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *TrashSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Agents != nil || s.Skills != nil
}

type TrashUpsertBuilder struct {
	*CreateBuilder[Trash, TrashSelect, TrashOmit]
}

func (b *TrashUpsertBuilder) Select(s TrashSelect) *TrashUpsertBuilder {
	b.selects = &s
	return b
}

func (b *TrashUpsertBuilder) Omit(o TrashOmit) *TrashUpsertBuilder {
	b.omits = &o
	return b
}

type TrashCreateBuilder struct {
	*CreateBuilder[Trash, TrashSelect, TrashOmit]
}

func (b *TrashCreateBuilder) Select(s TrashSelect) *TrashCreateBuilder {
	b.selects = &s
	return b
}

func (b *TrashCreateBuilder) Omit(o TrashOmit) *TrashCreateBuilder {
	b.omits = &o
	return b
}

func (b *TrashCreateBuilder) OnConflict(target UniqueConstraintTarget) *TrashConflictBuilder[TrashUpsertBuilder] {
	upsertBuilder := &TrashUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &TrashConflictBuilder[TrashUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *TrashCreateBuilder) SetId(v int32) *TrashCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *TrashCreateBuilder) SetCreatedAt(v time.Time) *TrashCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}
func (b *TrashCreateBuilder) SetContent(v string) *TrashCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *TrashCreateBuilder) Assignments(assignments ...FieldAssignmentOf[Trash]) *TrashCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *TrashDelegate) Create() *TrashCreateBuilder {
	return &TrashCreateBuilder{
		CreateBuilder: &CreateBuilder[Trash, TrashSelect, TrashOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedTrashId        uint64 = 1 << 0
	providedTrashCreatedAt uint64 = 1 << 1
	providedTrashContent   uint64 = 1 << 2
)

func assignmentsToTrashCreate(assignments []FieldAssignment) (TrashCreate, error) {
	var input TrashCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedTrashId
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
				errs.ValidateInt32("id", v, "")
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "createdAt":
			provided |= providedTrashCreatedAt
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			} else {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		case "content":
			provided |= providedTrashContent
			if v, ok := a.Val.(string); ok {
				input.Content = v
				errs.ValidateString("content", v, true, 0, false, false)
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		}
	}
	if provided&providedTrashContent == 0 {
		errs.Add("content", "", "required", "field Content is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *TrashCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 3)
	vals = make([]any, 0, 3)
	if s.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, *s.Id)
	}
	cols = append(cols, "createdAt")
	if s.CreatedAt != nil {
		vals = append(vals, *s.CreatedAt)
	} else {
		vals = append(vals, time.Now())
	}
	cols = append(cols, "content")
	vals = append(vals, s.Content)
	return
}

func partitionTrashInputs(dialect Dialect, inputs []*TrashCreate) [][]*TrashCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*TrashCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*TrashCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*TrashCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*TrashCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*TrashCreate{inputs}
}

func (d *TrashDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *TrashSelect, omits *TrashOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Trash, error) {
	input, err := assignmentsToTrashCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectTrashCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashCreateArgs) (*Trash, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectTrashCols(a.Select, omits)
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
			curr = func(c context.Context, a *TrashCreateArgs) (*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type TrashCreateManyBuilder struct {
	*CreateManyBuilder[Trash]
}

func (b *TrashCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *TrashConflictBuilder[TrashCreateManyBuilder] {
	return &TrashConflictBuilder[TrashCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type TrashCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Trash, TrashSelect, TrashOmit]
}

func (b *TrashCreateManyAndReturnBuilder) Select(s TrashSelect) *TrashCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *TrashCreateManyAndReturnBuilder) Omit(o TrashOmit) *TrashCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *TrashCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *TrashConflictBuilder[TrashCreateManyAndReturnBuilder] {
	return &TrashConflictBuilder[TrashCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToTrashRecordInputs(builders []*TrashCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *TrashDelegate) CreateMany(builders ...*TrashCreateBuilder) *TrashCreateManyBuilder {
	return &TrashCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Trash]{
			records:  createBuildersToTrashRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *TrashDelegate) CreateManyAndReturn(builders ...*TrashCreateBuilder) *TrashCreateManyAndReturnBuilder {
	return &TrashCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Trash, TrashSelect, TrashOmit]{
			records:  createBuildersToTrashRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToTrashCreateInputs(records []RecordInput) ([]*TrashCreate, error) {
	structs := make([]TrashCreate, len(records))
	inputs := make([]*TrashCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToTrashCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *TrashDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToTrashCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &TrashCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *TrashSelect, omits *TrashOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Trash, error) {
	inputs, err := recordsToTrashCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashCreateManyAndReturnArgs) ([]*Trash, error) {
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
			curr = func(c context.Context, a *TrashCreateManyAndReturnArgs) ([]*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *TrashSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Trash, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *Trash
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Trash.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.Trash.loadRelations(ctx, []*Trash{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "Trash", cols, returningCols, trashPKCols, conflictTarget, conflictAction, len(vals))
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

		var res Trash
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, trashPKCols)
}

func (d *TrashDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*Trash, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Trash") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Trash")
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

	var res Trash
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *TrashDelegate) buildBulkInsertSQL(q *Queries, batch []*TrashCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 3)
	for i, c := range trashDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Trash")
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
				if input.Id != nil {
					vals = append(vals, *input.Id)
				} else {
					writeDefault = true
				}
			case "createdAt":
				if input.CreatedAt != nil {
					vals = append(vals, *input.CreatedAt)
				} else {
					vals = append(vals, time.Now())
				}
			case "content":
				vals = append(vals, input.Content)
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

func applyTrashConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanTrashRows(rows *sql.Rows, returningCols []string) ([]*Trash, error) {
	var records []*Trash
	for rows.Next() {
		var res Trash
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

func (d *TrashDelegate) runCreateMany(ctx context.Context, inputs []*TrashCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionTrashInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashConflictClause(d.client.dialect, queryStr, vals, cols, trashPKCols, conflictTarget, conflictAction)

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

func (d *TrashDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*TrashCreate,
	selects *TrashSelect,
	omits *TrashOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Trash, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionTrashInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*Trash
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.Trash.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.Trash.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectTrashCols(selects, omits, trashPKCols...)
	recordsOut := make([]*Trash, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashConflictClause(d.client.dialect, queryStr, vals, cols, trashPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanTrashRows(rows, returningCols)
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

func (d *TrashDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*TrashCreate,
	selects *TrashSelect,
	omits *TrashOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Trash, error) {
	batches := partitionTrashInputs(d.client.dialect, inputs)
	returningCols := selectTrashCols(selects, omits, trashPKCols...)
	recordsOut := make([]*Trash, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashConflictClause(d.client.dialect, queryStr, vals, cols, trashPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("Trash") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "Trash")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, trashPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, trashPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanTrashRows(rows, returningCols)
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

type TrashConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *TrashConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *TrashConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *TrashConflictBuilder[B]) Update(fn func(u *TrashUpsert)) *B {
	cb.setAction(*TrashConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// TrashConflictUpdate creates a custom ConflictAction for Trash upsert conflicts.
func TrashConflictUpdate(fn func(u *TrashUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newTrashUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type TrashUpsert struct {
	Id        numericFieldUpsert[int32]
	CreatedAt fieldUpsert[time.Time]
	Content   fieldUpsert[string]
}

func newTrashUpsert(up *ConflictUpdate) *TrashUpsert {
	return &TrashUpsert{
		Id: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "id", update: up},
			tableName:   "Trash",
		},
		CreatedAt: fieldUpsert[time.Time]{column: "createdAt", update: up},
		Content:   fieldUpsert[string]{column: "content", update: up},
	}
}

type TrashUpdateBuilder struct {
	*UpdateBuilder[Trash, TrashSelect, TrashOmit]
}

type TrashUpdateManyBuilder struct {
	*UpdateManyBuilder[Trash]
}

type TrashUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[Trash, TrashSelect, TrashOmit]
}

func (b *TrashUpdateBuilder) SetId(v int32) *TrashUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *TrashUpdateManyBuilder) SetId(v int32) *TrashUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *TrashUpdateManyAndReturnBuilder) SetId(v int32) *TrashUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *TrashUpdateBuilder) SetCreatedAt(v time.Time) *TrashUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}

func (b *TrashUpdateManyBuilder) SetCreatedAt(v time.Time) *TrashUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}

func (b *TrashUpdateManyAndReturnBuilder) SetCreatedAt(v time.Time) *TrashUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}
func (b *TrashUpdateBuilder) SetContent(v string) *TrashUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *TrashUpdateManyBuilder) SetContent(v string) *TrashUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *TrashUpdateManyAndReturnBuilder) SetContent(v string) *TrashUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}

func (b *TrashUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[Trash]) *TrashUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *TrashUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[Trash]) *TrashUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *TrashUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[Trash]) *TrashUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *TrashDelegate) Update(where UniquePredicate[Trash], additional ...PredicateOf[Trash]) *TrashUpdateBuilder {
	return &TrashUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[Trash, TrashSelect, TrashOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *TrashDelegate) UpdateMany(preds ...PredicateOf[Trash]) *TrashUpdateManyBuilder {
	return &TrashUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[Trash]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *TrashDelegate) UpdateManyAndReturn(preds ...PredicateOf[Trash]) *TrashUpdateManyAndReturnBuilder {
	return &TrashUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[Trash, TrashSelect, TrashOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *TrashDelegate) buildUpdateSQL(preds []PredicateOf[Trash], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "Trash")
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

func (d *TrashDelegate) executeUpdate(ctx context.Context, where UniquePredicate[Trash], additional []PredicateOf[Trash], assignments []FieldAssignment, selects *TrashSelect, omits *TrashOmit) (*Trash, error) {
	allWhere := make([]PredicateOf[Trash], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToTrashUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashUpdateArgs) (*Trash, error) {
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
			curr = func(c context.Context, a *TrashUpdateArgs) (*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) runUpdate(ctx context.Context, preds []PredicateOf[Trash], cols []string, vals []any, selects *TrashSelect, omits *TrashOmit) (*Trash, error) {
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
		var res *Trash
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Trash.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Trash.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectTrashCols(selects, omits, trashPKCols...)
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
		return nil, &NotFoundError{Model: "Trash"}
	}

	var res Trash
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Trash{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *TrashDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[Trash], cols []string, vals []any) (int64, error) {
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

func (d *TrashDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[Trash], cols []string, vals []any, selects *TrashSelect, omits *TrashOmit) (*Trash, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "Trash"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *TrashDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[Trash], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToTrashUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &TrashUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *TrashUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *TrashDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Trash], assignments []FieldAssignment, selects *TrashSelect, omits *TrashOmit) ([]*Trash, error) {
	input, err := assignmentsToTrashUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashUpdateManyAndReturnArgs) ([]*Trash, error) {
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
			curr = func(c context.Context, a *TrashUpdateManyAndReturnArgs) ([]*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Trash], cols []string, vals []any, selects *TrashSelect, omits *TrashOmit) ([]*Trash, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[Trash]{Where: preds}, selects, omits)
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
		var res []*Trash
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Trash.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Trash.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectTrashCols(selects, omits, trashPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanTrashRows(rows, returningCols)
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

func (d *TrashDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[Trash], cols []string, vals []any, selects *TrashSelect, omits *TrashOmit) ([]*Trash, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*Trash{}, nil
	}
	return d.runFindMany(ctx, QueryParams[Trash]{Where: preds}, selects, omits)
}
func (d *TrashDelegate) FindUnique(where UniquePredicate[Trash], additional ...PredicateOf[Trash]) *FindUniqueBuilder[Trash, TrashSelect, TrashOmit] {
	return &FindUniqueBuilder[Trash, TrashSelect, TrashOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *TrashDelegate) FindFirst(preds ...PredicateOf[Trash]) *FindFirstBuilder[Trash, TrashSelect, TrashOmit] {
	return &FindFirstBuilder[Trash, TrashSelect, TrashOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *TrashDelegate) FindMany(preds ...PredicateOf[Trash]) *FindManyBuilder[Trash, TrashSelect, TrashOmit] {
	return &FindManyBuilder[Trash, TrashSelect, TrashOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *TrashDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Trash], additional []PredicateOf[Trash], selects *TrashSelect, omits *TrashOmit) (*Trash, error) {
	allWhere := make([]PredicateOf[Trash], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashFindUniqueArgs) (*Trash, error) {
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
			curr = func(c context.Context, a *TrashFindUniqueArgs) (*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Trash],
	selects *TrashSelect,
	omits *TrashOmit,
) (*Trash, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *TrashFindFirstArgs) (*Trash, error) {
		return d.runFindFirst(c, QueryParams[Trash]{
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
			curr = func(c context.Context, a *TrashFindFirstArgs) (*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Trash],
	selects *TrashSelect,
	omits *TrashOmit,
) ([]*Trash, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *TrashFindManyArgs) ([]*Trash, error) {
		return d.runFindMany(c, QueryParams[Trash]{
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
			curr = func(c context.Context, a *TrashFindManyArgs) ([]*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) runFindUnique(ctx context.Context, where []PredicateOf[Trash], selects *TrashSelect, omits *TrashOmit) (*Trash, error) {
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
	returningCols := selectTrashCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Trash{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *TrashDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Trash],
	selects *TrashSelect,
	omits *TrashOmit,
) (*Trash, error) {
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

func (d *TrashDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Trash],
	selects *TrashSelect,
	omits *TrashOmit,
) ([]*Trash, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, trashPKCols, trashUniqueCols, "Trash", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, trashPKCols, trashUniqueCols, isCursorQuery, params.Take)
	returningCols := selectTrashCols(selects, omits)

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

func (d *TrashDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*Trash, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Trash", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res Trash
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *TrashDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Trash, error) {
	query := buildSelectSQL(d.client, "Trash", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*Trash, 0)
	for rows.Next() {
		var res Trash
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
func (d *TrashDelegate) DeleteMany(preds ...PredicateOf[Trash]) *DeleteManyBuilder[Trash] {
	return &DeleteManyBuilder[Trash]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *TrashDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Trash]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &TrashDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *TrashDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Trash]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Trash")
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

func (d *TrashDelegate) Delete(where UniquePredicate[Trash], additional ...PredicateOf[Trash]) *DeleteBuilder[Trash, TrashSelect, TrashOmit] {
	return &DeleteBuilder[Trash, TrashSelect, TrashOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *TrashDelegate) executeDelete(ctx context.Context, where UniquePredicate[Trash], additional []PredicateOf[Trash], selects *TrashSelect, omits *TrashOmit) (*Trash, error) {
	allWhere := make([]PredicateOf[Trash], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashSelect()
	}

	args := &TrashDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashDeleteArgs) (*Trash, error) {
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
			curr = func(c context.Context, a *TrashDeleteArgs) (*Trash, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) runDelete(ctx context.Context, where []PredicateOf[Trash], selects *TrashSelect, omits *TrashOmit) (*Trash, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectTrashCols(selects, omits, trashPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Trash
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Trash.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "Trash"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Trash")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Trash]
			pkPreds = append(pkPreds, Predicate[Trash]{
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
	d.client.dialect.WriteQuotedIdent(&sb, "Trash")

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
		return nil, &NotFoundError{Model: "Trash"}
	}

	var row Trash
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *TrashDelegate) Count(preds ...PredicateOf[Trash]) *CountBuilder[Trash] {
	return &CountBuilder[Trash]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *TrashDelegate) executeCount(ctx context.Context, params QueryParams[Trash]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &TrashCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *TrashCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[Trash]{
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
			curr = func(c context.Context, a *TrashCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashDelegate) runCount(ctx context.Context, params QueryParams[Trash]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Trash")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Trash")
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
func (d *TrashDelegate) loadRelations(ctx context.Context, records []*Trash, selects *TrashSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Agents != nil {
		relationSelects, relationOmits, relationParams := selects.Agents.GetRelationParams()
		returningCols := selectTrashAgentCols(relationSelects, relationOmits, "trashId")
		// Inverse holds the FK: TrashAgent.trashId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Trash) int32 { return p.Id }),
			"TrashAgent",
			"trashId",
			returningCols,
			scanInto(returningCols, (*TrashAgent).ScanFields),
			directKey(func(c *TrashAgent) int32 { return c.TrashId }),
			appendMany(func(p *Trash) *[]*TrashAgent { return &p.Agents }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading agents: %w", err)
		}
		if err := d.client.TrashAgent.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Skills != nil {
		relationSelects, relationOmits, relationParams := selects.Skills.GetRelationParams()
		returningCols := selectTrashSkillCols(relationSelects, relationOmits, "trashId")
		// Inverse holds the FK: TrashSkill.trashId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Trash) int32 { return p.Id }),
			"TrashSkill",
			"trashId",
			returningCols,
			scanInto(returningCols, (*TrashSkill).ScanFields),
			directKey(func(c *TrashSkill) int32 { return c.TrashId }),
			appendMany(func(p *Trash) *[]*TrashSkill { return &p.Skills }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading skills: %w", err)
		}
		if err := d.client.TrashSkill.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
