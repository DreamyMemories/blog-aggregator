package types

import "github.com/DreamyMemories/blog-aggregator/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
