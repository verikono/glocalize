package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// TrashAgent represents the database model
type TrashAgent struct {
	TrashId int32  `db:"trashId" json:"trashId"`
	AgentId string `db:"agentId" json:"agentId"`
	Trash   *Trash `json:"trash,omitempty"`
	Agent   *Agent `json:"agent,omitempty"`
}

// TrashAgentCreate contains model input fields for TrashAgent creation operations.
//
// Fields for TrashAgent:
//
//	trashId int32  required
//	agentId string required
type TrashAgentCreate struct {
	TrashId int32  `json:"trashId"`
	AgentId string `json:"agentId"`
}

// colMask returns a bit mask of columns that are set
func (s *TrashAgentCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// TrashAgentUpdate contains model input fields for TrashAgent update operations.
type TrashAgentUpdate struct {
	TrashId *int32  `json:"trashId"`
	AgentId *string `json:"agentId"`
}

func (u *TrashAgentUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.TrashId != nil {
		cols = append(cols, "trashId")
		vals = append(vals, u.TrashId)
	}
	if u.AgentId != nil {
		cols = append(cols, "agentId")
		vals = append(vals, u.AgentId)
	}
	return cols, vals
}

func assignmentsToTrashAgentUpdate(assignments []FieldAssignment) (TrashAgentUpdate, error) {
	var input TrashAgentUpdate
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
		case "agentId":
			if v, ok := a.Val.(string); ok {
				input.AgentId = &v
				errs.ValidateString("agentId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.AgentId = v
			} else {
				errs.Add("agentId", a.Val, "type", "field agentId must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// TrashAgentSelect specifies which scalar and relation fields to select for TrashAgent.
//
// Selectable fields:
//
//	-- Scalars --
//	trashId (bool)
//	agentId (bool)
//	-- Relations --
//	trash   (Trash)
//	agent   (Agent)
type TrashAgentSelect struct {
	TrashId bool         `json:"trashId"`
	AgentId bool         `json:"agentId"`
	Trash   *TrashSelect `json:"trash,omitempty"`
	Agent   *AgentSelect `json:"agent,omitempty"`
}

var fullTrashAgentSelectVal = &TrashAgentSelect{
	TrashId: true,
	AgentId: true,
}

func fullTrashAgentSelect() *TrashAgentSelect {
	return fullTrashAgentSelectVal
}

func (s *TrashAgentSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.TrashId || s.AgentId
}

func (s *TrashAgentSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type TrashAgentOmit struct {
	TrashId bool `json:"trashId"`
	AgentId bool `json:"agentId"`
}

type TrashAgentSelectQuery interface {
	GetRelationParams() (*TrashAgentSelect, *TrashAgentOmit, QueryParams[TrashAgent])
}

func (s *TrashAgentSelect) GetRelationParams() (*TrashAgentSelect, *TrashAgentOmit, QueryParams[TrashAgent]) {
	return s, nil, QueryParams[TrashAgent]{}
}

type TrashAgentQueryBuilder struct {
	selects *TrashAgentSelect
	omits   *TrashAgentOmit
	where   []PredicateOf[TrashAgent]
	take    *int
	skip    *int
	orderBy []OrderBy[TrashAgent]
	cursor  UniquePredicate[TrashAgent]
}

func (b *TrashAgentQueryBuilder) Where(preds ...PredicateOf[TrashAgent]) *TrashAgentQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *TrashAgentQueryBuilder) Take(limit int) *TrashAgentQueryBuilder {
	b.take = &limit
	return b
}

func (b *TrashAgentQueryBuilder) Skip(offset int) *TrashAgentQueryBuilder {
	b.skip = &offset
	return b
}

func (b *TrashAgentQueryBuilder) OrderBy(orders ...OrderBy[TrashAgent]) *TrashAgentQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *TrashAgentQueryBuilder) Cursor(where UniquePredicate[TrashAgent]) *TrashAgentQueryBuilder {
	b.cursor = where
	return b
}

func (b *TrashAgentQueryBuilder) Select(s TrashAgentSelect) *TrashAgentQueryBuilder {
	b.selects = &s
	return b
}

func (b *TrashAgentQueryBuilder) Omit(o TrashAgentOmit) *TrashAgentQueryBuilder {
	b.omits = &o
	return b
}

func (b *TrashAgentQueryBuilder) GetRelationParams() (*TrashAgentSelect, *TrashAgentOmit, QueryParams[TrashAgent]) {
	if b == nil {
		return nil, nil, QueryParams[TrashAgent]{}
	}
	return b.selects, b.omits, QueryParams[TrashAgent]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// TrashAgentCreateArgs is the input argument passed to TrashAgent Create extension hooks.
//
// Fields for TrashAgent:
//
//	trashId int32  required
//	agentId string required
//
// Relations for TrashAgent:
//
//	trash   (Trash)
//	agent   (Agent)
type TrashAgentCreateArgs struct {
	// Data contains the model fields to insert.
	Data *TrashAgentCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *TrashAgentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// TrashAgentCreateManyArgs is the input argument passed to TrashAgent CreateMany extension hooks.
//
// Fields for TrashAgent:
//
//	trashId int32  required
//	agentId string required
type TrashAgentCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*TrashAgentCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *TrashAgentCreateManyArgs) AppendData(builders ...*TrashAgentCreateBuilder) *TrashAgentCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToTrashAgentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// TrashAgentCreateManyAndReturnArgs is the input argument passed to TrashAgent CreateManyAndReturn extension hooks.
//
// Fields for TrashAgent:
//
//	trashId int32  required
//	agentId string required
//
// Relations for TrashAgent:
//
//	trash   (Trash)
//	agent   (Agent)
type TrashAgentCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*TrashAgentCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *TrashAgentSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *TrashAgentCreateManyAndReturnArgs) AppendData(builders ...*TrashAgentCreateBuilder) *TrashAgentCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToTrashAgentCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// TrashAgentFindUniqueArgs is the input argument passed to TrashAgent FindUnique extension hooks.
type TrashAgentFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[TrashAgent]
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashAgentSelect
}

func (a *TrashAgentFindUniqueArgs) SetWhere(unique UniquePredicate[TrashAgent], additional ...PredicateOf[TrashAgent]) *TrashAgentFindUniqueArgs {
	a.Where = make([]PredicateOf[TrashAgent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashAgentFindFirstArgs is the input argument passed to TrashAgent FindFirst extension hooks.
type TrashAgentFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashAgent]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[TrashAgent]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[TrashAgent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashAgentSelect
}

func (a *TrashAgentFindFirstArgs) SetWhere(preds ...PredicateOf[TrashAgent]) *TrashAgentFindFirstArgs {
	a.Where = preds
	return a
}

func (a *TrashAgentFindFirstArgs) SetOrderBy(orders ...OrderBy[TrashAgent]) *TrashAgentFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *TrashAgentFindFirstArgs) SetCursor(cursor UniquePredicate[TrashAgent]) *TrashAgentFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *TrashAgentFindFirstArgs) SetSkip(n int) *TrashAgentFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *TrashAgentFindFirstArgs) SetTake(n int) *TrashAgentFindFirstArgs {
	a.Take = &n
	return a
}

// TrashAgentFindManyArgs is the input argument passed to TrashAgent FindMany extension hooks.
type TrashAgentFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashAgent]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[TrashAgent]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[TrashAgent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *TrashAgentSelect
}

func (a *TrashAgentFindManyArgs) SetWhere(preds ...PredicateOf[TrashAgent]) *TrashAgentFindManyArgs {
	a.Where = preds
	return a
}

func (a *TrashAgentFindManyArgs) SetOrderBy(orders ...OrderBy[TrashAgent]) *TrashAgentFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *TrashAgentFindManyArgs) SetCursor(cursor UniquePredicate[TrashAgent]) *TrashAgentFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *TrashAgentFindManyArgs) SetSkip(n int) *TrashAgentFindManyArgs {
	a.Skip = &n
	return a
}

func (a *TrashAgentFindManyArgs) SetTake(n int) *TrashAgentFindManyArgs {
	a.Take = &n
	return a
}

// TrashAgentCountArgs is the input argument passed to TrashAgent Count extension hooks.
type TrashAgentCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashAgent]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *TrashAgentCountArgs) SetWhere(preds ...PredicateOf[TrashAgent]) *TrashAgentCountArgs {
	a.Where = preds
	return a
}

func (a *TrashAgentCountArgs) SetSkip(n int) *TrashAgentCountArgs {
	a.Skip = &n
	return a
}

func (a *TrashAgentCountArgs) SetTake(n int) *TrashAgentCountArgs {
	a.Take = &n
	return a
}

// TrashAgentDeleteArgs is the input argument passed to TrashAgent Delete extension hooks.
type TrashAgentDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[TrashAgent]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *TrashAgentSelect
}

func (a *TrashAgentDeleteArgs) SetWhere(unique UniquePredicate[TrashAgent], additional ...PredicateOf[TrashAgent]) *TrashAgentDeleteArgs {
	a.Where = make([]PredicateOf[TrashAgent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashAgentDeleteManyArgs is the input argument passed to TrashAgent DeleteMany extension hooks.
type TrashAgentDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashAgent]
}

func (a *TrashAgentDeleteManyArgs) SetWhere(preds ...PredicateOf[TrashAgent]) *TrashAgentDeleteManyArgs {
	a.Where = preds
	return a
}

// TrashAgentUpdateArgs is the input argument passed to TrashAgent Update extension hooks.
type TrashAgentUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[TrashAgent]
	// Data contains the model fields to update.
	Data *TrashAgentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *TrashAgentSelect
}

func (a *TrashAgentUpdateArgs) SetWhere(unique UniquePredicate[TrashAgent], additional ...PredicateOf[TrashAgent]) *TrashAgentUpdateArgs {
	a.Where = make([]PredicateOf[TrashAgent], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// TrashAgentUpdateManyArgs is the input argument passed to TrashAgent UpdateMany extension hooks.
type TrashAgentUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashAgent]
	// Data contains the model fields to update.
	Data *TrashAgentUpdate
}

func (a *TrashAgentUpdateManyArgs) SetWhere(preds ...PredicateOf[TrashAgent]) *TrashAgentUpdateManyArgs {
	a.Where = preds
	return a
}

// TrashAgentUpdateManyAndReturnArgs is the input argument passed to TrashAgent UpdateManyAndReturn extension hooks.
type TrashAgentUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[TrashAgent]
	// Data contains the model fields to update.
	Data *TrashAgentUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *TrashAgentSelect
}

func (a *TrashAgentUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[TrashAgent]) *TrashAgentUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type TrashAgentCreateQuery = func(ctx context.Context, args *TrashAgentCreateArgs) (*TrashAgent, error)
type TrashAgentCreateManyQuery = func(ctx context.Context, args *TrashAgentCreateManyArgs) (int64, error)
type TrashAgentCreateManyAndReturnQuery = func(ctx context.Context, args *TrashAgentCreateManyAndReturnArgs) ([]*TrashAgent, error)
type TrashAgentFindUniqueQuery = func(ctx context.Context, args *TrashAgentFindUniqueArgs) (*TrashAgent, error)
type TrashAgentFindFirstQuery = func(ctx context.Context, args *TrashAgentFindFirstArgs) (*TrashAgent, error)
type TrashAgentFindManyQuery = func(ctx context.Context, args *TrashAgentFindManyArgs) ([]*TrashAgent, error)
type TrashAgentDeleteQuery = func(ctx context.Context, args *TrashAgentDeleteArgs) (*TrashAgent, error)
type TrashAgentDeleteManyQuery = func(ctx context.Context, args *TrashAgentDeleteManyArgs) (int64, error)
type TrashAgentCountQuery = func(ctx context.Context, args *TrashAgentCountArgs) (int64, error)
type TrashAgentUpdateQuery = func(ctx context.Context, args *TrashAgentUpdateArgs) (*TrashAgent, error)
type TrashAgentUpdateManyQuery = func(ctx context.Context, args *TrashAgentUpdateManyArgs) (int64, error)
type TrashAgentUpdateManyAndReturnQuery = func(ctx context.Context, args *TrashAgentUpdateManyAndReturnArgs) ([]*TrashAgent, error)

type TrashAgentExtension struct {
	Create              func(ctx context.Context, args *TrashAgentCreateArgs, next TrashAgentCreateQuery) (*TrashAgent, error)
	CreateMany          func(ctx context.Context, args *TrashAgentCreateManyArgs, next TrashAgentCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *TrashAgentCreateManyAndReturnArgs, next TrashAgentCreateManyAndReturnQuery) ([]*TrashAgent, error)
	FindUnique          func(ctx context.Context, args *TrashAgentFindUniqueArgs, next TrashAgentFindUniqueQuery) (*TrashAgent, error)
	FindFirst           func(ctx context.Context, args *TrashAgentFindFirstArgs, next TrashAgentFindFirstQuery) (*TrashAgent, error)
	FindMany            func(ctx context.Context, args *TrashAgentFindManyArgs, next TrashAgentFindManyQuery) ([]*TrashAgent, error)
	Delete              func(ctx context.Context, args *TrashAgentDeleteArgs, next TrashAgentDeleteQuery) (*TrashAgent, error)
	DeleteMany          func(ctx context.Context, args *TrashAgentDeleteManyArgs, next TrashAgentDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *TrashAgentCountArgs, next TrashAgentCountQuery) (int64, error)
	Update              func(ctx context.Context, args *TrashAgentUpdateArgs, next TrashAgentUpdateQuery) (*TrashAgent, error)
	UpdateMany          func(ctx context.Context, args *TrashAgentUpdateManyArgs, next TrashAgentUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *TrashAgentUpdateManyAndReturnArgs, next TrashAgentUpdateManyAndReturnQuery) ([]*TrashAgent, error)
}

type TrashAgentDelegate struct {
	client     *Queries
	extensions []TrashAgentExtension
}

func (d *TrashAgentDelegate) Use(exts ...TrashAgentExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *TrashAgent) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "trashId":
			targets[i] = &m.TrashId
		case "agentId":
			targets[i] = &m.AgentId
		}
	}
	return targets
}

var trashAgentDefaultCols = []string{
	"trashId",
	"agentId",
}

var trashAgentPKCols = []string{
	"trashId",
	"agentId",
}

var trashAgentUniqueCols = []string{}

func selectTrashAgentCols(selects *TrashAgentSelect, omits *TrashAgentOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return trashAgentDefaultCols
	}

	anySelected := selects != nil && (selects.TrashId || selects.AgentId || selects.Trash != nil || selects.Agent != nil)

	specs := []colSpec{
		{"trashId", selects != nil && selects.TrashId, omits != nil && omits.TrashId, selects != nil && selects.Trash != nil},
		{"agentId", selects != nil && selects.AgentId, omits != nil && omits.AgentId, selects != nil && selects.Agent != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *TrashAgentSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Trash != nil || s.Agent != nil
}

type TrashAgentUpsertBuilder struct {
	*CreateBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]
}

func (b *TrashAgentUpsertBuilder) Select(s TrashAgentSelect) *TrashAgentUpsertBuilder {
	b.selects = &s
	return b
}

func (b *TrashAgentUpsertBuilder) Omit(o TrashAgentOmit) *TrashAgentUpsertBuilder {
	b.omits = &o
	return b
}

type TrashAgentCreateBuilder struct {
	*CreateBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]
}

func (b *TrashAgentCreateBuilder) Select(s TrashAgentSelect) *TrashAgentCreateBuilder {
	b.selects = &s
	return b
}

func (b *TrashAgentCreateBuilder) Omit(o TrashAgentOmit) *TrashAgentCreateBuilder {
	b.omits = &o
	return b
}

func (b *TrashAgentCreateBuilder) OnConflict(target UniqueConstraintTarget) *TrashAgentConflictBuilder[TrashAgentUpsertBuilder] {
	upsertBuilder := &TrashAgentUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &TrashAgentConflictBuilder[TrashAgentUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *TrashAgentCreateBuilder) SetTrashId(v int32) *TrashAgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}
func (b *TrashAgentCreateBuilder) SetAgentId(v string) *TrashAgentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *TrashAgentCreateBuilder) Assignments(assignments ...FieldAssignmentOf[TrashAgent]) *TrashAgentCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *TrashAgentDelegate) Create() *TrashAgentCreateBuilder {
	return &TrashAgentCreateBuilder{
		CreateBuilder: &CreateBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedTrashAgentTrashId uint64 = 1 << 0
	providedTrashAgentAgentId uint64 = 1 << 1
)

func assignmentsToTrashAgentCreate(assignments []FieldAssignment) (TrashAgentCreate, error) {
	var input TrashAgentCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "trashId":
			provided |= providedTrashAgentTrashId
			if v, ok := a.Val.(int32); ok {
				input.TrashId = v
				errs.ValidateInt32("trashId", v, "")
			} else {
				errs.Add("trashId", a.Val, "type", "field trashId must be of type int32")
			}
		case "agentId":
			provided |= providedTrashAgentAgentId
			if v, ok := a.Val.(string); ok {
				input.AgentId = v
				errs.ValidateString("agentId", v, true, 0, false, false)
			} else {
				errs.Add("agentId", a.Val, "type", "field agentId must be of type string")
			}
		}
	}
	if provided&providedTrashAgentTrashId == 0 {
		errs.Add("trashId", nil, "required", "field TrashId is required")
	}
	if provided&providedTrashAgentAgentId == 0 {
		errs.Add("agentId", "", "required", "field AgentId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *TrashAgentCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "trashId")
	vals = append(vals, s.TrashId)
	cols = append(cols, "agentId")
	vals = append(vals, s.AgentId)
	return
}

func partitionTrashAgentInputs(dialect Dialect, inputs []*TrashAgentCreate) [][]*TrashAgentCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*TrashAgentCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*TrashAgentCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*TrashAgentCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*TrashAgentCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*TrashAgentCreate{inputs}
}

func (d *TrashAgentDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *TrashAgentSelect, omits *TrashAgentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*TrashAgent, error) {
	input, err := assignmentsToTrashAgentCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectTrashAgentCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashAgentCreateArgs) (*TrashAgent, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectTrashAgentCols(a.Select, omits)
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
			curr = func(c context.Context, a *TrashAgentCreateArgs) (*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type TrashAgentCreateManyBuilder struct {
	*CreateManyBuilder[TrashAgent]
}

func (b *TrashAgentCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *TrashAgentConflictBuilder[TrashAgentCreateManyBuilder] {
	return &TrashAgentConflictBuilder[TrashAgentCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type TrashAgentCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]
}

func (b *TrashAgentCreateManyAndReturnBuilder) Select(s TrashAgentSelect) *TrashAgentCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *TrashAgentCreateManyAndReturnBuilder) Omit(o TrashAgentOmit) *TrashAgentCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *TrashAgentCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *TrashAgentConflictBuilder[TrashAgentCreateManyAndReturnBuilder] {
	return &TrashAgentConflictBuilder[TrashAgentCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToTrashAgentRecordInputs(builders []*TrashAgentCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *TrashAgentDelegate) CreateMany(builders ...*TrashAgentCreateBuilder) *TrashAgentCreateManyBuilder {
	return &TrashAgentCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[TrashAgent]{
			records:  createBuildersToTrashAgentRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *TrashAgentDelegate) CreateManyAndReturn(builders ...*TrashAgentCreateBuilder) *TrashAgentCreateManyAndReturnBuilder {
	return &TrashAgentCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
			records:  createBuildersToTrashAgentRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToTrashAgentCreateInputs(records []RecordInput) ([]*TrashAgentCreate, error) {
	structs := make([]TrashAgentCreate, len(records))
	inputs := make([]*TrashAgentCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToTrashAgentCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *TrashAgentDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToTrashAgentCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &TrashAgentCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashAgentCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashAgentCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *TrashAgentSelect, omits *TrashAgentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*TrashAgent, error) {
	inputs, err := recordsToTrashAgentCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *TrashAgentCreateManyAndReturnArgs) ([]*TrashAgent, error) {
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
			curr = func(c context.Context, a *TrashAgentCreateManyAndReturnArgs) ([]*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *TrashAgentSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*TrashAgent, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *TrashAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.TrashAgent.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.TrashAgent.loadRelations(ctx, []*TrashAgent{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "TrashAgent", cols, returningCols, trashAgentPKCols, conflictTarget, conflictAction, len(vals))
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

		var res TrashAgent
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, trashAgentPKCols)
}

func (d *TrashAgentDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*TrashAgent, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("TrashAgent") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "TrashAgent")
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

	var res TrashAgent
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *TrashAgentDelegate) buildBulkInsertSQL(q *Queries, batch []*TrashAgentCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range trashAgentDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "TrashAgent")
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
			case "agentId":
				vals = append(vals, input.AgentId)
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

func applyTrashAgentConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanTrashAgentRows(rows *sql.Rows, returningCols []string) ([]*TrashAgent, error) {
	var records []*TrashAgent
	for rows.Next() {
		var res TrashAgent
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

func (d *TrashAgentDelegate) runCreateMany(ctx context.Context, inputs []*TrashAgentCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionTrashAgentInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashAgentConflictClause(d.client.dialect, queryStr, vals, cols, trashAgentPKCols, conflictTarget, conflictAction)

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

func (d *TrashAgentDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*TrashAgentCreate,
	selects *TrashAgentSelect,
	omits *TrashAgentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*TrashAgent, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionTrashAgentInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*TrashAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.TrashAgent.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.TrashAgent.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectTrashAgentCols(selects, omits, trashAgentPKCols...)
	recordsOut := make([]*TrashAgent, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashAgentConflictClause(d.client.dialect, queryStr, vals, cols, trashAgentPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanTrashAgentRows(rows, returningCols)
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

func (d *TrashAgentDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*TrashAgentCreate,
	selects *TrashAgentSelect,
	omits *TrashAgentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*TrashAgent, error) {
	batches := partitionTrashAgentInputs(d.client.dialect, inputs)
	returningCols := selectTrashAgentCols(selects, omits, trashAgentPKCols...)
	recordsOut := make([]*TrashAgent, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyTrashAgentConflictClause(d.client.dialect, queryStr, vals, cols, trashAgentPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("TrashAgent") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "TrashAgent")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, trashAgentPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, trashAgentPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanTrashAgentRows(rows, returningCols)
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

type TrashAgentConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *TrashAgentConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *TrashAgentConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *TrashAgentConflictBuilder[B]) Update(fn func(u *TrashAgentUpsert)) *B {
	cb.setAction(*TrashAgentConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// TrashAgentConflictUpdate creates a custom ConflictAction for TrashAgent upsert conflicts.
func TrashAgentConflictUpdate(fn func(u *TrashAgentUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newTrashAgentUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type TrashAgentUpsert struct {
	TrashId numericFieldUpsert[int32]
	AgentId fieldUpsert[string]
}

func newTrashAgentUpsert(up *ConflictUpdate) *TrashAgentUpsert {
	return &TrashAgentUpsert{
		TrashId: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "trashId", update: up},
			tableName:   "TrashAgent",
		},
		AgentId: fieldUpsert[string]{column: "agentId", update: up},
	}
}

type TrashAgentUpdateBuilder struct {
	*UpdateBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]
}

type TrashAgentUpdateManyBuilder struct {
	*UpdateManyBuilder[TrashAgent]
}

type TrashAgentUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]
}

func (b *TrashAgentUpdateBuilder) SetTrashId(v int32) *TrashAgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}

func (b *TrashAgentUpdateManyBuilder) SetTrashId(v int32) *TrashAgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}

func (b *TrashAgentUpdateManyAndReturnBuilder) SetTrashId(v int32) *TrashAgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "trashId", Val: v})
	return b
}
func (b *TrashAgentUpdateBuilder) SetAgentId(v string) *TrashAgentUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *TrashAgentUpdateManyBuilder) SetAgentId(v string) *TrashAgentUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *TrashAgentUpdateManyAndReturnBuilder) SetAgentId(v string) *TrashAgentUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *TrashAgentUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[TrashAgent]) *TrashAgentUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *TrashAgentUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[TrashAgent]) *TrashAgentUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *TrashAgentUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[TrashAgent]) *TrashAgentUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *TrashAgentDelegate) Update(where UniquePredicate[TrashAgent], additional ...PredicateOf[TrashAgent]) *TrashAgentUpdateBuilder {
	return &TrashAgentUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *TrashAgentDelegate) UpdateMany(preds ...PredicateOf[TrashAgent]) *TrashAgentUpdateManyBuilder {
	return &TrashAgentUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[TrashAgent]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *TrashAgentDelegate) UpdateManyAndReturn(preds ...PredicateOf[TrashAgent]) *TrashAgentUpdateManyAndReturnBuilder {
	return &TrashAgentUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *TrashAgentDelegate) buildUpdateSQL(preds []PredicateOf[TrashAgent], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "TrashAgent")
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

func (d *TrashAgentDelegate) executeUpdate(ctx context.Context, where UniquePredicate[TrashAgent], additional []PredicateOf[TrashAgent], assignments []FieldAssignment, selects *TrashAgentSelect, omits *TrashAgentOmit) (*TrashAgent, error) {
	allWhere := make([]PredicateOf[TrashAgent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToTrashAgentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashAgentUpdateArgs) (*TrashAgent, error) {
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
			curr = func(c context.Context, a *TrashAgentUpdateArgs) (*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) runUpdate(ctx context.Context, preds []PredicateOf[TrashAgent], cols []string, vals []any, selects *TrashAgentSelect, omits *TrashAgentOmit) (*TrashAgent, error) {
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
		var res *TrashAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.TrashAgent.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.TrashAgent.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectTrashAgentCols(selects, omits, trashAgentPKCols...)
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
		return nil, &NotFoundError{Model: "TrashAgent"}
	}

	var res TrashAgent
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*TrashAgent{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *TrashAgentDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[TrashAgent], cols []string, vals []any) (int64, error) {
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

func (d *TrashAgentDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[TrashAgent], cols []string, vals []any, selects *TrashAgentSelect, omits *TrashAgentOmit) (*TrashAgent, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "TrashAgent"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *TrashAgentDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[TrashAgent], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToTrashAgentUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &TrashAgentUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *TrashAgentUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashAgentUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *TrashAgentDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[TrashAgent], assignments []FieldAssignment, selects *TrashAgentSelect, omits *TrashAgentOmit) ([]*TrashAgent, error) {
	input, err := assignmentsToTrashAgentUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashAgentUpdateManyAndReturnArgs) ([]*TrashAgent, error) {
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
			curr = func(c context.Context, a *TrashAgentUpdateManyAndReturnArgs) ([]*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[TrashAgent], cols []string, vals []any, selects *TrashAgentSelect, omits *TrashAgentOmit) ([]*TrashAgent, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[TrashAgent]{Where: preds}, selects, omits)
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
		var res []*TrashAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.TrashAgent.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.TrashAgent.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectTrashAgentCols(selects, omits, trashAgentPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanTrashAgentRows(rows, returningCols)
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

func (d *TrashAgentDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[TrashAgent], cols []string, vals []any, selects *TrashAgentSelect, omits *TrashAgentOmit) ([]*TrashAgent, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*TrashAgent{}, nil
	}
	return d.runFindMany(ctx, QueryParams[TrashAgent]{Where: preds}, selects, omits)
}
func (d *TrashAgentDelegate) FindUnique(where UniquePredicate[TrashAgent], additional ...PredicateOf[TrashAgent]) *FindUniqueBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit] {
	return &FindUniqueBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *TrashAgentDelegate) FindFirst(preds ...PredicateOf[TrashAgent]) *FindFirstBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit] {
	return &FindFirstBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *TrashAgentDelegate) FindMany(preds ...PredicateOf[TrashAgent]) *FindManyBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit] {
	return &FindManyBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *TrashAgentDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[TrashAgent], additional []PredicateOf[TrashAgent], selects *TrashAgentSelect, omits *TrashAgentOmit) (*TrashAgent, error) {
	allWhere := make([]PredicateOf[TrashAgent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashAgentFindUniqueArgs) (*TrashAgent, error) {
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
			curr = func(c context.Context, a *TrashAgentFindUniqueArgs) (*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[TrashAgent],
	selects *TrashAgentSelect,
	omits *TrashAgentOmit,
) (*TrashAgent, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *TrashAgentFindFirstArgs) (*TrashAgent, error) {
		return d.runFindFirst(c, QueryParams[TrashAgent]{
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
			curr = func(c context.Context, a *TrashAgentFindFirstArgs) (*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[TrashAgent],
	selects *TrashAgentSelect,
	omits *TrashAgentOmit,
) ([]*TrashAgent, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *TrashAgentFindManyArgs) ([]*TrashAgent, error) {
		return d.runFindMany(c, QueryParams[TrashAgent]{
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
			curr = func(c context.Context, a *TrashAgentFindManyArgs) ([]*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) runFindUnique(ctx context.Context, where []PredicateOf[TrashAgent], selects *TrashAgentSelect, omits *TrashAgentOmit) (*TrashAgent, error) {
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
	returningCols := selectTrashAgentCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*TrashAgent{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *TrashAgentDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[TrashAgent],
	selects *TrashAgentSelect,
	omits *TrashAgentOmit,
) (*TrashAgent, error) {
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

func (d *TrashAgentDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[TrashAgent],
	selects *TrashAgentSelect,
	omits *TrashAgentOmit,
) ([]*TrashAgent, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, trashAgentPKCols, trashAgentUniqueCols, "TrashAgent", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, trashAgentPKCols, trashAgentUniqueCols, isCursorQuery, params.Take)
	returningCols := selectTrashAgentCols(selects, omits)

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

func (d *TrashAgentDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*TrashAgent, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "TrashAgent", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res TrashAgent
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *TrashAgentDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*TrashAgent, error) {
	query := buildSelectSQL(d.client, "TrashAgent", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*TrashAgent, 0)
	for rows.Next() {
		var res TrashAgent
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
func (d *TrashAgentDelegate) DeleteMany(preds ...PredicateOf[TrashAgent]) *DeleteManyBuilder[TrashAgent] {
	return &DeleteManyBuilder[TrashAgent]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *TrashAgentDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[TrashAgent]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &TrashAgentDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *TrashAgentDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *TrashAgentDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[TrashAgent]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "TrashAgent")
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

func (d *TrashAgentDelegate) Delete(where UniquePredicate[TrashAgent], additional ...PredicateOf[TrashAgent]) *DeleteBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit] {
	return &DeleteBuilder[TrashAgent, TrashAgentSelect, TrashAgentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *TrashAgentDelegate) executeDelete(ctx context.Context, where UniquePredicate[TrashAgent], additional []PredicateOf[TrashAgent], selects *TrashAgentSelect, omits *TrashAgentOmit) (*TrashAgent, error) {
	allWhere := make([]PredicateOf[TrashAgent], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullTrashAgentSelect()
	}

	args := &TrashAgentDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *TrashAgentDeleteArgs) (*TrashAgent, error) {
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
			curr = func(c context.Context, a *TrashAgentDeleteArgs) (*TrashAgent, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) runDelete(ctx context.Context, where []PredicateOf[TrashAgent], selects *TrashAgentSelect, omits *TrashAgentOmit) (*TrashAgent, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectTrashAgentCols(selects, omits, trashAgentPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *TrashAgent
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.TrashAgent.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "TrashAgent"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "TrashAgent")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[TrashAgent]
			pkPreds = append(pkPreds, Predicate[TrashAgent]{
				Data: PredicateData{
					Column:   "trashId",
					Operator: "=",
					Value:    res.TrashId,
				},
			})
			pkPreds = append(pkPreds, Predicate[TrashAgent]{
				Data: PredicateData{
					Column:   "agentId",
					Operator: "=",
					Value:    res.AgentId,
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
	d.client.dialect.WriteQuotedIdent(&sb, "TrashAgent")

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
		return nil, &NotFoundError{Model: "TrashAgent"}
	}

	var row TrashAgent
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *TrashAgentDelegate) Count(preds ...PredicateOf[TrashAgent]) *CountBuilder[TrashAgent] {
	return &CountBuilder[TrashAgent]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *TrashAgentDelegate) executeCount(ctx context.Context, params QueryParams[TrashAgent]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &TrashAgentCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *TrashAgentCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[TrashAgent]{
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
			curr = func(c context.Context, a *TrashAgentCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *TrashAgentDelegate) runCount(ctx context.Context, params QueryParams[TrashAgent]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "TrashAgent")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "TrashAgent")
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
func (d *TrashAgentDelegate) loadRelations(ctx context.Context, records []*TrashAgent, selects *TrashAgentSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Trash != nil {
		relationSelects, relationOmits, relationParams := selects.Trash.GetRelationParams()
		returningCols := selectTrashCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: TrashAgent.trashId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *TrashAgent) int32 { return p.TrashId }),
			"Trash",
			"id",
			returningCols,
			scanInto(returningCols, (*Trash).ScanFields),
			directKey(func(c *Trash) int32 { return c.Id }),
			setOne(func(p *TrashAgent, c *Trash) { p.Trash = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading trash: %w", err)
		}
		if err := d.client.Trash.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Agent != nil {
		relationSelects, relationOmits, relationParams := selects.Agent.GetRelationParams()
		returningCols := selectAgentCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: TrashAgent.agentId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *TrashAgent) string { return p.AgentId }),
			"Agent",
			"id",
			returningCols,
			scanInto(returningCols, (*Agent).ScanFields),
			directKey(func(c *Agent) string { return c.Id }),
			setOne(func(p *TrashAgent, c *Agent) { p.Agent = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading agent: %w", err)
		}
		if err := d.client.Agent.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
