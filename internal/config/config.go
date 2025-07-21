package config

// Config represents the application configuration
type Config struct {
	Port          string
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string
	LogLevel      string
	DevMode       bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Port:          "8080",
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
		LogLevel:      "info",
		DevMode:       false,
	}
}
