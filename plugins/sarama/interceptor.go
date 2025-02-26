package sarama

import (
	"strings"

	"github.com/IBM/sarama"

	"github.com/apache/skywalking-go/plugins/core/operator"
)

// https://github.com/apache/skywalking/blob/master/oap-server/server-starter/src/main/resources/component-libraries.yml
const (
	kafkaProducerComponentID = 40
	kafkaConsumerComponentID = 41
	kafkaPrefix              = "Kafka/"
	kafkaProducerSuffix      = "/Producer"
	kafkaConsumerSuffix      = "/Consumer"
)

type SaramaInterceptor struct {
}

// BeforeInvoke would be called before the target method invocation.
func (s *SaramaInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	addr := strings.Join(invocation.Args()[0].([]string), ",")
	config := invocation.Args()[1].(*sarama.Config)

	if len(config.Consumer.Interceptors) > 0 {
		config.Consumer.Interceptors = append(config.Consumer.Interceptors, newSaramaConsumeHook(addr))
	} else {
		config.Consumer.Interceptors = []sarama.ConsumerInterceptor{newSaramaConsumeHook(addr)}
	}
	if len(config.Producer.Interceptors) > 0 {
		config.Producer.Interceptors = append(config.Producer.Interceptors, newSaramaProduceHook(addr))
	} else {
		config.Producer.Interceptors = []sarama.ProducerInterceptor{newSaramaProduceHook(addr)}
	}
	return nil
}

// AfterInvoke would be called after the target method invocation.
func (s *SaramaInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {
	return nil
}
