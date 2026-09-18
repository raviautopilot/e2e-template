package swagger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSwaggerDefinitions(t *testing.T) {
	docPath := filepath.Join("..", "..", "local", "doc.json")
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		t.Skip("local/doc.json not found, skipping local swagger test")
	}

	spec, err := ParseFromFile(docPath)
	if err != nil {
		t.Fatalf("ParseFromFile failed: %v", err)
	}

	if len(spec.Definitions) == 0 {
		t.Fatalf("expected parsed definitions, got 0")
	}

	// Verify AnnouncementRequest
	annReq, ok := spec.Definitions["handler.AnnouncementRequest"]
	if !ok {
		t.Fatalf("expected handler.AnnouncementRequest definition")
	}
	if annReq.GoName != "AnnouncementRequest" {
		t.Errorf("expected GoName AnnouncementRequest, got %s", annReq.GoName)
	}
	if len(annReq.Fields) == 0 {
		t.Errorf("expected fields in AnnouncementRequest, got 0")
	}

	// Verify disambiguation: handler.DistrictBearer vs office_bearers.DistrictBearer
	handlerBearer, ok := spec.Definitions["handler.DistrictBearer"]
	if !ok {
		t.Fatalf("expected handler.DistrictBearer definition")
	}
	officeBearer, ok := spec.Definitions["office_bearers.DistrictBearer"]
	if !ok {
		t.Fatalf("expected office_bearers.DistrictBearer definition")
	}
	if handlerBearer.GoName == officeBearer.GoName {
		t.Errorf("expected distinct GoNames for colliding definitions, got %s for both", handlerBearer.GoName)
	}
	if handlerBearer.GoName != "HandlerDistrictBearer" {
		t.Errorf("expected HandlerDistrictBearer, got %s", handlerBearer.GoName)
	}
	if officeBearer.GoName != "OfficeBearersDistrictBearer" {
		t.Errorf("expected OfficeBearersDistrictBearer, got %s", officeBearer.GoName)
	}
}

func TestGenerateTestsWithModels(t *testing.T) {
	docPath := filepath.Join("..", "..", "local", "doc.json")
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		t.Skip("local/doc.json not found, skipping local generator test")
	}

	spec, err := ParseFromFile(docPath)
	if err != nil {
		t.Fatalf("ParseFromFile failed: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "swagger-gen-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := GenerateConfig{
		Spec:        spec,
		ServiceName: "tagasvc",
		BaseURL:     "http://localhost:8080",
		OutputDir:   tmpDir,
		ModulePath:  "e2e-template",
		Tags:        []string{"Admin Announcements"},
	}

	res, err := Generate(cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(res.Files) == 0 {
		t.Fatalf("expected generated files, got 0")
	}

	// Verify models.go was generated
	modelsPath := filepath.Join(tmpDir, "admin_announcements", "models.go")
	content, err := os.ReadFile(modelsPath)
	if err != nil {
		t.Fatalf("failed to read models.go: %v", err)
	}
	strContent := string(content)
	if !strings.Contains(strContent, "type AnnouncementRequest struct") {
		t.Errorf("models.go missing AnnouncementRequest struct")
	}
	if !strings.Contains(strContent, "type AnnouncementResponse struct") {
		t.Errorf("models.go missing AnnouncementResponse struct")
	}

	// Verify top-level models/models.go was generated
	topModelsPath := filepath.Join(tmpDir, "models", "models.go")
	topContent, err := os.ReadFile(topModelsPath)
	if err != nil {
		t.Fatalf("failed to read top-level models.go: %v", err)
	}
	if !strings.Contains(string(topContent), "type AnnouncementRequest struct") {
		t.Errorf("top-level models.go missing AnnouncementRequest struct")
	}
}
