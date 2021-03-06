package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateIgnoreAction actions.Action

// 修改屏蔽的Header
func (this *UpdateIgnoreAction) Run(params struct {
	From       string
	Server     string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string
	Name       string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["from"] = params.From
	this.Data["server"] = maps.Map{
		"filename": server.Filename,
	}
	this.Data["locationId"] = params.LocationId
	this.Data["rewriteId"] = params.RewriteId
	this.Data["fastcgiId"] = params.FastcgiId
	this.Data["backendId"] = params.BackendId
	this.Data["name"] = params.Name

	this.Show()
}

// 提交修改
func (this *UpdateIgnoreAction) RunPost(params struct {
	Server     string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string
	OldName    string
	Name       string
	Must       *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入Name")

	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	headerList, err := server.FindHeaderList(params.LocationId, params.BackendId, params.RewriteId, params.FastcgiId)
	if err != nil {
		this.Fail(err.Error())
	}
	headerList.UpdateIgnoreHeader(params.OldName, params.Name)

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
