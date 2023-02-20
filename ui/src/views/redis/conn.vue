<script setup lang="ts">
import type { DropdownOption, SelectOption, TreeOption } from "naive-ui";
import { Add, Refresh, Search, Settings, SettingsOutline,TerminalOutline } from "@vicons/ionicons5";
import type { Ref } from "vue";
import { useKeysHook } from "./hook/keys";

import { typesOptions } from "~/views/redis/redis_types";

import Key from "~/views/redis/key.vue";
import { useConn } from "./hook/conn";
import { useValueHook } from "./hook/value";
import AddValueModal from "./components/AddValueModal.vue";
import RedisConnSettingModal from "./components/RedisConnSettingModal.vue";
import RedisInfoModal from "./components/RedisConnInfoModal.vue";
import RedisCliModal from "./components/RedisCliModal.vue";

const route = useRoute();

const connId = route.query.id;

const { treeKeys, scanKeys } = useKeysHook();

const curKey: Ref<string> = ref("");

const { info, dbList, loadInfo, select, curDb, scan, connect, connected, connMeta } = useConn();

if (connected.value) {
  loadInfo("Keyspace");
} else {
  connect(connId as unknown as number).then(() => {
    loadInfo("Keyspace");
  });
}

const {
  type, ttl, load, value, len, list, set, del, expire, rename,
} = useValueHook();

const addValueModalRef = ref(null) as any;

function handleAddSelect(type: string) {
  addValueModalRef.value && addValueModalRef.value.open({ type });
}

function addField({ key, type, data }: { key: string; type: string; data?: Record<any, any> }) {
  addValueModalRef.value && addValueModalRef.value.open({ type, key, data });
}

function success(e: { key: string }) {
  keyRef.value && keyRef.value.load(e.key);
  loadInfo("Keyspace"); // 需统一到refresh
}

const keyRef = ref();

function renderDbLabel(option: SelectOption) {
  const leftChild = [
    h("span", option.label),
  ];
  // @ts-expect-error
  if (option.cnt > 0) {
    leftChild.push(h("span", {
      class: "ml-1",
      style: {
        fontSize: "12px",
        color: "gray",
      },
    }, `(${option.cnt})`));
  }
  const left = h("div", {}, leftChild);

  const dbAlias = connMeta.value.dbAlias || {};

  const right = h("div", {
    class: "ml-1",
    style: {
      fontSize: "12px",
      color: "gray",
    },
  }, dbAlias[option.label]);

  return h("div", {
    style: {
      "display": "flex",
      "justify-content": "space-between",
    },
  }, [
    left, right,
  ]);
}

function handleUpdateValue(value: string, option: SelectOption) {
  select(parseInt(value));
}

function onPatternChange(val: string) {
  if(!exactMatch.value){
    val = val.trim()+"*"
  }
  scan(val);
}

const pattern = ref("");

function refreshKey() {
  pattern.value = "";
  scan();
  loadInfo("Keyspace");
}

const showContextmenuRef = ref(false);
const optionsRef = ref<DropdownOption[]>([]);
const xRef = ref(0);
const yRef = ref(0);

async function handleSelect(e: string) {
  const actions: { [k: string]: () => void } = {
    del: async () => {
      let delKeys: string[] = [];

      if (multipleSelectRef.value) {
        for (const item of checkKeysWithMeta) {
          if (item.leaf) {
            delKeys.push(item.key);
          } else {
            const keys = await scanKeys(`${item.key}*`);
            delKeys = delKeys.concat(keys);
          }
        }
      } else {

        let k = contextMenuSelectedKey.value
        k = k.substring(4)
        if (contextMenuSelectedKeyIsLeaf.value) {
          delKeys = [k];
        } else {
          delKeys = await scanKeys(`${k}*`);

        }
      }

      if (delKeys.length) {
        del({
          key: delKeys,
        }).then(() => {
          refreshKey();
        });
      }

      console.debug(contextMenuSelectedKey.value, contextMenuSelectedKeyIsLeaf.value, multipleSelectRef.value);
      console.debug(e, checkedKeys.value);
    },
    multipleDel: async () => {
      console.debug(contextMenuSelectedKey.value, contextMenuSelectedKeyIsLeaf.value, multipleSelectRef.value);
      console.debug(e, checkedKeys.value);
    },
    select: async () => {
      multipleSelectRef.value = true;
    },
    cancel: async () => {
      multipleSelectRef.value = false;
    },
  };

  if (actions[e]) {
    await actions[e]();
  }

  showContextmenuRef.value = false;
}

function handleClickoutside() {
  showContextmenuRef.value = false;
}

const contextMenuSelectedKey = ref("");
const contextMenuSelectedKeyIsLeaf = ref(false);

const multipleSelectRef = ref(false);
const checkedKeys = ref([]);

const nodeProps = ({ option }: { option: TreeOption }) => {
  return {

    onClick() {
      // console.log("click",option)
      if (option.leaf) {
        curKey.value = option.key as string;
      }
    },
    onContextmenu(e: MouseEvent): void {
      // console.log("right click",option)
      optionsRef.value = [] as any;

      if (multipleSelectRef.value) {
        optionsRef.value.push({
          label: "Delete Select",
          key: "del",
        } as any);
        optionsRef.value.push({
          type: "divider",
        } as any);
        optionsRef.value.push({
          label: "Cancel",
          key: "cancel",
        } as any);
      } else {
        optionsRef.value.push({
          label: option.key,
          disabled: true,
        } as any);
        optionsRef.value.push({
          type: "divider",
        } as any);
        optionsRef.value.push({
          label: "Delete",
          key: "del",
        } as any);
        optionsRef.value.push({
          label: "Multiple Select",
          key: "select",
        } as any);
      }

      showContextmenuRef.value = true;
      xRef.value = e.clientX;
      yRef.value = e.clientY;
      e.preventDefault();

      contextMenuSelectedKey.value = option.key as string;
      contextMenuSelectedKeyIsLeaf.value = !!option.leaf as boolean;
    },
  };
};

let checkKeysWithMeta: any[] = [];

function onUpdateCheckedKeys(keys: any, keysWithMeta: any[]) {
  console.log(keys, keysWithMeta);
  checkedKeys.value = keys;
  checkKeysWithMeta = keysWithMeta;
}

function keysRenderLabel({ option }: { option: TreeOption }) {
  const child = [
    h("span", {}, option.label),
  ];
  if (option.children) {
    child.push(h("span", {
      style: {
        marginLeft: "8px",
        color: "gray",
        fontSize: "12px",
      },
    }, `(${option.childrenCnt})`));
  }
  return h("div", {}, child);
}

const redisConnSettingModalRef = ref(null);
const redisInfoModalRef = ref(null);
const redisCliModalRef = ref(null);
function setting() {
  redisConnSettingModalRef.value && redisConnSettingModalRef.value.open();
}
function redisInfo() {
  redisInfoModalRef.value && redisInfoModalRef.value.open();
}


function cli(){
  redisCliModalRef.value && redisCliModalRef.value.open();
}

function redisConnSettingSuccess() {

}


const exactMatch = ref(false)
</script>

<template>
  <n-layout style="height: calc(100vh - 40px)" class="bg-base-100" bordered>
    <n-layout position="absolute" style="padding-top: 4px;top: 2px; bottom: 32px" has-sider>
      <n-layout-sider
        content-style="padding: 0 4px 24px 4px;"
        :native-scrollbar="false"
        bordered
      >
        <div class="flex">
          <n-select
            size="small" :value="curDb" :options="dbList" :render-label="renderDbLabel"
            @update:value="handleUpdateValue"
          />

          <n-dropdown trigger="click" :options="typesOptions" @select="handleAddSelect">
            <n-button size="small" class="ml-1">
              <template #icon>
                <n-icon>
                  <Add />
                </n-icon>
              </template>
            </n-button>
          </n-dropdown>
        </div>

        <div class="mt-1 flex">
          <n-input v-model:value="pattern" placeholder="*" size="small" @change="onPatternChange" >
            <template #suffix>
              <n-checkbox title="exactMatch" v-model:checked="exactMatch"/>
            </template>
          </n-input>

          <n-button size="small" class="ml-1" @click="onPatternChange(pattern)">
            <template #icon>
              <n-icon>
                <Search />
              </n-icon>
            </template>
          </n-button>
          <n-button size="small" class="ml-1" @click="refreshKey">
            <template #icon>
              <n-icon>
                <Refresh />
              </n-icon>
            </template>
          </n-button>
        </div>

        <div style="height: 6px" />


        <div class="keys">
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
      <n-layout content-style="padding:0 12px;" :native-scrollbar="false">
        <Key v-if="curKey" ref="keyRef" :cur-key="curKey" @addField="addField" />
        <div v-else>
          please select a key from left
        </div>
        <AddValueModal ref="addValueModalRef" @success="success" />
      </n-layout>
    </n-layout>
    <n-layout-footer
      position="absolute"
      style="height: 32px; padding: 4px 12px"
      bordered
    >
      <n-grid cols="24">
        <n-gi :span="4">
          <text>
            {{ connMeta.title }}
          </text>
          <n-icon class="cursor-pointer align-middle  ml-1" :size="14" @click="setting">
            <SettingsOutline />
          </n-icon>
        </n-gi>
        <n-gi :span="4">
          <div class="cursor-pointer" @click="redisInfo">
            <n-divider vertical />
            <text title="used_memory_human">MEM: {{info.Memory.used_memory_human}}</text>
            <text class="ml-3" title="connected_clients">CN: {{info.Clients.connected_clients}}</text>
          </div>
        </n-gi>
        <n-gi :span="1"></n-gi>
        <n-gi :span="4">
          <div class="flex">
            <n-divider vertical />
            <div class="cursor-pointer" @click="cli">
              <n-icon class=" align-middle  ml-1" :size="14">
                <TerminalOutline />
              </n-icon>
              Terminal
            </div>
          </div>
        </n-gi>

      </n-grid>

    </n-layout-footer>
  </n-layout>

  <RedisConnSettingModal ref="redisConnSettingModalRef" @success="redisConnSettingSuccess" />
  <RedisInfoModal ref="redisInfoModalRef"/>
  <RedisCliModal ref="redisCliModalRef"/>
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
