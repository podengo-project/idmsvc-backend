package config

import (
	"testing"

	"github.com/openlyinc/pointy"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {
	assert.Panics(t, func() {
		setDefaults(nil)
	})

	v := viper.New()
	assert.NotPanics(t, func() {
		setDefaults(v)
	})

	assert.Equal(t, DefaultWebPort, v.Get("web.port"))
	assert.Equal(t, "info", v.Get("logging.level"))
	assert.Equal(t, "http://localhost:8010/api/inventory/v1", v.Get("clients.host_inventory_base_url"))
	assert.Equal(t, DefatulExpirationTime, v.Get("app.expiration_time"))
	assert.Equal(t, PaginationDefaultLimit, v.Get("app.pagination_default_limit"))
	assert.Equal(t, PaginationMaxLimit, v.Get("app.pagination_max_limit"))
}

func TestSetClowderConfiguration(t *testing.T) {
	// Panic for v = nil
	assert.PanicsWithValue(t, "'v' is nil", func() {
		setClowderConfiguration(nil, nil)
	})

	// Panic for cfg = nil
	v := viper.New()
	assert.PanicsWithValue(t, "'clowderConfig' is nil", func() {
		setClowderConfiguration(v, nil)
	})

	// Loading empty clowder.AppConfig
	cfg := clowder.AppConfig{}
	assert.NotPanics(t, func() {
		setClowderConfiguration(v, &cfg)
	})
	assert.Nil(t, v.Get("database.host"))
	assert.Nil(t, v.Get("database.port"))
	assert.Nil(t, v.Get("database.user"))
	assert.Nil(t, v.Get("database.password"))
	assert.Nil(t, v.Get("database.name"))
	assert.Nil(t, v.Get("database.ca_cert_path"))

	// Load database data (but RdsCa)
	cfg.Database = &clowder.DatabaseConfig{
		Hostname: "testhost",
		Port:     5432,
		Username: "testuser",
		Password: "testpassword",
		Name:     "testname",
	}
	assert.NotPanics(t, func() {
		setClowderConfiguration(v, &cfg)
	})
	assert.Equal(t, "testhost", v.Get("database.host"))
	assert.Equal(t, 5432, v.Get("database.port"))
	assert.Equal(t, "testuser", v.Get("database.user"))
	assert.Equal(t, "testpassword", v.Get("database.password"))
	assert.Equal(t, "testname", v.Get("database.name"))
	assert.Nil(t, v.Get("database.ca_cert_path"))

	// TODO Add test to cover RdsCa flow

	// Load cloudwatch data
	cfg.Logging.Cloudwatch = &clowder.CloudWatchConfig{
		Region:          "testregion",
		LogGroup:        "testgroup",
		SecretAccessKey: "testsecret",
		AccessKeyId:     "testaccesskeyid",
	}
	assert.NotPanics(t, func() {
		setClowderConfiguration(v, &cfg)
	})
	assert.Equal(t, "testregion", v.Get("cloudwatch.region"))
	assert.Equal(t, "testgroup", v.Get("cloudwatch.group"))
	assert.Equal(t, "testsecret", v.Get("cloudwatch.secret"))
	assert.Equal(t, "testaccesskeyid", v.Get("cloudwatch.key"))
}

func TestHasKafkaBrokerConfig(t *testing.T) {
	assert.False(t, hasKafkaBrokerConfig(nil))
	cfg := clowder.AppConfig{}
	assert.False(t, hasKafkaBrokerConfig(&cfg))
	cfg.Kafka = &clowder.KafkaConfig{}
	assert.False(t, hasKafkaBrokerConfig(&cfg))
	cfg.Kafka.Brokers = []clowder.BrokerConfig{}
	assert.False(t, hasKafkaBrokerConfig(&cfg))
	cfg.Kafka.Brokers = append(cfg.Kafka.Brokers, clowder.BrokerConfig{})
	assert.False(t, hasKafkaBrokerConfig(&cfg))
	cfg.Kafka.Brokers[0].Hostname = "test.kafka.svc.localdomain"
	assert.False(t, hasKafkaBrokerConfig(&cfg))
	cfg.Kafka.Brokers[0].Port = pointy.Int(3000)
	assert.True(t, hasKafkaBrokerConfig(&cfg))
}

func TestAddEventConfigDefaults(t *testing.T) {
	assert.PanicsWithValue(t, "'options' is nil", func() {
		addEventConfigDefaults(nil)
	})

	v := viper.New()
	addEventConfigDefaults(v)
	assert.Equal(t, 10000, v.Get("kafka.timeout"))
	assert.Equal(t, DefaultAppName, v.Get("kafka.group.id"))
	assert.Equal(t, "latest", v.Get("kafka.auto.offset.reset"))
	assert.Equal(t, 5000, v.Get("kafka.auto.commit.interval.ms"))
	assert.Equal(t, -1, v.Get("kafka.request.required.acks"))
	assert.Equal(t, 15, v.Get("kafka.message.send.max.retries"))
	assert.Equal(t, 100, v.Get("kafka.retry.backoff.ms"))

}
