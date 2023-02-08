<template>
  <n-modal
    v-model:show="showModal"
    class="custom-card"
    preset="card"
    :style="bodyStyle"
    title="Setting (DbAlias)"
    :bordered="false"
    :mask-closable="false"
    :segmented="segmented"
  >

    <div>
      <n-form
        ref="formRef"
        :label-width="80"
        :model="formValue"
        :rules="rules"
        label-placement="left"
      >
        <n-grid>
          <n-gi :span="12" v-for="item in dbAliasList" :key="item.label">
            <n-form-item  :label="item.label" :path="item.label">
              <n-input v-model:value="formValue[item.label]" :placeholder="item.label"/>
            </n-form-item>
          </n-gi>
        </n-grid>

      </n-form>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <n-button @click="showModal=false">close</n-button>
        <n-button class="ml-2" @click="handleValidateClick">submit</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script lang="ts">
import {FormInst} from "naive-ui";

import {apiJson} from "~/api";
import { useConn } from "../hook/conn";


export default {
  name: "RedisConnSettingModal",
  emits: ["success"],
  setup(_: any, {emit}: any) {

    const showModal = ref(false)
    const bodyStyle = {
      width: "600px"
    }
    const segmented = {
      content: 'soft',
      footer: 'soft'
    }

    const formRef = ref<FormInst | null>(null)
    const message = useMessage()
    const formValue = ref({} as Record<any, any>)

    const rules = {}

    const typeRef = ref('')



    function handleValidateClick(e: MouseEvent) {
      e.preventDefault()
      formRef.value?.validate((errors) => {
        if (!errors) {

          let dbAliasOri = formValue.value
          let dbAlias:Record<string, string> = {}
          Object.keys(dbAliasOri).forEach(k=>{

            dbAlias[k+""] = dbAliasOri[k]
          })


          apiJson.put({
            tag:"RedisConnection",
            RedisConnection:{
              dbAlias,
              id:connMeta.value.id+"",
            }
          }).then(()=>{
            console.log()
            message.success('Success')
            showModal.value = false
            formValue.value = {}
            emit("success")
          })

        } else {
          console.log(errors)
          message.error('Invalid')
        }
      })
    }

    function open() {
      showModal.value = true
      formValue.value = connMeta.value.dbAlias
    }



    const {dbList,connMeta} = useConn()
    const dbAliasList = computed(()=>{
      return dbList.value.map(item=>{
        return {
          label:item.label+""
        }
      })
    })
    return {
      open,
      showModal, rules, formValue, typeRef,
      bodyStyle, segmented, formRef,
      handleValidateClick,dbAliasList
    }
  }
}
</script>

<style scoped>

</style>
