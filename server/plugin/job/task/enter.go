package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/job/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/job/service"
	sys "github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	retryCount = 3
	t          timer.Timer
	option     []cron.Option
)

func ByFunc(name string, args string) error {
	switch name {
	case "clearTable":
		fallthrough
	default:
		return sys.ClearTable()
	}
}

func ByJob(name string, args string) interface{ Run() } {
	switch name {
	case "job":
		fallthrough
	default:
		return &Example{}
	}
}

// 初始化
func JobSetup() {
	option = append(option, cron.WithSeconds())
	t = global.GVA_Timer

	global.GVA_LOG.Info("[INFO] JobCore Starting...")
	sysJob := service.Service.SysJob
	jobs, err := sysJob.GetAllSysJob()
	if err != nil {
		global.GVA_LOG.Error("[ERROR] JobCore init error", zap.Error(err))
	}
	if len(jobs) == 0 {
		global.GVA_LOG.Info("[INFO] JobCore total:0")
	}

	err = sysJob.RemoveAllEntryID()
	if err != nil {
		global.GVA_LOG.Error("[ERROR] JobCore remove entry_id error", zap.Error(err))
	}

	for _, job := range jobs {
		global.GVA_LOG.Info("[INFO] JobCore add job:", zap.Any(job.JobGroup, job.JobName))
		var entryId cron.EntryID
		entryId, err = AddJob(job)
		if err != nil {
			global.GVA_LOG.Error("add timer error:", zap.Error(err))
		} else {
			// 更新EntryID到SysJob
			job.EntryId = int(entryId)
			err = sysJob.UpdateEntryID(job)
			if err != nil {
				global.GVA_LOG.Error("add timer error:", zap.Error(err))
			}
		}
	}
}

func AddJob(job model.SysJob) (entryId cron.EntryID, err error) {
	if job.JobType == 1 {
		// 接口
		if strings.HasPrefix(job.InvokeTarget, "http") {
			// 直接调用
			entryId, err = addHttp(job)
		} else {
			// 根据调用目标获得task下的接口
			entryId, err = addJob(job)
		}
	} else {
		// 函数
		entryId, err = addFunc(job)
	}
	return
}

func addHttp(job model.SysJob) (entryId cron.EntryID, err error) {
	startTime := time.Now()
	var obj = job.InvokeTarget
	if obj == "" {
		global.GVA_LOG.Warn("[Job] ExecJob Run job nil")
	} else {
		entryId, err = t.AddTaskByFunc(job.JobGroup, job.CronExpression, func() {
			startedTime := time.Now()
			var count = 0
			var str string
			/* 循环 */
		LOOP:
			if count < retryCount {
				/* 跳过迭代 */
				str, err = request.Get(job.InvokeTarget)
				if err != nil {
					// 如果失败暂停一段时间重试
					global.GVA_LOG.Error("[ERROR] mission failed! ", zap.Error(err))
					global.GVA_LOG.Info(fmt.Sprintf("[INFO] Retry after the task fails %d seconds! %s \n", (count+1)*5, str))
					time.Sleep(time.Duration(count+1) * 5 * time.Second)
					count = count + 1
					goto LOOP
				}
			}
			// 结束时间
			endedTime := time.Now()

			// 执行时间
			latencyedTime := endedTime.Sub(startedTime)

			global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s exec success , spend :%v", job.JobName, latencyedTime))
		}, job.JobName, option...)
		if err != nil {
			// 如果失败暂停一段时间重试
			global.GVA_LOG.Error("[ERROR] mission failed! ", zap.Error(err))
		}
		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//TODO: 待完善部分
		//str := time.Now().Format(timeFormat) + " [INFO] JobCore " + string(e.EntryId) + "exec success , spend :" + latencyTime.String()
		//ws.SendAll(str)
		global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s add success, spend: %.2fs ", job.JobName, latencyTime.Seconds()))
	}
	return
}

func addJob(job model.SysJob) (entryId cron.EntryID, err error) {
	startTime := time.Now()
	var obj = job.InvokeTarget
	if obj == "" {
		global.GVA_LOG.Warn("[Job] ExecJob Run job nil")
	} else {
		entryId, err = t.AddTaskByJob(job.JobGroup, job.CronExpression, ByJob(job.InvokeTarget, job.Args), job.JobName, option...)
		if err != nil {
			// 如果失败暂停一段时间重试
			global.GVA_LOG.Error("[ERROR] mission failed! ", zap.Error(err))
		}
		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//TODO: 待完善部分
		//str := time.Now().Format(timeFormat) + " [INFO] JobCore " + string(e.EntryId) + "exec success , spend :" + latencyTime.String()
		//ws.SendAll(str)
		global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s add success, spend: %.2fs ", job.JobName, latencyTime.Seconds()))
	}
	return
}

func addFunc(job model.SysJob) (entryId cron.EntryID, err error) {
	startTime := time.Now()
	var obj = job.InvokeTarget
	if obj == "" {
		global.GVA_LOG.Warn("[Job] ExecJob Run job nil")
	} else {
		entryId, err = t.AddTaskByFunc(job.JobGroup, job.CronExpression, func() {
			err = ByFunc(job.InvokeTarget, job.Args)
		}, job.JobName, option...)
		if err != nil {
			// 如果失败暂停一段时间重试
			global.GVA_LOG.Error("[ERROR] mission failed! ", zap.Error(err))
		}
		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//TODO: 待完善部分
		//str := time.Now().Format(timeFormat) + " [INFO] JobCore " + string(e.EntryId) + "exec success , spend :" + latencyTime.String()
		//ws.SendAll(str)
		global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s add success, spend: %.2fs ", job.JobName, latencyTime.Seconds()))
	}
	return
}

func RemoveJob(job model.SysJob) {
	startTime := time.Now()
	global.GVA_Timer.RemoveTaskByName(job.JobGroup, job.JobName)
	// 结束时间
	endTime := time.Now()
	// 执行时间
	latencyTime := endTime.Sub(startTime)
	global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s remove success, spend: %.2fs ", job.JobName, latencyTime.Seconds()))
}

func StartJob(job model.SysJob) {
	startTime := time.Now()

	global.GVA_Timer.StartCron(job.JobGroup)
	// 结束时间
	endTime := time.Now()
	// 执行时间
	latencyTime := endTime.Sub(startTime)
	global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s start success, spend: %.2fs ", job.JobName, latencyTime.Seconds()))
}

func StopJob(job model.SysJob) {
	startTime := time.Now()
	global.GVA_Timer.StopCron(job.JobGroup)
	// 结束时间
	endTime := time.Now()
	// 执行时间
	latencyTime := endTime.Sub(startTime)
	global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s stop success, spend: %.2fs ", job.JobName, latencyTime.Seconds()))
}

func ClearJob(job model.SysJob) {
	startTime := time.Now()
	global.GVA_Timer.Clear(job.JobGroup)
	// 结束时间
	endTime := time.Now()
	// 执行时间
	latencyTime := endTime.Sub(startTime)
	global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore %s clear success, spend: %.2fs ", job.JobName, latencyTime.Seconds()))
}

func CloseJob() {
	startTime := time.Now()
	global.GVA_Timer.Close()
	// 结束时间
	endTime := time.Now()
	// 执行时间
	latencyTime := endTime.Sub(startTime)
	global.GVA_LOG.Info(fmt.Sprintf("[Job] JobCore close all job success, spend: %.2fs ", latencyTime.Seconds()))
}
