<template>
    <div>
        <div class="gva-search-box">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
                <el-form-item label="会员卡名" prop="name">
                    <el-input v-model="searchInfo.name" placeholder="请输入会员卡名称" />
                </el-form-item>

                <el-form-item label="会员卡类型" prop="validMonth">
                    <el-select v-model="searchInfo.validMonth" clearable placeholder="请选择">
                        
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
                <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            </div>
            <el-table ref="table" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="id">
                <el-table-column type="index" width="55" label="序号"/>
                <el-table-column align="left" label="会员卡名称" prop="name" width="120" />
                <el-table-column align="left" label="会员卡金额(元)" prop="amount" width="120" />
                <el-table-column align="left" label="排序" prop="sortNum" width="120" />
                <el-table-column align="left" label="类型" prop="validMonth" width="120" />
                <el-table-column align="left" label="添加时间" prop="createTIme" width="120" />
                <el-table-column align="left" label="状态" prop="deleteFlag" width="120" />
                <el-table-column align="left" label="删除时间" prop="dateDelete" width="180">
                    <template #default="scope">{{ formatDate(scope.row.dateDelete) }}</template>
                </el-table-column>
                <el-table-column align="left" label="操作" fixed="right" min-width="240">
                    <template #default="scope">
                        <div>
                            <el-button type="primary" link class="table-button"
                                @click="gotoDetailPage(scope.row)">查看详情</el-button>
                        </div>
                        <div>
                            <el-button type="primary" link icon="edit" class="table-button"
                                @click="updateBeeOrderFunc(scope.row)">
                                变更
                            </el-button>
                            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
        
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {getDictDetailList} from '@/api/sysDictionaryDetail'
defineOptions({
    name: 'MemberCard'
})



// 自动化生成的字典（可能为空）以及字段
const searchInfo = ref({
    name: '',
    validMonth: 0,
    deleteFlag:-1,
})
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

onMounted(()=>{
    getDictDetailList().then(res=>{
        console.log(res,'会员卡类型')
    })
})
// 重置
const onReset = () => {
    searchInfo.value = {
        name: '',
    validMonth: 0,
    deleteFlag:-1,
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
    const table = await getBeeOrderList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
        tableData.value = table.data.list
        total.value = table.data.total
        page.value = table.data.page
        pageSize.value = table.data.pageSize
    }
}

// 删除行
// const deleteRow = (row) => {
//     ElMessageBox.confirm('确定要删除吗?', '提示', {
//         confirmButtonText: '确定',
//         cancelButtonText: '取消',
//         type: 'warning'
//     }).then(() => {
//         deleteBeeOrderFunc(row)
//     })
// }

</script>

<style></style>