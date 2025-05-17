<template>
    <div>
        <div class="gva-search-box">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
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
                <el-form-item>
                    <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
                    <el-button icon="refresh" @click="onReset">重置</el-button>
                </el-form-item>
            </el-form>
        </div>
        <div class="gva-table-box">
            <el-table ref="table" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="id">
                <el-table-column type="index" width="55" label="序号"/>
                <el-table-column align="left" label="会员卡名称" prop="name"  />
                <el-table-column align="left" label="会员卡金额(元)" prop="amount" width="120" />
                <el-table-column align="left" label="类型" prop="valid_month" width="120">
                    <template #default="scope">
                        {{ filterMemberCardType(scope.row.valid_month) }}
                    </template>
                    </el-table-column>
                <el-table-column align="left" label="添加时间" prop="create_time" :formatter="formatTime"/>
 
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
import { dayjs } from "element-plus";

import {getDictDetailList} from '@/api/sysDictionaryDetail'
import { userMemberCardPage } from '@/plugin/beeshop/api/beeUserMemberCard'

defineOptions({
    name: 'UseremberCard'
})



// 自动化生成的字典（可能为空）以及字段
const searchInfo = ref({
    name: '',
    validMonth: -1,
    deleteFlag: -1,
})
const elSearchFormRef = ref()
// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const memberCardTypes = ref([]);
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
const formatTime = (v)=>{
return dayjs(v.create_time).format('YYYY-MM-DD HH:mm:ss')
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
    const table = await userMemberCardPage({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    console.log(table)
    if (table.code === 0) {
        tableData.value = table.data.list
        total.value = table.data.total
    }
}


</script>

<style></style>