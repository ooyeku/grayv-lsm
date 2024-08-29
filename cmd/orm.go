/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/ooyeku/grav-lsm/internal/orm"
	"github.com/ooyeku/grav-lsm/pkg/config"
	"github.com/ooyeku/grav-lsm/pkg/utils"
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

func init() {
	ormCmd.AddCommand(queryCmd)
	ormCmd.AddCommand(createUserCmd)
	RootCmd.AddCommand(ormCmd)
	createUserCmd.Flags().String("username", "", "Username for the new user")
	createUserCmd.Flags().String("email", "", "Email for the new user")
	createUserCmd.Flags().String("password", "", "Password for the new user")
	createUserCmd.MarkFlagRequired("username")
	createUserCmd.MarkFlagRequired("email")
	createUserCmd.MarkFlagRequired("password")
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
