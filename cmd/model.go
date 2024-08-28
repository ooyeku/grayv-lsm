package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ooyeku/grav-lsm/internal/model"
	"github.com/spf13/cobra"
	"os"
)

// modelManager is a pointer to an instance of the ModelManager struct. ModelManager is responsible for managing model definitions,
// including creating, updating, deleting, retrieving, and listing models. It provides functionalities to validate fields and generate
// SQL migration scripts based on a model's definition. The manager uses a map to store the models, where the key is the model's name
// and the value is a pointer to a ModelDefinition struct. The manager can save and load models from a JSON file.
var modelManager *model.ModelManager

// modelCmd is a variable of type *cobra.Command that represents a command for managing data models.
// It has two properties: `Use` which specifies the command string to use, and `Short` which provides a brief description of the command.
// This variable is a member variable of an undisclosed package.
//
// Example usage:
//
//	modelCmd := &cobra.Command{
//	  Use:   "model",
//	  Short: "Manage data models",
//	}
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage data models",
}

// createModelCmd is a variable of type *cobra.Command.
// It represents a command-line command for creating a new model.
// The command has the name "create" and expects one argument: [name].
// When executed, the command calls the runCreateModel function.
// The command is added to the modelCmd command.
//
// The createModelCmd variable is used in the init function to configure
// the command, set the flags, and attach it to the modelCmd command.
// Other commands, such as updateModelCmd, generateModelCmd, and listModelsCmd,
// are also attached to the modelCmd command in the same way.
//
// The runCreateModel function is called when the command is executed.
// It extracts the model name and fields from the command arguments and flags.
// Then, it parses and validates the fields, creates a new model using the
// modelManager, and logs the result.
var createModelCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new model",
	Args:  cobra.ExactArgs(1),
	Run:   runCreateModel,
}

// updateModelCmd represents the command for updating an existing model. It requires one argument, which is the name of the model to be updated. The command includes flags to add or remove fields from the model. The `runUpdateModel` function is invoked when the command is executed. This function retrieves the model by name from the `modelManager`, updates the fields based on the provided flags, and then updates the model in the `modelManager`. Finally, it generates the updated model file and logs the success message.
var updateModelCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update an existing model",
	Args:  cobra.ExactArgs(1),
	Run:   runUpdateModel,
}

// listModelsCmd represents the command for listing all models.
// It is a variable of type *cobra.Command.
// It has the following fields:
//   - Use:   "list"
//   - Short: "List all models"
//   - Run:   runListModels
//
// Usage example:
//
//	modelCmd.AddCommand(listModelsCmd)
//
// Note: This command does not take any arguments.
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

// init initializes the model manager and registers the commands and flags related to model management.
//
// This function should be called once at the start of the program to set up the necessary components.
// It performs the following tasks:
// - Initializes the model manager using model.NewModelManager().
// - Sets up the flags for the createModelCmd, updateModelCmd, and generateModelCmd commands.
// - Registers the createModelCmd, updateModelCmd, and generateModelCmd commands under the modelCmd command.
// - Registers the modelCmd command under the RootCmd command.
//
// Example usage:
//
//	init()
//
// Note: This function is not intended to be called directly by users of this package.
func init() {
	modelManager = model.NewModelManager()

	createModelCmd.Flags().StringSlice("fields", []string{}, "Comma-separated list of fields in the format name:type")
	updateModelCmd.Flags().StringSlice("add-fields", []string{}, "Comma-separated list of fields to add in the format name:type")
	updateModelCmd.Flags().StringSlice("remove-fields", []string{}, "Comma-separated list of field names to remove")

	generateModelCmd.Flags().String("app", "", "Name of the Grav app to generate the model in")

	modelCmd.AddCommand(createModelCmd)
	modelCmd.AddCommand(updateModelCmd)
	RootCmd.AddCommand(modelCmd)
	modelCmd.AddCommand(listModelsCmd)
	modelCmd.AddCommand(generateModelCmd)
}

// runCreateModel creates a new model with the given name and fields. It parses the fields, creates a new model definition,
// and adds it to the model manager's models map. It then saves the models to the storage file. If there are any errors
// during the process, it logs an error message.
//
// Parameters:
// - cmd: The cobra.Command object representing the command.
// - args: The command arguments, where args[0] is the model name.
//
// Example usage:
//
//	runCreateModel(cmd, args)
//
// Note: This function is not intended to be called directly by users of this package.
func runCreateModel(cmd *cobra.Command, args []string) {
	modelName := args[0]
	fields, _ := cmd.Flags().GetStringSlice("fields")

	modelFields, err := parseFields(fields)
	if err != nil {
		log.WithError(err).Error("Failed to parse fields")
		return
	}

	err = modelManager.CreateModel(modelName, modelFields)
	if err != nil {
		log.WithError(err).Errorf("Failed to create model %s", modelName)
		return
	}

	log.Infof("Model %s created successfully", modelName)
}

// runUpdateModel updates the fields of an existing model.
// It retrieves the model with the given name from the model manager.
// If the model does not exist, an error is logged and the function returns.
// If addFields is provided, the new fields are parsed and appended to the model definition.
// If removeFields is provided, the fields are removed from the model definition.
// The updated model is then passed to the model manager to update the model's fields.
// If there is an error during the update, an error is logged and the function returns.
// After updating the model, a new model file is generated for the updated model definition.
// If there is an error during the generation of the model file, an error is logged and the function returns.
// Finally, a success message is logged indicating that the model was updated successfully.
func runUpdateModel(cmd *cobra.Command, args []string) {
	modelName := args[0]
	addFields, _ := cmd.Flags().GetStringSlice("add-fields")
	removeFields, _ := cmd.Flags().GetStringSlice("remove-fields")

	modelDef, err := modelManager.GetModel(modelName)
	if err != nil {
		log.WithError(err).Errorf("Failed to get model %s", modelName)
		return
	}

	if len(addFields) > 0 {
		newFields, err := parseFields(addFields)
		if err != nil {
			log.WithError(err).Error("Failed to parse new fields")
			return
		}
		modelDef.Fields = append(modelDef.Fields, newFields...)
	}

	if len(removeFields) > 0 {
		modelDef.Fields = removeFieldsFromModel(modelDef.Fields, removeFields)
	}

	err = modelManager.UpdateModel(modelName, modelDef.Fields)
	if err != nil {
		log.WithError(err).Errorf("Failed to update model %s", modelName)
		return
	}

	err = model.GenerateModelFile(modelDef)
	if err != nil {
		log.WithError(err).Errorf("Failed to generate updated model file for %s", modelName)
		return
	}

	log.Infof("Model %s updated successfully", modelName)
}

// runListModels lists all available models. It retrieves the list of models from the ModelManager and
// logs them in the output. If no models are found, it logs a message indicating that no models were found.
func runListModels(cmd *cobra.Command, args []string) {
	models := modelManager.ListModels()
	if len(models) == 0 {
		log.Info("No models found.")
		return
	}

	log.Info("Available models:")
	for _, m := range models {
		log.Infof("- %s", m)
	}
}

// runGenerateModel generates a model file based on the provided model name and app name.
// It retrieves the model definition from the model manager, sets the output directory if the app name is provided,
// and generates the model file using the model.GenerateModelFile method. If successful, it logs the
// success message along with the output directory if applicable.
func runGenerateModel(cmd *cobra.Command, args []string) {
	modelName := args[0]
	appName, _ := cmd.Flags().GetString("app")

	modelDef, err := modelManager.GetModel(modelName)
	if err != nil {
		log.WithError(err).Errorf("Failed to get model %s", modelName)
		return
	}

	var outputDir string
	if appName != "" {
		outputDir = filepath.Join(appName+"_grav", "internal", "models")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.WithError(err).Errorf("Failed to create directory for app %s", appName)
			return
		}

		// Set the output directory for the model generator
		modelDef.SetOutputDir(outputDir)
	}

	err = model.GenerateModelFile(modelDef)
	if err != nil {
		log.WithError(err).Errorf("Failed to generate model file for %s", modelName)
		return
	}

	if appName != "" {
		log.Infof("Model file for %s generated successfully in %s", modelName, outputDir)
	} else {
		log.Infof("Model file for %s generated successfully", modelName)
	}
}

// parseFields parses the given list of fields and returns a slice of model.Field.
// It splits each field by ":" and creates a model.NewField with the parts.
// If a field does not have exactly two parts, it returns an error with the invalid field format message.
// The returned model.Field is populated with the name, fieldType, tag, isNull, and isPrimary values.
// The name is the first part of the field, the fieldType is the second part, and the tag is generated
// by using the name converted to lowercase as the json tag value. The isNull flag is set to false,
// and the isPrimary flag is set to true if the name is either "ID", "Id", or "id".
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
