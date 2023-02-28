package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"appDog/lib"
)

type StatusPayloadData struct {
	UniqueId string `json:"UniqueId"`
	Status   bool   `json:"Status"`
}

type ActPayloadData struct {
	UniqueId string `json:"UniqueId"`
	Act      int    `json:"Act"` // 1 代表打开 0 代表关闭
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// 检查进程是否已经存在，否则直接退出
	logger := lib.StartLogger()
	defer logger.StopLogger()
	logger.Info("程序监控程序开始")
	logger.Info("加载配置")
	config, err := lib.LoadConfig()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	processName := config.MainProcessName
	logger.Info("检查进程是否存在")

	if processName == "" {
		logger.Error("进程名称配置错误（为空）")
		os.Exit(1)
	}
	if lib.ProcessExists(processName) {
		logger.Error("进程已经存在，无法一次运行多个")
		os.Exit(1)
	}

	logger.Info("链接MQTT")

	pushStatusFunc := func(client mqtt.Client, UniqueId string, Status bool) {
		payloadData := StatusPayloadData{
			UniqueId: UniqueId,
			Status:   Status,
		}

		payloadJson, _ := json.Marshal(payloadData)

		token := client.Publish(config.StatusMqttTopic, 0, false, string(payloadJson))
		token.Wait()
		if token.Error() != nil {
			logger.Error(errors.New("Failed to subscribe to MQTT topic: " + token.Error().Error()))
			return
		}
		logger.Info(string(payloadJson) + "消息推送成功")
	}

	connectFunc := func(client mqtt.Client) {
		logger.Info("mqtt 已经连接成功")
		logger.Info("监听mqtt主题")
		var ActPayloadData ActPayloadData
		err := client.Subscribe(config.ActMqttTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
			logger.Info(fmt.Sprintf("Received message: %s from topic: %s\n", string(msg.Payload()), msg.Topic()))

			json.Unmarshal(msg.Payload(), &ActPayloadData)

			for _, item := range config.Apps {

				if item.UniqueId == ActPayloadData.UniqueId {
					// 检查状态
					/*appStatus := lib.ProcessExists(item.ProcessName)*/

					switch ActPayloadData.Act {
					case 0:
						err := lib.KillProcessByName(item.ProcessName)
						if err != nil {
							logger.Error("关闭软件出错：")
							logger.Error(err.Error())
							return
						}
						logger.Info(fmt.Sprintf("成功关闭 %s", item.Name))
						pushStatusFunc(client, item.ProcessName, false)

						break
					case 1:
						err := lib.OpenProcessByShortcut(item.ShortcutName)
						if err != nil {
							logger.Error("打开软件出错：")
							logger.Error(err.Error())
							return
						}
						logger.Info(fmt.Sprintf("成功打开 %s", item.Name))
						pushStatusFunc(client, item.ProcessName, true)
						break
					}
					break
				}
			}

		})
		if err != nil {
			logger.Error(err)
		}
	}

	connectLostFunc := func(client mqtt.Client, err error) {
		logger.Error(err)
		logger.Info("等待自动重连")
	}

	// 开始MQTT 链接
	mqttClient, err := lib.NewMQTTClient(config.Mqtt.Host, config.Mqtt.ClientId, config.Mqtt.Username, config.Mqtt.Password, config.Mqtt.CleanSession, connectFunc, connectLostFunc)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// 订阅MQTT 主题
	logger.Info("起一个协程上报程序状态")
	go func() {
		for {
			appList := config.Apps
			for _, item := range appList {
				// 这个地方可能会存在瓶颈，因为这里会多次穷举进程信息
				status := lib.ProcessExists(item.ProcessName)
				logger.Info(fmt.Sprintf("检查程序%s的状态,当前状态为%t", item.Name, status))
				pushStatusFunc(mqttClient.Client, item.UniqueId, status)
			}
			waitSeconds := config.ProcessStatusCheckRate // 等待秒数
			waitTime := time.Duration(waitSeconds) * time.Second
			time.Sleep(waitTime)
		}
	}()
	<-c
	mqttClient.Client.Disconnect(250)

}
