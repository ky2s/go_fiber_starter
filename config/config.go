package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jameskeane/bcrypt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper

	errorHandler fiber.ErrorHandler
	fiber        *fiber.Config
}

var AppConfig *Config

var defaultErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Set error message
	// message := err.Error()

	// Check if it's a fiber.Error type
	// if e, ok := err.(*fiber.Error); ok {
	// 	code = e.Code
	// 	message = e.Message
	// }

	// TODO: Check return type for the client, JSON, HTML, YAML or any other (API vs web)

	// Return HTTP response
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	c.Status(code)

	// Render default error view
	return err
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.SetConfigName(".env")
	config.SetConfigType("dotenv")
	config.AddConfigPath(".")

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	config.SetErrorHandler(defaultErrorHandler)

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file

	// Set Fiber configurations
	config.setFiberConfig()

	AppConfig = config

	return config
}

func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
	config.errorHandler = errorHandler
}

func (config *Config) ConnectDB() (db *sqlx.DB, err error) {

	dbDriver := config.GetString("DB_CONNECTION")
	dbHost := config.GetString("DB_HOST")
	dbUser := config.GetString("DB_USERNAME")
	dbPass := config.GetString("DB_PASSWORD")
	dbPort := config.GetString("DB_PORT")
	dbName := config.GetString("DB_DATABASE")

	dbTimeout := string(config.GetString("DB_TIMEOUT"))
	dbMaxConn, err := strconv.Atoi(config.GetString("DB_MAX_CONN"))
	if err != nil {
		fmt.Println("error conv max conn :", err.Error())
		return nil, err
	}

	dbMaxIdleConn, err := strconv.Atoi(config.GetString("DB_MAX_IDLE_CONN"))
	if err != nil {
		fmt.Println("error conv max idle conn :", err.Error())
		return nil, err
	}

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?timeout=" + dbTimeout + "s&parseTime=True&loc=Asia%2FJakarta"

	db, err = sqlx.Open(dbDriver, dsn)
	if err != nil {
		fmt.Println("error conn :", err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("error conn ping", err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(dbMaxConn)
	db.SetMaxIdleConns(dbMaxIdleConn)

	return db, nil
}

func (config *Config) setDefaults() {
	// Set default App configuration
	config.SetDefault("APP_URL", config.GetString("APP_URL"))
	config.SetDefault("SERVER_PORT", config.GetString("SERVER_PORT"))
	config.SetDefault("APP_ENV", config.GetString("APP_ENV"))

	// Set default database configuration
	config.SetDefault("DB_CONNECTION", config.GetString("DB_CONNECTION"))
	config.SetDefault("DB_HOST", config.GetString("DB_HOST"))
	config.SetDefault("DB_USERNAME", config.GetString("DB_USERNAME"))
	config.SetDefault("DB_PASSWORD", config.GetString("DB_PASSWORD"))
	config.SetDefault("DB_PORT", config.GetString("DB_PORT"))
	config.SetDefault("DB_DATABASE", config.GetString("DB_DATABASE"))
	config.SetDefault("DB_MAX_CONN", config.GetString("DB_MAX_CONN"))
	config.SetDefault("DB_TIMEOUT", config.GetString("DB_TIMEOUT"))
	config.SetDefault("DB_MAX_IDLE_CONN", config.GetString("DB_MAX_IDLE_CONN"))
	config.SetDefault("MAX_CHUNK_HEARTBEAT", 12000)
	config.SetDefault("HOUSEKEEPING_HEARTBEAT", 3)

	// Set default hasher configuration
	config.SetDefault("HASHER_DRIVER", "argon2id")
	config.SetDefault("HASHER_MEMORY", 131072)
	config.SetDefault("HASHER_ITERATIONS", 4)
	config.SetDefault("HASHER_PARALLELISM", 4)
	config.SetDefault("HASHER_SALTLENGTH", 16)
	config.SetDefault("HASHER_KEYLENGTH", 32)
	config.SetDefault("HASHER_ROUNDS", bcrypt.DefaultRounds)

	// Set default session configuration
	// config.SetDefault("SESSION_PROVIDER", "mysql")
	// config.SetDefault("SESSION_KEYPREFIX", "session")
	// config.SetDefault("SESSION_HOST", "localhost")
	// config.SetDefault("SESSION_PORT", 3306)
	// config.SetDefault("SESSION_USERNAME", "fiber")
	// config.SetDefault("SESSION_PASSWORD", "secret")
	// config.SetDefault("SESSION_DATABASE", "boilerplate")
	// config.SetDefault("SESSION_TABLENAME", "sessions")
	// config.SetDefault("SESSION_LOOKUP", "cookie:session_id")
	// config.SetDefault("SESSION_DOMAIN", "")
	// config.SetDefault("SESSION_SAMESITE", "Lax")
	// config.SetDefault("SESSION_EXPIRATION", "12h")
	// config.SetDefault("SESSION_SECURE", false)
	// config.SetDefault("SESSION_GCINTERVAL", "1m")

	// Set default Fiber configuration
	config.SetDefault("FIBER_PREFORK", false)
	config.SetDefault("FIBER_SERVERHEADER", "")
	config.SetDefault("FIBER_STRICTROUTING", false)
	config.SetDefault("FIBER_CASESENSITIVE", false)
	config.SetDefault("FIBER_IMMUTABLE", false)
	config.SetDefault("FIBER_UNESCAPEPATH", false)
	config.SetDefault("FIBER_ETAG", false)
	config.SetDefault("FIBER_BODYLIMIT", 10*1024*1024)
	config.SetDefault("FIBER_CONCURRENCY", 262144)
	config.SetDefault("FIBER_VIEWS", "html")
	config.SetDefault("FIBER_VIEWS_DIRECTORY", "resources/views")
	config.SetDefault("FIBER_VIEWS_RELOAD", false)
	config.SetDefault("FIBER_VIEWS_DEBUG", false)
	config.SetDefault("FIBER_VIEWS_LAYOUT", "embed")
	config.SetDefault("FIBER_VIEWS_DELIMS_L", "{{")
	config.SetDefault("FIBER_VIEWS_DELIMS_R", "}}")
	config.SetDefault("FIBER_READTIMEOUT", 0)
	config.SetDefault("FIBER_WRITETIMEOUT", 0)
	config.SetDefault("FIBER_IDLETIMEOUT", 0)
	config.SetDefault("FIBER_READBUFFERSIZE", 4096)
	config.SetDefault("FIBER_WRITEBUFFERSIZE", 4096)
	config.SetDefault("FIBER_COMPRESSEDFILESUFFIX", ".fiber.gz")
	config.SetDefault("FIBER_PROXYHEADER", "")
	config.SetDefault("FIBER_GETONLY", false)
	config.SetDefault("FIBER_DISABLEKEEPALIVE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTDATE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTCONTENTTYPE", false)
	config.SetDefault("FIBER_DISABLEHEADERNORMALIZING", false)
	config.SetDefault("FIBER_DISABLESTARTUPMESSAGE", false)
	config.SetDefault("FIBER_REDUCEMEMORYUSAGE", false)

	// Set default Custom Access Logger middleware configuration
	// config.SetDefault("MW_ACCESS_LOGGER_ENABLED", true)
	// config.SetDefault("MW_ACCESS_LOGGER_TYPE", "console")
	// config.SetDefault("MW_ACCESS_LOGGER_FILENAME", "access.log")
	// config.SetDefault("MW_ACCESS_LOGGER_MAXSIZE", 500)
	// config.SetDefault("MW_ACCESS_LOGGER_MAXAGE", 28)
	// config.SetDefault("MW_ACCESS_LOGGER_MAXBACKUPS", 3)
	// config.SetDefault("MW_ACCESS_LOGGER_LOCALTIME", false)
	// config.SetDefault("MW_ACCESS_LOGGER_COMPRESS", false)

	// Set default Force HTTPS middleware configuration
	config.SetDefault("MW_FORCE_HTTPS_ENABLED", false)

	// Set default Force trailing slash middleware configuration
	config.SetDefault("MW_FORCE_TRAILING_SLASH_ENABLED", false)

	// Set default HSTS middleware configuration
	config.SetDefault("MW_HSTS_ENABLED", false)
	config.SetDefault("MW_HSTS_MAXAGE", 31536000)
	config.SetDefault("MW_HSTS_INCLUDESUBDOMAINS", true)
	config.SetDefault("MW_HSTS_PRELOAD", false)

	// Set default Suppress WWW middleware configuration
	config.SetDefault("MW_SUPPRESS_WWW_ENABLED", true)

	// Set default Fiber Cache middleware configuration
	config.SetDefault("MW_FIBER_CACHE_ENABLED", false)
	config.SetDefault("MW_FIBER_CACHE_EXPIRATION", "1m")
	config.SetDefault("MW_FIBER_CACHE_CACHECONTROL", false)

	// Set default Fiber Compress middleware configuration
	config.SetDefault("MW_FIBER_COMPRESS_ENABLED", false)
	config.SetDefault("MW_FIBER_COMPRESS_LEVEL", 0)

	// Set default Fiber CORS middleware configuration
	config.SetDefault("MW_FIBER_CORS_ENABLED", true)
	config.SetDefault("MW_FIBER_CORS_ALLOWORIGINS", "*")
	config.SetDefault("MW_FIBER_CORS_ALLOWMETHODS", "GET,POST,HEAD,PUT,DELETE,PATCH")
	config.SetDefault("MW_FIBER_CORS_ALLOWHEADERS", "Origin, Content-Type, Accept")
	// config.SetDefault("MW_FIBER_CORS_ALLOWCREDENTIALS", false)
	config.SetDefault("MW_FIBER_CORS_ALLOWCREDENTIALS", true)
	config.SetDefault("MW_FIBER_CORS_EXPOSEHEADERS", "")
	config.SetDefault("MW_FIBER_CORS_MAXAGE", 0)

	// Set default Fiber CSRF middleware configuration
	config.SetDefault("MW_FIBER_CSRF_ENABLED", false)
	config.SetDefault("MW_FIBER_CSRF_TOKENLOOKUP", "header:X-CSRF-Token")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_NAME", "_csrf")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_SAMESITE", "Strict")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_EXPIRES", "24h")
	config.SetDefault("MW_FIBER_CSRF_CONTEXTKEY", "csrf")

	// Set default Fiber ETag middleware configuration
	config.SetDefault("MW_FIBER_ETAG_ENABLED", false)
	config.SetDefault("MW_FIBER_ETAG_WEAK", false)

	// Set default Fiber Expvar middleware configuration
	config.SetDefault("MW_FIBER_EXPVAR_ENABLED", false)

	// Set default Fiber Favicon middleware configuration
	config.SetDefault("MW_FIBER_FAVICON_ENABLED", false)
	config.SetDefault("MW_FIBER_FAVICON_FILE", "")
	config.SetDefault("MW_FIBER_FAVICON_CACHECONTROL", "public, max-age=31536000")

	// Set default Fiber Limiter middleware configuration
	config.SetDefault("MW_FIBER_LIMITER_ENABLED", false)
	config.SetDefault("MW_FIBER_LIMITER_MAX", 5)
	config.SetDefault("MW_FIBER_LIMITER_DURATION", "1m")

	// Set default Fiber Monitor middleware configuration
	config.SetDefault("MW_FIBER_MONITOR_ENABLED", false)

	// Set default Fiber Pprof middleware configuration
	config.SetDefault("MW_FIBER_PPROF_ENABLED", false)

	// Set default Fiber Recover middleware configuration
	config.SetDefault("MW_FIBER_RECOVER_ENABLED", true)

	// Set default Fiber RequestID middleware configuration
	config.SetDefault("MW_FIBER_REQUESTID_ENABLED", false)
	config.SetDefault("MW_FIBER_REQUESTID_HEADER", "X-Request-ID")
	config.SetDefault("MW_FIBER_REQUESTID_CONTEXTKEY", "requestid")
}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		Prefork:                   config.GetBool("FIBER_PREFORK"),
		ServerHeader:              config.GetString("FIBER_SERVERHEADER"),
		StrictRouting:             config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive:             config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:                 config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:              config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:                      config.GetBool("FIBER_ETAG"),
		BodyLimit:                 config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:               config.GetInt("FIBER_CONCURRENCY"),
		ReadTimeout:               config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:              config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:               config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:            config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:           config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:      config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:               config.GetString("FIBER_PROXYHEADER"),
		GETOnly:                   config.GetBool("FIBER_GETONLY"),
		ErrorHandler:              config.errorHandler,
		DisableKeepalive:          config.GetBool("FIBER_DISABLEKEEPALIVE"),
		DisableDefaultDate:        config.GetBool("FIBER_DISABLEDEFAULTDATE"),
		DisableDefaultContentType: config.GetBool("FIBER_DISABLEDEFAULTCONTENTTYPE"),
		DisableHeaderNormalizing:  config.GetBool("FIBER_DISABLEHEADERNORMALIZING"),
		DisableStartupMessage:     config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:         config.GetBool("FIBER_REDUCEMEMORYUSAGE"),
		JSONEncoder:               json.Marshal,
		JSONDecoder:               json.Unmarshal,
	}
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}
