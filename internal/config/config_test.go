package config

import (
	"bytes"
	"log/slog"
	"os"
	"strconv"
	"testing"

	validator "github.com/go-playground/validator/v10"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/secrets"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.openly.dev/pointy"
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

func TestLog(t *testing.T) {
	// Given
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))

	cfg := Config{
		Web: Web{
			Port: 8080,
		},
		Database: Database{
			User:     "testdbuser",
			Password: "testdbpassword",
		},
		Logging: Logging{
			Level: "debug",
			Cloudwatch: Cloudwatch{
				Secret: "testcloudwatchsecret",
			},
		},
		Metrics: Metrics{
			Path: "/metrics",
		},
		Clients: Clients{
			RbacBaseURL:        "http://localhost:8080/api/rbac/v1",
			PendoAPIKey:        "testpendoapikey",
			PendoTrackEventKey: "testpendotrackeventkey",
		},
		Application: Application{
			Name:       "appname",
			MainSecret: "testmainsecret",
		},
	}

	// When
	cfg.Log(logger)
	loggedStr := buf.String()

	// Then
	assert.Contains(t, loggedStr, "Web.Port=8080")
	assert.Contains(t, loggedStr, "Database.User=testdbuser")
	assert.Contains(t, loggedStr, "Database.Password=***")
	assert.Contains(t, loggedStr, "Logging.Level=debug")
	assert.Contains(t, loggedStr, "Logging.Cloudwatch.Secret=***")
	assert.Contains(t, loggedStr, "Metrics.Path=/metrics")
	assert.Contains(t, loggedStr, "Clients.RbacBaseURL=http://localhost:8080/api/rbac/v1")
	assert.Contains(t, loggedStr, "Clients.PendoAPIKey=***")
	assert.Contains(t, loggedStr, "Clients.PendoTrackEventKey=***")
	assert.Contains(t, loggedStr, "Application.Name=appname")
	assert.Contains(t, loggedStr, "Application.MainSecret=***")

	// No password in the log
	assert.NotContains(t, loggedStr, "testdbpassword")
	assert.NotContains(t, loggedStr, "testcloudwatchsecret")
	assert.NotContains(t, loggedStr, "testpendoapikey")
	assert.NotContains(t, loggedStr, "testpendotrackeventkey")
	assert.NotContains(t, loggedStr, "testmainsecret")
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
			IdleTimeout:                 DefaultIdleTimeout,
			ReadTimeout:                 DefaultReadTimeout,
			WriteTimeout:                DefaultWriteTimeout,
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

func TestGuardProcessPublicEndpoint(t *testing.T) {
	assert.PanicsWithValue(t, "'serviceName' is an empty string", func() {
		guardUpdateServiceBasePath("", "", nil, nil)
	})
	assert.PanicsWithValue(t, "'version' is an empty string", func() {
		guardUpdateServiceBasePath("rbac", "", nil, nil)
	})
	assert.PanicsWithValue(t, "'target' is nil", func() {
		guardUpdateServiceBasePath("rbac", "v1", nil, nil)
	})

	v := viper.New()
	assert.PanicsWithValue(t, "'clowderConfig' is nil", func() {
		guardUpdateServiceBasePath("rbac", "v1", v, nil)
	})

	clowderConfig := clowder.AppConfig{}
	assert.NotPanics(t, func() {
		guardUpdateServiceBasePath("rbac", "v1", v, &clowderConfig)
	})
}

func TestUpdateServiceBasePath(t *testing.T) {
	const (
		hostname                      = "rbac-service.ephemeral-zzym5j.svc"
		viperPath                     = "clients.rbac_base_url"
		origBaseURL                   = "http://localhost:8080/api/rbac/v1"
		processedWithNoTLSPortBaseURL = "http://" + hostname + ":8000/api/rbac/v1"
		processedWithTLSPortBaseURL   = "https://" + hostname + ":8080/api/rbac/v1"
	)
	v := viper.New()
	v.Set(viperPath, origBaseURL)
	clowderConfig := clowder.AppConfig{}
	clowderConfig.Endpoints = []clowder.DependencyEndpoint{
		{
			App:      "rbac",
			Name:     "service",
			ApiPath:  "/api/rbac/",
			Hostname: hostname,
			Port:     8000,
			TlsPort:  nil,
		},
	}

	clowderConfig.Endpoints[0].ApiPath = ""
	clowderConfig.Endpoints[0].ApiPaths = []string{"/api/rbac/", "/api/rbac-2/"}
	updateServiceBasePath("rbac", "v1", v, &clowderConfig)
	assert.Equal(t, processedWithNoTLSPortBaseURL, v.GetString(viperPath))

	clowderConfig.Endpoints[0].TlsPort = pointy.Int(0)
	updateServiceBasePath("rbac", "v1", v, &clowderConfig)
	assert.Equal(t, processedWithNoTLSPortBaseURL, v.GetString(viperPath))

	clowderConfig.Endpoints[0].TlsPort = pointy.Int(8080)
	updateServiceBasePath("rbac", "v1", v, &clowderConfig)
	assert.Equal(t, processedWithTLSPortBaseURL, v.GetString(viperPath))
}

func TestGetEndpoint(t *testing.T) {
	clowderConfig := clowder.AppConfig{}
	assert.Nil(t, getEndpoint("rbac", &clowderConfig),
		"Return nil when the service is not found in the DependencyEndpoints")

	clowderConfig.Endpoints = []clowder.DependencyEndpoint{
		{
			App: "rbac",
		},
	}
	assert.Equal(t, &clowder.DependencyEndpoint{
		App: "rbac",
	}, getEndpoint("rbac", &clowderConfig),
		"Return the matched endpoint")

	assert.Nil(t, getEndpoint("host-inventory", &clowderConfig),
		"Return nil when no matching endpoint is found")
}

func TestBuildClientBaseURL(t *testing.T) {
	const (
		hostname = "rbac-service.ephemeral-zzym5j.svc"
	)
	dependencyEndpoint := clowder.DependencyEndpoint{
		App:      "rbac",
		Name:     "service",
		Hostname: hostname,
	}

	dependencyEndpoint.TlsPort = nil
	assert.Equal(t, "", buildClientBaseURL(&dependencyEndpoint, "v1"),
		"No ports information returns an empty string")

	dependencyEndpoint.Port = 8000
	assert.Equal(t, "", buildClientBaseURL(&dependencyEndpoint, "v1"),
		"No path information returns an empty string")

	dependencyEndpoint.ApiPaths = []string{"/api/rbac-1", "/api/rbac-2"}
	dependencyEndpoint.TlsPort = pointy.Int(8080)
	assert.Equal(t, "https://"+hostname+":"+strconv.Itoa(*dependencyEndpoint.TlsPort)+"/api/rbac-1/v1", buildClientBaseURL(&dependencyEndpoint, "v1"),
		"build string using the TLS port and the first apiPaths item")
}

func TestBuildClientBaseURLSchemaHostPort(t *testing.T) {
	const (
		hostname = "rbac-service.ephemeral-zzym5j.svc"
	)
	dependencyEndpoint := clowder.DependencyEndpoint{
		App:      "rbac",
		Name:     "service",
		Hostname: hostname,
	}

	assert.Equal(t, "", buildClientBaseURLSchemaHostPort(&dependencyEndpoint),
		"return empty when no ports information are present")

	dependencyEndpoint.Port = 8000
	assert.Equal(t, "http://"+hostname+":"+strconv.Itoa(dependencyEndpoint.Port), buildClientBaseURLSchemaHostPort(&dependencyEndpoint),
		"return the no TLS port when it is present and no TLSPort information")

	dependencyEndpoint.TlsPort = pointy.Int(8080)
	assert.Equal(t, "https://"+hostname+":"+strconv.Itoa(*dependencyEndpoint.TlsPort), buildClientBaseURLSchemaHostPort(&dependencyEndpoint),
		"when TLS port is present, it is preferred")
}

func TestBuildClientBaseURLPath(t *testing.T) {
	const (
		hostname = "rbac-service.ephemeral-zzym5j.svc"
	)
	dependencyEndpoint := clowder.DependencyEndpoint{
		App:      "rbac",
		Name:     "service",
		Hostname: hostname,
	}

	assert.Equal(t, "", buildClientBaseURLPath(&dependencyEndpoint, "v1"),
		"return empty for no ApiPath nor ApiPaths")

	dependencyEndpoint.ApiPaths = []string{"/api/rbac-2", "/api/rbac-3"}
	assert.Equal(t, "/api/rbac-2/v1", buildClientBaseURLPath(&dependencyEndpoint, "v1"),
		"check ApiPaths has priority and use the first item")
}

func TestHasEndpointPort(t *testing.T) {
	assert.False(t, hasEndpointPort(nil))
	assert.False(t, hasEndpointPort(&clowder.DependencyEndpoint{Port: 0}))
	assert.True(t, hasEndpointPort(&clowder.DependencyEndpoint{Port: 8080}))
}

func TestHasEndpointTLSPort(t *testing.T) {
	assert.False(t, hasEndpointTLSPort(nil))
	assert.False(t, hasEndpointTLSPort(&clowder.DependencyEndpoint{TlsPort: nil}))
	assert.False(t, hasEndpointTLSPort(&clowder.DependencyEndpoint{TlsPort: pointy.Int(0)}))
	assert.True(t, hasEndpointTLSPort(&clowder.DependencyEndpoint{TlsPort: pointy.Int(8090)}))
}

func TestCheckPathInList(t *testing.T) {
	assert.True(t, checkPathInList("/cdapp/certs", "/etc/ssl/certs:/etc/pki/tls/certs:/system/etc/security/cacerts:/cdapp/certs"))
	assert.False(t, checkPathInList("/cdapp/certs", "/etc/ssl/certs:/etc/pki/tls/certs:/system/etc/security/cacerts"))

	assert.False(t, checkPathInList("/cdapp/certs", ""))
	assert.True(t, checkPathInList("", "/etc/ssl/certs:/etc/pki/tls/certs:/system/etc/security/cacerts"))
	assert.True(t, checkPathInList("", ""))
}

func TestInitSSLCertDir(t *testing.T) {
	const (
		pathCerts      = "/cdapp/certs"
		sampleCertDirs = "/etc/ssl/certs:/etc/pki/tls/certs:/system/etc/security/cacerts"
	)
	var clowderConfig *clowder.AppConfig

	// Backup/Restore env value
	oldValue := os.Getenv(EnvSSLCertDirectory)
	defer func() {
		os.Setenv(EnvSSLCertDirectory, oldValue) //nolint:all
	}()

	os.Setenv(EnvSSLCertDirectory, "") //nolint:all
	assert.False(t, hasTLSCAPath(clowderConfig))
	assert.NotPanics(t, func() {
		initSSLCertDir(clowderConfig)
	})
	assert.Equal(t, "", os.Getenv(EnvSSLCertDirectory),
		"No changes when hasTLSCAPath is false")

	os.Setenv(EnvSSLCertDirectory, "") //nolint:all
	clowderConfig = &clowder.AppConfig{
		TlsCAPath: pointy.String(pathCerts),
	}
	assert.NotPanics(t, func() {
		initSSLCertDir(clowderConfig)
	})
	assert.Equal(t, pathCerts, os.Getenv(EnvSSLCertDirectory),
		"Update for an empty SSL_CERT_DIR")

	os.Setenv(EnvSSLCertDirectory, sampleCertDirs) //nolint:all
	clowderConfig = &clowder.AppConfig{
		TlsCAPath: pointy.String(pathCerts),
	}
	assert.NotPanics(t, func() {
		initSSLCertDir(clowderConfig)
	})
	assert.Equal(t, sampleCertDirs+":"+pathCerts, os.Getenv(EnvSSLCertDirectory),
		"Update for a non empty SSL_CERT_DIR")
}

func TestHasTLSCAPath(t *testing.T) {
	assert.False(t, hasTLSCAPath(nil),
		"False when clowderConfig is nil")
	assert.False(t, hasTLSCAPath(&clowder.AppConfig{TlsCAPath: nil}),
		"False when clowderConfig.TlsCAPath is nil")
	assert.False(t, hasTLSCAPath(&clowder.AppConfig{TlsCAPath: pointy.String("")}),
		"False when clowderConfig.TlsCAPath is an empty string")
	assert.True(t, hasTLSCAPath(&clowder.AppConfig{TlsCAPath: pointy.String("/cdapp/certs")}),
		"True when clowderConfig.TlsCAPath is not an empty string")
}
