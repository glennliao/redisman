<script setup lang="ts">
import Key from "~/views/redis/key.vue";
import type {DropdownOption, SelectOption, TreeOption} from 'naive-ui'
import {Add, Refresh, Search,Settings,SettingsOutline} from '@vicons/ionicons5'
import {useKeysHook} from './hook/keys'

const route = useRoute()

const connId = route.query.id


import {typesOptions} from "~/views/redis/redis_types";

import {Ref} from "vue";
import {useInfo} from "~/views/redis/hook/conn";
import AddValueModal from "~/views/redis/components/AddValueModal.vue";
import {useValueHook} from "~/views/redis/hook/value";
import RedisConnSettingModal from "~/views/redis/components/RedisConnSettingModal.vue";

const {treeKeys, scanKeys} = useKeysHook()

const curKey: Ref<string> = ref('')

const {info, dbList, loadInfo, select, curDb, scan, connect, connected, connMeta} = useInfo()

if (connected.value) {
  loadInfo()
} else {
  connect(connId as unknown as number).then(() => {
    loadInfo()
  })
}


const {
  type, ttl, load, value, len, list, set, del, expire, rename
} = useValueHook()

const addValueModalRef = ref(null) as any

function handleAddSelect(type: string) {

  addValueModalRef.value && addValueModalRef.value.open({type})
}


function addField({key, type, data}: { key: string, type: string, data?: Record<any, any> }) {
  addValueModalRef.value && addValueModalRef.value.open({type, key, data})
}

function success(e: { key: string; }) {
  keyRef.value && keyRef.value.load(e.key)
  loadInfo("Keyspace") // 需统一到refresh
}

const keyRef = ref()


function renderDbLabel(option: SelectOption) {

  let leftChild = [
    h("span", option.label)
  ]
  // @ts-ignore
  if (option.cnt > 0) {
    leftChild.push(h("span", {
      class: "ml-1",
      style: {
        fontSize: '12px',
        color: 'gray'
      }
    }, `(${option.cnt})`),)
  }
  let left = h("div", {}, leftChild)


  let dbAlias = connMeta.value.dbAlias || {}

  let right = h("div", {
    class: "ml-1",
    style: {
      fontSize: '12px',
      color: 'gray',

    }
  }, dbAlias[option.label])


  return h('div', {
    style: {
      display: "flex",
      "justify-content": "space-between"
    }
  }, [
    left, right
  ])
}

function handleUpdateValue(value: string, option: SelectOption) {
  select(parseInt(value))
}


function onPatternChange(val: string) {
  scan(val)
}

const pattern = ref('')

function refreshKey() {
  pattern.value = ""
  scan()
  loadInfo("Keyspace")
}


const showContextmenuRef = ref(false)
const optionsRef = ref<DropdownOption[]>([])
const xRef = ref(0)
const yRef = ref(0)

async function handleSelect(e: string) {

  let actions: { [k: string]: () => void } = {
    "del": async () => {

      let delKeys: string[] = []

      if (multipleSelectRef.value) {

        for (const item of checkKeysWithMeta) {
          if (item.leaf) {
            delKeys.push(item.key)
          } else {
            let keys = await scanKeys(item.key + "*")
            delKeys = delKeys.concat(keys)
          }
        }


      } else {
        if (contextMenuSelectedKeyIsLeaf.value) {
          delKeys = [contextMenuSelectedKey.value]
        } else {
          delKeys = await scanKeys(contextMenuSelectedKey.value + "*")
          console.log(delKeys)
        }
      }
      console.log(delKeys)
      if (delKeys.length) {
        del({
          key: delKeys
        }).then(() => {
          refreshKey()
        })
      }

      console.log(contextMenuSelectedKey.value, contextMenuSelectedKeyIsLeaf.value, multipleSelectRef.value)
      console.log(e, checkedKeys.value)
    },
    "multipleDel": async () => {

      console.log(contextMenuSelectedKey.value, contextMenuSelectedKeyIsLeaf.value, multipleSelectRef.value)
      console.log(e, checkedKeys.value)
    },
    "select": async () => {
      multipleSelectRef.value = true
    },
    "cancel": async () => {
      multipleSelectRef.value = false
    }
  }

  if (actions[e]) {
    await actions[e]()
  }

  showContextmenuRef.value = false
}

function handleClickoutside() {
  showContextmenuRef.value = false
}

const contextMenuSelectedKey = ref('')
const contextMenuSelectedKeyIsLeaf = ref(false)

const multipleSelectRef = ref(false)
const checkedKeys = ref([])

const nodeProps = ({option}: { option: TreeOption }) => {
  return {

    onClick() {
      // console.log("click",option)
      if (option.leaf) {
        curKey.value = option.key as string
      }
    },
    onContextmenu(e: MouseEvent): void {
      // console.log("right click",option)
      optionsRef.value = [] as any

      if (multipleSelectRef.value) {
        optionsRef.value.push({
          label: "Delete Select",
          key: "del"
        } as any)
        optionsRef.value.push({
          type: "divider"
        } as any)
        optionsRef.value.push({
          label: "Cancel",
          key: "cancel"
        } as any)
      } else {
        optionsRef.value.push({
          label: option.key,
          disabled: true,
        } as any)
        optionsRef.value.push({
          type: "divider"
        } as any)
        optionsRef.value.push({
          label: "Delete",
          key: "del"
        } as any)
        optionsRef.value.push({
          label: "Multiple Select",
          key: "select"
        } as any)
      }


      showContextmenuRef.value = true
      xRef.value = e.clientX
      yRef.value = e.clientY
      e.preventDefault()

      contextMenuSelectedKey.value = option.key as string
      contextMenuSelectedKeyIsLeaf.value = !!option.leaf as boolean
    }
  }
}

let checkKeysWithMeta: any[] = []

function onUpdateCheckedKeys(keys: any, keysWithMeta: any[]) {
  console.log(keys, keysWithMeta)
  checkedKeys.value = keys
  checkKeysWithMeta = keysWithMeta
}

function keysRenderLabel({option}: { option: TreeOption }) {
  let child = [
    h("span", {}, option.label)
  ]
  if (option.children) {
    child.push(h("span", {
      style: {
        marginLeft: '8px',
        color: 'gray',
        fontSize: "12px"
      }
    }, `(${option.childrenCnt})`))
  }
  return h("div", {}, child)
}



const redisConnSettingModalRef = ref(null)
function setting(){
  redisConnSettingModalRef.value && redisConnSettingModalRef.value.open()
}
function redisConnSettingSuccess(){

}

</script>

<template>
  <n-layout style="height: calc(100vh - 40px)" class="bg-base-100" bordered>
<!--    <n-layout-header style="height: 2px;" bordered>-->
<!--    </n-layout-header>-->
    <n-layout position="absolute" style="padding-top: 4px;top: 2px; bottom: 32px" has-sider>

      <n-layout-sider
        content-style="padding: 0 4px 24px 4px;"
        :native-scrollbar="false"
        bordered
      >

        <div class="flex">
          <n-select size="small" :value="curDb" :options="dbList" :render-label="renderDbLabel"
                    @update:value="handleUpdateValue"></n-select>



          <n-dropdown trigger="click" :options="typesOptions" @select="handleAddSelect">
            <n-button size="small" class="ml-1">
              <template #icon>
                <n-icon>
                  <Add/>
                </n-icon>
              </template>
            </n-button>
          </n-dropdown>
        </div>

        <div class="flex mt-1">
          <n-input placeholder="*" size="small" v-model:value="pattern" @change="onPatternChange"></n-input>

          <n-button size="small" class="ml-1" @click="onPatternChange(pattern)">
            <template #icon>
              <n-icon>
                <Search/>
              </n-icon>
            </template>
          </n-button>
          <n-button size="small" class="ml-1" @click="refreshKey">
            <template #icon>
              <n-icon>
                <Refresh/>
              </n-icon>
            </template>
          </n-button>
        </div>
        <div style="height: 12px"></div>
       <div  class="keys">
         <n-tree

           virtual-scroll
           block-line
           :render-label="keysRenderLabel"
           :data="treeKeys"
           expand-on-click
           style="height: calc(100vh - 64px - 32px - 124px)"
           selectable
           :node-props="nodeProps"
           :checkable="multipleSelectRef"
           :checked-keys="checkedKeys"
           check-on-click
           @update:checked-keys="onUpdateCheckedKeys"
         />
       </div>
        <n-dropdown
          trigger="manual"
          placement="bottom-start"
          :show="showContextmenuRef"
          :options="optionsRef"
          :x="xRef"
          :y="yRef"
          @select="handleSelect"
          @clickoutside="handleClickoutside"
        />

      </n-layout-sider>
      <n-layout content-style="padding:0 24px;" :native-scrollbar="false">
        <Key v-if="curKey" :curKey="curKey" ref="keyRef" @addField="addField"/>
        <div v-else>
          please select a key from left
        </div>
        <AddValueModal ref="addValueModalRef" @success="success"/>
      </n-layout>
    </n-layout>
    <n-layout-footer
      position="absolute"
      style="height: 32px; padding: 4px 12px"
      bordered
    >
      <div>
        {{ connMeta.title }} <n-icon @click="setting" class="pt-0.5 cursor-pointer" :size="12"> <SettingsOutline/></n-icon>
      </div>
    </n-layout-footer>
  </n-layout>

  <RedisConnSettingModal ref="redisConnSettingModalRef" @success="redisConnSettingSuccess"/>



</template>

<style>

.db-select-option .n-base-select-option__content {
  width: 100% !important;
}


.keys .n-tree .n-tree-node-wrapper {
  padding: 0;
}
.keys .n-tree .n-tree-node{
  padding: 2px 0;
}
.keys .n-tree .n-tree-node-switcher{
  width: 10px;
}
</style>
