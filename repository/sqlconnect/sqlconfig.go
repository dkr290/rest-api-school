// Package sqlconnect - for connecting and testing connection to the mariadb database server
package sqlconnect

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/config"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(conf *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		conf.DBUser,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
		conf.DBDatabase,
	)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	var count int
	for {

		err := db.Ping()
		if err != nil {
			fmt.Println("Error connection to mariadb", err)
			if count > 10 {
				return nil, fmt.Errorf("error connection for more then 10 times %v", err)
			}
			count++
			time.Sleep(2 * time.Second)
			continue
		} else {
			fmt.Println("Connected to mariadb")
			break
		}

	}
	return db, nil
}

func CreateTables() {
	createExecTable := `
    CREATE TABLE IF NOT EXISTS execs (
      id INT AUTO_INCREMENT PRIMARY KEY,
	    first_name VARCHAR(255) NOT NULL,
	    last_name VARCHAR(255) NOT NULL,
	    email VARCHAR(50) NOT NULL UNIQUE,
	    username VARCHAR(50) NOT NULL UNIQUE,
	    password VARCHAR(255) NOT NULL,
	    password_changed_at VARCHAR(255),
	    user_created_at TIMESTQAMP DEFAULT CURRENT_TIMESTAMP,
	    password_reset_token VARCHAR(255),
	    inactive_status BOOLEAN NOT NULL,
	    role VARCHAR(255) NOT NULL,
	    INDEX idx_email (email),
	    INDEX idx_username (username)
	);
	`
	createTeachersTable := `
    CREATE TABLE IF NOT EXISTS teachers (
      id INT AUTO_INCREMENT PRIMARY KEY,
	    first_name VARCHAR(255) NOT NULL,
	    last_name VARCHAR(255) NOT NULL,
	    email VARCHAR(255) NOT NULL UNIQUE,
	    class VARCHAR(255) NOT NULL,
	    subject VARCHAR(255) NOT NULL,
	    INDEX (email)
	) AUTO_INCREMENT=100;
	`
	createStudentsTable := `
   CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
	  first_name VARCHAR(255) NOT NULL,
	  last_name VARCHAR(255) NOT NULL,
	  email   VARCHAR(255) NOT NULL UNIQUE,
	  class VARCHAR(255) OT NULL,
	  INDEX (email),
	  FOREIGN KEY (class) REFERENCES teachers(class)
	)AUTO_INCREMENT=100;
	`
}
