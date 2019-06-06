// This file has been modified for the original found at
// https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/firestore/mock_test.go
//
// Copyright 2017 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mockfs

// A simple mock server.

import (
	"context"
	"fmt"

	empty "github.com/golang/protobuf/ptypes/empty"
	pb "google.golang.org/genproto/googleapis/firestore/v1beta2"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// GetDocument overrides the FirestoreServer GetDocument method
func (s *MockServer) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.Document, error) {
	res, err := s.popRPC(req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.Document), nil
}

// Commit overrides the FirestoreServer Commit method
func (s *MockServer) Commit(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	res, err := s.popRPC(req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.CommitResponse), nil
}

// BatchGetDocuments overrides the FirestoreServer BatchGetDocuments method
func (s *MockServer) BatchGetDocuments(req *pb.BatchGetDocumentsRequest, bs pb.Firestore_BatchGetDocumentsServer) error {
	res, err := s.popRPC(req)
	if err != nil {
		return err
	}
	responses := res.([]interface{})
	for _, res := range responses {
		switch res := res.(type) {
		case *pb.BatchGetDocumentsResponse:
			if err := bs.Send(res); err != nil {
				return err
			}
		case error:
			return res
		default:
			panic(fmt.Sprintf("mockfs.BatchGetDocuments: Bad response type: %+v", res))
		}
	}
	return nil
}

// RunQuery overrides the FirestoreServer RunQuery method
func (s *MockServer) RunQuery(req *pb.RunQueryRequest, qs pb.Firestore_RunQueryServer) error {
	res, err := s.popRPC(req)
	if err != nil {
		return err
	}
	responses := res.([]interface{})
	for _, res := range responses {
		switch res := res.(type) {
		case *pb.RunQueryResponse:
			if err := qs.Send(res); err != nil {
				return err
			}
		case error:
			return res
		default:
			panic(fmt.Sprintf("mockfs.RunQuery: Bad response type: %+v", res))
		}
	}
	return nil
}

// BeginTransaction overrides the FirestoreServer BeginTransaction method
func (s *MockServer) BeginTransaction(ctx context.Context, req *pb.BeginTransactionRequest) (*pb.BeginTransactionResponse, error) {
	res, err := s.popRPC(req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.BeginTransactionResponse), nil
}

// Rollback overrides the FirestoreServer Rollback method
func (s *MockServer) Rollback(ctx context.Context, req *pb.RollbackRequest) (*empty.Empty, error) {
	res, err := s.popRPC(req)
	if err != nil {
		return nil, err
	}
	return res.(*empty.Empty), nil
}

// Listen overrides the FirestoreServer Listen method
func (s *MockServer) Listen(stream pb.Firestore_ListenServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	responses, err := s.popRPC(req)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			panic(err)
		}
		return err
	}
	for _, res := range responses.([]interface{}) {
		if err, ok := res.(error); ok {
			return err
		}
		if err := stream.Send(res.(*pb.ListenResponse)); err != nil {
			return err
		}
	}
	return nil
}
