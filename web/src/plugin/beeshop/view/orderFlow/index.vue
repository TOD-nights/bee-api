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
        <el-form-item label="创建时间" prop="dateAdd">
          <template #label>
            <span>
              创建时间
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

        <el-form-item label="订单号" prop="orderNo">
          <el-input v-model="searchInfo.orderNo" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="第三方订单号" prop="thirdOrderNo">
          <el-input v-model="searchInfo.thirdOrderNo" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="支付状态" prop="status">
          <el-select
            v-model="searchInfo.status"
            clearable
            placeholder="请选择"
            @clear="
              () => {
                searchInfo.status = undefined;
              }
            "
          >
            <el-option
              v-for="(item, key) in beePayLogStatus"
              :key="key"
              :label="item.label"
              :value="parseInt(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="用户id" prop="uid">
          <el-input v-model.number="searchInfo.uid" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="门店" prop="shopid">
          <el-select
            v-model="searchInfo.shopid"
            clearable
            placeholder="请选择门店"
            @clear="
              () => {
                searchInfo.shopid = undefined;
              }
            "
          >
            <el-option
              v-for="(item, index) in shops"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品名" prop="goodsName">
          <el-input v-model="searchInfo.goodsName" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="支付状态" prop="status">
          <el-select
            v-model="searchInfo.status"
            clearable
            placeholder="请选择"
            @clear="
              () => {
                searchInfo.status = undefined;
              }
            "
          >
            <el-option
              v-for="(item, key) in beePayLogStatus"
              :key="key"
              :label="item.label"
              :value="parseInt(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="用户id" prop="uid">
          <el-input v-model.number="searchInfo.uid" placeholder="搜索条件" />
        </el-form-item>
        <!-- <el-form-item label="门店" prop="shopid">
          <el-select
            v-model="searchInfo.shopid"
            filterable
            clearable
            placeholder="请输入门店名称"
            @clear="
              () => {
                searchInfo.shopid = undefined;
              }
            "
          >
            <el-option
              v-for="(item, index) in shops"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item> -->







        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
          <el-form-item label="总金额" prop="money">
            <el-input
              v-model.number="searchInfo.startMoney"
              placeholder="最小值"
            />
            —
            <el-input
              v-model.number="searchInfo.endMoney"
              placeholder="最大值"
            />
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
        @sort-change="sortChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column
          sortable
          align="left"
          label="id字段"
          prop="id"
          width="120"
        />
        <el-table-column
          sortable
          align="left"
          label="用户id"
          prop="uid"
          width="120"
        />
        <el-table-column
          sortable
          align="left"
          label="总金额"
          prop="money"
          width="120"
        />
        <el-table-column
          align="left"
          label="门店"
          prop="shopName"
          width="120"
        />
        <el-table-column
          align="left"
          label="商品"
          prop="goodsName"
          width="120"
        />
        <el-table-column
          align="left"
          label="商品数量"
          prop="goods_number"
          width="120"
        />
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
          label="next_action"
          prop="nextAction"
          width="120"
        />
       
        <el-table-column
          align="left"
          label="订单号"
          prop="orderNo"
          width="120"
        />
        <el-table-column
          align="left"
          label="第三方订单号"
          prop="thirdOrderNo"
          width="120"
        />
        <el-table-column
          align="left"
          label="支付网关"
          prop="payGate"
          width="120"
        />
        <el-table-column align="left" label="备注" prop="remark" width="120" />
        <el-table-column
          sortable
          align="left"
          label="支付状态"
          prop="status"
          width="120"
        >
          <template #default="scope">
            {{ formatEnum(scope.row.status, beePayLogStatus) }}
          </template>
        </el-table-column>
       
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
              @click="updateBeePayLogFunc(scope.row)"
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
        <el-form-item label="总金额:" prop="money">
          <el-input-number
            v-model="formData.money"
            style="width: 100%"
            :precision="2"
            :clearable="true"
          />
        </el-form-item>
        <el-form-item label="next_action:" prop="nextAction">
          <el-input
            v-model="formData.nextAction"
            :clearable="true"
            placeholder="请输入next_action"
          />
        </el-form-item>
        <el-form-item label="订单号:" prop="orderNo">
          <el-input
            v-model="formData.orderNo"
            :clearable="true"
            placeholder="请输入订单号"
          />
        </el-form-item>
        <el-form-item label="第三方订单号:" prop="thirdOrderNo">
          <el-input
            v-model="formData.thirdOrderNo"
            :clearable="true"
            placeholder="请输入第三方订单号"
          />
        </el-form-item>
        <el-form-item label="支付网关:" prop="payGate">
          <el-input
            v-model="formData.payGate"
            :clearable="true"
            placeholder="请输入支付网关"
          />
        </el-form-item>
        <el-form-item label="备注:" prop="remark">
          <el-input
            v-model="formData.remark"
            :clearable="true"
            placeholder="请输入remark"
          />
        </el-form-item>
        <el-form-item label="支付状态:" prop="status">
          <el-select
            v-model="formData.status"
            placeholder="请选择支付状态"
            style="width: 100%"
            :clearable="true"
          >
            <el-option
              v-for="(item, key) in beePayLogStatus"
              :key="key"
              :label="item.label"
              :value="parseInt(item.value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="用户id:" prop="uid">
          <el-input
            v-model.number="formData.uid"
            :clearable="true"
            placeholder="请输入用户id"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBeePayLog,
  deleteBeePayLog,
  deleteBeePayLogByIds,
  updateBeePayLog,
  findBeePayLog,
  getBeePayLogList,
} from "@/plugin/beeshop/api/beePayLog";

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
import { getAllMyBeeShopInfos } from "@/plugin/beeshop/api/beeShopInfo";

import { ElMessage, ElMessageBox } from "element-plus";
import { ref, reactive,onMounted } from "vue";

defineOptions({
  name: "BeePayLog",
});

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false);

// 自动化生成的字典（可能为空）以及字段
const beePayLogStatus = ref([]);
const formData = ref({
  id: undefined,
  isDeleted: false,
  dateAdd: new Date(),
  dateUpdate: new Date(),
  dateDelete: undefined,
  money: 0,
  nextAction: "",
  orderNo: "",
  thirdOrderNo: "",
  payGate: "",
  remark: "",
  status: "",
  uid: undefined,
});

// 验证规则
const rule = reactive({});

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error("请填写结束日期"));
        } else if (
          !searchInfo.value.startCreatedAt &&
          searchInfo.value.endCreatedAt
        ) {
          callback(new Error("请填写开始日期"));
        } else if (
          searchInfo.value.startCreatedAt &&
          searchInfo.value.endCreatedAt &&
          (searchInfo.value.startCreatedAt.getTime() ===
            searchInfo.value.endCreatedAt.getTime() ||
            searchInfo.value.startCreatedAt.getTime() >
              searchInfo.value.endCreatedAt.getTime())
        ) {
          callback(new Error("开始日期应当早于结束日期"));
        } else {
          callback();
        }
      },
      trigger: "change",
    },
  ],
  dateAdd: [
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
  dateUpdate: [
    {
      validator: (rule, value, callback) => {
        if (
          searchInfo.value.startDateUpdate &&
          !searchInfo.value.endDateUpdate
        ) {
          callback(new Error("请填写结束日期"));
        } else if (
          !searchInfo.value.startDateUpdate &&
          searchInfo.value.endDateUpdate
        ) {
          callback(new Error("请填写开始日期"));
        } else if (
          searchInfo.value.startDateUpdate &&
          searchInfo.value.endDateUpdate &&
          (searchInfo.value.startDateUpdate.getTime() ===
            searchInfo.value.endDateUpdate.getTime() ||
            searchInfo.value.startDateUpdate.getTime() >
              searchInfo.value.endDateUpdate.getTime())
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
const shops = ref([])
// =========== 表格控制部分 ===========
const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const tableData = ref([]);
const searchInfo = ref({
  sort: "id",
  order: "descending",
});

onMounted(()=>{
  getAllMyBeeShopInfos().then((res) => {
    if (res.code == 0) {
      shops.value = res.data.list;
    }
  });
});
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
    id: "id",
    dateAdd: "date_add",
    dateUpdate: "date_update",
    money: "money",
    status: "status",
    uid: "uid",
  };

  let sort = sortMap[prop];
  if (!sort) {
    sort = prop.replace(/[A-Z]/g, (match) => `_${match.toLowerCase()}`);
  }

  searchInfo.value.sort = sort;
  searchInfo.value.order = order;
  getTableData();
};

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
  const table = await getBeePayLogList({
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
const setOptions = async () => {
  beePayLogStatus.value = await getDictFunc("BeePayLogStatus");
};

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
    deleteBeePayLogFunc(row);
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
    const res = await deleteBeePayLogByIds({ ids });
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
const updateBeePayLogFunc = async (row) => {
  const res = await findBeePayLog({ id: row.id });
  type.value = "update";
  if (res.code === 0) {
    formData.value = res.data;
    dialogFormVisible.value = true;
  }
};

// 删除行
const deleteBeePayLogFunc = async (row) => {
  const res = await deleteBeePayLog({ id: row.id });
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
    money: 0,
    nextAction: "",
    orderNo: "",
    thirdOrderNo: "",
    payGate: "",
    remark: "",
    status: "",
    uid: undefined,
  };
};
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return;
    let res;
    switch (type.value) {
      case "create":
        res = await createBeePayLog(formData.value);
        break;
      case "update":
        res = await updateBeePayLog(formData.value);
        break;
      default:
        res = await createBeePayLog(formData.value);
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
