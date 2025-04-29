package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SysJob 定时任务 结构体
type SysJob struct {
	global.GVA_MODEL
	JobName        string `json:"jobName" form:"jobName" gorm:"column:job_name;" binding:"required"`                      //任务名称
	JobGroup       string `json:"jobGroup" form:"jobGroup" gorm:"column:job_group;"`                                      //任务分组
	JobType        int    `json:"jobType" form:"jobType" gorm:"column:job_type;"`                                         //调用类型
	CronExpression string `json:"cronExpression" form:"cronExpression" gorm:"column:cron_expression;" binding:"required"` //定时表达式
	InvokeTarget   string `json:"invokeTarget" form:"invokeTarget" gorm:"column:invoke_target;" binding:"required"`       //调用目标
	Args           string `json:"args" form:"args" gorm:"column:args;"`                                                   //参数
	MisfirePolicy  int    `json:"misfirePolicy" form:"misfirePolicy" gorm:"column:misfire_policy;"`                       //执行策略
	Concurrent     int    `json:"concurrent" form:"concurrent" gorm:"column:concurrent;"`                                 //是否并发
	Status         int    `json:"status" form:"status" gorm:"column:status;" binding:"required"`                          //状态
	EntryId        int    `json:"entryId" form:"entryId" gorm:"column:entry_id;"`                                         //入口编码
}

// TableName 定时任务 SysJob自定义表名 sys_job
func (SysJob) TableName() string {
	return "sys_job"
}
