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
	Addr string
	data map[string][]datum
}

type datum struct {
	req    proto.Message
	adjust AdjustFunc
	resp   interface{}
}

// AdjustFunc is the signature for functions that adjust requests before querying matching data.
// Ensure inside AdjustFunc functions, you modify and return a copy of the original requests,
// not pointers to the modified originals. This allows for running AdjustFunc multiple times
// against multiple possible results
type AdjustFunc func(wantReq proto.Message, gotReq proto.Message) (wantReqAdj proto.Message)

func newServer() (*MockServer, error) {
	srv, err := gsrv.NewServer()
	if err != nil {
		return nil, err
	}
	mock := &MockServer{Addr: srv.Addr}
	mock.data = make(map[string][]datum)
	pb.RegisterFirestoreServer(srv.Gsrv, mock)
	srv.Start()
	return mock, nil
}

// Reset returns the MockServer to an empty state.
func (s *MockServer) Reset() {
	s.data = make(map[string][]datum)
}

// AddData adds a (request, response) pair for a given rpc function name to the
// server's list of expected interactions. The server will compare the incoming
// request with wantReq using proto.Equal. The response can be a message or an
// error.
//
// For the Listen RPC, resp should be a interface{}, where the element
// is either ListenResponse or an error.
//
// Passing nil for wantReq disables the request check.
func (s *MockServer) AddData(rpc string, wantReq proto.Message, resp interface{}) {
	s.AddDataAdjust(rpc, wantReq, resp, nil)
}

// AddDataAdjust is like AddData, but accepts a function that can be used
// to tweak the requests before comparison, for example to adjust for
// randomness.
func (s *MockServer) AddDataAdjust(rpc string, wantReq proto.Message, resp interface{}, adjust AdjustFunc) {
	_, ok := s.data[rpc]
	if !ok {
		s.data[rpc] = []datum{}
	}
	s.data[rpc] = append(s.data[rpc], datum{wantReq, adjust, resp})
}

// getData compares the request with the next expected (request, response) pair.
// It returns the response, or an error if the request doesn't match what
// was expected or there are no expected data.
func (s *MockServer) getData(rpc string, wantReq proto.Message) (interface{}, error) {
	_, ok := s.data[rpc]
	if !ok {
		return nil, errors.NewNotFoundError("The RPC could not be found.")
	}

	for _, got := range s.data[rpc] {

		if got.req != nil {

			// Handle request adjustments with proper scope
			wantReqAdj := wantReq
			if got.adjust != nil {
				wantReqAdj = got.adjust(wantReq, got.req)
			}

			// Sort FieldTransforms by FieldPath, since slice order is undefined and proto.Equal
			// is strict about order.
			switch wantReqTyped := wantReqAdj.(type) {
			case *pb.CommitRequest:
				for _, w := range wantReqTyped.Writes {
					switch opTyped := w.Operation.(type) {
					case *pb.Write_Transform:
						sort.Sort(byFieldPath(opTyped.Transform.FieldTransforms))
					}
				}
			}

			if !proto.Equal(wantReqAdj, got.req) {
				continue
			}
		}
		if err, ok := got.resp.(error); ok {
			return nil, err
		}
		return got.resp, nil
	}

	return nil, errors.NewNotFoundError("The requested response object was not found.")
}

type byFieldPath []*pb.DocumentTransform_FieldTransform

func (a byFieldPath) Len() int           { return len(a) }
func (a byFieldPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFieldPath) Less(i, j int) bool { return a[i].FieldPath < a[j].FieldPath }
