<template>

  <n-grid y-gap="8" x-gap="8" cols="24"  item-responsive>
    <n-gi span="24 320:24 560:12 800:12">
      <n-input-group>
        <n-input-group-label size="small">Key</n-input-group-label>
        <n-input size="small" v-model:value="keyVal">
          <template #suffix>
            <n-icon title="save" class="cursor-pointer" @click="changeKey">
              <SaveOutline/>
            </n-icon>
          </template>
        </n-input>
      </n-input-group>
    </n-gi>
    <n-gi span="24 320:16 560:8">
      <n-input-group>
        <n-input-group-label size="small">TTL</n-input-group-label>
        <n-input-number size="small" :show-button="false" v-model:value="ttlVal">
          <template #suffix>
            <div title="forever" class="cursor-pointer" @click="changeTTL('-')">
              X
            </div>
            <n-icon title="save" class="cursor-pointer ml-2" @click="changeTTL">
              <SaveOutline/>
            </n-icon>
          </template>
        </n-input-number>
      </n-input-group>
    </n-gi>
    <n-gi span="24 320:8 560:4">
      <n-button-group size="small">
        <n-popconfirm
          @positive-click="delKey"
        >
          <template #trigger>
            <n-button ghost title="删除">
              <template #icon>
                <n-icon>
                  <RemoveCircleOutline/>
                </n-icon>
              </template>
            </n-button>
          </template>
          del {{curKey}}?
        </n-popconfirm>

        <n-button ghost @click="refresh(keyVal)">
          <template #icon>
            <n-icon>
              <Refresh/>
            </n-icon>
          </template>

        </n-button>
      </n-button-group>
    </n-gi>
  </n-grid>

  <n-grid class="mt-2" y-gap="8" x-gap="8" cols="24" >
    <n-gi span="12">
      <n-tag round :bordered="false"  type="success">
        <div class="text-center" style="min-width: 30px">
          {{ type }}
        </div>
      </n-tag>
    </n-gi>
    <n-gi span="8">
      <n-tag round :bordered="false"  type="info">
        <div class="text-center" :title="ttlHumanizerDetail" style="max-width: 120px;overflow: hidden;">
          {{ ttlHumanizer }}
        </div>
      </n-tag>
    </n-gi>
    <n-gi span="4">
      <n-button size="small" v-if="![RedisTypes.String].includes(type)" strong  ghost round type="success" @click="addField">
        <template #icon>
          <n-icon >
            <Add/>
          </n-icon>
        </template>
      </n-button>
    </n-gi>
  </n-grid>


  <div class="mt-2">
    <div v-if="![RedisTypes.String].includes(type)">
      <n-data-table
        class="mt-2"
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :bordered="false"
        size="small"
      />
    </div>

    <KeyValue v-else :value="value" @save="saveVal"></KeyValue>
  </div>


</template>

<script setup>
import {useValueHook} from "~/views/redis/hook/value";
import KeyValue from "~/views/redis/value.vue";
import {Refresh, RemoveCircleOutline,SaveOutline,Add} from '@vicons/ionicons5'
import {NButton,useMessage } from 'naive-ui'
import {useKeysHook} from "~/views/redis/hook/keys";
import {useConn} from "~/views/redis/hook/conn";
import {RedisTypes} from "~/views/redis/redis_types";
import humanizeDuration from 'humanize-duration'

const shortEnglishHumanizer = humanizeDuration.humanizer({
  language: "shortEn",
  languages: {
    shortEn: {
      y: () => "y",
      mo: () => "mo",
      w: () => "w",
      d: () => "d",
      h: () => "h",
      m: () => "m",
      s: () => "s",
      ms: () => "ms",
    },
  },
});

const ttlHumanizer = computed(()=>{
  if (ttl.value == -1){
    return "forever"
  }
  return shortEnglishHumanizer(ttl.value*1000,{units:["h","m","s"]})
})
const ttlHumanizerDetail = computed(()=>{
  if (ttl.value == -1){
    return "forever"
  }
  return shortEnglishHumanizer(ttl.value*1000)
})



const props = defineProps({
  curKey: {
    type:String
  }
})

const {
  type, ttl, load, value, len, list,set,del, expire,persist,rename
} = useValueHook()

const {scan} = useKeysHook()

const ttlVal = ref(0)
const keyVal = ref(props.curKey)
let oriKey = props.curKey

watch(() => props.curKey, (val) => {
  load(val)
  keyVal.value = val
  oriKey = val
}, {immediate: true})



watch(ttl,(val)=>{
  ttlVal.value = parseInt(val)
},{immediate:true})


function editField(row){

  emit("addField", {
    key:props.curKey,
    type:type.value,
    data:row,
  })

}

function delField(row){

  del({
    key:keyVal.value,
    type:type.value,
    data:row
  }).then(()=>{
    refresh(keyVal.value)
    message.success("Success")
  })
}

const columns = computed(()=>{
  let col = [
    {
      title: 'No',
      key: 'no',
      width:60,
    },
  ]

  if(type.value === RedisTypes.Hash){
    col.push({
      title: 'field',
      key: 'field'
    })
  }

  if(type.value ===  RedisTypes.Stream){
    col.push({
      title: 'id',
      key: 'id',

    })
  }

  if(type.value === RedisTypes.ZSet){
    col.push({
      title: 'score',
      key: 'score'
    })
  }

  if(type.value === RedisTypes.ZSet){
    col.push({
      title: 'member',
      key: 'member',
      ellipsis: {
        tooltip: true
      }
    })
  }else{
    col.push({
      title: 'value',
      key: 'value',
      ellipsis: {
        tooltip: true
      }
    })
  }



  col.push({
    title: 'Action',
    key: 'actions',
    render(row) {

      let actions = [
        h(NButton,
          {
            strong: true,
            ghost:true,
            size: 'small',
            type:"error",
            round:true,
            onClick: () => delField(row)
          },
          {default: () => 'del'})
      ]

      if(type.value !== RedisTypes.Stream){
        actions.unshift(h(NButton,
          {
            strong: false,
            size: 'small',
            round:true,
            onClick: () => editField(row),
            class:"mr-1"
          },
          {default: () => 'edit'}))
      }

      return h(
        "div",actions
      )
    }
  })

  return col
})



const pagination = computed(()=>{
  return {
    itemCount:len.value,
    prefix:()=>h("div","total: "+len.value),
    size:"small"
  }
})

const data = computed(() => {

  return list.value.map((item, index) => {
    return {
      no: index + 1,
      ...item
    }
  })
})

const message = useMessage()


function saveVal(val){
  set(props.curKey, type.value,{value:val}).then(()=>{
    message.success(
      "Save"
    )
  })
}

function changeTTL(v){
  let ttl = v || ttlVal.value
  if (ttl !== "-"){
    expire(props.curKey, ttl).then(()=>{
      message.success("Success")
      refresh(keyVal.value)
    })
  }else{
    persist(props.curKey).then(()=>{
      message.success("Success")
      refresh(keyVal.value)
    })
  }

}

function changeKey(){
  rename(oriKey, keyVal.value).then(()=>{
    message.success(
      "Success"
    )
    refresh(keyVal.value)
    scan()
  })
}

const info = useConn()

function refresh(k){

  load(k)
  oriKey = k
}

function delKey(){
  del({key:props.curKey}).then(()=>{
    load()
    scan()
    info.loadInfo("KeySpace")
    message.success(
      "Success"
    )
  })
}

const emit = defineEmits(["addField"])
function addField(){
  emit("addField", {
    key:props.curKey,
    type:type.value,
  })
}


defineExpose({
  load
})

</script>

<style scoped>

</style>
