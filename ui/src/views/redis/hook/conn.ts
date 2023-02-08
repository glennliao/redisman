import { useKeysHook } from "./keys";
import { parseInfo } from "~/utils/redis";
import { redis } from "~/api";

const info = ref<Record<any, any>>({});
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
    databases.value = parseInt(data[1].databases);
  });
}

const dbList = computed(() => {
  const dbList = [];

  const keyspace = toRaw(info.value.Keyspace || {});
  Object.keys(keyspace).forEach((k) => {
    const sps = keyspace[k].split(",");
    keyspace[k.substring(2)] = {
      keys: sps[0].substring("keys=".length),
    };
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

function connect(id: number) {
  return redis.conn({ id }).then((data) => {
    data.id = id;
    connMeta.value = data;
    connected.value = true;
    select(curDb.value);

  });
}

export function useConn() {
  return {
    info, loadInfo, dbList, curDb, select, scan, connect, connected, connMeta,
  };
}
