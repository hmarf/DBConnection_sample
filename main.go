package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func initDB() *sql.DB {
	db, err := sql.Open("mysql", "user:password@tcp(0.0.0.0)/sampleDB?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {

	DB := initDB()
	defer DB.Close()

	// connection数 の制限
	DB.SetMaxOpenConns(9)

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		result, err := DB.Exec("INSERT INTO user(name,createdAt) VALUES(?,?)", "aaa", time.Now())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(time.Now())
		w.WriteHeader(200)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}

}
