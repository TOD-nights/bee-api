<template>

  
  <div
    class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 py-2 gap-4 md:gap-2 gva-container2"
  >
    
  <gva-card custom-class="col-span-1 lg:col-span-6" title="消费流水">
  <!-- 修改搜索栏布局 -->
  <el-row :gutter="10" class="mb-4">
    <!-- 在手机上占满宽度，在大屏幕上占4格 -->
    <el-col :xs="24" :sm="24" :md="8" :lg="6" class="mb-2">
      <el-select 
        v-model="shopId" 
        @change="shopIdChangeHandler"
        placeholder="请选择门店"
        style="width: 100%"
      >
        <el-option
          v-for="(item, index) in shops"
          :key="index"
          :value="item.id"
          :label="item.name"
        >{{ item.name }}</el-option>
      </el-select>
    </el-col>
    
    <!-- 在手机上占满宽度，在大屏幕上占8格 -->
    <el-col :xs="24" :sm="24" :md="12" :lg="10" class="mb-2">
      <el-date-picker
        v-model="dates"
        type="daterange"
        style="width: 100%"
        range-separator="到"
        start-placeholder="开始"
        end-placeholder="结束"
      />
    </el-col>
    
    <!-- 搜索按钮 -->
    <el-col :xs="24" :sm="24" :md="4" :lg="2" class="mb-2">
      <el-button 
        @click="refreshOrders"
        style="width: 100%"
      >
        <el-icon><Search /></el-icon>
        <span>搜索</span>
      </el-button>
    </el-col>
  </el-row>

  <!-- 表格部分 -->
  <el-table
    :data="orders"
    stripe
    style="width: 100%;height: 300px"
    :span-method="spanMethodHandler"
    :default-sort="{ prop: 'dateAdd', order: 'descending' }"
  >
    <!-- 在手机上隐藏订单号列 -->
    <el-table-column 
      prop="orderNumber" 
      label="订单号" 
      width="170" 
      :visible="false"
      class="hidden sm:table-cell"
    />
    <el-table-column
      prop="amountReal"
      label="订单金额"
      align="center"
      show-overflow-tooltip
    />
    <el-table-column 
      prop="dateAdd" 
      label="支付时间"
      min-width="120"
      sortable
    >
      <template #default="scope">
        {{ dateFmt(scope.row.dateAdd) }}
      </template>
    </el-table-column>
  </el-table>
</gva-card>


<!-- 添加今日流水卡片 -->
 <!-- 全部商店-->
 <el-card 
 v-if="isAdmin" 
 @click="gotoPage('/bee_index/beeFinancialManager/beePayLog')">
  <el-statistic 
    :value="todayAmount" 
    :precision="2"
  >
    <template #title>
      <div class="font-bold">今日全部商店消费金额</div>
    </template>
    <template #suffix>
      <span class="text-sm">元</span>
    </template>
  </el-statistic>
</el-card>

<!-- 单个商店-->
    <el-card @click="gotoPage('/bee_index/beeFinancialManager/beePayLog')">
  <el-statistic 
    :value="todayAmountSelect" 
    :precision="2"
  >
    <template #title>
      <div class="font-bold">今日选中商店消费金额</div>
    </template>
    <template #suffix>
      <span class="text-sm">元</span>
    </template>
  </el-statistic>
</el-card>


    
    <!-- 添加今日订单数卡片 -->


    <el-card 
    v-if="isAdmin" 
    @click="gotoPage('/bee_index/shop-order-admin/beeOrder')">
  <el-statistic 
    :value="todayOrderCount" 
  >
    <template #title>
      <div class="font-bold">今日全部商店订单数</div>
    </template>
    <template #suffix>
      <span class="text-sm">张</span>
    </template>
  </el-statistic>
</el-card>


<!-- 添加选中商店今日订单数卡片 -->

<el-card @click="gotoPage('/bee_index/shop-order-admin/beeOrder')">
  <el-statistic 
    :value="todayOrderCountSelect" 
  >
    <template #title>
      <div class="font-bold">今日选择商店订单数</div>
    </template>
    <template #suffix>
      <span class="text-sm">张</span>
    </template>
  </el-statistic>
</el-card>



<!-- 添加今日充值卡片 -->



<el-card 
v-if="isAdmin" 
@click="gotoPage('/bee_index/beeFinancialManager/beePayLog')">
  <el-statistic 
    :value="todayRechargeTotal" 
    :precision="2"
  >
    <template #title>
      <div class="font-bold">今日全部商店充值金额</div>
    </template>
    <template #suffix>
      <span class="text-sm">元</span>
    </template>
  </el-statistic>
</el-card>



<!-- 添加选中商店今日充值卡片-->



<el-card @click="gotoPage('/bee_index/beeFinancialManager/beePayLog')">
  <el-statistic 
    :value="todayRechargeSelect" 
    :precision="2"
  >
    <template #title>
      <div class="font-bold">今日选择商店充值金额</div>
    </template>
    <template #suffix>
      <span class="text-sm">元</span>
    </template>
  </el-statistic>
</el-card>






  </div>
</template>

<script setup>

import { reactive, ref, watch } from "vue";
import { GvaChart, GvaCard } from "./componenst";
import { dayjs } from "element-plus";
import { getBeeUserList } from "@/plugin/beeshop/api/beeUser";
import {
  getBeePayLogList,
  getBeePayTotal,
} from "@/plugin/beeshop/api/beePayLog";
import { getBeeUserBalanceLogList,getBeeUserBalanceLogCount } from "@/plugin/beeshop/api/beeUserBalanceLog";
import { getBeeOrderList, orderList,orderStatistic } from "@/plugin/beeshop/api/beeOrder";
import { getBeeShopGoodsList } from "@/plugin/beeshop/api/beeShopGoods";
import { formatDate, formatEnum, getDictFunc } from "@/utils/format";
import { useRoute, useRouter } from "vue-router";
import { getAllMyBeeShopInfos } from "@/plugin/beeshop/api/beeShopInfo";

//判断是否为管理员
import { computed } from 'vue'
import { useUserStore } from '@/pinia/modules/user'

const route = useRoute();
const router = useRouter();

const beeOrderStatus = ref([]);
const userData = ref([]);
const goodsData = ref([]);
const orderData = ref([]);
const payData = ref([]);
const payNumData = ref([]);
const orderTodoData = ref([]);
const orderTodo = ref({});
const shops = ref([]);
const orders = ref([]);
const totalOrders = ref(0);
const orderSum = ref(0);
const pageNum = ref(0);
const pageSize = ref(10);
const dates = ref([]);
const shopId = ref("");

// 添加新的变量来分别存储充值和支付数据
const todayRechargeTotal = ref(0);    // 今日充值总额
const todayPaymentTotal = ref(0);     // 今日支付总额
const todayRechargeSelect = ref(0);   // 选中商店今日充值总额
const todayPaymentSelect = ref(0);    // 选中商店今日支付总额


// 添加今日流水相关的数据
const todayAmount = ref(0);
// 添加今日订单数相关的数据
const todayOrderCount = ref(0);

// 添加新的变量
const todayAmountSelect = ref(0);
const todayOrderCountSelect = ref(0);


const orderStatisticHandler = async()=>{
  const today = dayjs().startOf('day').toDate();
  const tonight = dayjs().endOf('day').toDate();
  const orderStatisticRes = await orderStatistic({
    status: 1,
    shopId: shopId.value,
    startDateAdd:today,
    endDateAdd:tonight
  });
  if (orderStatisticRes.code == 0) {
    // console.log(orderStatisticRes.data, "orderStatisticRes.data");
    todayAmount.value = orderStatisticRes.data.data.sum || 0;
    todayOrderCount.value = orderStatisticRes.data.data.count || 0;
    todayAmountSelect.value = orderStatisticRes.data.data.todaySum || 0;
    todayOrderCountSelect.value = orderStatisticRes.data.data.todayCount || 0;
    todayRechargeTotal.value = orderStatisticRes.data.data.todayRecharge || 0;
    todayRechargeSelect.value = orderStatisticRes.data.data.todayRechargeSelected || 0;
    todayPaymentTotal.value = orderStatisticRes.data.data.todayPayment || 0;
    todayPaymentSelect.value = orderStatisticRes.data.data.todayPaymentSelected || 0;
  }
}
const shopIdChangeHandler = (v)=>{
  console.log(v, "shopIdChangeHandler");
}
// 修改 watch 监听器
watch(shopId, async (newVal) => {
  if (newVal) {  // 只有当有选择值时才调用
    await orderStatisticHandler();
  } else {
    todayAmountSelect.value = 0; // 没有选择时清零
    todayOrderCountSelect.value = 0;
    todayRechargeSelect.value = 0;   // 使用新的变量
    todayPaymentSelect.value = 0;     // 使用新的变
  }
}, { immediate: true });  // 添加 immediate: true 确保初始化时也会触发



const init = async () => {
  await orderStatisticHandler();

  beeOrderStatus.value = await getDictFunc("OrderStatus");
  
  await refreshOrders();
  getAllMyBeeShopInfos().then((res) => {
    if (res.code == 0) {
      shops.value = res.data.list;
    }
  });
  const orderTable = await getBeeOrderList({
    page: 1,
    pageSize: 0,
    status: 1,
  });
  if (orderTable.code === 0) {
    orderTodoData.value.push(orderTable.data.total);
    orderTodoData.value.push(orderTable.data.total);
    orderTodo.value = orderTable.data;
  }
};
const spanMethodHandler = (row, column, rowIndex, columnIndex) => {
  // console.log(row, column, rowIndex, columnIndex);
  // if(row - 1 == rowIndex && columnIndex == 0){
  //   return [1,2]
  // }else {
  //   return [0,0]
  // }
};
const dateFmt = (time) => {
  if (time) {
    return dayjs(time).format("YYYY-MM-DD hh:mm:ss");
  }
  return "";
};
init();
// 加载订单流水
const refreshOrders = async () => {
  await orderStatisticHandler();
  const ordersRes =  await orderList({
    page: pageNum.value,
    pageSize: 0,
    shopId: shopId.value,
    startDateAdd:dates.value[0],
    endDateAdd:dates.value[1]
  });
  if (ordersRes.code == 0) {
    orders.value = [
      ...ordersRes.data.list,
      { amountReal: ordersRes.data.sum, orderNumber: "合计" },
    ];
    totalOrders.value = ordersRes.data.total;
  }

};
const gotoPage = async (page) => {
  page = "/layout" + page;
  await router.push(page);
};
// const formatDate = (time)=>{
//   return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
// }
const gotoTodoPage = (row) => {
  gotoPage("/bee_index/shop-order-admin/beeOrderDetail?id=" + row.id);
};


const userStore = useUserStore()

// 判断是否是管理员
const isAdmin = computed(() => {
  return userStore.userInfo.authorities?.some(role => role.admin === 1) || false
})

</script>

<style lang="scss" scoped></style>


