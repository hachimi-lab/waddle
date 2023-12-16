package main

import (
	"github.com/dobyte/due"
	gatecluster "github.com/dobyte/due/cluster/gate"
	nodecluster "github.com/dobyte/due/cluster/node"
	"github.com/hachimi-lab/waddle/server/gate/business"

	"github.com/dobyte/due/locate/redis"
	"github.com/dobyte/due/log"
	"github.com/dobyte/due/mode"
	"github.com/dobyte/due/network/ws"
	"github.com/dobyte/due/registry/etcd"
	"github.com/dobyte/due/transport/grpc"
	"net/http"
)

func main() {
	mode.SetMode(mode.DebugMode)
	container := due.NewContainer()
	server := ws.NewServer()
	server.OnUpgrade(func(w http.ResponseWriter, r *http.Request) bool {
		log.Infof("upgrade connection: %s", r.URL.String())
		return true
	})
	locator := redis.NewLocator()
	registry := etcd.NewRegistry()
	transporter := grpc.NewTransporter()
	gate := gatecluster.NewGate(
		gatecluster.WithServer(server),
		gatecluster.WithLocator(locator),
		gatecluster.WithRegistry(registry),
		gatecluster.WithTransporter(transporter),
	)
	node := nodecluster.NewNode(
		nodecluster.WithLocator(locator),
		nodecluster.WithRegistry(registry),
		nodecluster.WithTransporter(transporter),
	)
	business.Register(node.Proxy())

	container.Add(gate, node)
	container.Serve()
}
