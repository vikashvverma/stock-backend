package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/vikashvverma/stock-backend/constants"
)

// Config holds the application configuration
type Config struct {
	appPort int

	dbUsername   string
	dbPassword   string
	dbServer     string
	dbPORT       int
	dbConnection string

	APIKey string

	logPath  string
	logFile  io.Writer
	logLevel int

	stock string
	data  string
}

type args struct {
	AppPort string `json:"appPort"`

	DBUsername string `json:"dbUsername"`
	DBPassword string `json:"dbPassword"`
	DBServer   string `json:"dbServer"`
	DBPort     string `json:"dbPort"`

	APIKey string `json:"apiKey"`

	LogPath  string `json:"logPath"`
	LogLevel string `json:"logLevel"`

	Stock string `json:"stock"`
	Data  string `json:"data"`
}

// New creates application configuration from the given args
func New(a *args) (*Config, error) {
	err := validate(a)
	if err != nil {
		return nil, fmt.Errorf("config initialization failed: %s", err)
	}

	appPort, err := strconv.Atoi(a.AppPort)
	if err != nil {
		return nil, fmt.Errorf("invalid value %q supplied for appPort: %s", a.AppPort, err)
	}

	dbPort, err := strconv.Atoi(a.DBPort)
	if err != nil {
		return nil, fmt.Errorf("invalid value %q supplied for dbPort: %s", a.DBPort, err)
	}

	connectionString := fmt.Sprintf("%s://%s:%s/%s",
		constants.DBTypeMongo,
		a.DBServer,
		a.DBPort,
		constants.Database,
	)

	var dbConnection string
	u, err := url.Parse(connectionString)
	if err == nil {
		u.User = url.UserPassword(a.DBUsername, a.DBPassword)
		dbConnection = u.String()
	}

	c := Config{
		appPort:      appPort,
		APIKey:       a.APIKey,
		dbUsername:   a.DBUsername,
		dbPassword:   a.DBPassword,
		dbServer:     a.DBServer,
		dbPORT:       dbPort,
		dbConnection: dbConnection,
		logPath:      a.LogPath,
		logFile:      logFile(a.LogPath, "stock.log"),
		logLevel:     parseLevel(a.LogLevel),
		data:         a.Data,
		stock:        a.Stock,
	}

	return &c, nil
}

// FromFile reads Config from a JSON file.
func FromFile(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file '%s': %s", path, err)
	}

	args := &args{}
	err = json.Unmarshal(content, args)
	if err != nil {
		return nil, fmt.Errorf("config file not valid: %s", err)
	}

	return New(args)
}

// FromFlags reads Config from a flags.
func FromFlags(cmdArgs []string) (*Config, error) {
	flagSet := flag.NewFlagSet(cmdArgs[0], flag.ContinueOnError)
	a := &args{}

	flagSet.StringVar(&a.AppPort, "app_port", "9000", "Application Port")
	flagSet.StringVar(&a.APIKey, "api_key", "", "API Key")
	flagSet.StringVar(&a.DBUsername, "db_username", "", "DB Username")
	flagSet.StringVar(&a.DBPassword, "db_password", "", "DB Password")
	flagSet.StringVar(&a.DBServer, "db_server", "", "DB Server")
	flagSet.StringVar(&a.DBPort, "db_port", "27017", "DB Port")
	flagSet.StringVar(&a.LogPath, "log_path", "", "Log Path")
	flagSet.StringVar(&a.LogLevel, "log_level", "info", "Log Level")
	flagSet.StringVar(&a.Stock, "seating", "data/stock.csv", "Stock csv")
	flagSet.StringVar(&a.Data, "data", "data/data.csv", "data csv")

	err := flagSet.Parse(cmdArgs[1:])

	if err != nil {
		return nil, err
	}

	return New(a)
}

// AppPort for the service to listen to.
func (config Config) AppPort() int {
	return config.appPort
}

// DBConnection for the database.
func (config Config) DBConnection() string {
	return config.dbConnection
}

// LogLevel returns log level for the application.
func (config Config) LogLevel() int {
	return config.logLevel
}

// LogFile returns file where the log should be logged.
func (config Config) LogFile() io.Writer {
	return config.logFile
}

// Stock config for the table.
func (config Config) Stock() string {
	return config.stock
}

// Data config for the table.
func (config Config) Data() string {
	return config.data
}

func validate(a *args) error {
	if a == nil {
		return fmt.Errorf("empty args supplied")
	}

	var missing []string

	if a.AppPort == "" {
		missing = append(missing, "appPort")
	}

	if a.DBServer == "" {
		missing = append(missing, "dbServer")
	}

	if a.DBPort == "" {
		missing = append(missing, "dbPort")
	}

	if len(missing) > 0 {
		return fmt.Errorf("%s not found", strings.Join(missing, ", "))
	}

	return nil
}

func logFile(path, name string) *os.File {
	if path == "" {
		return os.Stdout
	}

	file, err := os.OpenFile(fmt.Sprintf("%s%s", path, name), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("logFile: failed to create log to file, using default stdout %s", err)
		return os.Stdout
	}

	return file
}

func parseLevel(level string) int {
	switch strings.ToLower(level) {
	case "error":
		return 2
	case "warn", "warning":
		return 3
	case "info":
		return 4
	case "debug":
		return 5
	default:
		return 4
	}
}
