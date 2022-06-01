package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/alphaframework/alpha/aconfig"
	"github.com/alphaframework/alpha/alog"
)

func MustKafkaConfigWith(portName aconfig.PortName, appConfig *aconfig.Application, saramaConfig *aconfig.SaramaConfig) (*sarama.Config, []string, string) {
	c, brokers, topic, err := NewKafkaConfigWith(portName, appConfig, saramaConfig)
	if err != nil {
		panic(err)
	}

	return c, brokers, topic
}

func NewKafkaConfigWith(portName aconfig.PortName, appConfig *aconfig.Application, saramaConfig *aconfig.SaramaConfig) (*sarama.Config, []string, string, error) {
	location := appConfig.GetMatchedPrimaryPortLocation(portName)
	if location == nil {
		return nil, nil, "", fmt.Errorf("missing matched primaryport location(%s)", portName)
	}
	options := appConfig.GetSecondaryPort(portName).Options
	if options == nil {
		return nil, nil, "", fmt.Errorf("missing options for secondary port (%s)", portName)
	}

	return NewKafkaConfig(location, options, saramaConfig)
}

func NewKafkaConfig(location *aconfig.Location, options aconfig.KV, saramaConfig *aconfig.SaramaConfig) (*sarama.Config, []string, string, error) {
	alog.Sugar.Infof("kafka info, brokers: %s ,default_topic :%s", location.Address, options.GetString("topic"))

	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = saramaConfig.ProducerRetryMax

	conf.ClientID = saramaConfig.ClientID
	conf.ChannelBufferSize = saramaConfig.ChannelBufferSize
	conf.Producer.Return.Successes = saramaConfig.ProducerReturnSuccesses
	conf.Producer.Retry.Max = saramaConfig.ProducerRetryMax
	conf.Producer.RequiredAcks = saramaConfig.ProducerRequiredAcks
	conf.Producer.Timeout = time.Duration(saramaConfig.ProducerTimeout) * time.Second
	conf.Producer.Flush.Frequency = time.Duration(saramaConfig.ProducerFlushFrequency) * time.Millisecond

	conf.Net.TLS.Enable = saramaConfig.NetTLSEnable
	conf.Net.TLS.Config = saramaConfig.NetTLSConfig

	return conf, strings.Split(location.Address, ","), options.GetString("topic"), nil
}
