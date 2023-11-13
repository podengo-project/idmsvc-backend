// The scope of this file is:
// - Define the configuration struct.
// - Set default configuration values.
// - Map the data so viper can load the configuration there.
// See: https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
// See: https://consoledot.pages.redhat.com/docs/dev/getting-started/migration/config.html
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	validator "github.com/go-playground/validator/v10"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"k8s.io/utils/env"
)

const (
	// DefaultAppName is used to compose the route paths
	DefaultAppName = "idmsvc"
	// API URL path prefix
	DefaultPathPrefix = "/api/idmsvc/v1"
	// DefaultExpirationTime is used for the default token expiration period
	// expressed in seconds. The default value is set to 7200 (2 hours)
	DefaultTokenExpirationTimeSeconds = 7200
	// DefaultWebPort is the default port where the public API is listening
	DefaultWebPort = 8000

	// https://github.com/project-koku/koku/blob/main/koku/api/common/pagination.py

	// PaginationDefaultLimit is the default limit for the pagination
	PaginationDefaultLimit = 10
	// PaginationMaxLimit is the default max limit for the pagination
	PaginationMaxLimit = 1000

	// DefaultAcceptXRHFakeIdentity is disabled
	DefaultAcceptXRHFakeIdentity = false
	// DefaultValidateAPI is true
	DefaultValidateAPI = true
)

type Config struct {
	Loaded      bool
	Web         Web
	Database    Database
	Logging     Logging
	Kafka       Kafka
	Cloudwatch  Cloudwatch
	Metrics     Metrics
	Clients     Clients
	Application Application `mapstructure:"app"`
}

type Web struct {
	Port int16
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	// https://stackoverflow.com/questions/54844546/how-to-unmarshal-golang-viper-snake-case-values
	CACertPath string `mapstructure:"ca_cert_path"`
}

type Logging struct {
	Level   string
	Console bool
}

type Cloudwatch struct {
	Region  string
	Key     string
	Secret  string
	Session string
	Group   string
	Stream  string
}

type Kafka struct {
	Timeout int
	Group   struct {
		Id string
	}
	Auto struct {
		Offset struct {
			Reset string
		}
		Commit struct {
			Interval struct {
				Ms int
			}
		}
	}
	Bootstrap struct {
		Servers string
	}
	Topics []string
	Sasl   struct {
		Username  string
		Password  string
		Mechanism string
		Protocol  string
	}
	Request struct {
		Timeout struct {
			Ms int
		}
		Required struct {
			Acks int
		}
	}
	Capath  string
	Message struct {
		Send struct {
			Max struct {
				Retries int
			}
		}
	}
	Retry struct {
		Backoff struct {
			Ms int
		}
	}
}

type Metrics struct {
	// Defines the path to the metrics server that the app should be configured to
	// listen on for metric traffic.
	Path string `mapstructure:"path"`

	// Defines the metrics port that the app should be configured to listen on for
	// metric traffic.
	Port int `mapstructure:"port"`
}

// Clients gather all the client settings for the
type Clients struct {
	Inventory InventoryClient
}

type InventoryClient struct {
	// Define the base url for the host inventory service
	BaseUrl string `mapstructure:"base_url"`
}

// Application hold specific application settings
type Application struct {
	// API URL's path prefix, e.g. /api/idmsvc/v1
	PathPrefix string `mapstructure:"url_path_prefix" validate:"required"`
	// This is the default expiration time for the token
	// generated when a RHEL IDM domain is created
	TokenExpirationTimeSeconds int `mapstructure:"token_expiration_seconds" validate:"gte=600,lte=86400"`
	// Indicate the default pagination limit when it is 0 or not filled
	PaginationDefaultLimit int `mapstructure:"pagination_default_limit"`
	// Indicate the max pagination limit when it is grather
	PaginationMaxLimit int `mapstructure:"pagination_max_limit"`
	// AcceptXRHFakeIdentity define when the fake middleware is added to the route
	// to process the x-rh-fake-identity
	AcceptXRHFakeIdentity bool `mapstructure:"accept_x_rh_fake_identity"`
	// ValidateAPI indicate when the middleware to validate the API
	// requests and responses is disabled; by default it is enabled.
	ValidateAPI bool `mapstructure:"validate_api"`
	// secret for various MAC and encryptions like domain registration
	// token and encrypted private JWKs. "random" generates an ephemeral secret.
	// Secrets are derived with HKDF-SHA256.
	MainSecret string `mapstructure:"secret" validate:"required,base64rawurl"`
}

var config *Config = nil

func setDefaults(v *viper.Viper) {
	if v == nil {
		panic("viper instance cannot be nil")
	}
	// Web
	v.SetDefault("web.port", DefaultWebPort)

	// Database
	v.SetDefault("database.host", "")
	v.SetDefault("database.port", "")
	v.SetDefault("database.user", "")
	v.SetDefault("database.password", "")
	v.SetDefault("database.name", "")
	v.SetDefault("database.ca_cert_path", "")

	// Kafka
	addEventConfigDefaults(v)

	// Clowdwatch

	// Miscelanea
	v.SetDefault("logging.level", "info")

	// Clients
	v.SetDefault("clients.host_inventory_base_url", "http://localhost:8010/api/inventory/v1")

	// Application specific

	// Set default value for application expiration time for
	// the token created by the RHEL IDM domains
	v.SetDefault("app.token_expiration_seconds", DefaultTokenExpirationTimeSeconds)
	v.SetDefault("app.pagination_default_limit", PaginationDefaultLimit)
	v.SetDefault("app.pagination_max_limit", PaginationMaxLimit)
	v.SetDefault("app.accept_x_rh_fake_identity", DefaultAcceptXRHFakeIdentity)
	v.SetDefault("app.validate_api", DefaultValidateAPI)
	v.SetDefault("app.url_path_prefix", DefaultPathPrefix)
	v.SetDefault("app.secret", "")
	v.SetDefault("app.debug", false)
}

func setClowderConfiguration(v *viper.Viper, clowderConfig *clowder.AppConfig) {
	if v == nil {
		panic("'v' is nil")
	}
	if clowderConfig == nil {
		panic("'clowderConfig' is nil")
	}

	// Web
	v.Set("web.port", clowderConfig.PublicPort)

	// Database
	var rdsCertPath string
	if clowderConfig.Database != nil && clowderConfig.Database.RdsCa != nil {
		var err error
		if rdsCertPath, err = clowderConfig.RdsCa(); err != nil {
			slog.Warn("Cannot read RDS CA cert", slog.Any("error", err))
		}
	}
	if clowderConfig.Database != nil {
		v.Set("database.host", clowderConfig.Database.Hostname)
		v.Set("database.port", clowderConfig.Database.Port)
		v.Set("database.user", clowderConfig.Database.Username)
		v.Set("database.password", clowderConfig.Database.Password)
		v.Set("database.name", clowderConfig.Database.Name)
		if rdsCertPath != "" {
			v.Set("database.ca_cert_path", rdsCertPath)
		}
	}

	// Clowdwatch
	if clowderConfig.Logging.Cloudwatch != nil {
		v.Set("cloudwatch.region", clowderConfig.Logging.Cloudwatch.Region)
		v.Set("cloudwatch.group", clowderConfig.Logging.Cloudwatch.LogGroup)
		v.Set("cloudwatch.secret", clowderConfig.Logging.Cloudwatch.SecretAccessKey)
		v.Set("cloudwatch.key", clowderConfig.Logging.Cloudwatch.AccessKeyId)
	}

	// Metrics configuration
	v.Set("metrics.path", clowderConfig.MetricsPath)
	v.Set("metrics.port", clowderConfig.MetricsPort)
}

func Load(cfg *Config) *viper.Viper {
	var (
		err error
	)

	if cfg == nil {
		panic("'cfg' is nil")
	}

	v := viper.New()
	v.AddConfigPath(env.GetString("CONFIG_PATH", "./configs"))
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)
	if clowder.IsClowderEnabled() {
		setClowderConfiguration(v, clowder.LoadedConfig)
	}

	if err = v.ReadInConfig(); err != nil {
		slog.Warn("Not using config.yaml", slog.Any("error", err))
	}
	if err = v.Unmarshal(cfg); err != nil {
		slog.Warn("Mapping to configuration", slog.Any("error", err))
	}

	return v
}

func reportError(err error) {
	for _, err := range err.(validator.ValidationErrors) {
		slog.Error(
			"Configuration validation error",
			slog.String("namespace", err.Namespace()),
			slog.Group("rule",
				slog.String("tag", err.Tag()),
				slog.Any("value", err.Value),
			),
			slog.String("got", err.Param()),
			slog.String("type", err.Kind().String()),
		)
	}
}

func Validate(cfg *Config) (err error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}

// Get is a singleton to get the global loaded configuration.
func Get() *Config {
	if config != nil {
		return config
	}
	config = &Config{}
	v := Load(config)

	// Dump configuration as JSON
	if config.Logging.Level == "debug" {
		c := v.AllSettings()
		b, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

	if err := Validate(config); err != nil {
		reportError(err)
		panic("Invalid configuration")
	}
	return config
}

func hasKafkaBrokerConfig(cfg *clowder.AppConfig) bool {
	if cfg == nil || cfg.Kafka == nil || cfg.Kafka.Brokers == nil || len(cfg.Kafka.Brokers) <= 0 {
		return false
	}
	broker := cfg.Kafka.Brokers[0]
	if broker.Hostname == "" || broker.Port == nil {
		return false
	}
	return true
}

func addEventConfigDefaults(options *viper.Viper) {
	if options == nil {
		panic("'options' is nil")
	}
	options.SetDefault("kafka.timeout", 10000)
	options.SetDefault("kafka.group.id", DefaultAppName)
	options.SetDefault("kafka.auto.offset.reset", "latest")
	options.SetDefault("kafka.auto.commit.interval.ms", 5000)
	options.SetDefault("kafka.request.required.acks", -1) // -1 == "all"
	options.SetDefault("kafka.message.send.max.retries", 15)
	options.SetDefault("kafka.retry.backoff.ms", 100)

	if !clowder.IsClowderEnabled() {
		// If clowder is not present, set defaults to local configuration
		TopicTranslationConfig = NewTopicTranslationWithDefaults()
		options.SetDefault("kafka.bootstrap.servers", readEnv("KAFKA_BOOTSTRAP_SERVERS", ""))
		options.SetDefault("kafka.topics", "platform."+DefaultAppName+".domain-created")
		return
	}

	// Settings for clowder
	cfg := clowder.LoadedConfig
	TopicTranslationConfig = NewTopicTranslationWithClowder(cfg)
	options.SetDefault("kafka.bootstrap.servers", strings.Join(clowder.KafkaServers, ","))

	// Prepare topics
	topics := []string{}
	for _, value := range clowder.KafkaTopics {
		topics = append(topics, value.Name)
	}
	options.SetDefault("kafka.topics", strings.Join(topics, ","))

	if !hasKafkaBrokerConfig(cfg) {
		return
	}

	if cfg.Kafka.Brokers[0].Cacert != nil {
		// This method is writing only the first CA but if
		// that behavior changes in the future, nothing will
		// be changed here
		caPath, err := cfg.KafkaCa(cfg.Kafka.Brokers...)
		if err != nil {
			panic("Kafka CA failed to write")
		}
		options.Set("kafka.capath", caPath)
	}

	broker := cfg.Kafka.Brokers[0]
	if broker.Authtype != nil {
		options.Set("kafka.sasl.username", *broker.Sasl.Username)
		options.Set("kafka.sasl.password", *broker.Sasl.Password)
		options.Set("kafka.sasl.mechanism", *broker.Sasl.SaslMechanism)
		options.Set("kafka.sasl.protocol", *broker.Sasl.SecurityProtocol)
	}
}

func readEnv(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = def
	}
	return value
}
