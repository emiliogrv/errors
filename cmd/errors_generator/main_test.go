package main

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewGenerator tests the New constructor.
func TestNewGenerator(t *testing.T) {
	t.Parallel()

	// given: nothing
	// when: creating a new generator
	gen := New()

	// then: generator should be properly initialized
	assert.NotNil(t, gen)
	assert.NotNil(t, gen.templates)
	assert.Equal(t, Version, gen.data.Version)
	assert.True(t, gen.data.WithGenHeader)
	assert.Equal(t, TestGenNone, gen.TestGenLevel)
	assert.NotEmpty(t, gen.data.Date)
	assert.Equal(t, []string{"attr", "common", "error", "join", "json", "map", "string", "wrap"}, gen.Formats)
}

// TestValidateTestGenLevel tests the validateTestGenLevel method.
func TestValidateTestGenLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		level         string
		expectedLevel string
		expectError   bool
	}{
		{
			name:          "valid_none_level",
			level:         TestGenNone,
			expectedLevel: TestGenNone,
			expectError:   false,
		},
		{
			name:          "valid_flex_level",
			level:         TestGenFlex,
			expectedLevel: TestGenFlex,
			expectError:   false,
		},
		{
			name:          "valid_strict_level",
			level:         TestGenStrict,
			expectedLevel: TestGenStrict,
			expectError:   false,
		},
		{
			name:        "invalid_level",
			level:       "invalid",
			expectError: true,
		},
		{
			name:        "empty_level",
			level:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator
				gen := New()

				// when: validating test generation level
				err := gen.validateTestGenLevel(test.level)

				// then: error should match expectation
				if test.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, test.expectedLevel, gen.TestGenLevel)
				}
			},
		)
	}
}

// TestLoadFormats tests the loadFormats method.
func TestLoadFormats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		formats         string
		initialFormats  []string
		expectedFormats []string
		expectNil       bool
	}{
		{
			name:            "empty_formats_keeps_defaults",
			formats:         "",
			initialFormats:  []string{"attr", "common"},
			expectedFormats: []string{"attr", "common"},
			expectNil:       false,
		},
		{
			name:           "all_formats_clears_defaults",
			formats:        "all",
			initialFormats: []string{"attr", "common"},
			expectNil:      true,
		},
		{
			name:            "single_format_appends",
			formats:         "custom",
			initialFormats:  []string{"attr", "common"},
			expectedFormats: []string{"attr", "common", "custom"},
			expectNil:       false,
		},
		{
			name:            "multiple_formats_appends",
			formats:         "custom1,custom2,custom3",
			initialFormats:  []string{"attr"},
			expectedFormats: []string{"attr", "custom1", "custom2", "custom3"},
			expectNil:       false,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator with initial formats
				gen := New()
				gen.Formats = test.initialFormats

				// when: loading formats
				gen.loadFormats(test.formats)

				// then: formats should match expectation
				if test.expectNil {
					assert.Nil(t, gen.Formats)
				} else {
					assert.Equal(t, test.expectedFormats, gen.Formats)
				}
			},
		)
	}
}

// TestDiscoverTemplateFormats tests the discoverTemplateFormats method.
func TestDiscoverTemplateFormats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		templateNames     []string
		expectedFormats   []string
		minExpectedLength int
	}{
		{
			name:              "no_templates",
			templateNames:     []string{},
			expectedFormats:   []string{},
			minExpectedLength: 0,
		},
		{
			name:              "single_template",
			templateNames:     []string{"error.tmpl"},
			expectedFormats:   []string{"error"},
			minExpectedLength: 1,
		},
		{
			name:              "multiple_templates",
			templateNames:     []string{"error.tmpl", "wrap.tmpl", "join.tmpl"},
			minExpectedLength: 3,
		},
		{
			name:              "excludes_test_templates",
			templateNames:     []string{"error.tmpl", "error_test.tmpl", "wrap.tmpl"},
			minExpectedLength: 2,
		},
		{
			name:              "excludes_non_tmpl_files",
			templateNames:     []string{"error.tmpl", "readme.md", "wrap.tmpl"},
			minExpectedLength: 2,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator with templates
				gen := New()
				for _, name := range test.templateNames {
					gen.templates[name] = template.New(name)
				}

				// when: discovering template formats
				formats := gen.discoverTemplateFormats()

				// then: formats should match expectations
				assert.Len(t, formats, test.minExpectedLength)

				if len(test.expectedFormats) > 0 {
					assert.ElementsMatch(t, test.expectedFormats, formats)
				}
			},
		)
	}
}

// TestHasTemplate tests the hasTemplate method.
func TestHasTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		templates    map[string]*template.Template
		name         string
		templateName string
		expected     bool
	}{
		{
			name:         "template_exists",
			templateName: "error.tmpl",
			templates: map[string]*template.Template{
				"error.tmpl": template.New("error.tmpl"),
			},
			expected: true,
		},
		{
			name:         "template_does_not_exist",
			templateName: "missing.tmpl",
			templates: map[string]*template.Template{
				"error.tmpl": template.New("error.tmpl"),
			},
			expected: false,
		},
		{
			name:         "empty_templates",
			templateName: "error.tmpl",
			templates:    map[string]*template.Template{},
			expected:     false,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator with templates
				gen := New()
				gen.templates = test.templates

				// when: checking if template exists
				result := gen.hasTemplate(test.templateName)

				// then: result should match expectation
				assert.Equal(t, test.expected, result)
			},
		)
	}
}

// TestLoadEmbeddedTemplates tests the loadEmbeddedTemplates method.
func TestLoadEmbeddedTemplates(t *testing.T) {
	t.Parallel()

	// given: a new generator
	gen := New()

	// when: loading embedded templates
	err := gen.loadEmbeddedTemplates()

	// then: templates should be loaded without error
	require.NoError(t, err)
	assert.NotEmpty(t, gen.templates)
}

// TestLoadUserTemplates tests the loadUserTemplates method.
func TestLoadUserTemplates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setupDir    func(*testing.T) string
		name        string
		expectError bool
	}{
		{
			name: "valid_user_templates",
			setupDir: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				tmplPath := filepath.Join(dir, "custom.tmpl")
				err := os.WriteFile(tmplPath, []byte("package {{.PackageName}}"), 0o600)
				require.NoError(t, err)

				return dir
			},
			expectError: false,
		},
		{
			name: "empty_directory",
			setupDir: func(t *testing.T) string {
				t.Helper()

				return t.TempDir()
			},
			expectError: false,
		},
		{
			name: "nonexistent_directory",
			setupDir: func(t *testing.T) string {
				t.Helper()

				return filepath.Join(t.TempDir(), "nonexistent")
			},
			expectError: true,
		},
		{
			name: "invalid_template_syntax",
			setupDir: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				tmplPath := filepath.Join(dir, "invalid.tmpl")
				err := os.WriteFile(tmplPath, []byte("{{.Invalid"), 0o600)
				require.NoError(t, err)

				return dir
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator and a user template directory
				gen := New()
				dir := test.setupDir(t)

				// when: loading user templates
				err := gen.loadUserTemplates(dir)

				// then: error should match expectation
				if test.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

// TestGenerateFile tests the generateFile method.
func TestGenerateFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setupGen     func(*Generator)
		validateFile func(*testing.T, string)
		name         string
		templateName string
		outputName   string
		expectError  bool
	}{
		{
			name:         "successful_generation",
			templateName: "test.tmpl",
			outputName:   "test.go",
			setupGen: func(gen *Generator) {
				tmpl := template.Must(template.New("test.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["test.tmpl"] = tmpl
				gen.data.PackageName = "testpkg"
			},
			expectError: false,
			validateFile: func(t *testing.T, path string) {
				t.Helper()

				content, err := os.ReadFile(path) //nolint:gosec // security is not a concern here
				require.NoError(t, err)
				assert.Equal(t, "package testpkg", string(content))
			},
		},
		{
			name:         "template_not_found",
			templateName: "missing.tmpl",
			outputName:   "test.go",
			setupGen:     func(*Generator) {},
			expectError:  true,
		},
		{
			name:         "template_execution_error",
			templateName: "error.tmpl",
			outputName:   "test.go",
			setupGen: func(gen *Generator) {
				tmpl := template.Must(template.New("error.tmpl").Parse("{{.NonExistent}}"))
				gen.templates["error.tmpl"] = tmpl
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator with templates and target directory
				gen := New()
				gen.OutputDir = t.TempDir()
				test.setupGen(gen)

				// when: generating a file
				err := gen.generateFile(test.templateName, test.outputName)

				// then: error should match expectation
				if test.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					if test.validateFile != nil {
						outputPath := filepath.Join(gen.OutputDir, test.outputName)
						test.validateFile(t, outputPath)
					}
				}
			},
		)
	}
}

// TestGenerateFormat tests the generateFormat method.
func TestGenerateFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setupGen     func(*Generator)
		name         string
		format       string
		testGenLevel string
		expectFiles  []string
		expectError  bool
	}{
		{
			name:         "none_level_no_test",
			format:       "error",
			testGenLevel: TestGenNone,
			setupGen: func(gen *Generator) {
				tmpl := template.Must(template.New("error.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["error.tmpl"] = tmpl
			},
			expectError: false,
			expectFiles: []string{"error.go"},
		},
		{
			name:         "flex_level_with_test_template",
			format:       "error",
			testGenLevel: TestGenFlex,
			setupGen: func(gen *Generator) {
				tmpl := template.Must(template.New("error.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["error.tmpl"] = tmpl
				testTmpl := template.Must(template.New("error_test.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["error_test.tmpl"] = testTmpl
			},
			expectError: false,
			expectFiles: []string{"error.go", "error_test.go"},
		},
		{
			name:         "flex_level_without_test_template",
			format:       "error",
			testGenLevel: TestGenFlex,
			setupGen: func(gen *Generator) {
				tmpl := template.Must(template.New("error.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["error.tmpl"] = tmpl
			},
			expectError: false,
			expectFiles: []string{"error.go"},
		},
		{
			name:         "strict_level_with_test_template",
			format:       "error",
			testGenLevel: TestGenStrict,
			setupGen: func(gen *Generator) {
				tmpl := template.Must(template.New("error.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["error.tmpl"] = tmpl
				testTmpl := template.Must(template.New("error_test.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["error_test.tmpl"] = testTmpl
			},
			expectError: false,
			expectFiles: []string{"error.go", "error_test.go"},
		},
		{
			name:         "strict_level_without_test_template",
			format:       "error",
			testGenLevel: TestGenStrict,
			setupGen: func(gen *Generator) {
				tmpl := template.Must(template.New("error.tmpl").Parse("package {{.PackageName}}"))
				gen.templates["error.tmpl"] = tmpl
			},
			expectError: true,
		},
		{
			name:         "missing_main_template",
			format:       "error",
			testGenLevel: TestGenNone,
			setupGen:     func(*Generator) {},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator with specific test generation level
				gen := New()
				gen.OutputDir = t.TempDir()
				gen.TestGenLevel = test.testGenLevel
				gen.data.PackageName = "testpkg"
				test.setupGen(gen)

				// when: generating format
				err := gen.generateFormat(test.format)

				// then: error and files should match expectations
				if test.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					for _, expectedFile := range test.expectFiles {
						filePath := filepath.Join(gen.OutputDir, expectedFile)
						_, err := os.Stat(filePath)
						assert.NoError(t, err, "expected file %s to exist", expectedFile)
					}
				}
			},
		)
	}
}

// TestRun tests the Run method.
func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setupGen    func(*Generator)
		name        string
		expectError bool
	}{
		{
			name: "successful_run_with_defaults",
			setupGen: func(gen *Generator) {
				gen.OutputDir = t.TempDir()
				gen.Formats = []string{}
			},
			expectError: false,
		},
		{
			name: "successful_run_with_formats",
			setupGen: func(gen *Generator) {
				gen.OutputDir = t.TempDir()
				gen.Formats = nil // Will trigger discovery
			},
			expectError: false,
		},
		{
			name: "invalid_target_directory",
			setupGen: func(gen *Generator) {
				gen.OutputDir = string([]byte{0})
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator
				gen := New()
				test.setupGen(gen)

				// when: running the generator
				err := gen.Run()

				// then: error should match expectation
				if test.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

// TestExportTemplates tests the exportTemplates method.
func TestExportTemplates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setupGen         func(*Generator)
		validateExport   func(*testing.T, string)
		name             string
		expectError      bool
		minExpectedFiles int
	}{
		{
			name: "successful_export_to_new_directory",
			setupGen: func(gen *Generator) {
				gen.ExportDir = filepath.Join(t.TempDir(), "templates")
			},
			expectError:      false,
			minExpectedFiles: 1,
			validateExport: func(t *testing.T, exportDir string) {
				t.Helper()

				// Verify directory was created
				info, err := os.Stat(exportDir)
				require.NoError(t, err)
				assert.True(t, info.IsDir())

				// Verify at least some template files exist
				entries, err := os.ReadDir(exportDir)
				require.NoError(t, err)
				assert.NotEmpty(t, entries)

				// Verify files have .tmpl extension
				for _, entry := range entries {
					assert.Equal(t, ".tmpl", filepath.Ext(entry.Name()), "expected .tmpl file, got %s", entry.Name())
				}
			},
		},
		{
			name: "successful_export_to_existing_directory",
			setupGen: func(gen *Generator) {
				dir := t.TempDir()
				gen.ExportDir = dir
			},
			expectError:      false,
			minExpectedFiles: 1,
			validateExport: func(t *testing.T, exportDir string) {
				t.Helper()

				entries, err := os.ReadDir(exportDir)
				require.NoError(t, err)
				assert.NotEmpty(t, entries)
			},
		},
		{
			name: "export_creates_nested_directory",
			setupGen: func(gen *Generator) {
				gen.ExportDir = filepath.Join(t.TempDir(), "nested", "path", "templates")
			},
			expectError:      false,
			minExpectedFiles: 1,
			validateExport: func(t *testing.T, exportDir string) {
				t.Helper()

				info, err := os.Stat(exportDir)
				require.NoError(t, err)
				assert.True(t, info.IsDir())
			},
		},
		{
			name: "export_overwrites_existing_files",
			setupGen: func(gen *Generator) {
				dir := t.TempDir()
				gen.ExportDir = dir

				// Create a pre-existing file
				existingFile := filepath.Join(dir, "error.tmpl")
				err := os.WriteFile(existingFile, []byte("old content"), 0o600)
				require.NoError(t, err)
			},
			expectError:      false,
			minExpectedFiles: 1,
			validateExport: func(t *testing.T, exportDir string) {
				t.Helper()

				// Verify the file was overwritten with new content
				errorTmpl := filepath.Join(exportDir, "error.tmpl")
				content, err := os.ReadFile(errorTmpl) //nolint:gosec // security is not a concern here
				require.NoError(t, err)
				assert.NotEqual(t, "old content", string(content))
				assert.NotEmpty(t, content)
			},
		},
		{
			name: "invalid_export_directory",
			setupGen: func(gen *Generator) {
				gen.ExportDir = string([]byte{0})
			},
			expectError: true,
		},
		{
			name: "verify_file_permissions",
			setupGen: func(gen *Generator) {
				gen.ExportDir = t.TempDir()
			},
			expectError:      false,
			minExpectedFiles: 1,
			validateExport: func(t *testing.T, exportDir string) {
				t.Helper()

				entries, err := os.ReadDir(exportDir)
				require.NoError(t, err)
				require.NotEmpty(t, entries)

				// Check that files were created successfully
				firstFile := filepath.Join(exportDir, entries[0].Name())
				info, err := os.Stat(firstFile)
				require.NoError(t, err)

				// Verify file is readable
				assert.False(t, info.IsDir())
				assert.Positive(t, info.Size())
			},
		},
		{
			name: "verify_all_embedded_templates_exported",
			setupGen: func(gen *Generator) {
				gen.ExportDir = t.TempDir()
			},
			expectError:      false,
			minExpectedFiles: 8, // At least the core templates
			validateExport: func(t *testing.T, exportDir string) {
				t.Helper()

				entries, err := os.ReadDir(exportDir)
				require.NoError(t, err)

				// Verify common templates exist
				expectedTemplates := []string{"error.tmpl", "wrap.tmpl", "join.tmpl", "common.tmpl"}

				fileNames := make(map[string]bool)
				for _, entry := range entries {
					fileNames[entry.Name()] = true
				}

				for _, expected := range expectedTemplates {
					assert.True(t, fileNames[expected], "expected template %s to be exported", expected)
				}
			},
		},
		{
			name: "verify_template_content_integrity",
			setupGen: func(gen *Generator) {
				gen.ExportDir = t.TempDir()
			},
			expectError:      false,
			minExpectedFiles: 1,
			validateExport: func(t *testing.T, exportDir string) {
				t.Helper()

				// Read exported error.tmpl
				errorTmpl := filepath.Join(exportDir, "error.tmpl")
				content, err := os.ReadFile(errorTmpl) //nolint:gosec // security is not a concern here
				require.NoError(t, err)

				// Verify it's valid template content
				assert.NotEmpty(t, content)
				assert.Contains(t, string(content), "{{", "template should contain template syntax")
			},
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given: a generator with export directory
				gen := New()
				test.setupGen(gen)

				// when: exporting templates
				err := gen.exportTemplates()

				// then: error should match expectation
				if test.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)

					if test.validateExport != nil {
						test.validateExport(t, gen.ExportDir)
					}

					if test.minExpectedFiles > 0 {
						entries, err := os.ReadDir(gen.ExportDir)
						require.NoError(t, err)
						assert.GreaterOrEqual(t, len(entries), test.minExpectedFiles,
							"expected at least %d files, got %d", test.minExpectedFiles, len(entries))
					}
				}
			},
		)
	}
}
