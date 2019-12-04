package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/carlescere/scheduler"
	_ "github.com/go-sql-driver/mysql"
)

type insertData struct {
	name      string
	createdAT time.Time
}

func initDB() *sql.DB {
	db, err := sql.Open("mysql", "user:password@tcp(0.0.0.0)/sampleDB?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func insertDB(ch *chan insertData, db *sql.DB) {

	rescInterface := []interface{}{}
	stmt := "INSERT INTO user(name, createdAt) VALUES"
	insertFlag := false
LOOP:
	for {
		select {
		case data, ok := <-*ch:
			if ok {
				insertFlag = true
				stmt += "(?,?),"
				rescInterface = append(rescInterface, data.name)
				rescInterface = append(rescInterface, data.createdAT)
			}
		default:
			break LOOP
		}
	}

	if insertFlag {
		stmt = strings.TrimRight(stmt, ",")

		_, err := db.Exec(stmt, rescInterface...)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

func main() {

	DB := initDB()
	defer DB.Close()

	// connection数 の制限
	DB.SetMaxOpenConns(9)

	channel := make(chan insertData, 100000)

	_, _ = scheduler.Every(5).Seconds().NotImmediately().Run(func() { insertDB(&channel, DB) })

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		channel <- insertData{
			name:      "test user",
			createdAT: time.Now(),
		}
		w.WriteHeader(200)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}

}
