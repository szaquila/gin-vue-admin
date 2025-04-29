package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	base "github.com/flipped-aurora/gin-vue-admin/server/plugin/job/task"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"go.uber.org/zap"

	"github.com/robfig/cron/v3"
)

func Timer() {
	go func() {
		base.JobSetup()

		var option []cron.Option
		option = append(option, cron.WithSeconds())
		t := global.GVA_Timer
		if _, running := t.FindTask("SYSTEM", "定时清理数据库【日志, 黑名单】内容"); !running {
			global.GVA_LOG.Info("添加定时任务: 定时清理数据库【日志, 黑名单】内容")
			// 清理DB定时任务
			_, err := global.GVA_Timer.AddTaskByFunc("定时清理数据库【日志, 黑名单】内容", "@daily", func() {
				err := task.ClearTable() // 定时任务方法定在task文件包中
				if err != nil {
					global.GVA_LOG.Error("timer error:", zap.Error(err))
				}
			}, "SYSTEM", option...)
			if err != nil {
				global.GVA_LOG.Error("add timer error:", zap.Error(err))
			}
		}

		// 其他定时任务定在这里 参考上方使用方法

		//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "corn表达式", func() {
		//	具体执行内容...
		//  ......
		//}, option...)
		//if err != nil {
		//	global.GVA_LOG.Error("add timer error:", zap.Error(err))
		//}
	}()
}
