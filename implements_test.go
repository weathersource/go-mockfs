package mockfs

import (
	"context"
	"fmt"
	"testing"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	errors "github.com/weathersource/go-errors"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/grpc"
)

type BatchGetDocumentsServer struct {
	grpc.ServerStream
	resp *pb.BatchGetDocumentsResponse
}

func (s *BatchGetDocumentsServer) Send(resp *pb.BatchGetDocumentsResponse) error {
	s.resp = resp
	return nil
}

type BatchGetDocumentsServerError struct {
	grpc.ServerStream
	resp *pb.BatchGetDocumentsResponse
}

func (s *BatchGetDocumentsServerError) Send(resp *pb.BatchGetDocumentsResponse) error {
	return errors.NewInternalError("")
}

type RunQueryServer struct {
	grpc.ServerStream
	resp *pb.RunQueryResponse
}

func (s *RunQueryServer) Send(resp *pb.RunQueryResponse) error {
	s.resp = resp
	return nil
}

type RunQueryServerError struct {
	grpc.ServerStream
	resp *pb.RunQueryResponse
}

func (s *RunQueryServerError) Send(resp *pb.RunQueryResponse) error {
	return errors.NewInternalError("")
}

type ListenServer struct {
	grpc.ServerStream
	req  *pb.ListenRequest
	resp *pb.ListenResponse
}

func (s *ListenServer) Send(resp *pb.ListenResponse) error {
	s.resp = resp
	return nil
}

func (s *ListenServer) Recv() (*pb.ListenRequest, error) {
	return s.req, nil
}

type ListenServerRError struct {
	grpc.ServerStream
	req  *pb.ListenRequest
	resp *pb.ListenResponse
}

func (s *ListenServerRError) Send(resp *pb.ListenResponse) error {
	s.resp = resp
	return nil
}

func (s *ListenServerRError) Recv() (*pb.ListenRequest, error) {
	return nil, errors.NewInternalError("")
}

type ListenServerSError struct {
	grpc.ServerStream
	req  *pb.ListenRequest
	resp *pb.ListenResponse
}

func (s *ListenServerSError) Send(resp *pb.ListenResponse) error {
	return errors.NewInternalError("")
}

func (s *ListenServerSError) Recv() (*pb.ListenRequest, error) {
	return s.req, nil
}

func TestGetDocument(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	_, srv, err := New()
	assert.Nil(err)

	// test valid response
	srv.AddRPC(
		nil,
		&pb.Document{},
	)
	resp, err := srv.GetDocument(ctx, &pb.GetDocumentRequest{})
	assert.Nil(err)
	assert.NotNil(resp)

	// test error response
	srv.AddRPC(
		nil,
		errors.NewInternalError(""),
	)
	_, err = srv.GetDocument(ctx, &pb.GetDocumentRequest{})
	assert.NotNil(err)
}

func TestCommit(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	_, srv, err := New()
	assert.Nil(err)

	// test valid response
	srv.AddRPC(
		nil,
		&pb.CommitResponse{},
	)
	resp, err := srv.Commit(ctx, &pb.CommitRequest{})
	assert.Nil(err)
	assert.NotNil(resp)

	// test error response
	srv.AddRPC(
		nil,
		errors.NewInternalError(""),
	)
	_, err = srv.Commit(ctx, &pb.CommitRequest{})
	assert.NotNil(err)
}

func TestBatchGetDocuments(t *testing.T) {
	assert := assert.New(t)
	_, srv, err := New()
	assert.Nil(err)

	bs := BatchGetDocumentsServer{}
	bse := BatchGetDocumentsServerError{}

	// test valid response
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.BatchGetDocumentsResponse{},
		},
	)
	err = srv.BatchGetDocuments(&pb.BatchGetDocumentsRequest{}, &bs)
	assert.Nil(err)
	assert.NotNil(bs.resp)

	// test error send
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.BatchGetDocumentsResponse{},
		},
	)
	err = srv.BatchGetDocuments(&pb.BatchGetDocumentsRequest{}, &bse)
	assert.NotNil(err)

	// test error response
	srv.AddRPC(
		nil,
		errors.NewInternalError(""),
	)
	err = srv.BatchGetDocuments(&pb.BatchGetDocumentsRequest{}, &bs)
	assert.NotNil(err)

	// test error response in batch
	srv.AddRPC(
		nil,
		[]interface{}{
			errors.NewInternalError(""),
		},
	)
	err = srv.BatchGetDocuments(&pb.BatchGetDocumentsRequest{}, &bs)
	assert.NotNil(err)

	// test wrong type in batch
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.GetDocumentRequest{},
		},
	)
	assert.Panics(func() {
		srv.BatchGetDocuments(&pb.BatchGetDocumentsRequest{}, &bs)
	})
}

func TestRunQuery(t *testing.T) {
	assert := assert.New(t)
	_, srv, err := New()
	assert.Nil(err)

	qs := RunQueryServer{}
	qse := RunQueryServerError{}

	// test valid response
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.RunQueryResponse{},
		},
	)
	err = srv.RunQuery(&pb.RunQueryRequest{}, &qs)
	assert.Nil(err)
	assert.NotNil(qs.resp)

	// test error send
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.RunQueryResponse{},
		},
	)
	err = srv.RunQuery(&pb.RunQueryRequest{}, &qse)
	assert.NotNil(err)

	// test error response
	srv.AddRPC(
		nil,
		errors.NewInternalError(""),
	)
	err = srv.RunQuery(&pb.RunQueryRequest{}, &qs)
	assert.NotNil(err)

	// test error response in batch
	srv.AddRPC(
		nil,
		[]interface{}{
			errors.NewInternalError(""),
		},
	)
	err = srv.RunQuery(&pb.RunQueryRequest{}, &qs)
	assert.NotNil(err)

	// test wrong type in batch
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.GetDocumentRequest{},
		},
	)
	assert.Panics(func() {
		srv.RunQuery(&pb.RunQueryRequest{}, &qs)
	})
}

func TestBeginTransaction(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	_, srv, err := New()
	assert.Nil(err)

	// test valid response
	srv.AddRPC(
		nil,
		&pb.BeginTransactionResponse{},
	)
	resp, err := srv.BeginTransaction(ctx, &pb.BeginTransactionRequest{})
	assert.Nil(err)
	assert.NotNil(resp)

	// test error response
	srv.AddRPC(
		nil,
		errors.NewInternalError(""),
	)
	_, err = srv.BeginTransaction(ctx, &pb.BeginTransactionRequest{})
	assert.NotNil(err)
}

func TestRollback(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	_, srv, err := New()
	assert.Nil(err)

	// test valid response
	srv.AddRPC(
		nil,
		&empty.Empty{},
	)
	resp, err := srv.Rollback(ctx, &pb.RollbackRequest{})
	assert.Nil(err)
	assert.NotNil(resp)

	// test error response
	srv.AddRPC(
		nil,
		errors.NewInternalError(""),
	)
	_, err = srv.Rollback(ctx, &pb.RollbackRequest{})
	assert.NotNil(err)
}

func TestListen(t *testing.T) {
	assert := assert.New(t)
	_, srv, err := New()
	assert.Nil(err)

	ls := ListenServer{req: &pb.ListenRequest{}}
	lsre := ListenServerRError{req: &pb.ListenRequest{}}
	lsse := ListenServerSError{req: &pb.ListenRequest{}}

	// test valid response
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.ListenResponse{},
		},
	)
	err = srv.Listen(&ls)
	assert.Nil(err)
	assert.NotNil(ls.resp)

	// test error recv
	err = srv.Listen(&lsre)
	assert.NotNil(err)

	// test error send
	srv.AddRPC(
		nil,
		[]interface{}{
			&pb.ListenResponse{},
		},
	)
	err = srv.Listen(&lsse)
	assert.NotNil(err)

	// test error response
	srv.AddRPC(
		nil,
		errors.NewInternalError(""),
	)
	err = srv.Listen(&ls)
	assert.NotNil(err)

	// test error response in stream
	srv.AddRPC(
		nil,
		[]interface{}{
			errors.NewInternalError(""),
		},
	)
	err = srv.Listen(&ls)
	assert.NotNil(err)

	// test panic on unknown error
	fmt.Println("FOOBAR")
	srv.AddRPC(
		nil,
		errors.NewUnknownError(""),
	)
	assert.Panics(func() {
		srv.Listen(&ls)
	})
}
