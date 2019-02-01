// This file has been modified for the oringinal found at
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
	"fmt"
	"sort"

	proto "github.com/golang/protobuf/proto"
	errors "github.com/weathersource/go-errors"
	gsrv "github.com/weathersource/go-gsrv"
	pb "google.golang.org/genproto/googleapis/firestore/v1beta1"
)

// MockServer mocks the pb.FirestoreServer interface
// (https://godoc.org/google.golang.org/genproto/googleapis/firestore/v1beta1#FirestoreServer)
type MockServer struct {
	pb.FirestoreServer
	Addr     string
	reqItems []reqItem
	resps    []interface{}
}

type reqItem struct {
	wantReq proto.Message
	adjust  func(gotReq proto.Message)
}

func newServer() (*MockServer, error) {
	srv, err := gsrv.NewServer()
	if err != nil {
		return nil, err
	}
	mock := &MockServer{Addr: srv.Addr}
	pb.RegisterFirestoreServer(srv.Gsrv, mock)
	srv.Start()
	return mock, nil
}

// Reset returns the MockServer to an empty state.
func (s *MockServer) Reset() {
	s.reqItems = nil
	s.resps = nil
}

// AddRPC adds a (request, response) pair to the server's list of expected
// interactions. The server will compare the incoming request with wantReq
// using proto.Equal. The response can be a message or an error.
//
// For the Listen RPC, resp should be a []interface{}, where each element
// is either ListenResponse or an error.
//
// Passing nil for wantReq disables the request check.
func (s *MockServer) AddRPC(wantReq proto.Message, resp interface{}) {
	s.AddRPCAdjust(wantReq, resp, nil)
}

// AddRPCAdjust is like AddRPC, but accepts a function that can be used
// to tweak the requests before comparison, for example to adjust for
// randomness.
func (s *MockServer) AddRPCAdjust(wantReq proto.Message, resp interface{}, adjust func(gotReq proto.Message)) {
	s.reqItems = append(s.reqItems, reqItem{wantReq, adjust})
	s.resps = append(s.resps, resp)
}

// popRPC compares the request with the next expected (request, response) pair.
// It returns the response, or an error if the request doesn't match what
// was expected or there are no expected rpcs.
func (s *MockServer) popRPC(gotReq proto.Message) (interface{}, error) {
	if len(s.reqItems) == 0 || len(s.resps) == 0 {
		panic("mockfs.popRPC: Out of RPCs.")
	}
	ri := s.reqItems[0]
	resp := s.resps[0]
	s.reqItems = s.reqItems[1:]
	s.resps = s.resps[1:]
	if ri.wantReq != nil {
		if ri.adjust != nil {
			ri.adjust(gotReq)
		}

		// Sort FieldTransforms by FieldPath, since slice order is undefined and proto.Equal
		// is strict about order.
		switch gotReqTyped := gotReq.(type) {
		case *pb.CommitRequest:
			for _, w := range gotReqTyped.Writes {
				switch opTyped := w.Operation.(type) {
				case *pb.Write_Transform:
					sort.Sort(byFieldPath(opTyped.Transform.FieldTransforms))
				}
			}
		}

		if !proto.Equal(gotReq, ri.wantReq) {
			return nil, errors.NewInternalError(fmt.Sprintf("mockfs.popRPC: Bad request\ngot:  %T\n%s\nwant: %T\n%s",
				gotReq, proto.MarshalTextString(gotReq),
				ri.wantReq, proto.MarshalTextString(ri.wantReq)))
		}
	}
	if err, ok := resp.(error); ok {
		return nil, err
	}
	return resp, nil
}

type byFieldPath []*pb.DocumentTransform_FieldTransform

func (a byFieldPath) Len() int           { return len(a) }
func (a byFieldPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFieldPath) Less(i, j int) bool { return a[i].FieldPath < a[j].FieldPath }
