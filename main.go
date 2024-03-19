package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	username := os.Getenv("SQLSERVER_USER")
	password := os.Getenv("SQLSERVER_PASS")
	hostname := os.Getenv("SQLSERVER_HOST")
	port := os.Getenv("SQLSERVER_PORT")
	db := os.Getenv("SQLSERVER_DB")
	root := os.Getenv("SQLSERVER_BACKUP_ROOT")

	bakName := fmt.Sprintf("%s-%s.bak", db, time.Now().Format(time.RFC3339))
	destPath := path.Join(root, bakName)

	cmd := fmt.Sprintf(`
USE %s;
GO
BACKUP DATABASE %s
   TO DISK = '%s'
   WITH FORMAT,
      MEDIANAME = 'BackupMedia',
      NAME = 'Full Backup of %s';
`, db, db, destPath, db)

	query := url.Values{}
	query.Add("app name", "db_backup")

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, "fake_password"),
		Host:     fmt.Sprintf("%s:%s", hostname, port),
		RawQuery: query.Encode(),
	}

	log.Printf("Connecting: %v", u)

	u.User = url.UserPassword(username, password)

	conn, err := sql.Open("sqlserver", u.String())
	if err != nil {
		log.Fatal(err)
	}

	clauses := strings.Split(cmd, "GO")

	for _, clause := range clauses {
		log.Printf("Executing:\n%s\n", clause)
		if res, err := conn.Exec(clause); err != nil {
			log.Fatalf("Error executing clause:\n%s\n%+v\n%+v\n", clause, res, err)
		}
	}

	log.Println("Success!")
}
