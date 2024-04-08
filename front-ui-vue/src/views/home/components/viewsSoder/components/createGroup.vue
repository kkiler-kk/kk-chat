<template>
  <div>
    <a-modal v-model:visible="open" @ok="handleOk">
      <a-transfer
        v-model:target-keys="formState.membersId"
        v-model:selected-keys="selectedKeys"
        :data-source="data"
        show-search
        :filter-option="filterOption"
        :titles="['选择', '目标']"
        :render="(item) => `${item.title}(${item.description})`"
        @change="handleChange"
      />
      <a-input v-model:value="formState.name" placeholder="请输入群名称" />
      <!-- @selectChange="handleSelectChange"
        @scroll="handleScroll" -->
    </a-modal>

  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { listFriend } from "@/api/userFriend/userFriend";
import {groupAdd} from '@/api/group/index'
import { Session } from "@/utils/storage";
import { message } from "ant-design-vue";
const open = ref<boolean>(false);
const handleOk = (e: MouseEvent) => {
  if (formState.value.name == "") {
    return
  }
  let userInfo = Session.get("userInfo")
  formState.value.ownerId = userInfo.id
  groupAdd(formState.value).then((res) => {
    console.log("创建成功")
  }).catch((err) => {
    message.error("创建群聊失败", err)
    console.log(err)
  })
  open.value = false
};
interface MockData {
  key: string;
  title: string;
  description: string;
  chosen: boolean;
}
const showModal = () => {
  open.value = true;
};
const hideModal = () => {
  open.value = false;
};
const formState = ref<any>({
  name: '',
  ownerId: 0,
  membersId: [],
})
onMounted(() => {
  handleListFriend();
});
// const oriTargetKeys = ref(string[])()

const targetKeys = ref<number[]>([]);
const selectedKeys = ref<number[]>([]);
const data = ref<MockData[] | undefined>();
const handleChange = (
  nextTargetKeys: string[],
  direction: string,
  moveKeys: number[]
) => {
  console.log("nextTargetKeys", nextTargetKeys)
  console.log("direction", direction)
  console.log("moveKeys", moveKeys)
  formState.value.membersId = moveKeys
};
onMounted(() => {
  handleListFriend();
});
const handleListFriend = () => {
  const mData = [];
  listFriend().then((res: any) => {
    for (let i =0; i < res.length; i++) {
      let dataTemp = {
        key: res[i].id,
        title: res[i].name,
        description: res[i].identity,
        disabled: false,
      }
      mData.push(dataTemp)
    }
    data.value = mData
  });
};
const filterOption = (inputValue: string, option: any) => {
  return option.name.indexOf(inputValue) > -1;
};
//暴露state和play方法
defineExpose({
  showModal,
  hideModal,
});
</script>
<style scoped></style>
