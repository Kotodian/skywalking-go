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

package sarama

import (
	"errors"
	"strings"

	"github.com/IBM/sarama"
	"github.com/apache/skywalking-go/plugins/core/tracing"
)

type SaramaConsumeHook struct {
	Addr          string
	excludeTopics []string
}

func newSaramaConsumeHook(addr string) *SaramaConsumeHook {
	return &SaramaConsumeHook{
		Addr:          addr,
		excludeTopics: strings.Split(config.ExcludeTopics, ","),
	}
}

func (i *SaramaConsumeHook) OnConsume(message *sarama.ConsumerMessage) {
	operationName := kafkaPrefix + message.Topic + kafkaConsumerSuffix
	for _, topic := range i.excludeTopics {
		if topic == message.Topic {
			return
		}
	}
	s, err := tracing.CreateEntrySpan(operationName, func(headerKey string) (string, error) {
		for _, header := range message.Headers {
			if string(header.Key) == headerKey {
				return string(header.Value), nil
			}
		}
		return "", errors.New("header not found")
	},
		tracing.WithLayer(tracing.SpanLayerMQ),
		tracing.WithComponent(kafkaConsumerComponentID),
		tracing.WithTag(tracing.TagMQBroker, i.Addr),
		tracing.WithTag(tracing.TagMQTopic, message.Topic),
	)
	if err != nil {
		return
	}
	s.SetPeer(i.Addr)
	s.End()
}

type SaramaProduceHook struct {
	Addr          string
	excludeTopics []string
}

func newSaramaProduceHook(addr string) *SaramaProduceHook {
	return &SaramaProduceHook{
		Addr:          addr,
		excludeTopics: strings.Split(config.ExcludeTopics, ","),
	}
}

func (i *SaramaProduceHook) OnSend(message *sarama.ProducerMessage) {
	operationName := kafkaPrefix + message.Topic + kafkaProducerSuffix
	for _, topic := range i.excludeTopics {
		if topic == message.Topic {
			return
		}
	}
	s, err := tracing.CreateExitSpan(
		// operationName
		operationName,

		// peer
		i.Addr,

		// injector
		func(k, v string) error {
			if len(message.Headers) == 0 {
				message.Headers = []sarama.RecordHeader{
					{Key: []byte(k), Value: []byte(v)},
				}
			} else {
				for _, header := range message.Headers {
					if string(header.Key) == k {
						return errors.New("header already exists")
					}
				}
				message.Headers = append(message.Headers, sarama.RecordHeader{Key: []byte(k), Value: []byte(v)})
			}
			return nil
		},

		// opts
		tracing.WithLayer(tracing.SpanLayerMQ),
		tracing.WithComponent(kafkaProducerComponentID),
		tracing.WithTag(tracing.TagMQTopic, message.Topic),
		tracing.WithTag(tracing.TagMQBroker, i.Addr),
	)
	if err != nil {
		return
	}
	s.SetPeer(i.Addr)
	s.End()
}
