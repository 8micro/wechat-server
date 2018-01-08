package server

import "gopkg.in/chanxuehong/wechat.v2/mp/core"
import "gopkg.in/chanxuehong/wechat.v2/mp/message/callback/request"
import "gopkg.in/chanxuehong/wechat.v2/mp/message/callback/response"

import (
	"fmt"
)

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
