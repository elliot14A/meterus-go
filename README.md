# Meterus Go Client Documentation

## Introduction

The Meterus Go client is a Go package that provides a convenient way to interact with the Meterus gRPC server, a metering and billing solution. This client simplifies the process of sending events, managing meters, subjects, and querying data from the Meterus service.

## Installation

To use the Meterus Go client in your project, you need to install it using Go modules. Add the following import to your Go file:

```go
import "github.com/elliot14A/meterus-go/client"
```

Then run:

```
go mod tidy
```

## Connecting to Meterus

To create a new Meterus client, use the `NewMeterusClient` function:

```go
meterusClient, err := client.NewMeterusClient("address:port", "your-api-key")
if err != nil {
    // Handle error
}
defer meterusClient.Close()
```

Replace `"address:port"` with the address of your Meterus server and `"your-api-key"` with your Meterus API key.

## Core Concepts

### CloudEvent

The `CloudEvent` struct represents an event in the Meterus system. It contains fields such as Id, Source, SpecVersion, Type, Time, Subject, and Data.

### Services

The client is now organized into different services:

- MeteringService
- SubjectService
- ValidationService

## Basic Usage

### Metering Service

#### Creating a Metering Service

```go
meteringService := meterusClient.NewMeteringService()
```

#### Ingesting Events

```go
event, err := client.NewCloudEvent(
    // Meterus creates new id if empty id is passed
    "",
    "event-source",
    "1.0",
    "event-type",
    time.Now(),
    "event-subject",
    map[string]any{"key": "value"},
)
if err != nil {
    // Handle error
}

err = meteringService.Ingest(context.Background(), event)
if err != nil {
    // Handle error
}
```

#### Listing Meters

```go
response, err := meteringService.ListMeters(context.Background(), 10, 1)
if err != nil {
    // Handle error
}
// Process the response
```

#### Creating a Meter

```go
meter, err := meteringService.CreateMeter(context.Background(), &meter.CreateMeterRequest{
    Slug:        "new-meter",
    Description: proto.String("A new meter"),
    Aggregation: meter.Aggregation_AGGREGATION_COUNT,
    // Set other fields as needed
})
if err != nil {
    // Handle error
}
// Use the created meter
```

#### Querying a Meter

```go
response, err := meteringService.QueryMeter(context.Background(), &meter.QueryMeterRequest{
    MeterIdOrSlug: "meter-id-or-slug",
    From:          timestamppb.Now(),
    To:            timestamppb.Now(),
    // Set other query parameters
})
if err != nil {
    // Handle error
}
// Process the query response
```

### Subject Service

#### Creating a Subject Service

```go
subjectService := meterusClient.NewSubjectService()
```

#### Creating a Subject

```go
subject, err := subjectService.Create(context.Background(), "subject-id", proto.String("Display Name"))
if err != nil {
    // Handle error
}
// Use the created subject
```

#### Listing Subjects

```go
subjects, err := subjectService.ListById(context.Background(), 1, 10)
if err != nil {
    // Handle error
}
// Process the subjects
```

### Validation Service

#### Creating a Validation Service

```go
validationService := meterusClient.NewValidationService()
```

#### Validating an API Key

```go
isValid, subject, additionalAttributes, err := validationService.ValidateApiKey(context.Background(), "your-api-key", []string{"required-scope"})
if err != nil {
    // Handle error
}
// Use the validation results
```

## Advanced Usage

### Custom gRPC Dial Options

You can pass custom gRPC dial options when creating a new client:

```go
opts := []grpc.DialOption{
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    // Add other options as needed
}
client, err := client.NewMeterusClient("address:port", "your-api-key", opts...)
```

### Error Handling

All methods that communicate with the Meterus server return errors. Always check and handle these errors appropriately in your application.

### Context Usage

The client methods accept a `context.Context` parameter. Use this to set timeouts, deadlines, or cancel operations:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Use the context in client method calls
response, err := meteringService.ListMeters(ctx, 10, 1)
```

## Best Practices

1. Always close the client connection when you're done using it.
2. Use appropriate timeouts and contexts to prevent long-running operations.
3. Handle errors returned by the client methods.
4. Use the `NewCloudEvent` helper function to create properly formatted CloudEvents.
5. Organize your code by using the specific services (MeteringService, SubjectService, ValidationService) for related operations.

## Conclusion

The Meterus Go client provides a simple and efficient way to interact with the Meterus metering and billing service. By following this documentation, you should be able to integrate Meterus into your Go applications effectively.
For more detailed information about the Meterus service itself, please refer to the [Meterus GitHub repository](https://github.com/factly/meterus).
