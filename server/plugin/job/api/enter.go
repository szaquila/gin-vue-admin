package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/job/service"

var (
	Api           = new(api)
	serviceSysJob = service.Service.SysJob
)

type api struct{ SysJob sysJob }
