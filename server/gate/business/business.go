package business

import (
	"github.com/dobyte/due/cluster/node"
	"github.com/hachimi-lab/waddle/internal/logic"
)

func Register(proxy *node.Proxy) {
	logic.LoginModule(proxy).Init()
}
