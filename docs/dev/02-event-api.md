# Event API

The Event API is implemented on top of Kafka infrastructure

> Just make aware that exist an ADR that provide a central place
> to keep all the event schemas, and to provide language bindings
> so that we update the schema once and we get the language bindings
> got using into the different services. That intent would be using
> [cloudevents](https://cloudevents.io).

## What we have

- Integration with kafka broker.
  - Design first approach (very likely can be enhanced).
  - Generate golang types from schema definition
  - Validate schemas consumed before dispatch on the matching handler.
  - Provide a general propose producer that can be used to integrate
    with existing kafka topics and the owned topics.
  - Provides a routed to dispatch the events.
  - Local infrastructure to check kafka messages before go to
    ephemeral environment, stage, prod.
  - Support to deploy the service in EE.
  - Support to connect to kafka in EE.
  - Support for SASL.

TODO:

- Serialize message to don't loose them if something fail.
- Job to recover not processed events or error
- Define schema for headers?
- Error handling for kafka consumers.
- Metrics for kafka events

## Some notes about kafka

- A group-id define a set of agents which work together
  to consume messages from one or several topics.
- Partitions define how many clients can connect to
  consume messages. This is a keypoint value for scaling
  the kafka consuming processes. Keep this in mind when
  automatic horizontal scaling is defined.
- One topic has one schema associated.
- x-rh-identity and x-rh-insights-request-id must be
  added to the produced message.
- The key and validate the consumed messages are as
  important as validate the request input for the
  http handler.

## I want to create a topic that own our service

- Define the new message schema for the topic at:
  `api/event`, and name it as the name of your topic.
- Generate types by: `make generate-event`. They will be
  generated at: `internal/api/event/` directory.
- Add the new structure at `internal/api/event/schemas.go`:

  ```golang
  const TopicMyTopicName = "platform.idmsvc.my-topic-name"
  ...
  var AllowedTopic = []string{
    TopicTodoCreated,
    TopicMyTopicName,
  }
  //go:embed "my_topic_name.event.json"
  var schemaEventMyTopicName string
  ...
  var (
    schemaKey2JsonSpec = map[string]string{
      TopicTodoCreated: schemaEventTodoCreated,
      TopicMyTopicName: schemaEventMyTopicName,
    }
  )
  ...
  ```

- Create a handler for it at: `internal/handler/impl/my_topic_name_event_handler.go`:

  ```golang
  type myTopicNameEventHandler struct {
	db *gorm.DB
  }

  func NewTodoCreatedEventHandler(db *gorm.DB) event.Eventable {
    if db == nil {
	  return nil
    }
    return &myTopicNameEventHandler{
	  db: db,
    }
  }

  func (h *myTopicNameEventHandler) OnMessage(msg *kafka.Message) error {
    return fmt.Errorf("Not implemented")
  }

  ```
- Add your event handler to the kafka router at: `internal/infrastructure/service/impl/kafka.go`:

  ```golang
  ...
  eventRouter.Add(api_event.TopicMyTopicName, impl.NewMyTopicNameEventHandler(s.db))
  ...
  ```

- Before promote to stage, you need to declare the topic
  at: https://github.com/RedHatInsights/platform-mq

## References

- https://consoledot.pages.redhat.com/docs/dev/services/kafka.html
- https://json-schema.org/specification.html
- https://github.com/RedHatInsights/platform-mq
  (see: https://github.com/RedHatInsights/platform-mq/pull/246/files)
