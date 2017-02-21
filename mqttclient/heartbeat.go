package mqttclient

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//Heartbeat for MQTT Client
//Will Subscribe and Publish to MQTT
type Heartbeat struct {
	Topic             string
	Client            *MQTTClient
	OnHeartbeatMissed func(error)
	LastHeartbeat     time.Time
	stopped           bool
	MaxDurationMs     int64
}

func (h *Heartbeat) receivedHeartbeat(client MQTT.Client, msg MQTT.Message) {
	h.LastHeartbeat = time.Now()
}

//Start the heartbeat
func (h *Heartbeat) Start() error {
	if h.MaxDurationMs == 0 {
		h.MaxDurationMs = 5000
	}
	h.stopped = false
	client := h.Client.MQTTClient
	if !client.IsConnected() {
		return fmt.Errorf("Can't start heartbeat. MQTT Client is not connected!")
	}
	if token := client.Subscribe(h.Topic, uint8(0), h.receivedHeartbeat); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	h.LastHeartbeat = time.Now()

	go func() {
		for !h.stopped {
			token := h.Client.MQTTClient.Publish(h.Topic, uint8(2), false, "OK")
			token.Wait()
			err := token.Error()
			if err != nil {
				h.OnHeartbeatMissed(err)
				h.stopped = true
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		for !h.stopped {
			duration := time.Now().Sub(h.LastHeartbeat).Nanoseconds() / 1000000
			if duration > h.MaxDurationMs {
				h.stopped = true
				h.OnHeartbeatMissed(fmt.Errorf("Timeout in heartbeat: %d.", duration))
				return
			}
			time.Sleep(500 * time.Millisecond)
		}

	}()
	return nil
}
