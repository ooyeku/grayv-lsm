package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ooyeku/grayv-lsm/internal/model"
	"github.com/ooyeku/grayv-lsm/internal/orm"
	"github.com/ooyeku/grayv-lsm/pkg/config"
	"github.com/spf13/cobra"
)

var modelManager *model.ModelManager

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage data models",
}

var createModelCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new model",
	Args:  cobra.ExactArgs(1),
	Run:   runCreateModel,
}

var updateModelCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update an existing model",
	Args:  cobra.ExactArgs(1),
	Run:   runUpdateModel,
}

var listModelsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all models",
	Run:   runListModels,
}

var generateModelCmd = &cobra.Command{
	Use:   "generate [name]",
	Short: "Generate Go code for an existing model",
	Args:  cobra.ExactArgs(1),
	Run:   runGenerateModel,
}

func init() {
	modelManager = model.NewModelManager()

	createModelCmd.Flags().StringSlice("fields", []string{}, "Comma-separated list of fields in the format name:type")
	updateModelCmd.Flags().StringSlice("add-fields", []string{}, "Comma-separated list of fields to add in the format name:type")
	updateModelCmd.Flags().StringSlice("remove-fields", []string{}, "Comma-separated list of field names to remove")

	generateModelCmd.Flags().String("app", "", "Name of the Grayv app to generate the model in")

	modelCmd.AddCommand(createModelCmd)
	modelCmd.AddCommand(updateModelCmd)
	RootCmd.AddCommand(modelCmd)
	modelCmd.AddCommand(listModelsCmd)
	modelCmd.AddCommand(generateModelCmd)
}

func runCreateModel(cmd *cobra.Command, args []string) {
	modelName := args[0]
	fields, _ := cmd.Flags().GetStringSlice("fields")

	modelFields, err := parseFields(fields)
	if err != nil {
		log.WithError(err).Error("Failed to parse fields")
		return
	}

	conn, err := getDBConnection()
	if err != nil {
		log.WithError(err).Error("Failed to get database connection")
		return
	}
	defer conn.Close()

	fieldsJSON, err := json.Marshal(modelFields)
	if err != nil {
		log.WithError(err).Error("Failed to marshal model fields")
		return
	}

	query := "INSERT INTO models (name, fields) VALUES ($1, $2)"
	_, err = conn.Query(query, modelName, fieldsJSON)
	if err != nil {
		log.WithError(err).Errorf("Failed to create model %s", modelName)
		return
	}

	log.Infof("Model %s created successfully", modelName)
}

func runUpdateModel(cmd *cobra.Command, args []string) {
	modelName := args[0]
	addFields, _ := cmd.Flags().GetStringSlice("add-fields")
	removeFields, _ := cmd.Flags().GetStringSlice("remove-fields")

	conn, err := getDBConnection()
	if err != nil {
		log.WithError(err).Error("Failed to get database connection")
		return
	}
	defer conn.Close()

	var fieldsJSON []byte
	rows, err := conn.Query("SELECT fields FROM models WHERE name = $1", modelName)
	if err != nil {
		log.WithError(err).Errorf("Failed to get model %s", modelName)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&fieldsJSON)
		if err != nil {
			log.WithError(err).Error("Failed to scan model fields")
			return
		}

		var modelFields []model.Field
		err = json.Unmarshal(fieldsJSON, &modelFields)
		if err != nil {
			log.WithError(err).Error("Failed to unmarshal model fields")
			return
		}

		if len(addFields) > 0 {
			newFields, err := parseFields(addFields)
			if err != nil {
				log.WithError(err).Error("Failed to parse new fields")
				return
			}
			modelFields = append(modelFields, newFields...)
		}

		if len(removeFields) > 0 {
			modelFields = removeFieldsFromModel(modelFields, removeFields)
		}

		updatedFieldsJSON, err := json.Marshal(modelFields)
		if err != nil {
			log.WithError(err).Error("Failed to marshal updated model fields")
			return
		}

		_, err = conn.Query("UPDATE models SET fields = $1 WHERE name = $2", updatedFieldsJSON, modelName)
		if err != nil {
			log.WithError(err).Errorf("Failed to update model %s", modelName)
			return
		}

		log.Infof("Model %s updated successfully", modelName)
	}
}

func runListModels(cmd *cobra.Command, args []string) {
	conn, err := getDBConnection()
	if err != nil {
		log.WithError(err).Error("Failed to get database connection")
		return
	}
	defer conn.Close()

	models, err := listModelsFromDB(conn)
	if err != nil {
		log.WithError(err).Error("Failed to list models")
		return
	}

	if len(models) == 0 {
		log.Info("No models found.")
	} else {
		log.Info("Available models:")
		for _, m := range models {
			log.Infof("- %s", m)
		}
	}
}

func listModelsFromDB(conn *orm.Connection) ([]string, error) {
	query := "SELECT name FROM models"
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		models = append(models, name)
	}

	return models, rows.Err()
}

func runGenerateModel(cmd *cobra.Command, args []string) {
	modelName := args[0]

	conn, err := getDBConnection()
	if err != nil {
		log.WithError(err).Error("Failed to get database connection")
		return
	}
	defer conn.Close()

	var fieldsJSON []byte
	rows, err := conn.Query("SELECT fields FROM models WHERE name = $1", modelName)
	if err != nil {
		log.WithError(err).Errorf("Failed to get model %s from database", modelName)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&fieldsJSON)
		if err != nil {
			log.WithError(err).Error("Failed to scan model fields")
			return
		}

		var modelFields []model.Field
		err = json.Unmarshal(fieldsJSON, &modelFields)
		if err != nil {
			log.WithError(err).Error("Failed to unmarshal model fields")
			return
		}

		modelDef := &model.ModelDefinition{
			Name:   modelName,
			Fields: modelFields,
		}

		err = model.GenerateModelFile(modelDef)
		if err != nil {
			log.WithError(err).Errorf("Failed to generate model file for %s", modelName)
			return
		}

		log.Infof("Model %s generated successfully", modelName)
	}
}

// parseFields parses the given list of fields and returns a slice of model.Field.
// If no error occurs, it returns the slice of model.Field and a nil error. Otherwise, it returns nil and an error.
func parseFields(fields []string) ([]model.Field, error) {
	var modelFields []model.Field
	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field format: %s", field)
		}
		name := parts[0]
		fieldType := parts[1]
		tag := fmt.Sprintf(`json:"%s"`, strings.ToLower(name))
		isNull := false
		isPrimary := name == "ID" || name == "Id" || name == "id"
		modelFields = append(modelFields, model.NewField(name, fieldType, tag, isNull, isPrimary))
	}
	return modelFields, nil
}

// removeFieldsFromModel removes specified fields from a list of model fields and returns the updated list.
//
// Parameters:
// - fields: The list of model fields to remove from.
// - fieldsToRemove: The list of field names to remove.
//
// Returns:
// - updatedFields: The list of model fields after removing the specified fields.
func removeFieldsFromModel(fields []model.Field, fieldsToRemove []string) []model.Field {
	var updatedFields []model.Field
	for _, field := range fields {
		if !contains(fieldsToRemove, field.Name) {
			updatedFields = append(updatedFields, field)
		}
	}
	return updatedFields
}

// contains checks if a string item is present in a string slice.
// It returns true if the item is found, and false otherwise.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getDBConnection() (*orm.Connection, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	conn, err := orm.NewConnection(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return conn, nil
}
