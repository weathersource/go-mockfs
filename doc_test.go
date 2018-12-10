package mockfs_test

import (
	"context"
	"testing"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	assert "github.com/stretchr/testify/assert"
	mockfs "github.com/weathersource/go-mockfs"
	pb "google.golang.org/genproto/googleapis/firestore/v1beta1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

func Example_error() {
	var t *testing.T

	// Get a firestore client and mock firestore server
	client, server, err := mockfs.New()
	assert.NotNil(t, client)
	assert.NotNil(t, server)
	assert.Nil(t, err)

	// Populate a mock document "a" in collection "C"
	var (
		aTime         = time.Date(2017, 1, 26, 0, 0, 0, 0, time.UTC)
		aTimestamp, _ = ptypes.TimestampProto(aTime)
		dbPath        = "projects/projectID/databases/(default)"
		path          = "projects/projectID/databases/(default)/documents/C/a"
	)
	server.AddRPC(
		&pb.BatchGetDocumentsRequest{
			Database:  dbPath,
			Documents: []string{path},
		},
		[]interface{}{
			&pb.BatchGetDocumentsResponse{
				Result:   &pb.BatchGetDocumentsResponse_Missing{path},
				ReadTime: aTimestamp,
			},
		},
	)

	// Get document "a" in collection "C"
	_, err2 := client.Collection("C").Doc("a").Get(context.Background())

	// Test the response
	assert.Equal(t, codes.NotFound, grpc.Code(err2))
}

func Example_success() {
	var t *testing.T

	// Get a firestore client and mock firestore server
	client, server, err := mockfs.New()
	assert.NotNil(t, client)
	assert.NotNil(t, server)
	assert.Nil(t, err)

	// Populate a mock document "b" in collection "C"
	var (
		aTime          = time.Date(2017, 1, 26, 0, 0, 0, 0, time.UTC)
		aTime2         = time.Date(2017, 2, 5, 0, 0, 0, 0, time.UTC)
		aTimestamp, _  = ptypes.TimestampProto(aTime)
		aTimestamp2, _ = ptypes.TimestampProto(aTime2)
		dbPath         = "projects/projectID/databases/(default)"
		path           = "projects/projectID/databases/(default)/documents/C/b"
		pdoc           = &pb.Document{
			Name:       path,
			CreateTime: aTimestamp,
			UpdateTime: aTimestamp,
			Fields:     map[string]*pb.Value{"f": {ValueType: &pb.Value_IntegerValue{int64(1)}}},
		}
	)
	server.AddRPC(
		&pb.BatchGetDocumentsRequest{
			Database:  dbPath,
			Documents: []string{path},
		},
		[]interface{}{
			&pb.BatchGetDocumentsResponse{
				Result:   &pb.BatchGetDocumentsResponse_Found{pdoc},
				ReadTime: aTimestamp2,
			},
		},
	)

	// Get document "b" in collection "C"
	ref := client.Collection("C").Doc("b")
	gotDoc, err := ref.Get(context.Background())

	// Test the response
	assert.Nil(t, err)
	if assert.NotNil(t, gotDoc) {
		assert.Equal(t, ref, gotDoc.Ref)
		assert.Equal(t, aTime, gotDoc.CreateTime)
		assert.Equal(t, aTime, gotDoc.UpdateTime)
		assert.Equal(t, aTime2, gotDoc.ReadTime)
	}
}
