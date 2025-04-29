<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="任务名称:" prop="jobName">
          <el-input v-model="formData.jobName" :clearable="true" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务分组:" prop="jobGroup">
          <el-select v-model="formData.jobGroup" placeholder="请选择任务分组" style="width: 100%" :clearable="true">
            <el-option v-for="(item, key) in jobGroupOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="调用类型:" prop="jobType">
          <el-input v-model.number="formData.jobType" :clearable="true" placeholder="请输入" />
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
          <el-input v-model.number="formData.misfirePolicy" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="是否并发:" prop="concurrent">
          <el-input v-model.number="formData.concurrent" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="状态:" prop="status">
          <el-input v-model.number="formData.status" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
  import { createSysJob, updateSysJob, findSysJob } from '@/api/job/index'

  defineOptions({
    name: 'SysJobForm'
  })

  // 自动获取字典
  import { getDictFunc } from '@/utils/format'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { ref, reactive } from 'vue'

  const route = useRoute()
  const router = useRouter()

  // 提交按钮loading
  const btnLoading = ref(false)

  const type = ref('')
  const jobGroupOptions = ref([])
  const formData = ref({
    jobName: '',
    jobGroup: '',
    jobType: undefined,
    cronExpression: '',
    invokeTarget: '',
    args: '',
    misfirePolicy: undefined,
    concurrent: undefined,
    status: undefined
  })
  // 验证规则
  const rule = reactive({
    jobName: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      }
    ],
    cronExpression: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      }
    ],
    invokeTarget: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      }
    ],
    status: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      }
    ]
  })

  const elFormRef = ref()

  // 初始化方法
  const init = async () => {
    // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSysJob({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    jobGroupOptions.value = await getDictFunc('jobGroup')
  }

  init()
  // 保存按钮
  const save = async () => {
    btnLoading.value = true
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return (btnLoading.value = false)
      let res
      switch (type.value) {
        case 'create':
          res = await createSysJob(formData.value)
          break
        case 'update':
          res = await updateSysJob(formData.value)
          break
        default:
          res = await createSysJob(formData.value)
          break
      }
      btnLoading.value = false
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '创建/更改成功'
        })
      }
    })
  }

  // 返回按钮
  const back = () => {
    router.go(-1)
  }
</script>

<style></style>
