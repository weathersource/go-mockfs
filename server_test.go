package mockfs

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
	assert "github.com/stretchr/testify/assert"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
)

var (
	aTime       = time.Date(2017, 1, 26, 0, 0, 0, 0, time.UTC)
	aTime2      = time.Date(2017, 2, 5, 0, 0, 0, 0, time.UTC)
	aTime3      = time.Date(2017, 3, 20, 0, 0, 0, 0, time.UTC)
	aTimestamp  = mustTimestampProto(aTime)
	aTimestamp2 = mustTimestampProto(aTime2)
	aTimestamp3 = mustTimestampProto(aTime3)
)

func mustTimestampProto(t time.Time) *tspb.Timestamp {
	ts := tspb.New(t)
	return ts
}

func TestNewServer(t *testing.T) {
	assert := assert.New(t)

	server, err := newServer()
	assert.NotNil(server)
	assert.Nil(err)
}

func TestReset(t *testing.T) {
	assert := assert.New(t)

	_, srv, err := New()
	assert.Nil(err)
	srv.AddRPC(
		&pb.BatchGetDocumentsRequest{},
		[]interface{}{},
	)
	srv.Reset()
	assert.Nil(srv.reqItems)
	assert.Nil(srv.resps)
}

// modified from https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/firestore/docref_test.go
func TestAddRPC(t *testing.T) {
	assert := assert.New(t)

	_, srv, err := New()
	assert.Nil(err)
	srv.AddRPC(
		&pb.BatchGetDocumentsRequest{},
		[]interface{}{},
	)
	assert.NotNil(srv.resps)
	assert.NotNil(srv.reqItems[0].wantReq)
	assert.Nil(srv.reqItems[0].adjust)
}

// modified from https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/firestore/docref_test.go
func TestAddRPCAdjust(t *testing.T) {
	assert := assert.New(t)

	_, srv, err := New()
	assert.Nil(err)
	srv.AddRPCAdjust(
		&pb.BatchGetDocumentsRequest{},
		[]interface{}{},
		func(req proto.Message) {},
	)
	assert.NotNil(srv.resps)
	assert.NotNil(srv.reqItems[0].wantReq)
	assert.NotNil(srv.reqItems[0].adjust)
}

func TestPopRPC(t *testing.T) {
	assert := assert.New(t)

	_, srv, err := New()
	assert.NotNil(srv)
	assert.Nil(err)

	// test no RPCs
	panicTest := func() {
		srv.popRPC(nil)
	}
	assert.Panics(panicTest)

	// test success adjust commit
	wantReq := &pb.CommitRequest{
		Database: "projects/projectID/databases/(default)",
		Writes: []*pb.Write{
			{
				Operation: &pb.Write_Update{
					Update: &pb.Document{
						Name:   "projects/projectID/databases/(default)/documents/C/d",
						Fields: map[string]*pb.Value{"a": {ValueType: &pb.Value_IntegerValue{int64(1)}}},
					},
				},
			},
		},
	}
	w := wantReq.Writes[0]
	w.CurrentDocument = &pb.Precondition{
		ConditionType: &pb.Precondition_Exists{false},
	}
	srv.AddRPCAdjust(
		wantReq,
		&pb.CommitResponse{
			WriteResults: []*pb.WriteResult{{UpdateTime: aTimestamp}},
		},
		func(gotReq proto.Message) {
			// We can't know the doc ID before Add is called, so we take it from the request.
			w.Operation.(*pb.Write_Update).Update.Name = gotReq.(*pb.CommitRequest).Writes[0].Operation.(*pb.Write_Update).Update.Name
		},
	)
	resp, err := srv.popRPC(wantReq)
	assert.Nil(err)
	assert.NotNil(resp)

	// test error non matching requests
	path := "projects/projectID/databases/(default)/documents/C/a"
	pdoc := &pb.Document{
		Name:       path,
		CreateTime: aTimestamp,
		UpdateTime: aTimestamp,
		Fields:     map[string]*pb.Value{"f": {ValueType: &pb.Value_IntegerValue{int64(1)}}},
	}
	srv.AddRPC(
		&pb.BatchGetDocumentsRequest{
			Database:  "projects/projectID/databases/(default)",
			Documents: []string{"projects/projectID/databases/(default)/documents/C/a"},
		}, []interface{}{
			&pb.BatchGetDocumentsResponse{
				Result:   &pb.BatchGetDocumentsResponse_Found{pdoc},
				ReadTime: aTimestamp2,
			},
		},
	)
	_, err = srv.popRPC(&pb.BatchGetDocumentsRequest{
		Database:  "projects/projectID/databases/(default)",
		Documents: []string{"projects/projectID/databases/(default)/documents/C/b"},
	})
	assert.NotNil(err)
}

// type byFieldPath []*pb.DocumentTransform_FieldTransform
// func (a byFieldPath) Len() int           { return len(a) }
// func (a byFieldPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a byFieldPath) Less(i, j int) bool { return a[i].FieldPath < a[j].FieldPath }

func TestLen(t *testing.T) {
	bfp := byFieldPath{
		&pb.DocumentTransform_FieldTransform{FieldPath: "a"},
		&pb.DocumentTransform_FieldTransform{FieldPath: "b"},
	}
	assert.Equal(t, 2, bfp.Len())
}

func TestSwap(t *testing.T) {
	bfp := byFieldPath{
		&pb.DocumentTransform_FieldTransform{FieldPath: "a"},
		&pb.DocumentTransform_FieldTransform{FieldPath: "b"},
	}
	bfp.Swap(0, 1)
	assert.Equal(t, "b", bfp[0].FieldPath)
	assert.Equal(t, "a", bfp[1].FieldPath)
}

func TestLess(t *testing.T) {
	bfp := byFieldPath{
		&pb.DocumentTransform_FieldTransform{FieldPath: "a"},
		&pb.DocumentTransform_FieldTransform{FieldPath: "b"},
	}
	assert.True(t, bfp.Less(0, 1))
	assert.False(t, bfp.Less(1, 0))
}
