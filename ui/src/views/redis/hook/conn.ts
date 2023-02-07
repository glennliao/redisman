import {parseInfo} from "~/utils/redis";
import {conn, executeCommand} from "~/api/redis";
import {useKeysHook} from './keys'

const info = ref<Record<any, any>>({})
const databases = ref(0)
const curDb = ref(0)

function loadInfo(section:string="default"){
  executeCommand([
    ["info",section],
    ["config", "get" ,"databases"]
  ]).then(data=>{
    let infoStr = data[0]
    info.value  = {
      ...toRaw(info.value),
      ...parseInfo(infoStr)
    }
    databases.value = parseInt(data[1].databases)
  })
}

const dbList = computed(()=>{
  let dbList = []

  let keyspace = toRaw(info.value.Keyspace || {})
  Object.keys(keyspace).forEach(k=>{
    let sps = keyspace[k].split(",")
    keyspace[k.substring(2)] = {
      keys:sps[0].substring("keys=".length)
    }
  })
  console.debug(keyspace)

  for (let i = 0; i < databases.value; i++) {
    dbList.push({
      label: i,
      cnt: (keyspace[i]||{keys:0}).keys,
      value: i,
      class:"db-select-option"
    })
  }
  return dbList
})


const {scan}  = useKeysHook()


function select(db:number){
  executeCommand([
    ["select",db+""],
  ]).then(data=>{
    if (data[0] == "OK"){
      curDb.value = db
      scan()
    }
  })
}


const connected = ref(false)
const connMeta = ref({
  title:"",
  dbAlias:{},
  id:0,
})

function connect(id:number){
  return conn({id}).then((data)=>{
    data.id = id
    connMeta.value = data
    connected.value = true
    select(curDb.value)
    scan()
  })
}


export function useInfo() {
  return {
    info, loadInfo,dbList,curDb,select,scan,connect,connected,connMeta
  }
}
