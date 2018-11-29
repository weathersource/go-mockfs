package mockfs

import (
	"context"
	"testing"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	assert "github.com/stretchr/testify/assert"
	mockfs "github.com/weathersource/go-mockfs"
	pb "google.golang.org/genproto/googleapis/firestore/v1beta1"
)

func Example_error() {
	var t *testing.T

	// Get a firestore client and mock firestore server
	client, server := mockfs.New(t)

	// Populate a mock document "b" in collection "C"
	var (
		aTime      = time.Date(2017, 1, 26, 0, 0, 0, 0, time.UTC)
		aTimestamp = ptypes.TimestampProto(aTime)
		dbPath     = "projects/projectID/databases/(default)"
		path       = "projects/projectID/databases/(default)/documents/C/b"
	)
	srv.AddRPC(
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
	_, err := client.Collection("C").Doc("b").Get(context.Background())

	// Test the response
	assert.Equal(codes.NotFound, grpc.Code(err))
}

func Example_success() {
	var t *testing.T

	// Get a firestore client and mock firestore server
	client, server := mockfs.New(t)

	// Populate a mock document "a" in collection "C"
	var (
		aTime          = time.Date(2017, 1, 26, 0, 0, 0, 0, time.UTC)
		aTime2         = time.Date(2017, 2, 5, 0, 0, 0, 0, time.UTC)
		aTimestamp, _  = ptypes.TimestampProto(aTime)
		aTimestamp2, _ = ptypes.TimestampProto(aTime2)
		dbPath         = "projects/projectID/databases/(default)"
		path           = "projects/projectID/databases/(default)/documents/C/a"
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

	// Get document "a" in collection "C"
	gotDoc, err := client.Collection("C").Doc("a").Get(context.Background())

	// Test the response
	assert.Nil(t, err)
	if assert.NotNil(t, ref, gotDoc) {
		assert.Equal(t, ref, gotDoc.Ref)
		assert.Equal(t, aTime, gotDoc.CreateTime)
		assert.Equal(t, aTime, gotDoc.UpdateTime)
		assert.Equal(t, aTime2, gotDoc.ReadTime)
	}
}
