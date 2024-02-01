package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

type DatabaseConfig struct {
	Host   string
	Port   uint16
	DBName string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("User is ?")
	ok := scanner.Scan()
	if !ok {
		fmt.Fprintf(os.Stderr, "Unable read username:\n")
		os.Exit(1)
	}
	username := scanner.Text()
	fmt.Println("User is " + username)

	fmt.Println("Password is ?")
	ok = scanner.Scan()
	if !ok {
		fmt.Fprintf(os.Stderr, "Unable read user password:\n")
		os.Exit(1)
	}
	password := scanner.Text()
	fmt.Println("User password is " + password)

	connConfig, err := ReadParse("connecting.conf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading connect.conf: %v\n", err)
		os.Exit(1)
	}

	urlExample := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, connConfig.Host, connConfig.Port, connConfig.DBName)
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var version string
	err = conn.QueryRow(context.Background(), "select VERSION()").Scan(&version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(version)
}
func ReadParse(filePath string) (DatabaseConfig, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	uri := strings.TrimSpace(string(fileContent))

	connConfig, err := pgx.ParseConfig(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URI: %v\n", err)
		os.Exit(1)
	}

	dbConfig := DatabaseConfig{
		Host:   connConfig.Host,
		Port:   connConfig.Port,
		DBName: connConfig.Database,
	}

	return dbConfig, err
}
