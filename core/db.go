package core

import (
	"database/sql"
	"pixiv-tg-bot/cmd"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 初始化数据库
func init() {
	db, err := sql.Open("sqlite3", cmd.DB_PATH)
	checkErr(err)
	DB = db

	createTable()
}

func createTable() {
	sqlTable := `
CREATE TABLE IF NOT EXISTS "novels" (
   "id" VARCHAR(20) PRIMARY KEY,
   "title" VARCHAR(100) NULL,
   "update_date" VARCHAR(50) NULL
);

CREATE TABLE IF NOT EXISTS "content_titles" (
	"id" VARCHAR(20) PRIMARY KEY,
	"title" VARCHAR(100) NULL,
	"novel_id" VARCHAR(20)
 );
`

	_, err := DB.Exec(sqlTable)
	checkErr(err)
}
