// The scope of this file is:
// - Define the configuration struct.
// - Set default configuration values.
// - Map the data so viper can load the configuration there.
// See: https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
// See: https://consoledot.pages.redhat.com/docs/dev/getting-started/migration/config.html
package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/secrets"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
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
	// HostconfJWKs expire after 90 days and get renewed when the last
	// token expires in less than 30 days.
	DefaultHostconfJwkValidity         = time.Duration(90 * 24 * time.Hour)
	DefaultHostconfJwkRenewalThreshold = time.Duration(30 * 24 * time.Hour)
	// DefaultWebPort is the default port where the public API is listening
	DefaultWebPort = 8000
	// DefaultEnableRBAC is true
	DefaultEnableRBAC = true

	// DefaultDatabaseMaxOpenConn is the default for max open database connections
	DefaultDatabaseMaxOpenConn = 30

	// DefaultIdleTimeout 5 mins by default
	DefaultIdleTimeout = time.Duration(5 * time.Minute)
	// DefaultReadTimeout 3 seconds by default
	DefaultReadTimeout = time.Duration(3 * time.Second)
	// DefaultWriteTimeout 3 seconds by default
	DefaultWriteTimeout = time.Duration(3 * time.Second)

	// https://github.com/project-koku/koku/blob/main/koku/api/common/pagination.py

	// PaginationDefaultLimit is the default limit for the pagination
	PaginationDefaultLimit = 10
	// PaginationMaxLimit is the default max limit for the pagination
	PaginationMaxLimit = 1000

	// DefaultAcceptXRHFakeIdentity is disabled
	DefaultAcceptXRHFakeIdentity = false
	// DefaultValidateAPI is true
	DefaultValidateAPI = true

	// EnvSSLCertDirectory environment variable that provides
	// the paths for the CA certificates
	EnvSSLCertDirectory = "SSL_CERT_DIR"
)

var (
	// DefaultSizeLimitRequestHeader in bytes. Default 32KB
	DefaultSizeLimitRequestHeader = (32 * 1024)
	// DefaultSizeLimitRequestBody in bytes. Default 128KB
	DefaultSizeLimitRequestBody = (128 * 1024)
)

type Config struct {
	Loaded      bool
	Web         Web
	Database    Database
	Logging     Logging
	Kafka       Kafka
	Metrics     Metrics
	Clients     Clients
	Application Application `mapstructure:"app"`
	// Secrets is an untagged field and filled out on load
	Secrets secrets.AppSecrets `mapstructure:"-" json:"-"`
}

type Web struct {
	Port int16
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string `json:"-"`
	Name     string
	// https://stackoverflow.com/questions/54844546/how-to-unmarshal-golang-viper-snake-case-values
	CACertPath   string `mapstructure:"ca_cert_path"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type Cloudwatch struct {
	Region  string
	Key     string
	Secret  string `json:"-"`
	Session string
	Group   string
	Stream  string
}

type Logging struct {
	Level      string
	Console    bool
	Location   bool
	Type       string
	Cloudwatch Cloudwatch
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
		Password  string `json:"-"`
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

// Clients gather all the configuration to properly setup
// the third party services that idmsvc need to interact with.
type Clients struct {
	// RbacBaseURL is the base endpoint to launch RBAC requests.
	RbacBaseURL string `mapstructure:"rbac_base_url"`
	// PendoBaseURL is the base url to reach out the pendo API.
	PendoBaseURL string `mapstructure:"pendo_base_url"`
	// PendoAPIKey indicates the shared key to communicate with the API.
	PendoAPIKey string `mapstructure:"pendo_api_key" json:"-"`
	// PendoTrackEventKey indicates the shared key to communicate with the API
	// for track events.
	PendoTrackEventKey string `mapstructure:"pendo_track_event_key" json:"-"`
	// PendoRequestTimeoutSecs indicates the timeout for every request.
	PendoRequestTimeoutSecs int `mapstructure:"pendo_request_timeout_secs"`
}

// Application hold specific application settings
type Application struct {
	// Name is the internal application name
	Name string `validate:"required"`
	// API URL's path prefix, e.g. /api/idmsvc/v1
	PathPrefix string `mapstructure:"url_path_prefix" validate:"required"`
	// This is the default expiration time for the token
	// generated when a RHEL IDM domain is created
	TokenExpirationTimeSeconds int `mapstructure:"token_expiration_seconds" validate:"gte=600,lte=86400"`
	// Expiration and renewal duration for hostconf JWKs
	// TODO: short gte for local testing
	HostconfJwkValidity         time.Duration `mapstructure:"hostconf_jwk_validity" validate:"gte=1m,lte=8760h"`
	HostconfJwkRenewalThreshold time.Duration `mapstructure:"hostconf_jwk_renewal_threshold" validate:"gte=1m,lte=2160h"`
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
	// token and encrypted private JWKs.
	// Secrets are derived with HKDF-SHA256.
	MainSecret string `mapstructure:"secret" validate:"required,base64rawurl" json:"-"`
	// Flag to enable/disable rbac
	EnableRBAC bool `mapstructure:"enable_rbac"`
	// IdleTimeout for the API endpoints.
	IdleTimeout time.Duration `mapstructure:"idle_timeout" validate:"gte=1ms,lte=5m"`
	// ReadTimeout for the API endpoints.
	ReadTimeout time.Duration `mapstructure:"read_timeout" validate:"gte=1ms,lte=10s"`
	// WriteTimeout for the API endpoints.
	WriteTimeout time.Duration `mapstructure:"write_timeout" validate:"gte=1ms,lte=10s"`
	// SizeLimitRequestHeader for the API endpoints.
	SizeLimitRequestHeader int `mapstructure:"size_limit_request_header"`
	// SizeLimitRequestBody for the API endpoints.
	SizeLimitRequestBody int `mapstructure:"size_limit_request_body"`
}

var config *Config = nil

func DefaultCloudwatchStream() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "idmsvc"
	}
	return hostname
}

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
	v.SetDefault("database.max_open_conns", DefaultDatabaseMaxOpenConn)

	// Kafka
	addEventConfigDefaults(v)

	// Miscelanea
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.console", true)
	v.SetDefault("logging.location", false)
	v.SetDefault("logging.type", "null")

	// Cloudwatch
	v.SetDefault("logging.cloudwatch.region", "")
	v.SetDefault("logging.cloudwatch.group", "")
	v.SetDefault("logging.cloudwatch.stream", DefaultCloudwatchStream())
	v.SetDefault("logging.cloudwatch.key", "")
	v.SetDefault("logging.cloudwatch.secret", "")
	v.SetDefault("logging.cloudwatch.session", "")

	// Clients
	v.SetDefault("clients.rbac_base_url", "")
	v.SetDefault("clients.pendo_base_url", "")
	v.SetDefault("clients.pendo_api_key", "")
	v.SetDefault("clients.pendo_track_event_key", "")
	v.SetDefault("clients.pendo_request_timeout_secs", 0)

	// Application specific

	// Set default value for application expiration time for
	// the token created by the RHEL IDM domains
	v.SetDefault("app.token_expiration_seconds", DefaultTokenExpirationTimeSeconds)
	v.SetDefault("app.hostconf_jwk_validity", DefaultHostconfJwkValidity)
	v.SetDefault("app.hostconf_jwk_renewal_threshold", DefaultHostconfJwkRenewalThreshold)
	v.SetDefault("app.pagination_default_limit", PaginationDefaultLimit)
	v.SetDefault("app.pagination_max_limit", PaginationMaxLimit)
	v.SetDefault("app.accept_x_rh_fake_identity", DefaultAcceptXRHFakeIdentity)
	v.SetDefault("app.validate_api", DefaultValidateAPI)
	v.SetDefault("app.enable_rbac", DefaultEnableRBAC)
	v.SetDefault("app.url_path_prefix", DefaultPathPrefix)
	v.SetDefault("app.secret", "")
	v.SetDefault("app.debug", false)

	// Timeouts and limits
	v.SetDefault("app.idle_timeout", DefaultIdleTimeout)
	v.SetDefault("app.read_timeout", DefaultReadTimeout)
	v.SetDefault("app.write_timeout", DefaultWriteTimeout)
	v.SetDefault("app.size_limit_request_header", DefaultSizeLimitRequestHeader)
	v.SetDefault("app.size_limit_request_body", DefaultSizeLimitRequestBody)
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
	v.Set("logging.type", clowderConfig.Logging.Type)
	if clowderConfig.Logging.Cloudwatch != nil {
		v.Set("logging.cloudwatch.key", clowderConfig.Logging.Cloudwatch.AccessKeyId)
		v.Set("logging.cloudwatch.group", clowderConfig.Logging.Cloudwatch.LogGroup)
		v.Set("logging.cloudwatch.region", clowderConfig.Logging.Cloudwatch.Region)
		v.Set("logging.cloudwatch.secret", clowderConfig.Logging.Cloudwatch.SecretAccessKey)
		// TODO Delete the block below when the below PR is merged
		// See: https://github.com/RedHatInsights/clowder/pull/627
		if clowderConfig.Logging.Type == "" && len(clowderConfig.Logging.Cloudwatch.SecretAccessKey) > 0 {
			v.Set("logging.type", "cloudwatch")
		}
	}

	// Metrics configuration
	v.Set("metrics.path", clowderConfig.MetricsPath)
	v.Set("metrics.port", clowderConfig.MetricsPort)

	// Override client base url configuration from clowder when available
	updateServiceBasePath("rbac", "v1", v, clowderConfig)
}

// guardUpdateServiceBasePath raise a panic when some of the arguments for
// updateServiceBasePath is not provided.
func guardUpdateServiceBasePath(serviceName, version string, target *viper.Viper, clowderConfig *clowder.AppConfig) {
	if serviceName == "" {
		panic("'serviceName' is an empty string")
	}
	if version == "" {
		panic("'version' is an empty string")
	}
	if target == nil {
		panic("'target' is nil")
	}
	if clowderConfig == nil {
		panic("'clowderConfig' is nil")
	}
}

// updateServiceBasePath overrides the client base url when an endpoint is
// found for the serviceName, then try to build the base url, and if success
// overrides the base url with the value from clowder configuration.
// serviceName is the service we want to update the endpoint (it should exists
// in the dependencies field of clowderApp).
// target is the viper instance where the new base url calculated will be written
// if the endpoint exists in the configuration.
// clowderConfig represent the configuration injected for clowder which is the
// source of information to update the base url.
func updateServiceBasePath(serviceName, version string, target *viper.Viper, clowderConfig *clowder.AppConfig) {
	guardUpdateServiceBasePath(serviceName, version, target, clowderConfig)
	paramPath := "clients." + serviceName + "_base_url"
	if serviceEndpoint := getEndpoint(serviceName, clowderConfig); serviceEndpoint != nil {
		if serviceBaseURLString := buildClientBaseURL(serviceEndpoint, version); serviceBaseURLString != "" {
			slog.Debug("override base url for '" + serviceName + "' service to '" + serviceBaseURLString + "' from clowder endpoints")
			target.Set(paramPath, serviceBaseURLString)
			return
		}
	}
	slog.Debug("base url for '" + serviceName + "' service to '" + target.GetString(paramPath) + "'")
}

// getEndpoint search for the serviceName in the slice of Endpoints.
// Return nil when not found, else the reference to the clowder.DependencyEndpoint
// which match the serviceName criteria.
func getEndpoint(serviceName string, clowderConfig *clowder.AppConfig) *clowder.DependencyEndpoint {
	for _, ep := range clowderConfig.Endpoints {
		if ep.App == serviceName {
			return &ep
		}
	}
	return nil
}

// buildClientBaseURL construct a string to be used as base url for the service
// represented by serviceEndpoint and the specific version.
// serviceEndpoint is a structure retrieved from clowder.AppConfig.Endpoints
// that match with the 3rd party service that this one depends on. See getEndpoint
// version is the string for the 3rd party service api version to use; it adds the final
// suffix to the returned base URL
// Return the base
func buildClientBaseURL(serviceEndpoint *clowder.DependencyEndpoint, version string) string {
	// No checks on arguments as it is expected to be called from higher level where
	// the checks have been actually made.

	serviceBaseURLString := buildClientBaseURLSchemaHostPort(serviceEndpoint)
	if serviceBaseURLString == "" {
		return ""
	}

	serviceBaseURLPath := buildClientBaseURLPath(serviceEndpoint, version)
	if serviceBaseURLPath == "" {
		return ""
	}

	return serviceBaseURLString + serviceBaseURLPath
}

// buildClientBaseURLSchemaHostPort build the schema, hostname and port
// to reach out for the given serviceEndpoint. If TLS port is defined,
// then the schema 'https' and the referenced port are used; if no TLS
// port is indicated 'http' port and the no TLS port are used.
// Return the first base url part '[http|https]://{hostname}:[TlsPort|Port]
func buildClientBaseURLSchemaHostPort(serviceEndpoint *clowder.DependencyEndpoint) string {
	// Add schema, hostname and port
	if hasEndpointTLSPort(serviceEndpoint) {
		return "https://" + serviceEndpoint.Hostname +
			":" + strconv.Itoa(*serviceEndpoint.TlsPort)
	} else if hasEndpointPort(serviceEndpoint) {
		return "http://" + serviceEndpoint.Hostname +
			":" + strconv.Itoa(serviceEndpoint.Port)
	}
	slog.Warn("No Port nor TLSPort found for the service endpoint",
		slog.String("service", serviceEndpoint.App))
	return ""
}

// buildClientBaseURLPath build the second part 'api path + version'
// If ApiPaths has some item, the first item is used, else if the deprecated
// ApiPath is set, it is used, else return ""
func buildClientBaseURLPath(serviceEndpoint *clowder.DependencyEndpoint, version string) string {
	if len(serviceEndpoint.ApiPaths) > 0 {
		// Use the first path in the slice
		return strings.TrimSuffix(serviceEndpoint.ApiPaths[0], "/") +
			"/" + version
	}
	return ""
}

// hasEndpointTLSPort return if the DependencyEndpoint has a valid TlsPort field
// that can be used.
// Return true if TlsPort can be used, else false.
func hasEndpointTLSPort(serviceEndpoint *clowder.DependencyEndpoint) bool {
	return serviceEndpoint != nil && serviceEndpoint.TlsPort != nil && *serviceEndpoint.TlsPort > 0
}

// hasEndpointPort return if the DependencyEndpoint has a valid Port field
// that can be used.
// Return true if Port can be used, else false.
func hasEndpointPort(serviceEndpoint *clowder.DependencyEndpoint) bool {
	return serviceEndpoint != nil && serviceEndpoint.Port > 0
}

func Load(cfg *Config) *viper.Viper {
	var err error

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
		slog.Debug("clowder is enabled")
		setClowderConfiguration(v, clowder.LoadedConfig)
	} else {
		slog.Debug("clowder not enabled")
	}

	if err = v.ReadInConfig(); err != nil {
		slog.Warn("Not using config.yaml", slog.Any("error", err))
	}
	if err = v.Unmarshal(cfg); err != nil {
		slog.Warn("Mapping to configuration", slog.Any("error", err))
	}

	return v
}

// Log logs the configuration
func (c *Config) Log(logger *slog.Logger) {
	logger.Info(
		"Configuration",
		slog.Group("Web",
			slog.Int("Port", int(c.Web.Port)),
		),
		slog.Group("Database",
			slog.String("Host", c.Database.Host),
			slog.Int("Port", c.Database.Port),
			slog.String("User", c.Database.User),
			slog.String("Password", obfuscateSecret(c.Database.Password)),
			slog.String("Name", c.Database.Name),
			slog.String("CACertPath", c.Database.CACertPath),
			slog.Int("MaxOpenConns", c.Database.MaxOpenConns),
		),
		slog.Group("Logging",
			slog.String("Level", c.Logging.Level),
			slog.Bool("Console", c.Logging.Console),
			slog.Bool("Location", c.Logging.Location),
			slog.String("Type", c.Logging.Type),
			slog.Group("Cloudwatch",
				slog.String("Region", c.Logging.Cloudwatch.Region),
				slog.String("Key", c.Logging.Cloudwatch.Key),
				slog.String("Secret", obfuscateSecret(c.Logging.Cloudwatch.Secret)),
				slog.String("Session", c.Logging.Cloudwatch.Session),
				slog.String("Group", c.Logging.Cloudwatch.Group),
				slog.String("Stream", c.Logging.Cloudwatch.Stream),
			),
		),
		slog.Group("Metrics",
			slog.String("Path", c.Metrics.Path),
			slog.Int("Port", c.Metrics.Port),
		),
		slog.Group("Clients",
			slog.String("RbacBaseURL", c.Clients.RbacBaseURL),
			slog.String("PendoBaseURL", c.Clients.PendoBaseURL),
			slog.String("PendoAPIKey", obfuscateSecret(c.Clients.PendoAPIKey)),
			slog.String("PendoTrackEventKey", obfuscateSecret(c.Clients.PendoTrackEventKey)),
			slog.Int("PendoRequestTimeoutSecs", c.Clients.PendoRequestTimeoutSecs),
		),
		slog.Group("Application",
			slog.String("Name", c.Application.Name),
			slog.String("PathPrefix", c.Application.PathPrefix),
			slog.Int("TokenExpirationTimeSeconds", c.Application.TokenExpirationTimeSeconds),
			slog.Duration("HostconfJwkValidity", c.Application.HostconfJwkValidity),
			slog.Duration("HostconfJwkRenewalThreshold", c.Application.HostconfJwkRenewalThreshold),
			slog.Int("PaginationDefaultLimit", c.Application.PaginationDefaultLimit),
			slog.Int("PaginationMaxLimit", c.Application.PaginationMaxLimit),
			slog.Bool("AcceptXRHFakeIdentity", c.Application.AcceptXRHFakeIdentity),
			slog.Bool("ValidateAPI", c.Application.ValidateAPI),
			slog.String("MainSecret", obfuscateSecret(c.Application.MainSecret)),
			slog.Bool("EnableRBAC", c.Application.EnableRBAC),
			slog.Duration("IdleTimeout", c.Application.IdleTimeout),
			slog.Duration("ReadTimeout", c.Application.ReadTimeout),
			slog.Duration("WriteTimeout", c.Application.WriteTimeout),
			slog.Int("SizeLimitRequestHeader", c.Application.SizeLimitRequestHeader),
			slog.Int("SizeLimitRequestBody", c.Application.SizeLimitRequestBody),
		),
	)
}

func obfuscateSecret(value string) string {
	if value == "" {
		return ""
	}
	return "***"
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
	_ = Load(config)

	if err := Validate(config); err != nil {
		reportError(err)
		panic("Invalid configuration")
	}

	sec, err := secrets.NewAppSecrets(config.Application.MainSecret)
	if err != nil {
		panic(err)
	}
	config.Secrets = *sec

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

// initSSLCertDir update SSL_CERT_DIR to add TLSCAPath if not found
// in the list of directories.
// clowderConfig
func initSSLCertDir(clowderConfig *clowder.AppConfig) {
	envSSLCertDir := os.Getenv(EnvSSLCertDirectory)
	if hasTLSCAPath(clowderConfig) {
		if !checkPathInList(*clowderConfig.TlsCAPath, envSSLCertDir) {
			if envSSLCertDir == "" {
				envSSLCertDir = *clowderConfig.TlsCAPath
			} else {
				envSSLCertDir = envSSLCertDir + ":" + *clowderConfig.TlsCAPath
			}
			if err := os.Setenv(EnvSSLCertDirectory, envSSLCertDir); err != nil {
				panic(err.Error())
			}
			slog.Info(EnvSSLCertDirectory + " was updated")
			return
		}
	}
	slog.Info(EnvSSLCertDirectory + " not updated")
}

// hasTLSCAPath check the condition when TlsCAPath has
// some content to use.
func hasTLSCAPath(clowderConfig *clowder.AppConfig) bool {
	return clowderConfig != nil && clowderConfig.TlsCAPath != nil && *clowderConfig.TlsCAPath != ""
}

func checkPathInList(path, pathList string) bool {
	if path == "" {
		// Nothing to check
		return true
	}
	path = strings.TrimSuffix(path, "/")
	pathListItems := strings.Split(pathList, ":")
	for _, item := range pathListItems {
		if strings.TrimSuffix(item, "/") == path {
			return true
		}
	}
	return false
}

//nolint:all
func init() {
	// NOTE
	//
	// Linter disabled to allow this exception
	// I do not recommend the leverage of func init
	// because provide "black magic" that could evoke
	// not expected behaviors difficult to debug. Be
	// aware their execution happens before the
	// first line of the `main` function is reached out.
	initSSLCertDir(clowder.LoadedConfig)
}
