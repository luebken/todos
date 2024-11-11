package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/template/html/v2"
	"github.com/segmentio/kafka-go"

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

	defer db.Close()

	// Kafka
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "todos-topic",
		Balancer: &kafka.LeastBytes{},
	})

	defer kafkaWriter.Close()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db, kafkaWriter)
	})
	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db, kafkaWriter)
	})
	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db, kafkaWriter)
	})
	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db, kafkaWriter)
	})
	app.Get("/logs", func(c *fiber.Ctx) error {
		return logsHandler(c)
	})
	app.Static("/", "./public")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB, kafkaWriter *kafka.Writer) error {
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

	// Send log to Kafka
	kafkaWriter.WriteMessages(c.Context(), kafka.Message{
		Key:   []byte("get"),
		Value: []byte(fmt.Sprintf("User %s fetched todos at %s", username, time.Now().Format(time.RFC3339))),
	})

	return c.Render("index", fiber.Map{
		"Todos": todos,
	})
}

type todo struct {
	Item     string
	Username string
}

func postHandler(c *fiber.Ctx, db *sql.DB, kafkaWriter *kafka.Writer) error {
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

	// Send log to Kafka
	kafkaWriter.WriteMessages(c.Context(), kafka.Message{
		Key:   []byte("post"),
		Value: []byte(fmt.Sprintf("User %s added todo %s at %s", newTodo.Username, newTodo.Item, time.Now().Format(time.RFC3339))),
	})

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB, kafkaWriter *kafka.Writer) error {
	log.Println("put")

	olditem := c.Query("olditem")
	newitem := c.Query("newitem")
	username := c.Query("username")

	log.Printf("Old item: %v, New item: %v", olditem, newitem)
	_, err := db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)
	if err != nil {
		log.Fatalf("An error occurred while executing query: %v", err)
	}

	// Send log to Kafka
	kafkaWriter.WriteMessages(c.Context(), kafka.Message{
		Key:   []byte("put"),
		Value: []byte(fmt.Sprintf("User %s updated todo from %s to %s at %s", username, olditem, newitem, time.Now().Format(time.RFC3339))),
	})

	return c.Status(fiber.StatusOK).SendString("Item updated")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB, kafkaWriter *kafka.Writer) error {
	log.Println("delete")

	todoToDelete := c.Query("item")
	username := c.Query("username")

	_, err := db.Exec("DELETE from todos WHERE item=$1", todoToDelete)
	if err != nil {
		log.Fatalf("An error occurred while executing query: %v", err)
	}

	// Send log to Kafka
	kafkaWriter.WriteMessages(c.Context(), kafka.Message{
		Key:   []byte("delete"),
		Value: []byte(fmt.Sprintf("User %s deleted todo %s at %s", username, todoToDelete, time.Now().Format(time.RFC3339))),
	})
	return c.SendString("deleted")
}

// TODO change to streaming logs to the ui
func logsHandler(c *fiber.Ctx) error {
	log.Println("logs")

	// Create reader with latest offset
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "todos-topic",
		GroupID: "log-consumer-group",
	})
	defer r.Close()

	// Create timeout context
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	messages := make([]string, 0, 20)

	// Read messages with timeout
	for i := 0; i < 20; i++ {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				log.Println("timeout reached")
				break
			}
			if err == io.EOF {
				log.Println("no more messages")
				break
			}
			log.Printf("error reading message: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to read messages")
		}
		messages = append(messages, string(m.Value))
	}

	// Return empty array if no messages
	if len(messages) == 0 {
		return c.JSON([]string{})
	}

	return c.JSON(messages)
}
