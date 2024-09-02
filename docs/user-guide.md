# Grayv LSM User Guide

Grayv LSM (Lifecycle Management) is a CLI tool for managing the lifecycle of Grayv Apps. Grayv apps are lightweight backend components consisting of a containerized database, a model/schema generator, and an ORM system.

## Table of Contents

- [Gryav LSM User Guide](#grav-lsm-user-guide)
  - [Table of Contents](#table-of-contents)
  - [1. Installation](#1-installation)
  - [2. Configuration](#2-configuration)
  - [3. Managing Apps](#3-managing-apps)
  - [4. Database Management](#4-database-management)
  - [5. Model Management](#5-model-management)
  - [6. Migrations and Seeding](#6-migrations-and-seeding)
  - [7. ORM Management](#7-orm-management)

## 1. Installation

To install Grayv LSM, follow these steps:

1. Ensure you have Go installed on your system.
2. Run the following command to install Grayv LSM:
   ```
   go install github.com/ooyeku/grayv-lsm@latest
   ```
3. The `grayv-lsm` command should now be available in your terminal.

## 2. Configuration

Grayv LSM uses a configuration file to manage database and server settings. The default configuration is embedded in the application, but you can override it by creating a `config.json` file in the same directory as the executable.

Example `config.json`:

```json
{
    "Database": {
        "Driver": "postgres",
        "Host": "localhost",
        "Port": 5432,
        "User": "postgres",
        "Password": "postgres",
        "Name": "gravorm",
        "SSLMode": "disable"
    }
}
```

Configuration file can also be set using environment variables. The following environment variables are supported:

- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_HOST`
- `DB_PORT`
- `DB_SSLMODE`

Furthermore, the config command can be used to get and set the config values.

```
grayv-lsm config get database.host
grayv-lsm config set database.host 127.0.0.1
```


## 3. Managing Apps

Grayv LSM allows you to create, list, and delete Grav apps.

- Create a new app:
  ```
  grayv-lsm app create myapp
  ```

- List all apps:
  ```
  grayv-lsm app list
  ```

- Delete an app:
  ```
  grayv-lsm app delete myapp
  ```

## 4. Database Management

Grayv LSM provides commands to manage the database lifecycle.

- Build the database Docker image:
  ```
  grayv-lsm db build
  ```

- Start the database container:
  ```
  grayv-lsm db start
  ```

- Stop the database container:
  ```
  grayv-lsm db stop
  ```

- Remove the database container:
  ```
  grayv-lsm db remove
  ```

- Check database status:
  ```
  grayv-lsm db status
  ```

- List database tables:
  ```
  grayv-lsm db list-tables
  ```

## 5. Model Management

Grayv LSM allows you to create, update, and generate models.

- Create a new model:
  ```
  grayv-lsm model create User --fields "name:string,email:string,age:int"
  ```

- Update an existing model:
  ```
  grayv-lsm model update User --add-fields "address:string" --remove-fields "age"
  ```

- List all models:
  ```
  grayv-lsm model list
  ```

- Generate Go code for a model:
  ```
  grayv-lsm model generate User --app myapp
  ```

## 6. Migrations and Seeding

Grayv LSM supports database migrations and seeding.

- Run migrations:
  ```
  grayv-lsm db migrate
  ```

- Rollback migrations:
  ```
  grayv-lsm db rollback [steps]
  ```

- Seed the database:
  ```
  grayv-lsm db seed
  ```

## 7. ORM Management

Grayv LSM allows you to manage the ORM system.

- Create a new user:
  ```
  grayv-lsm orm create-user --username "admin" --email "admin@example.com" --password "admin"
  ```

- List all users:
  ```
  grayv-lsm orm list-users
  ```

- Delete a user:
  ```
  grayv-lsm orm delete-user --id 1
  ```

- Update a user:
  ```
  grayv-lsm orm update-user --id 1 --username "admin123" --email "admin123@example.com"
  ```

- Raw SQL query:
  ```
  grayv-lsm orm query "SELECT * FROM users"
  ```

Remember to run `grayv-lsm --help` or `grayv-lsm [command] --help` for more information on available commands and their usage.
