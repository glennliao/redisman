<script setup lang="ts">
import {AddOutline, TrashOutline,CreateOutline} from "@vicons/ionicons5";
import {apiJson} from "~/api";
import {useConn} from "~/views/redis/hook/conn";
import AddConnectionModal from "~/views/redis/components/AddConnectionModal.vue";

const {connect} = useConn();

const connectionListRef = ref([]);

const message = useMessage();
const dialog = useDialog();

function loadList() {
  apiJson.get({
    "RedisConnection[]": {},
  }).then((data) => {
    connectionListRef.value = data["RedisConnection[]"];
    console.log(data, connectionListRef.value)
  });
}

loadList();
const router = useRouter();

function connRedis(id: number) {
  showLoadingRef.value = true;
  connect(id).then(() => {
    router.push(`/redis/conn?id=${id}`);
  }).catch((err) => {
    console.log(err);
    dialog.error({
      title: "err",
      content: err.msg,
      positiveText: "Ok",
    });
  }).finally(() => {
    showLoadingRef.value = false;
  });
}

const addConnectionModalRef = ref(null) as any;

function success(e: any) {
  console.log(e);
  loadList();
}

function add() {
  addConnectionModalRef.value && addConnectionModalRef.value.open({});
}

const showLoadingRef = ref(false);

function delConn(id: number) {
  console.log("del", id);
  apiJson.delete({
    tag: "RedisConnection",
    RedisConnection: {
      id: `${id}`,
    },
  }).then(() => {
    loadList();
  });
}

function updateConn(id: number) {
  addConnectionModalRef.value && addConnectionModalRef.value.open({id});
}
</script>

<template>
  <n-spin :show="showLoadingRef">
    <n-layout style="height: calc(100vh - 64px - 16px)" class="bg-base-100">
      <div class="flex flex-wrap p-2">
        <div v-for="conn in connectionListRef" :key="conn.id" class="conn-item cursor-pointer">
          <n-card size="small" embedded :title="conn.title" hoverable>
            <div class="conn-content" @click="connRedis(conn.id)">
              <div>
                {{ conn.host }}:{{ conn.port }}
              </div>
              <div>
                {{ conn.username }}
              </div>
              <div>
                {{ conn.createdAt }}
              </div>
            </div>
            <template #header>
              <div  @click="connRedis(conn.id)">
                {{conn.title}}
              </div>
            </template>
            <template #action>
              <div class="text-center">
                <n-button-group>
                  <n-button size="small" round @click="updateConn(conn.id)">
                    <template #icon>
                      <n-icon><CreateOutline /></n-icon>
                    </template>
                  </n-button>
                  <n-popconfirm
                    @positive-click="delConn(conn.id)"
                    @negative-click=""
                  >
                    <template #trigger>
                      <n-button size="small" round >
                        <template #icon>
                          <n-icon >
                            <TrashOutline/>
                          </n-icon>
                        </template>
                      </n-button>
                    </template>
                    Del ?
                  </n-popconfirm>
                </n-button-group>
              </div>
            </template>
          </n-card>
        </div>
        <div class="conn-item cursor-pointer">
          <n-card embedded hoverable size="small" style="width: 160px" @click="add">
            <div class="conn-content text-center">
              <n-icon size="32">
                <AddOutline/>
              </n-icon>
            </div>
          </n-card>
        </div>
      </div>
    </n-layout>
    <n-layout-footer>
      <div
        class="footer footer-center bottom-0 border-t border-base-100 bg-base-200 px-4 py-1 text-base-content opacity-90"
        style="height: 32px;width:unset">
        RedisMan v0.2.0
      </div>
    </n-layout-footer>
    <AddConnectionModal ref="addConnectionModalRef" @success="success"/>
  </n-spin>
</template>

<style scoped>
.conn-item {

  margin: 6px;
  padding: 6px;
}

.conn-content {
  /*width: 160px;*/
  /*height: 160px;*/
}
</style>
