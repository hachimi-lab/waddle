package logic

import (
	"github.com/dobyte/due/cluster/node"
	"github.com/dobyte/due/log"
	"github.com/gogo/protobuf/proto"
	"github.com/hachimi-lab/waddle/internal/pb"
	"reflect"
)

type Login struct {
	*node.Proxy
}

func LoginModule(proxy *node.Proxy) *Login {
	return &Login{Proxy: proxy}
}

func (slf *Login) Init() {
	slf.Router().Group(func(group *node.RouterGroup) {
		group.AddRouteHandler(int32(pb.MsgId_Login), false, wrap(slf.doLogin))
		group.AddRouteHandler(int32(pb.MsgId_Greet), false, wrap(slf.doGreet))
	})
}

func wrap[T, S proto.Message](fn func(ctx *node.Context, in T, out S)) func(ctx *node.Context) {
	return func(ctx *node.Context) {
		reqType := reflect.TypeOf((*T)(nil)).Elem().Elem()
		ackType := reflect.TypeOf((*S)(nil)).Elem().Elem()
		noneType := reflect.TypeOf((*pb.None)(nil)).Elem()

		req := reflect.New(reqType).Interface().(T)
		ack := reflect.New(ackType).Interface().(S)

		if reqType != noneType {
			if err := ctx.Request.Parse(req); err != nil {
				log.Errorf("parse message failed: %v", err)
				_ = ctx.Disconnect(true)
				return
			}
		}
		if ackType != noneType {
			defer func() {
				if err := ctx.Response(ack); err != nil {
					log.Errorf("response message failed: %v", err)
				}
			}()
		}

		fn(ctx, req, ack)
	}
}

func (slf *Login) doLogin(ctx *node.Context, req *pb.LoginReq, ack *pb.LoginAck) {
	if req.GetAccount() != 100000 && req.GetPassword() != "123456" {
		ack.Code = pb.LoginCode_IncorrectAccountOrPassword
		return
	}
	if err := ctx.BindGate(req.Account); err != nil {
		log.Errorf("bind gate failed: %v", err)
		ack.Code = pb.LoginCode_Failed
		return
	}
	ack.Code = pb.LoginCode_OK
	ack.Account = req.Account
	ack.Token = "token"
}

func (slf *Login) doGreet(ctx *node.Context, req *pb.GreetReq, ack *pb.None) {
	log.Infof("hi %s", req.GetName())
}
