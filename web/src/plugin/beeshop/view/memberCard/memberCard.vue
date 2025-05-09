<template>
    <div>
        <div class="gva-search-box">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
                @keyup.enter="onSubmit">
                <el-form-item label="会员卡名" prop="name">
                    <el-input v-model="searchInfo.name" placeholder="请输入会员卡名称" />
                </el-form-item>

                <el-form-item label="会员卡类型" prop="validMonth">
                    <el-select v-model="searchInfo.validMonth" clearable placeholder="请选择">
                        <el-option label="全部" :value="-1" />
                        <el-option v-for="item, index in memberCardTypes" :key="index" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item label="会员卡状态" prop="deleteFlag">
                    <el-select v-model="searchInfo.deleteFlag" clearable placeholder="请选择">
                        <el-option label="全部" :value="-1" />
                        <el-option label="未删除" :value="0" />
                        <el-option label="已删除" :value="1" />
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
                    <el-button icon="refresh" @click="onReset">重置</el-button>
                </el-form-item>
            </el-form>
        </div>
        <div class="gva-table-box">
            <div class="gva-btn-list">
                <el-button type="primary" icon="plus" @click="addHandler">新增</el-button>
            </div>
            <el-table ref="table" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="id">
                <el-table-column type="index" width="55" label="序号" />
                <el-table-column align="left" label="会员卡名称" prop="name" width="180" />
                <el-table-column align="left" label="会员卡金额(元)" prop="amount" width="120" />
                <el-table-column align="left" label="排序" prop="sortNum" width="120" />
                <el-table-column align="left" label="类型" prop="validMonth" width="120">
                    <template #default="scope">
                        {{ filterMemberCardType(scope.row.validMonth) }}
                    </template>
                </el-table-column>
                <el-table-column align="left" label="添加时间" prop="create_time" width="240">
                    <template #default="scope">
                        {{ formatDate(scope.row.create_time) }}
                    </template>
                </el-table-column>
                <el-table-column align="left" label="状态" prop="deleteFlag" width="120">
                    <template #default="scope">
                        <el-tag v-if="scope.row.deleteFlag == 1" type="warning">已删除 </el-tag>
                        <el-tag v-if="scope.row.deleteFlag == 0" type="primary">未删除 </el-tag>
                    </template>
                </el-table-column>
                <el-table-column align="left" label="删除时间" prop="delete_time" width="180">
                    <template #default="scope">{{ scope.row.deleteFlag == 1 ? formatDate(scope.row.delete_time) : ''
                    }}</template>
                </el-table-column>
                <el-table-column align="left" label="操作" fixed="right">
                    <template #default="scope">
                        <div>
                            <el-button type="primary" link class="table-button"
                                @click="info(scope.row)">查看详情</el-button>

                            <el-button type="primary" link icon="edit" class="table-button" @click="edit(scope.row)">
                                变更
                            </el-button>
                            <el-button v-if="scope.row.deleteFlag == 0" type="primary" link icon="delete"
                                @click="deleteRow(scope.row)">删除</el-button>
                                <el-button v-if="scope.row.deleteFlag == 1" type="primary" link icon="recovery"
                                @click="recover(scope.row)">恢复</el-button>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <div class="gva-pagination">
                <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page"
                    :page-size="pageSize" :page-sizes="[10, 30, 50, 100]" :total="total"
                    @current-change="handleCurrentChange" @size-change="handleSizeChange" />
            </div>
        </div>
        <MemberCardFormComponent :id="selectedId" v-model:open="openForm" :view="view" @refresh="refreshData" />
    </div>
</template>

<script setup>
import { formatDate } from '@/utils/format'
import { ref, onMounted } from 'vue'
import { getDictDetailList } from '@/api/sysDictionaryDetail'
import { memberCardPage, deleteOneById ,recoverOneById } from '@/plugin/beeshop/api/beeMemberCard'
import MemberCardFormComponent from './form.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
defineOptions({
    name: 'MemberCard'
})



// 自动化生成的字典（可能为空）以及字段
const searchInfo = ref({
    name: '',
    validMonth: -1,
    deleteFlag: -1,
})
const elSearchFormRef = ref()
const openForm = ref(false)
const selectedId = ref(0)
// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const memberCardTypes = ref([]);
const view = ref(true)
onMounted(() => {
    getDictDetailList('member_card_type').then(res => {
        if (res.code == 0) {
            memberCardTypes.value = res.data
        }
    })
    getTableData();
})
const filterMemberCardType = (validMonth) => {
    const types = memberCardTypes.value.filter(item => item.value == validMonth)
    if (types && types.length > 0) {
        return types[0].label
    }
}
const addHandler = () => {
    openForm.value = true
    view.value = false
    selectedId.value = 0
}
// 重置
const onReset = () => {
    searchInfo.value = {
        name: '',
        validMonth: -1,
        deleteFlag: -1,
    }
    getTableData()
}

// 搜索
const onSubmit = () => {
    page.value = 1,
        pageSize.value = 10
    getTableData();
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
    const table = await memberCardPage({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    console.log(table)
    if (table.code === 0) {
        tableData.value = table.data.list
        total.value = table.data.total
    }
}

const refreshData = () => {
    getTableData();
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        deleteOneById(row.id).then(res => {
            if (res.code == 0) {
                getTableData()
                ElMessage.success("删除成功");
            } else {
                ElMessage.error("删除失败");
            }
        })
    })
}


const recover = (row) => {
    ElMessageBox.confirm('确定要恢复吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        recoverOneById(row.id).then(res => {
            if (res.code == 0) {
                getTableData()
                ElMessage.success("恢复成功");
            } else {
                ElMessage.error("恢复失败");
            }
        })
    })
}


const edit = (row) => {
    openForm.value = true
    view.value = false
    selectedId.value = row.id
}

const info = (row) => {
    openForm.value = true
    view.value = true
    selectedId.value = row.id
}

</script>

<style></style>