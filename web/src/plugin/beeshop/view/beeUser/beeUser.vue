<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="id字段" prop="id">
            
             <el-input v-model.number="searchInfo.id" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="用户show_id" prop="showUid">
            
             <el-input v-model.number="searchInfo.showUid" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="会员" prop="vipLevel">
         <el-input v-model.number="searchInfo.vipLevel" placeholder="搜索条件" />

        </el-form-item>

          <!-- 将需要控制显示状态的查询条件添加到此范围内 
            <el-form-item label="是否虚拟人" prop="isVirtual">
            <el-select v-model="searchInfo.isVirtual" clearable placeholder="请选择">
                <el-option
                    key="true"
                    label="是"
                    value="true">
                </el-option>
                <el-option
                    key="false"
                    label="否"
                    value="false">
                </el-option>
            </el-select>
            </el-form-item>-->
        <el-form-item label="状态" prop="status">
            
             <el-input v-model.number="searchInfo.status" placeholder="搜索条件" />

        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="id"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        
        <el-table-column sortable align="left" label="id字段" prop="id" width="120" />
        <!--   <el-table-column align="left" label="用户show_id" prop="showUid" width="120" />-->
        <el-table-column align="left" label="头像" prop="avatarUrl" width="120" >
          <template #default="scope">
            <CustomPic
                pic-type="file"
                style="margin-top:8px"
                :pic-src="scope.row.avatarUrl"
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="会员卡号" prop="cardNumber" width="120" />
         <!--  <el-table-column align="left" label="城市" prop="city" width="120" />-->
         <el-table-column align="left" label="登录时间" prop="dateLogin" width="180">
            <template #default="scope">{{ formatDate(scope.row.dateLogin) }}</template>
         </el-table-column>
  <!--       <el-table-column align="left" label="性别" prop="gender" width="120" />
        <el-table-column align="left" label="注册ip" prop="ipAdd" width="120" />
        <el-table-column align="left" label="登录ip" prop="ipLogin" width="120" />-->
<!--        <el-table-column align="left" label="人脸识别" prop="isFaceCheck" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isFaceCheck) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="身份证识别" prop="isIdcardCheck" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isIdcardCheck) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="管理员" prop="isManager" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isManager) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="销售人员" prop="isSeller" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isSeller) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="已发放注册优惠券" prop="isSendRegisterCoupons" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isSendRegisterCoupons) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="店长" prop="isShopManager" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isShopManager) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="团长" prop="isTeamLeader" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isTeamLeader) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="团员" prop="isTeamMember" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isTeamMember) }}</template>-->
<!--        </el-table-column>-->
<!--        <el-table-column align="left" label="关注用户" prop="isUserAttendant" width="120">-->
<!--            <template #default="scope">{{ formatBoolean(scope.row.isUserAttendant) }}</template>-->
<!--        </el-table-column>
        <el-table-column align="left" label="是否虚拟人" prop="isVirtual" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.isVirtual) }}</template>
        </el-table-column>-->
        <el-table-column align="left" label="名字" prop="nick" width="120" />
        <el-table-column align="left" label="所在省" prop="province" width="120" />
        <el-table-column sortable align="left" label="vip等级" prop="vipLevel" width="120" />
        <el-table-column align="left" label="注册来源" prop="source" width="120" />
        <el-table-column align="left" label="状态" prop="status" width="120" />
        <el-table-column align="left" label="已删除" prop="isDeleted" width="120">
          <template #default="scope">{{ formatBoolean(scope.row.isDeleted) }}</template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" prop="dateAdd" width="180">
          <template #default="scope">{{ formatDate(scope.row.dateAdd) }}</template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" prop="dateUpdate" width="180">
          <template #default="scope">{{ formatDate(scope.row.dateUpdate) }}</template>
        </el-table-column>
        <el-table-column align="left" label="删除时间" prop="dateDelete" width="180">
          <template #default="scope">{{ formatDate(scope.row.dateDelete) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateBeeUserFunc(scope.row)">变更</el-button>
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
                <span class="text-lg">{{type==='create'?'添加':'修改'}}</span>
                <div>
                  <el-button type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="id字段:"  prop="id" >
              <el-input v-model.number="formData.id" :clearable="true" placeholder="请输入id字段" />
            </el-form-item>
            <!--            <el-form-item label="已删除:"  prop="isDeleted" >-->
<!--              <el-switch v-model="formData.isDeleted" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="创建时间:"  prop="dateAdd" >-->
<!--              <el-date-picker v-model="formData.dateAdd" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />-->
<!--            </el-form-item>-->
<!--            <el-form-item label="更新时间:"  prop="dateUpdate" >-->
<!--              <el-date-picker v-model="formData.dateUpdate" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />-->
<!--            </el-form-item>-->
<!--            <el-form-item label="删除时间:"  prop="dateDelete" >-->
<!--              <el-date-picker v-model="formData.dateDelete" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />-->
<!--            </el-form-item>-->
            <el-form-item label="用户show_id:"  prop="showUid" >
              <el-input v-model.number="formData.showUid" :clearable="true" placeholder="请输入用户show_id" />
            </el-form-item>
            <el-form-item label="头像:"  prop="avatarUrl" >
              <SelectImage
                  v-model="formData.avatarUrl"
              />
              <el-input v-model="formData.avatarUrl" :clearable="true"  placeholder="请输入头像" />
            </el-form-item>
            <el-form-item label="会员卡号:"  prop="cardNumber" >
              <el-input v-model="formData.cardNumber" :clearable="true"  placeholder="请输入会员卡号" />
            </el-form-item>
            <el-form-item label="城市:"  prop="city" >
              <el-input v-model="formData.city" :clearable="true"  placeholder="请输入城市" />
            </el-form-item>
            <el-form-item label="登录时间:"  prop="dateLogin" >
              <el-date-picker v-model="formData.dateLogin" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
            </el-form-item>
            <el-form-item label="性别:"  prop="gender" >
              <el-input v-model.number="formData.gender" :clearable="true" placeholder="请输入性别" />
            </el-form-item>
            <el-form-item label="注册ip:"  prop="ipAdd" >
              <el-input v-model="formData.ipAdd" :clearable="true"  placeholder="请输入注册ip" />
            </el-form-item>
            <el-form-item label="登录ip:"  prop="ipLogin" >
              <el-input v-model="formData.ipLogin" :clearable="true"  placeholder="请输入登录ip" />
            </el-form-item>
<!--            <el-form-item label="人脸识别:"  prop="isFaceCheck" >-->
<!--              <el-switch v-model="formData.isFaceCheck" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="身份证识别:"  prop="isIdcardCheck" >-->
<!--              <el-switch v-model="formData.isIdcardCheck" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="管理员:"  prop="isManager" >-->
<!--              <el-switch v-model="formData.isManager" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="销售人员:"  prop="isSeller" >-->
<!--              <el-switch v-model="formData.isSeller" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="已发放注册优惠券:"  prop="isSendRegisterCoupons" >-->
<!--              <el-switch v-model="formData.isSendRegisterCoupons" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="店长:"  prop="isShopManager" >-->
<!--              <el-switch v-model="formData.isShopManager" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="团长:"  prop="isTeamLeader" >-->
<!--              <el-switch v-model="formData.isTeamLeader" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="团员:"  prop="isTeamMember" >-->
<!--              <el-switch v-model="formData.isTeamMember" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="关注用户:"  prop="isUserAttendant" >-->
<!--              <el-switch v-model="formData.isUserAttendant" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>-->
<!--            </el-form-item>-->
        <!--    <el-form-item label="是否虚拟人:"  prop="isVirtual" >
              <el-switch v-model="formData.isVirtual" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
            </el-form-item>-->
            <el-form-item label="名字:"  prop="nick" >
              <el-input v-model="formData.nick" :clearable="true"  placeholder="请输入名字" />
            </el-form-item>
            <el-form-item label="所在省:"  prop="province" >
              <el-input v-model="formData.province" :clearable="true"  placeholder="请输入所在省" />
            </el-form-item>
            <el-form-item label="vip等级:"  prop="vipLevel" >
              <el-input v-model.number="formData.vipLevel" :clearable="true" placeholder="请输入vip等级" />
            </el-form-item>
            <el-form-item label="注册来源:"  prop="source" >
              <el-input v-model.number="formData.source" :clearable="true" placeholder="请输入注册来源" />
            </el-form-item>
            <el-form-item label="状态:"  prop="status" >
              <el-input v-model.number="formData.status" :clearable="true" placeholder="请输入状态" />
            </el-form-item>
          </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBeeUser,
  deleteBeeUser,
  deleteBeeUserByIds,
  updateBeeUser,
  findBeeUser,
  getBeeUserList
} from '@/plugin/beeshop/api/beeUser'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import SelectImage from "@/components/selectImage/selectImage.vue";
import CustomPic from "@/components/customPic/index.vue";

defineOptions({
    name: 'BeeUser'
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        id: undefined,
        isDeleted: false,
        dateAdd: new Date(),
        dateUpdate: new Date(),
        dateDelete: undefined,
        showUid: undefined,
        avatarUrl: '',
        cardNumber: '',
        city: '',
        dateLogin: new Date(),
        gender: undefined,
        ipAdd: '',
        ipLogin: '',
        isFaceCheck: false,
        isIdcardCheck: false,
        isManager: false,
        isSeller: false,
        isSendRegisterCoupons: false,
        isShopManager: false,
        isTeamLeader: false,
        isTeamMember: false,
        isUserAttendant: false,
        isVirtual: false,
        nick: '',
        province: '',
        vipLevel: undefined,
        source: undefined,
        status: undefined,
        })



// 验证规则
const rule = reactive({
})

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
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
})
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
            id: 'id',
  }

  let sort = sortMap[prop]
  if(!sort){
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
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    if (searchInfo.value.isDeleted === ""){
        searchInfo.value.isDeleted=null
    }
    if (searchInfo.value.isFaceCheck === ""){
        searchInfo.value.isFaceCheck=null
    }
    if (searchInfo.value.isIdcardCheck === ""){
        searchInfo.value.isIdcardCheck=null
    }
    if (searchInfo.value.isManager === ""){
        searchInfo.value.isManager=null
    }
    if (searchInfo.value.isSeller === ""){
        searchInfo.value.isSeller=null
    }
    if (searchInfo.value.isSendRegisterCoupons === ""){
        searchInfo.value.isSendRegisterCoupons=null
    }
    if (searchInfo.value.isShopManager === ""){
        searchInfo.value.isShopManager=null
    }
    if (searchInfo.value.isTeamLeader === ""){
        searchInfo.value.isTeamLeader=null
    }
    if (searchInfo.value.isTeamMember === ""){
        searchInfo.value.isTeamMember=null
    }
    if (searchInfo.value.isUserAttendant === ""){
        searchInfo.value.isUserAttendant=null
    }
    if (searchInfo.value.isVirtual === ""){
        searchInfo.value.isVirtual=null
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
const getTableData = async() => {
  const table = await getBeeUserList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
const setOptions = async () =>{
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
            deleteBeeUserFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
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
      const res = await deleteBeeUserByIds({ ids })
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
const updateBeeUserFunc = async(row) => {
    const res = await findBeeUser({ id: row.id })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteBeeUserFunc = async (row) => {
    const res = await deleteBeeUser({ id: row.id })
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
        showUid: undefined,
        avatarUrl: '',
        cardNumber: '',
        city: '',
        dateLogin: new Date(),
        gender: undefined,
        ipAdd: '',
        ipLogin: '',
        isFaceCheck: false,
        isIdcardCheck: false,
        isManager: false,
        isSeller: false,
        isSendRegisterCoupons: false,
        isShopManager: false,
        isTeamLeader: false,
        isTeamMember: false,
        isUserAttendant: false,
        isVirtual: false,
        nick: '',
        province: '',
        vipLevel: undefined,
        source: undefined,
        status: undefined,
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createBeeUser(formData.value)
                  break
                case 'update':
                  res = await updateBeeUser(formData.value)
                  break
                default:
                  res = await createBeeUser(formData.value)
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

</script>

<style>

</style>
