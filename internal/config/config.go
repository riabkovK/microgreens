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
		PasswordSalt string
	}

	PostgresConfig struct {
		Username string `mapstructure:"username" yaml:"username"`
		Password string
		Host     string `mapstructure:"host" yaml:"host"`
		Port     string `mapstructure:"port" yaml:"port"`
		DBName   string `mapstructure:"DBName" yaml:"DBName"`
		SSLMode  string `mapstructure:"SSLMode" yaml:"SSLMode"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL" yaml:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL" yaml:"refreshTokenTTL"`
		SigningKey      string
	}

	ServerConfig struct {
		Port               string        `mapstructure:"port" yaml:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout" yaml:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout" yaml:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes" yaml:"maxHeaderMegabytes"`
		EnablePrintRoutes  bool          `mapstructure:"enablePrintRoutes" yaml:"enablePrintRoutes"`
	}

	LoggerConfig struct {
		Format        string        `mapstructure:"format" yaml:"format"`
		TimeFormat    string        `mapstructure:"timeFormat" yaml:"timeFormat"`
		TimeZone      string        `mapstructure:"timeZone" yaml:"timeZone"`
		TimeInterval  time.Duration `mapstructure:"timeInterval" yaml:"timeInterval"`
		Output        io.Writer     `mapstructure:"output" yaml:"output"`
		DisableColors bool          `mapstructure:"disableColors" yaml:"disableColors"`
	}

	LimiterConfig struct {
		KeyGenerator func(*fiber.Ctx) string `mapstructure:"keyGenerator" yaml:"keyGenerator"`
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

	// Gives the opportunity to override from environment variables
	viperInst.AutomaticEnv()
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

func setFromDotEnv(viperInst *viper.Viper, cfg *Config) {
	cfg.Postgres.Password = viperInst.GetString("POSTGRES_PASSWORD")
	cfg.Auth.PasswordSalt = viperInst.GetString("AUTH_PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = viperInst.GetString("AUTH_JWT_SIGNING_KEY")
}

// populateDefault defines defaults values of envs that changes rarely
func populateDefault(viperInst *viper.Viper) {
	// Server
	viperInst.SetDefault("server.port", defaultServerPort)
	viperInst.SetDefault("server.readTimeout", defaultServerRWTimeout)
	viperInst.SetDefault("server.writeTimeout", defaultServerRWTimeout)
	viperInst.SetDefault("server.maxHeaderMegabytes", defaultServerMaxHeaderMegabytes)

	// Auth
	viperInst.SetDefault("auth.jwt.accessTokenTTL", defaultAuthAccessTokenTTL)
	viperInst.SetDefault("auth.jwt.refreshTokenTTL", defaultAuthRefreshTokenTTL)

	// Logger
	viperInst.SetDefault("logger.format", defaultLoggerFormat)
	viperInst.SetDefault("logger.timeFormat", defaultLoggerTimeFormat)
	viperInst.SetDefault("logger.timeZone", defaultLoggerTimeZone)
	viperInst.SetDefault("logger.timeInterval", defaultLoggerTimeInterval)
	viperInst.SetDefault("logger.output", defaultLoggerOutput)
	viperInst.SetDefault("logger.disableColors", defaultLoggerDisableColorsFlag)

	// Limiter
	viperInst.SetDefault("limiter.keyGenerator", defaultLimiterKeyGenerator)
	viperInst.SetDefault("limiter.max", defaultLimiterMaxConnections)
	viperInst.SetDefault("limiter.expiration", defaultLimiterExpiration)
}
