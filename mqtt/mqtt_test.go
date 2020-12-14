package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"
)

func TestMQTT(t *testing.T) {
	url := "tcp://172.16.4.17:1833"

	go func() {
		clientID := "s1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		topic := "/a"
		filters := map[string]byte{
			topic: 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(clientID, topic, string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(clientID, topic, "subscribe")
	}()

	go func() {
		clientID := "s2"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		topic := "/a/b"
		filters := map[string]byte{
			topic: 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(clientID, topic, string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(clientID, topic, "subscribe")
	}()

	go func() {
		clientID := "s3"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		topic := "/a/+"
		filters := map[string]byte{
			topic: 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(clientID, topic, string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(clientID, topic, "subscribe")
	}()

	go func() {
		clientID := "s4"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		topic := "/a/+/c"
		filters := map[string]byte{
			topic: 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(clientID, topic, string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(clientID, topic, "subscribe")
	}()

	go func() {
		clientID := "s5"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		topic := "/a/#"
		filters := map[string]byte{
			topic: 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(clientID, topic, string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(clientID, topic, "subscribe")
	}()

	time.Sleep(1 * time.Second)
	go func() {
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID("p1")
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		topic := "/a/b/d"
		if token := client.Publish(topic, 1, false, "bar"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("p1", topic, "publish")
	}()

	select {}
}

func TestMQTT1(t *testing.T) {
	url := "tcp://172.16.4.17:1883"

	go func() {
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID("s1")
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters := map[string]byte{
			"/foo":  1,
			"/foo1": 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters = map[string]byte{
			"/foo":  0,
			"/foo1": 0,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), "sub")
	}()

	time.Sleep(1 * time.Second)
	go func() {
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID("p1")
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Publish("/foo", 1, false, "bar"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), "pub")
	}()

	select {}
}

func TestMqttRetain(t *testing.T) {
	url := "tcp://172.16.4.17:1883"

	/*go func() {
		clientID := "s1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters := map[string]byte{
			"#":     1,
			"/+":    1,
			"/foo1": 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), clientID, message.Topic(), string(message.Payload()))
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "sub")
	}()
	time.Sleep(1 * time.Second)*/

	go func() {
		clientID := "p1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Publish("/foo", 1, true, "bar"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Publish("/foo1", 1, true, "bar"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "pub")
	}()
	time.Sleep(1 * time.Second)

	go func() {
		clientID := "s2"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters := map[string]byte{
			"#":     1,
			"/+":    1,
			"/foo1": 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), clientID, message.Topic(), string(message.Payload()), message.Retained())
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "sub")
	}()

	select {}
}

func TestMqttWill1(t *testing.T) {
	url := "tcp://172.16.4.17:1883"

	go func() {
		clientID := "w1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		opts.SetWill("will", "88", 1, false)
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "connect")
	}()

	select {}
}

func TestMqttWill2(t *testing.T) {
	url := "tcp://172.16.4.17:1883"

	go func() {
		clientID := "s1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters := map[string]byte{
			"will": 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), clientID, message.Topic(), string(message.Payload()), message.Retained())
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "connect")
	}()
	time.Sleep(1 * time.Second)

	go func() {
		clientID := "w1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		opts.SetWill("will", "88", 1, false)
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "connect")
		client.Disconnect(250)
	}()

	select {}
}

func TestMqttSub(t *testing.T) {
	url := "tcp://172.16.4.11:1883"

	go func() {
		clientID := "s1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		// opts.SetCleanSession(false)
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters := map[string]byte{
			"foo":  1,
			"foo1": 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), clientID, message.Topic(), string(message.Payload()), message.Retained())
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "sub")
	}()

	select {}
}

func TestMqttPub(t *testing.T) {
	url := "tcp://172.16.4.17:1883"

	go func() {
		clientID := "p"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Publish("foo", 1, false, "bar"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Publish("foo1", 1, false, "bar1"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "pub")
	}()

	select {}
}

func TestMqttGroup(t *testing.T) {
	url := "tcp://172.16.4.17:1883"

	go func() {
		clientID := "s1"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters := map[string]byte{
			"$share/group/#":   1,
			"$share/group/foo": 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), clientID, message.Topic(), string(message.Payload()), message.Retained())
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "sub")
	}()

	go func() {
		clientID := "s2"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		filters := map[string]byte{
			"$share/group/foo": 1,
		}
		if token := client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println(time.Now().String(), clientID, message.Topic(), string(message.Payload()), message.Retained())
		}); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "sub")
	}()

	time.Sleep(time.Second)
	go func() {
		clientID := "p"
		opts := mqtt.NewClientOptions().AddBroker(url)
		opts.SetClientID(clientID)
		opts.SetUsername("u")
		opts.SetPassword("p")
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Publish("foo", 1, false, "bar"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println(time.Now().String(), clientID, "pub")
	}()

	select {}
}