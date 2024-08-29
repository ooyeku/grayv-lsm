# Grav LSM User Guide

Grav LSM (Lifecycle Management) is a CLI tool for managing the lifecycle of Grav Apps. Grav apps are lightweight backend components consisting of a containerized database, a model/schema generator, and an ORM system.

## Table of Contents

- [Grav LSM User Guide](#grav-lsm-user-guide)
  - [Table of Contents](#table-of-contents)
  - [1. Installation](#1-installation)
  - [2. Configuration](#2-configuration)
  - [3. Managing Apps](#3-managing-apps)
  - [4. Database Management](#4-database-management)
  - [5. Model Management](#5-model-management)
  - [6. Migrations and Seeding](#6-migrations-and-seeding)
  - [7. ORM Management](#7-orm-management)

## 1. Installation

To install Grav LSM, follow these steps:

1. Ensure you have Go installed on your system.
2. Run the following command:
   ```
   go get github.com/ooyeku/grav-lsm
   ```
3. The `grav-lsm` command should now be available in your terminal.

## 2. Configuration

Grav LSM uses a configuration file to manage database and server settings. The default configuration is embedded in the application, but you can override it by creating a `config.json` file in the same directory as the executable.

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

## 3. Managing Apps

Grav LSM allows you to create, list, and delete Grav apps.

- Create a new app:
  ```
  grav-lsm app create myapp
  ```

- List all apps:
  ```
  grav-lsm app list
  ```

- Delete an app:
  ```
  grav-lsm app delete myapp
  ```

## 4. Database Management

Grav LSM provides commands to manage the database lifecycle.

- Build the database Docker image:
  ```
  grav-lsm db build
  ```

- Start the database container:
  ```
  grav-lsm db start
  ```

- Stop the database container:
  ```
  grav-lsm db stop
  ```

- Remove the database container:
  ```
  grav-lsm db remove
  ```

- Check database status:
  ```
  grav-lsm db status
  ```

- List database tables:
  ```
  grav-lsm db list-tables
  ```

## 5. Model Management

Grav LSM allows you to create, update, and generate models.

- Create a new model:
  ```
  grav-lsm model create User --fields "name:string,email:string,age:int"
  ```

- Update an existing model:
  ```
  grav-lsm model update User --add-fields "address:string" --remove-fields "age"
  ```

- List all models:
  ```
  grav-lsm model list
  ```

- Generate Go code for a model:
  ```
  grav-lsm model generate User --app myapp
  ```

## 6. Migrations and Seeding

Grav LSM supports database migrations and seeding.

- Run migrations:
  ```
  grav-lsm db migrate
  ```

- Rollback migrations:
  ```
  grav-lsm db rollback [steps]
  ```

- Seed the database:
  ```
  grav-lsm db seed
  ```

## 7. ORM Management

Grav LSM allows you to manage the ORM system.

- Create a new user:
  ```
  grav-lsm orm create-user --username "admin" --email "admin@example.com" --password "admin"
  ```

- Raw SQL query:
  ```
  grav-lsm orm query "SELECT * FROM users"
  ```

Remember to run `grav-lsm --help` or `grav-lsm [command] --help` for more information on available commands and their usage.