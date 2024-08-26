local config = {
    Database = {
        Driver = "postgres",
        Host = "localhost",
        Port = 5432,
        User = "postgres",
        Password = "postgres",
        Name = "gravorm",
        SSLMode = "disable"
    },
    Server = {
        Host = "0.0.0.0",
        Port = 8080
    },
    Logging = {
        Level = "info",
        File = "app.log"
    }
}

local function set_env_vars()
    os.execute('export DB_USER=' .. config.Database.User)
    os.execute('export DB_PASSWORD=' .. config.Database.Password)
    os.execute('export DB_NAME=' .. config.Database.Name)
end

local function file_exists(name)
    local f = io.open(name, "r")
    if f then f:close() end
    return f ~= nil
end

local function run_command(command)
    local handle = io.popen(command)
    if handle then
        local result = handle:read("*a")
        handle:close()
        return result
    else
        return nil, "Failed to execute command"
    end
end

function Build_image()
    set_env_vars()
    print("Starting the build process for the database Docker image...")
    if not file_exists("./internal/database/lsm/Dockerfile") then
        error("Dockerfile not found!")
    end
    local build_command = "docker build -f ./internal/database/lsm/Dockerfile -t gravorm-db --build-arg DB_USER=$DB_USER --build-arg DB_PASSWORD=$DB_PASSWORD --build-arg DB_NAME=$DB_NAME ."
    local build_result = os.execute(build_command)
    if build_result == 0 then
        print("Database Docker image built successfully.")
    else
        error("Failed to build the database Docker image.")
    end
end

function Start_container()
    set_env_vars()
    print("Starting the database Docker container...")
    
    -- Check if the container already exists
    local container_exists = os.execute("docker ps -aq -f name=gravorm-db")
    if container_exists then
        print("Container gravorm-db already exists. Removing it...")
        os.execute("docker rm -f gravorm-db")
    end
    
    -- Check if the image exists locally
    local image_exists = os.execute("docker images -q gravorm-db")
    if not image_exists then
        error("Docker image gravorm-db not found. Please build the image first.")
    end
    
    -- Start the Docker container
    local start_command = "docker run -d --name gravorm-db -e POSTGRES_USER=" .. config.Database.User .. " -e POSTGRES_PASSWORD=" .. config.Database.Password .. " -e POSTGRES_DB=" .. config.Database.Name .. " -p 5432:5432 gravorm-db"
    local start_result = os.execute(start_command)
    if start_result == 0 then
        print("Database Docker container started successfully.")
        
        -- Verify the container is running
        local container_running = os.execute("docker ps -q -f name=gravorm-db")
        if not container_running then
            error("Database Docker container is not running.")
        end
        
        -- Verify environment variables inside the container
        local verify_command = "docker exec gravorm-db env | grep POSTGRES"
        local verify_result = os.execute(verify_command)
        if verify_result == 0 then
            print("Environment variables are set correctly in the container.")
        else
            error("Failed to verify environment variables in the container.")
        end
    else
        error("Failed to start the database Docker container.")
    end
end

function Stop_container()
    print("Stopping the database Docker container...")
    local stop_result = os.execute("docker stop gravorm-db")
    if stop_result == 0 then
        print("Database Docker container stopped successfully.")
    else
        error("Failed to stop the database Docker container.")
    end
end

function Remove_container()
    print("Removing the database Docker container...")
    local remove_result = os.execute("docker rm gravorm-db")
    if remove_result == 0 then
        print("Database Docker container removed successfully.")
    else
        error("Failed to remove the database Docker container.")
    end
end

-- Mark unused functions to avoid warnings
local _ = set_env_vars
local _ = file_exists
local _ = run_command
local _ = Build_image
local _ = Start_container
local _ = Stop_container
local _ = Remove_container

return {
    Build_image = Build_image,
    Start_container = Start_container,
    Stop_container = Stop_container,
    Remove_container = Remove_container
}