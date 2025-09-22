package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("database_url is not set")
	}

	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return conn, nil
}

// Template from neondb site

// err := godotenv.Load("../../backend-go.env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	connStr := os.Getenv("DATABASE_URL")
// 	conn, err := pgx.Connect(context.Background(), connStr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close(context.Background())
// 	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS playing_with_neon(id SERIAL PRIMARY KEY, name TEXT NOT NULL, value REAL);")
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, err = conn.Exec(context.Background(), "INSERT INTO playing_with_neon(name, value) SELECT LEFT(md5(i::TEXT), 10), random() FROM generate_series(1, 10) s(i);")
// 	if err != nil {
// 		panic(err)
// 	}
// 	rows, err := conn.Query(context.Background(), "SELECT * FROM playing_with_neon")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int32
// 		var name string
// 		var value float32
// 		if err := rows.Scan(&id, &name, &value); err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("%d | %s | %f\n", id, name, value)
// 	}
