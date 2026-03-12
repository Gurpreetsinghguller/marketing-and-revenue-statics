package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/logger"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/event"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// MQTTBroker implements EventBroker interface for MQTT
type MQTTBroker struct {
	client    mqtt.Client
	brokerURL string
	topic     string
	clientID  string
	log       *logrus.Logger
	QoS       int
}

// liskov substitution principle: MQTTBroker can be used wherever EventBroker is expected
var _ event.EventBroker = (*MQTTBroker)(nil)

func NewBroker(brokerURL, clientID, topic string, qos int) event.EventBroker {
	return &MQTTBroker{
		brokerURL: brokerURL,
		clientID:  clientID,
		topic:     topic,
		log:       logger.Get(),
		QoS:       qos,
	}
}

func (m *MQTTBroker) Start(ctx context.Context, handler event.EventHandler) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.brokerURL)
	opts.SetClientID(m.clientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Second * 10)
	opts.SetConnectRetryInterval(time.Second * 3)
	opts.SetMessageChannelDepth(1000)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		m.processMessage(client, msg, handler)
	})

	// Re-subscribe after every reconnect
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		m.log.Info("MQTT connected, subscribing to topic...")
		if token := client.Subscribe(m.topic, uint8(m.QoS), nil); token.Wait() && token.Error() != nil {
			m.log.WithError(token.Error()).Error("failed to subscribe to topic")
		} else {
			m.log.WithField("topic", m.topic).Info("successfully subscribed to topic")
		}
	})

	// Log disconnections
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		m.log.WithError(err).Warn("MQTT connection lost, will auto-reconnect...")
	})

	m.client = mqtt.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	m.log.WithFields(logrus.Fields{
		"broker_url": m.brokerURL,
		"topic":      m.topic,
		"client_id":  m.clientID,
	}).Info("MQTT broker connected and listening")

	// Wait briefly to ensure OnConnectHandler subscription completes
	time.Sleep(500 * time.Millisecond)

	// Verify subscription succeeded
	if !m.IsConnected() {
		return fmt.Errorf("MQTT client connection lost after startup")
	}

	return nil
}

func (m *MQTTBroker) IsConnected() bool {
	return m.client != nil && m.client.IsConnected()
}

// Close gracefully closes the MQTT broker connection
func (m *MQTTBroker) Close() error {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250)
		m.log.WithField("client_id", m.clientID).Info("MQTT broker disconnected")
	}
	return nil
}

// Health checks MQTT broker connectivity
func (m *MQTTBroker) Health(ctx context.Context) error {
	if m.client == nil || !m.client.IsConnected() {
		return fmt.Errorf("MQTT broker is not connected")
	}
	return nil
}

func (m *MQTTBroker) processMessage(client mqtt.Client, msg mqtt.Message, handler event.EventHandler) {
	var event event.Event
	m.log.WithFields(logrus.Fields{
		"topic":      m.topic,
		"client_id":  m.clientID,
		"message_id": msg.MessageID(),
		"payload":    string(msg.Payload()),
	}).Info("Received MQTT message")
	if err := json.Unmarshal(msg.Payload(), &event); err != nil {
		m.log.WithError(err).Error("failed to unmarshal MQTT message")
		return
	}

	// Set broker metadata
	event.Source = "mqtt"
	event.MessageID = fmt.Sprintf("%d", msg.MessageID())
	event.ProcessedAt = time.Now().UTC().Format(time.RFC3339)

	// Process event
	if err := handler(context.Background(), &event); err != nil {
		m.log.WithError(err).Error("failed to handle MQTT event")
	}
}
