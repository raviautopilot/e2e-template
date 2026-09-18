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

// Generate creates Go API test suites from a SwaggerSpec according to GenerateConfig.
func Generate(cfg GenerateConfig) (*GenerateResult, error) {
	spec := cfg.Spec
	if spec == nil {
		return nil, fmt.Errorf("SwaggerSpec cannot be nil")
	}

	result := &GenerateResult{}

	// Filter groups by tags if specified
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

	if len(groups) == 0 {
		return nil, fmt.Errorf("no matching endpoint groups found for tags: %v", cfg.Tags)
	}

	// 1. Generate top-level shared models package (models/models.go) containing all spec definitions
	if len(spec.Definitions) > 0 {
		modelsDir := filepath.Join(cfg.OutputDir, "models")
		if err := os.MkdirAll(modelsDir, 0755); err == nil {
			allModelsPath := filepath.Join(modelsDir, "models.go")
			allModelsList := make([]Model, 0, len(spec.Definitions))
			for _, m := range spec.Definitions {
				allModelsList = append(allModelsList, m)
			}
			sort.Slice(allModelsList, func(i, j int) bool {
				return allModelsList[i].GoName < allModelsList[j].GoName
			})
			allModelData := struct {
				GenerateConfig
				PackageName string
				Models      []Model
			}{cfg, "models", allModelsList}
			if err := writeTemplate(allModelsPath, modelsTmpl, allModelData); err == nil {
				result.Files = append(result.Files, allModelsPath)
			}
		}
	}

	// 2. For each endpoint group, create a dedicated package directory with models & single-test files
	for _, group := range groups {
		groupDirName := sanitizeFileName(group.Tag)
		groupDir := filepath.Join(cfg.OutputDir, groupDirName)
		if err := os.MkdirAll(groupDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create package directory %s: %w", groupDir, err)
		}

		packageName := fmt.Sprintf("%s_%s", cfg.ServiceName, groupDirName)

		// main_test.go in group package
		mainPath := filepath.Join(groupDir, "main_test.go")
		mainData := struct {
			GenerateConfig
			PackageName string
		}{cfg, packageName}
		if err := writeTemplate(mainPath, mainTestTmpl, mainData); err != nil {
			return nil, fmt.Errorf("failed to generate %s: %w", mainPath, err)
		}
		result.Files = append(result.Files, mainPath)

		// models.go in group package (transitively includes all models used by this group)
		modelsPath := filepath.Join(groupDir, "models.go")
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

		// Test files in group package (one test per file)
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
					PackageName     string
					Group           EndpointGroup
					Endpoint        *Endpoint
					ModelName       string
					RequestModel    string
					ResponseModel   string
					ResponseIsArray bool
					Definitions     map[string]Model
				}{
					GenerateConfig:  cfg,
					PackageName:     packageName,
					Group:           group,
					Endpoint:        ep,
					ModelName:       modelName,
					RequestModel:    ep.RequestModelName(spec.Definitions),
					ResponseModel:   ep.ResponseModelName(spec.Definitions),
					ResponseIsArray: ep.ResponseIsArray(),
					Definitions:     spec.Definitions,
				}
				if err := writeTemplate(listPath, listTestTmpl, data); err != nil {
					return nil, fmt.Errorf("failed to generate %s: %w", listPath, err)
				}
				result.Files = append(result.Files, listPath)
			}

			// 02_create_test.go
			if ep := group.CreateEndpoint(); ep != nil {
				createPath := filepath.Join(groupDir, "02_create_test.go")
				reqModel := ep.RequestModelName(spec.Definitions)
				if reqModel == "" {
					reqModel = modelName
				}
				data := struct {
					GenerateConfig
					PackageName     string
					Group           EndpointGroup
					Endpoint        *Endpoint
					ModelName       string
					RequestModel    string
					ResponseModel   string
					ResponseIsArray bool
					Definitions     map[string]Model
				}{
					GenerateConfig:  cfg,
					PackageName:     packageName,
					Group:           group,
					Endpoint:        ep,
					ModelName:       reqModel,
					RequestModel:    reqModel,
					ResponseModel:   ep.ResponseModelName(spec.Definitions),
					ResponseIsArray: ep.ResponseIsArray(),
					Definitions:     spec.Definitions,
				}
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
					PackageName     string
					Group           EndpointGroup
					Endpoint        *Endpoint
					ModelName       string
					RequestModel    string
					ResponseModel   string
					ResponseIsArray bool
					Definitions     map[string]Model
				}{
					GenerateConfig:  cfg,
					PackageName:     packageName,
					Group:           group,
					Endpoint:        ep,
					ModelName:       modelName,
					RequestModel:    ep.RequestModelName(spec.Definitions),
					ResponseModel:   ep.ResponseModelName(spec.Definitions),
					ResponseIsArray: ep.ResponseIsArray(),
					Definitions:     spec.Definitions,
				}
				if err := writeTemplate(getPath, getTestTmpl, data); err != nil {
					return nil, fmt.Errorf("failed to generate %s: %w", getPath, err)
				}
				result.Files = append(result.Files, getPath)
			}

			// 04_update_test.go
			if ep := group.UpdateEndpoint(); ep != nil {
				updatePath := filepath.Join(groupDir, "04_update_test.go")
				reqModel := ep.RequestModelName(spec.Definitions)
				if reqModel == "" {
					reqModel = modelName
				}
				data := struct {
					GenerateConfig
					PackageName     string
					Group           EndpointGroup
					Endpoint        *Endpoint
					ModelName       string
					RequestModel    string
					ResponseModel   string
					ResponseIsArray bool
					Definitions     map[string]Model
				}{
					GenerateConfig:  cfg,
					PackageName:     packageName,
					Group:           group,
					Endpoint:        ep,
					ModelName:       reqModel,
					RequestModel:    reqModel,
					ResponseModel:   ep.ResponseModelName(spec.Definitions),
					ResponseIsArray: ep.ResponseIsArray(),
					Definitions:     spec.Definitions,
				}
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
					PackageName     string
					Group           EndpointGroup
					Endpoint        *Endpoint
					ModelName       string
					RequestModel    string
					ResponseModel   string
					ResponseIsArray bool
					Definitions     map[string]Model
				}{
					GenerateConfig:  cfg,
					PackageName:     packageName,
					Group:           group,
					Endpoint:        ep,
					ModelName:       modelName,
					RequestModel:    ep.RequestModelName(spec.Definitions),
					ResponseModel:   ep.ResponseModelName(spec.Definitions),
					ResponseIsArray: ep.ResponseIsArray(),
					Definitions:     spec.Definitions,
				}
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
		return fmt.Errorf("failed to parse template: %w", err)
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
	var queue []string

	for _, g := range groups {
		for _, ep := range g.Endpoints {
			if ep.RequestBody != nil && ep.RequestBody.Resolved != "" {
				queue = append(queue, ep.RequestBody.Resolved)
			}
			for _, resp := range ep.Responses {
				if resp.Schema != nil && resp.Schema.Resolved != "" {
					queue = append(queue, resp.Schema.Resolved)
				}
			}
			for _, p := range ep.Parameters {
				if p.Schema != nil && p.Schema.Resolved != "" {
					queue = append(queue, p.Schema.Resolved)
				}
			}
		}
	}

	for defName := range defs {
		if strings.HasSuffix(strings.ToLower(defName), "errorresponse") {
			queue = append(queue, defName)
		}
	}

	var models []Model
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]

		if seen[name] {
			continue
		}
		seen[name] = true

		if m, ok := findModel(name, defs); ok {
			models = append(models, m)
			for _, f := range m.Fields {
				if f.RefModel != "" && !seen[f.RefModel] {
					queue = append(queue, f.RefModel)
				}
			}
		}
	}

	// If no models were explicitly referenced by group endpoints, include all definitions
	if len(models) == 0 && len(defs) > 0 {
		for _, m := range defs {
			models = append(models, m)
		}
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i].GoName < models[j].GoName
	})

	return models
}

func findModel(name string, defs map[string]Model) (Model, bool) {
	if m, ok := defs[name]; ok {
		return m, true
	}
	goName := swaggerNameToGoName(name)
	for _, m := range defs {
		if m.GoName == goName || m.SwaggerName == name {
			return m, true
		}
	}
	return Model{}, false
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
	if len(g.Endpoints) != 1 {
		return false
	}
	ep := g.Endpoints[0]
	path := strings.ToLower(ep.Path)
	return strings.Contains(path, "health") || strings.Contains(path, "ping") || strings.Contains(path, "status")
}

// buildExampleBody creates a Go struct literal or map literal string from a model's fields.
func buildExampleBody(modelName string, defs map[string]Model) string {
	model, ok := findModel(modelName, defs)
	if !ok || len(model.Fields) == 0 {
		return "map[string]interface{}{}"
	}

	var lines []string
	for _, f := range model.Fields {
		if isAutoField(f.JSONName) {
			continue
		}
		val := exampleValue(f)
		lines = append(lines, fmt.Sprintf("\t\t\t\t%s: %s,", f.GoName, val))
	}

	if len(lines) == 0 {
		return model.GoName + "{}"
	}

	return model.GoName + "{\n" + strings.Join(lines, "\n") + "\n\t\t\t}"
}

// buildUpdateBody creates a Go struct literal or map literal string for updating.
func buildUpdateBody(modelName string, defs map[string]Model) string {
	model, ok := findModel(modelName, defs)
	if !ok || len(model.Fields) == 0 {
		return "map[string]interface{}{\"updated\": true}"
	}

	var lines []string
	for _, f := range model.Fields {
		if isAutoField(f.JSONName) {
			continue
		}
		val := exampleValueUpdated(f)
		lines = append(lines, fmt.Sprintf("\t\t\t\t%s: %s,", f.GoName, val))
	}

	if len(lines) == 0 {
		return model.GoName + "{}"
	}

	return model.GoName + "{\n" + strings.Join(lines, "\n") + "\n\t\t\t}"
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
			if f.GoType == "int" || f.GoType == "int64" {
				return fmt.Sprintf("%d", int64(v))
			}
			return fmt.Sprintf("%v", v)
		case bool:
			return fmt.Sprintf("%v", v)
		}
	}

	switch f.GoType {
	case "string":
		return fmt.Sprintf("%q", "sample_"+f.JSONName)
	case "int", "int64":
		return "1"
	case "float64":
		return "1.0"
	case "bool":
		return "true"
	case "[]string":
		return fmt.Sprintf("[]string{%q}", "sample_"+f.JSONName)
	case "[]int":
		return "[]int{1}"
	case "map[string]interface{}":
		return "map[string]interface{}{}"
	default:
		if strings.HasPrefix(f.GoType, "[]") {
			return "nil"
		}
		if strings.HasPrefix(f.GoType, "*") {
			return "nil"
		}
		return "nil"
	}
}

func exampleValueUpdated(f Field) string {
	if f.Example != nil {
		switch v := f.Example.(type) {
		case string:
			return fmt.Sprintf("%q", v+"_updated")
		case float64:
			if f.GoType == "int" || f.GoType == "int64" {
				return fmt.Sprintf("%d", int64(v)+1)
			}
			return fmt.Sprintf("%v", v)
		case bool:
			return fmt.Sprintf("%v", !v)
		}
	}

	switch f.GoType {
	case "string":
		return fmt.Sprintf("%q", "updated_"+f.JSONName)
	case "int", "int64":
		return "2"
	case "float64":
		return "2.0"
	case "bool":
		return "false"
	case "[]string":
		return fmt.Sprintf("[]string{%q}", "updated_"+f.JSONName)
	case "[]int":
		return "[]int{2}"
	case "map[string]interface{}":
		return "map[string]interface{}{}"
	default:
		if strings.HasPrefix(f.GoType, "[]") {
			return "nil"
		}
		if strings.HasPrefix(f.GoType, "*") {
			return "nil"
		}
		return "nil"
	}
}
