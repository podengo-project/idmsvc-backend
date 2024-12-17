# Logging

A contextual log is added to every request that allow
to link messages based on tracking the `request-id` field.

The request scoped log is stored in a `context.Context`
instance and we can recover the log instance from it by
using the primitive `LogFromCtx` function; for instance:

```go
    // Recover request log instance
    logger := app_context.LogFromCtx(ctx)

    // Print directly an error by
    logger.Error("my error message")

    // Print informative message by
    logger.Info("my action was successful")
```

The above add to every message the `request-id` recovered
from the request headers, and any other information aggregated
by the middlewares; this simplify the error message and allow
to focus on the message and any additional from the current
context.

## Guidelines

- Only the http handler will use the framework dependent
  context.
- Beyond the http handler, we will inject the `context.Context`
  instance because it is framework independent.
- For the deep levels, print out the error returned from the
  called library/framework. For higher levels, print out a
  description of the context where it happens. The reason is
  when we see the logs we will be told the history of
  what was wrong, starting with the root reason returned
  by the `err` instance.
- Enable the `location: true` to print out information about
  the exact location where the log message is printed out.
