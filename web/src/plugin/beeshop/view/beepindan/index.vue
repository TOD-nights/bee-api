<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
               @keyup.enter="onSubmit">
        <el-form-item label="订单id" prop="id">
          <el-input v-model="searchInfo.id" placeholder="请输入订单id"/>
        </el-form-item>
        <el-form-item label="" prop="orderNumber">
          <el-input v-model="searchInfo.orderNumber" placeholder="请输入订单号"/>
        </el-form-item>
        <el-form-item label="订单状态" prop="status">

          <el-select v-model="searchInfo.status" clearable placeholder="请选择" >
            <el-option label="全部" :value="-1"/>
            <el-option label="未支付" :value="0"/>
            <el-option label="未取单" :value="1"/>
            <el-option label="已完成" :value="4"/>
          </el-select>

        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
          ref="multipleTable"
          style="width: 100%"
          tooltip-effect="dark"
          :data="tableData"
          row-key="id"      >

        <el-table-column  align="left" label="id字段" prop="id" width="120"/>
        <el-table-column align="left" label="拼单号" prop="pindanId" width="120"/>
        <el-table-column align="left" label="订单号" prop="orderNumber" width="120"/>
        <el-table-column align="left" label="用户id" prop="userId" width="120"/>
        <el-table-column align="left" label="商品金额" prop="amount" width="120"/>
        <el-table-column align="left" label="商品Vip金额" prop="amountVip" width="120"/>
        <el-table-column align="left" label="付款金额" prop="amountReal" width="120"/>
 
        <el-table-column align="left" label="支付订单时间" prop="datePay" width="180">
          <template #default="scope">{{ formatDate(scope.row.datePay) }}</template>
        </el-table-column>
        <el-table-column align="left" label="商品总数量" prop="goodsNumber" width="120"/>
       
        <el-table-column align="left" label="取单码" prop="qudama" width="120"/>
        <el-table-column align="left" label="下单ip" prop="ip" width="120"/>
       
        <el-table-column align="left" label="订单状态" prop="status" width="120">
          <template #default="scope">{{ 
            scope.row.status == 0?'未支付':(scope.row.status == 1?'已支付,未取单':'已取单')
            }}</template>
        </el-table-column>
        <el-table-column align="left" label="需要配送" prop="deliveryStatus" width="120">
          <template #default="scope">{{ formatBoolean(scope.row.deliveryStatus) }}</template>
        </el-table-column>
        <el-table-column align="left" label="是否已经支付" prop="isPay" width="120">
          <template #default="scope">{{ formatBoolean(scope.row.isPay) }}</template>
        </el-table-column>
        
        <el-table-column align="left" label="商店id" prop="shopId" width="120"/>
        <el-table-column align="left" label="自提商店id" prop="shopIdZt" width="120"/>
        <el-table-column align="left" label="自提商店名称" prop="shopNameZt" width="120"/>
        <el-table-column align="left" label="配送类型" prop="peisongType" width="120"/>
        
        <el-table-column align="left" label="创建时间" prop="dateAdd" width="180">
          <template #default="scope">{{ formatDate(scope.row.dateAdd) }}</template>
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
  </div>
</template>

<script setup>
import {
  createBeeOrder,
  deleteBeeOrder,
  deleteBeeOrderByIds,
  updateBeeOrder,
  findBeeOrder,
  getBeeOrderList
} from '@/plugin/beeshop/api/beeOrder'

// 全量引入格式化工具 请按需保留
import {
  getDictFunc,
  formatDate,
  formatBoolean,
  filterDict,
  filterDataSource,
  ReturnArrImg,
  onDownloadFile,
  formatEnum
} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {ref, reactive} from 'vue'
import {useRoute, useRouter} from "vue-router";

defineOptions({
  name: 'BeePindanOrder'
})

const route = useRoute()
const router = useRouter()
const beeOrderStatus = ref([])
const beeOrderType = ref([])


// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  id: undefined,
  isDeleted: false,
  dateAdd: new Date(),
  dateUpdate: new Date(),
  dateDelete: undefined,
  amount: 0,
  amountCard: 0,
  amountCoupons: 0,
  amountLogistics: 0,
  amountReal: 0,
  amountRefundTotal: 0,
  amountTax: 0,
  amountTaxGst: 0,
  amountTaxService: 0,
  autoDeliverStatus: undefined,
  dateClose: new Date(),
  datePay: new Date(),
  goodsNumber: undefined,
  hasRefund: false,
  hxNumber: '',
  ip: '',
  isCanHx: false,
  isDelUser: false,
  isEnd: false,
  isHasBenefit: false,
  isNeedLogistics: false,
  peisongType: '',
  isPay: false,
  isScoreOrder: false,
  isSuccessPingtuan: false,
  orderNumber: '',
  orderType: undefined,
  pid: undefined,
  qudanhao: '',
  refundStatus: undefined,
  remark: '',
  score: undefined,
  scoreDeduction: undefined,
  shopId: undefined,
  shopIdZt: undefined,
  shopNameZt: '',
  status: undefined,
  trips: undefined,
  type: undefined,
  uid: undefined,
})


// 验证规则
const rule = reactive({})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
  dateAdd: [{
    validator: (rule, value, callback) => {
      if (searchInfo.value.startDateAdd && !searchInfo.value.endDateAdd) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startDateAdd && searchInfo.value.endDateAdd) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startDateAdd && searchInfo.value.endDateAdd && (searchInfo.value.startDateAdd.getTime() === searchInfo.value.endDateAdd.getTime() || searchInfo.value.startDateAdd.getTime() > searchInfo.value.endDateAdd.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change'
  }],
  dateUpdate: [{
    validator: (rule, value, callback) => {
      if (searchInfo.value.startDateUpdate && !searchInfo.value.endDateUpdate) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startDateUpdate && searchInfo.value.endDateUpdate) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startDateUpdate && searchInfo.value.endDateUpdate && (searchInfo.value.startDateUpdate.getTime() === searchInfo.value.endDateUpdate.getTime() || searchInfo.value.startDateUpdate.getTime() > searchInfo.value.endDateUpdate.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change'
  }],
  dateDelete: [{
    validator: (rule, value, callback) => {
      if (searchInfo.value.startDateDelete && !searchInfo.value.endDateDelete) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startDateDelete && searchInfo.value.endDateDelete) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startDateDelete && searchInfo.value.endDateDelete && (searchInfo.value.startDateDelete.getTime() === searchInfo.value.endDateDelete.getTime() || searchInfo.value.startDateDelete.getTime() > searchInfo.value.endDateDelete.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change'
  }],
  dateClose: [{
    validator: (rule, value, callback) => {
      if (searchInfo.value.startDateClose && !searchInfo.value.endDateClose) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startDateClose && searchInfo.value.endDateClose) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startDateClose && searchInfo.value.endDateClose && (searchInfo.value.startDateClose.getTime() === searchInfo.value.endDateClose.getTime() || searchInfo.value.startDateClose.getTime() > searchInfo.value.endDateClose.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change'
  }],
  datePay: [{
    validator: (rule, value, callback) => {
      if (searchInfo.value.startDatePay && !searchInfo.value.endDatePay) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startDatePay && searchInfo.value.endDatePay) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startDatePay && searchInfo.value.endDatePay && (searchInfo.value.startDatePay.getTime() === searchInfo.value.endDatePay.getTime() || searchInfo.value.startDatePay.getTime() > searchInfo.value.endDatePay.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change'
  }],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({
  sort: 'id',
  order: 'descending',
  status: -1
})
// 排序
const sortChange = ({prop, order}) => {
  const sortMap = {
    id: 'id',
    dateAdd: 'date_add',
  }

  let sort = sortMap[prop]
  if (!sort) {
    sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}

// 重置
const onReset = () => {
  searchInfo.value = {
    sort: 'id',
    order: 'descending',
  }
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    if (searchInfo.value.isDeleted === "") {
      searchInfo.value.isDeleted = null
    }
    if (searchInfo.value.hasRefund === "") {
      searchInfo.value.hasRefund = null
    }
    if (searchInfo.value.isCanHx === "") {
      searchInfo.value.isCanHx = null
    }
    if (searchInfo.value.isDelUser === "") {
      searchInfo.value.isDelUser = null
    }
    if (searchInfo.value.isEnd === "") {
      searchInfo.value.isEnd = null
    }
    if (searchInfo.value.isHasBenefit === "") {
      searchInfo.value.isHasBenefit = null
    }
    if (searchInfo.value.isNeedLogistics === "") {
      searchInfo.value.isNeedLogistics = null
    }
    if (searchInfo.value.isPay === "") {
      searchInfo.value.isPay = null
    }
    if (searchInfo.value.isScoreOrder === "") {
      searchInfo.value.isScoreOrder = null
    }
    if (searchInfo.value.isSuccessPingtuan === "") {
      searchInfo.value.isSuccessPingtuan = null
    }
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
  const table = await getBeeOrderList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

const init = async () => {
  if (route.query.status !== undefined) {
    searchInfo.value.status = route.query.status
  }
  if (route.query.order_id !== undefined) {
    searchInfo.value.id = route.query.order_id
  }
  beeOrderStatus.value = await getDictFunc('OrderStatus')
  beeOrderType.value = await getDictFunc('OrderType')
}
init()

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
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
    deleteBeeOrderFunc(row)
  })
}

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
    multipleSelection.value.map(item => {
      ids.push(item.id)
    })
    const res = await deleteBeeOrderByIds({ids})
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateBeeOrderFunc = async (row) => {
  const res = await findBeeOrder({id: row.id})
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteBeeOrderFunc = async (row) => {
  const res = await deleteBeeOrder({id: row.id})
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
    id: undefined,
    isDeleted: false,
    dateAdd: new Date(),
    dateUpdate: new Date(),
    dateDelete: undefined,
    amount: 0,
    amountCard: 0,
    amountCoupons: 0,
    amountLogistics: 0,
    amountReal: 0,
    amountRefundTotal: 0,
    amountTax: 0,
    amountTaxGst: 0,
    amountTaxService: 0,
    autoDeliverStatus: undefined,
    dateClose: new Date(),
    datePay: new Date(),
    goodsNumber: undefined,
    hasRefund: false,
    hxNumber: '',
    ip: '',
    isCanHx: false,
    isDelUser: false,
    isEnd: false,
    isHasBenefit: false,
    isNeedLogistics: false,
    isPay: false,
    isScoreOrder: false,
    isSuccessPingtuan: false,
    orderNumber: '',
    orderType: undefined,
    pid: undefined,
    qudanhao: '',
    refundStatus: undefined,
    remark: '',
    score: undefined,
    scoreDeduction: undefined,
    shopId: undefined,
    shopIdZt: undefined,
    shopNameZt: '',
    peisongType: '',
    status: undefined,
    trips: undefined,
    type: undefined,
    uid: undefined,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createBeeOrder(formData.value)
        break
      case 'update':
        res = await updateBeeOrder(formData.value)
        break
      default:
        res = await createBeeOrder(formData.value)
        break
    }
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

const gotoDetailPage = async(row) => {
  await router.push('beeOrderDetail?id='+row.id)
}
</script>

<style>

</style>
