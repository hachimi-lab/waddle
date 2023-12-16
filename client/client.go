package main

import (
	"github.com/dobyte/due"
	"github.com/dobyte/due/cluster"
	"github.com/dobyte/due/cluster/client"
	"github.com/dobyte/due/log"
	"github.com/dobyte/due/network/ws"
	"github.com/hachimi-lab/waddle/internal/pb"
)

func main() {
	container := due.NewContainer()
	component := client.NewClient(
		client.WithClient(ws.NewClient()),
	)
	initEvent(component.Proxy())
	initRoute(component.Proxy())
	container.Add(component)
	container.Serve()
}

func initEvent(proxy client.Proxy) {
	proxy.AddEventListener(cluster.Connect, onConnect)
	proxy.AddEventListener(cluster.Reconnect, onReconnect)
	proxy.AddEventListener(cluster.Disconnect, onDisconnect)
}

func initRoute(proxy client.Proxy) {
	proxy.AddRouteHandler(int32(pb.MsgId_Login), loginHandler)
}

func onConnect(proxy client.Proxy) {
	log.Infof("connection is opened")
	err := proxy.Push(0, int32(pb.MsgId_Login), &pb.LoginReq{
		Account:  100000,
		Password: "123456",
	})
	if err != nil {
		log.Errorf("push login message failed: %v", err)
	}
}

func onReconnect(proxy client.Proxy) {
	log.Infof("connection is reopened")
	err := proxy.Push(0, int32(pb.MsgId_Login), &pb.LoginReq{
		Account:  100000,
		Password: "123456",
	})
	if err != nil {
		log.Errorf("push login message failed: %v", err)
	}
}

func onDisconnect(proxy client.Proxy) {
	log.Infof("connection is closed")
	err := proxy.Reconnect()
	if err != nil {
		log.Errorf("reconnect failed: %v", err)
	}
}

func loginHandler(request client.Request) {
	ack := &pb.LoginAck{}
	if err := request.Parse(ack); err != nil {
		log.Errorf("parse login message failed: %v", err)
		return
	}
	switch ack.Code {
	case pb.LoginCode_Failed:
		log.Error("user login failed")
		return
	case pb.LoginCode_IncorrectAccountOrPassword:
		log.Error("incorrect account or password")
		return
	default:
		log.Infof("user login successful, token: %v", ack.Token)
	}
	if err := request.Proxy().Push(0, int32(pb.MsgId_Greet), &pb.GreetReq{Name: "Kaka"}); err != nil {
		log.Errorf("push greet message failed: %v", err)
	}
}
