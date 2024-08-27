package cmd

import (
	"os"

	"github.com/ooyeku/grav-lsm/internal/database/lsm"
	"github.com/ooyeku/grav-lsm/internal/database/seed"
	"github.com/ooyeku/grav-lsm/internal/orm"
	"github.com/ooyeku/grav-lsm/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dbManager *lsm.DBLifecycleManager
var log = logrus.New()
var cfg *config.Config

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	var err error
	cfg, err = config.LoadConfig("config.json")
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}
	dbManager = lsm.NewDBLifecycleManager(cfg)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the database lifecycle",
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the database Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dbManager.BuildImage(); err != nil {
			log.WithError(err).Error("Error building database image")
		} else {
			log.Info("Database image built successfully")
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dbManager.StartContainer(); err != nil {
			log.WithError(err).Error("Error starting database container")
		} else {
			log.Info("Database container started successfully")
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dbManager.StopContainer(); err != nil {
			log.WithError(err).Error("Error stopping database container")
		} else {
			log.Info("Database container stopped successfully")
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dbManager.RemoveContainer(); err != nil {
			log.WithError(err).Error("Error removing database container")
		} else {
			log.Info("Database container removed successfully")
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the health of the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := dbManager.GetStatus()
		if err != nil {
			log.WithError(err).Error("Error checking database status")
			return
		}
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := orm.NewConnection(&cfg.Database)
		if err != nil {
			log.WithError(err).Error("Error connecting to database")
			return
		}
		defer conn.Close()

		seeder := seed.NewSeeder(conn.GetDB())
		err = seeder.LoadSeeds("./seeds")
		if err != nil {
			log.WithError(err).Error("Error loading seeds")
			return
		}

		err = seeder.Seed()
		if err != nil {
			log.WithError(err).Error("Error seeding database")
		} else {
			log.Info("Database seeded successfully")
		}
	},
}

func init() {
	dbCmd.AddCommand(buildCmd)
	dbCmd.AddCommand(startCmd)
	dbCmd.AddCommand(stopCmd)
	dbCmd.AddCommand(removeCmd)
	dbCmd.AddCommand(statusCmd)
	dbCmd.AddCommand(seedCmd)
	RootCmd.AddCommand(dbCmd)
}
