package database

// Config provides database configuration
type Config struct {
	// Enabled used to run the server without the database, simply here to allow
	// the server to start since this is a structure demo, and not a functional
	// server
	Enabled bool

	// Database driver
	Driver string

	// Database connection string
	ConnectionString string

	// Log will enable or disable query logging
	Log bool

	// Check if there is a custom migrations folder
	MigrationFolder string
}
