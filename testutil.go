// This file has been modified for the oringinal found at
// https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/internal/testutil/server.go
//
// Copyright 2016 Google LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mockfs

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"

	grpc "google.golang.org/grpc"
)

// A server is an in-process gRPC server, listening on a system-chosen port on
// the local loopback interface. Servers are for testing only and are not
// intended to be used in production code.
//
// To create a server, make a new server, register your handlers, then call
// Start:
//
//  srv, err := NewServer()
//  ...
//  mypb.RegisterMyServiceServer(srv.Gsrv, &myHandler)
//  ....
//  srv.Start()
//
// Clients should connect to the server with no security:
//
//  conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure())
//  ...
type server struct {
	Addr string
	Port int
	l    net.Listener
	Gsrv *grpc.Server
}

// NewServer creates a new server. The server will be listening for gRPC connections
// at the address named by the Addr field, without TLS.
func newServer(opts ...grpc.ServerOption) (*server, error) {
	return newServerWithPort(0, opts...)
}

// NewServerWithPort creates a new server at a specific port. The server will be listening
// for gRPC connections at the address named by the Addr field, without TLS.
func newServerWithPort(port int, opts ...grpc.ServerOption) (*server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}
	s := &server{
		Addr: l.Addr().String(),
		Port: parsePort(l.Addr().String()),
		l:    l,
		Gsrv: grpc.NewServer(opts...),
	}
	return s, nil
}

// Start causes the server to start accepting incoming connections.
// Call Start after registering handlers.
func (s *server) start() {
	go func() {
		if err := s.Gsrv.Serve(s.l); err != nil {
			log.Printf("mockfs.Server.Start: %v", err)
		}
	}()
}

// Close shuts down the server.
func (s *server) close() {
	s.Gsrv.Stop()
	s.l.Close()
}

var portParser = regexp.MustCompile(`:[0-9]+`)

func parsePort(addr string) int {
	res := portParser.FindAllString(addr, -1)
	if len(res) == 0 {
		panic(fmt.Errorf("parsePort: found no numbers in %s", addr))
	}
	stringPort := res[0][1:] // strip the :
	p, err := strconv.ParseInt(stringPort, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(p)
}
