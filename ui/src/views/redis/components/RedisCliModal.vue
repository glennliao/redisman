<template>
  <n-modal
    v-model:show="showModal"
    class="custom-card"
    preset="card"
    :style="bodyStyle"
    title="Cli"
    :bordered="false"
    :segmented="segmented"
  >

    <div v-if="showModal">
      <terminal :init-log="[]" :show-header="false" context="redis" style="height: 74vh" name="my-terminal" @execCmd="onExecCmd"></terminal>
    </div>
    <template #footer>
      <div class="flex justify-between">
        <n-button @click="showModal=false">close</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script lang="ts">
import 'xterm/css/xterm.css'
import {redis} from '../../../api'
import Terminal from 'vue-web-terminal'
import splitargs from "redis-splitargs"

export default {
  name: "RedisCliModal",
  components: {Terminal},
  setup(_: any) {

    const showModal = ref(false)

    const bodyStyle = {
      width: "86vw",
      height: "86vh",
      overflowY: "auto",
    }
    const segmented = {
      content: 'soft',
      footer: 'soft'
    }

    function open() {
      showModal.value = true
    }

    function onExecCmd(key:string, command:string, success:(p:any)=>void, failed:(p:any)=>void){
      if (key === 'fail') {
        failed('Something wrong!!!')
      } else {
        let [comm,...args] = splitargs(command)
        redis.cliCommand([[comm,...args]]).then((data)=>{

          data = data[0]
          if(!Array.isArray(data)){
            data = [data]
          }

          data.forEach((item:string)=>{
            console.log(item)
            success({
              type: 'html',
              content: "<pre>"+item+"</pre>"
            })
          })
        })


      }
    }

    return {
      open,
      showModal,
      bodyStyle, segmented,
      onExecCmd
    }
  }
}
</script>

<style scoped>

</style>
