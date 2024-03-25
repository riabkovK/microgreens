package config

import (
	"flag"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

const (
	// Server
	defaultServerPort               = "8000"
	defaultServerRWTimeout          = 10 * time.Second
	defaultServerMaxHeaderMegabytes = 1

	// Auth
	defaultAuthAccessTokenTTL  = 15 * time.Minute
	defaultAuthRefreshTokenTTL = 24 * time.Hour * 30

	// Logger
	defaultLoggerFormat            = "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\\n"
	defaultLoggerTimeFormat        = "15:04:05"
	defaultLoggerTimeZone          = "Local"
	defaultLoggerTimeInterval      = 500 * time.Millisecond
	defaultLoggerDisableColorsFlag = false

	// Limiter
	defaultLimiterMaxConnections = 3
	defaultLimiterExpiration     = 1 * time.Minute

	DefaultConfigDir  = "configs"
	DefaultConfigFile = "config"
	DefaultDotEnvDir  = "."

	appEnvPrefix = "MG_CFG"
)

var (
	defaultLoggerOutput        = os.Stdout
	defaultLimiterKeyGenerator = func(c *fiber.Ctx) string { return c.IP() }
)

type (
	Config struct {
		Postgres PostgresConfig
		Server   ServerConfig
		Auth     AuthConfig
		Logger   LoggerConfig
		Limiter  LimiterConfig
	}

	AuthConfig struct {
		JWT          JWTConfig
		PasswordSalt string `mapstructure:"password_salt" yaml:"password_salt"`
	}

	PostgresConfig struct {
		Username string `mapstructure:"username" yaml:"username"`
		Password string `mapstructure:"password" yaml:"password"`
		Host     string `mapstructure:"host" yaml:"host"`
		Port     string `mapstructure:"port" yaml:"port"`
		DBName   string `mapstructure:"db_name" yaml:"db_name"`
		SSLMode  string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl" yaml:"access_token_ttl"`
		RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl" yaml:"refresh_token_ttl"`
		SigningKey      string        `mapstructure:"signing_key" yaml:"signing_key"`
	}

	ServerConfig struct {
		Port               string        `mapstructure:"port" yaml:"port"`
		ReadTimeout        time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
		WriteTimeout       time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
		MaxHeaderMegabytes int           `mapstructure:"max_header_megabytes" yaml:"max_header_megabytes"`
		EnablePrintRoutes  bool          `mapstructure:"enable_print_routes" yaml:"enable_print_routes"`
	}

	LoggerConfig struct {
		Format        string        `mapstructure:"format" yaml:"format"`
		TimeFormat    string        `mapstructure:"time_format" yaml:"time_format"`
		TimeZone      string        `mapstructure:"time_zone" yaml:"time_zone"`
		TimeInterval  time.Duration `mapstructure:"time_interval" yaml:"time_interval"`
		Output        io.Writer     `mapstructure:"output" yaml:"output"`
		DisableColors bool          `mapstructure:"disable_colors" yaml:"disable_colors"`
	}

	LimiterConfig struct {
		KeyGenerator func(*fiber.Ctx) string `mapstructure:"key_generator" yaml:"key_generator"`
		Max          int                     `mapstructure:"max" yaml:"max"`
		Expiration   time.Duration           `mapstructure:"expiration" yaml:"expiration"`
	}
)

func InitConfig() (*Config, error) {
	var (
		configDir  string
		configFile string
		dotEnvDir  string
	)

	flag.StringVar(&configDir, "configDir", DefaultConfigDir, "absolute path (or relative path from app binary) to the common config yaml file")
	flag.StringVar(&configFile, "configFile", DefaultConfigFile, "name of yaml config file in configDir")
	flag.StringVar(&dotEnvDir, "dotEnvDir", DefaultDotEnvDir, "absolute path (or relative path from app binary) to the .env file")
	flag.Parse()

	logrus.Infof("Common config file: %s/%s", configDir, configFile)
	logrus.Infof("Directory with .env file: %s", dotEnvDir)

	// Init viper instance
	viperInst := viper.New()
	viperInst.SetEnvPrefix(appEnvPrefix)

	// populate viper instance with defaults envs
	populateDefault(viperInst)

	// populate viper instance with envs from config
	if err := parseConfigFile(viperInst, configDir, configFile); err != nil {
		return nil, err
	}

	// populate viper instance with envs from .env
	if err := parseDotEnvFile(viperInst, dotEnvDir); err != nil {
		return nil, err
	}

	// Override viper envs from environment variables
	viperInst.AutomaticEnv()

	// unmarshal envs in Config struct
	var cfg Config
	if err := unmarshal(viperInst, &cfg); err != nil {
		return nil, err
	}

	// Set from .env, because viper cannot unmarshal
	setFromDotEnv(viperInst, &cfg)

	return &cfg, nil
}

func unmarshal(viperInst *viper.Viper, cfg *Config) error {
	if err := viperInst.UnmarshalKey("server", &cfg.Server); err != nil {
		return err
	}
	if err := viperInst.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}
	if err := viperInst.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	if err := viperInst.UnmarshalKey("auth.jwt", &cfg.Auth.JWT); err != nil {
		return err
	}
	if err := viperInst.UnmarshalKey("logger", &cfg.Logger); err != nil {
		return err
	}
	if err := viperInst.UnmarshalKey("limiter", &cfg.Limiter); err != nil {
		return err
	}

	return nil
}

func setFromDotEnv(viperInst *viper.Viper, cfg *Config) {
	cfg.Postgres.Password = viperInst.GetString("POSTGRES.PASSWORD")
	cfg.Auth.PasswordSalt = viperInst.GetString("AUTH.PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = viperInst.GetString("AUTH.JWT.SIGNING_KEY")
}

func parseConfigFile(viperInst *viper.Viper, configDir, configFile string) error {
	viperInst.AddConfigPath(configDir)
	viperInst.SetConfigName(configFile)

	return viperInst.MergeInConfig()
}

func parseDotEnvFile(viperInst *viper.Viper, dotEnvDir string) error {
	viperInst.AddConfigPath(dotEnvDir)
	viperInst.SetConfigName(".env")
	viperInst.SetConfigType("env")

	return viperInst.MergeInConfig()
}

// populateDefault defines defaults values of envs that changes rarely
func populateDefault(viperInst *viper.Viper) {
	// Server
	viperInst.SetDefault("server.port", defaultServerPort)
	viperInst.SetDefault("server.read_timeout", defaultServerRWTimeout)
	viperInst.SetDefault("server.write_timeout", defaultServerRWTimeout)
	viperInst.SetDefault("server.max_header_megabytes", defaultServerMaxHeaderMegabytes)

	// Auth
	viperInst.SetDefault("auth.jwt.access_token_ttl", defaultAuthAccessTokenTTL)
	viperInst.SetDefault("auth.jwt.refresh_token_ttl", defaultAuthRefreshTokenTTL)

	// Logger
	viperInst.SetDefault("logger.format", defaultLoggerFormat)
	viperInst.SetDefault("logger.time_format", defaultLoggerTimeFormat)
	viperInst.SetDefault("logger.time_zone", defaultLoggerTimeZone)
	viperInst.SetDefault("logger.time_interval", defaultLoggerTimeInterval)
	viperInst.SetDefault("logger.output", defaultLoggerOutput)
	viperInst.SetDefault("logger.disable_colors", defaultLoggerDisableColorsFlag)

	// Limiter
	viperInst.SetDefault("limiter.key_generator", defaultLimiterKeyGenerator)
	viperInst.SetDefault("limiter.max", defaultLimiterMaxConnections)
	viperInst.SetDefault("limiter.expiration", defaultLimiterExpiration)
}
