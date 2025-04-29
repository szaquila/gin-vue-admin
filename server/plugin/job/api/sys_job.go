package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/job/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/job/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/job/task"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var SysJob = new(sysJob)

type sysJob struct{}

// CreateSysJob 创建定时任务
// @Tags SysJob
// @Summary 创建定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysJob true "创建定时任务"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysJob/createSysJob [post]
func (a *sysJob) CreateSysJob(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var info model.SysJob
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceSysJob.CreateSysJob(ctx, &info)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	if info.Status == 1 {
		// 启动任务
		task.AddJob(info)
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteSysJob 删除定时任务
// @Tags SysJob
// @Summary 删除定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysJob true "删除定时任务"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysJob/deleteSysJob [delete]
func (a *sysJob) DeleteSysJob(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	resysJob, err := serviceSysJob.GetSysJob(ctx, ID)
	if err == nil {
		if resysJob.Status == 1 {
			// 停止任务
			task.RemoveJob(resysJob)
		}
	}
	err = serviceSysJob.DeleteSysJob(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSysJobByIds 批量删除定时任务
// @Tags SysJob
// @Summary 批量删除定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysJob/deleteSysJobByIds [delete]
func (a *sysJob) DeleteSysJobByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	var pageInfo request.SysJobSearch
	pageInfo.Page = 1
	pageInfo.PageSize = 99999
	list, _, err := serviceSysJob.GetSysJobInfoList(ctx, pageInfo)
	if err == nil {
		for _, v := range list {
			if v.Status == 1 {
				// 停止任务
				task.RemoveJob(v)
			}
		}
	}
	err = serviceSysJob.DeleteSysJobByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSysJob 更新定时任务
// @Tags SysJob
// @Summary 更新定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysJob true "更新定时任务"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysJob/updateSysJob [put]
func (a *sysJob) UpdateSysJob(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var info model.SysJob
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	resysJob, err := serviceSysJob.GetSysJob(ctx, strconv.FormatUint(uint64(info.ID), 10))
	if err == nil {
		if resysJob.Status == 1 {
			// 停止任务
			task.RemoveJob(resysJob)
		}
	}
	err = serviceSysJob.UpdateSysJob(ctx, info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	if info.Status == 1 {
		// 启动任务
		task.AddJob(info)
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysJob 用id查询定时任务
// @Tags SysJob
// @Summary 用id查询定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询定时任务"
// @Success 200 {object} response.Response{data=model.SysJob,msg=string} "查询成功"
// @Router /sysJob/findSysJob [get]
func (a *sysJob) FindSysJob(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	resysJob, err := serviceSysJob.GetSysJob(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(resysJob, c)
}

// GetSysJobList 分页获取定时任务列表
// @Tags SysJob
// @Summary 分页获取定时任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.SysJobSearch true "分页获取定时任务列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysJob/getSysJobList [get]
func (a *sysJob) GetSysJobList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo request.SysJobSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceSysJob.GetSysJobInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetSysJobPublic 不需要鉴权的定时任务接口
// @Tags SysJob
// @Summary 不需要鉴权的定时任务接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysJob/getSysJobPublic [get]
func (a *sysJob) GetSysJobPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	serviceSysJob.GetSysJobPublic(ctx)
	response.OkWithDetailed(gin.H{"info": "不需要鉴权的定时任务接口信息"}, "获取成功", c)
}
