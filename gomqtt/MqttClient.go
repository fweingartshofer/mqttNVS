package gomqtt

import (
	"fmt"
	"os"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

var (
	mClient mqtt.Client
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("[ ", time.Now(), "]", " Message received: ")
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func create(){
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	hostame, _ := os.Hostname()
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID(hostame)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	mClient = mqtt.NewClient(opts)

	if token := mClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}


func listening(){
	for mClient.IsConnected(){
		time.Sleep(1* time.Second)
	}
}

func sendMsg(topic string, payload string, qos byte, retained bool) mqtt.Token{
	token := mClient.Publish(topic, qos, retained, payload)
	token.Wait()
	return token
}

func subscribe(topic string) mqtt.Token{
	token := mClient.Subscribe("client/laptop", 0, nil)
	token.Wait()
	return token

}

func unsubscribe() mqtt.Token{
	token := mClient.Unsubscribe("client/laptop")
	token.Wait()
	return token
}

func exit(){
	mClient.Disconnect(250)
	os.Exit(1)
}
