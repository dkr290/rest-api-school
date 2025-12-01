package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port                       string
	DBUser                     string
	DBPassword                 string
	DBDatabase                 string
	DBPort                     string
	DBHost                     string
	Debug                      bool
	JWTSecret                  string
	JWTExpiresIn               time.Duration
	ExcludedAuthMiddlewarePath []string
	ResetTokenExpDuration      time.Duration
}

func LoadConfig() *Config {
	c := &Config{}
	c.GetFlags()
	return c
}

func (c *Config) GetFlags() {
	var JwtStringExpireValue string
	var exclPaths string
	var resetTokenExpDuration string
	flag.StringVar(
		&c.Port,
		"app-port",
		"8082",
		"Application port",
	)
	flag.StringVar(&c.DBUser, "database-user", "root", "Database user")
	flag.StringVar(&c.DBPassword, "database-password", "password", "Database password")
	flag.StringVar(&c.DBDatabase, "database", "school_test", "Database to use")
	flag.StringVar(&c.DBHost, "database-host", "mariadb", "The database host to use or ip")
	flag.StringVar(
		&c.DBPort,
		"database-port",
		"3306",
		"Database port for Mariadb database connection",
	)
	flag.BoolVar(&c.Debug, "debug", false, "Using debug true or false")
	flag.StringVar(&c.JWTSecret, "jwt-secret", "jwtsecret", "use the jwt secret")
	flag.StringVar(&JwtStringExpireValue, "jwt-expire", "60s", "expiration time of jwt token")
	flag.StringVar(
		&resetTokenExpDuration,
		"reset-tkn-exp",
		"10s",
		"expiry duration of reset token in s min etc",
	)
	flag.StringVar(
		&exclPaths,
		"login-path-to-exclude",
		"/docs,/openapi,/schemas,/execs/login,/execs/forgotpassword",
		"paths to exclude when making login middleware check",
	)

	flag.Parse()

	if dbhost := getEnv("DATABASE_HOST"); dbhost != "" {
		c.DBHost = dbhost
	}

	if port := getEnv("PORT"); port != "" {
		c.Port = port
	}
	if dbuser := getEnv("DB_USER"); dbuser != "" {
		c.DBUser = dbuser
	}

	if dbpass := getEnv("DB_PASSWORD"); dbpass != "" {
		c.DBPassword = dbpass
	}
	if database := getEnv("DATABASE"); database != "" {
		c.DBDatabase = database
	}
	if dbport := getEnv("DATABASE_PORT"); dbport != "" {
		c.DBPort = dbport
	}
	if jwtsecret := getEnv("JWT_SECRET"); jwtsecret != "" {
		c.JWTSecret = jwtsecret
	}

	if JwtStrExp := getEnv("JWT_EXPIRES_IN"); JwtStrExp != "" {
		d, err := time.ParseDuration(JwtStrExp)
		if err != nil {
			panic(err)
		}
		c.JWTExpiresIn = d
	} else {
		d, err := time.ParseDuration(JwtStringExpireValue)
		if err != nil {
			panic(err)
		}
		c.JWTExpiresIn = d
	}

	if resetTokenExp := getEnv("RESET_TOKEN_EXP_DURATION"); resetTokenExp != "" {
		d, err := time.ParseDuration(resetTokenExp)
		if err != nil {
			panic(err)
		}
		c.ResetTokenExpDuration = d
	} else {
		d, err := time.ParseDuration(resetTokenExpDuration)
		if err != nil {
			panic(err)
		}
		c.ResetTokenExpDuration = d
	}

	envPaths := getEnv("LOGIN_EXCLUDE_PATHS")

	if envPaths != "" {
		paths := strings.Split(envPaths, ",")
		c.ExcludedAuthMiddlewarePath = paths
	} else {
		c.ExcludedAuthMiddlewarePath = strings.Split(exclPaths, ",")
	}

	if debugFl := getEnv("DEBUG_FL"); debugFl != "" {
		if debug, err := strconv.ParseBool(debugFl); err == nil {
			c.Debug = debug
		}
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
