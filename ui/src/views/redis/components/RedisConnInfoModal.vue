<template>
  <n-modal
    v-model:show="showModal"
    class="custom-card"
    preset="card"
    :style="bodyStyle"
    title="Info"
    :bordered="false"
    :segmented="segmented"
  >

    <div v-if="info.Server">
      <n-tabs size="small" type="line" animated>
        <n-tab-pane name="info" tab="Info">
          <n-descriptions label-placement="left" class="p-2" :column="2">
            <n-descriptions-item label="redis_version">
              {{ info.Server.redis_version }}
            </n-descriptions-item>
            <n-descriptions-item label="pid">
              {{ info.Server.process_id }}
            </n-descriptions-item>
            <n-descriptions-item label="os">
              <span >{{ info.Server.os }}</span>
            </n-descriptions-item>

          </n-descriptions>

          <n-grid class="mt-4 p-2" cols="24" x-gap="12">
            <n-gi span="4">
              <n-statistic label="used_memory_peak" :value="info.Memory.used_memory_peak_human"/>
            </n-gi>
            <n-gi span="4">
              <n-statistic label="used_memory_lua" :value="info.Memory.used_memory_lua_human"/>
            </n-gi>
            <n-gi span="4">
              <n-statistic label="total_connections" :value="info.Stats.total_connections_received"/>
            </n-gi>
            <n-gi span="4">
              <n-statistic label="total_commands" :value="info.Stats.total_commands_processed"/>
            </n-gi>

          </n-grid>
<!--          <div class="mt-4">-->
<!--            <div class="flex flex-wrap p-2" style="border: 1px solid rgb(239, 239, 245)">-->
<!--              <div class="db-item shadow mx-2 rounded" v-for="item in dbList" :key="item.db">-->
<!--                {{ item.db }}-->
<!--                <n-number-animation ref="numberAnimationInstRef" :to="parseInt(item.keys||'0')"/>-->
<!--              </div>-->
<!--            </div>-->
<!--          </div>-->

          <div class="mt-2">
            <n-grid cols="24" x-gap="6" y-gap="6">
              <n-gi span="12">
                <use-memory-line :data="useMemoryData"/>
              </n-gi>
              <n-gi span="12">
                <client-num-line :data="clientNumData"/>
              </n-gi>
            </n-grid>
          </div>

        </n-tab-pane>
        <n-tab-pane v-for="(sectionVal,section) in info" :name="section" :tab="section">
          <div style="overflow-y: auto;height: calc(86vh - 230px)">
            <n-list>
              <n-list-item v-for="(val,key) in sectionVal" :key="key">
                <n-thing>
                  <div class="flex">
                    <div style="width:250px">{{ key }}</div>
                    {{ val }}
                  </div>
                </n-thing>
              </n-list-item>
            </n-list>
          </div>

        </n-tab-pane>
      </n-tabs>

    </div>
    <template #footer>
      <div class="flex justify-between">
        <div class="flex items-center justify-between" style="width:180px">
          <text>Refresh</text>
          <n-input-number size="small" v-model:value="second" placeholder="" :show-button="false" style="width:60px">
            <template #suffix>s</template>
          </n-input-number>
          <n-switch v-model:value="autoRefresh"/>
        </div>
        <n-button @click="showModal=false">close</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script lang="ts">

import {redis} from '../../../api'
import {parseInfo} from "~/utils/redis";
import {useConn} from "~/views/redis/hook/conn";
import UseMemoryLine from "~/components/info/UseMemoryLine.vue";
import ClientNumLine from "~/components/info/ClientNumLine.vue";

export default {
  name: "RedisConnInfoModal",
  components: {ClientNumLine, UseMemoryLine},
  emits: ["success"],
  setup(_: any, {emit}: any) {

    const info = ref({
      Memory: {used_memory_human: ""},
      Clients: {connected_clients: ''},
      Keyspace:{} as Record<string, string>,
    })

    const second = ref(5)
    const autoRefresh = ref(false)

    let timer: any = null
    watch(autoRefresh, (val) => {

      if (val) {
        timer = setInterval(loadInfo, second.value * 1000)
      } else {
        clearInterval(timer)
      }

    })


    function loadInfo() {
      return redis.command([["info", "all"]]).then(data => {
        info.value = {
          ...parseInfo(data[0]),
        } as any

        useMemoryData.value = useMemoryData.value.concat([{
          time: new Date().getTime(),
          value: info.value.Memory.used_memory_human
        }] as any)


        clientNumData.value = clientNumData.value.concat([{
          time: new Date().getTime(),
          value: info.value.Clients.connected_clients
        }] as any)
      })
    }


    const showModal = ref(false)

    watch(showModal, (val) => {
      if (!val) {
        timer && clearInterval(timer)
      }
    })

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
      autoRefresh.value = false
      useMemoryData.value = []
      clientNumData.value = []
      loadInfo()
    }

    let useMemoryData = ref([])
    let clientNumData = ref([])


    const conn = useConn()

    const dbList = computed(()=>{
      const keyspace = toRaw(info.value.Keyspace || {}) as unknown as Record<string, Record<string, any>>;
      Object.keys(keyspace).forEach((k) => {
        if(k.startsWith("db")){
          const sps = keyspace[k].split(",");
          keyspace[k.substring(2)] = {
            keys: sps[0].substring("keys=".length),
          };
        }
      });

      return conn.dbList.value.map(item=>{
        return {
          db:item.value,
          ...keyspace[item.value]||{}
        }
      })

    })


    return {
      open,
      showModal,
      bodyStyle, segmented,
      info, second, autoRefresh,dbList,
      useMemoryData,clientNumData
    }
  }
}
</script>

<style scoped>
.db-item{
  width:40px;
  height: 40px;

}
</style>
