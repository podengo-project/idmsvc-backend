# Frameworks and libraries

- It uses golang 1.18 to match with hmscontent version.
- For the Service API it uses [echo](https://echo.labstack.com/)
  framework that is what was learned from hmscontent experience.
  > not tested the current state to build the boilerplate for
  > other frameworks like gin

  ```raw
  High performance, extensible, minimalist Go web framework
  ```

- For the logging system it uses the [slog](https://go.dev/blog/slog) library

- The database is using gorm: https://gorm.io/docs/index.html

- Unit Testing:
  - Testify: https://pkg.go.dev/github.com/stretchr/testify
  - Mockery: https://github.com/vektra/mockery
  - SqlMock: https://github.com/DATA-DOG/go-sqlmock

## Contents

- [Service API](01-service-api.md)
- [Kafka Events](02-event-api.md)
- [Metrics](03-metrics.md)
- [Infrastructure](04-infrastructure.md)
- [RBAC](05-rbac.md)
- [Metrics](06-metrics.md)
- [Configurations](07-configs.md)
- [Logs](08-logs.md)
