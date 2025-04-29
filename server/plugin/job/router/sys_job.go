package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var SysJob = new(sysJob)

type sysJob struct{}

// Init 初始化 定时任务 路由信息
func (r *sysJob) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("sysJob").Use(middleware.OperationRecord())
		group.POST("createSysJob", apiSysJob.CreateSysJob)             // 新建定时任务
		group.DELETE("deleteSysJob", apiSysJob.DeleteSysJob)           // 删除定时任务
		group.DELETE("deleteSysJobByIds", apiSysJob.DeleteSysJobByIds) // 批量删除定时任务
		group.PUT("updateSysJob", apiSysJob.UpdateSysJob)              // 更新定时任务
	}
	{
		group := private.Group("sysJob")
		group.GET("findSysJob", apiSysJob.FindSysJob)       // 根据ID获取定时任务
		group.GET("getSysJobList", apiSysJob.GetSysJobList) // 获取定时任务列表
	}
	{
		group := public.Group("sysJob")
		group.GET("getSysJobPublic", apiSysJob.GetSysJobPublic) // 定时任务开放接口
	}
}
