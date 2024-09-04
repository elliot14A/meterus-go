# Meterus Go Client Documentation

## Introduction

The Meterus Go client is a Go package that provides a convenient way to interact with the Meterus gRPC server, a metering and billing solution. This client simplifies the process of sending events, managing meters, and querying data from the Meterus service.

## Installation

To use the Meterus Go client in your project, you need to install it using Go modules. Add the following import to your Go file:

```go
import "github.com/factly/meterus/meterus-go"
```

Then run:

```
go mod tidy
```

## Connecting to Meterus

To create a new Meterus client, use the `NewMeterusClient` function:

```go
client, err := client.NewMeterusClient("address:port", "your-api-key")
if err != nil {
    // Handle error
}
defer client.Close()
```

Replace `"address:port"` with the address of your Meterus server and `"your-api-key"` with your Meterus API key.

## Core Concepts

### CloudEvent

The `CloudEvent` struct represents an event in the Meterus system. It contains the following fields:

- `Id`: Unique identifier for the event
- `Source`: The source of the event
- `SpecVersion`: The version of the CloudEvents spec
- `Type`: The type of the event
- `Time`: The time the event occurred
- `Subject`: The subject of the event
- `Data`: Additional data associated with the event

## Basic Usage

### Ingesting Events

To ingest a CloudEvent into Meterus:

```go
event, err := client.NewCloudEvent(
    "event-id",
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

err = client.Ingest(context.Background(), event)
if err != nil {
    // Handle error
}
```

### Listing Meters

To retrieve a list of meters:

```go
response, err := client.ListMeters(context.Background(), 10, 1)
if err != nil {
    // Handle error
}
// Process the response
```

### Creating a Meter

To create a new meter:

```go
meter, err := client.CreateMeter(context.Background(), &pb.CreateMeterRequest{
    Slug:        "new-meter",
    Description: proto.String("A new meter"),
    Aggregation: pb.Aggregation_AGGREGATION_COUNT,
    // Set other fields as needed
})
if err != nil {
    // Handle error
}
// Use the created meter
```

### Querying a Meter

To query data from a meter:

```go
response, err := client.QueryMeter(context.Background(), &pb.QueryMeterRequest{
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
response, err := client.ListMeters(ctx, 10, 1)
```

## Best Practices

1. Always close the client connection when you're done using it.
2. Use appropriate timeouts and contexts to prevent long-running operations.
3. Handle errors returned by the client methods.
4. Use the `NewCloudEvent` helper function to create properly formatted CloudEvents.

## Conclusion

The Meterus Go client provides a simple and efficient way to interact with the Meterus metering and billing service. By following this documentation, you should be able to integrate Meterus into your Go applications effectively.

For more detailed information about the Meterus service itself, please refer to the [Meterus GitHub repository](https://github.com/factly/meterus).
