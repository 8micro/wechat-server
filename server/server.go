package server

import "github.com/8micro/wechat-server/security"
import "github.com/8micro/gounits/httpx"
import "gopkg.in/chanxuehong/wechat.v2/mp/core"
import "gopkg.in/chanxuehong/wechat.v2/mp/message/callback/request"

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
