package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	// please, do not define constants like this in production
	DbHost     = "db"
	DbUser     = "postgres-dev"
	DbPassword = "mysecretpassword"
	DbName     = "dev"
	Migration  = `CREATE TABLE IF NOT EXISTS ttodos (
id serial PRIMARY KEY,
title text NOT NULL,
created_at timestamp with time zone DEFAULT current_timestamp)`
)

// board's bulletin
type Todo struct {
	Id string `json:"id"`
	Title    string    `json:"title" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// global database connection
var db *sql.DB

func getAllTodo() ([]Todo, error) {
	const q = `SELECT id,title, created_at FROM ttodos ORDER BY created_at DESC LIMIT 100`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	results := make([]Todo, 0)

	for rows.Next() {
		var id string
		var title string
		var createAt time.Time
		// scanning the data from the returned rows
		err = rows.Scan(&id,&title, &createAt)
		if err != nil {
			return nil, err
		}
		// creating a new result
		results = append(results, Todo{id,title, createAt})
	}

	return results, nil
}

func addTodo(todo Todo) error {
	const q = `INSERT INTO ttodos(title, created_at) VALUES ($1, $2)`
	_, err := db.Exec(q, todo.Title, todo.CreatedAt)
	return err
}

func main() {
	var err error
	// create a router with a default configuration
	r := gin.Default()
	// endpoint to retrieve all posted bulletins
	r.GET("/api/getall", func(context *gin.Context) {
		results, err := getAllTodo()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
			return
		}
		context.JSON(http.StatusOK, results)
	})
	// endpoint to create a new bulletin
	r.POST("/api/add", func(context *gin.Context) {
		var t Todo
		// reading the request's body & parsing the json
		if context.Bind(&t) == nil {
			t.CreatedAt = time.Now()
			if err := addTodo(t); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
				return
			}
			context.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}
		// if binding was not successful, return an error
		context.JSON(http.StatusUnprocessableEntity, gin.H{"status": "invalid body"})
	})
	// open a connection to the database
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	// do not forget to close the connection
	defer db.Close()
	// ensuring the table is created
	_, err = db.Query(Migration)
	if err != nil {
		log.Println("failed to run migrations", err.Error())
		return
	}
	// running the http server
	log.Println("running..")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}