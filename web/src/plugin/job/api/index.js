import service from '@/utils/request'
// @Tags SysJob
// @Summary 创建定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysJob true "创建定时任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysJob/createSysJob [post]
export const createSysJob = (data) => {
  return service({
    url: '/sysJob/createSysJob',
    method: 'post',
    data
  })
}

// @Tags SysJob
// @Summary 删除定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysJob true "删除定时任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysJob/deleteSysJob [delete]
export const deleteSysJob = (params) => {
  return service({
    url: '/sysJob/deleteSysJob',
    method: 'delete',
    params
  })
}

// @Tags SysJob
// @Summary 批量删除定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除定时任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysJob/deleteSysJob [delete]
export const deleteSysJobByIds = (params) => {
  return service({
    url: '/sysJob/deleteSysJobByIds',
    method: 'delete',
    params
  })
}

// @Tags SysJob
// @Summary 更新定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysJob true "更新定时任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysJob/updateSysJob [put]
export const updateSysJob = (data) => {
  return service({
    url: '/sysJob/updateSysJob',
    method: 'put',
    data
  })
}

// @Tags SysJob
// @Summary 用id查询定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysJob true "用id查询定时任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysJob/findSysJob [get]
export const findSysJob = (params) => {
  return service({
    url: '/sysJob/findSysJob',
    method: 'get',
    params
  })
}

// @Tags SysJob
// @Summary 分页获取定时任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取定时任务列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysJob/getSysJobList [get]
export const getSysJobList = (params) => {
  return service({
    url: '/sysJob/getSysJobList',
    method: 'get',
    params
  })
}
// @Tags SysJob
// @Summary 不需要鉴权的定时任务接口
// @Accept application/json
// @Produce application/json
// @Param data query request.SysJobSearch true "分页获取定时任务列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysJob/getSysJobPublic [get]
export const getSysJobPublic = () => {
  return service({
    url: '/sysJob/getSysJobPublic',
    method: 'get'
  })
}
