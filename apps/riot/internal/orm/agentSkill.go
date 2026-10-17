package riot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// AgentSkill represents the database model
type AgentSkill struct {
	AgentId  string       `db:"agentId" json:"agentId"`
	SkillId  string       `db:"skillId" json:"skillId"`
	Priority PriorityType `db:"priority" json:"priority"`
	Agent    *Agent       `json:"agent,omitempty"`
	Skill    *Skill       `json:"skill,omitempty"`
}

// AgentSkillCreate contains model input fields for AgentSkill creation operations.
//
// Fields for AgentSkill:
//
//	agentId  string   required
//	skillId  string   required
//	priority Priority default: NORMAL
type AgentSkillCreate struct {
	AgentId  string        `json:"agentId"`
	SkillId  string        `json:"skillId"`
	Priority *PriorityType `json:"priority"`
}

// colMask returns a bit mask of columns that are set
func (s *AgentSkillCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	if s.Priority != nil {
		mask |= 1 << 2
	}
	return mask
}

// AgentSkillUpdate contains model input fields for AgentSkill update operations.
type AgentSkillUpdate struct {
	AgentId  *string       `json:"agentId"`
	SkillId  *string       `json:"skillId"`
	Priority *PriorityType `json:"priority"`
}

func (u *AgentSkillUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.AgentId != nil {
		cols = append(cols, "agentId")
		vals = append(vals, u.AgentId)
	}
	if u.SkillId != nil {
		cols = append(cols, "skillId")
		vals = append(vals, u.SkillId)
	}
	if u.Priority != nil {
		cols = append(cols, "priority")
		vals = append(vals, u.Priority)
	}
	return cols, vals
}

func assignmentsToAgentSkillUpdate(assignments []FieldAssignment) (AgentSkillUpdate, error) {
	var input AgentSkillUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "agentId":
			if v, ok := a.Val.(string); ok {
				input.AgentId = &v
				errs.ValidateString("agentId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.AgentId = v
			} else {
				errs.Add("agentId", a.Val, "type", "field agentId must be of type string")
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
		case "priority":
			if v, ok := a.Val.(PriorityType); ok {
				input.Priority = &v
				if !v.IsValid() {
					errs.Add("priority", v, "enum", fmt.Sprintf("invalid enum value %q for field priority", v))
				}
			} else if v, ok := a.Val.(*PriorityType); ok {
				input.Priority = v
				if v != nil && !v.IsValid() {
					errs.Add("priority", *v, "enum", fmt.Sprintf("invalid enum value %q for field priority", *v))
				}
			} else {
				errs.Add("priority", a.Val, "type", "field priority must be of type PriorityType")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// AgentSkillSelect specifies which scalar and relation fields to select for AgentSkill.
//
// Selectable fields:
//
//	-- Scalars --
//	agentId  (bool)
//	skillId  (bool)
//	priority (bool)
//	-- Relations --
//	agent    (Agent)
//	skill    (Skill)
type AgentSkillSelect struct {
	AgentId  bool         `json:"agentId"`
	SkillId  bool         `json:"skillId"`
	Priority bool         `json:"priority"`
	Agent    *AgentSelect `json:"agent,omitempty"`
	Skill    *SkillSelect `json:"skill,omitempty"`
}

var fullAgentSkillSelectVal = &AgentSkillSelect{
	AgentId:  true,
	SkillId:  true,
	Priority: true,
}

func fullAgentSkillSelect() *AgentSkillSelect {
	return fullAgentSkillSelectVal
}

func (s *AgentSkillSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.AgentId || s.SkillId || s.Priority
}

func (s *AgentSkillSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type AgentSkillOmit struct {
	AgentId  bool `json:"agentId"`
	SkillId  bool `json:"skillId"`
	Priority bool `json:"priority"`
}

type AgentSkillSelectQuery interface {
	GetRelationParams() (*AgentSkillSelect, *AgentSkillOmit, QueryParams[AgentSkill])
}

func (s *AgentSkillSelect) GetRelationParams() (*AgentSkillSelect, *AgentSkillOmit, QueryParams[AgentSkill]) {
	return s, nil, QueryParams[AgentSkill]{}
}

type AgentSkillQueryBuilder struct {
	selects *AgentSkillSelect
	omits   *AgentSkillOmit
	where   []PredicateOf[AgentSkill]
	take    *int
	skip    *int
	orderBy []OrderBy[AgentSkill]
	cursor  UniquePredicate[AgentSkill]
}

func (b *AgentSkillQueryBuilder) Where(preds ...PredicateOf[AgentSkill]) *AgentSkillQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *AgentSkillQueryBuilder) Take(limit int) *AgentSkillQueryBuilder {
	b.take = &limit
	return b
}

func (b *AgentSkillQueryBuilder) Skip(offset int) *AgentSkillQueryBuilder {
	b.skip = &offset
	return b
}

func (b *AgentSkillQueryBuilder) OrderBy(orders ...OrderBy[AgentSkill]) *AgentSkillQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *AgentSkillQueryBuilder) Cursor(where UniquePredicate[AgentSkill]) *AgentSkillQueryBuilder {
	b.cursor = where
	return b
}

func (b *AgentSkillQueryBuilder) Select(s AgentSkillSelect) *AgentSkillQueryBuilder {
	b.selects = &s
	return b
}

func (b *AgentSkillQueryBuilder) Omit(o AgentSkillOmit) *AgentSkillQueryBuilder {
	b.omits = &o
	return b
}

func (b *AgentSkillQueryBuilder) GetRelationParams() (*AgentSkillSelect, *AgentSkillOmit, QueryParams[AgentSkill]) {
	if b == nil {
		return nil, nil, QueryParams[AgentSkill]{}
	}
	return b.selects, b.omits, QueryParams[AgentSkill]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// AgentSkillCreateArgs is the input argument passed to AgentSkill Create extension hooks.
//
// Fields for AgentSkill:
//
//	agentId  string   required
//	skillId  string   required
//	priority Priority default: NORMAL
//
// Relations for AgentSkill:
//
//	agent    (Agent)
//	skill    (Skill)
type AgentSkillCreateArgs struct {
	// Data contains the model fields to insert.
	Data *AgentSkillCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *AgentSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// AgentSkillCreateManyArgs is the input argument passed to AgentSkill CreateMany extension hooks.
//
// Fields for AgentSkill:
//
//	agentId  string   required
//	skillId  string   required
//	priority Priority default: NORMAL
type AgentSkillCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AgentSkillCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AgentSkillCreateManyArgs) AppendData(builders ...*AgentSkillCreateBuilder) *AgentSkillCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToAgentSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AgentSkillCreateManyAndReturnArgs is the input argument passed to AgentSkill CreateManyAndReturn extension hooks.
//
// Fields for AgentSkill:
//
//	agentId  string   required
//	skillId  string   required
//	priority Priority default: NORMAL
//
// Relations for AgentSkill:
//
//	agent    (Agent)
//	skill    (Skill)
type AgentSkillCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*AgentSkillCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *AgentSkillSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *AgentSkillCreateManyAndReturnArgs) AppendData(builders ...*AgentSkillCreateBuilder) *AgentSkillCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToAgentSkillCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// AgentSkillFindUniqueArgs is the input argument passed to AgentSkill FindUnique extension hooks.
type AgentSkillFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AgentSkill]
	// Select specifies which scalar and relation fields to select and return.
	Select *AgentSkillSelect
}

func (a *AgentSkillFindUniqueArgs) SetWhere(unique UniquePredicate[AgentSkill], additional ...PredicateOf[AgentSkill]) *AgentSkillFindUniqueArgs {
	a.Where = make([]PredicateOf[AgentSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AgentSkillFindFirstArgs is the input argument passed to AgentSkill FindFirst extension hooks.
type AgentSkillFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AgentSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[AgentSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[AgentSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AgentSkillSelect
}

func (a *AgentSkillFindFirstArgs) SetWhere(preds ...PredicateOf[AgentSkill]) *AgentSkillFindFirstArgs {
	a.Where = preds
	return a
}

func (a *AgentSkillFindFirstArgs) SetOrderBy(orders ...OrderBy[AgentSkill]) *AgentSkillFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *AgentSkillFindFirstArgs) SetCursor(cursor UniquePredicate[AgentSkill]) *AgentSkillFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *AgentSkillFindFirstArgs) SetSkip(n int) *AgentSkillFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *AgentSkillFindFirstArgs) SetTake(n int) *AgentSkillFindFirstArgs {
	a.Take = &n
	return a
}

// AgentSkillFindManyArgs is the input argument passed to AgentSkill FindMany extension hooks.
type AgentSkillFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AgentSkill]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[AgentSkill]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[AgentSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *AgentSkillSelect
}

func (a *AgentSkillFindManyArgs) SetWhere(preds ...PredicateOf[AgentSkill]) *AgentSkillFindManyArgs {
	a.Where = preds
	return a
}

func (a *AgentSkillFindManyArgs) SetOrderBy(orders ...OrderBy[AgentSkill]) *AgentSkillFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *AgentSkillFindManyArgs) SetCursor(cursor UniquePredicate[AgentSkill]) *AgentSkillFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *AgentSkillFindManyArgs) SetSkip(n int) *AgentSkillFindManyArgs {
	a.Skip = &n
	return a
}

func (a *AgentSkillFindManyArgs) SetTake(n int) *AgentSkillFindManyArgs {
	a.Take = &n
	return a
}

// AgentSkillCountArgs is the input argument passed to AgentSkill Count extension hooks.
type AgentSkillCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AgentSkill]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *AgentSkillCountArgs) SetWhere(preds ...PredicateOf[AgentSkill]) *AgentSkillCountArgs {
	a.Where = preds
	return a
}

func (a *AgentSkillCountArgs) SetSkip(n int) *AgentSkillCountArgs {
	a.Skip = &n
	return a
}

func (a *AgentSkillCountArgs) SetTake(n int) *AgentSkillCountArgs {
	a.Take = &n
	return a
}

// AgentSkillDeleteArgs is the input argument passed to AgentSkill Delete extension hooks.
type AgentSkillDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AgentSkill]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *AgentSkillSelect
}

func (a *AgentSkillDeleteArgs) SetWhere(unique UniquePredicate[AgentSkill], additional ...PredicateOf[AgentSkill]) *AgentSkillDeleteArgs {
	a.Where = make([]PredicateOf[AgentSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AgentSkillDeleteManyArgs is the input argument passed to AgentSkill DeleteMany extension hooks.
type AgentSkillDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AgentSkill]
}

func (a *AgentSkillDeleteManyArgs) SetWhere(preds ...PredicateOf[AgentSkill]) *AgentSkillDeleteManyArgs {
	a.Where = preds
	return a
}

// AgentSkillUpdateArgs is the input argument passed to AgentSkill Update extension hooks.
type AgentSkillUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[AgentSkill]
	// Data contains the model fields to update.
	Data *AgentSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AgentSkillSelect
}

func (a *AgentSkillUpdateArgs) SetWhere(unique UniquePredicate[AgentSkill], additional ...PredicateOf[AgentSkill]) *AgentSkillUpdateArgs {
	a.Where = make([]PredicateOf[AgentSkill], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// AgentSkillUpdateManyArgs is the input argument passed to AgentSkill UpdateMany extension hooks.
type AgentSkillUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AgentSkill]
	// Data contains the model fields to update.
	Data *AgentSkillUpdate
}

func (a *AgentSkillUpdateManyArgs) SetWhere(preds ...PredicateOf[AgentSkill]) *AgentSkillUpdateManyArgs {
	a.Where = preds
	return a
}

// AgentSkillUpdateManyAndReturnArgs is the input argument passed to AgentSkill UpdateManyAndReturn extension hooks.
type AgentSkillUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[AgentSkill]
	// Data contains the model fields to update.
	Data *AgentSkillUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *AgentSkillSelect
}

func (a *AgentSkillUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[AgentSkill]) *AgentSkillUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type AgentSkillCreateQuery = func(ctx context.Context, args *AgentSkillCreateArgs) (*AgentSkill, error)
type AgentSkillCreateManyQuery = func(ctx context.Context, args *AgentSkillCreateManyArgs) (int64, error)
type AgentSkillCreateManyAndReturnQuery = func(ctx context.Context, args *AgentSkillCreateManyAndReturnArgs) ([]*AgentSkill, error)
type AgentSkillFindUniqueQuery = func(ctx context.Context, args *AgentSkillFindUniqueArgs) (*AgentSkill, error)
type AgentSkillFindFirstQuery = func(ctx context.Context, args *AgentSkillFindFirstArgs) (*AgentSkill, error)
type AgentSkillFindManyQuery = func(ctx context.Context, args *AgentSkillFindManyArgs) ([]*AgentSkill, error)
type AgentSkillDeleteQuery = func(ctx context.Context, args *AgentSkillDeleteArgs) (*AgentSkill, error)
type AgentSkillDeleteManyQuery = func(ctx context.Context, args *AgentSkillDeleteManyArgs) (int64, error)
type AgentSkillCountQuery = func(ctx context.Context, args *AgentSkillCountArgs) (int64, error)
type AgentSkillUpdateQuery = func(ctx context.Context, args *AgentSkillUpdateArgs) (*AgentSkill, error)
type AgentSkillUpdateManyQuery = func(ctx context.Context, args *AgentSkillUpdateManyArgs) (int64, error)
type AgentSkillUpdateManyAndReturnQuery = func(ctx context.Context, args *AgentSkillUpdateManyAndReturnArgs) ([]*AgentSkill, error)

type AgentSkillExtension struct {
	Create              func(ctx context.Context, args *AgentSkillCreateArgs, next AgentSkillCreateQuery) (*AgentSkill, error)
	CreateMany          func(ctx context.Context, args *AgentSkillCreateManyArgs, next AgentSkillCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *AgentSkillCreateManyAndReturnArgs, next AgentSkillCreateManyAndReturnQuery) ([]*AgentSkill, error)
	FindUnique          func(ctx context.Context, args *AgentSkillFindUniqueArgs, next AgentSkillFindUniqueQuery) (*AgentSkill, error)
	FindFirst           func(ctx context.Context, args *AgentSkillFindFirstArgs, next AgentSkillFindFirstQuery) (*AgentSkill, error)
	FindMany            func(ctx context.Context, args *AgentSkillFindManyArgs, next AgentSkillFindManyQuery) ([]*AgentSkill, error)
	Delete              func(ctx context.Context, args *AgentSkillDeleteArgs, next AgentSkillDeleteQuery) (*AgentSkill, error)
	DeleteMany          func(ctx context.Context, args *AgentSkillDeleteManyArgs, next AgentSkillDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *AgentSkillCountArgs, next AgentSkillCountQuery) (int64, error)
	Update              func(ctx context.Context, args *AgentSkillUpdateArgs, next AgentSkillUpdateQuery) (*AgentSkill, error)
	UpdateMany          func(ctx context.Context, args *AgentSkillUpdateManyArgs, next AgentSkillUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *AgentSkillUpdateManyAndReturnArgs, next AgentSkillUpdateManyAndReturnQuery) ([]*AgentSkill, error)
}

type AgentSkillDelegate struct {
	client     *Queries
	extensions []AgentSkillExtension
}

func (d *AgentSkillDelegate) Use(exts ...AgentSkillExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *AgentSkill) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "agentId":
			targets[i] = &m.AgentId
		case "skillId":
			targets[i] = &m.SkillId
		case "priority":
			targets[i] = &m.Priority
		}
	}
	return targets
}

var agentSkillDefaultCols = []string{
	"agentId",
	"skillId",
	"priority",
}

var agentSkillPKCols = []string{
	"agentId",
	"skillId",
}

var agentSkillUniqueCols = []string{}

func selectAgentSkillCols(selects *AgentSkillSelect, omits *AgentSkillOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return agentSkillDefaultCols
	}

	anySelected := selects != nil && (selects.AgentId || selects.SkillId || selects.Priority || selects.Agent != nil || selects.Skill != nil)

	specs := []colSpec{
		{"agentId", selects != nil && selects.AgentId, omits != nil && omits.AgentId, selects != nil && selects.Agent != nil},
		{"skillId", selects != nil && selects.SkillId, omits != nil && omits.SkillId, selects != nil && selects.Skill != nil},
		{"priority", selects != nil && selects.Priority, omits != nil && omits.Priority, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *AgentSkillSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Agent != nil || s.Skill != nil
}

type AgentSkillUpsertBuilder struct {
	*CreateBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]
}

func (b *AgentSkillUpsertBuilder) Select(s AgentSkillSelect) *AgentSkillUpsertBuilder {
	b.selects = &s
	return b
}

func (b *AgentSkillUpsertBuilder) Omit(o AgentSkillOmit) *AgentSkillUpsertBuilder {
	b.omits = &o
	return b
}

type AgentSkillCreateBuilder struct {
	*CreateBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]
}

func (b *AgentSkillCreateBuilder) Select(s AgentSkillSelect) *AgentSkillCreateBuilder {
	b.selects = &s
	return b
}

func (b *AgentSkillCreateBuilder) Omit(o AgentSkillOmit) *AgentSkillCreateBuilder {
	b.omits = &o
	return b
}

func (b *AgentSkillCreateBuilder) OnConflict(target UniqueConstraintTarget) *AgentSkillConflictBuilder[AgentSkillUpsertBuilder] {
	upsertBuilder := &AgentSkillUpsertBuilder{CreateBuilder: b.CreateBuilder}
	return &AgentSkillConflictBuilder[AgentSkillUpsertBuilder]{
		builder:        upsertBuilder,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *AgentSkillCreateBuilder) SetAgentId(v string) *AgentSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}
func (b *AgentSkillCreateBuilder) SetSkillId(v string) *AgentSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}
func (b *AgentSkillCreateBuilder) SetPriority(v PriorityType) *AgentSkillCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AgentSkillCreateBuilder) Assignments(assignments ...FieldAssignmentOf[AgentSkill]) *AgentSkillCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AgentSkillDelegate) Create() *AgentSkillCreateBuilder {
	return &AgentSkillCreateBuilder{
		CreateBuilder: &CreateBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedAgentSkillAgentId  uint64 = 1 << 0
	providedAgentSkillSkillId  uint64 = 1 << 1
	providedAgentSkillPriority uint64 = 1 << 2
)

func assignmentsToAgentSkillCreate(assignments []FieldAssignment) (AgentSkillCreate, error) {
	var input AgentSkillCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "agentId":
			provided |= providedAgentSkillAgentId
			if v, ok := a.Val.(string); ok {
				input.AgentId = v
				errs.ValidateString("agentId", v, true, 0, false, false)
			} else {
				errs.Add("agentId", a.Val, "type", "field agentId must be of type string")
			}
		case "skillId":
			provided |= providedAgentSkillSkillId
			if v, ok := a.Val.(string); ok {
				input.SkillId = v
				errs.ValidateString("skillId", v, true, 0, false, false)
			} else {
				errs.Add("skillId", a.Val, "type", "field skillId must be of type string")
			}
		case "priority":
			provided |= providedAgentSkillPriority
			if v, ok := a.Val.(PriorityType); ok {
				input.Priority = &v
				if !v.IsValid() {
					errs.Add("priority", v, "enum", fmt.Sprintf("invalid enum value %q for field priority", v))
				}
			} else {
				errs.Add("priority", a.Val, "type", "field priority must be of type PriorityType")
			}
		}
	}
	if provided&providedAgentSkillAgentId == 0 {
		errs.Add("agentId", "", "required", "field AgentId is required")
	}
	if provided&providedAgentSkillSkillId == 0 {
		errs.Add("skillId", "", "required", "field SkillId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *AgentSkillCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 3)
	vals = make([]any, 0, 3)
	cols = append(cols, "agentId")
	vals = append(vals, s.AgentId)
	cols = append(cols, "skillId")
	vals = append(vals, s.SkillId)
	if s.Priority != nil {
		cols = append(cols, "priority")
		vals = append(vals, *s.Priority)
	}
	return
}

func partitionAgentSkillInputs(dialect Dialect, inputs []*AgentSkillCreate) [][]*AgentSkillCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*AgentSkillCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*AgentSkillCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*AgentSkillCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*AgentSkillCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*AgentSkillCreate{inputs}
}

func (d *AgentSkillDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *AgentSkillSelect, omits *AgentSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*AgentSkill, error) {
	input, err := assignmentsToAgentSkillCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectAgentSkillCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AgentSkillCreateArgs) (*AgentSkill, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectAgentSkillCols(a.Select, omits)
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
			curr = func(c context.Context, a *AgentSkillCreateArgs) (*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type AgentSkillCreateManyBuilder struct {
	*CreateManyBuilder[AgentSkill]
}

func (b *AgentSkillCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *AgentSkillConflictBuilder[AgentSkillCreateManyBuilder] {
	return &AgentSkillConflictBuilder[AgentSkillCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type AgentSkillCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]
}

func (b *AgentSkillCreateManyAndReturnBuilder) Select(s AgentSkillSelect) *AgentSkillCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *AgentSkillCreateManyAndReturnBuilder) Omit(o AgentSkillOmit) *AgentSkillCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *AgentSkillCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *AgentSkillConflictBuilder[AgentSkillCreateManyAndReturnBuilder] {
	return &AgentSkillConflictBuilder[AgentSkillCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToAgentSkillRecordInputs(builders []*AgentSkillCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *AgentSkillDelegate) CreateMany(builders ...*AgentSkillCreateBuilder) *AgentSkillCreateManyBuilder {
	return &AgentSkillCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[AgentSkill]{
			records:  createBuildersToAgentSkillRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *AgentSkillDelegate) CreateManyAndReturn(builders ...*AgentSkillCreateBuilder) *AgentSkillCreateManyAndReturnBuilder {
	return &AgentSkillCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
			records:  createBuildersToAgentSkillRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToAgentSkillCreateInputs(records []RecordInput) ([]*AgentSkillCreate, error) {
	structs := make([]AgentSkillCreate, len(records))
	inputs := make([]*AgentSkillCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToAgentSkillCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *AgentSkillDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToAgentSkillCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &AgentSkillCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AgentSkillCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AgentSkillCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *AgentSkillSelect, omits *AgentSkillOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*AgentSkill, error) {
	inputs, err := recordsToAgentSkillCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *AgentSkillCreateManyAndReturnArgs) ([]*AgentSkill, error) {
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
			curr = func(c context.Context, a *AgentSkillCreateManyAndReturnArgs) ([]*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *AgentSkillSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*AgentSkill, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *AgentSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AgentSkill.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.AgentSkill.loadRelations(ctx, []*AgentSkill{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "AgentSkill", cols, returningCols, agentSkillPKCols, conflictTarget, conflictAction, len(vals))
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

		var res AgentSkill
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, agentSkillPKCols)
}

func (d *AgentSkillDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*AgentSkill, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("AgentSkill") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "AgentSkill")
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

	var res AgentSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *AgentSkillDelegate) buildBulkInsertSQL(q *Queries, batch []*AgentSkillCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 3)
	for i, c := range agentSkillDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "AgentSkill")
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
			case "agentId":
				vals = append(vals, input.AgentId)
			case "skillId":
				vals = append(vals, input.SkillId)
			case "priority":
				if input.Priority != nil {
					vals = append(vals, *input.Priority)
				} else {
					writeDefault = true
				}
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

func applyAgentSkillConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanAgentSkillRows(rows *sql.Rows, returningCols []string) ([]*AgentSkill, error) {
	var records []*AgentSkill
	for rows.Next() {
		var res AgentSkill
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

func (d *AgentSkillDelegate) runCreateMany(ctx context.Context, inputs []*AgentSkillCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionAgentSkillInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAgentSkillConflictClause(d.client.dialect, queryStr, vals, cols, agentSkillPKCols, conflictTarget, conflictAction)

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

func (d *AgentSkillDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*AgentSkillCreate,
	selects *AgentSkillSelect,
	omits *AgentSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*AgentSkill, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionAgentSkillInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*AgentSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.AgentSkill.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.AgentSkill.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			return nil
		})
		return res, err
	}

	returningCols := selectAgentSkillCols(selects, omits, agentSkillPKCols...)
	recordsOut := make([]*AgentSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAgentSkillConflictClause(d.client.dialect, queryStr, vals, cols, agentSkillPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanAgentSkillRows(rows, returningCols)
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

func (d *AgentSkillDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*AgentSkillCreate,
	selects *AgentSkillSelect,
	omits *AgentSkillOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*AgentSkill, error) {
	batches := partitionAgentSkillInputs(d.client.dialect, inputs)
	returningCols := selectAgentSkillCols(selects, omits, agentSkillPKCols...)
	recordsOut := make([]*AgentSkill, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyAgentSkillConflictClause(d.client.dialect, queryStr, vals, cols, agentSkillPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("AgentSkill") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "AgentSkill")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, agentSkillPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, agentSkillPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanAgentSkillRows(rows, returningCols)
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

type AgentSkillConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *AgentSkillConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *AgentSkillConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *AgentSkillConflictBuilder[B]) Update(fn func(u *AgentSkillUpsert)) *B {
	cb.setAction(*AgentSkillConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// AgentSkillConflictUpdate creates a custom ConflictAction for AgentSkill upsert conflicts.
func AgentSkillConflictUpdate(fn func(u *AgentSkillUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newAgentSkillUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type AgentSkillUpsert struct {
	AgentId  fieldUpsert[string]
	SkillId  fieldUpsert[string]
	Priority fieldUpsert[PriorityType]
}

func newAgentSkillUpsert(up *ConflictUpdate) *AgentSkillUpsert {
	return &AgentSkillUpsert{
		AgentId:  fieldUpsert[string]{column: "agentId", update: up},
		SkillId:  fieldUpsert[string]{column: "skillId", update: up},
		Priority: fieldUpsert[PriorityType]{column: "priority", update: up},
	}
}

type AgentSkillUpdateBuilder struct {
	*UpdateBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]
}

type AgentSkillUpdateManyBuilder struct {
	*UpdateManyBuilder[AgentSkill]
}

type AgentSkillUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]
}

func (b *AgentSkillUpdateBuilder) SetAgentId(v string) *AgentSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *AgentSkillUpdateManyBuilder) SetAgentId(v string) *AgentSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}

func (b *AgentSkillUpdateManyAndReturnBuilder) SetAgentId(v string) *AgentSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "agentId", Val: v})
	return b
}
func (b *AgentSkillUpdateBuilder) SetSkillId(v string) *AgentSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *AgentSkillUpdateManyBuilder) SetSkillId(v string) *AgentSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}

func (b *AgentSkillUpdateManyAndReturnBuilder) SetSkillId(v string) *AgentSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "skillId", Val: v})
	return b
}
func (b *AgentSkillUpdateBuilder) SetPriority(v PriorityType) *AgentSkillUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AgentSkillUpdateManyBuilder) SetPriority(v PriorityType) *AgentSkillUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AgentSkillUpdateManyAndReturnBuilder) SetPriority(v PriorityType) *AgentSkillUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "priority", Val: v})
	return b
}

func (b *AgentSkillUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[AgentSkill]) *AgentSkillUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AgentSkillUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[AgentSkill]) *AgentSkillUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *AgentSkillUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[AgentSkill]) *AgentSkillUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *AgentSkillDelegate) Update(where UniquePredicate[AgentSkill], additional ...PredicateOf[AgentSkill]) *AgentSkillUpdateBuilder {
	return &AgentSkillUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *AgentSkillDelegate) UpdateMany(preds ...PredicateOf[AgentSkill]) *AgentSkillUpdateManyBuilder {
	return &AgentSkillUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[AgentSkill]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *AgentSkillDelegate) UpdateManyAndReturn(preds ...PredicateOf[AgentSkill]) *AgentSkillUpdateManyAndReturnBuilder {
	return &AgentSkillUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *AgentSkillDelegate) buildUpdateSQL(preds []PredicateOf[AgentSkill], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "AgentSkill")
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

func (d *AgentSkillDelegate) executeUpdate(ctx context.Context, where UniquePredicate[AgentSkill], additional []PredicateOf[AgentSkill], assignments []FieldAssignment, selects *AgentSkillSelect, omits *AgentSkillOmit) (*AgentSkill, error) {
	allWhere := make([]PredicateOf[AgentSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToAgentSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentSkillUpdateArgs) (*AgentSkill, error) {
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
			curr = func(c context.Context, a *AgentSkillUpdateArgs) (*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) runUpdate(ctx context.Context, preds []PredicateOf[AgentSkill], cols []string, vals []any, selects *AgentSkillSelect, omits *AgentSkillOmit) (*AgentSkill, error) {
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
		var res *AgentSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.AgentSkill.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.AgentSkill.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAgentSkillCols(selects, omits, agentSkillPKCols...)
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
		return nil, &NotFoundError{Model: "AgentSkill"}
	}

	var res AgentSkill
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*AgentSkill{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *AgentSkillDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[AgentSkill], cols []string, vals []any) (int64, error) {
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

func (d *AgentSkillDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[AgentSkill], cols []string, vals []any, selects *AgentSkillSelect, omits *AgentSkillOmit) (*AgentSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "AgentSkill"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *AgentSkillDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[AgentSkill], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToAgentSkillUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &AgentSkillUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *AgentSkillUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AgentSkillUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *AgentSkillDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[AgentSkill], assignments []FieldAssignment, selects *AgentSkillSelect, omits *AgentSkillOmit) ([]*AgentSkill, error) {
	input, err := assignmentsToAgentSkillUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentSkillUpdateManyAndReturnArgs) ([]*AgentSkill, error) {
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
			curr = func(c context.Context, a *AgentSkillUpdateManyAndReturnArgs) ([]*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[AgentSkill], cols []string, vals []any, selects *AgentSkillSelect, omits *AgentSkillOmit) ([]*AgentSkill, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[AgentSkill]{Where: preds}, selects, omits)
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
		var res []*AgentSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.AgentSkill.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.AgentSkill.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectAgentSkillCols(selects, omits, agentSkillPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanAgentSkillRows(rows, returningCols)
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

func (d *AgentSkillDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[AgentSkill], cols []string, vals []any, selects *AgentSkillSelect, omits *AgentSkillOmit) ([]*AgentSkill, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*AgentSkill{}, nil
	}
	return d.runFindMany(ctx, QueryParams[AgentSkill]{Where: preds}, selects, omits)
}
func (d *AgentSkillDelegate) FindUnique(where UniquePredicate[AgentSkill], additional ...PredicateOf[AgentSkill]) *FindUniqueBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit] {
	return &FindUniqueBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *AgentSkillDelegate) FindFirst(preds ...PredicateOf[AgentSkill]) *FindFirstBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit] {
	return &FindFirstBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *AgentSkillDelegate) FindMany(preds ...PredicateOf[AgentSkill]) *FindManyBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit] {
	return &FindManyBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *AgentSkillDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[AgentSkill], additional []PredicateOf[AgentSkill], selects *AgentSkillSelect, omits *AgentSkillOmit) (*AgentSkill, error) {
	allWhere := make([]PredicateOf[AgentSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentSkillFindUniqueArgs) (*AgentSkill, error) {
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
			curr = func(c context.Context, a *AgentSkillFindUniqueArgs) (*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[AgentSkill],
	selects *AgentSkillSelect,
	omits *AgentSkillOmit,
) (*AgentSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AgentSkillFindFirstArgs) (*AgentSkill, error) {
		return d.runFindFirst(c, QueryParams[AgentSkill]{
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
			curr = func(c context.Context, a *AgentSkillFindFirstArgs) (*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[AgentSkill],
	selects *AgentSkillSelect,
	omits *AgentSkillOmit,
) ([]*AgentSkill, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *AgentSkillFindManyArgs) ([]*AgentSkill, error) {
		return d.runFindMany(c, QueryParams[AgentSkill]{
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
			curr = func(c context.Context, a *AgentSkillFindManyArgs) ([]*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) runFindUnique(ctx context.Context, where []PredicateOf[AgentSkill], selects *AgentSkillSelect, omits *AgentSkillOmit) (*AgentSkill, error) {
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
	returningCols := selectAgentSkillCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*AgentSkill{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *AgentSkillDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[AgentSkill],
	selects *AgentSkillSelect,
	omits *AgentSkillOmit,
) (*AgentSkill, error) {
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

func (d *AgentSkillDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[AgentSkill],
	selects *AgentSkillSelect,
	omits *AgentSkillOmit,
) ([]*AgentSkill, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, agentSkillPKCols, agentSkillUniqueCols, "AgentSkill", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, agentSkillPKCols, agentSkillUniqueCols, isCursorQuery, params.Take)
	returningCols := selectAgentSkillCols(selects, omits)

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

func (d *AgentSkillDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*AgentSkill, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "AgentSkill", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res AgentSkill
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *AgentSkillDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*AgentSkill, error) {
	query := buildSelectSQL(d.client, "AgentSkill", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*AgentSkill, 0)
	for rows.Next() {
		var res AgentSkill
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
func (d *AgentSkillDelegate) DeleteMany(preds ...PredicateOf[AgentSkill]) *DeleteManyBuilder[AgentSkill] {
	return &DeleteManyBuilder[AgentSkill]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *AgentSkillDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[AgentSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &AgentSkillDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *AgentSkillDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *AgentSkillDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[AgentSkill]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "AgentSkill")
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

func (d *AgentSkillDelegate) Delete(where UniquePredicate[AgentSkill], additional ...PredicateOf[AgentSkill]) *DeleteBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit] {
	return &DeleteBuilder[AgentSkill, AgentSkillSelect, AgentSkillOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *AgentSkillDelegate) executeDelete(ctx context.Context, where UniquePredicate[AgentSkill], additional []PredicateOf[AgentSkill], selects *AgentSkillSelect, omits *AgentSkillOmit) (*AgentSkill, error) {
	allWhere := make([]PredicateOf[AgentSkill], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullAgentSkillSelect()
	}

	args := &AgentSkillDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *AgentSkillDeleteArgs) (*AgentSkill, error) {
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
			curr = func(c context.Context, a *AgentSkillDeleteArgs) (*AgentSkill, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) runDelete(ctx context.Context, where []PredicateOf[AgentSkill], selects *AgentSkillSelect, omits *AgentSkillOmit) (*AgentSkill, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectAgentSkillCols(selects, omits, agentSkillPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *AgentSkill
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AgentSkill.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "AgentSkill"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "AgentSkill")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[AgentSkill]
			pkPreds = append(pkPreds, Predicate[AgentSkill]{
				Data: PredicateData{
					Column:   "agentId",
					Operator: "=",
					Value:    res.AgentId,
				},
			})
			pkPreds = append(pkPreds, Predicate[AgentSkill]{
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
	d.client.dialect.WriteQuotedIdent(&sb, "AgentSkill")

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
		return nil, &NotFoundError{Model: "AgentSkill"}
	}

	var row AgentSkill
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *AgentSkillDelegate) Count(preds ...PredicateOf[AgentSkill]) *CountBuilder[AgentSkill] {
	return &CountBuilder[AgentSkill]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *AgentSkillDelegate) executeCount(ctx context.Context, params QueryParams[AgentSkill]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &AgentSkillCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *AgentSkillCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[AgentSkill]{
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
			curr = func(c context.Context, a *AgentSkillCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *AgentSkillDelegate) runCount(ctx context.Context, params QueryParams[AgentSkill]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "AgentSkill")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "AgentSkill")
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
func (d *AgentSkillDelegate) loadRelations(ctx context.Context, records []*AgentSkill, selects *AgentSkillSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Agent != nil {
		relationSelects, relationOmits, relationParams := selects.Agent.GetRelationParams()
		returningCols := selectAgentCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: AgentSkill.agentId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *AgentSkill) string { return p.AgentId }),
			"Agent",
			"id",
			returningCols,
			scanInto(returningCols, (*Agent).ScanFields),
			directKey(func(c *Agent) string { return c.Id }),
			setOne(func(p *AgentSkill, c *Agent) { p.Agent = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading agent: %w", err)
		}
		if err := d.client.Agent.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Skill != nil {
		relationSelects, relationOmits, relationParams := selects.Skill.GetRelationParams()
		returningCols := selectSkillCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: AgentSkill.skillId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *AgentSkill) string { return p.SkillId }),
			"Skill",
			"id",
			returningCols,
			scanInto(returningCols, (*Skill).ScanFields),
			directKey(func(c *Skill) string { return c.Id }),
			setOne(func(p *AgentSkill, c *Skill) { p.Skill = c }),
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
