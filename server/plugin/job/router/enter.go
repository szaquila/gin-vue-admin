package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/job/api"

var (
	Router    = new(router)
	apiSysJob = api.Api.SysJob
)

type router struct{ SysJob sysJob }
