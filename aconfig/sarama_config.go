package aconfig

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/Shopify/sarama"
	"github.com/dartagnanli/alpha/alog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// type RequiredAcks int16

const (
	// config for sarama
	defaultKafkaClientID          = "sarama"
	defaultKafkaChannelBufferSize = 256

	// NoResponse doesn't send any response, the TCP ACK is all you get.
	NoResponse sarama.RequiredAcks = 0
	// WaitForLocal waits for only the local commit to succeed before responding.
	WaitForLocal sarama.RequiredAcks = 1
	// WaitForAll waits for all in-sync replicas to commit before responding.
	// The minimum number of in-sync replicas is configured on the broker via
	// the `min.insync.replicas` configuration key.
	WaitForAll sarama.RequiredAcks = -1
)

type SaramaConfig struct {
	ClientID                string              `json:"client_id,omitempty"`
	ChannelBufferSize       int                 `json:"channel_buffer_size,omitempty"`
	ProducerReturnSuccesses bool                `json:"producer_return_successes,omitempty"`
	ProducerRetryMax        int                 `json:"producer_retry_max,omitempty"`
	ProducerRequiredAcks    sarama.RequiredAcks `json:"producer_required_acks,omitempty"`
	ProducerTimeout         int                 `json:"producer_timeout,omitempty"`
	ProducerFlushFrequency  int                 `json:"producer_flush_frequency,omitempty"`

	NetTLSEnable       bool        `json:"net_tls_enable,omitempty"`
	NetTLSConfig       *tls.Config `json:"net_tls_config,omitempty"`
	InsecureSkipVerify bool        `json:"insecure_skip_verify,omitempty"`
	CertFile           string      `json:"cert_file,omitempty"`
	KeyFile            string      `json:"key_file,omitempty"`
	CaFile             string      `json:"ca_file,omitempty"`
	K8sNamespace       string      `json:"k8s_namespace,omitempty"`
	K8sSecret          string      `json:"k8s_secret,omitempty"`
}

func (sc *SaramaConfig) complete() {

	if sc.ClientID == "" {
		sc.ClientID = defaultKafkaClientID
	}
	if sc.ChannelBufferSize == 0 {
		sc.ChannelBufferSize = defaultKafkaChannelBufferSize
	}
	if !sc.ProducerReturnSuccesses {
		sc.ProducerReturnSuccesses = false
	}
	if sc.ProducerRetryMax == 0 {
		sc.ProducerRetryMax = 3
	}
	if sc.ProducerRequiredAcks == 0 {
		sc.ProducerRequiredAcks = WaitForLocal
	}
	if sc.InsecureSkipVerify {
		sc.InsecureSkipVerify = true
	}
	if sc.NetTLSEnable {
		tlsConfig := sc.createTLSConfigurationWithK8s()
		if tlsConfig != nil {
			sc.NetTLSConfig = tlsConfig
		}
	}
	if sc.ProducerTimeout == 0 {
		sc.ProducerTimeout = 10
	}
	if sc.ProducerFlushFrequency == 0 {
		sc.ProducerFlushFrequency = 500
	}
}

func (sc *SaramaConfig) createTLSConfiguration() (t *tls.Config) {
	if sc.CertFile != "" && sc.KeyFile != "" && sc.CaFile != "" {
		cert, err := tls.LoadX509KeyPair(sc.CertFile, sc.KeyFile)
		if err != nil {
			panic(err)
		}

		caCert, err := os.ReadFile(sc.CaFile)
		if err != nil {
			panic(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: sc.InsecureSkipVerify,
		}
	}
	// will be nil by default if nothing is provided
	return t
}

func (sc *SaramaConfig) createTLSConfigurationWithK8s() (t *tls.Config) {

	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	secret := sc.K8sSecret
	namespace := sc.K8sNamespace
	mySecret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			alog.Sugar.Errorf("Secret %s in namespace %s not found\n", secret, namespace)
		}
	}
	cert, err := tls.X509KeyPair(mySecret.Data["cert.pem"], mySecret.Data["key.pem"])
	if err != nil {
		panic(err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(mySecret.Data["ca_cert.pem"])

	t = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: sc.InsecureSkipVerify,
	}

	return t
}
