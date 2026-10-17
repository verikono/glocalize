package graphql

import (
	"fmt"
	"strings"

	GT "github.com/verikono/glocalize/types/go"
)

// ParsePrisma converts a Prisma schema document into a GraphQL schema document.
func ParsePrisma(schema GT.PrismaSchema) (*GT.GraphQLSchema, error) {

	models, err := extractModels(schema)
	if err != nil {
		return nil, err
	}
	enums := extractEnums(schema)

	ctx := parserContext{models: models, enums: enums}
	parsedEnums := parseEnums(ctx)
	parsedModels, err := parseModels(ctx)
	if err != nil {
		return nil, err
	}

	parsedScalars := parseScalars()
	parts := make([]string, 0, len(parsedScalars)+len(parsedEnums)+len(parsedModels))
	parts = append(parts, parsedScalars...)
	parts = append(parts, parsedEnums...)
	parts = append(parts, parsedModels...)

	out := strings.Join(parts, "\n\n")

	return &out, nil
}

// parseScalars emits the fixed GraphQL scalar declarations.
func parseScalars() []string {
	scalars := make([]string, 0, len(scalarOrder))
	for _, scalar := range scalarOrder {
		scalars = append(scalars, "scalar "+scalar)
	}
	return scalars
}

// extractModels walks a Prisma document and returns its model blocks.
func extractModels(pris string) ([]*ParsableModelShape, error) {
	lines := strings.Split(pris, "\n")
	var models []*ParsableModelShape
	var current *ParsableModelShape
	var docblocks []string
	var propblocks []string
	inmodel := false

	for _, line := range lines {
		if modelStartRe.MatchString(line) {
			c := strings.ReplaceAll(line, " ", "")
			name := c[5:]
			if open := strings.Index(c, "{"); open > 5 {
				name = c[5:open]
			}
			current = &ParsableModelShape{Name: name, Docs: cloneStrings(docblocks)}
			models = append(models, current)
			inmodel = true
			docblocks = nil
			continue
		}

		// Model-level documentation: every comment line accumulates here, so a
		// comment before a model becomes its docs. Comments inside a model are
		// cleared again at the model's closing brace.
		if commentLine.MatchString(line) {
			docblocks = append(docblocks, stripLeading(line, "/ "))
		}

		if modelEndRe.MatchString(line) {
			docblocks = nil
			inmodel = false
		}

		if attributeRe.MatchString(line) {
			continue
		}

		if inmodel && line != "" && propertyDirectiveRe.MatchString(line) {
			fields := strings.Fields(line)
			rest := ""
			if len(fields) > 2 {
				rest = strings.Join(fields[2:], " ")
			}
			prop, err := parsePropertyLine("  "+rest, nil)
			if err != nil {
				return nil, err
			}
			current.Props = append(current.Props, prop)
			propblocks = nil
		}

		if inmodel && line != "" {
			if commentLine.MatchString(line) {
				propblocks = append(propblocks, stripLeading(line, "/ "))
				continue
			}
			if strings.TrimSpace(line) == "" {
				continue
			}
			prop, err := parsePropertyLine(line, propblocks)
			if err != nil {
				return nil, err
			}
			propblocks = nil
			current.Props = append(current.Props, prop)
		}
	}
	return models, nil
}

// extractEnums walks a Prisma document and returns its enum blocks.
func extractEnums(pris string) []*ParsableEnumShape {
	lines := strings.Split(pris, "\n")
	var enums []*ParsableEnumShape
	var current *ParsableEnumShape
	inenum := false

	for _, line := range lines {
		if enumStartRe.MatchString(line) {
			c := strings.ReplaceAll(line, " ", "")
			name := c[4:]
			if open := strings.Index(c, "{"); open > 4 {
				name = c[4:open]
			}
			current = &ParsableEnumShape{Name: name}
			enums = append(enums, current)
			inenum = true
			continue
		}

		if modelEndRe.MatchString(line) {
			inenum = false
		}

		if inenum && line != "" {
			if commentLine.MatchString(line) {
				continue
			}
			option := dropSpaces(line)
			if option != "" {
				current.Options = append(current.Options, option)
			}
		}
	}
	return enums
}

// parsePropertyLine parses a single Prisma field line into a property shape.
func parsePropertyLine(line string, propblocks []string) (ParsablePropertyShape, error) {
	directives := []string{}
	fieldLine := line

	if strings.Contains(line, "///") {
		parts := strings.Split(line, "///")
		fieldLine = parts[0]
		directives = strings.Fields(strings.TrimSpace(strings.Join(parts[1:], "///")))
	}

	fields := strings.Fields(fieldLine)
	if len(fields) == 0 {
		return ParsablePropertyShape{}, fmt.Errorf("invalid name extracted")
	}
	if len(fields) == 1 {
		return ParsablePropertyShape{}, fmt.Errorf("invalid type extracted")
	}

	name := fields[0]
	typ := fields[1]

	prismaDirectives := make([]string, 0, len(fields)-2)
	for _, d := range fields[2:] {
		if strings.HasPrefix(d, "@") {
			prismaDirectives = append(prismaDirectives, d)
		}
	}
	allDirectives := append(append([]string{}, prismaDirectives...), directives...)

	optional := strings.HasSuffix(typ, "?")
	multi := strings.HasSuffix(typ, "[]") || containsString(allDirectives, "parse:JSONARRAY")

	scalar := ""
	for _, d := range allDirectives {
		if strings.HasPrefix(d, "parse:SCALAR") {
			scalar = d
			break
		}
	}

	params := make([]string, 0)
	for _, d := range allDirectives {
		if strings.HasPrefix(d, "parse:ARGUE") {
			params = append(params, d)
		}
	}

	finalType := typ
	if scalar != "" {
		// "parse:SCALAR(WEI)" -> "WEI"
		segs := strings.FieldsFunc(scalar, func(r rune) bool { return r == '(' || r == ')' })
		if len(segs) > 1 {
			finalType = segs[1]
		}
	} else if strings.HasSuffix(finalType, "?") {
		finalType = finalType[:len(finalType)-1]
	} else if strings.HasSuffix(finalType, "[]") {
		finalType = finalType[:len(finalType)-2]
	}

	return ParsablePropertyShape{
		Name:       name,
		Type:       finalType,
		Optional:   optional,
		Multi:      multi,
		Params:     params,
		Directives: allDirectives,
		Docblocks:  cloneStrings(propblocks),
	}, nil
}

// parseModels renders every model block as a GraphQL type.
func parseModels(ctx parserContext) ([]string, error) {
	parsed := make([]string, 0, len(ctx.models))
	for _, model := range ctx.models {
		stack := make([]string, 0, len(model.Props)+3)

		if len(model.Docs) > 0 {
			stack = append(stack, `"""`)
			stack = append(stack, model.Docs...)
			stack = append(stack, `"""`)
		}
		stack = append(stack, "type "+model.Name+" {")

		for _, prop := range model.Props {
			if skipNameRe.MatchString(prop.Name) {
				continue
			}
			rendered, err := parseModelProp(prop, ctx)
			if err != nil {
				return nil, err
			}
			stack = append(stack, "\t"+strings.Join(rendered, "\n\t"))
		}

		stack = append(stack, "}")
		parsed = append(parsed, strings.Join(stack, "\n"))
	}
	return parsed, nil
}

// parseModelProp renders a single field into its GraphQL lines.
func parseModelProp(prop ParsablePropertyShape, ctx parserContext) ([]string, error) {
	gqlType := gqlTypeMap[strings.ToLower(prop.Type)]

	var modelName string
	for _, m := range ctx.models {
		if m.Name == prop.Type {
			modelName = m.Name
			break
		}
	}
	var enumName string
	for _, e := range ctx.enums {
		if e.Name == prop.Type {
			enumName = e.Name
			break
		}
	}

	if gqlType == "" && modelName == "" && enumName == "" {
		return nil, fmt.Errorf("unknown model prop %q", prop.Type)
	}

	finalType := gqlType
	if finalType == "" {
		finalType = modelName
	}
	if finalType == "" {
		finalType = enumName
	}

	var stack []string
	if len(prop.Docblocks) > 0 {
		stack = append(stack, `"""`)
		for _, block := range prop.Docblocks {
			stack = append(stack, stripLeading(block, "/ \t"))
		}
		stack = append(stack, `"""`)
	}

	argsSignature := ""
	for _, param := range prop.Params {
		if strings.HasPrefix(strings.ToLower(param), "parse:argue") {
			argsSignature = renderArgs(param)
		}
	}

	multiOpen, multiClose := "", ""
	if prop.Multi {
		multiOpen, multiClose = "[", "]"
	}
	bang := ""
	if !prop.Optional {
		bang = "!"
	}

	stack = append(stack, fmt.Sprintf("%s%s: %s%s%s%s", prop.Name, argsSignature, multiOpen, finalType, multiClose, bang))
	return stack, nil
}

// renderArgs turns a "parse:ARGUE[...]" directive into a GraphQL argument signature.
func renderArgs(param string) string {
	open := strings.Index(param, "[")
	closeIdx := strings.LastIndex(param, "]")
	if open < 0 || closeIdx < open {
		return ""
	}
	var parts []string
	for _, group := range strings.Split(param[open+1:closeIdx], "][") {
		parts = append(parts, "("+renderArgSegment(group)+")")
	}
	return strings.Join(parts, "")
}

// renderArgSegment renders one bracket group of "name:type" arguments.
func renderArgSegment(seg string) string {
	pairs := strings.Split(seg, ",")
	rendered := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			continue
		}
		name := kv[0]
		typ := kv[1]
		mandatory := strings.HasSuffix(typ, "!") || strings.HasSuffix(name, "!")
		typ = strings.TrimSuffix(typ, "!")
		name = strings.TrimSuffix(name, "!")
		gql := gqlTypeMap[strings.ToLower(typ)]
		if gql == "" {
			gql = typ
		}
		suffix := ""
		if mandatory {
			suffix = "!"
		}
		rendered = append(rendered, name+":"+gql+suffix)
	}
	return strings.Join(rendered, ",")
}

// parseEnums renders every enum block as a GraphQL enum.
func parseEnums(ctx parserContext) []string {
	parsed := make([]string, 0, len(ctx.enums))
	for _, en := range ctx.enums {
		stack := []string{"enum " + en.Name + " {"}
		for _, option := range en.Options {
			if skipNameRe.MatchString(option) {
				continue
			}
			stack = append(stack, "\t"+option)
		}
		stack = append(stack, "}")
		parsed = append(parsed, strings.Join(stack, "\n"))
	}
	return parsed
}

// stripLeading returns s with the leading run of characters in excluded
// removed. When every character is excluded it returns s unchanged, matching
// the reference's String.prototype.substring(-1) behaviour.
func stripLeading(s, excluded string) string {
	for i := 0; i < len(s); i++ {
		if !strings.ContainsRune(excluded, rune(s[i])) {
			return s[i:]
		}
	}
	return s
}

// cloneStrings returns a defensive copy of s.
func cloneStrings(s []string) []string {
	if s == nil {
		return nil
	}
	return append([]string(nil), s...)
}

// containsString reports whether list contains target.
func containsString(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// dropSpaces returns line with every space rune removed.
func dropSpaces(line string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' {
			return -1
		}
		return r
	}, line)
}

// ParseGQL converts a GraphQL schema document back into a Prisma schema
// document. It is the inverse of ParsePrisma: GraphQL types become Prisma
// models, GraphQL enums become Prisma enums, argument signatures become
// "/// parse:ARGUE[...]" directives, and the WEI scalar becomes a Decimal
// field carrying "/// parse:SCALAR(WEI)". Prisma attributes (@id, @relation,
// @@map, ...) are not recoverable from GraphQL and are therefore omitted.
func ParseGQL(schema GT.GraphQLSchema) (*GT.PrismaSchema, error) {
	lines := strings.Split(schema, "\n")

	var blocks []gqlBlock
	var pendingDocs []string

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		switch {
		case trimmed == `"""`:
			var docs []string
			i++
			for i < len(lines) && strings.TrimSpace(lines[i]) != `"""` {
				docs = append(docs, strings.TrimSpace(lines[i]))
				i++
			}
			i++ // closing """
			pendingDocs = append(pendingDocs, docs...)

		case strings.HasPrefix(trimmed, "scalar "):
			i++

		case strings.HasPrefix(trimmed, "enum ") && strings.HasSuffix(trimmed, "{"):
			name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "enum "), "{"))
			b := gqlBlock{kind: "enum", name: name, docs: cloneStrings(pendingDocs)}
			i++
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "}") {
				t := strings.TrimSpace(lines[i])
				if t != "" && t != `"""` {
					b.options = append(b.options, t)
				}
				i++
			}
			i++ // closing }
			blocks = append(blocks, b)
			pendingDocs = nil

		case strings.HasPrefix(trimmed, "type ") && strings.HasSuffix(trimmed, "{"):
			name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "type "), "{"))
			b := gqlBlock{kind: "model", name: name, docs: cloneStrings(pendingDocs)}
			pendingDocs = nil
			var fieldDocs []string
			i++
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "}") {
				t := strings.TrimSpace(lines[i])
				switch {
				case t == `"""`:
					fieldDocs = nil
					i++
					for i < len(lines) && strings.TrimSpace(lines[i]) != `"""` {
						fieldDocs = append(fieldDocs, strings.TrimSpace(lines[i]))
						i++
					}
					i++ // closing """
				case t != "":
					b.fields = append(b.fields, gqlField{docs: cloneStrings(fieldDocs), line: t})
					fieldDocs = nil
					i++
				default:
					i++
				}
			}
			i++ // closing }
			blocks = append(blocks, b)

		default:
			i++
		}
	}

	parts := make([]string, 0, len(blocks))
	for _, b := range blocks {
		switch b.kind {
		case "enum":
			parts = append(parts, renderPrismaEnum(b))
		case "model":
			rendered, err := renderPrismaModel(b)
			if err != nil {
				return nil, err
			}
			parts = append(parts, rendered)
		}
	}

	out := strings.Join(parts, "\n\n")
	return &out, nil
}

// renderPrismaEnum renders a GraphQL enum as a Prisma enum block.
func renderPrismaEnum(b gqlBlock) string {
	var sb strings.Builder
	sb.WriteString("enum " + b.name + " {\n")
	for _, o := range b.options {
		sb.WriteString("\t" + o + "\n")
	}
	sb.WriteString("}")
	return sb.String()
}

// renderPrismaModel renders a GraphQL type as a Prisma model block.
func renderPrismaModel(b gqlBlock) (string, error) {
	var sb strings.Builder
	for _, d := range b.docs {
		sb.WriteString("// " + d + "\n")
	}
	sb.WriteString("model " + b.name + " {\n")
	for _, f := range b.fields {
		rendered, err := renderPrismaField(f)
		if err != nil {
			return "", err
		}
		for _, l := range rendered {
			sb.WriteString("\t" + l + "\n")
		}
	}
	sb.WriteString("}")
	return sb.String(), nil
}

// renderPrismaField renders a GraphQL field as Prisma field lines.
func renderPrismaField(f gqlField) ([]string, error) {
	line := strings.TrimSpace(f.line)

	var out []string
	for _, d := range f.docs {
		out = append(out, "// "+d)
	}

	argsRaw := ""
	rest := line
	if open := strings.Index(line, "("); open >= 0 {
		closeIdx := strings.Index(line[open:], ")")
		if closeIdx < 0 {
			return nil, fmt.Errorf("prismagql: malformed field %q", line)
		}
		argsRaw = line[open : open+closeIdx+1]
		rest = strings.TrimSpace(line[:open] + line[open+closeIdx+1:])
	}

	colon := strings.Index(rest, ":")
	if colon < 0 {
		return nil, fmt.Errorf("prismagql: malformed field %q", line)
	}
	name := strings.TrimSpace(rest[:colon])
	typ := strings.TrimSpace(rest[colon+1:])

	required := strings.HasSuffix(typ, "!")
	typ = strings.TrimSuffix(typ, "!")

	multi := false
	if strings.HasPrefix(typ, "[") && strings.HasSuffix(typ, "]") {
		multi = true
		typ = typ[1 : len(typ)-1]
	}

	prismaType, directive := prismaTypeForGql(typ)

	field := name + " " + prismaType
	if multi {
		field += "[]"
	} else if !required {
		field += "?"
	}
	if directive != "" {
		field += " " + directive
	}
	if argsRaw != "" {
		rendered, err := renderPrismaArgs(argsRaw)
		if err != nil {
			return nil, err
		}
		field += " " + rendered
	}

	out = append(out, field)
	return out, nil
}

// prismaTypeForGql maps a GraphQL type to its Prisma type and optional directive.
func prismaTypeForGql(typ string) (string, string) {
	if typ == "WEI" {
		return "Decimal", "/// parse:SCALAR(WEI)"
	}
	if p, ok := prismaGqlTypeMap[typ]; ok {
		return p, ""
	}
	return typ, ""
}

// renderPrismaArgs turns a GraphQL argument signature into a "parse:ARGUE[...]" directive.
func renderPrismaArgs(argsRaw string) (string, error) {
	inner := strings.TrimSpace(argsRaw)
	inner = strings.TrimPrefix(inner, "(")
	inner = strings.TrimSuffix(inner, ")")

	parts := make([]string, 0)
	for _, seg := range strings.Split(inner, ",") {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		kv := strings.SplitN(seg, ":", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("prismagql: malformed argument %q", seg)
		}
		name := strings.TrimSpace(kv[0])
		typ := strings.TrimSpace(kv[1])
		mandatory := strings.HasSuffix(typ, "!")
		typ = strings.TrimSuffix(typ, "!")
		prisma := prismaArgTypeMap[typ]
		if prisma == "" {
			prisma = typ
		}
		if mandatory {
			name += "!"
		}
		parts = append(parts, name+":"+prisma)
	}

	return "/// parse:ARGUE[" + strings.Join(parts, ",") + "]", nil
}
