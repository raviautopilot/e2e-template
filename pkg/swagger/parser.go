package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// SwaggerSpec is the parsed, normalized representation of a Swagger 2.0 spec.
type SwaggerSpec struct {
	Title       string
	Description string
	Version     string
	BasePath    string
	Host        string
	Groups      []EndpointGroup          // endpoints grouped by tag
	Definitions map[string]Model         // resolved model definitions
	RawPaths    map[string]PathItem      // raw path items (for advanced usage)
}

// EndpointGroup is a collection of endpoints sharing the same swagger tag.
type EndpointGroup struct {
	Tag       string
	Endpoints []Endpoint
}

// Endpoint represents a single API endpoint (method + path).
type Endpoint struct {
	Method      string // GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
	Path        string // e.g. /api/v1/config/dependencies/{id}
	Summary     string
	Description string
	Tag         string
	OperationID string
	Consumes    []string
	Produces    []string
	Parameters  []Parameter
	Responses   map[string]Response // keyed by status code string
	RequestBody *ModelRef           // resolved request body schema (if POST/PUT/PATCH)
	SuccessCode string              // primary success status code
	SuccessRef  string              // $ref for success response model
}

// Parameter represents a swagger parameter.
type Parameter struct {
	Name        string
	In          string // path, query, body, header
	Required    bool
	Type        string // swagger type (string, integer, boolean)
	GoType      string // resolved Go type
	Description string
	Schema      *ModelRef
}

// Response represents a swagger response definition.
type Response struct {
	StatusCode  string
	Description string
	Schema      *ModelRef
}

// ModelRef is a reference to a model definition.
type ModelRef struct {
	Ref      string // e.g. #/definitions/model.CfgDependency
	Resolved string // e.g. model.CfgDependency
	IsArray  bool
}

// Model represents a Go struct derived from a swagger definition.
type Model struct {
	SwaggerName string  // e.g. model.CfgDependency
	GoName      string  // e.g. CfgDependency
	Fields      []Field
}

// Field is a single field in a Model.
type Field struct {
	JSONName string      // e.g. parent_value_code
	GoName   string      // e.g. ParentValueCode
	GoType   string      // e.g. string, int, bool
	JSONTag  string      // e.g. `json:"parent_value_code,omitempty"`
	Example  interface{} // swagger example value
}

// ── Raw swagger JSON structures ─────────────────────────────────────────────

type rawSwagger struct {
	Swagger     string                        `json:"swagger"`
	Info        rawInfo                       `json:"info"`
	Host        string                        `json:"host"`
	BasePath    string                        `json:"basePath"`
	Paths       map[string]map[string]rawOp   `json:"paths"`
	Definitions map[string]rawDefinition      `json:"definitions"`
}

type rawInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type rawOp struct {
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	OperationID string            `json:"operationId"`
	Tags        []string          `json:"tags"`
	Consumes    []string          `json:"consumes"`
	Produces    []string          `json:"produces"`
	Parameters  []rawParameter    `json:"parameters"`
	Responses   map[string]rawResp `json:"responses"`
}

type rawParameter struct {
	Name        string          `json:"name"`
	In          string          `json:"in"`
	Required    bool            `json:"required"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema"`
}

type rawResp struct {
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema"`
}

type rawDefinition struct {
	Type       string                    `json:"type"`
	Properties map[string]rawProperty    `json:"properties"`
}

type rawProperty struct {
	Type    string      `json:"type"`
	Example interface{} `json:"example"`
}

type rawSchemaRef struct {
	Ref   string       `json:"$ref"`
	Type  string       `json:"type"`
	Items *rawSchemaRef `json:"items"`
}

// PathItem is a raw path item for advanced access.
type PathItem = map[string]rawOp

// ── Public API ──────────────────────────────────────────────────────────────

// ParseFromURL fetches a swagger JSON from a URL and parses it.
// If the URL points to a Swagger UI HTML page, it automatically detects and resolves the JSON spec.
func ParseFromURL(rawURL string) (*SwaggerSpec, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch swagger from %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("swagger URL returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read swagger response: %w", err)
	}

	// Check if the response is HTML rather than JSON (e.g. Swagger UI index.html)
	trimmed := bytes.TrimSpace(data)
	if bytes.HasPrefix(trimmed, []byte("<")) || strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return tryResolveSwaggerJSON(rawURL, data)
	}

	return parseJSON(data)
}

// tryResolveSwaggerJSON attempts to locate the swagger JSON specification
// when given a Swagger UI HTML page URL (e.g. /swagger/index.html).
func tryResolveSwaggerJSON(htmlURL string, htmlData []byte) (*SwaggerSpec, error) {
	parsed, err := url.Parse(htmlURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %s: %w", htmlURL, err)
	}

	var candidateURLs []string

	// 1. Look for swagger-initializer.js in the HTML
	initRegex := regexp.MustCompile(`["']([^"']*swagger-initializer\.js[^"']*)["']`)
	if match := initRegex.FindSubmatch(htmlData); len(match) > 1 {
		initRel, err := url.Parse(string(match[1]))
		if err == nil {
			initURL := parsed.ResolveReference(initRel).String()
			if initResp, err := http.Get(initURL); err == nil && initResp.StatusCode == 200 {
				initData, _ := io.ReadAll(initResp.Body)
				initResp.Body.Close()
				urlRegex := regexp.MustCompile(`url:\s*["']([^"']+)["']`)
				if urlMatch := urlRegex.FindSubmatch(initData); len(urlMatch) > 1 {
					specRel, err := url.Parse(string(urlMatch[1]))
					if err == nil {
						candidateURLs = append(candidateURLs, parsed.ResolveReference(specRel).String())
					}
				}
			}
		}
	}

	// 2. Look for direct spec url in HTML
	urlRegex := regexp.MustCompile(`url:\s*["']([^"']+\.json[^"']*)["']`)
	if match := urlRegex.FindSubmatch(htmlData); len(match) > 1 {
		specRel, err := url.Parse(string(match[1]))
		if err == nil {
			candidateURLs = append(candidateURLs, parsed.ResolveReference(specRel).String())
		}
	}

	// 3. Common relative candidate paths
	for _, rel := range []string{"doc.json", "swagger.json", "openapi.json", "v2/api-docs"} {
		relURL, _ := url.Parse(rel)
		candidateURLs = append(candidateURLs, parsed.ResolveReference(relURL).String())
	}

	// 4. If path ended with index.html or index.htm
	cleanPath := strings.TrimSuffix(parsed.Path, "index.html")
	cleanPath = strings.TrimSuffix(cleanPath, "index.htm")
	cleanPath = strings.TrimSuffix(cleanPath, "/")
	candidateURLs = append(candidateURLs, fmt.Sprintf("%s://%s%s/doc.json", parsed.Scheme, parsed.Host, cleanPath))
	candidateURLs = append(candidateURLs, fmt.Sprintf("%s://%s%s/swagger.json", parsed.Scheme, parsed.Host, cleanPath))
	candidateURLs = append(candidateURLs, fmt.Sprintf("%s://%s/swagger/doc.json", parsed.Scheme, parsed.Host))
	candidateURLs = append(candidateURLs, fmt.Sprintf("%s://%s/swagger/swagger.json", parsed.Scheme, parsed.Host))

	// Deduplicate preserving order
	seen := make(map[string]bool)
	var uniqueCandidates []string
	for _, c := range candidateURLs {
		if c != "" && !seen[c] && c != htmlURL {
			seen[c] = true
			uniqueCandidates = append(uniqueCandidates, c)
		}
	}

	for _, cand := range uniqueCandidates {
		resp, err := http.Get(cand)
		if err != nil || resp.StatusCode != 200 {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			continue
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		trimmed := bytes.TrimSpace(data)
		if bytes.HasPrefix(trimmed, []byte("<")) {
			continue
		}

		spec, err := parseJSON(data)
		if err == nil && spec != nil {
			fmt.Printf("ℹ Auto-detected Swagger UI HTML page at %s\n", htmlURL)
			fmt.Printf("ℹ Resolved Swagger JSON spec at: %s\n", cand)
			return spec, nil
		}
	}

	return nil, fmt.Errorf("URL returned HTML (Swagger UI page) instead of Swagger JSON.\nTo fix: please provide the direct Swagger JSON URL (e.g. %s://%s/swagger/doc.json)", parsed.Scheme, parsed.Host)
}

// ParseFromFile reads a swagger JSON from a local file and parses it.
func ParseFromFile(path string) (*SwaggerSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read swagger file %s: %w", path, err)
	}
	return parseJSON(data)
}

// ParseFromReader reads swagger JSON from any io.Reader.
func ParseFromReader(r io.Reader) (*SwaggerSpec, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read swagger input: %w", err)
	}
	return parseJSON(data)
}

// ── Internal parsing ────────────────────────────────────────────────────────

func parseJSON(data []byte) (*SwaggerSpec, error) {
	trimmed := bytes.TrimSpace(data)
	if bytes.HasPrefix(trimmed, []byte("<")) {
		return nil, fmt.Errorf("received HTML/XML instead of Swagger JSON (starts with '<'); please provide a Swagger JSON endpoint")
	}

	var raw rawSwagger
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse swagger JSON: %w", err)
	}

	spec := &SwaggerSpec{
		Title:       raw.Info.Title,
		Description: raw.Info.Description,
		Version:     raw.Info.Version,
		BasePath:    raw.BasePath,
		Host:        raw.Host,
		Definitions: make(map[string]Model),
		RawPaths:    make(map[string]PathItem),
	}

	// 1. Parse definitions into Models
	for defName, rawDef := range raw.Definitions {
		model := Model{
			SwaggerName: defName,
			GoName:      swaggerNameToGoName(defName),
			Fields:      make([]Field, 0, len(rawDef.Properties)),
		}

		// Collect and sort field names for deterministic output
		fieldNames := make([]string, 0, len(rawDef.Properties))
		for name := range rawDef.Properties {
			fieldNames = append(fieldNames, name)
		}
		sort.Strings(fieldNames)

		for _, jsonName := range fieldNames {
			prop := rawDef.Properties[jsonName]
			goType := swaggerTypeToGo(prop.Type)

			// Fields like "id" that are integers should stay int, not become string
			field := Field{
				JSONName: jsonName,
				GoName:   snakeToPascal(jsonName),
				GoType:   goType,
				JSONTag:  fmt.Sprintf("`json:\"%s,omitempty\"`", jsonName),
				Example:  prop.Example,
			}
			model.Fields = append(model.Fields, field)
		}

		spec.Definitions[defName] = model
	}

	// 2. Parse paths into endpoints, grouped by tag
	tagMap := make(map[string][]Endpoint)

	for path, methods := range raw.Paths {
		spec.RawPaths[path] = methods
		for method, op := range methods {
			tag := "General"
			if len(op.Tags) > 0 {
				tag = op.Tags[0]
			}

			ep := Endpoint{
				Method:      strings.ToUpper(method),
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
				Tag:         tag,
				OperationID: op.OperationID,
				Consumes:    op.Consumes,
				Produces:    op.Produces,
				Parameters:  make([]Parameter, 0, len(op.Parameters)),
				Responses:   make(map[string]Response),
			}

			// Parse parameters
			for _, rp := range op.Parameters {
				param := Parameter{
					Name:        rp.Name,
					In:          rp.In,
					Required:    rp.Required,
					Type:        rp.Type,
					GoType:      swaggerTypeToGo(rp.Type),
					Description: rp.Description,
				}

				// Handle body parameter with schema $ref
				if rp.In == "body" && rp.Schema != nil {
					ref := parseSchemaRef(rp.Schema)
					if ref != nil {
						param.Schema = ref
						ep.RequestBody = ref
					}
				}

				ep.Parameters = append(ep.Parameters, param)
			}

			// Parse responses
			for code, rr := range op.Responses {
				resp := Response{
					StatusCode:  code,
					Description: rr.Description,
				}
				if rr.Schema != nil {
					resp.Schema = parseSchemaRef(rr.Schema)
				}
				ep.Responses[code] = resp

				// Identify primary success code and model
				if isSuccessCode(code) {
					ep.SuccessCode = code
					if resp.Schema != nil {
						ep.SuccessRef = resp.Schema.Resolved
					}
				}
			}

			tagMap[tag] = append(tagMap[tag], ep)
		}
	}

	// 3. Sort and build groups
	tagNames := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tagNames = append(tagNames, tag)
	}
	sort.Strings(tagNames)

	for _, tag := range tagNames {
		eps := tagMap[tag]
		// Sort endpoints within group: by path, then by method order
		sort.Slice(eps, func(i, j int) bool {
			if eps[i].Path == eps[j].Path {
				return methodOrder(eps[i].Method) < methodOrder(eps[j].Method)
			}
			return eps[i].Path < eps[j].Path
		})
		spec.Groups = append(spec.Groups, EndpointGroup{
			Tag:       tag,
			Endpoints: eps,
		})
	}

	return spec, nil
}

func parseSchemaRef(data json.RawMessage) *ModelRef {
	var schema rawSchemaRef
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil
	}

	ref := &ModelRef{}
	if schema.Ref != "" {
		ref.Ref = schema.Ref
		ref.Resolved = resolveRef(schema.Ref)
		return ref
	}
	if schema.Type == "array" && schema.Items != nil {
		ref.IsArray = true
		if schema.Items.Ref != "" {
			ref.Ref = schema.Items.Ref
			ref.Resolved = resolveRef(schema.Items.Ref)
		}
		return ref
	}
	return nil
}

func resolveRef(ref string) string {
	// #/definitions/model.CfgDependency → model.CfgDependency
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

// ── Type conversion helpers ─────────────────────────────────────────────────

func swaggerTypeToGo(sType string) string {
	switch sType {
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "string":
		return "string"
	case "array":
		return "[]interface{}"
	default:
		return "interface{}"
	}
}

func swaggerNameToGoName(name string) string {
	// model.CfgDependency → CfgDependency
	// handler.ErrorResponse → ErrorResponse
	parts := strings.Split(name, ".")
	last := parts[len(parts)-1]
	if len(last) > 0 {
		return strings.ToUpper(last[:1]) + last[1:]
	}
	return last
}

func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}
	return result.String()
}

func isSuccessCode(code string) bool {
	return strings.HasPrefix(code, "2")
}

func methodOrder(method string) int {
	switch method {
	case "GET":
		return 0
	case "POST":
		return 1
	case "PUT":
		return 2
	case "PATCH":
		return 3
	case "DELETE":
		return 4
	case "HEAD":
		return 5
	case "OPTIONS":
		return 6
	default:
		return 99
	}
}

// ── Query helpers ───────────────────────────────────────────────────────────

// HasCRUD checks if an endpoint group has the standard list, create, get, update, delete pattern.
func (g EndpointGroup) HasCRUD() bool {
	methods := make(map[string]bool)
	for _, ep := range g.Endpoints {
		methods[ep.Method] = true
	}
	return methods["GET"] && methods["POST"] && (methods["PUT"] || methods["PATCH"]) && methods["DELETE"]
}

// ListEndpoint returns the GET endpoint that lists resources (no path params).
func (g EndpointGroup) ListEndpoint() *Endpoint {
	for i := range g.Endpoints {
		ep := &g.Endpoints[i]
		if ep.Method == "GET" && !strings.Contains(ep.Path, "{") {
			return ep
		}
	}
	return nil
}

// CreateEndpoint returns the POST endpoint.
func (g EndpointGroup) CreateEndpoint() *Endpoint {
	for i := range g.Endpoints {
		ep := &g.Endpoints[i]
		if ep.Method == "POST" {
			return ep
		}
	}
	return nil
}

// GetByIDEndpoint returns the GET endpoint with a path parameter.
func (g EndpointGroup) GetByIDEndpoint() *Endpoint {
	for i := range g.Endpoints {
		ep := &g.Endpoints[i]
		if ep.Method == "GET" && strings.Contains(ep.Path, "{") {
			return ep
		}
	}
	return nil
}

// UpdateEndpoint returns the PUT or PATCH endpoint.
func (g EndpointGroup) UpdateEndpoint() *Endpoint {
	for i := range g.Endpoints {
		ep := &g.Endpoints[i]
		if (ep.Method == "PUT" || ep.Method == "PATCH") && strings.Contains(ep.Path, "{") {
			return ep
		}
	}
	return nil
}

// DeleteEndpoint returns the DELETE endpoint.
func (g EndpointGroup) DeleteEndpoint() *Endpoint {
	for i := range g.Endpoints {
		ep := &g.Endpoints[i]
		if ep.Method == "DELETE" {
			return ep
		}
	}
	return nil
}

// GetRequestModel returns the swagger model name used in POST/PUT body for this group.
func (g EndpointGroup) GetRequestModel() string {
	if ep := g.CreateEndpoint(); ep != nil && ep.RequestBody != nil {
		return ep.RequestBody.Resolved
	}
	if ep := g.UpdateEndpoint(); ep != nil && ep.RequestBody != nil {
		return ep.RequestBody.Resolved
	}
	return ""
}

