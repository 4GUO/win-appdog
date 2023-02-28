package lib

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient 是MQTT客户端结构体
type MQTTClient struct {
	Client mqtt.Client
}

// NewMQTTClient 创建一个新的MQTT客户端
func NewMQTTClient(brokerURL string, clientID string, username string, password string, cleanSession bool, connectFunc mqtt.OnConnectHandler, connectLostFunc mqtt.ConnectionLostHandler) (*MQTTClient, error) {
	// 创建MQTT客户端连接选项
	opts := mqtt.NewClientOptions()
	// 设置MQTT代理地址
	opts.AddBroker(brokerURL)
	// 设置MQTT客户端ID
	opts.SetClientID(clientID)

	// 设置MQTT客户端身份验证信息
	if username != "" && password != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}
	opts.SetKeepAlive(60 * time.Second)
	// 设置消息回调处理函数
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10000)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetCleanSession(cleanSession)
	opts.SetConnectTimeout(30 * time.Second)
	opts.OnConnect = connectFunc
	opts.OnConnectionLost = connectLostFunc
	// 创建MQTT客户端
	client := mqtt.NewClient(opts)

	// 连接MQTT代理
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, errors.New("Failed to connect to MQTT broker: " + token.Error().Error())
	}

	return &MQTTClient{Client: client}, nil
}
