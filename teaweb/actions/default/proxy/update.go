package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 修改代理服务信息
func (this *UpdateAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}
	this.Data["proxy"] = server
	this.Data["filename"] = server.Filename
	this.Data["selectedTab"] = "basic"

	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Show()
}

// 保存提交
func (this *UpdateAction) RunPost(params struct {
	HttpOn      bool
	Server      string
	Description string
	Name        []string
	Listen      []string
	Root        string
	Charset     string
	Index       []string
	Must        *actions.Must
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	params.Must.
		Field("description", params.Description).
		Require("代理服务名称不能为空")

	server.Http = params.HttpOn
	server.Description = params.Description
	server.Name = params.Name
	server.Listen = params.Listen
	server.Root = params.Root
	server.Charset = params.Charset
	server.Index = params.Index
	err = server.Validate()
	if err != nil {
		this.Fail("校验失败：" + err.Error())
	}

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重启
	proxyutils.NotifyChange()

	this.Success()
}
