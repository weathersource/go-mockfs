// This file has been modified for the oringinal found at
// https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/firestore/util_test.go
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

import (
	"context"

	firestore "cloud.google.com/go/firestore"
	errors "github.com/weathersource/go-errors"
	option "google.golang.org/api/option"
	grpc "google.golang.org/grpc"
)

// New creates a new Firestore Client and MockServer
func New() (*firestore.Client, *MockServer, error) {
	srv, err := newServer()
	if err != nil {
		return nil, nil, errors.NewUnknownError("Failed to create Firestore server.")
	}
	conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, errors.NewUnknownError("Failed to create Firestore connection.")
	}
	client, err := firestore.NewClient(context.Background(), "projectID", option.WithGRPCConn(conn))
	if err != nil {
		return nil, nil, errors.NewUnknownError("Failed to create Firestore client.")
	}
	return client, srv, nil
}
