package main

import (
	"fmt"
	"os"
	"time"
	"github.com/eclipse/paho.mqtt.golang"
	"flag"
	"log"
)


var (
	mClient mqtt.Client
)

func main() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("[ ", time.Now(), "]", " Message received: ")
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	listen := flag.String("l","#","Subscribe to a  topic and listen to it")
	flag.Parse()

	hostame, _ := os.Hostname()
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID(hostame)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	mClient = mqtt.NewClient(opts)

	if token := mClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if listen!=nil && *listen != ""{
		subscribe(*listen)
		listening()
	}

	exit()
}

func listening(){

	for mClient.IsConnected(){
		time.Sleep(1* time.Second)

	}
	fmt.Println("Connection lost")
}

func sendMsg(topic string, payload string, qos byte, retained bool){
	fmt.Println("Sending Message..")
	token := mClient.Publish(topic, qos, retained, payload)
	token.Wait()
	if token.Error() != nil{
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Println("Message sent")
}

func subscribe(topic string){
	if token := mClient.Subscribe("client/laptop", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func unsubscribe(){
	if token := mClient.Unsubscribe("client/laptop"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func exit(){
	fmt.Println("Sub exit")
	mClient.Disconnect(250)
	os.Exit(1)
}