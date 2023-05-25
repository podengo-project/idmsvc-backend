// The scope of this file is:
// - Define the configuration struct.
// - Set default configuration values.
// - Map the data so viper can load the configuration there.
// See: https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
// See: https://consoledot.pages.redhat.com/docs/dev/getting-started/migration/config.html
package config

import (
	"os"
	"strings"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"k8s.io/utils/env"
)

const (
	// DefaultAppName is used to compose the route paths
	DefaultAppName = "hmsidm"
	// DefatulExpirationTime is used for the default token expiration period
	DefatulExpirationTime = 15
	// DefaultWebPort is the default port where the public API is listening
	DefaultWebPort = 8000

	// https://github.com/project-koku/koku/blob/main/koku/api/common/pagination.py

	// PaginationDefaultLimit is the default limit for the pagination
	PaginationDefaultLimit = 10
	// PaginationMaxLimit is the default max limit for the pagination
	PaginationMaxLimit = 1000
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

type Clients struct {
	Inventory InventoryClient
}

type InventoryClient struct {
	// Define the base url for the host inventory service
	BaseUrl string `mapstructure:"base_url"`
}

// Application hold specific application settings
type Application struct {
	// This is the default expiration time for the token
	// generated when a RHEL IDM domain is created
	ExpirationTime int `mapstructure:"expiration_time"`
	// Indicate the default pagination limit when it is 0 or not filled
	PaginationDefaultLimit int `mapstructure:"pagination_default_limit"`
	// Indicate the max pagination limit when it is grather
	PaginationMaxLimit int `mapstructure:"pagination_max_limit"`
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
	v.SetDefault("app.expiration_time", DefatulExpirationTime)
	v.SetDefault("app.pagination_default_limit", PaginationDefaultLimit)
	v.SetDefault("app.pagination_max_limit", PaginationMaxLimit)
}

func setClowderConfiguration(v *viper.Viper, clowderConfig *clowder.AppConfig) {
	if v == nil {
		panic("'v' is nil")
	}
	if clowderConfig == nil {
		panic("'clowderConfig' is nil")
	}
	var rdsCertPath string
	if clowderConfig.Database != nil && clowderConfig.Database.RdsCa != nil {
		var err error
		if rdsCertPath, err = clowderConfig.RdsCa(); err != nil {
			log.Warn().Err(err).Msg("Cannot read RDS CA cert")
		}
	}

	// Web
	v.Set("web.port", clowderConfig.PublicPort)

	// Database
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

func Load(cfg *Config) *Config {
	var (
		err error
	)

	if cfg == nil {
		panic("cfg is nil")
	}

	v := viper.New()
	v.AddConfigPath(env.GetString("CONFIG_PATH", "./configs"))
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	setDefaults(v)
	if clowder.IsClowderEnabled() {
		setClowderConfiguration(v, clowder.LoadedConfig)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err = v.ReadInConfig(); err != nil {
		log.Warn().Msgf("Not using config.yaml: %s", err.Error())
	}
	if err = v.Unmarshal(cfg); err != nil {
		log.Warn().Msgf("Mapping to configuration: %s", err.Error())
	}

	return cfg
}

// Get is a singleton to get the global loaded configuration.
func Get() *Config {
	if config != nil {
		return config
	}
	config = &Config{}
	return Load(config)
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
