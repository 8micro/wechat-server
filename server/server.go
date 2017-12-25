package server

import "github.com/8micro/wechat-server/security"
import "github.com/8micro/gutils/httpx"
import "gopkg.in/chanxuehong/wechat.v2/mp/core"
import "gopkg.in/chanxuehong/wechat.v2/mp/message/callback/request"
import "gopkg.in/chanxuehong/wechat.v2/mp/message/callback/response"

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

//CoreServer is exported
type CoreServer struct {
	core.ErrorHandler
	Key               string
	client            *httpx.HttpClient
	accessTokenServer *core.DefaultAccessTokenServer
	server            *core.Server
}

//NewCoreServer is exported
func NewCoreServer(key string, securityData *security.SecurityData) *CoreServer {

	coreServer := &CoreServer{
		Key: key,
		client: httpx.NewClient().
			SetTransport(&http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext,
				DisableKeepAlives:     false,
				MaxIdleConns:          25,
				MaxIdleConnsPerHost:   25,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout,
				ExpectContinueTimeout: http.DefaultTransport.(*http.Transport).ExpectContinueTimeout,
			}),
	}

	serveMux := core.NewServeMux()
	serveMux.MsgHandleFunc(request.MsgTypeText, coreServer.msgTextHandleFunc)
	serveMux.EventHandleFunc(request.EventTypeSubscribe, coreServer.eventSubscribeHandleFunc)
	serveMux.EventHandleFunc(request.EventTypeUnsubscribe, coreServer.eventUnsubscribeHandleFunc)
	coreServer.server = core.NewServer(securityData.OriginID, securityData.AppID, securityData.Token, securityData.AESKey, serveMux, coreServer)
	coreServer.accessTokenServer = core.NewDefaultAccessTokenServer(securityData.AppID, securityData.AppSecret, coreServer.client.RawClient())
	return coreServer
}

//Signature is exported
//return access wechat server signature string.
func (coreServer *CoreServer) Signature() (string, error) {

	return coreServer.accessTokenServer.Token()
}

//Server is exported
//return wechat core server.
func (coreServer *CoreServer) Server() *core.Server {

	return coreServer.server
}

//ServeError is exported
func (coreServer *CoreServer) ServeError(w http.ResponseWriter, r *http.Request, err error) {

	fmt.Printf("### ServeError:%s\n", err)
}

//msgTextHandleFunc is exported
//text message handle
func (coreServer *CoreServer) msgTextHandleFunc(context *core.Context) {

	//根据用户输入, search 相关关键字的文章并返回.
	reqText := request.GetText(context.MixedMsg)
	fmt.Printf("#msg -> %s:%s\n", reqText.FromUserName, reqText.Content)
	respText := response.NewText(reqText.FromUserName, reqText.ToUserName, reqText.CreateTime, "欢迎! Hello World.")
	context.RawResponse(respText) // 明文模式回复
}

func (coreServer *CoreServer) eventSubscribeHandleFunc(context *core.Context) {

	//用户微信关注事件，获取用户信息，记录用户信息到数据库
	//若wechat表已存在则更新一次.
	//subscribeEvent := request.GetSubscribeEvent(context.MixedMsg)

}

func (coreServer *CoreServer) eventUnsubscribeHandleFunc(context *core.Context) {

	//用户微信取消关注事件，
	//unsubscribeEvent := request.GetUnsubscribeEvent(context.MixedMsg)
	//不删除微信信息，只修改状态.
}
