<template>
  <div style="cursor: pointer" class="flex" @click="showModal = true">
    RedisMan {{version}}
   <n-badge v-if="version !== latest" dot>
      (new)
   </n-badge>
  </div>
  <n-modal
    v-model:show="showModal"
    class="custom-card"
    preset="card"
    :style="bodyStyle"
    title="RedisMan"
    :bordered="false"
    :mask-closable="false"
    :segmented="segmented"
  >

    <div>
     <div>
       curVersion : {{version}}
     </div>
      <div >
        latest: {{latest}}
      </div>

     <div>
       <a href="https://github.com/glennliao/redisman" target="_blank">to download</a>
     </div>
      <n-button @click="checkNewVersion">check</n-button>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <n-button @click="showModal=false">close</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script lang="ts">

import {checkVersion} from "~/api";


export default {
  name: "VersionInfoModal",
  setup(_: any, {emit}: any) {

    const showModal = ref(false)
    const bodyStyle = {
      width: "600px"
    }
    const segmented = {
      content: 'soft',
      footer: 'soft'
    }

    const version = import.meta.env.VITE_app_version
    const latest = ref(version)

    function checkNewVersion(){
      checkVersion().then(data=>{
        latest.value = data.latest.Version
      })
    }

    checkNewVersion()
    return {
      open,
      showModal,
      bodyStyle, segmented,
      checkNewVersion,version,latest

    }
  }
}
</script>

<style scoped>

</style>
