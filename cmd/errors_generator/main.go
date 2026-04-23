// Package main provides a command-line interface to generate code based on templates.
//
// This package is responsible for generating code based on templates.
// It uses the following packages to fulfill its task:
//   - github.com/emiliogrv/errors/internal/template: To embed the default templates into the binary.
//
// The command-line interface is divided into two main parts:
//   - The first part is used to specify the output directory and the package name.
//   - The second part is used to specify the formats to generate.
//
// The formats can be either specified as a comma-separated list of formats,
// or the special value "all" to generate all formats.
//
// The generator will generate code for the main file, and test file, based on the formats specified.
//
// The generator will also generate code for the format specific files, if a user template is specified.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"

	internaltemplate "github.com/emiliogrv/errors/internal/template"
)

type (
	Generator struct {
		InputDir     string
		OutputDir    string
		ExportDir    string
		Formats      []string
		TestGenLevel string
		templates    map[string]*template.Template
		data         TemplateData
	}

	TemplateData struct {
		PackageName   string
		Date          string
		Version       string
		WithGenHeader bool
	}
)

const (
	TestGenNone   = "none"
	TestGenFlex   = "flex"
	TestGenStrict = "strict"

	Version = "0.0.1"

	folderPermissions = 0o750
	filePermissions   = 0o600
	emptyString       = ""

	zero = 0
	one  = 1
	two  = 2
	three = 3
	dotSeparator = "."
)

func New() *Generator {
	return &Generator{
		templates: make(map[string]*template.Template),
		data: TemplateData{
			Date:          time.Now().Format(time.RFC3339),
			Version:       Version,
			WithGenHeader: true,
		},
		Formats:      []string{"attr", "common", "error", "join", "json", "map", "string", "wrap"},
		TestGenLevel: TestGenNone,
	}
}

func main() {
	generator := New()

	flag.StringVar(
		&generator.InputDir,
		"input-dir",
		emptyString,
		"Path to user templates directory (optional)",
	)
	flag.StringVar(&generator.OutputDir, "output-dir", emptyString, "Output directory for generated files")
	flag.StringVar(
		&generator.data.PackageName,
		"package",
		"errors",
		"Package name for generated code (default: errors)",
	)
	flag.BoolVar(
		&generator.data.WithGenHeader,
		"with-gen-header",
		true,
		"Include generated message in generated code (default: true)",
	)
	flag.StringVar(
		&generator.ExportDir,
		"export-dir",
		emptyString,
		"Export default templates to the specified directory and exit",
	)
	formats := flag.String(
		"formats",
		emptyString,
		"Comma-separated list of formats to generate, or 'all' to generate all formats (default: core)",
	)
	testGen := flag.String("test-gen", TestGenNone, "Test generation level: none, flex, strict (default: none)")
	versionFlag := flag.String("version", emptyString,
		"Version to use for generated code (auto-detects from git if not specified)")
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(zero)
	}

	// Handle export-dir flag
	if generator.ExportDir != emptyString {
		err := generator.exportTemplates()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Default templates exported to: " + generator.ExportDir)
		os.Exit(zero)
	}

	if generator.OutputDir == emptyString {
		flag.Usage()
		os.Exit(one)
	}

	// Set version based on flag or git detection
	if *versionFlag != emptyString {
		generator.data.Version = *versionFlag
	} else {
		gitVersion, err := getLastGitVersion()
		if err != nil {
			log.Fatalf("Error: could not detect git version and no version provided. "+
			"Use -version flag or ensure git tags are available: %v", err)
		}

		if gitVersion == emptyString {
			log.Fatal("Error: no semver git tags found and no version provided. Use -version flag or create a git tag")
		}

		generator.data.Version = gitVersion
	}

	err := generator.validateTestGenLevel(*testGen)
	if err != nil {
		log.Fatalln(err)
	}

	generator.loadFormats(*formats)

	err = generator.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func (receiver *Generator) Run() error {
	// Load embedded templates first
	err := receiver.loadEmbeddedTemplates()
	if err != nil {
		return fmt.Errorf("loading embedded templates: %w", err)
	}

	// Load and override with user templates if specified
	if receiver.InputDir != emptyString {
		err = receiver.loadUserTemplates(receiver.InputDir)
		if err != nil {
			return fmt.Errorf("loading user templates: %w", err)
		}
	}

	// If formats is set to "all", discover all available formats from templates
	if receiver.Formats == nil {
		receiver.Formats = receiver.discoverTemplateFormats()
	}

	// Create target directory if it doesn't exist
	err = os.MkdirAll(receiver.OutputDir, folderPermissions)
	if err != nil {
		return fmt.Errorf("creating target directory: %w", err)
	}

	// Generate files for each requested format
	for _, formatName := range receiver.Formats {
		err = receiver.generateFormat(formatName)
		if err != nil {
			return fmt.Errorf("generating format %s: %w", formatName, err)
		}
	}

	if receiver.TestGenLevel != TestGenNone {
		err = receiver.generateFile("compatibility_test.tmpl", "compatibility_test.go")
		if err != nil {
			return fmt.Errorf("generating compatibility test file: %w", err)
		}
	}

	return nil
}

func (receiver *Generator) validateTestGenLevel(level string) error {
	switch level {
	case TestGenNone, TestGenFlex, TestGenStrict:
		receiver.TestGenLevel = level

		return nil
	default:
		//nolint:err113 // dynamic is expected
		return fmt.Errorf("invalid test generation level: %s. Must be one of: none, flex, strict", level)
	}
}

func (receiver *Generator) loadFormats(formats string) {
	if formats == emptyString {
		return
	}

	if formats == "all" {
		// Clear default formats as we'll load all available ones
		receiver.Formats = nil

		return
	}

	// Append user-specified formats
	receiver.Formats = append(receiver.Formats, strings.Split(formats, ",")...)
}

func (receiver *Generator) discoverTemplateFormats() []string {
	formats := make(map[string]struct{})

	// Collect all formats from templates
	for name := range receiver.templates {
		if strings.HasSuffix(name, ".tmpl") && !strings.HasSuffix(name, "_test.tmpl") {
			formatName := strings.TrimSuffix(name, ".tmpl")
			formats[formatName] = struct{}{}
		}
	}

	// Convert to slice
	result := make([]string, zero, len(formats))
	for formatName := range formats {
		result = append(result, formatName)
	}

	return result
}

func (receiver *Generator) loadEmbeddedTemplates() error {
	entries, err := fs.ReadDir(internaltemplate.DefaultTemplates, ".")
	if err != nil {
		return fmt.Errorf("reading embedded templates dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		content, errRF := fs.ReadFile(internaltemplate.DefaultTemplates, name)
		if errRF != nil {
			return fmt.Errorf("reading embedded template %s: %w", name, errRF)
		}

		tmpl, errN := template.New(name).Funcs(sprig.FuncMap()).Parse(string(content))
		if errN != nil {
			return fmt.Errorf("parsing embedded template %s: %w", name, errN)
		}

		receiver.templates[name] = tmpl
	}

	return nil
}

func (receiver *Generator) loadUserTemplates(dir string) error {
	err := filepath.WalkDir(
		dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf("getting relative path for user template %s: %w", path, err)
			}

			content, err := os.ReadFile(path) //nolint:gosec // security is not a concern here
			if err != nil {
				return fmt.Errorf("reading user template %s: %w", relPath, err)
			}

			tmpl, err := template.New(relPath).Funcs(sprig.FuncMap()).Parse(string(content))
			if err != nil {
				return fmt.Errorf("parsing user template %s: %w", relPath, err)
			}

			// User templates always override embedded ones
			receiver.templates[relPath] = tmpl

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("walking user templates dir: %w", err)
	}

	return nil
}

func (receiver *Generator) generateFormat(formatName string) error {
	// Generate main file
	err := receiver.generateFile(formatName+".tmpl", formatName+".go")
	if err != nil {
		return fmt.Errorf("generating main file: %w", err)
	}

	// Handle test file generation based on level
	testTemplate := formatName + "_test.tmpl"
	hasTestTemplate := receiver.hasTemplate(testTemplate)

	switch receiver.TestGenLevel {
	case TestGenNone:
		return nil
	case TestGenFlex:
		if hasTestTemplate {
			err = receiver.generateFile(testTemplate, formatName+"_test.go")
			if err != nil {
				return fmt.Errorf("generating test file: %w", err)
			}
		}

		return nil
	case TestGenStrict:
		if !hasTestTemplate {
			//nolint:err113 // dynamic is expected
			return fmt.Errorf("test template not found for format %s (required in strict mode)", formatName)
		}

		err = receiver.generateFile(testTemplate, formatName+"_test.go")
		if err != nil {
			return fmt.Errorf("generating test file: %w", err)
		}

		return nil
	}

	return nil
}

func (receiver *Generator) generateFile(templateName, outputName string) (err error) {
	tmpl, ok := receiver.templates[templateName]
	if !ok {
		return fmt.Errorf("template not found: %s", templateName) //nolint:err113 // dynamic is expected
	}

	// Execute template to buffer first
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, receiver.data)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting generated code: %w", err)
	}

	// Prepare output file
	outputPath := filepath.Join(receiver.OutputDir, outputName)

	outputFile, err := os.Create(outputPath) //nolint:gosec // security is not a concern here
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer func(outputFile *os.File) {
		errC := outputFile.Close()
		if err != nil && errC != nil {
			err = fmt.Errorf("closing output file: %w", errC)
		}
	}(outputFile)

	// Write formatted content to file
	_, err = outputFile.Write(formatted)
	if err != nil {
		return fmt.Errorf("writing formatted output: %w", err)
	}

	return nil
}

func (receiver *Generator) hasTemplate(templateName string) bool {
	_, exists := receiver.templates[templateName]

	return exists
}

func (receiver *Generator) exportTemplates() error {
	// Create export directory if it doesn't exist
	err := os.MkdirAll(receiver.ExportDir, folderPermissions)
	if err != nil {
		return fmt.Errorf("creating export directory: %w", err)
	}

	// Read all embedded templates
	entries, err := fs.ReadDir(internaltemplate.DefaultTemplates, ".")
	if err != nil {
		return fmt.Errorf("reading embedded templates: %w", err)
	}

	// Export each template file
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Read template content from embedded FS
		content, errRF := fs.ReadFile(internaltemplate.DefaultTemplates, name)
		if errRF != nil {
			return fmt.Errorf("reading template %s: %w", name, errRF)
		}

		// Write to export directory
		outputPath := filepath.Join(receiver.ExportDir, name)

		errW := os.WriteFile(outputPath, content, filePermissions)
		if errW != nil {
			return fmt.Errorf("writing template %s: %w", name, errW)
		}
	}

	return nil
}

// getLastGitVersion attempts to get the last version from git tags
// It handles semver versions and strips any suffixes like -alpha, -beta, etc.
func getLastGitVersion() (string, error) {
	// Get all git tags sorted by version
	cmd := exec.CommandContext(context.Background(), "git", "tag", "--sort=-v:refname")

	output, err := cmd.Output()
	if err != nil {
		return emptyString, fmt.Errorf("failed to get git tags: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == zero || (len(tags) == one && tags[zero] == emptyString) {
		return emptyString, nil
	}

	// Semver regex to match versions like v1.2.3, 1.2.3, etc
	semverRegex := regexp.MustCompile(`^v?\d+\.\d+\.\d+`)

	// Find the first tag that matches semver pattern
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == emptyString {
			continue
		}

		match := semverRegex.FindString(tag)
		if match != emptyString {
			// Remove 'v' prefix if present
			match = strings.TrimPrefix(match, "v")

			// Parse and increment the last segment
			parts := strings.Split(match, dotSeparator)
			if len(parts) == three {
				// Convert last part to int and increment
				lastPart := parts[two]

				var lastNum int

				_, err = fmt.Sscanf(lastPart, "%d", &lastNum)
				if err == nil {
					lastNum++
					// Reconstruct version
					return fmt.Sprintf("%s.%s.%d", parts[0], parts[1], lastNum), nil
				}
			}

			return match, nil
		}
	}

	return emptyString, nil
}
