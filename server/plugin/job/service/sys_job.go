package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/job/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/job/model/request"
)

var SysJob = new(sysJob)

type sysJob struct{}

// CreateSysJob 创建定时任务记录
// Author [yourname](https://github.com/yourname)
func (s *sysJob) CreateSysJob(ctx context.Context, sysJob *model.SysJob) (err error) {
	err = global.GVA_DB.Create(sysJob).Error
	return err
}

// DeleteSysJob 删除定时任务记录
// Author [yourname](https://github.com/yourname)
func (s *sysJob) DeleteSysJob(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&model.SysJob{}, "id = ?", ID).Error
	return err
}

// DeleteSysJobByIds 批量删除定时任务记录
// Author [yourname](https://github.com/yourname)
func (s *sysJob) DeleteSysJobByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]model.SysJob{}, "id in ?", IDs).Error
	return err
}

// UpdateSysJob 更新定时任务记录
// Author [yourname](https://github.com/yourname)
func (s *sysJob) UpdateSysJob(ctx context.Context, sysJob model.SysJob) (err error) {
	err = global.GVA_DB.Model(&model.SysJob{}).Where("id = ?", sysJob.ID).Updates(&sysJob).Error
	return err
}

// GetSysJob 根据ID获取定时任务记录
// Author [yourname](https://github.com/yourname)
func (s *sysJob) GetSysJob(ctx context.Context, ID string) (sysJob model.SysJob, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysJob).Error
	return
}

// GetSysJobInfoList 分页获取定时任务记录
// Author [yourname](https://github.com/yourname)
func (s *sysJob) GetSysJobInfoList(ctx context.Context, info request.SysJobSearch) (list []model.SysJob, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.SysJob{})
	var sysJobs []model.SysJob
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.JobName != nil && *info.JobName != "" {
		db = db.Where("job_name LIKE ?", "%"+*info.JobName+"%")
	}
	if info.JobGroup != nil && *info.JobGroup != "" {
		db = db.Where("job_group = ?", *info.JobGroup)
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&sysJobs).Error
	return sysJobs, total, err
}

func (s *sysJob) GetSysJobPublic(ctx context.Context) {
}

// GetAllSysJob 获取所有有效的定时任务记录
func (s *sysJob) GetAllSysJob() (list []model.SysJob, err error) {
	err = global.GVA_DB.Model(&model.SysJob{}).Where("status=?", "1").Find(&list).Error
	return
}

func (s *sysJob) RemoveAllEntryID() (err error) {
	err = global.GVA_DB.Model(&model.SysJob{}).Where("entry_id > ?", 0).Update("entry_id", 0).Error
	return
}

// UpdateEntryID 更新定时任务记录
// Author [yourname](https://github.com/yourname)
func (s *sysJob) UpdateEntryID(sysJob model.SysJob) (err error) {
	err = global.GVA_DB.Model(&model.SysJob{}).Where("id = ?", sysJob.ID).Updates(&sysJob).Error
	return err
}
