module github.com/podengo-project/idmsvc-backend

// When updating go version, update the below too:
//   .github/workflows/main.yml
//   build/package/Dockerfile
go 1.20

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/achiku/planter v0.0.0-20230118102009-4847a8211e35
	github.com/atombender/go-jsonschema v0.13.1
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/deepmap/oapi-codegen v1.14.0
	github.com/getkin/kin-openapi v0.120.0
	github.com/go-playground/validator/v10 v10.16.0
	github.com/golang-migrate/migrate/v4 v4.16.2
	github.com/google/uuid v1.4.0
	github.com/labstack/echo/v4 v4.11.3
	github.com/labstack/gommon v0.4.1
	github.com/lestrrat-go/jwx/v2 v2.0.17
	github.com/lib/pq v1.10.9
	github.com/loov/goda v0.5.7
	github.com/mikefarah/yq/v4 v4.34.2
	github.com/oapi-codegen/runtime v1.1.0
	github.com/openlyinc/pointy v1.2.1
	github.com/prometheus/client_golang v1.17.0
	github.com/qri-io/jsonschema v0.2.1
	github.com/redhatinsights/app-common-go v1.6.7
	github.com/redhatinsights/platform-go-middlewares v0.20.0
	github.com/rs/zerolog v1.31.0
	github.com/spf13/viper v1.17.0
	github.com/stretchr/testify v1.8.4
	github.com/subpop/xrhidgen v0.0.0-20221013174335-ca910188f4c9
	github.com/vektra/mockery/v2 v2.38.0
	golang.org/x/crypto v0.15.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b
)

require (
	github.com/a8m/envsubst v1.4.2 // indirect
	github.com/alecthomas/kingpin v2.2.6+incompatible // indirect
	github.com/alecthomas/participle/v2 v2.1.0 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chigopher/pathlib v0.15.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/elliotchance/orderedmap v1.5.0 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.11.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.4 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/gomega v1.24.1 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/peterbourgon/ff/v3 v3.1.2 // indirect
	github.com/pioz/faker v1.7.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/redhatinsights/module-update-router v0.0.0-20220708181923-3d8d459361dd // indirect
	github.com/sagikazarmark/locafero v0.3.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sanity-io/litter v1.5.5 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sgreben/flagvar v1.10.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/term v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.4.0 // indirect
	golang.org/x/tools v0.15.0 // indirect
	golang.org/x/tools/go/vcs v0.1.0-deprecated // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
