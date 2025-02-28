package generator

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/springwell/cli/pkg/config"
	"github.com/springwell/cli/pkg/util"
)

// EntityGenerator generates entity-related code
type EntityGenerator struct {
	Config     *config.Config
	ProjectDir string
}

// NewEntityGenerator creates a new EntityGenerator
func NewEntityGenerator(config *config.Config, projectDir string) *EntityGenerator {
	return &EntityGenerator{
		Config:     config,
		ProjectDir: projectDir,
	}
}

// GenerateEntity generates an entity and its related components
func (g *EntityGenerator) GenerateEntity(name, fieldsStr, relationsStr, tableName string, audit, lombok, generateDto, generateRepo, generateService, generateController bool) error {
	// Parse fields and relations
	fields, err := util.ParseFieldDefinitions(fieldsStr)
	if err != nil {
		return err
	}

	relations, err := util.ParseRelationships(relationsStr)
	if err != nil {
		return err
	}

	// If table name is not provided, generate it from entity name
	if tableName == "" {
		tableName = util.ToDatabaseTableName(name)
	}

	// Create template data
	data := map[string]interface{}{
		"name":       name,
		"nameCamel":  util.ToJavaVariableName(name),
		"namePlural": util.ToJavaVariableName(name) + "s", // Simple pluralization, can be improved
		"package":    g.Config.Project.Package,
		"tableName":  tableName,
		"fields":     fields,
		"relations":  relations,
		"audit":      audit,
		"lombok":     lombok,
	}

	// Generate entity
	if err := g.generateFromTemplate("entity/entity.tmpl", filepath.Join(g.ProjectDir, "src/main/java", strings.ReplaceAll(g.Config.Project.Package, ".", "/"), "domain/entity", name+".java"), data); err != nil {
		return err
	}

	// Generate repository
	if generateRepo {
		if err := g.generateFromTemplate("entity/repository.tmpl", filepath.Join(g.ProjectDir, "src/main/java", strings.ReplaceAll(g.Config.Project.Package, ".", "/"), "repository", name+"Repository.java"), data); err != nil {
			return err
		}
	}

	// Generate service
	if generateService {
		if err := g.generateFromTemplate("entity/service.tmpl", filepath.Join(g.ProjectDir, "src/main/java", strings.ReplaceAll(g.Config.Project.Package, ".", "/"), "service", name+"Service.java"), data); err != nil {
			return err
		}
	}

	// Generate controller
	if generateController {
		if err := g.generateFromTemplate("entity/controller.tmpl", filepath.Join(g.ProjectDir, "src/main/java", strings.ReplaceAll(g.Config.Project.Package, ".", "/"), "controller", name+"Controller.java"), data); err != nil {
			return err
		}
	}

	// Generate DTO
	if generateDto {
		if err := g.generateFromTemplate("entity/dto.tmpl", filepath.Join(g.ProjectDir, "src/main/java", strings.ReplaceAll(g.Config.Project.Package, ".", "/"), "dto", name+"DTO.java"), data); err != nil {
			return err
		}
	}

	return nil
}

// generateFromTemplate generates a file from a template
func (g *EntityGenerator) generateFromTemplate(templatePath, outputPath string, data map[string]interface{}) error {
	// Determine the template directory
	var templateDir string
	if _, err := os.Stat(filepath.Join(g.ProjectDir, g.Config.Templates.Directory)); err == nil {
		// Use project-specific templates if available
		templateDir = filepath.Join(g.ProjectDir, g.Config.Templates.Directory)
	} else {
		// Use built-in templates
		templateDir = filepath.Join("pkg", "templates")
	}

	// Read the template file
	templateContent, err := ioutil.ReadFile(filepath.Join(templateDir, templatePath))
	if err != nil {
		return err
	}

	// Create a new template
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"toLowerCase": strings.ToLower,
	}).Parse(string(templateContent))
	if err != nil {
		return err
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// Write the generated file
	return util.WriteFile(outputPath, buf.String())
}
