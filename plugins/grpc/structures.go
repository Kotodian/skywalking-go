// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package grpc

import (
	"context"
)

//skywalking:native google.golang.org/grpc/internal/transport Stream
type nativeStream struct {
	ctx    context.Context
	method string
}

func (s *nativeStream) Method() string {
	return s.method
}

func (s *nativeStream) Context() context.Context {
	return s.ctx
}

//skywalking:native google.golang.org/grpc/internal/transport ServerStream
type nativeServerStream struct {
	s *nativeStream
}

func (s *nativeServerStream) Method() string {
	return s.s.Method()
}

func (s *nativeServerStream) Context() context.Context {
	return s.s.Context()
}

//skywalking:native google.golang.org/grpc ClientConn
type nativeClientConn struct {
}

func (cc *nativeClientConn) Target() string {
	return ""
}

//skywalking:native google.golang.org/grpc clientStream
type nativeclientStream struct {
	callHdr *nativeCallHdr
}

//skywalking:native google.golang.org/grpc/internal/transport Stream
type nativeCallHdr struct {
	Method string
}

//skywalking:native google.golang.org/grpc serverStream
type nativeserverStream struct {
	s *nativeStream
}
