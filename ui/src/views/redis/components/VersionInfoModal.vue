<template>
  <div style="cursor: pointer" class="flex" @click="showModal = true">
    RedisMan {{ version }}
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
      <div class="flex">
        <div>
          https://github.com/glennliao/redisman
        </div>
        <div class="ml-2">
          <div><a href="https://github.com/glennliao/redisman" target="_blank" style="text-decoration: none">
            <n-button size="small" type="primary">download</n-button>
          </a></div>
        </div>
      </div>
      <div class="flex mt-2">
        <div>
          curVersion :
          <n-tag type="info">
            {{ version }}
        </n-tag>
        </div>
        <div class="ml-2">
          latest: {{ latest }}
          <n-tag type="success">
            {{ latest }}
          </n-tag>
        </div>
        <div class="ml-4">
          <n-button size="small" :loading="checking" @click="checkNewVersion" type="info">
            check
          </n-button>
        </div>
      </div>


      <div class="flex mt-2">

      </div>
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

    const version = import.meta.env.VITE_app_version.trim()
    const latest = ref(version)

    const checking = ref(false)

    function checkNewVersion() {
      checking.value = true
      checkVersion().then(data => {
        latest.value = data.latest.Version?.trim()
      }).finally(()=>{
        checking.value = false
      })
    }

    checkNewVersion()
    return {
      open,
      showModal,
      bodyStyle, segmented,
      checkNewVersion, version, latest,checking

    }
  }
}
</script>

<style scoped>

</style>
