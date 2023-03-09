import { useKeysHook } from "./keys";
import { parseInfo } from "~/utils/redis";
import { redis } from "~/api";

const info = ref<Record<any, any>>({
  Clients:{},
  Memory:{}
});
const databases = ref(0);
const curDb = ref(0);

function loadInfo(section = "default") {
  redis.command([
    ["info", section],
    ["config", "get", "databases"],
  ]).then((data) => {
    const infoStr = data[0];
    info.value = {
      ...toRaw(info.value),
      ...parseInfo(infoStr),
    };

    let config = data[1]
    if(Array.isArray(config)){
      databases.value = parseInt(data[1][1])
    }else{
      databases.value = parseInt(data[1].databases);
    }

  });
}

const dbList = computed(() => {
  const dbList = [];

  const keyspace = toRaw(info.value.Keyspace || {});
  Object.keys(keyspace).forEach((k) => {
    if(k.startsWith("db")){
      const sps = keyspace[k].split(",");
      keyspace[k.substring(2)] = {
        keys: sps[0].substring("keys=".length),
      };
    }
  });

  for (let i = 0; i < databases.value; i++) {
    dbList.push({
      label: i,
      cnt: (keyspace[i] || { keys: 0 }).keys,
      value: i,
      class: "db-select-option",
    });
  }
  return dbList;
});

const { scan } = useKeysHook();

function select(db: number) {
  redis.command([
    ["select", `${db}`],
  ]).then((data) => {
    if (data[0] == "OK") {
      curDb.value = db;
      scan();
    }
  });
}

const connected = ref(false);
const connMeta = ref({
  title: "",
  dbAlias: {},
  id: 0,
});

function status(){
  redis.command([
    ["info", "Memory"],
    ["info", "Clients"],
  ]).then((data) => {
    info.value = {
      ...toRaw(info.value),
      ...parseInfo(data[0]),
      ...parseInfo(data[1]),
    };
  });
}

let statusTimer:any = null

function connect(id: number) {
  return redis.conn({ id }).then((data) => {
    data.id = id;
    connMeta.value = data;
    connected.value = true;
    select(curDb.value);
    status()
    if(statusTimer !== null){
      clearInterval(statusTimer)
      statusTimer = null
    }
    statusTimer = setInterval(()=>{
      status()
    }, 30 * 1000)
  }).catch((err)=>{
    if(statusTimer !== null){
      clearInterval(statusTimer)
      statusTimer = null
    }
    throw err
  });
}

export function useConn() {
  return {
    info, loadInfo, dbList, curDb, select, scan, connect, connected, connMeta,
  };
}
