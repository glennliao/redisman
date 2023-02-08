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
    style="--n-padding-bottom:6px"
  >

    <div>
      <n-form
        ref="formRef"
        :label-width="80"
        :model="formValue"
        :rules="rules"
      >

        <n-tabs type="line" animated>

          <n-tab-pane name="general" tab="General">
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

          </n-tab-pane>
          <n-tab-pane name="ssh" tab="SSH" display-directive="show">
            <template #tab>
              SSH
              <n-badge class="ml-1" size="small" :dot="enableSSH" type="info"/>
            </template>
            <n-switch v-model:value="formValue.options.ssh.enable">
              <template #checked>
                SSH
              </template>
              <template #unchecked>
                SSH
              </template>
            </n-switch>
            <div class="mt-4">
              <n-grid :x-gap="4">
                <n-gi :span="12">
                  <n-form-item required label="host" path="options.ssh.host">
                    <n-input :disabled="!enableSSH" v-model:value="formValue.options.ssh.host" placeholder="127.0.0.1"/>
                  </n-form-item>
                </n-gi>
                <n-gi :span="12">
                  <n-form-item label="port" path="ssh.port">
                    <n-input :disabled="!enableSSH" :allow-input="(value: string) => !value || /^\d+$/.test(value)" :show-button="false" v-model:value="formValue.options.ssh.port" placeholder="22"/>
                  </n-form-item>
                </n-gi>
              </n-grid>

              <n-form-item required label="username" path="options.ssh.username">
                <n-input :disabled="!enableSSH" v-model:value="formValue.options.ssh.username" placeholder="username"/>
              </n-form-item>
              <n-form-item  label="authType" >
                <n-select :disabled="!enableSSH" default-value="password" :options="[{label:'password',value:'password'},{label:'privateKey',value:'privateKey'}]" v-model:value="formValue.options.ssh.authType"/>
              </n-form-item>
              <n-form-item required v-if="formValue.options.ssh.authType === 'password'" label="password" path="options.ssh.password">
                <n-input :disabled="!enableSSH" v-model:value="formValue.options.ssh.password" placeholder="password"/>
              </n-form-item>
              <n-form-item required v-if="formValue.options.ssh.authType === 'privateKey'" label="privateKey" path="options.ssh.privateKey">
                <n-input :disabled="!enableSSH" type="textarea" v-model:value="formValue.options.ssh.privateKey" placeholder="privateKey"/>
              </n-form-item>
              <n-form-item v-if="formValue.options.ssh.authType === 'privateKey'" label="passphrase" path="options.ssh.passphrase">
                <n-input :disabled="!enableSSH" v-model:value="formValue.options.ssh.passphrase" placeholder="passphrase"/>
              </n-form-item>
            </div>
          </n-tab-pane>

        </n-tabs>
      </n-form>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <div>
          <n-button ghost type="success" @click="connTest" :loading="testing">Test</n-button>
        </div>
        <div>
          <n-button @click="showModal=false">close</n-button>
          <n-button type="primary" class="ml-2" @click="handleValidateClick">submit</n-button>
        </div>
      </div>
    </template>
  </n-modal>
</template>

<script lang="ts">
import {FormInst} from "naive-ui";
import cloneDeep from 'lodash-es/cloneDeep'
import {apiJson,redis} from "~/api";

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

    const enableSSH = computed(()=>{
      return formValue.value.options.ssh.enable
    })

    const generialRules = {
      "title":[{required:true}]
    }

    const sshRules = {
      "options.ssh.host":[{required:true,message:"required"}],
      "options.ssh.username":[{required:true,message:"required"}],
    }

    const rules = computed(()=>{
      let rules:Record<string, any> = {
        ...generialRules,
      }
      if(enableSSH.value){
        Object.assign(rules, sshRules)
        if(formValue.value.options.ssh.authType === "password"){
          rules["options.ssh.password"] = [{required:true,message:"required"}]
        }else{
          rules["options.ssh.privateKey"] = [{required:true,message:"required"}]
        }
      }


      return rules
    })


    const typeRef = ref('')



    function handleValidateClick(e: MouseEvent) {
      e.preventDefault()
      formRef.value?.validate((errors) => {
        if (!errors) {

          let formInfo = fillDefaultField(cloneDeep(formValue.value))

          let api = formInfo.id ? apiJson.put : apiJson.post
          if(!formInfo.id){
            formInfo.dbAlias = {}
          }

          api({
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
          message.error('Form Invalid')
        }
      })
    }

    function open({id}: { id?: string }) {
      formValue.value = cloneDeep(defaultInfoValue)
      showModal.value = true
      if(id){
        apiJson.get({
          "RedisConnection":{
            id
          }
        }).then(data=>{
          let connection = data.RedisConnection
          formValue.value = {
            id:id+"",
            title:connection.title,
            host:connection.host,
            port:connection.port,
            username:connection.username,
            password:connection.password,
            db:connection.db,
            options:JSON.parse(connection.options),
          }
        })
      }

    }


    const testing = ref(false)
    function connTest(){
      testing.value = true
      let formInfo = fillDefaultField(cloneDeep(formValue.value))
      redis.connTest(formInfo).then((ret)=>{
        message.success('Success')
      }).catch(err=>{
        message.error(err.msg)
      }).finally(()=>{
        testing.value = false
      })
    }

    function fillDefaultField(info:Record<string, any>){
      info.host = info.host || "127.0.0.1"
      info.port = info.port || "6379"
      if(info.options.ssh.enable){
        info.options.ssh.port = info.options.ssh.port || "22"
      }
      return info
    }

    return {
      open,
      showModal, rules, formValue, typeRef,
      bodyStyle, segmented, formRef,enableSSH,
      handleValidateClick,
      connTest,testing
    }
  }
}
</script>

<style scoped>

</style>
