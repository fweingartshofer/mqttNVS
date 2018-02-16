package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	subscribe string
	mClient   mqtt.Client
)

var ch mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topics := strings.Split(msg.Topic(), "/")
	fmt.Printf("[ %s ] \n", time.Now().Format(time.RFC822Z))
	fmt.Printf("Message received from %s\n", topics[1])
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("Qos: %v\n", msg.Qos())
	fmt.Println("______________________________")
	fmt.Printf("MSG:\n%s\n", msg.Payload())
	fmt.Println("==============================")
	if len(topics) > 2 && topics[2] == "img" {
		//saveImg(msg)
	}
}

func saveImg(msg mqtt.Message) {
	imgstr := strings.Split(string(msg.Payload()), ",")
	fmt.Println(len(imgstr))
	myimg := image.NewRGBA64(image.Rect(0, 0, 1920, 1080))
	for i := 0; i < len(imgstr); i++ {
		k, _ := strconv.ParseUint(imgstr[i], 10, 8)
		myimg.Pix[i] = uint8(k)
	}
	f, err := os.Create("./test.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, myimg)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func sub(topic string) mqtt.Token {
	token := mClient.Subscribe(topic, 0, nil)
	token.Wait()
	return token

}

func createClient() {
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	hostame, _ := os.Hostname()
	hostame = "Master_" + hostame
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID(hostame)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(ch)
	opts.SetPingTimeout(1 * time.Second)
	opts.AutoReconnect = true
	mClient = mqtt.NewClient(opts)

	if token := mClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func main() {
	subscribe = "client/#"
	createClient()
	err := sub(subscribe)
	if (err).Error() != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for {
		time.Sleep(1 * time.Second)
		if !mClient.IsConnected() {
			mClient.Connect()
		}
	}
	fmt.Println("Press Enter to Exit.")
	os.Stdin.Read([]byte{0})
	mClient.Disconnect(10)
}
