package riot

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var _ = time.Time{}
var _ = json.RawMessage{}
var _ = strings.Join
var _ = slices.Clone[[]any]
var _ = rand.Read
var _ = strconv.AppendUint

//go:embed migrations/*.sql
var migrationsFS embed.FS

func generateCUID() string {
	now := uint64(time.Now().UnixMilli())
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	buf := make([]byte, 1, 33)
	buf[0] = 'c'
	buf = strconv.AppendUint(buf, now, 16)

	const hextable = "0123456789abcdef"
	for _, v := range b {
		buf = append(buf, hextable[v>>4], hextable[v&0x0f])
	}
	return string(buf)
}

// URIType represents the "URI" enum in the schema.
//
// Values:
//
//	URI_STATIC -> "static"
//	URI_URL    -> "url"
type URIType string

const (
	// URI_STATIC maps to "static" in the database.
	URI_STATIC URIType = "static"
	// URI_URL maps to "url" in the database.
	URI_URL URIType = "url"
)

// IsValid reports whether e is a valid URIType enum value.
func (e URIType) IsValid() bool {
	switch e {
	case URI_STATIC, URI_URL:
		return true
	}
	return false
}

// ReferenceTypeType represents the "ReferenceType" enum in the schema.
//
// Values:
//
//	ReferenceType_LITERATURE    -> "literature"
//	ReferenceType_IMAGE         -> "image"
//	ReferenceType_VIDEO         -> "video"
//	ReferenceType_COMMANDS      -> "commands"
//	ReferenceType_SCRIPT        -> "script"
//	ReferenceType_CONFIGURATION -> "configuration"
type ReferenceTypeType string

const (
	// ReferenceType_LITERATURE maps to "literature" in the database.
	ReferenceType_LITERATURE ReferenceTypeType = "literature"
	// ReferenceType_IMAGE maps to "image" in the database.
	ReferenceType_IMAGE ReferenceTypeType = "image"
	// ReferenceType_VIDEO maps to "video" in the database.
	ReferenceType_VIDEO ReferenceTypeType = "video"
	// ReferenceType_COMMANDS maps to "commands" in the database.
	ReferenceType_COMMANDS ReferenceTypeType = "commands"
	// ReferenceType_SCRIPT maps to "script" in the database.
	ReferenceType_SCRIPT ReferenceTypeType = "script"
	// ReferenceType_CONFIGURATION maps to "configuration" in the database.
	ReferenceType_CONFIGURATION ReferenceTypeType = "configuration"
)

// IsValid reports whether e is a valid ReferenceTypeType enum value.
func (e ReferenceTypeType) IsValid() bool {
	switch e {
	case ReferenceType_LITERATURE, ReferenceType_IMAGE, ReferenceType_VIDEO, ReferenceType_COMMANDS, ReferenceType_SCRIPT, ReferenceType_CONFIGURATION:
		return true
	}
	return false
}

// ThemeModeType represents the "ThemeMode" enum in the schema.
//
// Values:
//
//	ThemeMode_LIGHT     -> "light"
//	ThemeMode_DARK      -> "dark"
//	ThemeMode_PRIMITIVE -> "primitive"
type ThemeModeType string

const (
	// ThemeMode_LIGHT maps to "light" in the database.
	ThemeMode_LIGHT ThemeModeType = "light"
	// ThemeMode_DARK maps to "dark" in the database.
	ThemeMode_DARK ThemeModeType = "dark"
	// ThemeMode_PRIMITIVE maps to "primitive" in the database.
	ThemeMode_PRIMITIVE ThemeModeType = "primitive"
)

// IsValid reports whether e is a valid ThemeModeType enum value.
func (e ThemeModeType) IsValid() bool {
	switch e {
	case ThemeMode_LIGHT, ThemeMode_DARK, ThemeMode_PRIMITIVE:
		return true
	}
	return false
}

// StyleTypeType represents the "StyleType" enum in the schema.
//
// Values:
//
//	StyleType_COLOR     -> "color"
//	StyleType_SPATIAL   -> "spatial"
//	StyleType_SIZE      -> "size"
//	StyleType_DIMENSION -> "dimension"
//	StyleType_FONT      -> "font"
//	StyleType_SHADOW    -> "shadow"
//	StyleType_ZINDEX    -> "zindex"
//	StyleType_GRADIENT  -> "gradient"
//	StyleType_WEIGHT    -> "weight"
//	StyleType_CSS       -> "css"
type StyleTypeType string

const (
	// StyleType_COLOR maps to "color" in the database.
	StyleType_COLOR StyleTypeType = "color"
	// StyleType_SPATIAL maps to "spatial" in the database.
	StyleType_SPATIAL StyleTypeType = "spatial"
	// StyleType_SIZE maps to "size" in the database.
	StyleType_SIZE StyleTypeType = "size"
	// StyleType_DIMENSION maps to "dimension" in the database.
	StyleType_DIMENSION StyleTypeType = "dimension"
	// StyleType_FONT maps to "font" in the database.
	StyleType_FONT StyleTypeType = "font"
	// StyleType_SHADOW maps to "shadow" in the database.
	StyleType_SHADOW StyleTypeType = "shadow"
	// StyleType_ZINDEX maps to "zindex" in the database.
	StyleType_ZINDEX StyleTypeType = "zindex"
	// StyleType_GRADIENT maps to "gradient" in the database.
	StyleType_GRADIENT StyleTypeType = "gradient"
	// StyleType_WEIGHT maps to "weight" in the database.
	StyleType_WEIGHT StyleTypeType = "weight"
	// StyleType_CSS maps to "css" in the database.
	StyleType_CSS StyleTypeType = "css"
)

// IsValid reports whether e is a valid StyleTypeType enum value.
func (e StyleTypeType) IsValid() bool {
	switch e {
	case StyleType_COLOR, StyleType_SPATIAL, StyleType_SIZE, StyleType_DIMENSION, StyleType_FONT, StyleType_SHADOW, StyleType_ZINDEX, StyleType_GRADIENT, StyleType_WEIGHT, StyleType_CSS:
		return true
	}
	return false
}

// TokenTypeType represents the "TokenType" enum in the schema.
//
// Values:
//
//	TokenType_SEMANTIC  -> "semantic"
//	TokenType_APPLIED   -> "applied"
//	TokenType_PRIMITIVE -> "primitive"
type TokenTypeType string

const (
	// TokenType_SEMANTIC maps to "semantic" in the database.
	TokenType_SEMANTIC TokenTypeType = "semantic"
	// TokenType_APPLIED maps to "applied" in the database.
	TokenType_APPLIED TokenTypeType = "applied"
	// TokenType_PRIMITIVE maps to "primitive" in the database.
	TokenType_PRIMITIVE TokenTypeType = "primitive"
)

// IsValid reports whether e is a valid TokenTypeType enum value.
func (e TokenTypeType) IsValid() bool {
	switch e {
	case TokenType_SEMANTIC, TokenType_APPLIED, TokenType_PRIMITIVE:
		return true
	}
	return false
}

// PriorityType represents the "Priority" enum in the schema.
//
// Values:
//
//	Priority_LOW       -> "LOW"
//	Priority_NORMAL    -> "NORMAL"
//	Priority_HIGH      -> "HIGH"
//	Priority_ESSENTIAL -> "ESSENTIAL"
type PriorityType string

const (
	// Priority_LOW maps to "LOW" in the database.
	Priority_LOW PriorityType = "LOW"
	// Priority_NORMAL maps to "NORMAL" in the database.
	Priority_NORMAL PriorityType = "NORMAL"
	// Priority_HIGH maps to "HIGH" in the database.
	Priority_HIGH PriorityType = "HIGH"
	// Priority_ESSENTIAL maps to "ESSENTIAL" in the database.
	Priority_ESSENTIAL PriorityType = "ESSENTIAL"
)

// IsValid reports whether e is a valid PriorityType enum value.
func (e PriorityType) IsValid() bool {
	switch e {
	case Priority_LOW, Priority_NORMAL, Priority_HIGH, Priority_ESSENTIAL:
		return true
	}
	return false
}

// Standard Sentinel Errors
var (
	// ErrNotFound is returned when a query expects a record but finds none.
	ErrNotFound = errors.New("phi: record not found")

	// ErrNoRowsAffected is returned when an update or delete operation affects zero rows when a target was expected.
	ErrNoRowsAffected = errors.New("phi: no rows affected")

	// ErrConstraint is the base error for database constraint violations.
	ErrConstraint = errors.New("phi: constraint violation")

	// ErrUniqueConstraint is returned when a unique index or primary key constraint is violated.
	ErrUniqueConstraint = errors.New("phi: unique constraint violation")

	// ErrFKConstraint is returned when a foreign key constraint is violated.
	ErrFKConstraint = errors.New("phi: foreign key constraint violation")

	// ErrNotNullConstraint is returned when a NOT NULL column constraint is violated.
	ErrNotNullConstraint = errors.New("phi: not null constraint violation")

	// ErrCheckConstraint is returned when a CHECK constraint is violated.
	ErrCheckConstraint = errors.New("phi: check constraint violation")

	// ErrDeadlock is returned when a database deadlock is detected.
	ErrDeadlock = errors.New("phi: deadlock detected")

	// ErrLockTimeout is returned when a database lock wait timeout is exceeded.
	ErrLockTimeout = errors.New("phi: lock wait timeout exceeded")

	// ErrSerialization is returned when a transaction fails due to serialization / concurrent update conflicts.
	ErrSerialization = errors.New("phi: serialization failure")

	// ErrTxDone is returned when attempting an operation on a transaction that has already been committed or rolled back.
	ErrTxDone = errors.New("phi: transaction already committed or rolled back")

	// ErrConnClosed is returned when the underlying database connection is closed.
	ErrConnClosed = errors.New("phi: connection closed")
)

// Helper functions for inspecting error kinds (Ent-style)

// IsNotFound returns true if the error indicates a record was not found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsNoRowsAffected returns true if no rows were modified by an update or delete execution.
func IsNoRowsAffected(err error) bool {
	return errors.Is(err, ErrNoRowsAffected)
}

// IsConstraintError returns true if the error is any type of database constraint violation.
func IsConstraintError(err error) bool {
	return errors.Is(err, ErrConstraint) ||
		errors.Is(err, ErrUniqueConstraint) ||
		errors.Is(err, ErrFKConstraint) ||
		errors.Is(err, ErrNotNullConstraint) ||
		errors.Is(err, ErrCheckConstraint)
}

// IsUniqueConstraint returns true if the error is a unique constraint violation.
func IsUniqueConstraint(err error) bool {
	return errors.Is(err, ErrUniqueConstraint)
}

// IsFKConstraint returns true if the error is a foreign key constraint violation.
func IsFKConstraint(err error) bool {
	return errors.Is(err, ErrFKConstraint)
}

// IsNotNullConstraint returns true if the error is a NOT NULL constraint violation.
func IsNotNullConstraint(err error) bool {
	return errors.Is(err, ErrNotNullConstraint)
}

// IsCheckConstraint returns true if the error is a CHECK constraint violation.
func IsCheckConstraint(err error) bool {
	return errors.Is(err, ErrCheckConstraint)
}

// IsDeadlock returns true if the error is a database deadlock failure.
func IsDeadlock(err error) bool {
	return errors.Is(err, ErrDeadlock)
}

// IsLockTimeout returns true if the error is a lock timeout failure.
func IsLockTimeout(err error) bool {
	return errors.Is(err, ErrLockTimeout)
}

// IsSerialization returns true if the error is a transaction serialization failure.
func IsSerialization(err error) bool {
	return errors.Is(err, ErrSerialization)
}

// IsTxDone returns true if the error indicates the transaction was already completed.
func IsTxDone(err error) bool {
	return errors.Is(err, ErrTxDone)
}

// IsValidationError returns true if the error is a client-side validation failure.
func IsValidationError(err error) bool {
	var ve ValidationError
	if errors.As(err, &ve) {
		return true
	}
	var vePtr *ValidationError
	return errors.As(err, &vePtr)
}

// Concrete Error Structs

// NotFoundError represents a record not found error with model context.
type NotFoundError struct {
	Model string
	Cause error
}

func (e *NotFoundError) Error() string {
	if e.Model != "" {
		return fmt.Sprintf("phi: %s record not found", e.Model)
	}
	return "phi: record not found"
}

func (e *NotFoundError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return ErrNotFound
}

// ConstraintError represents a detailed database constraint error.
type ConstraintError struct {
	Kind       error  // e.g. ErrUniqueConstraint, ErrFKConstraint
	Table      string // Table name where constraint occurred (if available)
	Constraint string // Constraint name (if available)
	Cause      error  // Underlying driver error
}

func (e *ConstraintError) Error() string {
	if e.Constraint != "" {
		return fmt.Sprintf("%v: constraint %q violated on table %q: %v", e.Kind, e.Constraint, e.Table, e.Cause)
	}
	if e.Table != "" {
		return fmt.Sprintf("%v: constraint violated on table %q: %v", e.Kind, e.Table, e.Cause)
	}
	return fmt.Sprintf("%v: %v", e.Kind, e.Cause)
}

func (e *ConstraintError) Unwrap() error {
	if e.Kind != nil {
		return e.Kind
	}
	return ErrConstraint
}

// TranslateDBError normalizes driver-specific errors and standard library sql errors into Phi domain errors.
func TranslateDBError(err error) error {
	if err == nil {
		return nil
	}

	// Idempotency check: If err is ALREADY a normalized Phi domain error, return as-is to avoid duplicate wrapping!
	if errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrConstraint) ||
		errors.Is(err, ErrTxDone) ||
		errors.Is(err, ErrConnClosed) ||
		errors.Is(err, ErrDeadlock) ||
		errors.Is(err, ErrLockTimeout) ||
		errors.Is(err, ErrSerialization) {
		return err
	}

	// STDLIB sql error mappings
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}
	if errors.Is(err, sql.ErrTxDone) {
		return fmt.Errorf("%w: %w", ErrTxDone, err)
	}
	if errors.Is(err, sql.ErrConnDone) {
		return fmt.Errorf("%w: %w", ErrConnClosed, err)
	}
	// SQLite Error Translation
	type codeCoder interface {
		Code() int
	}
	if cErr, ok := err.(codeCoder); ok {
		code := cErr.Code()
		switch code {
		case 2067, 1555:
			return &ConstraintError{Kind: ErrUniqueConstraint, Cause: err}
		case 787:
			return &ConstraintError{Kind: ErrFKConstraint, Cause: err}
		case 1299:
			return &ConstraintError{Kind: ErrNotNullConstraint, Cause: err}
		case 275:
			return &ConstraintError{Kind: ErrCheckConstraint, Cause: err}
		}
	}
	type sqliteExtendedCoder interface {
		ExtendedCode() int
	}
	if sqExt, ok := err.(sqliteExtendedCoder); ok {
		code := sqExt.ExtendedCode()
		switch code {
		case 2067, 1555:
			return &ConstraintError{Kind: ErrUniqueConstraint, Cause: err}
		case 787:
			return &ConstraintError{Kind: ErrFKConstraint, Cause: err}
		case 1299:
			return &ConstraintError{Kind: ErrNotNullConstraint, Cause: err}
		case 275:
			return &ConstraintError{Kind: ErrCheckConstraint, Cause: err}
		}
	}
	msg := err.Error()
	if strings.Contains(msg, "UNIQUE constraint failed") || strings.Contains(msg, "PRIMARY KEY constraint failed") {
		return &ConstraintError{Kind: ErrUniqueConstraint, Cause: err}
	}
	if strings.Contains(msg, "FOREIGN KEY constraint failed") {
		return &ConstraintError{Kind: ErrFKConstraint, Cause: err}
	}
	if strings.Contains(msg, "NOT NULL constraint failed") {
		return &ConstraintError{Kind: ErrNotNullConstraint, Cause: err}
	}
	if strings.Contains(msg, "CHECK constraint failed") {
		return &ConstraintError{Kind: ErrCheckConstraint, Cause: err}
	}

	return err
}

// FieldError represents a single validation failure on a specific field.
type FieldError struct {
	Field string
	Value any
	Rule  string
	Msg   string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field %s: %s (value: %v, rule: %s)", e.Field, e.Msg, e.Value, e.Rule)
}

// ValidationError collects multiple validation errors during an operation.
type ValidationError struct {
	Errors []FieldError
}

func (e ValidationError) Error() string {
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, "; "))
}

func (e *ValidationError) Add(field string, value any, rule string, msg string) {
	e.Errors = append(e.Errors, FieldError{
		Field: field,
		Value: value,
		Rule:  rule,
		Msg:   msg,
	})
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

func (errs *ValidationError) ValidateString(fieldName string, val string, isRequired bool, maxLen int, isBit bool, isInet bool) {
	if isRequired && val == "" {
		errs.Add(fieldName, val, "required", fmt.Sprintf("field %s is required", fieldName))
	}
	if strings.Contains(val, "\x00") {
		errs.Add(fieldName, val, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(val) {
		errs.Add(fieldName, val, "safety", "string must be valid UTF-8")
	}
	if maxLen > 0 && utf8.RuneCountInString(val) > maxLen {
		errs.Add(fieldName, val, "length", fmt.Sprintf("string exceeds maximum length of %d characters", maxLen))
	}
	if isBit {
		if strings.IndexFunc(val, func(r rune) bool { return r != '0' && r != '1' }) >= 0 {
			errs.Add(fieldName, val, "format", "bit string must contain only '0' and '1'")
		}
	}
}

func (errs *ValidationError) ValidateInt32(fieldName string, val int32, rule string) {
	switch rule {
	case "SmallInt":
		if val < -32768 || val > 32767 {
			errs.Add(fieldName, val, "range", "value is out of range for SmallInt (-32768 to 32767)")
		}
	case "TinyInt":
		if val < -128 || val > 127 {
			errs.Add(fieldName, val, "range", "value is out of range for TinyInt (-128 to 127)")
		}
	case "Oid":
		if val < 0 {
			errs.Add(fieldName, val, "range", "value is out of range for Oid (must be non-negative)")
		}
	}
}

func (errs *ValidationError) ValidateInt64(fieldName string, val int64, rule string) {
	switch rule {
	case "SmallInt":
		if val < -32768 || val > 32767 {
			errs.Add(fieldName, val, "range", "value is out of range for SmallInt (-32768 to 32767)")
		}
	case "TinyInt":
		if val < -128 || val > 127 {
			errs.Add(fieldName, val, "range", "value is out of range for TinyInt (-128 to 127)")
		}
	case "Oid":
		if val < 0 {
			errs.Add(fieldName, val, "range", "value is out of range for Oid (must be non-negative)")
		}
	}
}

func (errs *ValidationError) ValidateInt(fieldName string, val int, rule string) {
	switch rule {
	case "SmallInt":
		if val < -32768 || val > 32767 {
			errs.Add(fieldName, val, "range", "value is out of range for SmallInt (-32768 to 32767)")
		}
	case "TinyInt":
		if val < -128 || val > 127 {
			errs.Add(fieldName, val, "range", "value is out of range for TinyInt (-128 to 127)")
		}
	case "Oid":
		if val < 0 {
			errs.Add(fieldName, val, "range", "value is out of range for Oid (must be non-negative)")
		}
	}
}

type UniqueConstraintTarget interface {
	UniqueColumns() []string
}

type CompositeUniqueConstraint[M any] struct {
	Name    string
	Columns []string
}

func (c CompositeUniqueConstraint[M]) UniqueColumns() []string {
	return c.Columns
}

type ConflictActionType int

const (
	ConflictActionIgnore ConflictActionType = iota
	ConflictActionUpdateNewValues
	ConflictActionUpdateCustom
)

type ConflictAction struct {
	Type        ConflictActionType
	Assignments []string
	Args        []any
}

func ConflictDoNothing() *ConflictAction {
	return &ConflictAction{Type: ConflictActionIgnore}
}

func ConflictUpdateNewValues() *ConflictAction {
	return &ConflictAction{Type: ConflictActionUpdateNewValues}
}

func (a *ConflictAction) IsDoNothing() bool {
	return a != nil && a.Type == ConflictActionIgnore
}

func (a *ConflictAction) IsUpdateNewValues() bool {
	return a != nil && a.Type == ConflictActionUpdateNewValues
}

func (a *ConflictAction) IsUpdateCustom() bool {
	return a != nil && a.Type == ConflictActionUpdateCustom
}

func (a *ConflictAction) String() string {
	if a == nil {
		return "ConflictAction(NONE)"
	}
	switch a.Type {
	case ConflictActionIgnore:
		return "ConflictAction(DO NOTHING)"
	case ConflictActionUpdateNewValues:
		return "ConflictAction(UPDATE NEW VALUES)"
	case ConflictActionUpdateCustom:
		return fmt.Sprintf("ConflictAction(UPDATE CUSTOM: %v)", a.Assignments)
	default:
		return "ConflictAction(UNKNOWN)"
	}
}

type ConflictUpdate struct {
	assignments []string
	args        []any
}

type fieldUpsert[T any] struct {
	column string
	update *ConflictUpdate
}

func (f fieldUpsert[T]) Update() {
	f.update.assignments = append(f.update.assignments, fmt.Sprintf(`"%s" = EXCLUDED."%s"`, f.column, f.column))
}

func (f fieldUpsert[T]) Set(val T) {
	f.update.assignments = append(f.update.assignments, fmt.Sprintf(`"%s" = ?`, f.column))
	f.update.args = append(f.update.args, val)
}

type numericFieldUpsert[T any] struct {
	fieldUpsert[T]
	tableName string
}

func (f numericFieldUpsert[T]) Increment(val T) {
	f.update.assignments = append(f.update.assignments, fmt.Sprintf(`"%s" = "%s"."%s" + ?`, f.column, f.tableName, f.column))
	f.update.args = append(f.update.args, val)
}

func (f numericFieldUpsert[T]) Decrement(val T) {
	f.update.assignments = append(f.update.assignments, fmt.Sprintf(`"%s" = "%s"."%s" - ?`, f.column, f.tableName, f.column))
	f.update.args = append(f.update.args, val)
}

type Dialect struct {
	QuoteChar               byte
	PlaceholderFmt          string // "" means "?"
	SupportsInsertReturning bool
	SupportsUpdateReturning bool
	SupportsDeleteReturning bool
	SupportsLimitMinusOne   bool
	SupportsBulkInsert      bool
	SupportsDefaultKeyword  bool
	RequiresConflictTarget  bool
	ConflictKeyword         string // "ON CONFLICT" (Pg/SQLite) or "ON DUPLICATE KEY" (MySQL)
	ConflictIgnore          string // "DO NOTHING" or ""
	ConflictUpdate          string // "DO UPDATE SET" or "UPDATE"
	ConflictExcluded        string // "EXCLUDED." or "VALUES("
	ConflictExcludedEnd     string // "" or ")"
	ILikeFmt                string // "{col} ILIKE {val}" or "LOWER({col}) LIKE LOWER({val})"
	ArrayHasFmt             string // "{val} = ANY({col})" or "JSON_CONTAINS({col}, {val})"
	ArrayHasEveryFmt        string // "{col} @> ARRAY[{vals}]" (empty = compose from ArrayHasFmt)
	ArrayHasSomeFmt         string // "{col} && ARRAY[{vals}]" (empty = compose from ArrayHasFmt)
	ArrayIsEmptyFmt         string // "(cardinality({col}) = 0 OR {col} IS NULL)" or "({col} IS NULL OR json_array_length({col}) = 0)"
}

func (d Dialect) WriteQuotedIdent(sb *strings.Builder, ident string) {
	sb.WriteByte(d.QuoteChar)
	sb.WriteString(ident)
	sb.WriteByte(d.QuoteChar)
}

func (d Dialect) WritePlaceholder(sb *strings.Builder, idx int) {
	if d.PlaceholderFmt == "" {
		sb.WriteByte('?')
	} else {
		sb.WriteString(d.PlaceholderFmt)
		sb.WriteString(strconv.Itoa(idx))
	}
}

func (d Dialect) Quote(ident string) string {
	var sb strings.Builder
	sb.Grow(len(ident) + 2)
	d.WriteQuotedIdent(&sb, ident)
	return sb.String()
}

func (d Dialect) BindVar(idx int) string {
	if d.PlaceholderFmt == "" {
		return "?"
	}
	var sb strings.Builder
	sb.Grow(8)
	sb.WriteString(d.PlaceholderFmt)
	sb.WriteString(strconv.Itoa(idx))
	return sb.String()
}

func (d Dialect) FormatLimitOffset(take *int, skip *int) string {
	if take != nil {
		if skip != nil {
			return fmt.Sprintf(" LIMIT %d OFFSET %d", *take, *skip)
		}
		return fmt.Sprintf(" LIMIT %d", *take)
	}
	if skip != nil {
		if d.SupportsLimitMinusOne {
			return fmt.Sprintf(" LIMIT -1 OFFSET %d", *skip)
		}
		return fmt.Sprintf(" OFFSET %d", *skip)
	}
	return ""
}

func (d Dialect) FormatILike(col string, placeholder string) string {
	r := strings.NewReplacer("{col}", d.Quote(col), "{val}", placeholder)
	return r.Replace(d.ILikeFmt)
}

func (d Dialect) FormatArrayHas(col string, placeholder string) string {
	r := strings.NewReplacer("{col}", d.Quote(col), "{val}", placeholder)
	return r.Replace(d.ArrayHasFmt)
}

func (d Dialect) FormatArrayHasEvery(col string, placeholders []string) string {
	if d.ArrayHasEveryFmt != "" {
		r := strings.NewReplacer("{col}", d.Quote(col), "{vals}", strings.Join(placeholders, ", "))
		return r.Replace(d.ArrayHasEveryFmt)
	}
	var parts []string
	for _, p := range placeholders {
		parts = append(parts, d.FormatArrayHas(col, p))
	}
	if len(parts) == 0 {
		return "1=1"
	}
	return "(" + strings.Join(parts, " AND ") + ")"
}

func (d Dialect) FormatArrayHasSome(col string, placeholders []string) string {
	if d.ArrayHasSomeFmt != "" {
		r := strings.NewReplacer("{col}", d.Quote(col), "{vals}", strings.Join(placeholders, ", "))
		return r.Replace(d.ArrayHasSomeFmt)
	}
	var parts []string
	for _, p := range placeholders {
		parts = append(parts, d.FormatArrayHas(col, p))
	}
	if len(parts) == 0 {
		return "1=0"
	}
	return "(" + strings.Join(parts, " OR ") + ")"
}

func (d Dialect) FormatArrayIsEmpty(col string) string {
	r := strings.NewReplacer("{col}", d.Quote(col))
	return r.Replace(d.ArrayIsEmptyFmt)
}

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type FieldAssignment struct {
	Col string
	Val any
}

type FieldAssignmentOf[M any] struct {
	Col string
	Val any
}

type RecordInput struct {
	Assignments []FieldAssignment
}

type OrderDirection string

const (
	Asc  OrderDirection = "ASC"
	Desc OrderDirection = "DESC"
)

type OrderBy[M any] struct {
	Field     string
	Direction OrderDirection
}

type QueryParams[M any] struct {
	Where   []PredicateOf[M]
	Take    *int
	Skip    *int
	OrderBy []OrderBy[M]
	Cursor  UniquePredicate[M]
}

func computeNonConflictCols(allCols []string, conflictCols []string, pkCols []string) []string {
	exclude := make(map[string]bool, len(conflictCols)+len(pkCols))
	for _, c := range conflictCols {
		exclude[c] = true
	}
	for _, c := range pkCols {
		exclude[c] = true
	}
	var result []string
	for _, c := range allCols {
		if !exclude[c] {
			result = append(result, c)
		}
	}
	return result
}

func (d Dialect) BuildConflictClause(conflictCols []string, action *ConflictAction, nonConflictCols []string, startParamIndex int) (string, []any) {
	if action == nil {
		return "", nil
	}
	var colsStr string
	if len(conflictCols) > 0 && d.RequiresConflictTarget {
		var quoted []string
		for _, col := range conflictCols {
			quoted = append(quoted, d.Quote(col))
		}
		colsStr = " (" + strings.Join(quoted, ", ") + ")"
	}
	switch action.Type {
	case ConflictActionIgnore:
		if d.ConflictIgnore == "" {
			return "", nil
		}
		return " " + d.ConflictKeyword + colsStr + " " + d.ConflictIgnore, nil
	case ConflictActionUpdateNewValues:
		if len(nonConflictCols) == 0 {
			if d.ConflictIgnore == "" {
				return "", nil
			}
			return " " + d.ConflictKeyword + colsStr + " " + d.ConflictIgnore, nil
		}
		var sets []string
		for _, col := range nonConflictCols {
			qc := d.Quote(col)
			sets = append(sets, fmt.Sprintf(`%s = %s%s%s`, qc, d.ConflictExcluded, qc, d.ConflictExcludedEnd))
		}
		return " " + d.ConflictKeyword + colsStr + " " + d.ConflictUpdate + " " + strings.Join(sets, ", "), nil
	case ConflictActionUpdateCustom:
		if len(action.Assignments) == 0 {
			if d.ConflictIgnore == "" {
				return "", nil
			}
			return " " + d.ConflictKeyword + colsStr + " " + d.ConflictIgnore, nil
		}
		var assignments []string
		paramIdx := startParamIndex
		for _, assignment := range action.Assignments {
			var sb strings.Builder
			src := assignment
			for {
				idx := strings.Index(src, "?")
				if idx == -1 {
					sb.WriteString(src)
					break
				}
				sb.WriteString(src[:idx])
				sb.WriteString(d.BindVar(paramIdx))
				paramIdx++
				src = src[idx+1:]
			}
			assignments = append(assignments, sb.String())
		}
		return " " + d.ConflictKeyword + colsStr + " " + d.ConflictUpdate + " " + strings.Join(assignments, ", "), action.Args
	}
	return "", nil
}

func newDialect() Dialect {
	return Dialect{
		QuoteChar:               '"',
		PlaceholderFmt:          "", // means ?
		SupportsInsertReturning: true,
		SupportsUpdateReturning: true,
		SupportsDeleteReturning: true,
		SupportsLimitMinusOne:   true,
		SupportsBulkInsert:      true,
		SupportsDefaultKeyword:  false,
		RequiresConflictTarget:  true,
		ConflictKeyword:         "ON CONFLICT",
		ConflictIgnore:          "DO NOTHING",
		ConflictUpdate:          "DO UPDATE SET",
		ConflictExcluded:        "EXCLUDED.",
		ConflictExcludedEnd:     "",
		ILikeFmt:                "LOWER({col}) LIKE LOWER({val})",
		ArrayHasFmt:             "{val} IN (SELECT value FROM json_each({col}))",
		ArrayHasEveryFmt:        "",
		ArrayHasSomeFmt:         "",
		ArrayIsEmptyFmt:         "({col} IS NULL OR json_array_length({col}) = 0)",
	}
}

type Queries struct {
	db        DBTX
	baseDB    DBTX
	provider  string
	dialect   Dialect
	stmtCache map[string]*sql.Stmt
	mu        *sync.RWMutex
	// Agent provides CRUD operations for Agent.
	//
	//   id           string    default: cuid()
	//   name         string    required
	//   description  string    required
	//   content      string    default:
	//   allowedTools string    required
	//   metadata     string    required
	//   createdAt    time.Time default: now()
	//   updatedAt    time.Time required
	Agent *AgentDelegate
	// AgentSkill provides CRUD operations for AgentSkill.
	//
	//   agentId  string   required
	//   skillId  string   required
	//   priority Priority default: NORMAL
	AgentSkill *AgentSkillDelegate
	// App provides CRUD operations for App.
	//
	//   name       string required
	//   project_id string required
	App *AppDelegate
	// AppSkill provides CRUD operations for AppSkill.
	//
	//   app_id   string   required
	//   skillId  string   required
	//   priority Priority default: NORMAL
	AppSkill *AppSkillDelegate
	// Document provides CRUD operations for Document.
	//
	//   id        string required
	//   kind      string default:
	//   body      []byte required
	//   updatedAt int32  required
	Document *DocumentDelegate
	// Reference provides CRUD operations for Reference.
	//
	//   name         string        required
	//   description  string        required
	//   accessType   URI           required
	//   resourceType ReferenceType required
	//   content      string        optional
	Reference *ReferenceDelegate
	// SkillReference provides CRUD operations for SkillReference.
	//
	//   skillId       string required
	//   referenceName string required
	SkillReference *SkillReferenceDelegate
	// Project provides CRUD operations for Project.
	//
	//   name       string required
	//   skillsPath string required
	//   agentsPath string required
	Project *ProjectDelegate
	// ProjectSkill provides CRUD operations for ProjectSkill.
	//
	//   projectName string required
	//   skillId     string required
	ProjectSkill *ProjectSkillDelegate
	// ProjectSpecialization provides CRUD operations for ProjectSpecialization.
	//
	//   generalizesName string required
	//   specializesName string required
	ProjectSpecialization *ProjectSpecializationDelegate
	// Skill provides CRUD operations for Skill.
	//
	//   id          string    default: cuid()
	//   name        string    required
	//   description string    required
	//   url         string    required
	//   content     string    required
	//   keywords    string    default:
	//   createdAt   time.Time default: now()
	//   updatedAt   time.Time required
	Skill *SkillDelegate
	// SkillDependency provides CRUD operations for SkillDependency.
	//
	//   dependencyId string required
	//   dependentId  string required
	SkillDependency *SkillDependencyDelegate
	// StyleToken provides CRUD operations for StyleToken.
	//
	//   style    string    required
	//   value    string    required
	//   mode     ThemeMode required
	//   type     StyleType required
	//   token    TokenType required
	//   app_id   string    required
	//   variable string    optional
	StyleToken *StyleTokenDelegate
	// Trash provides CRUD operations for Trash.
	//
	//   id        int32     default: autoincrement()
	//   createdAt time.Time default: now()
	//   content   string    required
	Trash *TrashDelegate
	// TrashAgent provides CRUD operations for TrashAgent.
	//
	//   trashId int32  required
	//   agentId string required
	TrashAgent *TrashAgentDelegate
	// TrashSkill provides CRUD operations for TrashSkill.
	//
	//   trashId int32  required
	//   skillId string required
	TrashSkill *TrashSkillDelegate
	// User provides CRUD operations for User.
	//
	//   id   string required
	//   name string required
	User *UserDelegate
}

type DB struct {
	*Queries
	sqlDB *sql.DB
}

func Open(provider, dataSourceName string) (*DB, error) {
	sqlDB, err := sql.Open(provider, dataSourceName)
	if err != nil {
		return nil, err
	}

	q := &Queries{
		db:        sqlDB,
		baseDB:    sqlDB,
		provider:  provider,
		dialect:   newDialect(),
		stmtCache: make(map[string]*sql.Stmt),
		mu:        &sync.RWMutex{},
	}
	q.initDelegates()
	return &DB{
		Queries: q,
		sqlDB:   sqlDB,
	}, nil
}

func (q *Queries) initDelegates() {
	q.Agent = &AgentDelegate{client: q}
	q.AgentSkill = &AgentSkillDelegate{client: q}
	q.App = &AppDelegate{client: q}
	q.AppSkill = &AppSkillDelegate{client: q}
	q.Document = &DocumentDelegate{client: q}
	q.Reference = &ReferenceDelegate{client: q}
	q.SkillReference = &SkillReferenceDelegate{client: q}
	q.Project = &ProjectDelegate{client: q}
	q.ProjectSkill = &ProjectSkillDelegate{client: q}
	q.ProjectSpecialization = &ProjectSpecializationDelegate{client: q}
	q.Skill = &SkillDelegate{client: q}
	q.SkillDependency = &SkillDependencyDelegate{client: q}
	q.StyleToken = &StyleTokenDelegate{client: q}
	q.Trash = &TrashDelegate{client: q}
	q.TrashAgent = &TrashAgentDelegate{client: q}
	q.TrashSkill = &TrashSkillDelegate{client: q}
	q.User = &UserDelegate{client: q}
}

func (q *Queries) copyHooksFrom(other *Queries) {
	q.Agent.extensions = make([]AgentExtension, len(other.Agent.extensions))
	copy(q.Agent.extensions, other.Agent.extensions)
	q.AgentSkill.extensions = make([]AgentSkillExtension, len(other.AgentSkill.extensions))
	copy(q.AgentSkill.extensions, other.AgentSkill.extensions)
	q.App.extensions = make([]AppExtension, len(other.App.extensions))
	copy(q.App.extensions, other.App.extensions)
	q.AppSkill.extensions = make([]AppSkillExtension, len(other.AppSkill.extensions))
	copy(q.AppSkill.extensions, other.AppSkill.extensions)
	q.Document.extensions = make([]DocumentExtension, len(other.Document.extensions))
	copy(q.Document.extensions, other.Document.extensions)
	q.Reference.extensions = make([]ReferenceExtension, len(other.Reference.extensions))
	copy(q.Reference.extensions, other.Reference.extensions)
	q.SkillReference.extensions = make([]SkillReferenceExtension, len(other.SkillReference.extensions))
	copy(q.SkillReference.extensions, other.SkillReference.extensions)
	q.Project.extensions = make([]ProjectExtension, len(other.Project.extensions))
	copy(q.Project.extensions, other.Project.extensions)
	q.ProjectSkill.extensions = make([]ProjectSkillExtension, len(other.ProjectSkill.extensions))
	copy(q.ProjectSkill.extensions, other.ProjectSkill.extensions)
	q.ProjectSpecialization.extensions = make([]ProjectSpecializationExtension, len(other.ProjectSpecialization.extensions))
	copy(q.ProjectSpecialization.extensions, other.ProjectSpecialization.extensions)
	q.Skill.extensions = make([]SkillExtension, len(other.Skill.extensions))
	copy(q.Skill.extensions, other.Skill.extensions)
	q.SkillDependency.extensions = make([]SkillDependencyExtension, len(other.SkillDependency.extensions))
	copy(q.SkillDependency.extensions, other.SkillDependency.extensions)
	q.StyleToken.extensions = make([]StyleTokenExtension, len(other.StyleToken.extensions))
	copy(q.StyleToken.extensions, other.StyleToken.extensions)
	q.Trash.extensions = make([]TrashExtension, len(other.Trash.extensions))
	copy(q.Trash.extensions, other.Trash.extensions)
	q.TrashAgent.extensions = make([]TrashAgentExtension, len(other.TrashAgent.extensions))
	copy(q.TrashAgent.extensions, other.TrashAgent.extensions)
	q.TrashSkill.extensions = make([]TrashSkillExtension, len(other.TrashSkill.extensions))
	copy(q.TrashSkill.extensions, other.TrashSkill.extensions)
	q.User.extensions = make([]UserExtension, len(other.User.extensions))
	copy(q.User.extensions, other.User.extensions)
}

func (db *DB) Close() error {
	return db.sqlDB.Close()
}

// Raw returns the underlying *sql.DB connection pool.
func (db *DB) Raw() *sql.DB {
	return db.sqlDB
}

// RunMigrations runs all pending migrations from the embedded folder.
func (db *DB) RunMigrations(ctx context.Context) error {
	log.Println("Running migrations...")
	if err := goose.SetDialect(db.provider); err != nil {
		return err
	}
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrationsFS)
	err := goose.UpContext(ctx, db.sqlDB, "migrations")
	if err != nil {
		log.Printf("Migrations failed: %v", err)
		return err
	}
	log.Println("Migrations completed successfully.")
	return nil
}

func (q *Queries) prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	q.mu.RLock()
	stmt, ok := q.stmtCache[query]
	q.mu.RUnlock()

	if !ok {
		var prepDB DBTX = q.db
		if q.baseDB != nil {
			prepDB = q.baseDB
		}
		var err error
		stmt, err = prepDB.PrepareContext(ctx, query)
		if err != nil {
			return nil, TranslateDBError(err)
		}

		q.mu.Lock()
		if existing, found := q.stmtCache[query]; found {
			stmt.Close()
			stmt = existing
		} else {
			q.stmtCache[query] = stmt
		}
		q.mu.Unlock()
	}

	if tx, ok := q.db.(*sql.Tx); ok {
		return tx.StmtContext(ctx, stmt), nil
	}
	return stmt, nil
}

func (q *Queries) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Printf("[%s] SQL Query: %s | Args: %v", strings.ToUpper(q.provider), query, args)
	stmt, err := q.prepare(ctx, query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		log.Printf("[%s] SQL Error: %v | Query: %s | Args: %v", strings.ToUpper(q.provider), err, query, args)
		return nil, TranslateDBError(err)
	}
	return res, nil
}

func (q *Queries) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Printf("[%s] SQL Exec: %s | Args: %v", strings.ToUpper(q.provider), query, args)
	stmt, err := q.prepare(ctx, query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		log.Printf("[%s] SQL Error: %v | Query: %s | Args: %v", strings.ToUpper(q.provider), err, query, args)
		return nil, TranslateDBError(err)
	}
	return res, nil
}

func (q *Queries) inTx() bool {
	_, ok := q.db.(*sql.Tx)
	return ok
}

func (q *Queries) transaction(ctx context.Context, fn func(txQ *Queries) error) error {
	if q.inTx() {
		return fn(q)
	}

	starter, ok := q.db.(interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	})
	if !ok {
		return fn(q)
	}
	log.Printf("[%s] SQL Begin Transaction", strings.ToUpper(q.provider))
	tx, err := starter.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Printf("[%s] SQL Rollback Transaction", strings.ToUpper(q.provider))
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txQueries := &Queries{
		db:        tx,
		baseDB:    q.baseDB,
		provider:  q.provider,
		dialect:   q.dialect,
		stmtCache: q.stmtCache,
		mu:        q.mu,
	}
	txQueries.initDelegates()
	txQueries.copyHooksFrom(q)

	if err := fn(txQueries); err != nil {
		log.Printf("[%s] SQL Rollback Transaction", strings.ToUpper(q.provider))
		_ = tx.Rollback()
		return err
	}
	log.Printf("[%s] SQL Commit Transaction", strings.ToUpper(q.provider))
	return tx.Commit()
}

type PredicateData struct {
	Column    string
	Operator  string
	Value     any
	IsLogical bool
	Children  []PredicateData
}

type ChildPredicate struct {
	Column string
	Value  any
}

type PredicateOf[M any] interface {
	ToPredicateData() PredicateData
	Validate() error
	Column() string
	Value() any
	Children() []ChildPredicate
	phantom(M)
}

type Predicate[M any] struct {
	Data PredicateData
}

func (p Predicate[M]) ToPredicateData() PredicateData {
	return p.Data
}

func (p Predicate[M]) Validate() error {
	return p.Data.Validate()
}

func (p Predicate[M]) Column() string {
	return p.Data.Column
}

func (p Predicate[M]) Value() any {
	return p.Data.Value
}

func (p Predicate[M]) Children() []ChildPredicate {
	if !p.Data.IsLogical || len(p.Data.Children) == 0 {
		return nil
	}
	res := make([]ChildPredicate, len(p.Data.Children))
	for i, c := range p.Data.Children {
		res[i] = ChildPredicate{Column: c.Column, Value: c.Value}
	}
	return res
}

func (p Predicate[M]) phantom(M) {}

type UniquePredicate[M any] struct {
	Data PredicateData
}

func (p UniquePredicate[M]) ToPredicateData() PredicateData {
	return p.Data
}

func (p UniquePredicate[M]) Validate() error {
	if p.Data.Column == "" && len(p.Data.Children) == 0 {
		return fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	return p.Data.Validate()
}

func (p UniquePredicate[M]) Column() string {
	return p.Data.Column
}

func (p UniquePredicate[M]) Value() any {
	return p.Data.Value
}

func (p UniquePredicate[M]) Children() []ChildPredicate {
	if !p.Data.IsLogical || len(p.Data.Children) == 0 {
		return nil
	}
	res := make([]ChildPredicate, len(p.Data.Children))
	for i, c := range p.Data.Children {
		res[i] = ChildPredicate{Column: c.Column, Value: c.Value}
	}
	return res
}

func (p UniquePredicate[M]) phantom(M) {}

func validateValue(col string, val any) error {
	switch v := val.(type) {
	case string:
		if strings.Contains(v, "\x00") {
			return &ValidationError{
				Errors: []FieldError{
					{Field: col, Value: v, Rule: "safety", Msg: "string cannot contain null bytes"},
				},
			}
		}
		if !utf8.ValidString(v) {
			return &ValidationError{
				Errors: []FieldError{
					{Field: col, Value: v, Rule: "safety", Msg: "string must be valid UTF-8"},
				},
			}
		}
	case []string:
		for _, s := range v {
			if err := validateValue(col, s); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range v {
			if err := validateValue(col, item); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pd PredicateData) Validate() error {
	if pd.IsLogical {
		for _, child := range pd.Children {
			if err := child.Validate(); err != nil {
				return err
			}
		}
		return nil
	}
	return validateValue(pd.Column, pd.Value)
}

func And[M any](preds ...PredicateOf[M]) PredicateOf[M] {
	var children []PredicateData
	for _, p := range preds {
		if p != nil {
			children = append(children, p.ToPredicateData())
		}
	}
	return Predicate[M]{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "AND",
			Children:  children,
		},
	}
}

func Or[M any](preds ...PredicateOf[M]) PredicateOf[M] {
	var children []PredicateData
	for _, p := range preds {
		if p != nil {
			children = append(children, p.ToPredicateData())
		}
	}
	return Predicate[M]{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "OR",
			Children:  children,
		},
	}
}

func Not[M any](pred PredicateOf[M]) PredicateOf[M] {
	var children []PredicateData
	if pred != nil {
		children = append(children, pred.ToPredicateData())
	}
	return Predicate[M]{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "NOT",
			Children:  children,
		},
	}
}

type Field[M any, T any] struct {
	Column string
}

func (f Field[M, T]) Set(val T) FieldAssignmentOf[M] {
	return FieldAssignmentOf[M]{Col: f.Column, Val: val}
}

func (f Field[M, T]) EQ(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) NEQ(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) GT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f Field[M, T]) GTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) LT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f Field[M, T]) LTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) In(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f Field[M, T]) NotIn(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "NOT IN",
			Value:    vals,
		},
	}
}

func (f Field[M, T]) Between(min T, max T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "BETWEEN",
			Value:    []any{min, max},
		},
	}
}

type OptionalField[M any, T any] struct {
	Field[M, T]
}

func (f OptionalField[M, T]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f OptionalField[M, T]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f Field[M, T]) Asc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Asc}
}

func (f Field[M, T]) Desc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Desc}
}

type UniqueField[M any, T any] struct {
	Column string
}

func (f UniqueField[M, T]) UniqueColumns() []string {
	return []string{f.Column}
}

func (f UniqueField[M, T]) Set(val T) FieldAssignmentOf[M] {
	return FieldAssignmentOf[M]{Col: f.Column, Val: val}
}

func (f UniqueField[M, T]) EQ(val T) UniquePredicate[M] {
	return UniquePredicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) NEQ(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) GT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) GTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) LT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) LTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) In(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f UniqueField[M, T]) NotIn(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "NOT IN",
			Value:    vals,
		},
	}
}

func (f UniqueField[M, T]) Between(min T, max T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "BETWEEN",
			Value:    []any{min, max},
		},
	}
}

type OptionalUniqueField[M any, T any] struct {
	UniqueField[M, T]
}

func (f OptionalUniqueField[M, T]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f OptionalUniqueField[M, T]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f UniqueField[M, T]) Asc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Asc}
}

func (f UniqueField[M, T]) Desc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Desc}
}

type StringField[M any] struct {
	Column string
}

func (f StringField[M]) Set(val string) FieldAssignmentOf[M] {
	return FieldAssignmentOf[M]{Col: f.Column, Val: val}
}

func (f StringField[M]) EQ(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f StringField[M]) NEQ(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f StringField[M]) GT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f StringField[M]) GTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f StringField[M]) LT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f StringField[M]) LTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f StringField[M]) In(vals []string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f StringField[M]) NotIn(vals []string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "NOT IN",
			Value:    vals,
		},
	}
}

func (f StringField[M]) Between(min string, max string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "BETWEEN",
			Value:    []any{min, max},
		},
	}
}

func (f StringField[M]) Like(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    val,
		},
	}
}

func (f StringField[M]) ILike(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "ILIKE",
			Value:    val,
		},
	}
}

func (f StringField[M]) HasPrefix(prefix string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    prefix + "%",
		},
	}
}

func (f StringField[M]) HasSuffix(suffix string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + suffix,
		},
	}
}

func (f StringField[M]) Contains(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + val + "%",
		},
	}
}

type OptionalStringField[M any] struct {
	StringField[M]
}

func (f OptionalStringField[M]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f OptionalStringField[M]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f StringField[M]) Asc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Asc}
}

func (f StringField[M]) Desc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Desc}
}

type StringUniqueField[M any] struct {
	Column string
}

func (f StringUniqueField[M]) UniqueColumns() []string {
	return []string{f.Column}
}

func (f StringUniqueField[M]) Set(val string) FieldAssignmentOf[M] {
	return FieldAssignmentOf[M]{Col: f.Column, Val: val}
}

func (f StringUniqueField[M]) EQ(val string) UniquePredicate[M] {
	return UniquePredicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) NEQ(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) GT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) GTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) LT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) LTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) In(vals []string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f StringUniqueField[M]) NotIn(vals []string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "NOT IN",
			Value:    vals,
		},
	}
}

func (f StringUniqueField[M]) Between(min string, max string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "BETWEEN",
			Value:    []any{min, max},
		},
	}
}

func (f StringUniqueField[M]) Like(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) ILike(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "ILIKE",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) HasPrefix(prefix string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    prefix + "%",
		},
	}
}

func (f StringUniqueField[M]) HasSuffix(suffix string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + suffix,
		},
	}
}

func (f StringUniqueField[M]) Contains(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + val + "%",
		},
	}
}

type OptionalStringUniqueField[M any] struct {
	StringUniqueField[M]
}

func (f OptionalStringUniqueField[M]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f OptionalStringUniqueField[M]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f StringUniqueField[M]) Asc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Asc}
}

func (f StringUniqueField[M]) Desc() OrderBy[M] {
	return OrderBy[M]{Field: f.Column, Direction: Desc}
}

type ArrayField[M any, T any] struct {
	Column string
}

func (f ArrayField[M, T]) Set(vals []T) FieldAssignmentOf[M] {
	return FieldAssignmentOf[M]{Col: f.Column, Val: vals}
}

func (f ArrayField[M, T]) Has(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "ARRAY_HAS",
			Value:    val,
		},
	}
}

func (f ArrayField[M, T]) HasEvery(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "ARRAY_HAS_EVERY",
			Value:    vals,
		},
	}
}

func (f ArrayField[M, T]) HasSome(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "ARRAY_HAS_SOME",
			Value:    vals,
		},
	}
}

func (f ArrayField[M, T]) IsEmpty() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "ARRAY_IS_EMPTY",
		},
	}
}

func CompilePredicates[M any](dialect Dialect, preds []PredicateOf[M], startBindIdx ...int) (string, []any, int) {
	bindIdx := 1
	if len(startBindIdx) > 0 && startBindIdx[0] > 0 {
		bindIdx = startBindIdx[0]
	}
	if len(preds) == 0 {
		return "", nil, bindIdx
	}
	var data []PredicateData
	for _, p := range preds {
		if p != nil {
			data = append(data, p.ToPredicateData())
		}
	}
	return CompilePredicateData(dialect, data, bindIdx)
}

func CompilePredicateData(dialect Dialect, data []PredicateData, startBindIdx ...int) (string, []any, int) {
	bindIdx := 1
	if len(startBindIdx) > 0 && startBindIdx[0] > 0 {
		bindIdx = startBindIdx[0]
	}
	if len(data) == 0 {
		return "", nil, bindIdx
	}
	var parts []string
	var args []any

	var compile func(p PredicateData) string
	compile = func(p PredicateData) string {
		if p.IsLogical {
			if len(p.Children) == 0 {
				return ""
			}
			if p.Operator == "NOT" {
				sub := compile(p.Children[0])
				if sub == "" {
					return ""
				}
				return fmt.Sprintf("NOT (%s)", sub)
			}
			var subParts []string
			for _, child := range p.Children {
				sub := compile(child)
				if sub != "" {
					subParts = append(subParts, sub)
				}
			}
			if len(subParts) == 0 {
				return ""
			}
			if len(subParts) == 1 {
				return subParts[0]
			}
			return fmt.Sprintf("(%s)", strings.Join(subParts, " "+p.Operator+" "))
		}

		switch p.Operator {
		case "IS NULL", "IS NOT NULL":
			return fmt.Sprintf("%s %s", dialect.Quote(p.Column), p.Operator)
		case "IN":
			valSlice := unpackSlice(p.Value)
			if len(valSlice) == 0 {
				return "1=0"
			}
			var placeHolders []string
			for range valSlice {
				placeHolders = append(placeHolders, dialect.BindVar(bindIdx))
				bindIdx++
			}
			args = append(args, valSlice...)
			return fmt.Sprintf("%s IN (%s)", dialect.Quote(p.Column), strings.Join(placeHolders, ", "))
		case "NOT IN":
			valSlice := unpackSlice(p.Value)
			if len(valSlice) == 0 {
				return "1=1"
			}
			var placeHolders []string
			for range valSlice {
				placeHolders = append(placeHolders, dialect.BindVar(bindIdx))
				bindIdx++
			}
			args = append(args, valSlice...)
			return fmt.Sprintf("%s NOT IN (%s)", dialect.Quote(p.Column), strings.Join(placeHolders, ", "))
		case "BETWEEN":
			valSlice := unpackSlice(p.Value)
			if len(valSlice) < 2 {
				return ""
			}
			p1 := dialect.BindVar(bindIdx)
			bindIdx++
			p2 := dialect.BindVar(bindIdx)
			bindIdx++
			args = append(args, valSlice[0], valSlice[1])
			return fmt.Sprintf("%s BETWEEN %s AND %s", dialect.Quote(p.Column), p1, p2)
		case "ILIKE":
			placeholder := dialect.BindVar(bindIdx)
			bindIdx++
			args = append(args, p.Value)
			return dialect.FormatILike(p.Column, placeholder)
		case "ARRAY_HAS":
			placeholder := dialect.BindVar(bindIdx)
			bindIdx++
			args = append(args, p.Value)
			return dialect.FormatArrayHas(p.Column, placeholder)
		case "ARRAY_HAS_EVERY":
			valSlice := unpackSlice(p.Value)
			if len(valSlice) == 0 {
				return "1=1"
			}
			var placeHolders []string
			for range valSlice {
				placeHolders = append(placeHolders, dialect.BindVar(bindIdx))
				bindIdx++
			}
			args = append(args, valSlice...)
			return dialect.FormatArrayHasEvery(p.Column, placeHolders)
		case "ARRAY_HAS_SOME":
			valSlice := unpackSlice(p.Value)
			if len(valSlice) == 0 {
				return "1=0"
			}
			var placeHolders []string
			for range valSlice {
				placeHolders = append(placeHolders, dialect.BindVar(bindIdx))
				bindIdx++
			}
			args = append(args, valSlice...)
			return dialect.FormatArrayHasSome(p.Column, placeHolders)
		case "ARRAY_IS_EMPTY":
			return dialect.FormatArrayIsEmpty(p.Column)
		default:
			placeholder := dialect.BindVar(bindIdx)
			bindIdx++
			args = append(args, p.Value)
			return fmt.Sprintf("%s %s %s", dialect.Quote(p.Column), p.Operator, placeholder)
		}
	}

	for _, p := range data {
		part := compile(p)
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return "", nil, bindIdx
	}
	return strings.Join(parts, " AND "), args, bindIdx
}

func unpackSlice(val any) []any {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case []string:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []int:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []int32:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []int64:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []float32:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []float64:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []bool:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []time.Time:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case [][]byte:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []json.RawMessage:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []any:
		return v
	case []URIType:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []ReferenceTypeType:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []ThemeModeType:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []StyleTypeType:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []TokenTypeType:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []PriorityType:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	default:
		return []any{val}
	}
}

type Tx struct {
	*Queries
	tx *sql.Tx
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	sqlTx, err := db.sqlDB.BeginTx(ctx, opts)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	log.Printf("[%s] SQL Begin Transaction", strings.ToUpper(db.provider))
	q := &Queries{
		db:        sqlTx,
		baseDB:    db.Queries.baseDB,
		provider:  db.provider,
		dialect:   db.dialect,
		stmtCache: db.Queries.stmtCache,
		mu:        db.Queries.mu,
	}
	q.initDelegates()
	q.copyHooksFrom(db.Queries)
	return &Tx{
		Queries: q,
		tx:      sqlTx,
	}, nil
}

func (tx *Tx) Commit() error {
	log.Printf("[%s] SQL Commit Transaction", strings.ToUpper(tx.provider))
	return TranslateDBError(tx.tx.Commit())
}

func (tx *Tx) Rollback() error {
	log.Printf("[%s] SQL Rollback Transaction", strings.ToUpper(tx.provider))
	return TranslateDBError(tx.tx.Rollback())
}

// Raw returns the underlying *sql.Tx transaction handle.
func (tx *Tx) Raw() *sql.Tx {
	return tx.tx
}

// Transaction starts a transaction and runs the provided function.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (db *DB) Transaction(ctx context.Context, fn func(tx *Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w (rollback failed: %v)", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateBuilder[M any, S any, O any] struct {
	assignments    []FieldAssignment
	selects        *S
	omits          *O
	execFunc       func(ctx context.Context, assignments []FieldAssignment, s *S, o *O, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*M, error)
	conflictAction *ConflictAction
	conflictTarget UniqueConstraintTarget
}

func (b *CreateBuilder[M, S, O]) Select(s S) *CreateBuilder[M, S, O] {
	b.selects = &s
	return b
}

func (b *CreateBuilder[M, S, O]) Omit(o O) *CreateBuilder[M, S, O] {
	b.omits = &o
	return b
}

func (b *CreateBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.assignments, b.selects, b.omits, b.conflictTarget, b.conflictAction)
}

type CreateManyBuilder[M any] struct {
	records        []RecordInput
	execFunc       func(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error)
	conflictAction *ConflictAction
	conflictTarget UniqueConstraintTarget
}

func (b *CreateManyBuilder[M]) SkipDuplicates() *CreateManyBuilder[M] {
	b.conflictAction = &ConflictAction{Type: ConflictActionIgnore}
	return b
}

func (b *CreateManyBuilder[M]) Exec(ctx context.Context) (int64, error) {
	return b.execFunc(ctx, b.records, b.conflictTarget, b.conflictAction)
}

type CreateManyAndReturnBuilder[M any, S any, O any] struct {
	records        []RecordInput
	selects        *S
	omits          *O
	execFunc       func(ctx context.Context, records []RecordInput, s *S, o *O, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*M, error)
	conflictAction *ConflictAction
	conflictTarget UniqueConstraintTarget
}

func (b *CreateManyAndReturnBuilder[M, S, O]) SkipDuplicates() *CreateManyAndReturnBuilder[M, S, O] {
	b.conflictAction = &ConflictAction{Type: ConflictActionIgnore}
	return b
}

func (b *CreateManyAndReturnBuilder[M, S, O]) Select(s S) *CreateManyAndReturnBuilder[M, S, O] {
	b.selects = &s
	return b
}

func (b *CreateManyAndReturnBuilder[M, S, O]) Omit(o O) *CreateManyAndReturnBuilder[M, S, O] {
	b.omits = &o
	return b
}

func (b *CreateManyAndReturnBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.execFunc(ctx, b.records, b.selects, b.omits, b.conflictTarget, b.conflictAction)
}

type UpdateBuilder[M any, S any, O any] struct {
	where       UniquePredicate[M]
	additional  []PredicateOf[M]
	assignments []FieldAssignment
	selects     *S
	omits       *O
	execFunc    func(ctx context.Context, where UniquePredicate[M], additional []PredicateOf[M], assignments []FieldAssignment, s *S, o *O) (*M, error)
}

func (b *UpdateBuilder[M, S, O]) Select(selects S) *UpdateBuilder[M, S, O] {
	b.selects = &selects
	return b
}

func (b *UpdateBuilder[M, S, O]) Omit(omits O) *UpdateBuilder[M, S, O] {
	b.omits = &omits
	return b
}

func (b *UpdateBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.where, b.additional, b.assignments, b.selects, b.omits)
}

type UpdateManyBuilder[M any] struct {
	where       []PredicateOf[M]
	assignments []FieldAssignment
	execFunc    func(ctx context.Context, where []PredicateOf[M], assignments []FieldAssignment) (int64, error)
}

func (b *UpdateManyBuilder[M]) Exec(ctx context.Context) (int64, error) {
	return b.execFunc(ctx, b.where, b.assignments)
}

type UpdateManyAndReturnBuilder[M any, S any, O any] struct {
	where       []PredicateOf[M]
	assignments []FieldAssignment
	selects     *S
	omits       *O
	execFunc    func(ctx context.Context, where []PredicateOf[M], assignments []FieldAssignment, s *S, o *O) ([]*M, error)
}

func (b *UpdateManyAndReturnBuilder[M, S, O]) Select(selects S) *UpdateManyAndReturnBuilder[M, S, O] {
	b.selects = &selects
	return b
}

func (b *UpdateManyAndReturnBuilder[M, S, O]) Omit(omits O) *UpdateManyAndReturnBuilder[M, S, O] {
	b.omits = &omits
	return b
}

func (b *UpdateManyAndReturnBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.execFunc(ctx, b.where, b.assignments, b.selects, b.omits)
}

type FindUniqueBuilder[M any, S any, O any] struct {
	where      UniquePredicate[M]
	additional []PredicateOf[M]
	selects    *S
	omits      *O
	execFunc   func(ctx context.Context, where UniquePredicate[M], additional []PredicateOf[M], s *S, o *O) (*M, error)
}

func (b *FindUniqueBuilder[M, S, O]) Select(s S) *FindUniqueBuilder[M, S, O] {
	b.selects = &s
	return b
}

func (b *FindUniqueBuilder[M, S, O]) Omit(o O) *FindUniqueBuilder[M, S, O] {
	b.omits = &o
	return b
}

func (b *FindUniqueBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.where, b.additional, b.selects, b.omits)
}

type FindFirstBuilder[M any, S any, O any] struct {
	where    []PredicateOf[M]
	take     *int
	skip     *int
	orderBy  []OrderBy[M]
	cursor   UniquePredicate[M]
	selects  *S
	omits    *O
	execFunc func(ctx context.Context, params QueryParams[M], s *S, o *O) (*M, error)
}

func (b *FindFirstBuilder[M, S, O]) Take(limit int) *FindFirstBuilder[M, S, O] {
	b.take = &limit
	return b
}

func (b *FindFirstBuilder[M, S, O]) Skip(offset int) *FindFirstBuilder[M, S, O] {
	b.skip = &offset
	return b
}

func (b *FindFirstBuilder[M, S, O]) OrderBy(orders ...OrderBy[M]) *FindFirstBuilder[M, S, O] {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *FindFirstBuilder[M, S, O]) Cursor(where UniquePredicate[M]) *FindFirstBuilder[M, S, O] {
	b.cursor = where
	return b
}

func (b *FindFirstBuilder[M, S, O]) Select(s S) *FindFirstBuilder[M, S, O] {
	b.selects = &s
	return b
}

func (b *FindFirstBuilder[M, S, O]) Omit(o O) *FindFirstBuilder[M, S, O] {
	b.omits = &o
	return b
}

func (b *FindFirstBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	params := QueryParams[M]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
	return b.execFunc(ctx, params, b.selects, b.omits)
}

type FindManyBuilder[M any, S any, O any] struct {
	where    []PredicateOf[M]
	take     *int
	skip     *int
	orderBy  []OrderBy[M]
	cursor   UniquePredicate[M]
	selects  *S
	omits    *O
	execFunc func(ctx context.Context, params QueryParams[M], s *S, o *O) ([]*M, error)
}

func (b *FindManyBuilder[M, S, O]) Take(limit int) *FindManyBuilder[M, S, O] {
	b.take = &limit
	return b
}

func (b *FindManyBuilder[M, S, O]) Skip(offset int) *FindManyBuilder[M, S, O] {
	b.skip = &offset
	return b
}

func (b *FindManyBuilder[M, S, O]) OrderBy(orders ...OrderBy[M]) *FindManyBuilder[M, S, O] {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *FindManyBuilder[M, S, O]) Cursor(where UniquePredicate[M]) *FindManyBuilder[M, S, O] {
	b.cursor = where
	return b
}

func (b *FindManyBuilder[M, S, O]) Select(s S) *FindManyBuilder[M, S, O] {
	b.selects = &s
	return b
}

func (b *FindManyBuilder[M, S, O]) Omit(o O) *FindManyBuilder[M, S, O] {
	b.omits = &o
	return b
}

func (b *FindManyBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	params := QueryParams[M]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
	return b.execFunc(ctx, params, b.selects, b.omits)
}

type DeleteManyBuilder[M any] struct {
	where    []PredicateOf[M]
	execFunc func(ctx context.Context, where []PredicateOf[M]) (int64, error)
}

func (b *DeleteManyBuilder[M]) Exec(ctx context.Context) (int64, error) {
	return b.execFunc(ctx, b.where)
}

type DeleteBuilder[M any, S any, O any] struct {
	where      UniquePredicate[M]
	additional []PredicateOf[M]
	selects    *S
	omits      *O
	execFunc   func(ctx context.Context, where UniquePredicate[M], additional []PredicateOf[M], selects *S, omits *O) (*M, error)
}

func (b *DeleteBuilder[M, S, O]) Select(selects S) *DeleteBuilder[M, S, O] {
	b.selects = &selects
	return b
}

func (b *DeleteBuilder[M, S, O]) Omit(omits O) *DeleteBuilder[M, S, O] {
	b.omits = &omits
	return b
}

func (b *DeleteBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.where, b.additional, b.selects, b.omits)
}

type CountBuilder[M any] struct {
	where    []PredicateOf[M]
	take     *int
	skip     *int
	execFunc func(ctx context.Context, params QueryParams[M]) (int64, error)
}

func (b *CountBuilder[M]) Take(limit int) *CountBuilder[M] {
	b.take = &limit
	return b
}

func (b *CountBuilder[M]) Skip(offset int) *CountBuilder[M] {
	b.skip = &offset
	return b
}

func (b *CountBuilder[M]) Exec(ctx context.Context) (int64, error) {
	params := QueryParams[M]{
		Where: b.where,
		Take:  b.take,
		Skip:  b.skip,
	}
	return b.execFunc(ctx, params)
}
func directKey[T any, K any](get func(*T) K) func(*T) (string, bool) {
	return func(t *T) (string, bool) {
		return fmt.Sprint(get(t)), true
	}
}

func optionalKey[T any, K any](get func(*T) *K) func(*T) (string, bool) {
	return func(t *T) (string, bool) {
		if p := get(t); p != nil {
			return fmt.Sprint(*p), true
		}
		return "", false
	}
}

func setOne[P any, C any](set func(*P, *C)) func(*P, []*C) {
	return func(p *P, children []*C) {
		if len(children) > 0 {
			set(p, children[0])
		}
	}
}

func appendMany[P any, C any](get func(*P) *[]*C) func(*P, []*C) {
	return func(p *P, children []*C) {
		if s := get(p); s != nil {
			*s = append(*s, children...)
		}
	}
}

func scanInto[C any](cols []string, scan func(*C, []string) []any) func(*sql.Rows, *C) error {
	return func(rows *sql.Rows, c *C) error {
		return rows.Scan(scan(c, cols)...)
	}
}

type colSpec struct {
	col      string
	selected bool
	omitted  bool
	forceIn  bool
}

func computeCols(specs []colSpec, hasSelects, anySelected bool) []string {
	var cols []string
	for _, s := range specs {
		include := true
		if hasSelects {
			include = !anySelected || s.selected || s.forceIn
		} else if s.omitted {
			include = false
		}
		if include {
			cols = append(cols, s.col)
		}
	}
	if len(cols) == 0 {
		for _, s := range specs {
			cols = append(cols, s.col)
		}
	}
	return cols
}

func reverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func formatOrderBySQL[M any](dialect Dialect, orders []OrderBy[M], pkCols []string, uniqueCols []string, isCursorQuery bool, take *int) string {
	isNegativeTake := take != nil && *take < 0
	if len(orders) == 0 {
		if !isCursorQuery || len(pkCols) == 0 {
			return ""
		}
		var sb strings.Builder
		sb.WriteString(" ORDER BY ")
		for i, pkCol := range pkCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			dialect.WriteQuotedIdent(&sb, pkCol)
			if isNegativeTake {
				sb.WriteString(" DESC")
			} else {
				sb.WriteString(" ASC")
			}
		}
		return sb.String()
	}

	var sb strings.Builder
	usedCols := make(map[string]bool, len(orders)+len(pkCols))
	hasUnique := false
	uniqueMap := make(map[string]bool, len(uniqueCols))
	for _, u := range uniqueCols {
		uniqueMap[u] = true
	}

	count := 0
	for _, ord := range orders {
		if ord.Field == "" {
			continue
		}
		if count == 0 {
			sb.WriteString(" ORDER BY ")
		} else {
			sb.WriteString(", ")
		}
		dialect.WriteQuotedIdent(&sb, ord.Field)
		sb.WriteString(" ")
		dir := ord.Direction
		if isNegativeTake {
			if dir == Asc {
				dir = Desc
			} else {
				dir = Asc
			}
		}
		sb.WriteString(string(dir))
		usedCols[ord.Field] = true
		if uniqueMap[ord.Field] {
			hasUnique = true
		}
		count++
	}

	if isCursorQuery && !hasUnique {
		for _, pkCol := range pkCols {
			if !usedCols[pkCol] {
				if count == 0 {
					sb.WriteString(" ORDER BY ")
				} else {
					sb.WriteString(", ")
				}
				dialect.WriteQuotedIdent(&sb, pkCol)
				if isNegativeTake {
					sb.WriteString(" DESC")
				} else {
					sb.WriteString(" ASC")
				}
				usedCols[pkCol] = true
				count++
			}
		}
	}
	return sb.String()
}

func compileCursorClause[M any](
	dialect Dialect,
	cursor UniquePredicate[M],
	orders []OrderBy[M],
	pkCols []string,
	uniqueCols []string,
	tableName string,
	nextIdx int,
	take *int,
) (cursorClause string, cursorVals []any, err error) {
	if cursor.Data.Column == "" && len(cursor.Data.Children) == 0 {
		return "", nil, nil
	}
	if err := cursor.Validate(); err != nil {
		return "", nil, err
	}
	cursorWhere, cVals, _ := CompilePredicates(dialect, []PredicateOf[M]{cursor}, nextIdx)
	if cursorWhere == "" {
		return "", nil, nil
	}

	uniqueMap := make(map[string]bool, len(uniqueCols))
	for _, u := range uniqueCols {
		uniqueMap[u] = true
	}

	var orderedCols []string
	used := make(map[string]bool, len(orders)+len(pkCols))
	hasUnique := false
	for _, ord := range orders {
		if ord.Field != "" && !used[ord.Field] {
			orderedCols = append(orderedCols, ord.Field)
			used[ord.Field] = true
			if uniqueMap[ord.Field] {
				hasUnique = true
			}
		}
	}

	if !hasUnique {
		for _, pk := range pkCols {
			if !used[pk] {
				orderedCols = append(orderedCols, pk)
				used[pk] = true
			}
		}
	}

	var subSb strings.Builder
	subSb.WriteString("(")
	for i, col := range orderedCols {
		if i > 0 {
			subSb.WriteString(", ")
		}
		dialect.WriteQuotedIdent(&subSb, col)
	}
	isNegativeTake := take != nil && *take < 0
	cursorOp := ">"
	if len(orders) > 0 && orders[0].Direction == Desc {
		cursorOp = "<"
	}
	if isNegativeTake {
		if cursorOp == ">" {
			cursorOp = "<"
		} else {
			cursorOp = ">"
		}
	}

	subSb.WriteString(") ")
	subSb.WriteString(cursorOp)
	subSb.WriteString(" (SELECT ")
	for i, col := range orderedCols {
		if i > 0 {
			subSb.WriteString(", ")
		}
		dialect.WriteQuotedIdent(&subSb, col)
	}
	subSb.WriteString(" FROM ")
	dialect.WriteQuotedIdent(&subSb, tableName)
	subSb.WriteString(" WHERE ")
	subSb.WriteString(cursorWhere)
	subSb.WriteString(")")

	return subSb.String(), cVals, nil
}

func buildSelectSQL(q *Queries, table string, cols []string, whereClause string, orderByClause string, take *int, skip *int) string {
	var sb strings.Builder
	sb.Grow(len(table) + len(whereClause) + len(orderByClause) + len(cols)*18 + 48)
	sb.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WriteQuotedIdent(&sb, col)
	}
	sb.WriteString(" FROM ")
	q.dialect.WriteQuotedIdent(&sb, table)
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if take != nil {
		sb.WriteString(" LIMIT ")
		limit := *take
		if limit < 0 {
			limit = -limit
		}
		sb.WriteString(strconv.Itoa(limit))
		if skip != nil {
			sb.WriteString(" OFFSET ")
			sb.WriteString(strconv.Itoa(*skip))
		}
	} else if skip != nil {
		if q.dialect.SupportsLimitMinusOne {
			sb.WriteString(" LIMIT -1 OFFSET ")
		} else {
			sb.WriteString(" OFFSET ")
		}
		sb.WriteString(strconv.Itoa(*skip))
	}
	return sb.String()
}

func buildSingleInsertSQL(
	q *Queries,
	table string,
	cols []string,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
	valsCount int,
) (string, []any) {
	var sb strings.Builder
	sb.Grow(128 + len(table)*2 + len(cols)*15 + len(returningCols)*15)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, table)
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WriteQuotedIdent(&sb, col)
	}
	sb.WriteString(") VALUES (")
	for i := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WritePlaceholder(&sb, i+1)
	}
	sb.WriteString(")")

	var conflictCols []string
	if conflictTarget != nil {
		conflictCols = conflictTarget.UniqueColumns()
	}
	var nonConflictCols []string
	if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
		nonConflictCols = computeNonConflictCols(cols, conflictCols, pkCols)
	}
	clause, clauseArgs := q.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, valsCount+1)
	sb.WriteString(clause)

	if q.dialect.SupportsInsertReturning && len(returningCols) > 0 {
		sb.WriteString(" RETURNING ")
		for i, col := range returningCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			q.dialect.WriteQuotedIdent(&sb, col)
		}
	}

	return sb.String(), clauseArgs
}

func loadRelation[P any, C any](
	ctx context.Context,
	q *Queries,
	parents []*P,
	parentKey func(*P) (string, bool),
	table string,
	fkCol string,
	returningCols []string,
	scan func(*sql.Rows, *C) error,
	childKey func(*C) (string, bool),
	assign func(*P, []*C),
	params QueryParams[C],
) ([]*C, error) {
	var parentKeys []any
	for _, p := range parents {
		if p == nil {
			continue
		}
		if key, ok := parentKey(p); ok {
			parentKeys = append(parentKeys, key)
		}
	}
	if len(parentKeys) == 0 {
		return nil, nil
	}

	// Prepend parent ID checks to filters using Predicate[C]
	allPreds := append([]PredicateOf[C]{
		Predicate[C]{
			Data: PredicateData{
				Column:    fkCol,
				Operator:  "IN",
				Value:     parentKeys,
				IsLogical: false,
			},
		},
	}, params.Where...)

	whereClause, vals, nextIdx := CompilePredicates(q.dialect, allPreds)
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if isCursorQuery {
		cClause, cVals, err := compileCursorClause(q.dialect, params.Cursor, params.OrderBy, []string{"id"}, nil, table, nextIdx, params.Take)
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

	query := compileRelationSQL(q.dialect, table, fkCol, returningCols, whereClause, params)

	rows, err := q.query(ctx, query, vals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	childMap := make(map[string][]*C, len(parents))
	allChildren := make([]*C, 0, len(parents))

	for rows.Next() {
		var child C
		if err := scan(rows, &child); err != nil {
			return nil, TranslateDBError(err)
		}
		if key, ok := childKey(&child); ok {
			childMap[key] = append(childMap[key], &child)
		}
		allChildren = append(allChildren, &child)
	}
	if err := rows.Err(); err != nil {
		return nil, TranslateDBError(err)
	}

	for _, p := range parents {
		if p == nil {
			continue
		}
		if key, ok := parentKey(p); ok {
			assign(p, childMap[key])
		}
	}

	return allChildren, nil
}

// m2mParentCol is the output alias Phi appends to a many-to-many preload so the
// join table's parent key travels alongside the target columns and rows can be
// grouped back onto the parents they belong to.
const m2mParentCol = "_m2m_parent"

// loadManyToMany preloads a many-to-many relation that Phi materialized as a
// join table. It selects the target rows, joins the join table on the target's
// primary key, filters on the join table's parent column, and groups the
// children back onto each parent.
//
// Target columns are referenced unqualified (matching loadRelation) and the
// join table is aliased, so only the join columns are qualified. Because the
// query fans out over parents, the parent key is selected as an extra trailing
// column and the target row is scanned together with it. Pagination is applied
// per parent with a window function, exactly as loadRelation does.
func loadManyToMany[P any, C any](
	ctx context.Context,
	q *Queries,
	parents []*P,
	parentKey func(*P) (string, bool),
	targetTable string,
	targetKeyCol string,
	joinTable string,
	joinParentCol string,
	joinTargetCol string,
	returningCols []string,
	scan func(*C, []string) []any,
	assign func(*P, []*C),
	params QueryParams[C],
) ([]*C, error) {
	var parentKeys []any
	for _, p := range parents {
		if p == nil {
			continue
		}
		if key, ok := parentKey(p); ok {
			parentKeys = append(parentKeys, key)
		}
	}
	if len(parentKeys) == 0 {
		return nil, nil
	}

	const joinAlias = "_m2m"

	// Parent filter lives on the join table, so it is qualified; user filters
	// address the target table and stay unqualified. The parent bind values come
	// first, so user predicates are compiled from the next bind index.
	var whereSb strings.Builder
	whereSb.WriteString(q.dialect.Quote(joinAlias))
	whereSb.WriteString(".")
	whereSb.WriteString(q.dialect.Quote(joinParentCol))
	whereSb.WriteString(" IN (")
	for i := range parentKeys {
		if i > 0 {
			whereSb.WriteString(", ")
		}
		whereSb.WriteString(q.dialect.BindVar(i + 1))
	}
	whereSb.WriteString(")")

	vals := make([]any, 0, len(parentKeys)+4)
	vals = append(vals, parentKeys...)

	if len(params.Where) > 0 {
		userClause, userVals, _ := CompilePredicates(q.dialect, params.Where, len(parentKeys)+1)
		if userClause != "" {
			whereSb.WriteString(" AND (")
			whereSb.WriteString(userClause)
			whereSb.WriteString(")")
			vals = append(vals, userVals...)
		}
	}

	query := compileManyToManySQL(q.dialect, targetTable, targetKeyCol, joinTable, joinAlias, joinParentCol, joinTargetCol, returningCols, " WHERE "+whereSb.String(), params)

	rows, err := q.query(ctx, query, vals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	childMap := make(map[string][]*C, len(parents))
	allChildren := make([]*C, 0, len(parents))

	for rows.Next() {
		var child C
		var parentVal string
		targets := make([]any, 0, len(returningCols)+1)
		targets = append(targets, scan(&child, returningCols)...)
		targets = append(targets, &parentVal)
		if err := rows.Scan(targets...); err != nil {
			return nil, TranslateDBError(err)
		}
		childMap[parentVal] = append(childMap[parentVal], &child)
		allChildren = append(allChildren, &child)
	}
	if err := rows.Err(); err != nil {
		return nil, TranslateDBError(err)
	}

	for _, p := range parents {
		if p == nil {
			continue
		}
		if key, ok := parentKey(p); ok {
			assign(p, childMap[key])
		}
	}

	return allChildren, nil
}

// compileManyToManySQL builds the join query behind loadManyToMany. Without
// pagination it is a plain SELECT over the join; with pagination it wraps a
// ROW_NUMBER() window partitioned by the join table's parent column so the
// limit applies to each parent rather than the whole result set.
func compileManyToManySQL[M any](dialect Dialect, targetTable, targetKeyCol, joinTable, joinAlias, joinParentCol, joinTargetCol string, cols []string, where string, params QueryParams[M]) string {
	prefixed := func(col string) string {
		return dialect.Quote(joinAlias) + "." + dialect.Quote(col)
	}
	ownCol := prefixed(joinParentCol)

	// Order columns address the target model, so they are qualified by the
	// target table to stay unambiguous next to the joined join table.
	orderBy := ""
	if len(params.OrderBy) > 0 {
		var ob strings.Builder
		ob.WriteString(" ORDER BY ")
		for i, ord := range params.OrderBy {
			if i > 0 {
				ob.WriteString(", ")
			}
			ob.WriteString(dialect.Quote(targetTable))
			ob.WriteString(".")
			ob.WriteString(dialect.Quote(ord.Field))
			ob.WriteString(" ")
			ob.WriteString(string(ord.Direction))
		}
		orderBy = ob.String()
	}

	joinClause := " FROM " + dialect.Quote(targetTable) +
		" JOIN " + dialect.Quote(joinTable) + " AS " + dialect.Quote(joinAlias) +
		" ON " + prefixed(joinTargetCol) + " = " + dialect.Quote(targetTable) + "." + dialect.Quote(targetKeyCol)

	if params.Take == nil && params.Skip == nil {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		for i, col := range cols {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(dialect.Quote(col))
		}
		sb.WriteString(", ")
		sb.WriteString(ownCol)
		sb.WriteString(" AS ")
		sb.WriteString(dialect.Quote(m2mParentCol))
		sb.WriteString(joinClause)
		sb.WriteString(where)
		sb.WriteString(orderBy)
		return sb.String()
	}

	rowOrder := orderBy
	if rowOrder == "" {
		rowOrder = " ORDER BY " + ownCol + " ASC"
	}

	var inner strings.Builder
	inner.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			inner.WriteString(", ")
		}
		inner.WriteString(dialect.Quote(col))
	}
	inner.WriteString(", ")
	inner.WriteString(ownCol)
	inner.WriteString(" AS ")
	inner.WriteString(dialect.Quote(m2mParentCol))
	inner.WriteString(", ROW_NUMBER() OVER (PARTITION BY ")
	inner.WriteString(ownCol)
	inner.WriteString(rowOrder)
	inner.WriteString(") AS row_num")
	inner.WriteString(joinClause)
	inner.WriteString(where)

	var outer strings.Builder
	outer.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			outer.WriteString(", ")
		}
		outer.WriteString(dialect.Quote(col))
	}
	outer.WriteString(", ")
	outer.WriteString(dialect.Quote(m2mParentCol))
	outer.WriteString(" FROM (")
	outer.WriteString(inner.String())
	outer.WriteString(") t WHERE ")

	if params.Take != nil && params.Skip != nil {
		outer.WriteString(fmt.Sprintf("row_num > %d AND row_num <= %d", *params.Skip, *params.Skip+*params.Take))
	} else if params.Take != nil {
		outer.WriteString(fmt.Sprintf("row_num <= %d", *params.Take))
	} else {
		outer.WriteString(fmt.Sprintf("row_num > %d", *params.Skip))
	}
	return outer.String()
}

func compileRelationSQL[M any](dialect Dialect, table, fkCol string, cols []string, where string, params QueryParams[M]) string {
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if params.Take != nil || params.Skip != nil || isCursorQuery {
		return compilePartitionedRelationSQL(dialect, table, fkCol, cols, where, params)
	}
	return compileSimpleRelationSQL(dialect, table, cols, where, params)
}

func compilePartitionedRelationSQL[M any](dialect Dialect, table, fkCol string, cols []string, where string, params QueryParams[M]) string {
	var innerSb strings.Builder
	innerSb.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			innerSb.WriteString(", ")
		}
		innerSb.WriteString(dialect.Quote(col))
	}
	innerSb.WriteString(", ROW_NUMBER() OVER (PARTITION BY ")
	innerSb.WriteString(dialect.Quote(fkCol))
	innerSb.WriteString(" ORDER BY ")
	if len(params.OrderBy) > 0 {
		for i, ord := range params.OrderBy {
			if i > 0 {
				innerSb.WriteString(", ")
			}
			innerSb.WriteString(dialect.Quote(ord.Field))
			innerSb.WriteString(" ")
			innerSb.WriteString(string(ord.Direction))
		}
	} else {
		innerSb.WriteString(dialect.Quote("id"))
		innerSb.WriteString(" ASC")
	}
	innerSb.WriteString(") as row_num FROM ")
	innerSb.WriteString(dialect.Quote(table))
	innerSb.WriteString(where)

	var outerSb strings.Builder
	outerSb.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			outerSb.WriteString(", ")
		}
		outerSb.WriteString(dialect.Quote(col))
	}
	outerSb.WriteString(" FROM (")
	outerSb.WriteString(innerSb.String())
	outerSb.WriteString(") t WHERE ")

	if params.Take != nil && params.Skip != nil {
		outerSb.WriteString(fmt.Sprintf("row_num > %d AND row_num <= %d", *params.Skip, *params.Skip+*params.Take))
	} else if params.Take != nil {
		outerSb.WriteString(fmt.Sprintf("row_num <= %d", *params.Take))
	} else if params.Skip != nil {
		outerSb.WriteString(fmt.Sprintf("row_num > %d", *params.Skip))
	}
	return outerSb.String()
}

func compileSimpleRelationSQL[M any](dialect Dialect, table string, cols []string, where string, params QueryParams[M]) string {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dialect.Quote(col))
	}
	sb.WriteString(" FROM ")
	sb.WriteString(dialect.Quote(table))
	sb.WriteString(where)
	if len(params.OrderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		for i, ord := range params.OrderBy {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(dialect.Quote(ord.Field))
			sb.WriteString(" ")
			sb.WriteString(string(ord.Direction))
		}
	}
	return sb.String()
}
