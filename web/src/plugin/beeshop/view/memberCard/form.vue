<template>
  <el-drawer destroy-on-close size="800" :model-value="props.open" :show-close="true" :before-close="closeDialog">
    <template #header>
      <div v-if="!props.view" class="flex justify-between items-center">
        <span class="text-lg">{{ id == 0 ? '添加' : '修改' }}</span>
        <div>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
          <el-button @click="closeDialog">取 消</el-button>
        </div>
      </div>
    </template>

    <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rules" label-width="80px" :disabled="view">
      <el-form-item label="会员卡名称:" prop="name">
        <el-input v-model="formData.name" :clearable="true" placeholder="请输入会员卡名称" />
      </el-form-item>
      <el-form-item label="会员卡金额:" prop="amount">
        <el-input v-model.number="formData.amount" type="digit" :clearable="true" placeholder="请输入会员卡金额">
          <template #append>
            <span>元</span>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="排序:" prop="sortNum">
        <el-input v-model.number="formData.sortNum" type="number" :clearable="true" placeholder="请输入会员卡排序" />
      </el-form-item>
      <el-form-item label="会员卡类型:" prop="validMonth">
        <el-select v-model="formData.validMonth" placeholder="请选择" :clearable="false">
          <el-option v-for="(item, key) in memberCardTypes" :key="key" :label="item.label"
            :value="parseInt(item.value)" />
        </el-select>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script setup>

// 全量引入格式化工具 请按需保留
import { getDictDetailList } from '@/api/sysDictionaryDetail'
import { memberCardSave, infoById } from '@/plugin/beeshop/api/beeMemberCard'
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
defineOptions({
  name: 'MemberCardForm'
})
const emits = defineEmits(['update:open', 'refresh']);
const props = defineProps({
  open: {
    type: Boolean,
    default: false,
    required: true
  },
  id: {
    type: Number,
    default: 0,
    required: false
  },
  view: {
    type: Boolean,
    required: false,
    default: true
  }
})
watch(()=>props.open, (newV) => {
  if (newV && props.id > 0) {
    infoById(props.id).then(res => {
      if (res.code == 0) {
        formData.value = res.data
      }
    })
  } else {
    formData.value = {}
  }
});
const memberCardTypes = ref([])

const formData = ref({
  amount: 0,
  name: "",
  sortNum: 0,
  validMonth: 0
})


// 验证规则
const rules = reactive({
  name: [{
    required: true, trigger: 'blur', message: '会员卡名称不能为空'
  }],
  amount: [{
    required: true, trigger: 'blur', message: '会员卡金额不能为空'
  }, {
    validator: (rule, value, callback) => {
      if (value < 0) {
        callback('会员卡金额不能小于0');
      } else {
        callback()
      }
    }
  }],
  validMonth: [{ required: true, trigger: 'blur', message: '会员卡类型不能为空' }, {
    validator: (rule, value, callback) => {
      if (value < 0) {
        callback('会员卡分类不能小于0');
      } else {
        callback()
      }
    }
  }]
})

const elFormRef = ref(null)

onMounted(() => {
  getDictDetailList('member_card_type').then(res => {
    if (res.code == 0) {
      memberCardTypes.value = res.data
    }
  })
});

// 弹窗控制标记
const enterDialog = async () => {
  console.log(formData)
  const valid = await elFormRef.value.validate()
  if (valid) {
    const res = await memberCardSave(formData.value)
    if (res.code == 0) {
      ElMessage.success("提交成功")
      closeDialog()
      emits('refresh')
    } else {
      ElMessage.success("提交失败")
    }
  } else {
    ElMessage.error('请完善表单')
  }
}
const closeDialog = () => {
  emits('update:open', false)
}
</script>

<style></style>