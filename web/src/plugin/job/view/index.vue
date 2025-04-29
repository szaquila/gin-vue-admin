<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="任务名称" prop="jobName">
          <el-input v-model="searchInfo.jobName" placeholder="搜索条件" />
        </el-form-item>

        <el-form-item label="任务分组" prop="jobGroup">
          <el-select v-model="searchInfo.jobGroup" clearable filterable placeholder="请选择" @clear="searchInfo.jobGroup = undefined">
            <el-option v-for="(item, key) in jobGroupOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" clearable filterable placeholder="请选择" @clear="searchInfo.status = undefined">
            <el-option v-for="(item, key) in statusOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button icon="delete" style="margin-left: 10px" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="任务名称" prop="jobName" width="280" />
        <el-table-column align="left" label="任务分组" prop="jobGroup" width="120">
          <template #default="scope">
            {{ filterDict(scope.row.jobGroup, jobGroupOptions) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="定时表达式" prop="cronExpression" width="120" />
        <el-table-column align="left" label="调用目标" prop="invokeTarget" width="160" />
        <el-table-column align="left" label="状态" prop="status" width="120">
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              inline-prompt
              :active-value="1"
              active-text="正常"
              :inactive-value="2"
              inactive-text="停用"
              @change="switchEnable(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateSysJobFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="任务名称:" prop="jobName">
          <el-input v-model="formData.jobName" :clearable="true" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务分组:" prop="jobGroup">
          <el-select v-model="formData.jobGroup" placeholder="请选择任务分组" style="width: 100%" filterable :clearable="true">
            <el-option v-for="(item, key) in jobGroupOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="调用类型:" prop="jobType">
          <el-radio-group v-model="formData.jobType" placeholder="请选择调用类型" style="width: 100%" filterable :clearable="true">
            <el-radio-button v-for="(item, key) in jobTypeOptions" :key="key" :label="item.label" :value="item.value" />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="定时表达式:" prop="cronExpression">
          <el-input v-model="formData.cronExpression" :clearable="true" placeholder="请输入定时表达式" />
        </el-form-item>
        <el-form-item label="调用目标:" prop="invokeTarget">
          <el-input v-model="formData.invokeTarget" :clearable="true" placeholder="请输入调用目标" />
        </el-form-item>
        <el-form-item label="参数:" prop="args">
          <el-input v-model="formData.args" :clearable="true" placeholder="请输入参数" />
        </el-form-item>
        <el-form-item label="执行策略:" prop="misfirePolicy">
          <el-radio-group v-model="formData.misfirePolicy" placeholder="请选择执行策略" style="width: 100%" filterable :clearable="true">
            <el-radio-button v-for="(item, key) in misfirePolicyOptions" :key="key" :label="item.label" :value="item.value" />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="是否并发:" prop="concurrent">
          <el-radio-group v-model="formData.concurrent" placeholder="请选择是否并发" style="width: 100%" filterable :clearable="true">
            <el-radio-button v-for="(item, key) in concurrentOptions" :key="key" :label="item.label" :value="item.value" />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态:" prop="status">
          <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%" filterable :clearable="true">
            <el-option v-for="(item, key) in statusOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="任务名称">
          {{ detailFrom.jobName }}
        </el-descriptions-item>
        <el-descriptions-item label="任务分组">
          {{ filterDict(detailFrom.jobGroup, jobGroupOptions) }}
        </el-descriptions-item>
        <el-descriptions-item label="调用类型">
          {{ filterDict(detailFrom.jobType.toString(), jobTypeOptions) }}
        </el-descriptions-item>
        <el-descriptions-item label="定时表达式">
          {{ detailFrom.cronExpression }}
        </el-descriptions-item>
        <el-descriptions-item label="调用目标">
          {{ detailFrom.invokeTarget }}
        </el-descriptions-item>
        <el-descriptions-item label="参数">
          {{ detailFrom.args }}
        </el-descriptions-item>
        <el-descriptions-item label="执行策略">
          {{ filterDict(detailFrom.misfirePolicy.toString(), misfirePolicyOptions) }}
        </el-descriptions-item>
        <el-descriptions-item label="是否并发">
          {{ filterDict(detailFrom.concurrent.toString(), concurrentOptions) }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          {{ filterDict(detailFrom.status.toString(), statusOptions) }}
        </el-descriptions-item>
        <el-descriptions-item label="入口编码">
          {{ detailFrom.entryId }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
  import { createSysJob, deleteSysJob, deleteSysJobByIds, updateSysJob, findSysJob, getSysJobList } from '@/plugin/job/api/index'

  // 全量引入格式化工具 请按需保留
  import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { nextTick, ref, reactive } from 'vue'

  defineOptions({
    name: 'SysJob'
  })

  // 提交按钮loading
  const btnLoading = ref(false)

  // 控制更多查询条件显示/隐藏状态
  const showAllQuery = ref(false)

  // 自动化生成的字典（可能为空）以及字段
  const jobGroupOptions = ref([])
  const jobTypeOptions = ref([])
  const concurrentOptions = ref([])
  const misfirePolicyOptions = ref([])
  const statusOptions = ref([])
  const formData = ref({
    jobName: '',
    jobGroup: 'DEFAULT',
    jobType: '1',
    cronExpression: '',
    invokeTarget: '',
    args: '',
    misfirePolicy: '1',
    concurrent: '1',
    status: '1'
  })

  // 验证规则
  const rule = reactive({
    jobName: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: '不能只输入空格',
        trigger: ['input', 'blur']
      }
    ],
    cronExpression: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: '不能只输入空格',
        trigger: ['input', 'blur']
      }
    ],
    invokeTarget: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: '不能只输入空格',
        trigger: ['input', 'blur']
      }
    ]
  })

  const elFormRef = ref()
  const elSearchFormRef = ref()

  // =========== 表格控制部分 ===========
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})
  // 重置
  const onReset = () => {
    searchInfo.value = {}
    getTableData()
  }

  // 搜索
  const onSubmit = () => {
    elSearchFormRef.value?.validate(async (valid) => {
      if (!valid) return
      page.value = 1
      getTableData()
    })
  }

  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  // 修改页面容量
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // 查询
  const getTableData = async () => {
    const table = await getSysJobList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  // ============== 表格控制部分结束 ===============

  // 获取需要的字典 可能为空 按需保留
  const setOptions = async () => {
    jobGroupOptions.value = await getDictFunc('jobGroup')
    jobTypeOptions.value = await getDictFunc('jobType')
    concurrentOptions.value = await getDictFunc('concurrent')
    misfirePolicyOptions.value = await getDictFunc('misfirePolicy')
    statusOptions.value = await getDictFunc('status')
  }

  // 获取需要的字典 可能为空 按需保留
  setOptions()

  // 多选数据
  const multipleSelection = ref([])
  // 多选
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  // 删除行
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      deleteSysJobFunc(row)
    })
  }

  // 多选删除
  const onDelete = async () => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map((item) => {
          IDs.push(item.ID)
        })
      const res = await deleteSysJobByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }

  // 行为控制标记（弹窗内部需要增还是改）
  const type = ref('')

  // 更新行
  const updateSysJobFunc = async (row) => {
    const res = await findSysJob({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      res.data.jobType = res.data.jobType.toString()
      res.data.concurrent = res.data.concurrent.toString()
      res.data.misfirePolicy = res.data.misfirePolicy.toString()
      res.data.status = res.data.status.toString()
      formData.value = res.data
      dialogFormVisible.value = true
    }
  }

  // 删除行
  const deleteSysJobFunc = async (row) => {
    const res = await deleteSysJob({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }

  // 弹窗控制标记
  const dialogFormVisible = ref(false)

  // 打开弹窗
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  // 关闭弹窗
  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
      jobName: '',
      jobGroup: 'DEFAULT',
      jobType: '1',
      cronExpression: '',
      invokeTarget: '',
      args: '',
      misfirePolicy: '1',
      concurrent: '1',
      status: '1'
    }
  }
  // 弹窗确定
  const enterDialog = async () => {
    btnLoading.value = true
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return (btnLoading.value = false)
      let data = formData.value
      data.jobType = Number(data.jobType)
      data.concurrent = Number(data.concurrent)
      data.misfirePolicy = Number(data.misfirePolicy)
      data.status = Number(data.status)
      let res
      switch (type.value) {
        case 'create':
          res = await createSysJob(data)
          break
        case 'update':
          res = await updateSysJob(data)
          break
        default:
          res = await createSysJob(data)
          break
      }
      btnLoading.value = false
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '创建/更改成功'
        })
        closeDialog()
        getTableData()
      }
    })
  }

  const detailFrom = ref({})

  // 查看详情控制标记
  const detailShow = ref(false)

  // 打开详情弹窗
  const openDetailShow = () => {
    detailShow.value = true
  }

  // 打开详情
  const getDetails = async (row) => {
    // 打开弹窗
    const res = await findSysJob({ ID: row.ID })
    if (res.code === 0) {
      res.data.jobType = res.data.jobType.toString()
      res.data.concurrent = res.data.concurrent.toString()
      res.data.misfirePolicy = res.data.misfirePolicy.toString()
      res.data.status = res.data.status.toString()
      detailFrom.value = res.data
      openDetailShow()
    }
  }

  // 关闭详情弹窗
  const closeDetailShow = () => {
    detailShow.value = false
    detailFrom.value = {}
  }

  const switchEnable = async (row) => {
    formData.value = JSON.parse(JSON.stringify(row))
    await nextTick()
    const data = {
      ...formData.value
    }
    const res = await updateSysJob(data)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: `${data.enable === 2 ? '禁用' : '启用'}成功`
      })
      await getTableData()
    }
  }
</script>

<style></style>
