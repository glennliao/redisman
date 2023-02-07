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
    auto-focus

  >

    <div>
      <n-form
        ref="formRef"
        :label-width="80"
        :model="formValue"
        :rules="rules"
      >
        <n-form-item required label="title" path="title">
          <n-input v-model:value="formValue.title" placeholder="title"/>
        </n-form-item>
        <n-grid :x-gap="4">
          <n-gi :span="12">
            <n-form-item  label="host" path="host">
              <n-input v-model:value="formValue.host" placeholder="127.0.0.1"/>
            </n-form-item>
          </n-gi>
          <n-gi :span="12">
            <n-form-item label="port" path="port">
              <n-input :allow-input="(value: string) => !value || /^\d+$/.test(value)"  v-model:value="formValue.port" placeholder="6379"/>
            </n-form-item>
          </n-gi>
        </n-grid>
        <n-grid :x-gap="4">
          <n-gi :span="12">
            <n-form-item label="username" path="username">
              <n-input v-model:value="formValue.username" placeholder="username"/>
            </n-form-item>
          </n-gi>
          <n-gi :span="12">
            <n-form-item label="password" path="password">
              <n-input v-model:value="formValue.password" placeholder="password"/>
            </n-form-item>
          </n-gi>
        </n-grid>



        <div>
          <n-grid :cols="24">
            <n-gi :span="5">
              <n-switch disabled>
                <template #checked>
                  Readonly
                </template>
                <template #unchecked>
                  Readonly
                </template>
              </n-switch>
            </n-gi>
            <n-gi :span="4">

              <n-switch v-model:value="formValue.options.ssh.enable">
                <template #checked>
                  SSH
                </template>
                <template #unchecked>
                  SSH
                </template>
              </n-switch>
            </n-gi>
            <n-gi :span="4">
              <n-switch disabled>
                <template #checked>
                  SSL
                </template>
                <template #unchecked>
                  SSL
                </template>
              </n-switch>
            </n-gi>
            <n-gi :span="5">
              <n-switch disabled>
                <template #checked>
                  Sentinel
                </template>
                <template #unchecked>
                  Sentinel
                </template>
              </n-switch>
            </n-gi>
            <n-gi :span="5">
              <n-switch disabled>
                <template #checked>
                  Cluster
                </template>
                <template #unchecked>
                  Cluster
                </template>
              </n-switch>
            </n-gi>
          </n-grid>
        </div>

        <div>
          <div v-if="formValue.options.ssh.enable">
            <n-divider title-placement="left">
              SSH
            </n-divider>
            <div>
              <n-grid :x-gap="4">
                <n-gi :span="12">
                  <n-form-item required label="host" path="formValue.options.ssh.host">
                    <n-input v-model:value="formValue.options.ssh.host" placeholder="127.0.0.1"/>
                  </n-form-item>
                </n-gi>
                <n-gi :span="12">
                  <n-form-item required label="port" path="formValue.options.ssh.port">
                    <n-input :allow-input="(value: string) => !value || /^\d+$/.test(value)" :show-button="false" v-model:value="formValue.options.ssh.port" placeholder="22"/>
                  </n-form-item>
                </n-gi>
              </n-grid>

              <n-form-item required label="username" path="formValue.options.ssh.username">
                <n-input v-model:value="formValue.options.ssh.username" placeholder="username"/>
              </n-form-item>
              <n-form-item  label="authType" path="formValue.options.ssh.authType">
                <n-select default-value="password" :options="[{label:'password',value:'password'},{label:'privateKey',value:'privateKey'}]" v-model:value="formValue.options.ssh.authType" placeholder="authType"/>
              </n-form-item>
              <n-form-item required v-if="formValue.options.ssh.authType === 'password'" label="password" path="formValue.options.ssh.password">
                <n-input v-model:value="formValue.options.ssh.password" placeholder="password"/>
              </n-form-item>
              <n-form-item required v-if="formValue.options.ssh.authType === 'privateKey'" label="privateKey" path="formValue.options.ssh.privateKey">
                <n-input type="textarea" v-model:value="formValue.options.ssh.privateKey" placeholder="privateKey"/>
              </n-form-item>
              <n-form-item v-if="formValue.options.ssh.authType === 'privateKey'" label="passphrase" path="formValue.options.ssh.passphrase">
                <n-input v-model:value="formValue.options.ssh.passphrase" placeholder="passphrase"/>
              </n-form-item>

            </div>
          </div>
        </div>
      </n-form>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <div>
          <n-button @click="connTest">Test</n-button>
        </div>
        <div>
          <n-button @click="showModal=false">close</n-button>
          <n-button class="ml-2" @click="handleValidateClick">submit</n-button>
        </div>
      </div>
    </template>
  </n-modal>
</template>

<script lang="ts">
import {FormInst} from "naive-ui";
import cloneDeep from 'lodash-es/cloneDeep'
import {apijson,testConn} from "~/api/redis";

export default {
  name: "AddConnectionModal",
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

    const defaultInfoValue = {
      options:{
        ssh:{
          enable:false,
          authType:"password"
        }
      }
    }

    const formValue = ref(cloneDeep(defaultInfoValue) as Record<any, any>)

    const rules = {
      "title":[{required:true}]
    }

    const typeRef = ref('')



    function handleValidateClick(e: MouseEvent) {
      e.preventDefault()
      formRef.value?.validate((errors) => {
        if (!errors) {

          let formInfo = cloneDeep(formValue.value)
          formInfo.dbAlias = {}
          formInfo.host = formInfo.host || "127.0.0.1"
          formInfo.port = formInfo.port || "6379"

          apijson.post({
            tag:"RedisConnection",
            RedisConnection:formInfo
          }).then(()=>{
            message.success('Success')
            showModal.value = false
            formValue.value = cloneDeep(defaultInfoValue)
            emit("success")
          })

        } else {
          console.log(errors)
          message.error('Invalid')
        }
      })
    }

    let addFieldKey = ''
    let oldData: Record<any, any> | undefined = {}

    function open({type, key, data}: { type: string, key: string, data?: Record<any, any> }) {
      showModal.value = true
      typeRef.value = type
      formValue.value = cloneDeep(defaultInfoValue)

      oldData = {}
      if (key) {
        addFieldKey = key
        if (data) {
          Object.keys(data).forEach((k) => {
            // @ts-ignore
            formValue.value[k] = data[k]
          })
        }
        oldData = data
        // @ts-ignore
        formValue.value["key"] = key
      }
    }


    function connTest(){
      let formInfo = cloneDeep(formValue.value)
      formInfo.dbAlias = {}
      formInfo.host = formInfo.host || "127.0.0.1"
      formInfo.port = formInfo.port || "6379"


      testConn(formInfo).then((ret)=>{
        message.success('Success')
      }).catch(err=>{
        message.error(err.msg)
      })
    }


    return {
      open,
      showModal, rules, formValue, typeRef,
      bodyStyle, segmented, formRef,
      handleValidateClick,
      connTest
    }
  }
}
</script>

<style scoped>

</style>
