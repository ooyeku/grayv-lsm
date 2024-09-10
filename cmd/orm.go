/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/ooyeku/grayv-lsm/internal/orm"
	"github.com/ooyeku/grayv-lsm/pkg/config"
	"github.com/ooyeku/grayv-lsm/pkg/utils"
	"github.com/spf13/cobra"
)

// ormCmd represents the orm command
var ormCmd = &cobra.Command{
	Use:   "orm",
	Short: "Perform ORM operations",
}

var queryCmd = &cobra.Command{
	Use:   "query [SQL]",
	Short: "Execute a SQL query",
	Args:  cobra.ExactArgs(1),
	Run:   runQuery,
}

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user in the database",
	Run:   runCreateUser,
}

var updateUserCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Update an existing user in the database",
	Run:   runUpdateUser,
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete a user from the database",
	Run:   runDeleteUser,
}

var listUsersCmd = &cobra.Command{
	Use:   "list-users",
	Short: "List all users in the database",
	Run:   runListUsers,
}

func init() {
	ormCmd.AddCommand(queryCmd)
	ormCmd.AddCommand(createUserCmd)
	ormCmd.AddCommand(updateUserCmd)
	ormCmd.AddCommand(deleteUserCmd)
	ormCmd.AddCommand(listUsersCmd)
	RootCmd.AddCommand(ormCmd)

	// Existing flags for createUserCmd...

	updateUserCmd.Flags().Int("id", 0, "ID of the user to update")
	updateUserCmd.Flags().String("username", "", "New username for the user")
	updateUserCmd.Flags().String("email", "", "New email for the user")
	updateUserCmd.Flags().String("password", "", "New password for the user")
	updateUserCmd.MarkFlagRequired("id")

	deleteUserCmd.Flags().Int("id", 0, "ID of the user to delete")
	deleteUserCmd.MarkFlagRequired("id")

	createUserCmd.Flags().String("username", "", "Username for the new user")
	createUserCmd.Flags().String("email", "", "Email for the new user")
	createUserCmd.Flags().String("password", "", "Password for the new user")
	createUserCmd.MarkFlagRequired("username")
	createUserCmd.MarkFlagRequired("email")
	createUserCmd.MarkFlagRequired("password")

	// ... (existing code)
}

func runQuery(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}

	conn, err := orm.NewConnection(&cfg.Database)
	if err != nil {
		log.WithError(err).Error("Error connecting to database")
		return
	}
	defer conn.Close()

	query := args[0]
	rows, err := conn.Query(query)
	if err != nil {
		log.WithError(err).Error("Error executing query")
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.WithError(err).Error("Error getting column names")
		return
	}

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.WithError(err).Error("Error scanning row")
			continue
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}

		fmt.Println(rowData)
	}

	if err := rows.Err(); err != nil {
		log.WithError(err).Error("Error iterating over rows")
	}
}

func runCreateUser(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}

	conn, err := orm.NewConnection(&cfg.Database)
	if err != nil {
		log.WithError(err).Error("Error connecting to database")
		return
	}
	defer conn.Close()

	username, _ := cmd.Flags().GetString("username")
	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.WithError(err).Error("Error hashing password")
		return
	}

	query := "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)"
	_, err = conn.Query(query, username, email, hashedPassword)
	if err != nil {
		log.WithError(err).Error("Error creating new user")
		return
	}

	log.Info("New user created successfully")
}

func runUpdateUser(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}

	conn, err := orm.NewConnection(&cfg.Database)
	if err != nil {
		log.WithError(err).Error("Error connecting to database")
		return
	}
	defer conn.Close()

	id, _ := cmd.Flags().GetInt("id")
	username, _ := cmd.Flags().GetString("username")
	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")

	updateFields := make(map[string]interface{})
	if username != "" {
		updateFields["username"] = username
	}
	if email != "" {
		updateFields["email"] = email
	}
	if password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.WithError(err).Error("Error hashing password")
			return
		}
		updateFields["password_hash"] = hashedPassword
	}

	if len(updateFields) == 0 {
		log.Error("No fields to update")
		return
	}

	query := "UPDATE users SET "
	var values []interface{}
	i := 0
	for field, value := range updateFields {
		if i > 0 {
			query += ", "
		}
		query += field + " = $" + fmt.Sprintf("%d", i+1)
		values = append(values, value)
		i++
	}
	query += " WHERE id = $" + fmt.Sprintf("%d", i+1)
	values = append(values, id)

	_, err = conn.GetDB().Exec(query, values...)
	if err != nil {
		log.WithError(err).Error("Error updating user")
		return
	}

	log.Info("User updated successfully")
}

func runDeleteUser(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}

	conn, err := orm.NewConnection(&cfg.Database)
	if err != nil {
		log.WithError(err).Error("Error connecting to database")
		return
	}
	defer conn.Close()

	id, _ := cmd.Flags().GetInt("id")

	query := "DELETE FROM users WHERE id = $1"
	_, err = conn.GetDB().Exec(query, id)
	if err != nil {
		log.WithError(err).Error("Error deleting user")
		return
	}

	log.Info("User deleted successfully")
}

func runListUsers(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}

	conn, err := orm.NewConnection(&cfg.Database)
	if err != nil {
		log.WithError(err).Error("Error connecting to database")
		return
	}
	defer conn.Close()

	query := "SELECT id, username, email FROM users"
	rows, err := conn.GetDB().Query(query)
	if err != nil {
		log.WithError(err).Error("Error querying users")
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, email string
		if err := rows.Scan(&id, &username, &email); err != nil {
			log.WithError(err).Error("Error scanning user row")
			continue
		}
		users = append(users, map[string]interface{}{
			"id":       id,
			"username": username,
			"email":    email,
		})
	}

	if len(users) == 0 {
		log.Info("No users found")
	} else {
		log.Info("Users:")
		for _, user := range users {
			log.Infof("ID: %d, Username: %s, Email: %s", user["id"], user["username"], user["email"])
		}
	}
}
