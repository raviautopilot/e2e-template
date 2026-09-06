package swagger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// GenerateConfig holds all options for test generation.
type GenerateConfig struct {
	Spec        *SwaggerSpec
	ServiceName string   // e.g. "configsvc" — used as Go package name & directory
	BaseURL     string   // e.g. "http://localhost:1705"
	OutputDir   string   // e.g. "tests/api/configsvc"
	Tags        []string // empty = all tags; otherwise filter
	ModulePath  string   // Go module path, e.g. "e2e-template"
}

// GenerateResult summarizes what was generated.
type GenerateResult struct {
	Files []string // list of created file paths
}

// Generate writes Go test files from a parsed swagger spec into per-group packages.
func Generate(cfg GenerateConfig) (*GenerateResult, error) {
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory %s: %w", cfg.OutputDir, err)
	}

	result := &GenerateResult{}
	spec := cfg.Spec

	// Filter groups by tag if specified
	groups := spec.Groups
	if len(cfg.Tags) > 0 {
		tagSet := make(map[string]bool)
		for _, t := range cfg.Tags {
			tagSet[strings.ToLower(strings.TrimSpace(t))] = true
		}
		var filtered []EndpointGroup
		for _, g := range groups {
			if tagSet[strings.ToLower(g.Tag)] {
				filtered = append(filtered, g)
			}
		}
		groups = filtered
	}

	// For each endpoint group, create a dedicated package directory with single-test files
	for _, group := range groups {
		groupDirName := sanitizeFileName(group.Tag)
		groupDir := filepath.Join(cfg.OutputDir, groupDirName)
		if err := os.MkdirAll(groupDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create package directory %s: %w", groupDir, err)
		}

		packageName := fmt.Sprintf("%s_%s", cfg.ServiceName, groupDirName)

		// 1. main_test.go in group package
		mainPath := filepath.Join(groupDir, "main_test.go")
		mainData := struct {
			GenerateConfig
			PackageName string
		}{cfg, packageName}
		if err := writeTemplate(mainPath, mainTestTmpl, mainData); err != nil {
			return nil, fmt.Errorf("failed to generate %s: %w", mainPath, err)
		}
		result.Files = append(result.Files, mainPath)

		// 2. models_test.go in group package
		modelsPath := filepath.Join(groupDir, "models_test.go")
		referencedModels := collectReferencedModels([]EndpointGroup{group}, spec.Definitions)
		modelData := struct {
			GenerateConfig
			PackageName string
			Models      []Model
		}{cfg, packageName, referencedModels}
		if err := writeTemplate(modelsPath, modelsTmpl, modelData); err != nil {
			return nil, fmt.Errorf("failed to generate %s: %w", modelsPath, err)
		}
		result.Files = append(result.Files, modelsPath)

		// 3. Test files in group package (one test per file)
		if strings.ToLower(group.Tag) == "system" || isHealthOnlyGroup(group) {
			healthPath := filepath.Join(groupDir, "01_health_test.go")
			data := struct {
				GenerateConfig
				PackageName string
				Group       EndpointGroup
				Definitions map[string]Model
			}{cfg, packageName, group, spec.Definitions}
			if err := writeTemplate(healthPath, healthTestTmpl, data); err != nil {
				return nil, fmt.Errorf("failed to generate %s: %w", healthPath, err)
			}
			result.Files = append(result.Files, healthPath)
		} else {
			modelName := group.GetRequestModel()

			// 01_list_test.go
			if ep := group.ListEndpoint(); ep != nil {
				listPath := filepath.Join(groupDir, "01_list_test.go")
				data := struct {
					GenerateConfig
					PackageName string
					Group       EndpointGroup
					Endpoint    *Endpoint
					ModelName   string
					Definitions map[string]Model
				}{cfg, packageName, group, ep, modelName, spec.Definitions}
				if err := writeTemplate(listPath, listTestTmpl, data); err != nil {
					return nil, fmt.Errorf("failed to generate %s: %w", listPath, err)
				}
				result.Files = append(result.Files, listPath)
			}

			// 02_create_test.go
			if ep := group.CreateEndpoint(); ep != nil {
				createPath := filepath.Join(groupDir, "02_create_test.go")
				data := struct {
					GenerateConfig
					PackageName string
					Group       EndpointGroup
					Endpoint    *Endpoint
					ModelName   string
					Definitions map[string]Model
				}{cfg, packageName, group, ep, modelName, spec.Definitions}
				if err := writeTemplate(createPath, createTestTmpl, data); err != nil {
					return nil, fmt.Errorf("failed to generate %s: %w", createPath, err)
				}
				result.Files = append(result.Files, createPath)
			}

			// 03_get_test.go
			if ep := group.GetByIDEndpoint(); ep != nil {
				getPath := filepath.Join(groupDir, "03_get_test.go")
				data := struct {
					GenerateConfig
					PackageName string
					Group       EndpointGroup
					Endpoint    *Endpoint
					ModelName   string
					Definitions map[string]Model
				}{cfg, packageName, group, ep, modelName, spec.Definitions}
				if err := writeTemplate(getPath, getTestTmpl, data); err != nil {
					return nil, fmt.Errorf("failed to generate %s: %w", getPath, err)
				}
				result.Files = append(result.Files, getPath)
			}

			// 04_update_test.go
			if ep := group.UpdateEndpoint(); ep != nil {
				updatePath := filepath.Join(groupDir, "04_update_test.go")
				data := struct {
					GenerateConfig
					PackageName string
					Group       EndpointGroup
					Endpoint    *Endpoint
					ModelName   string
					Definitions map[string]Model
				}{cfg, packageName, group, ep, modelName, spec.Definitions}
				if err := writeTemplate(updatePath, updateTestTmpl, data); err != nil {
					return nil, fmt.Errorf("failed to generate %s: %w", updatePath, err)
				}
				result.Files = append(result.Files, updatePath)
			}

			// 05_delete_test.go
			if ep := group.DeleteEndpoint(); ep != nil {
				deletePath := filepath.Join(groupDir, "05_delete_test.go")
				data := struct {
					GenerateConfig
					PackageName string
					Group       EndpointGroup
					Endpoint    *Endpoint
					ModelName   string
					Definitions map[string]Model
				}{cfg, packageName, group, ep, modelName, spec.Definitions}
				if err := writeTemplate(deletePath, deleteTestTmpl, data); err != nil {
					return nil, fmt.Errorf("failed to generate %s: %w", deletePath, err)
				}
				result.Files = append(result.Files, deletePath)
			}

			// 06_parameterized_test.go (table-driven test)
			paramPath := filepath.Join(groupDir, "06_parameterized_test.go")
			data := struct {
				GenerateConfig
				PackageName string
				Group       EndpointGroup
				ModelName   string
				Definitions map[string]Model
			}{cfg, packageName, group, modelName, spec.Definitions}
			if err := writeTemplate(paramPath, parameterizedTestTmpl, data); err != nil {
				return nil, fmt.Errorf("failed to generate %s: %w", paramPath, err)
			}
			result.Files = append(result.Files, paramPath)
		}
	}

	return result, nil
}

func writeTemplate(path string, tmplStr string, data interface{}) error {
	funcMap := template.FuncMap{
		"toLower":          strings.ToLower,
		"toUpper":          strings.ToUpper,
		"toPascal":         snakeToPascal,
		"sanitize":         sanitizeTestName,
		"sanitizeFileName": sanitizeFileName,
		"hasPathParam":     hasPathParam,
		"pathParamName":    pathParamName,
		"swaggerToGoName":  swaggerNameToGoName,
		"buildExampleBody": buildExampleBody,
		"buildUpdateBody":  buildUpdateBody,
		"pad":              func(n int) string { return fmt.Sprintf("%02d", n) },
		"add":              func(a, b int) int { return a + b },
		"quote":            func(s string) string { return fmt.Sprintf("%q", s) },
		"join":             strings.Join,
		"pathParamReplace": pathParamReplace,
	}

	tmpl, err := template.New("gen").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func collectReferencedModels(groups []EndpointGroup, defs map[string]Model) []Model {
	seen := make(map[string]bool)
	var models []Model

	for _, g := range groups {
		for _, ep := range g.Endpoints {
			// Collect from request body
			if ep.RequestBody != nil && ep.RequestBody.Resolved != "" {
				if !seen[ep.RequestBody.Resolved] {
					if m, ok := defs[ep.RequestBody.Resolved]; ok {
						models = append(models, m)
						seen[ep.RequestBody.Resolved] = true
					}
				}
			}
			// Collect from responses
			for _, resp := range ep.Responses {
				if resp.Schema != nil && resp.Schema.Resolved != "" {
					if !seen[resp.Schema.Resolved] {
						if m, ok := defs[resp.Schema.Resolved]; ok {
							models = append(models, m)
							seen[resp.Schema.Resolved] = true
						}
					}
				}
			}
		}
	}

	// Also add ErrorResponse if present
	if !seen["handler.ErrorResponse"] {
		if m, ok := defs["handler.ErrorResponse"]; ok {
			models = append(models, m)
		}
	}

	// Sort for deterministic output
	sort.Slice(models, func(i, j int) bool {
		return models[i].GoName < models[j].GoName
	})

	return models
}

func sanitizeTestName(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "-", "_")
	var result strings.Builder
	for _, r := range s {
		if r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func sanitizeFileName(s string) string {
	return strings.ToLower(sanitizeTestName(s))
}

func hasPathParam(path string) bool {
	return strings.Contains(path, "{")
}

func pathParamName(path string) string {
	start := strings.Index(path, "{")
	end := strings.Index(path, "}")
	if start >= 0 && end > start {
		return path[start+1 : end]
	}
	return "id"
}

// pathParamReplace converts /api/v1/resource/{id} to /api/v1/resource/%d for fmt.Sprintf.
func pathParamReplace(path string) string {
	result := path
	for {
		start := strings.Index(result, "{")
		if start < 0 {
			break
		}
		end := strings.Index(result, "}")
		if end < 0 {
			break
		}
		result = result[:start] + "%d" + result[end+1:]
	}
	return result
}

func isHealthOnlyGroup(g EndpointGroup) bool {
	for _, ep := range g.Endpoints {
		if strings.Contains(strings.ToLower(ep.Path), "health") {
			return true
		}
	}
	return len(g.Endpoints) == 1 && g.Endpoints[0].Method == "GET"
}

// buildExampleBody creates a Go map literal string from a model's example values.
func buildExampleBody(modelName string, defs map[string]Model) string {
	model, ok := defs[modelName]
	if !ok {
		return "map[string]interface{}{}"
	}

	var lines []string
	for _, f := range model.Fields {
		// Skip auto-generated fields
		if isAutoField(f.JSONName) {
			continue
		}
		val := exampleValue(f)
		lines = append(lines, fmt.Sprintf("\t\t%q: %s,", f.JSONName, val))
	}

	if len(lines) == 0 {
		return "map[string]interface{}{}"
	}

	return "map[string]interface{}{\n" + strings.Join(lines, "\n") + "\n\t}"
}

// buildUpdateBody creates a Go map literal for updating (modifies string values).
func buildUpdateBody(modelName string, defs map[string]Model) string {
	model, ok := defs[modelName]
	if !ok {
		return "map[string]interface{}{}"
	}

	var lines []string
	for _, f := range model.Fields {
		if isAutoField(f.JSONName) {
			continue
		}
		val := exampleValueUpdated(f)
		lines = append(lines, fmt.Sprintf("\t\t%q: %s,", f.JSONName, val))
	}

	if len(lines) == 0 {
		return "map[string]interface{}{}"
	}

	return "map[string]interface{}{\n" + strings.Join(lines, "\n") + "\n\t}"
}

func isAutoField(name string) bool {
	auto := map[string]bool{
		"id": true, "created_at": true, "created_by": true,
		"updated_at": true, "updated_by": true, "deleted_at": true,
	}
	return auto[name]
}

func exampleValue(f Field) string {
	if f.Example != nil {
		switch v := f.Example.(type) {
		case string:
			return fmt.Sprintf("%q", v)
		case float64:
			return fmt.Sprintf("%v", int(v))
		case bool:
			return fmt.Sprintf("%v", v)
		default:
			return fmt.Sprintf("%q", fmt.Sprintf("%v", v))
		}
	}
	// Generate unique sensible defaults by type
	switch f.GoType {
	case "string":
		return fmt.Sprintf("fmt.Sprintf(\"%s-%%d\", time.Now().UnixNano()%%100000)", f.JSONName)
	case "int":
		return "1"
	case "bool":
		return "true"
	default:
		return "nil"
	}
}

func exampleValueUpdated(f Field) string {
	if f.Example != nil {
		switch v := f.Example.(type) {
		case string:
			return fmt.Sprintf("%q", v+"-updated")
		case float64:
			return fmt.Sprintf("%v", int(v)+1)
		case bool:
			return fmt.Sprintf("%v", !v)
		default:
			return fmt.Sprintf("%q", fmt.Sprintf("%v-updated", v))
		}
	}
	switch f.GoType {
	case "string":
		return fmt.Sprintf("fmt.Sprintf(\"upd-%s-%%d\", time.Now().UnixNano()%%100000)", f.JSONName)
	case "int":
		return "2"
	case "bool":
		return "false"
	default:
		return "nil"
	}
}
