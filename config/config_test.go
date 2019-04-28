package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	config, err := New(&args{AppPort: "9000", DBUsername: "foo", DBPassword: "bar",
		DBServer: "baz", DBPort: "27017", Stock: "./data/stock.csv", Data: "./data/data.csv"})
	require.NoError(t, err, "Expected no error")
	require.Nil(t, err, "Expected err to be nil")

	expectedConfig := &Config{
		appPort:      9000,
		dbUsername:   "foo",
		dbPassword:   "bar",
		dbServer:     "baz",
		dbPORT:       27017,
		dbConnection: "mongodb://foo:bar@baz:27017/trading",
		logLevel:     4,
		logFile:      os.Stdout,
		stock:        "./data/stock.csv",
		data:         "./data/data.csv",
	}

	assert.Equal(t, expectedConfig, config)
}

func TestNewFailsWhenInvalidAppPort(t *testing.T) {
	config, err := New(&args{AppPort: "SOME_PORT", DBUsername: "foo", DBPassword: "bar",
		DBServer: "baz", DBPort: "5432"})
	require.Nil(t, config, "Expected config to be nil")

	assert.Contains(t, err.Error(), "invalid value \"SOME_PORT\" supplied for appPort:")
}

func TestNewFailsWhenInvalidDBPort(t *testing.T) {
	config, err := New(&args{AppPort: "900", DBUsername: "foo", DBPassword: "bar",
		DBServer: "baz", DBPort: "SOME_PORT"})
	require.Nil(t, config, "Expected config to be nil")

	assert.Contains(t, err.Error(), "invalid value \"SOME_PORT\" supplied for dbPort:")
}

func TestNewFailWhenArgsMissing(t *testing.T) {
	config, err := New(&args{})
	require.Nil(t, config, "Expected config to be nil")

	assert.Equal(t, "config initialization failed: appPort, dbServer, dbPort not found", err.Error())
}

func TestValidateWhenMissingArgs(t *testing.T) {
	err := validate(&args{})
	require.Error(t, err, "Expected an error")

	assert.Equal(t, "appPort, dbServer, dbPort not found", err.Error())
}

func TestValidateWhenNoArgs(t *testing.T) {
	err := validate(nil)
	require.Error(t, err, "Expected an error")

	assert.Equal(t, "empty args supplied", err.Error())
}

func TestNewConfigFromFile(t *testing.T) {
	content := []byte(`
		{
  			"appPort": "9000",
  			"dbUsername": "foo",
  			"dbPassword": "bar",
  			"dbServer": "baz",
  			"dbPort": "5432",
  			"database": "demo",
			"logPath": "/foo/bar/"
		}`)

	tmpFileName := createTemporaryFile(t, content)
	defer os.Remove(tmpFileName)

	config, err := FromFile(tmpFileName)
	require.NoError(t, err, "Expected no error")

	assert.Equal(t, 9000, config.appPort)
	assert.Equal(t, "foo", config.dbUsername)
	assert.Equal(t, "bar", config.dbPassword)
	assert.Equal(t, "baz", config.dbServer)
	assert.Equal(t, 5432, config.dbPORT)
	assert.Equal(t, "mongodb://foo:bar@baz:5432/trading", config.dbConnection)
	assert.Equal(t, "/foo/bar/", config.logPath)
	assert.IsType(t, &os.File{}, config.logFile)
	assert.Equal(t, 4, config.logLevel)
}

func TestLoadUsingFlags(t *testing.T) {
	cmdArgs := []string{"cmd",
		"-app_port=9000",
		"-db_username=foo",
		"-db_password=bar",
		"-db_server=baz",
		"-db_port=5432",
	}

	config, err := FromFlags(cmdArgs)
	require.Nil(t, err, "Expected err to be nil")

	assert.Equal(t, 9000, config.appPort)
	assert.Equal(t, "foo", config.dbUsername)
	assert.Equal(t, "bar", config.dbPassword)
	assert.Equal(t, "baz", config.dbServer)
	assert.Equal(t, 5432, config.dbPORT)
	assert.Equal(t, "mongodb://foo:bar@baz:5432/trading", config.dbConnection)
	assert.Equal(t, "", config.logPath)
	assert.IsType(t, &os.File{}, config.logFile)
	assert.Equal(t, 4, config.logLevel)
}

func TestLoadUsingFlagsFailsWhenExtraFlags(t *testing.T) {
	cmdArgs := []string{"cmd",
		"-app_port=9000",
		"-db_username=foo",
		"-db_password=bar",
		"-db_server=baz",
		"-db_port=5432",
		"-foo=bar",
	}

	config, err := FromFlags(cmdArgs)
	require.Nil(t, config, "Expected config to be nil")

	assert.Equal(t, "flag provided but not defined: -foo", err.Error())

}

func TestNewConfigFromFileWhenFileMissing(t *testing.T) {
	_, err := FromFile("/some/path/which/does/not/exist")
	require.Error(t, err, "Expected error for missing file")

	assert.Contains(t, err.Error(), "unable to open config file '/some/path/which/does/not/exist':")
}

func TestNewConfigFromFileWhenMalformed(t *testing.T) {
	content := []byte("malformed JSON")
	tmpFileName := createTemporaryFile(t, content)
	defer os.Remove(tmpFileName)

	_, err := FromFile(tmpFileName)
	require.Error(t, err, "Expected error for missing file")

	assert.Contains(t, err.Error(), "config file not valid:")
}

func TestAppPort(t *testing.T) {
	c := &Config{appPort: 8000}
	assert.Equal(t, 8000, c.AppPort())
}

func TestDBConnection(t *testing.T) {
	c := &Config{dbConnection: "foo://bar@baz"}
	assert.Equal(t, "foo://bar@baz", c.DBConnection())
}

func TestLogLevel(t *testing.T) {
	c := &Config{logLevel: 2}
	assert.Equal(t, 2, c.LogLevel())
}

func TestStock(t *testing.T) {
	c := &Config{stock: "./data/stock.csv"}
	assert.Equal(t, "./data/stock.csv", c.Stock())
}

func TestData(t *testing.T) {
	c := &Config{data: "./data/data.csv"}
	assert.Equal(t, "./data/data.csv", c.Data())
}

func TestFile(t *testing.T) {
	c := &Config{logFile: os.Stdout}
	assert.Equal(t, os.Stdout, c.LogFile())
}

func TestLogFile(t *testing.T) {
	file := logFile("./", "foo.bar")
	defer os.Remove("./foo.bar")

	assert.Equal(t, "./foo.bar", file.Name())
	assert.IsType(t, &os.File{}, file)
}

func TestParseLevel(t *testing.T) {
	assert.Equal(t, 2, parseLevel("error"))
	assert.Equal(t, 3, parseLevel("warn"))
	assert.Equal(t, 3, parseLevel("warning"))
	assert.Equal(t, 4, parseLevel("info"))
	assert.Equal(t, 4, parseLevel("wrong"))
	assert.Equal(t, 5, parseLevel("debug"))
}

func createTemporaryFile(t *testing.T, content []byte) string {
	name := "config.json"
	tmpFile, err := ioutil.TempFile(".", name)
	require.NoError(t, err, "Expected no error")

	_, err = tmpFile.Write(content)
	require.NoError(t, err, "Expected no error writing to temporary file")

	err = tmpFile.Close()
	require.NoError(t, err, "Expected no error closing file")

	return tmpFile.Name()
}
