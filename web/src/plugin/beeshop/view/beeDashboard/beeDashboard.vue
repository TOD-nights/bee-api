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
    style="width: 100%"
    :span-method="spanMethodHandler"
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
      prop="amount"
      label="订单金额"
      align="center"
      show-overflow-tooltip
    />
    <el-table-column 
      prop="datePay" 
      label="支付时间"
      min-width="120"
    >
      <template #default="scope">
        {{ dateFmt(scope.row.dateAdd) }}
      </template>
    </el-table-column>
  </el-table>
</gva-card>


<!-- 添加今日流水卡片 -->
 <!-- 全部商店-->
 <el-card @click="gotoPage('/bee_index/beeFinancialManager/beePayLog')">
  <el-statistic 
    :value="todayAmount" 
    :precision="2"
  >
    <template #title>
      <div class="font-bold">今日全部商店充值和消费金额</div>
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
      <div class="font-bold">今日选中商店充值和消费金额</div>
    </template>
    <template #suffix>
      <span class="text-sm">元</span>
    </template>
  </el-statistic>
</el-card>


    
    <!-- 添加今日订单数卡片 -->


    <el-card @click="gotoPage('/bee_index/shop-order-admin/beeOrder')">
  <el-statistic 
    :value="todayOrderCount" 
    :precision="2"
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
    :precision="2"
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



<el-card @click="gotoPage('/bee_index/beeFinancialManager/beePayLog')">
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



<!-- 添加选中商店今日充值卡片 -->


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





<!-- 添加今日支付卡片 -->


<el-card @click="gotoPage('/bee_index/beeFinancialManager/beeUserBalanceLog')">
  <el-statistic 
    :value="todayPaymentTotal" 
    :precision="2"
  >
    <template #title>
      <div class="font-bold">今日全部商店余额支付金额</div>
    </template>
    <template #suffix>
      <span class="text-sm">元</span>
    </template>
  </el-statistic>
</el-card>



<!-- 添加选中商店今日支付卡片 -->


<el-card @click="gotoPage('/bee_index/beeFinancialManager/beeUserBalanceLog')">
  <el-statistic 
    :value="todayPaymentSelect" 
    :precision="2"
  >
    <template #title>
      <div class="font-bold">今日选择商店余额支付金额</div>
    </template>
    <template #suffix>
      <span class="text-sm">元</span>
    </template>
  </el-statistic>
</el-card>



    <!--
    <gva-card
      custom-class="col-span-1 lg:col-span-2 h-32"
      @click="gotoPage('/bee_index/shop-user-admin/beeUser')"
    >
      <gva-chart :data="userData" title="用户数" />
    </gva-card>


    <gva-card
      custom-class="col-span-1 lg:col-span-2 h-32"
      @click="gotoPage('/bee_index/shop-order-admin/beeOrder')"
    >
      <gva-chart :data="orderData" title="订单统计" />
    </gva-card>
    
    
    <gva-card
      custom-class="col-span-1 lg:col-span-2 h-32"
      @click="gotoPage('/bee_index/beeOrderTodo')"
    >
      <gva-chart :data="orderTodoData" title="订单待办" />
    </gva-card>


    <gva-card
      custom-class="col-span-1 lg:col-span-2 h-32"
      @click="gotoPage('/bee_index/beeFinancialManager/beePayLog')"
    >
      <gva-chart :data="payData" title="总支付金额" />
    </gva-card>

 添加今日流水卡片 

    <gva-card
      custom-class="col-span-1 lg:col-span-2 h-32"
      @click="gotoPage('/bee_index/beeFinancialManager/beePayLog')"
    >
      <gva-chart :data="payNumData" title="支付人数" />
    </gva-card>-->


    <gva-card
      title="待办事项"
      custom-class="col-span-1 md:col-span-6 lg:col-span-6 row-span-2"
    >
      <div>
        <el-table
          :data="orderTodo.list"
          stripe
          style="width: 100% ; height: 300px "
          @row-click="gotoTodoPage"
        >
          <el-table-column prop="orderNumber" label="订单号" width="200" />
          <el-table-column
            prop="id"
            label="事项"
            align="center"
            show-overflow-tooltip
          >
            <template #default="scope">
              {{ formatEnum(scope.row.status, beeOrderStatus) }}
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="金额" width="100" />
          <el-table-column prop="dateAdd" label="支付时间">
            <template #default="scope">{{
              formatDate(scope.row.dateAdd)
            }}</template>
          </el-table-column>
        </el-table>
      </div>
    </gva-card>
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
import { getBeeOrderList, orderList } from "@/plugin/beeshop/api/beeOrder";
import { getBeeShopGoodsList } from "@/plugin/beeshop/api/beeShopGoods";
import { formatDate, formatEnum, getDictFunc } from "@/utils/format";
import { useRoute, useRouter } from "vue-router";
import { getAllMyBeeShopInfos } from "@/plugin/beeshop/api/beeShopInfo";
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

// 在 init 函数中添加获取今日流水的逻辑
const getTodayAmount = async () => {
  const today = dayjs().startOf('day').toDate();
  const tonight = dayjs().endOf('day').toDate();
    // 调用API获取今日支付总额
  const todayPayment = await getBeePayTotal({
    page: 1,
    pageSize: 1,
    startDateAdd: today,
    endDateAdd: tonight,
    sum: "money"
  });
  // 如果请求成功，更新todayAmount的值
  if (todayPayment.code === 0) {
    todayAmount.value = todayPayment.data.total || 0;
  }
};

// 添加新的函数来获取选中商店的今日流水
// 修改获取选中商店的今日流水的函数
const getTodayAmountSelect = async () => {
  const today = dayjs().startOf('day').toDate();
  const tonight = dayjs().endOf('day').toDate();
  
  // 使用 orderList 而不是 getBeePayTotal
  const todayOrders = await orderList({
    page: 1,
    pageSize: 1000,  // 设置较大的数值以获取所有订单
    shopId: shopId.value,
    startDateAdd: today,
    endDateAdd: tonight
  });

  if (todayOrders.code === 0) {
    // 计算总金额
    todayAmountSelect.value = todayOrders.data.sum || 0;
  }
};

// 添加新函数来获取选中商店的今日订单数
const getTodayOrderCountSelect = async () => {
  const today = dayjs().startOf('day').toDate();
  const tonight = dayjs().endOf('day').toDate();
  
  const todayOrders = await orderList({
    page: 1,
    pageSize: 1,
    shopId: shopId.value,
    startDateAdd: today,
    endDateAdd: tonight
  });

  if (todayOrders.code === 0) {
    todayOrderCountSelect.value = todayOrders.data.total || 0;
  }
};


// 获取今日订单数的方法
const getTodayOrderCount = async () => {
  const today = dayjs().startOf('day').toDate();
  const tonight = dayjs().endOf('day').toDate();
  
  const todayOrders = await getBeeOrderList({
    page: 1,
    pageSize: 1,
    startDateAdd: today,
    endDateAdd: tonight
  });

  if (todayOrders.code === 0) {
    todayOrderCount.value = todayOrders.data.total || 0;
  }
};


// 获取今日余额变动数据，区分充值和支付
const getTodayBalanceLog = async () => {
  const today = dayjs().startOf('day').toDate();
  const tonight = dayjs().endOf('day').toDate();
  
  const balanceLog = await getBeeUserBalanceLogList({
    page: 1,
    pageSize: 1000,
    startDateAdd: today,
    endDateAdd: tonight
  });
// 分别获取充值和支付记录
const rechargeLog = await getBeeUserBalanceLogCount({
    page: 1,
    pageSize: 1000,
    startDateAdd: today,
    endDateAdd: tonight,
    type: 'recharge'  // 获取充值记录
  });

  getBeeUserBalanceLogCount({
    startDateAdd: today,
    endDateAdd: tonight,
    type: 'recharge',  // 获取充值记录
    shopId: shopId.value
  }).then(res=>{
    if(res.code == 200){
      todayRechargeSelect.value = res.data
    }
  })

  const paymentLog = await getBeeUserBalanceLogList({
    page: 1,
    pageSize: 1000,
    startDateAdd: today,
    endDateAdd: tonight,
    type: 'payment'  // 获取支付记录
  });

  if (rechargeLog.code === 0) {
    console.log(rechargeLog,rechargeLog.data.data)

    todayRechargeTotal.value = rechargeLog.data.data;
  }

  if (paymentLog.code === 0) {
    todayPaymentTotal.value = Math.abs(paymentLog.data.list.reduce((sum, item) => sum + Number(item.num), 0));
  }
  
  if (balanceLog.code === 0) {
    // 分别计算充值和支付金额
    const rechargeList = balanceLog.data.list.filter(item => item.mark === '充值');
    const paymentList = balanceLog.data.list.filter(item => item.mark === '订单支付');
    
    // todayRechargeTotal.value = rechargeList.reduce((sum, item) => sum + Number(item.num), 0);
    // todayPaymentTotal.value = Math.abs(paymentList.reduce((sum, item) => sum + Number(item.num), 0));
  }
};
// 获取选中商店的余额变动数据，区分充值和支付
const getTodayBalanceLogSelect = async () => {
  const today = dayjs().startOf('day').toDate();
  const tonight = dayjs().endOf('day').toDate();
  
   // 打印 API 请求参数
   console.log('API 请求参数:', {
    page: 1,
    pageSize: 1000,
    startDateAdd: today,
    endDateAdd: tonight,
    shopId: shopId.value
    
  });


  const balanceLog = await getBeeUserBalanceLogList({
    page: 1,
    pageSize: 1000,
    startDateAdd: today,
    endDateAdd: tonight,
    shopId: shopId.value
  });
  console.log('API 返回的完整响应:', balanceLog); // 查看完整响应
  console.log('余额变动数据:', balanceLog.data.list);
  
  if (balanceLog.code === 0) {
      // 检查所有不同类型的 mark
      const uniqueMarks = [...new Set(balanceLog.data.list.map(item => item.mark))];
    console.log('所有 mark 类型:', uniqueMarks);

    const rechargeList = balanceLog.data.list.filter(item => item.mark === '充值');
    const paymentList = balanceLog.data.list.filter(item => item.mark === '订单支付');
    
    console.log('充值记录:', rechargeList);
    console.log('支付记录:', paymentList);

    
   
    // todayRechargeSelect.value = rechargeList.reduce((sum, item) => sum + Number(item.num), 0);
    todayPaymentSelect.value = Math.abs(paymentList.reduce((sum, item) => sum + Number(item.num), 0));

    console.log('计算后的充值金额:', todayRechargeSelect.value);
    console.log('计算后的支付金额:', todayPaymentSelect.value);
  }
};

// 修改 watch 监听器
watch(shopId, async (newVal) => {
  if (newVal) {  // 只有当有选择值时才调用
    await getTodayAmountSelect();
    await getTodayOrderCountSelect();
    await getTodayBalanceLogSelect();
  } else {
    todayAmountSelect.value = 0; // 没有选择时清零
    todayOrderCountSelect.value = 0;
    todayRechargeSelect.value = 0;   // 使用新的变量
    todayPaymentSelect.value = 0;     // 使用新的变
  }
}, { immediate: true });  // 添加 immediate: true 确保初始化时也会触发



const init = async () => {
  await getTodayAmount(); // 初始化时获取今日流水
  await getTodayAmountSelect(); // 初始化时获取今日选择商店流水
  await getTodayOrderCount(); // 获取今日订单数
  await getTodayOrderCountSelect();
  await getTodayBalanceLog();
  await getTodayBalanceLogSelect();

  beeOrderStatus.value = await getDictFunc("OrderStatus");
  for (let i = 7; i >= 0; i--) {
    // 获取最近7天数据
    const start = new Date(0);
    const end = dayjs()
      .add(-1 * i, "day")
      .endOf("day")
      .toDate();

    const funcs = [];
    funcs.push(
      getBeeUserList({
        page: 1,
        pageSize: 1,
        startDateAdd: start,
        endDateAdd: end,
      })
    );
    funcs.push(
      getBeePayLogList({
        page: 1,
        pageSize: 1,
        startDateAdd: start,
        endDateAdd: end,
        distinct: "uid",
      })
    );
    funcs.push(
      getBeeOrderList({
        page: 1,
        pageSize: 1,
        startDateAdd: start,
        endDateAdd: end,
      })
    );
    funcs.push(
      getBeeShopGoodsList({
        page: 1,
        pageSize: 1,
        startDateAdd: start,
        endDateAdd: end,
      })
    );
    funcs.push(
      getBeePayTotal({
        page: 1,
        pageSize: 1,
        startDateAdd: start,
        endDateAdd: end,
        sum: "money",
      })
    );
    const results = await Promise.all(funcs);
    const userTable = results[0];
    if (userTable.code === 0) {
      userData.value.push(userTable.data.total);
    }
    const payLogUidTable = results[1];
    if (payLogUidTable.code === 0) {
      payNumData.value.push(payLogUidTable.data.total);
    }
    const orderTable = await results[2];
    if (orderTable.code === 0) {
      orderData.value.push(orderTable.data.total);
    }
    const goodsTable = await results[3];
    if (goodsTable.code === 0) {
      goodsData.value.push(goodsTable.data.total);
    }
    const payDataTable = await results[4];
    if (payDataTable.code === 0) {
      payData.value.push(payDataTable.data.total);
    }
  }
  await refreshOrders();
  getAllMyBeeShopInfos().then((res) => {
    if (res.code == 0) {
      shops.value = res.data.list;
    }
  });
  const orderTable = await getBeeOrderList({
    page: 1,
    pageSize: 100,
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
  const ordersRes =  await orderList({
    page: pageNum.value,
    pageSize: 10,
    shopId: shopId.value,
    startDateAdd:dates.value[0],
    endDateAdd:dates.value[1]
  });
  if (ordersRes.code == 0) {
    console.log(ordersRes.data.list, "ordersRes.data.list");
    orders.value = [
      ...ordersRes.data.list,
      { amount: ordersRes.data.sum, orderNumber: "合计" },
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
</script>

<style lang="scss" scoped></style>
