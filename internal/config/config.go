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

const DefaultAppName = "hmsidm"

type Config struct {
	Loaded     bool
	Web        Web
	Database   Database
	Logging    Logging
	Kafka      Kafka
	Cloudwatch Cloudwatch
	Metrics    Metrics
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

var config *Config = nil

func setDefaults(v *viper.Viper) {
	if v == nil {
		panic("viper instance cannot be nil")
	}
	// Web
	v.SetDefault("web.port", 8000)

	// Database

	// Kafka
	addEventConfigDefaults(v)

	// Clowdwatch

	// Miscelanea
	v.SetDefault("logging.level", "info")
}

func setClowderConfiguration(v *viper.Viper) {
	if !clowder.IsClowderEnabled() {
		return
	}

	cfg := clowder.LoadedConfig
	if cfg == nil {
		log.Error().Msg("clowder.LoadedConfig is nil")
		return
	}
	var rdsCertPath string
	if cfg.Database != nil && cfg.Database.RdsCa != nil {
		var err error
		if rdsCertPath, err = cfg.RdsCa(); err != nil {
			log.Warn().Err(err).Msg("Cannot read RDS CA cert")
		}
	}

	// Web
	v.Set("web.port", cfg.PublicPort)

	// Database
	if cfg.Database != nil {
		v.Set("database.host", cfg.Database.Hostname)
		v.Set("database.port", cfg.Database.Port)
		v.Set("database.user", cfg.Database.Username)
		v.Set("database.password", cfg.Database.Password)
		v.Set("database.name", cfg.Database.Name)
		if rdsCertPath != "" {
			v.Set("database.ca_cert_path", rdsCertPath)
		}
	}

	// Clowdwatch
	v.Set("cloudwatch.region", cfg.Logging.Cloudwatch.Region)
	v.Set("cloudwatch.group", cfg.Logging.Cloudwatch.LogGroup)
	v.Set("cloudwatch.secret", cfg.Logging.Cloudwatch.SecretAccessKey)
	v.Set("cloudwatch.key", cfg.Logging.Cloudwatch.AccessKeyId)

	// Metrics configuration
	v.Set("metrics.path", cfg.MetricsPath)
	v.Set("metrics.port", cfg.MetricsPort)
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
	setClowderConfiguration(v)
	v.AutomaticEnv()
	if err = v.ReadInConfig(); err != nil {
		log.Warn().Msgf("Not using config.yaml: %s", err.Error())
	}
	if err = v.Unmarshal(cfg); err != nil {
		log.Warn().Msgf("Mapping to configuration: %s", err.Error())
	}

	return cfg
}

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
	return true
}

func addEventConfigDefaults(options *viper.Viper) {
	options.SetDefault("kafka.timeout", 10000)
	options.SetDefault("kafka.group.id", "hmsidm")
	options.SetDefault("kafka.auto.offset.reset", "latest")
	options.SetDefault("kafka.auto.commit.interval.ms", 5000)
	options.SetDefault("kafka.request.required.acks", -1) // -1 == "all"
	options.SetDefault("kafka.message.send.max.retries", 15)
	options.SetDefault("kafka.retry.backoff.ms", 100)

	if !clowder.IsClowderEnabled() {
		// If clowder is not present, set defaults to local configuration
		TopicTranslationConfig = NewTopicTranslationWithDefaults()
		options.SetDefault("kafka.bootstrap.servers", readEnv("KAFKA_BOOTSTRAP_SERVERS", ""))
		options.SetDefault("kafka.topics", "platform.hmsidm.todo-created")
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
