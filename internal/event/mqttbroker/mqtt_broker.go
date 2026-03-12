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
}

var _ event.EventBroker = (*MQTTBroker)(nil)

func NewBroker(brokerURL, clientID, topic string) event.EventBroker {
	return &MQTTBroker{
		brokerURL: brokerURL,
		clientID:  clientID,
		topic:     topic,
		log:       logger.Get(),
	}
}

func (m *MQTTBroker) Start(ctx context.Context, handler event.EventHandler) error {
	fmt.Println("MQTT CREDS", m.brokerURL, m.clientID)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.brokerURL)
	opts.SetClientID(m.clientID)
	opts.SetCleanSession(true) // ← ADD THIS: clears ghost sessions
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Second * 10)
	opts.SetConnectRetryInterval(time.Second * 3) // ← ADD THIS
	opts.SetMessageChannelDepth(1000)
	// opts.SetUsername("username")
	// opts.SetPassword("password")

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		m.processMessage(client, msg, handler)
	})

	// ← ADD THIS: re-subscribe after every reconnect
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		m.log.Info("MQTT connected, subscribing to topic...")
		if token := client.Subscribe(m.topic, 1, nil); token.Wait() && token.Error() != nil {
			m.log.WithError(token.Error()).Error("failed to re-subscribe after reconnect")
		}
	})

	// ← ADD THIS: log disconnections
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

	// NOTE: Remove Subscribe from here — OnConnectHandler handles it now
	// This avoids double-subscribing on first connect

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

// Publish sends an event to the MQTT broker
func (m *MQTTBroker) Publish(ctx context.Context, event *event.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	token := m.client.Publish(m.topic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish event: %w", token.Error())
	}

	m.log.WithFields(logrus.Fields{
		"topic":     m.topic,
		"client_id": m.clientID,
		"event_id":  event.ID,
	}).Info("Event published to MQTT broker")

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
