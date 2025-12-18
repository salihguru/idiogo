package query

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type V[T any] []T

type Conds struct {
	Key    string
	Values V[any]
	Skip   bool
}

const (
	AND = "AND"
	OR  = "OR"

	SortDesc = "DESC"
	SortAsc  = "ASC"
)

// Apply creates a GORM function that applies the given conditions to the query
// It builds the query string and values from the provided conditions and applies them to the GORM DB instance
// The conditions are combined using the specified operator (AND or OR)
// db.Clause(query.Apply(conds, opr...))
func Apply(conds []Conds, opr ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		conds, vals := Build(conds, opr...)
		if conds != "" {
			return db.Where(conds, vals...)
		}
		return db
	}
}

//	conds, vals := query.Build([]query.Conds{
//		{
//			Key:    "sender_id = ? OR receiver_id = ?",
//			Values: query.V[any]{accountId, accountId},
//			Skip:   false,
//		},
//	})
//
// Build([]Conds) will return a query string and a slice of values that can be used in the QueryContext method.
func Build(conds []Conds, opr ...string) (string, []interface{}) {
	if len(conds) == 0 {
		return "", nil
	}
	op := getOption(AND, opr...)
	var query string
	var values []interface{}
	var limitIdx int
	for idx, cond := range conds {
		if cond.Skip {
			continue
		}
		if strings.Contains(cond.Key, "LIMIT") {
			limitIdx = idx
			continue
		}
		if len(cond.Values) > 0 && cond.Values[0] != "" {
			query += fmt.Sprintf("%s %s ", cond.Key, op)
			values = append(values, cond.Values...)
		} else if !strings.Contains(cond.Key, "?") {
			query += fmt.Sprintf("%s %s ", cond.Key, op)
		}
	}

	if len(query) == 0 {
		return "", nil
	}
	if op == AND {
		query = query[:len(query)-5]
	} else {
		query = query[:len(query)-4]
	}
	if limitIdx > 0 {
		query += fmt.Sprintf(" %s", conds[limitIdx].Key)
		values = append(values, conds[limitIdx].Values...)
	}

	return query, values
}

func skipCond() Conds {
	return Conds{Skip: true}
}

func IntToBool(k string, i int) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s = ?", k),
		Values: V[any]{i == 1},
		Skip:   i == 0,
	}
}

func Like(k string, v string) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s LIKE ?", k),
		Values: V[any]{"%" + v + "%"},
		Skip:   v == "",
	}
}

func ILike(k string, v string) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s ILIKE ?", k),
		Values: V[any]{"%" + v + "%"},
		Skip:   v == "",
	}
}

// ILikeMulti creates an ILIKE condition for multiple fields with OR operator
// fields: slice of field names to search in
// value: single search value to search for in all fields
func ILikeMulti(fields []string, value string) Conds {
	if len(fields) == 0 || value == "" {
		return Conds{
			Key:    "",
			Values: V[any]{},
			Skip:   true,
		}
	}

	// Build the query like: "field1 ILIKE ? OR field2 ILIKE ? OR field3 ILIKE ?"
	conditions := make([]string, len(fields))
	values := make([]interface{}, len(fields))

	for i, field := range fields {
		conditions[i] = fmt.Sprintf("%s ILIKE ?", field)
		values[i] = "%" + value + "%"
	}

	return Conds{
		Key:    strings.Join(conditions, " OR "),
		Values: values,
		Skip:   false,
	}
}

// ILikeMultiValues creates an ILIKE condition for multiple fields with multiple values
// Each field will be searched against each value with OR operator
// fields: slice of field names to search in
// values: slice of values to search for
func ILikeMultiValues(fields []string, values []string) Conds {
	if len(fields) == 0 || len(values) == 0 {
		return Conds{
			Key:    "",
			Values: V[any]{},
			Skip:   true,
		}
	}

	// Build conditions for each field-value combination
	var conditions []string
	var queryValues []interface{}

	for _, field := range fields {
		for _, value := range values {
			if value != "" {
				conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
				queryValues = append(queryValues, "%"+value+"%")
			}
		}
	}

	if len(conditions) == 0 {
		return Conds{
			Key:    "",
			Values: V[any]{},
			Skip:   true,
		}
	}

	return Conds{
		Key:    strings.Join(conditions, " OR "),
		Values: queryValues,
		Skip:   false,
	}
}

func IntGreaterOrEqual(k string, i int) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s >= ?", k),
		Values: V[any]{i},
		Skip:   i == 0,
	}
}

func ReplacePlaceholder(q string, start ...int) string {
	s := 0
	if len(start) > 0 {
		s = start[0]
	}
	parts := strings.Split(q, "?")
	for i := 0; i < len(parts)-1; i++ {
		parts[i] += fmt.Sprintf("$%d", s+1)
		s++
	}
	return strings.Join(parts, "")
}

func Geo(k string, long, lat, radius float64) Conds {
	return Conds{
		Key:    fmt.Sprintf("ST_DWithin(%s, ST_Point(?, ?)::geography, ?)", k),
		Values: V[any]{long, lat, radius},
		Skip:   long == 0 || lat == 0 || radius == 0,
	}
}

func OrderGeo(k string, long, lat float64) string {
	return fmt.Sprintf("ST_Distance( %s, ST_Point(%f, %f))", k, long, lat)
}

func Text(k string, v string) Conds {
	return Conds{
		Key: fmt.Sprintf("to_tsvector('simple', %s) @@ to_tsquery('simple', ?)", k),
		Values: V[any]{
			strings.Join(strings.Fields(v), " & "),
		},
		Skip: v == "",
	}
}

// TextPrefix creates a text search condition with prefix matching support
// This allows partial word matching by adding :* to each search term
// Example: TextPrefix("title || ' ' || description", "vil") will match "villa"
func TextPrefix(k string, v string) Conds {
	if v == "" {
		return skipCond()
	}

	// Split search terms and add :* to each for prefix matching
	words := strings.Fields(v)
	prefixWords := make([]string, len(words))
	for i, word := range words {
		prefixWords[i] = word + ":*"
	}

	return Conds{
		Key: fmt.Sprintf("to_tsvector('simple', %s) @@ to_tsquery('simple', ?)", k),
		Values: V[any]{
			strings.Join(prefixWords, " & "),
		},
		Skip: false,
	}
}

func Eq(k string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s = ?", k),
		Values: V[any]{v},
		Skip:   getOption(v == nil, skip...),
	}
}

func StrArr(k string, v interface{}, skip ...bool) Conds {
	// Convert the slice to PostgreSQL array format
	if strSlice, ok := v.([]string); ok && len(strSlice) > 0 {
		// Create PostgreSQL array literal format: ARRAY['value1','value2']
		quotedValues := make([]string, len(strSlice))
		for i, s := range strSlice {
			// Escape single quotes in values
			quotedValues[i] = "'" + strings.ReplaceAll(s, "'", "''") + "'"
		}
		arrayLiteral := "ARRAY[" + strings.Join(quotedValues, ",") + "]::text[]"

		return Conds{
			Key:    fmt.Sprintf("%s && %s", k, arrayLiteral),
			Values: V[any]{},
			Skip:   getOption(false, skip...),
		}
	}

	return Conds{
		Key:    "",
		Values: V[any]{},
		Skip:   true,
	}
}

func IsEmptyUUID(id uuid.UUID) bool {
	return id == uuid.Nil || id.String() == ""
}

func Custom[T any](k string, v T, skip ...bool) Conds {
	return Conds{
		Key:    k,
		Values: V[any]{v},
		Skip:   getOption(false, skip...),
	}
}

func NotEq(k string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s != ?", k),
		Values: V[any]{v},
		Skip:   getOption(v == nil, skip...),
	}
}

func In[T any](k string, v []T, skip ...bool) Conds {
	if len(v) == 0 {
		return Conds{
			Key:    fmt.Sprintf("%s IS NULL", k),
			Values: V[any]{},
			Skip:   getOption(true, skip...),
		}
	}
	placeholders := strings.Repeat("?,", len(v)-1) + "?"
	values := make([]any, len(v))
	for i, val := range v {
		values[i] = val
	}
	return Conds{
		Key:    fmt.Sprintf("%s IN (%s)", k, placeholders),
		Values: values,
		Skip:   getOption(false, skip...),
	}
}

func InSeperated(k string, v []string, skip ...bool) Conds {
	if len(v) == 0 || v[0] == "" {
		return Conds{
			Key:    fmt.Sprintf("%s IS NULL", k),
			Values: V[any]{},
			Skip:   getOption(true, skip...),
		}
	}
	placeholders := strings.Repeat("?,", len(v)-1) + "?"
	return Conds{
		Key:    fmt.Sprintf("%s IN (%s)", k, placeholders),
		Values: V[any]{v},
		Skip:   getOption(false, skip...),
	}
}

func Min(k string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s >= ?", k),
		Values: V[any]{v},
		Skip:   getOption(v == nil, skip...),
	}
}

func Max(k string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s <= ?", k),
		Values: V[any]{v},
		Skip:   getOption(v == nil, skip...),
	}
}

func NotNull(k string, skip bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s IS NOT NULL", k),
		Values: V[any]{},
		Skip:   skip,
	}
}

func NotIn[T any](k string, v []T, skip ...bool) Conds {
	if len(v) == 0 {
		return Conds{
			Key:    "",
			Values: V[any]{},
			Skip:   true,
		}
	}
	placeholders := strings.Repeat("?,", len(v)-1) + "?"
	return Conds{
		Key:    fmt.Sprintf("%s NOT IN (%s)", k, placeholders),
		Values: V[any]{v},
		Skip:   getOption(false, skip...),
	}
}

func NotInSeperated(k string, v []string, skip ...bool) Conds {
	if len(v) == 0 || v[0] == "" {
		return Conds{
			Key:    "",
			Values: V[any]{},
			Skip:   getOption(true, skip...),
		}
	}
	placeholders := strings.Repeat("?,", len(v)-1) + "?"
	return Conds{
		Key:    fmt.Sprintf("%s NOT IN (%s)", k, placeholders),
		Values: V[any]{v},
		Skip:   getOption(false, skip...),
	}
}

// ArrayContains checks if a single value is contained in a PostgreSQL array field
// Example: "? = ANY(tags)" where tags is an array column
func ArrayContains(k string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("? = ANY(%s)", k),
		Values: V[any]{v},
		Skip:   getOption(v == nil || v == "", skip...),
	}
}

// ArrayEquals checks if a PostgreSQL array field is equal to a given value
// Example: "tags = ARRAY['tag1','tag2']"
func ArrayEquals(k string, v []string, skip ...bool) Conds {
	if len(v) == 0 {
		return skipCond()
	}
	varlen := strings.Repeat("?,", len(v)-1) + "?"
	values := make([]any, len(v))
	for i, val := range v {
		values[i] = val
	}
	return Conds{
		Key:    fmt.Sprintf("%s && ARRAY[%s]::text[]", k, varlen),
		Values: values,
		Skip:   getOption(len(v) == 0, skip...),
	}
}

// JsonbField creates a condition for JSONB field access
// Example: JsonbField("address", "city", "New York") => "address->>'city' = ?"
func JsonbField(k string, field string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s->>'%s' = ?", k, field),
		Values: V[any]{v},
		Skip:   getOption(v == nil || v == "", skip...),
	}
}

// JsonbFieldNullSafe creates a null-safe condition for JSONB field access
// Example: JsonbFieldNullSafe("address", "city", "New York") => "(address IS NOT NULL AND address->>'city' = ?)"
func JsonbFieldNullSafe(k string, field string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("(%s IS NOT NULL AND %s->>'%s' = ?)", k, k, field),
		Values: V[any]{v},
		Skip:   getOption(v == nil || v == "", skip...),
	}
}

// JsonbFieldILike creates an ILIKE condition for JSONB field
func JsonbFieldILike(k string, field string, v string, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("%s->>'%s' ILIKE ?", k, field),
		Values: V[any]{"%" + v + "%"},
		Skip:   getOption(v == "", skip...),
	}
}

// JsonbFieldILikeNullSafe creates a null-safe ILIKE condition for JSONB field
func JsonbFieldILikeNullSafe(k string, field string, v string, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("(%s IS NOT NULL AND %s->>'%s' ILIKE ?)", k, k, field),
		Values: V[any]{"%" + v + "%"},
		Skip:   getOption(v == "", skip...),
	}
}

// JsonbNestedField creates a condition for nested JSONB field access
// Example: JsonbNestedField("config", []string{"translation", "tr", "title"}, "value") => "config->'translation'->'tr'->>'title' = ?"
func JsonbNestedField(k string, path []string, v interface{}, skip ...bool) Conds {
	if len(path) == 0 {
		return skipCond()
	}

	var pathBuilder strings.Builder
	pathBuilder.WriteString(k)

	for i, segment := range path {
		if i == len(path)-1 {
			// Last element uses ->> for text extraction
			pathBuilder.WriteString("->>'")
			pathBuilder.WriteString(segment)
			pathBuilder.WriteString("'")
		} else {
			// Intermediate elements use -> for JSON navigation
			pathBuilder.WriteString("->'")
			pathBuilder.WriteString(segment)
			pathBuilder.WriteString("'")
		}
	}

	return Conds{
		Key:    fmt.Sprintf("%s = ?", pathBuilder.String()),
		Values: V[any]{v},
		Skip:   getOption(v == nil || v == "", skip...),
	}
}

// JsonbNestedFieldILike creates an ILIKE condition for nested JSONB field
func JsonbNestedFieldILike(k string, path []string, v string, skip ...bool) Conds {
	if len(path) == 0 {
		return skipCond()
	}

	var pathBuilder strings.Builder
	pathBuilder.WriteString(k)

	for i, segment := range path {
		if i == len(path)-1 {
			// Last element uses ->> for text extraction
			pathBuilder.WriteString("->>'")
			pathBuilder.WriteString(segment)
			pathBuilder.WriteString("'")
		} else {
			// Intermediate elements use -> for JSON navigation
			pathBuilder.WriteString("->'")
			pathBuilder.WriteString(segment)
			pathBuilder.WriteString("'")
		}
	}

	return Conds{
		Key:    fmt.Sprintf("%s ILIKE ?", pathBuilder.String()),
		Values: V[any]{"%" + v + "%"},
		Skip:   getOption(v == "", skip...),
	}
}

// JsonbMultiFieldsILike creates an OR condition for multiple JSONB paths with ILIKE
// Example: JsonbMultiFieldsILike("translation", [][]string{{"tr", "title"}, {"en", "title"}, {"tr", "description"}, {"en", "description"}}, "search")
// Results in: (translation->'tr'->>'title' ILIKE ? OR translation->'en'->>'title' ILIKE ? OR ...)
func JsonbMultiFieldsILike(k string, paths [][]string, v string, skip ...bool) Conds {
	if len(paths) == 0 || v == "" {
		return skipCond()
	}

	var conditions []string
	var values []interface{}

	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		var pathBuilder strings.Builder
		pathBuilder.WriteString(k)

		for i, segment := range path {
			if i == len(path)-1 {
				// Last element uses ->> for text extraction
				pathBuilder.WriteString("->>'")
				pathBuilder.WriteString(segment)
				pathBuilder.WriteString("'")
			} else {
				// Intermediate elements use -> for JSON navigation
				pathBuilder.WriteString("->'")
				pathBuilder.WriteString(segment)
				pathBuilder.WriteString("'")
			}
		}

		conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", pathBuilder.String()))
		values = append(values, "%"+v+"%")
	}

	if len(conditions) == 0 {
		return skipCond()
	}

	return Conds{
		Key:    "(" + strings.Join(conditions, " OR ") + ")",
		Values: values,
		Skip:   getOption(false, skip...),
	}
}

// JsonbArrayOverlap creates a condition to check if JSONB array contains any of the given values
// Example: JsonbArrayOverlap("audience", "interest", []string{"hiking", "camping"}) => "audience->'interest' ?| ARRAY['hiking','camping']"
func JsonbArrayOverlap(k string, field string, values []string, skip ...bool) Conds {
	if len(values) == 0 {
		return skipCond()
	}

	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = "'" + strings.ReplaceAll(v, "'", "''") + "'"
	}

	return Conds{
		Key:    fmt.Sprintf("%s->'%s' ?| ARRAY[%s]::text[]", k, field, strings.Join(quoted, ",")),
		Values: V[any]{},
		Skip:   getOption(false, skip...),
	}
}

// JsonbAgeRange creates a condition to check if a value falls within a JSONB age range array
// Example: JsonbAgeRange("audience", 25) => "(audience->'age_range'->0)::int <= 25 AND (audience->'age_range'->1)::int >= 25"
func JsonbAgeRange(k string, age int, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("(%s->'age_range'->0)::int <= %d AND (%s->'age_range'->1)::int >= %d", k, age, k, age),
		Values: V[any]{},
		Skip:   getOption(age == 0, skip...),
	}
}

// PersonalizationScore creates a scoring expression for personalization
type PersonalizationScoreField struct {
	Type   string      // "interests", "badges", "gender", "age"
	Value  interface{} // The value to match against
	Points int         // Points to award if matched
}

// BuildPersonalizationScore creates a complex personalization scoring expression
func BuildPersonalizationScore(jsonbField string, fields []PersonalizationScoreField) string {
	if len(fields) == 0 {
		return "0"
	}

	var expressions []string

	for _, field := range fields {
		switch field.Type {
		case "interest", "badges":
			if values, ok := field.Value.([]string); ok && len(values) > 0 {
				quoted := make([]string, len(values))
				for i, v := range values {
					quoted[i] = "'" + strings.ReplaceAll(v, "'", "''") + "'"
				}
				expr := fmt.Sprintf("CASE WHEN %s->'%s' ?| ARRAY[%s]::text[] THEN %d ELSE 0 END",
					jsonbField, field.Type, strings.Join(quoted, ","), field.Points)
				expressions = append(expressions, expr)
			}
		case "gender":
			if value, ok := field.Value.(string); ok && value != "" {
				expr := fmt.Sprintf("CASE WHEN %s->>'gender' = '%s' THEN %d ELSE 0 END",
					jsonbField, strings.ReplaceAll(value, "'", "''"), field.Points)
				expressions = append(expressions, expr)
			}
		case "age":
			if age, ok := field.Value.(int); ok && age > 0 {
				expr := fmt.Sprintf("CASE WHEN (%s->'age_range'->0)::int <= %d AND (%s->'age_range'->1)::int >= %d THEN %d ELSE 0 END",
					jsonbField, age, jsonbField, age, field.Points)
				expressions = append(expressions, expr)
			}
		}
	}

	if len(expressions) == 0 {
		return "0"
	}

	return "(" + strings.Join(expressions, " + ") + ")"
}

// JsonbNumericMin creates a condition for JSONB numeric field minimum value
// Example: JsonbNumericMin("review", "average_point", 4.5) => "(review->>'average_point')::float >= 4.5"
func JsonbNumericMin(k string, field string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("(%s->>'%s')::float >= ?", k, field),
		Values: V[any]{v},
		Skip:   getOption(v == nil, skip...),
	}
}

// JsonbNumericMax creates a condition for JSONB numeric field maximum value
// Example: JsonbNumericMax("review", "average_point", 4.5) => "(review->>'average_point')::float <= 4.5"
func JsonbNumericMax(k string, field string, v interface{}, skip ...bool) Conds {
	return Conds{
		Key:    fmt.Sprintf("(%s->>'%s')::float <= ?", k, field),
		Values: V[any]{v},
		Skip:   getOption(v == nil, skip...),
	}
}

// StrArrWithPrefix checks if a PostgreSQL array field contains any of the given values with a prefix
// This is useful for tags where database stores '#tag' but search terms come without '#'
// Example: StrArrWithPrefix("tags", []string{"sapanca", "bolu"}, "#") will match against ['#sapanca', '#bolu']
func StrArrWithPrefix(k string, v interface{}, prefix string, skip ...bool) Conds {
	// Convert the slice to PostgreSQL array format with prefix
	if strSlice, ok := v.([]string); ok && len(strSlice) > 0 {
		// Create PostgreSQL array literal format with prefix: ARRAY['#value1','#value2']
		quotedValues := make([]string, len(strSlice))
		for i, s := range strSlice {
			// Add prefix and escape single quotes in values
			quotedValues[i] = "'" + prefix + strings.ReplaceAll(s, "'", "''") + "'"
		}
		arrayLiteral := "ARRAY[" + strings.Join(quotedValues, ",") + "]::text[]"

		return Conds{
			Key:    fmt.Sprintf("%s && %s", k, arrayLiteral),
			Values: V[any]{},
			Skip:   getOption(false, skip...),
		}
	}
	return skipCond()
}

func getOption[V any](v V, opts ...V) V {
	if len(opts) > 0 {
		return opts[0]
	}
	return v
}
