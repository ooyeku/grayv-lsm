# Grav LSM (Lifecycle Management)

Grav LSM is a powerful CLI tool for managing the lifecycle of Grav Apps. It provides a comprehensive solution for creating and managing lightweight backend components, including a containerized database, model/schema generator, and ORM system.

## Features

- **App Management**: Create, list, and delete Grav apps with ease.
- **Database Management**: Build, start, stop, and manage PostgreSQL databases using Docker.
- **Model Management**: Create, update, and generate Go code for data models.
- **Migrations and Seeding**: Manage database schema changes and seed initial data.
- **ORM Integration**: Simplified database operations with built-in ORM functionality.

## Installation

Ensure you have Go installed on your system, then run:

```bash
go install github.com/ooyeku/grav-lsm
```

Ensure $GOPATH/bin is in your $PATH.
```bash
   export PATH=$PATH:$(go env GOPATH)/bin
```

## Quick Start

1. Create a new Grav app:
   ```bash
   grav-lsm app create myapp
   ```

2. Build and start the database:
   ```bash
   grav-lsm db build
   grav-lsm db start
   ```

3. Create a model:
   ```bash
   grav-lsm model create User --fields "name:string,email:string,age:int"
   ```

4. Generate model code:
   ```bash
   grav-lsm model generate User --app myapp
   ```

5. Run migrations:
   ```bash
   grav-lsm db migrate
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
        "Name": "gravorm",
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