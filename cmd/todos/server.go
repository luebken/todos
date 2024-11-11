package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/template/html/v2"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Database
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalf("Environment variable DATABASE_URL is not set. Aborting.")
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})
	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})
	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})
	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})
	app.Static("/", "./public")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	log.Println("get")
	username := c.Query("username")

	var item string
	var todos []todo

	if username != "" {
		rows, err := db.Query("SELECT item FROM todos WHERE username=$1", username)
		if err != nil {
			log.Fatalln(err)
			return c.JSON("An error occurred")
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&item)
			todos = append(todos, todo{Item: item, Username: username})
		}
	}

	return c.Render("index", fiber.Map{
		"Todos": todos,
	})
}

type todo struct {
	Item     string
	Username string
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	log.Println("post")

	newTodo := todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	log.Printf("%v", newTodo)
	if newTodo.Item != "" && newTodo.Username != "" {
		_, err := db.Exec("INSERT into todos (item, username) VALUES ($1, $2)", newTodo.Item, newTodo.Username)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	log.Println("put")

	olditem := c.Query("olditem")
	newitem := c.Query("newitem")
	db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)
	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	log.Println("delete")

	todoToDelete := c.Query("item")
	db.Exec("DELETE from todos WHERE item=$1", todoToDelete)
	return c.SendString("deleted")
}
