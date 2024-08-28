package cmd

import (
	"os"
	"strconv"

	"github.com/ooyeku/grav-lsm/internal/database/lsm"
	"github.com/ooyeku/grav-lsm/internal/database/migration"
	"github.com/ooyeku/grav-lsm/internal/database/seed"
	"github.com/ooyeku/grav-lsm/internal/orm"
	"github.com/ooyeku/grav-lsm/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

// dbManager is a pointer to an instance of the DBLifecycleManager struct.
// It represents a manager for the lifecycle of a database and is responsible
// for setting environment variables, checking file existence, running commands,
// building and starting a Docker container, stopping and removing the container,
// and getting the status of the container. The dbManager variable is of type *lsm.DBLifecycleManager.
var dbManager *lsm.DBLifecycleManager

// log is a variable of type *logrus.Logger in the logrus package.
// It is used for logging messages, errors, and informational data throughout the application.
// The variable is initialized with a new instance of logrus.Logger using the New() function.
// It is commonly used in various command functions to log messages about the status and progress of each command.
var log = logrus.New()

// cfg is a pointer to a config.Config object that holds the configuration for our program.
var cfg *config.Config

// init initializes the program by setting up the logging formatter, output, and log level.
// It loads the program configuration using the config.LoadConfig function and assigns it to the cfg variable.
// If an error occurs during the configuration loading, it logs the error and returns.
// It creates a new instance of the DBLifecycleManager using the lsm.NewDBLifecycleManager function
// and assigns it to the dbManager variable.
func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}
	dbManager = lsm.NewDBLifecycleManager(cfg)
}

// dbCmd is a variable of type *cobra.Command that represents the "db" command.
// It is used to manage the lifecycle of the database.
// It has subcommands for building, starting, stopping, removing, checking the status,
// seeding initial data, running migrations, rolling back migrations, and listing tables.
// The dbCmd variable should be added to the RootCmd command.
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the database lifecycle",
}

// buildCmd is a variable of type *cobra.Command that represents the "build" command.
// It is used to build the database Docker image.
// The command has a Run function that calls the BuildImage method of the dbManager object.
// If the build process fails, an error is logged. Otherwise, a success message is logged.
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

// startCmd is a Command variable that represents the "start" command.
// It is used to start the database Docker container.
// When executed, it calls the StartContainer method of the dbManager instance.
// If an error occurs, it logs an error message.
// If the container starts successfully, it logs a success message.
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

// stopCmd is a variable of type *cobra.Command that represents the "stop" command of a database Docker container.
// It stops the container by running the command "docker stop gravorm-db". If the container is stopped successfully,
// it logs a success message. If an error occurs during the stopping process, it logs the error and returns it.
// The stopCmd is part of a DBLifecycleManager instance used to manage the lifecycle of the database.
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

// removeCmd is a Cobra command that removes the database Docker container.
// It executes the RemoveContainer() method of the DBLifecycleManager type.
// If the removal process fails, an error is returned. Otherwise, a success message is logged.
// This command is part of the CLI tool's functionality for managing the lifecycle of a database.
// It does not have any additional options or arguments.
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

// statusCmd represents the status command.
// It checks the health and status of the database.
// If the status indicates that the container is running,
// it connects to the database and retrieves various metrics,
// such as number of tables, database size, active connections, uptime,
// transactions, cache hit ratio, and slow queries.
// It logs the metrics information using the logrus package.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the health and status of the database",
	Run: func(cmd *cobra.Command, args []string) {
		status, err := dbManager.GetStatus()
		if err != nil {
			log.WithError(err).Error("Error checking database status")
			return
		}

		if strings.Contains(status, "Container is running") {
			conn, err := orm.NewConnection(&cfg.Database)
			if err != nil {
				log.WithError(err).Error("Error connecting to database")
				return
			}
			defer func(conn *orm.Connection) {
				err := conn.Close()
				if err != nil {
					log.WithError(err).Error("Error closing database connection")
				}
			}(conn)

			metrics, err := conn.GetDatabaseMetrics()
			if err != nil {
				log.WithError(err).Error("Error fetching database metrics")
				return
			}

			log.Info("Database Metrics:")
			log.Infof("- Number of tables: %d", metrics.TableCount)
			log.Infof("- Database size: %s", metrics.DatabaseSize)
			log.Infof("- Active connections: %d", metrics.ActiveConnections)
			log.Infof("- Uptime: %s", metrics.Uptime)
			log.Infof("- Transactions (commits/rollbacks): %d/%d", metrics.Commits, metrics.Rollbacks)
			log.Infof("- Cache hit ratio: %.2f%%", metrics.CacheHitRatio)
			log.Infof("- Slow queries (last hour): %d", metrics.SlowQueryCount)
		}
	},
}

// seedCmd represents the "seed" command. It is used to seed the database with initial data. It internally calls the NewConnection() function to establish a connection to the database. Then, it creates a seeder using the NewSeeder() function, loads seed files using the LoadSeeds() function, and finally executes the seeds using the Seed() function. If any error occurs during these operations, it logs an error message.
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := orm.NewConnection(&cfg.Database)
		if err != nil {
			log.WithError(err).Error("Error connecting to database")
			return
		}
		defer func(conn *orm.Connection) {
			err := conn.Close()
			if err != nil {
				log.WithError(err).Error("Error closing database connection")
			}
		}(conn)

		seeder := seed.NewSeeder(conn.GetDB())
		err = seeder.LoadSeeds()
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

// migrateCmd is a command variable of type *cobra.Command. It is used to run database migrations.
// When executed, it connects to the database using a configured ORM connection obtained from orm.NewConnection.
// It then creates a new migrator using migration.NewMigrator, and loads the migrations using migrator.LoadMigrations.
// Finally, it runs the migrations using migrator.Migrate. If any errors occur during these steps, they will be logged.
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := orm.NewConnection(&cfg.Database)
		if err != nil {
			log.WithError(err).Error("Error connecting to database")
			return
		}
		defer func(conn *orm.Connection) {
			err := conn.Close()
			if err != nil {
				log.WithError(err).Error("Error closing database connection")
			}
		}(conn)

		migrator := migration.NewMigrator(conn.GetDB(), log)
		err = migrator.LoadMigrations()
		if err != nil {
			log.WithError(err).Error("Error loading migrations")
			return
		}

		err = migrator.Migrate()
		if err != nil {
			log.WithError(err).Error("Error running migrations")
		} else {
			log.Info("Database migrations completed successfully")
		}
	},
}

// rollbackCmd is a command that performs database rollback. It takes an optional argument "steps" that specifies the number
// of migrations to roll back. If the argument is not provided, it defaults to 1. The command connects to the database by using
// the NewConnection function to establish a connection based on the provided configuration. It creates a migrator instance using
// the NewMigrator function. Then it loads the migrations and rolls back the specified number of steps. If any error occurs,
// it logs the error message. If the rollback is successful, it logs the number of rolled back migrations.
var rollbackCmd = &cobra.Command{
	Use:   "rollback [steps]",
	Short: "Rollback database migrations",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		steps := 1
		if len(args) > 0 {
			var err error
			steps, err = strconv.Atoi(args[0])
			if err != nil {
				log.WithError(err).Error("Invalid number of steps")
				return
			}
		}

		conn, err := orm.NewConnection(&cfg.Database)
		if err != nil {
			log.WithError(err).Error("Error connecting to database")
			return
		}
		defer func(conn *orm.Connection) {
			err := conn.Close()
			if err != nil {
				log.WithError(err).Error("Error closing database connection")
			}
		}(conn)

		migrator := migration.NewMigrator(conn.GetDB(), log)
		err = migrator.LoadMigrations()
		if err != nil {
			log.WithError(err).Error("Error loading migrations")
			return
		}

		err = migrator.Rollback(steps)
		if err != nil {
			log.WithError(err).Error("Error rolling back migrations")
		} else {
			log.Infof("Rolled back %d migration(s) successfully", steps)
		}
	},
}

// listTablesCmd is a variable that defines a Cobra command for listing all tables in the database.
// It has a "Run" function that establishes a database connection using the orm.NewConnection function
// and retrieves a list of tables from the database using the conn.ListTables function.
// If there are no tables found, it logs a message. Otherwise, it logs each table name.
// The listTablesCmd variable is used in the init function to add the command to the dbCmd.
var listTablesCmd = &cobra.Command{
	Use:   "list-tables",
	Short: "List all tables in the database",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := orm.NewConnection(&cfg.Database)
		if err != nil {
			log.WithError(err).Error("Error connecting to database")
			return
		}
		defer func(conn *orm.Connection) {
			err := conn.Close()
			if err != nil {
				log.WithError(err).Error("Error closing database connection")
			}
		}(conn)

		tables, err := conn.ListTables()
		if err != nil {
			log.WithError(err).Error("Error listing tables")
			return
		}

		if len(tables) == 0 {
			log.Info("No tables found in the database")
		} else {
			log.Info("Tables in the database:")
			for _, table := range tables {
				log.Infof("- %s", table)
			}
		}
	},
}

// init is an initialization function that adds subcommands to the dbCmd command and the
// dbCmd command itself to the RootCmd command. It configures the CLI tool with commands
// for managing the lifecycle of a database, including building, starting, stopping, removing,
// and checking the status of a database Docker container, as well as running migrations,
// seeding initial data, and listing tables. These commands are set up by adding them to
// the appropriate variables, such as dbCmd, buildCmd, startCmd, stopCmd, removeCmd,
// statusCmd, seedCmd, migrateCmd, and listTablesCmd. These variables are then used to
// register the commands with the RootCmd command.
func init() {
	dbCmd.AddCommand(buildCmd)
	dbCmd.AddCommand(startCmd)
	dbCmd.AddCommand(stopCmd)
	dbCmd.AddCommand(removeCmd)
	dbCmd.AddCommand(statusCmd)
	dbCmd.AddCommand(seedCmd)
	dbCmd.AddCommand(migrateCmd)
	dbCmd.AddCommand(rollbackCmd)
	dbCmd.AddCommand(listTablesCmd)
	RootCmd.AddCommand(dbCmd)
}
