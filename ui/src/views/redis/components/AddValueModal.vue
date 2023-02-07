<template>
  <n-modal
    v-model:show="showModal"
    class="custom-card"
    preset="card"
    :style="bodyStyle"
    title="Add"
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
      >
        <n-form-item label="key" path="key">
          <n-input v-model:value="formValue.key" placeholder="key" />
        </n-form-item>
        <n-form-item v-if="typeRef === RedisTypes.Hash || typeRef === 'stream'" label="field" path="field">
          <n-input  v-model:value="formValue.field" placeholder="field" />
        </n-form-item>
        <n-form-item v-if="typeRef === RedisTypes.ZSet" label="score" path="score">
          <n-input  v-model:value="formValue.score" placeholder="score" />
        </n-form-item>
        <n-form-item label="value" path="value">
          <n-input type="textarea" v-model:value="formValue.value" placeholder="value" />
        </n-form-item>
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
import {useValueHook} from "~/views/redis/hook/value";
import {useKeysHook} from '../hook/keys'
import {RedisTypes} from "~/views/redis/redis_types";



export default {
  name: "AddValueModal",
  computed: {
    RedisTypes() {
      return RedisTypes
    }
  },
  emits:["success"],
  setup(_: any, {emit}: any){

    const showModal = ref(false)
    const bodyStyle = {
      width:"600px"
    }
    const segmented = {
      content: 'soft',
      footer: 'soft'
    }

    const formRef = ref<FormInst | null>(null)
    const message = useMessage()
    const formValue = ref({} as Record<any, any>)

    const rules = {

    }

    const typeRef = ref('')

    const useValue = useValueHook()

    const {scan} = useKeysHook()




    function handleValidateClick(e: MouseEvent) {
      e.preventDefault()
      formRef.value?.validate((errors) => {
        if (!errors) {

          let {key,field, member,score,value,id} = formValue.value as any

          useValue.set(key,typeRef.value,{field,member,score,value,id}, oldData as any)
            .then(()=>{
              scan()
              message.success('Success')
              showModal.value = false
              formValue.value = {}
              emit("success",{
                key:addFieldKey
              })
            })
        } else {
          console.log(errors)
          message.error('Invalid')
        }
      })
    }

    let addFieldKey = ''
    let oldData:Record<any, any>|undefined = {}

    function open({type,key,data}:{type:string,key:string,data?:Record<any, any>}){
      showModal.value = true
      typeRef.value = type
      oldData = {}
      if(key){
        addFieldKey = key
        if(data){
          Object.keys(data).forEach((k)=>{
            // @ts-ignore
            formValue.value[k] = data[k]
          })
        }
        oldData = data
        // @ts-ignore
        formValue.value["key"] = key
      }
    }


    return {
      open,
      showModal,rules,formValue,typeRef,
      bodyStyle,segmented,formRef,
      handleValidateClick
    }
  }
}
</script>

<style scoped>

</style>
