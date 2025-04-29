package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysJobSearch struct {
	JobName  *string `json:"jobName" form:"jobName"`
	JobGroup *string `json:"jobGroup" form:"jobGroup"`
	Status   *int    `json:"status" form:"status"`
	request.PageInfo
}
