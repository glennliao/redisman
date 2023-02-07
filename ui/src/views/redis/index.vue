<script setup lang="ts">
import {Add} from '@vicons/ionicons5'
import {apijson} from "~/api/redis";
import {useInfo} from "~/views/redis/hook/conn";
import AddConnectionModal from "~/views/redis/components/AddConnectionModal.vue";
const {connect} = useInfo()

const connectionListRef = ref([])


const message = useMessage()
const dialog = useDialog()
function loadList(){
  apijson.get({
    "RedisConnection[]":{

    }
  }).then(data=>{
    connectionListRef.value = data["RedisConnection[]"]
  })
}
loadList()
const router = useRouter()


function connRedis(id:number){
  showLoadingRef.value = true
  connect(id).then(()=>{
    router.push("/redis/conn?id="+id)
  }).catch(err=>{
    console.log(err)
    dialog.error({
      title: 'err',
      content: err.msg,
      positiveText: 'Ok'
    })
  }).finally(()=>{
    showLoadingRef.value = false
  })
}

const addConnectionModalRef = ref(null) as any

function success(e:any){
  console.log(e)
  loadList()
}

function add(type: string) {

  addConnectionModalRef.value && addConnectionModalRef.value.open({type})
}


const showLoadingRef = ref(false)

function delConn(id:number){
  console.log("del",id)
  apijson.delete({
    tag:"RedisConnection",
    RedisConnection:{
      id:id+""
    }
  }).then(()=>{
    loadList()
  })
}

function updateConn(id:number){
  console.log('update conn',id)
}


</script>

<template>

  <n-spin :show="showLoadingRef">
  <n-layout style="height: calc(100vh - 64px - 8px)" class="bg-base-100">

    <div class="flex flex-wrap p-2">
      <div v-for="conn in connectionListRef" :key="conn.id" class="conn-item cursor-pointer">
        <n-card size="small" embedded :title="conn.title" hoverable    >
          <div class="conn-content" @click="connRedis(conn.id)">
            <div>
              {{conn.host}}:{{conn.port}}
            </div>
            <div>
              {{conn.username}} - {{conn.password}}
            </div>

            <div>
              {{conn.createdAt}}
            </div>
          </div>
          <template #action>
            <n-button-group>
<!--              <n-button ghost round @click="updateConn(conn.id)">-->
<!--&lt;!&ndash;                <template #icon>&ndash;&gt;-->
<!--&lt;!&ndash;                  <n-icon><log-in-icon /></n-icon>&ndash;&gt;-->
<!--&lt;!&ndash;                </template>&ndash;&gt;-->
<!--                Edit-->
<!--              </n-button>-->

              <n-button round @click="delConn(conn.id)">
<!--                <template #icon>-->
<!--                  <n-icon><log-in-icon /></n-icon>-->
<!--                </template>-->
                Del
              </n-button>
            </n-button-group>
          </template>
        </n-card>
      </div>
      <div class="conn-item cursor-pointer">
        <n-card embedded hoverable size="small"  @click="add" style="width: 160px">
          <div class="conn-content text-center">
            <n-icon size="32">
              <Add />
            </n-icon>
          </div>
        </n-card>
      </div>

    </div>
  </n-layout>
  <n-layout-footer>
    <div class="footer footer-center bottom-0 border-t border-base-100 bg-base-200 px-4 py-1 text-base-content opacity-90" style="height: 32px" >
      RedisMan v0.1.0
    </div>
  </n-layout-footer>
  <AddConnectionModal ref="addConnectionModalRef" @success="success"/>
  </n-spin>
</template>

<style scoped>
.conn-item{

  margin: 6px;
  padding: 6px;
}

.conn-content{
  /*width: 160px;*/
  /*height: 160px;*/
}

</style>
