package config

import (
	"testing"

	validator "github.com/go-playground/validator/v10"
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/secrets"
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

	// TODO HMS-3609 There are many more properties to check; complete the test
	assert.Equal(t, DefaultWebPort, v.Get("web.port"))
	assert.Equal(t, "info", v.Get("logging.level"))
	assert.Equal(t, "http://localhost:8010/api/inventory/v1", v.Get("clients.host_inventory_base_url"))
	assert.Equal(t, DefaultTokenExpirationTimeSeconds, v.Get("app.token_expiration_seconds"))
	assert.Equal(t, PaginationDefaultLimit, v.Get("app.pagination_default_limit"))
	assert.Equal(t, PaginationMaxLimit, v.Get("app.pagination_max_limit"))

	assert.Equal(t, DefaultEnableRBAC, v.Get("app.enable_rbac"))
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

	// Load RDSACert with wrong path
	cfg.Database = &clowder.DatabaseConfig{}
	cfg.Database.RdsCa = pointy.String("/tmp/itdoesnotexist.pem")
	assert.NotPanics(t, func() {
		setClowderConfiguration(v, &cfg)
	})

	// Load RDSACert nil
	cfg.Database.RdsCa = nil
	assert.NotPanics(t, func() {
		setClowderConfiguration(v, &cfg)
	})

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
	assert.NotNil(t, v.Get("database.ca_cert_path"))

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
	assert.Equal(t, "testregion", v.Get("logging.cloudwatch.region"))
	assert.Equal(t, "testgroup", v.Get("logging.cloudwatch.group"))
	assert.Equal(t, "testsecret", v.Get("logging.cloudwatch.secret"))
	assert.Equal(t, "testaccesskeyid", v.Get("logging.cloudwatch.key"))
	assert.Equal(t, "cloudwatch", v.Get("logging.type"))
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

func TestLoad(t *testing.T) {
	// 'cfg' is nil panic
	assert.Panics(t, func() {
		Load(nil)
	}, "'cfg' is nil")

	// Success Load
	cfg := Config{}
	assert.NotPanics(t, func() {
		Load(&cfg)
	})
}

func TestValidateConfig(t *testing.T) {
	cfg := Config{
		Application: Application{
			Name:                        DefaultAppName,
			PathPrefix:                  DefaultPathPrefix,
			MainSecret:                  secrets.GenerateRandomMainSecret(),
			TokenExpirationTimeSeconds:  0,
			HostconfJwkValidity:         DefaultHostconfJwkValidity,
			HostconfJwkRenewalThreshold: DefaultHostconfJwkRenewalThreshold,
		},
	}
	err := Validate(&cfg)
	assert.Error(t, err)
	ve, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)
	// TODO Change the order; it should be (t, expected_value, current_calue)
	assert.Equal(t, len(ve), 1)
	assert.Equal(t, ve[0].Namespace(), "Config.Application.TokenExpirationTimeSeconds")
	assert.Equal(t, ve[0].Tag(), "gte")
}
