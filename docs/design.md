# Grav-ORM Project Structure

## Main Components

1. CLI Application (cmd/)
    - Implement commands for database lifecycle management
    - Implement commands for ORM operations
---------------------------------------------------------------------------------------------------------------------------
### 'db' cmd package
Example:
```
grav-orm db build
grav-orm db remove
grav-orm db start
grav-orm db status
grav-orm db stop
```

The 'db' cmd package is responsible for managing the database lifecycle. It includes commands for creating, configuring, migrating, and seeding the database.
1. 'build' command: Builds the database instance
2. 'remove' command: Removes the database instance
3. 'start' command: Starts the database instance
4. 'status' command: Checks the status of the database instance
5. 'stop' command: Stops the database instance
---------------------------------------------------------------------------------------------------------------------------
2. Internal Packages (internal/)
    - Database
      - lsm (lifecycle management)
      - migration
      - seeder
    - Model
      - model
      - generator
    - ORM core functionality
        - connection
        - query
        - crud

3. Public Packages (pkg/)
    - Configuration management

4. Main Application (main.go)
    - Entry point for the CLI application

5. Migrations and Seeds
    - The `migrations` directory is used to store database migration files. Migrations are a way to manage changes to your database schema over time.





# Migrations Directory

The `migrations` directory stores database migration files. These files define how your database schema changes over time.

## Purpose:
- Keep track of database schema changes
- Allow for version control of your database structure
- Enable easy deployment and rollback of database changes
- Facilitate collaboration among team members

## Typical Structure:
Migration files are usually named with a timestamp or version number, followed by a descriptive name:

```
migrations/
  ├── 20230101000000_create_users_table.sql
  ├── 20230102000000_add_email_to_users.sql
  ├── 20230103000000_create_posts_table.sql
  └── ...
```

## File Content:
Each migration file typically contains two sections:
1. "Up" migration: Commands to apply the change
2. "Down" migration: Commands to revert the change

Example:
```sql
-- Up migration
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Down migration
DROP TABLE users;
```

Your ORM's migration tool will use these files to manage the database schema.


2. Seeds Directory:

The `seeds` directory is used to store seed files. Seed files contain initial or sample data that you want to populate your database with.





# Seeds Directory

The `seeds` directory stores database seed files. These files contain initial or sample data to populate your database tables.

## Purpose:
- Provide initial data for development and testing
- Ensure consistent starting data across different environments
- Populate lookup tables or other essential data

## Typical Structure:
Seed files are often named after the table they populate or the type of data they insert:

```
seeds/
  ├── 01_users.sql
  ├── 02_categories.sql
  ├── 03_products.sql
  └── ...
```

## File Content:
Seed files typically contain INSERT statements to add data to your tables.

Example:
```sql
INSERT INTO users (name, email) VALUES
('John Doe', 'john@example.com'),
('Jane Smith', 'jane@example.com');

INSERT INTO categories (name) VALUES
('Electronics'),
('Books'),
('Clothing');
```

Your ORM's seeding tool will use these files to populate the database with initial data.


Using migrations and seeds in your ORM system provides several benefits:

1. Versioned Database Schema: Migrations allow you to version your database schema, making it easier to track changes over time and collaborate with team members.

2. Reproducible Environments: With migrations and seeds, you can easily set up identical database structures and initial data across different environments (development, staging, production).

3. Easier Testing: Seeds provide consistent test data, making it easier to write and run tests against your database.

4. Simplified Deployment: When deploying your application, you can run migrations to update the database schema and seeds to ensure necessary data is present.

To use these effectively in your ORM:

1. Implement a migration runner that can apply and rollback migrations from the `migrations` directory.
2. Create a seeder that can run the seed files from the `seeds` directory.
3. Integrate these into your CLI commands, allowing users to easily run migrations and seed data.

For example, you might have CLI commands like:

```
grav-orm migrate up
grav-orm migrate down
grav-orm seed
```

These would run the migrations (forward or backward) and seed the database, respectively.