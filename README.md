# Grayv LSM (Lifecycle Management)
   Grayv LSM is a cli tool for managing the lifecycle of small backend components.  Quickly create a Postgres database, create models, generate code, and run migrations.  This is currently useful for quickly creating a new backend for a static site, or for use in prototyping.  While useful on its own, it will be accompanied by a set of other tools encompasing the Grayv system.

## Features

- **App Management**: Create, list, and delete Grav apps with ease.
- **Database Management**: Build, start, stop, and manage PostgreSQL databases using Docker.
- **Model Management**: Create, update, and generate Go code for data models.
- **Migrations and Seeding**: Manage database schema changes and seed initial data.
- **ORM Integration**: Simplified database operations with built-in ORM functionality.

## Installation

Ensure you have Go installed on your system, then run:

```bash
go install github.com/ooyeku/grayv-lsm@latest
```

Ensure $GOPATH/bin is in your $PATH.
```bash
   export PATH=$PATH:$(go env GOPATH)/bin
```

## Quick Start

1. Create a new Grav app:
   ```bash
   grayv-lsm app create myapp
   ```

2. Build and start the database:
   ```bash
   grayv-lsm db build
   grayv-lsm db start
   ```

3. Create a model:
   ```bash
   grayv-lsm model create User --fields "name:string,email:string,age:int"
   ```

4. Generate model code:
   ```bash
   grayv-lsm model generate User --app myapp
   ```

5. Run migrations:
   ```bash
   grayv-lsm db migrate
   ```

## Configuration

Grav LSM uses a `config.json` file for configuration. You can create this file in the same directory as the executable to override default settings:

```json
{
    "Database": {
        "Driver": "postgres",
        "Host": "localhost",
        "Port": 5432,
        "User": "postgres",
        "Password": "postgres",
        "Name": "grayv-db",
        "SSLMode": "disable"
    }
}
```

## Usage

Run `grav-lsm --help` to see all available commands. Here are some common operations:

- Manage apps: `grav-lsm app [create|list|delete]`
- Manage database: `grav-lsm db [build|start|stop|remove|status]`
- Manage models: `grav-lsm model [create|update|list|generate]`
- Database operations: `grav-lsm db [migrate|rollback|seed]`

For detailed usage instructions, please refer to the [User Guide](docs/user-guide.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.