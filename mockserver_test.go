package mockfs

import (
	"context"
	"testing"
	"time"

	firestore "cloud.google.com/go/firestore"
	proto "github.com/golang/protobuf/proto"
	ptypes "github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	assert "github.com/stretchr/testify/assert"
	pb "google.golang.org/genproto/googleapis/firestore/v1beta1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

var (
	aTime                = time.Date(2017, 1, 26, 0, 0, 0, 0, time.UTC)
	aTime2               = time.Date(2017, 2, 5, 0, 0, 0, 0, time.UTC)
	aTime3               = time.Date(2017, 3, 20, 0, 0, 0, 0, time.UTC)
	aTimestamp           = mustTimestampProto(aTime)
	aTimestamp2          = mustTimestampProto(aTime2)
	aTimestamp3          = mustTimestampProto(aTime3)
	writeResultForSet    = &firestore.WriteResult{UpdateTime: aTime}
	testData             = map[string]interface{}{"a": 1}
	testFields           = map[string]*pb.Value{"a": {ValueType: &pb.Value_IntegerValue{int64(1)}}}
	commitResponseForSet = &pb.CommitResponse{
		WriteResults: []*pb.WriteResult{{UpdateTime: aTimestamp}},
	}
)

func mustTimestampProto(t time.Time) *tspb.Timestamp {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return ts
}

func TestNewMockServer(t *testing.T) {
	assert := assert.New(t)

	server, err := newMockServer()
	assert.NotNil(server)
	assert.Nil(err)
}

// modified from https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/firestore/docref_test.go
func TestAddRPC(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	c, srv := New(t)
	dbPath := "projects/projectID/databases/(default)"

	path := "projects/projectID/databases/(default)/documents/C/a"
	pdoc := &pb.Document{
		Name:       path,
		CreateTime: aTimestamp,
		UpdateTime: aTimestamp,
		Fields:     map[string]*pb.Value{"f": {ValueType: &pb.Value_IntegerValue{int64(1)}}},
	}
	srv.AddRPC(&pb.BatchGetDocumentsRequest{
		Database:  dbPath,
		Documents: []string{path},
	}, []interface{}{
		&pb.BatchGetDocumentsResponse{
			Result:   &pb.BatchGetDocumentsResponse_Found{pdoc},
			ReadTime: aTimestamp2,
		},
	})
	ref := c.Collection("C").Doc("a")
	gotDoc, err := ref.Get(ctx)
	assert.Nil(err)
	if assert.NotNil(gotDoc) {
		assert.Equal(ref, gotDoc.Ref)
		assert.Equal(aTime, gotDoc.CreateTime)
		assert.Equal(aTime, gotDoc.UpdateTime)
		assert.Equal(aTime2, gotDoc.ReadTime)
	}

	path2 := "projects/projectID/databases/(default)/documents/C/b"
	srv.AddRPC(
		&pb.BatchGetDocumentsRequest{
			Database:  dbPath,
			Documents: []string{path2},
		}, []interface{}{
			&pb.BatchGetDocumentsResponse{
				Result:   &pb.BatchGetDocumentsResponse_Missing{path2},
				ReadTime: aTimestamp3,
			},
		})
	_, err = c.Collection("C").Doc("b").Get(ctx)
	assert.Equal(codes.NotFound, grpc.Code(err))
}

// modified from https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/firestore/collref_test.go
func TestAddRPCAdjust(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	c, srv := New(t)

	wantReq := &pb.CommitRequest{
		Database: "projects/projectID/databases/(default)",
		Writes: []*pb.Write{
			{
				Operation: &pb.Write_Update{
					Update: &pb.Document{
						Name:   "projects/projectID/databases/(default)/documents/C/d",
						Fields: testFields,
					},
				},
			},
		},
	}
	w := wantReq.Writes[0]
	w.CurrentDocument = &pb.Precondition{
		ConditionType: &pb.Precondition_Exists{false},
	}
	srv.AddRPCAdjust(wantReq, commitResponseForSet, func(gotReq proto.Message) {
		// We can't know the doc ID before Add is called, so we take it from
		// the request.
		w.Operation.(*pb.Write_Update).Update.Name = gotReq.(*pb.CommitRequest).Writes[0].Operation.(*pb.Write_Update).Update.Name
	})
	_, wr, err := c.Collection("C").Add(ctx, testData)
	assert.Nil(err)
	assert.Equal(writeResultForSet, wr)
}
