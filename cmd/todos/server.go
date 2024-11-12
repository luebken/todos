package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/segmentio/kafka-go"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
)

// Add a helper function for Kafka connectivity check
func checkKafkaConnection(writer *kafka.Writer) error {
	// Try to write a test message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("health-check"),
		Value: []byte("ping"),
	})
	return err
}

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

	// Health check endpoints
	app.Get("/health/live", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/health/ready", func(c *fiber.Ctx) error {
		// Check database connection
		if err := db.PingContext(context.Background()); err != nil {
			log.Printf("Database health check failed: %v", err)
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}

		// Check Kafka connection
		if err := checkKafkaConnection(kafkaWriter); err != nil {
			log.Printf("Kafka health check failed: %v", err)
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	// websocket middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// routes
	app.Static("/", "./public")
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
	app.Get("/ws/logs", websocket.New(logsWebSocketHandler))

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

// Update the logsHandler to serve the logs page instead
func logsHandler(c *fiber.Ctx) error {
	return c.Render("logs", fiber.Map{})
}

// handleLogsWebSocket handles the WebSocket connection for streaming logs
func logsWebSocketHandler(c *websocket.Conn) {
	// Create reader with latest offset
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "todos-topic",
		GroupID: "log-consumer-group",
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		if err := c.WriteJSON(fiber.Map{
			"message": string(m.Value),
			"time":    time.Now().Format(time.RFC3339),
		}); err != nil {
			log.Printf("Error writing to websocket: %v", err)
			break
		}
	}
}
