package client_test

import (
	"context"
	"os"
	"testing"
	"time"

	client "github.com/elliot14A/meterus-go"
	pb "github.com/elliot14A/meterus-go/meters/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	testClient *client.MeterusClient
	testMeter  *pb.Meter
)

func TestMain(m *testing.M) {
	// Setup
	addr := os.Getenv("METERUS_ADDR")
	if addr == "" {
		addr = "localhost:8000" // Default address
	}
	apiKey := os.Getenv("METERUS_API_KEY")
	if apiKey == "" {
		panic("METERUS_API_KEY environment variable is not set")
	}

	var err error
	testClient, err = client.NewMeterusClient(addr, apiKey)
	if err != nil {
		panic("Failed to create Meterus client: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Teardown
	if testClient != nil {
		testClient.Close()
	}

	os.Exit(code)
}

func TestCreateAndGetMeter(t *testing.T) {
	ctx := context.Background()

	// Create a new meter
	createReq := &pb.CreateMeterRequest{
		Slug:        "test_meter",
		Description: stringPtr("Test meter for integration tests"),
		Aggregation: pb.Aggregation_AGGREGATION_COUNT,
		GroupBy:     []string{"user_id"},
		CreatedBy:   "integration-test",
		EventType:   "test.event",
	}

	createdMeter, err := testClient.CreateMeter(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, createdMeter)

	// Store the created meter for other tests
	testMeter = createdMeter

	// Get the created meter
	fetchedMeter, err := testClient.GetMeter(ctx, createdMeter.Id)
	require.NoError(t, err)
	require.NotNil(t, fetchedMeter)

	assert.Equal(t, createdMeter.Id, fetchedMeter.Id)
	assert.Equal(t, createReq.Slug, fetchedMeter.Slug)
	assert.Equal(t, *createReq.Description, *fetchedMeter.Description)
	assert.Equal(t, createReq.Aggregation, fetchedMeter.Aggregation)
	assert.Equal(t, createReq.GroupBy, fetchedMeter.GroupBy)
	assert.Equal(t, createReq.CreatedBy, fetchedMeter.CreatedBy)
	assert.Equal(t, createReq.EventType, fetchedMeter.EventType)
}

func TestListMeters(t *testing.T) {
	ctx := context.Background()

	// Ensure we have at least one meter (created in the previous test)
	require.NotNil(t, testMeter, "TestCreateAndGetMeter should be run before this test")

	// List meters
	response, err := testClient.ListMeters(ctx, 10, 1)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.NotEmpty(t, response.Meters)
	assert.Contains(t, response.Meters, testMeter)
}

func TestIngestAndQueryMeter(t *testing.T) {
	ctx := context.Background()

	// Ensure we have a meter to ingest events for
	require.NotNil(t, testMeter, "TestCreateAndGetMeter should be run before this test")

	// Create and ingest a test event
	eventData, err := structpb.NewStruct(map[string]interface{}{
		"user_id": "test-user",
		"action":  "test-action",
	})
	require.NoError(t, err)

	event := &pb.CloudEvent{
		Id:          "test-event-id",
		Source:      "integration-test",
		SpecVersion: "1.0",
		Type:        testMeter.EventType,
		Time:        timestamppb.Now(),
		Subject:     "test-subject",
		Data:        eventData,
	}

	err = testClient.Ingest(ctx, event)
	require.NoError(t, err)

	// Wait a bit for the event to be processed
	time.Sleep(2 * time.Second)

	// Query the meter
	now := time.Now()
	queryReq := &pb.QueryMeterRequest{
		MeterIdOrSlug: testMeter.Id,
		From:          timestamppb.New(now.Add(-1 * time.Hour)),
		To:            timestamppb.New(now),
		Subject:       []string{"test-subject"},
		GroupBy:       []string{"user_id"},
	}

	queryResp, err := testClient.QueryMeter(ctx, queryReq)
	require.NoError(t, err)
	require.NotNil(t, queryResp)

	assert.NotEmpty(t, queryResp.Data)
	assert.Equal(t, float64(1), queryResp.Data[0].Value) // We expect a count of 1 for our single event
}

// func TestListMeterSubjects(t *testing.T) {
// 	ctx := context.Background()
//
// 	// Ensure we have a meter with subjects
// 	require.NotNil(t, testMeter, "TestCreateAndGetMeter and TestIngestAndQueryMeter should be run before this test")
//
// 	// List subjects
// 	response, err := testClient.ListMeterSubjects(ctx, testMeter.Id)
// 	require.NoError(t, err)
// 	require.NotNil(t, response)
//
// 	assert.Contains(t, response.Subjects, "test-subject")
// }

func TestDeleteMeter(t *testing.T) {
	ctx := context.Background()

	// Ensure we have a meter to delete
	require.NotNil(t, testMeter, "TestCreateAndGetMeter should be run before this test")

	// Delete the meter
	err := testClient.DeleteMeter(ctx, testMeter.Slug)
	require.NoError(t, err)

	// Try to get the deleted meter (should fail)
	_, err = testClient.GetMeter(ctx, testMeter.Slug)
	assert.Error(t, err)
}

func TestValidateApiKey(t *testing.T) {
	ctx := context.Background()

	// Test with valid API key
	response, err := testClient.ValidateApiKey(ctx, os.Getenv("METERUS_API_KEY"), []string{"meterus:meterus"})
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.NotNil(t, response.Metadata)

	// Test with invalid API key
	_, err = testClient.ValidateApiKey(ctx, "invalid-api-key", []string{"read", "write"})
	assert.Error(t, err)
}

func stringPtr(s string) *string {
	return &s
}
