package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"e2e-template/pkg/swagger"
)

func main() {
	// CLI flags
	swaggerURL := flag.String("swagger-url", "", "URL to fetch swagger JSON from (e.g. http://localhost:1705/swagger/doc.json)")
	swaggerFile := flag.String("swagger-file", "", "Local path to swagger JSON file")
	serviceName := flag.String("service-name", "", "Service/package name for generated tests (e.g. configsvc)")
	baseURL := flag.String("base-url", "", "Base URL for API calls in generated tests (e.g. http://localhost:1705)")
	outputDir := flag.String("output-dir", "", "Output directory for generated test files (e.g. tests/api/configsvc)")
	modulePath := flag.String("module-path", "e2e-template", "Go module path from go.mod")
	tags := flag.String("tags", "", "Comma-separated list of swagger tags to generate (empty = all)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `API Test Generator — Generate Go API tests from Swagger 2.0 spec

Usage:
  generate-tests [flags]

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  # Generate from URL, all endpoints
  generate-tests --swagger-url http://localhost:1705/swagger/doc.json \
    --service-name configsvc --base-url http://localhost:1705 \
    --output-dir tests/api/configsvc

  # Generate from file, specific tags only
  generate-tests --swagger-file swagger.json \
    --service-name usersvc --base-url http://localhost:8080 \
    --output-dir tests/api/usersvc --tags "Users,Auth"
`)
	}

	flag.Parse()

	// Validate inputs
	if *swaggerURL == "" && *swaggerFile == "" {
		fmt.Fprintf(os.Stderr, "ERROR: Either --swagger-url or --swagger-file is required.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *serviceName == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --service-name is required.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *baseURL == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --base-url is required.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *outputDir == "" {
		*outputDir = "tests/api/" + *serviceName
	}

	// Parse swagger spec
	fmt.Println("=========================================")
	fmt.Println("  API Test Generator from Swagger")
	fmt.Println("=========================================")
	fmt.Println()

	var spec *swagger.SwaggerSpec
	var err error

	if *swaggerURL != "" {
		fmt.Printf("📡 Fetching swagger from: %s\n", *swaggerURL)
		spec, err = swagger.ParseFromURL(*swaggerURL)
	} else {
		fmt.Printf("📄 Reading swagger from: %s\n", *swaggerFile)
		spec, err = swagger.ParseFromFile(*swaggerFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse swagger: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Parsed: %s (v%s)\n", spec.Title, spec.Version)
	fmt.Printf("  Base Path: %s\n", spec.BasePath)
	fmt.Printf("  Definitions: %d models\n", len(spec.Definitions))
	fmt.Printf("  Tag Groups: %d\n", len(spec.Groups))
	for _, g := range spec.Groups {
		fmt.Printf("    • %s (%d endpoints)\n", g.Tag, len(g.Endpoints))
	}
	fmt.Println()

	// Parse tags filter
	var tagFilter []string
	if *tags != "" {
		for _, t := range strings.Split(*tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tagFilter = append(tagFilter, t)
			}
		}
		fmt.Printf("🔍 Filtering to tags: %v\n", tagFilter)
	} else {
		fmt.Println("🔍 Generating tests for ALL endpoints")
	}

	// Generate
	fmt.Printf("📁 Output directory: %s\n", *outputDir)
	fmt.Printf("📦 Package name: %s_test\n", *serviceName)
	fmt.Println()

	cfg := swagger.GenerateConfig{
		Spec:        spec,
		ServiceName: *serviceName,
		BaseURL:     *baseURL,
		OutputDir:   *outputDir,
		Tags:        tagFilter,
		ModulePath:  *modulePath,
	}

	result, err := swagger.Generate(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to generate tests: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=========================================")
	fmt.Println("  ✓ Generation Complete!")
	fmt.Println("=========================================")
	for _, f := range result.Files {
		fmt.Printf("  ✓ Created %s\n", f)
	}
	fmt.Println()
	fmt.Printf("Run with: ./run-api-tests.sh ./%s/...\n", *outputDir)
}
