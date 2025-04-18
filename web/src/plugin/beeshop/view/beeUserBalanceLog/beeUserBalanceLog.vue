<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="用户id" prop="uid">
          <el-input v-model.number="searchInfo.uid" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="门店" prop="shopId">
          <el-select
            v-model="searchInfo.shopId"
            filterable
            placeholder="请选择门店"
            style="width: 240px"
          >
            <el-option
              v-for="item in shops"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
          <el-form-item label="添加时间" prop="createTime">
            <template #label>
              <span>
                添加时间
                <el-tooltip
                  content="搜索范围是开始日期（包含）至结束日期（不包含）"
                >
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </span>
            </template>
            <el-date-picker
              v-model="searchInfo.startDateAdd"
              type="datetime"
              placeholder="开始日期"
              :disabled-date="
                (time) =>
                  searchInfo.endDateAdd
                    ? time.getTime() > searchInfo.endDateAdd.getTime()
                    : false
              "
            ></el-date-picker>
            —
            <el-date-picker
              v-model="searchInfo.endDateAdd"
              type="datetime"
              placeholder="结束日期"
              :disabled-date="
                (time) =>
                  searchInfo.startDateAdd
                    ? time.getTime() < searchInfo.startDateAdd.getTime()
                    : false
              "
            ></el-date-picker>
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >查询</el-button
          >
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button
            link
            type="primary"
            icon="arrow-down"
            @click="showAllQuery = true"
            v-if="!showAllQuery"
            >展开</el-button
          >
          <el-button
            link
            type="primary"
            icon="arrow-up"
            @click="showAllQuery = false"
            v-else
            >收起</el-button
          >
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog"
          >新增</el-button
        >
        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >删除</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column
          align="left"
          sortable
          label="id字段"
          prop="id"
          width="80"
        />
        <el-table-column
          align="left"
          label="门店"
          prop="shopName"
          width="120"
        />
        <el-table-column
          align="left"
          label="订单id"
          prop="orderId"
          width="120"
        />
        <el-table-column
          align="left"
          label="订单金额"
          prop="amount"
          width="120"
        />
        <el-table-column
          align="left"
          label="货币类型"
          prop="balanceType"
          width="120"
        >
          <template #default="scope">{{
            formatEnum(scope.row.balanceType, balanceTypeMap)
          }}</template>
        </el-table-column>
        <el-table-column align="left" label="数量" prop="num" width="120" />
        <el-table-column align="left" label="用户id" prop="uid" width="120" />
        <el-table-column align="left" label="备注" prop="mark" width="120" />
        <el-table-column
          align="left"
          label="已删除"
          prop="isDeleted"
          width="120"
        >
          <template #default="scope">{{
            formatBoolean(scope.row.isDeleted)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="创建时间"
          prop="dateAdd"
          width="180"
        >
          <template #default="scope">{{
            formatDate(scope.row.dateAdd)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="更新时间"
          prop="dateUpdate"
          width="180"
        >
          <template #default="scope">{{
            formatDate(scope.row.dateUpdate)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="删除时间"
          prop="dateDelete"
          width="180"
        >
          <template #default="scope">{{
            formatDate(scope.row.dateDelete)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateBeeUserBalanceLogFunc(scope.row)"
              >变更</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
              >删除</el-button
            >
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
    <el-drawer
      destroy-on-close
      size="800"
      v-model="dialogFormVisible"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === "create" ? "添加" : "修改" }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form
        :model="formData"
        label-position="top"
        ref="elFormRef"
        :rules="rule"
        label-width="80px"
      >
        <el-form-item label="id字段:" prop="id">
          <el-input
            v-model.number="formData.id"
            :clearable="true"
            placeholder="请输入id字段"
          />
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
        <el-form-item label="订单id:" prop="orderId">
          <el-input
            v-model="formData.orderId"
            :clearable="true"
            placeholder="请输入订单id"
          />
        </el-form-item>
        <el-form-item label="货币类型:" prop="balanceType">
          <el-select
            v-model="formData.balanceType"
            clearable
            placeholder="请选择"
            :clearable="false"
          >
            <el-option
              v-for="(item, key) in balanceTypeMap"
              :key="key"
              :label="item.label"
              :value="parseInt(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="数量:" prop="num">
          <el-input-number
            v-model="formData.num"
            style="width: 100%"
            :precision="2"
            :clearable="true"
          />
        </el-form-item>
        <el-form-item label="用户id:" prop="uid">
          <el-input
            v-model.number="formData.uid"
            :clearable="true"
            placeholder="请输入用户id"
          />
        </el-form-item>
        <el-form-item label="备注:" prop="mark">
          <el-input
            v-model="formData.mark"
            :clearable="true"
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBeeUserBalanceLog,
  deleteBeeUserBalanceLog,
  deleteBeeUserBalanceLogByIds,
  updateBeeUserBalanceLog,
  findBeeUserBalanceLog,
  getBeeUserBalanceLogList,
} from "@/plugin/beeshop/api/beeUserBalanceLog";
import {
  getAllMyBeeShopInfos
} from '@/plugin/beeshop/api/beeShopInfo'
// 全量引入格式化工具 请按需保留
import {
  getDictFunc,
  formatDate,
  formatBoolean,
  filterDict,
  filterDataSource,
  ReturnArrImg,
  onDownloadFile,
  formatEnum,
} from "@/utils/format";
import { ElMessage, ElMessageBox } from "element-plus";
import { ref, reactive } from "vue";

defineOptions({
  name: "BeeUserBalanceLog",
});

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false);
const shops = ref([])
const balanceTypeMap = ref([]);
const init = async () => {
  balanceTypeMap.value = await getDictFunc("BalanceType");
  console.log(balanceTypeMap);
  shops.value  = (await getAllMyBeeShopInfos()).data.list
  
};
init();

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  id: undefined,
  isDeleted: false,
  dateAdd: new Date(),
  dateUpdate: new Date(),
  dateDelete: undefined,
  orderId: "",
  balanceType: undefined,
  num: 0,
  uid: undefined,
  mark: "",
});

// 验证规则
const rule = reactive({});

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startDateAdd && !searchInfo.value.endDateAdd) {
          callback(new Error("请填写结束日期"));
        } else if (
          !searchInfo.value.startDateAdd &&
          searchInfo.value.endDateAdd
        ) {
          callback(new Error("请填写开始日期"));
        } else if (
          searchInfo.value.startDateAdd &&
          searchInfo.value.endDateAdd &&
          (searchInfo.value.startDateAdd.getTime() ===
            searchInfo.value.endDateAdd.getTime() ||
            searchInfo.value.startDateAdd.getTime() >
              searchInfo.value.endDateAdd.getTime())
        ) {
          callback(new Error("开始日期应当早于结束日期"));
        } else {
          callback();
        }
      },
      trigger: "change",
    },
  ],
});

const elFormRef = ref();
const elSearchFormRef = ref();

// =========== 表格控制部分 ===========
const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const tableData = ref([]);
const searchInfo = ref({
  sort: "id",
  order: "descending",
  shopId: ''
});

// 重置
const onReset = () => {
  searchInfo.value = {
    sort: "id",
    order: "descending",
  };
  getTableData();
};

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return;
    page.value = 1;
    pageSize.value = 10;
    if (searchInfo.value.isDeleted === "") {
      searchInfo.value.isDeleted = null;
    }
    getTableData();
  });
};

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
};

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
};

// 查询
const getTableData = async () => {
  const table = await getBeeUserBalanceLogList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  });
  if (table.code === 0) {
    tableData.value = table.data.list;
    total.value = table.data.total;
    page.value = table.data.page;
    pageSize.value = table.data.pageSize;
  }
};

getTableData();

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {};

// 获取需要的字典 可能为空 按需保留
setOptions();

// 多选数据
const multipleSelection = ref([]);
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val;
};

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm("确定要删除吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    deleteBeeUserBalanceLogFunc(row);
  });
};

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm("确定要删除吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    const ids = [];
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: "warning",
        message: "请选择要删除的数据",
      });
      return;
    }
    multipleSelection.value &&
      multipleSelection.value.map((item) => {
        ids.push(item.id);
      });
    const res = await deleteBeeUserBalanceLogByIds({ ids });
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: "删除成功",
      });
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--;
      }
      getTableData();
    }
  });
};

// 行为控制标记（弹窗内部需要增还是改）
const type = ref("");

// 更新行
const updateBeeUserBalanceLogFunc = async (row) => {
  const res = await findBeeUserBalanceLog({ id: row.id });
  type.value = "update";
  if (res.code === 0) {
    formData.value = res.data;
    dialogFormVisible.value = true;
  }
};

// 删除行
const deleteBeeUserBalanceLogFunc = async (row) => {
  const res = await deleteBeeUserBalanceLog({ id: row.id });
  if (res.code === 0) {
    ElMessage({
      type: "success",
      message: "删除成功",
    });
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--;
    }
    getTableData();
  }
};

// 弹窗控制标记
const dialogFormVisible = ref(false);

// 打开弹窗
const openDialog = () => {
  type.value = "create";
  dialogFormVisible.value = true;
};

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false;
  formData.value = {
    id: undefined,
    isDeleted: false,
    dateAdd: new Date(),
    dateUpdate: new Date(),
    dateDelete: undefined,
    orderId: "",
    balanceType: undefined,
    num: 0,
    uid: undefined,
    mark: "",
  };
};
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return;
    let res;
    switch (type.value) {
      case "create":
        res = await createBeeUserBalanceLog(formData.value);
        break;
      case "update":
        res = await updateBeeUserBalanceLog(formData.value);
        break;
      default:
        res = await createBeeUserBalanceLog(formData.value);
        break;
    }
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: "创建/更改成功",
      });
      closeDialog();
      getTableData();
    }
  });
};
</script>

<style></style>
