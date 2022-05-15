package core

import (
	"database/sql"
	"os"
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
	file, err := os.Stat(cmd.DB_PATH)
	if err != nil {
		panic(err)
	}
	if file.IsDir() {
		panic("数据库路径错误!")
	}

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
   "update_date" VARCHAR(50) NULL,
   "content" Text NULL
);
`

	_, err := DB.Exec(sqlTable)
	checkErr(err)
}
