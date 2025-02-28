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

package clickhousegov2

import (
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
)

// https://github.com/apache/skywalking/blob/master/oap-server/server-starter/src/main/resources/component-libraries.yml
const (
	ComponentClickhouse = 120
	DbTypeClickhouse    = "clickhouse"
)

type CkConnectSendQueryInterceptor struct {
}

func (i *CkConnectSendQueryInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	query := invocation.Args()[0].(string)
	connect := invocation.CallerInstance().(*nativeConnect)
	span, err := tracing.CreateExitSpan(query,
		connect.conn.RemoteAddr().String(),
		func(k, v string) error {
			return nil
		},
		tracing.WithLayer(tracing.SpanLayerDatabase),
		tracing.WithComponent(ComponentClickhouse),
		tracing.WithTag(tracing.TagDBType, DbTypeClickhouse),
		tracing.WithTag(tracing.TagDBInstance, connect.conn.RemoteAddr().String()),
		tracing.WithTag(tracing.TagDBStatement, query),
		tracing.WithTag("version", newVersion().String()),
	)
	if err != nil {
		return err
	}
	invocation.SetContext(span)
	return nil
}

func (i *CkConnectSendQueryInterceptor) AfterInvoke(invocation operator.Invocation, results ...interface{}) error {
	ctx := invocation.GetContext()
	if ctx == nil {
		return nil
	}
	if err, ok := results[0].(error); ok && err != nil {
		ctx.(tracing.Span).Error(err.Error())
	}
	ctx.(tracing.Span).End()
	return nil
}
