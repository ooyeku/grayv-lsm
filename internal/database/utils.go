package database

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ooyeku/grav-orm/pkg/config"
)

func PromptDatabaseConfig() config.DatabaseConfig {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter database driver (postgres/sqlite): ")
	driver, _ := reader.ReadString('\n')
	driver = strings.TrimSpace(driver)

	fmt.Print("Enter database host: ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)

	fmt.Print("Enter database port: ")
	portStr, _ := reader.ReadString('\n')
	port, _ := strconv.Atoi(strings.TrimSpace(portStr))

	fmt.Print("Enter database user: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print("Enter database password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Enter database name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter SSL mode (disable/enable): ")
	sslMode, _ := reader.ReadString('\n')
	sslMode = strings.TrimSpace(sslMode)

	return config.DatabaseConfig{
		Driver:   driver,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslMode,
	}
}
