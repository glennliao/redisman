import { redis } from "~/api";
import { RedisTypes, redisTypes } from "~/views/redis/redis_types";

export function useValueHook() {
  const type = ref("");
  const ttl = ref("0");

  const value = ref("");
  const list = ref([]);
  const len = ref(0);

  function loadMeta(key: string) {
    return redis.command([
      ["type", key],
      ["ttl", key],
    ]).then((data) => {
      type.value = data[0];
      ttl.value = `${data[1]}`;
    });
  }

  function loadValue(key: string, index: string | number | undefined = undefined) {
    redis.command([
      ["get", key],
    ]).then((data) => {
      value.value = data[0];
    });
  }

  function loadList(key: string, pattern = "*", cursor = "0", count = "1000") {
    const commands = [] as string[][];

    const redisType = redisTypes[type.value];
    if (redisType) {
      commands.push(...redisType.list(key, pattern, cursor, count));
    } else {
      console.warn("不支持的redisTypes:", type);
      return Promise.reject(`不支持的redisTypes:${type}`);
    }

    redis.command(commands).then((data) => {
      len.value = data[0];

      let ret = data[1];
      if (commands[commands.length - 1][0].endsWith("scan")) {
        ret = ret[1];
        cursor = ret[0]; // todo
      }

      ret = redisType.listResultHandler(ret);
      list.value = ret;
      return ret;
    });
  }

  function load(key: string) {
    if (!key) {
      return;
    }
    loadMeta(key).then(() => {
      value.value = "";
      len.value = 0;
      list.value = [];
      if (type.value === RedisTypes.String) {
        loadValue(key);
      } else {
        loadList(key, "*", "0");
      }
    });
  }

  interface Data {
    value: string;member: string;score: string;id: string;field: string;
  }

  function set(key: string, type: string, data: Data, oldData: Data) {
    const commands = [];
    if (type === RedisTypes.String) {
      commands.push(["set", key, data.value]);
    } else {
      const redisType = redisTypes[type];
      if (redisType) {
        commands.push(...redisType.set(key, data, oldData));
      } else {
        console.warn("不支持的redisTypes:", type);
        return Promise.reject(`不支持的redisTypes:${type}`);
      }
    }

    return redis.command(commands).then(() => {
      load(key);
    });
  }

  function del({ key, type, data }: { key: string | string[];type?: string;data?: Record<string, any> }) {
    const commands = [];
    if (data) {
      const redisType = redisTypes[type];
      if (redisType) {
        commands.push(redisType.del(key as string, data));
      } else {
        console.warn("不支持的redisTypes:", type);
        return Promise.reject(`不支持的redisTypes:${type}`);
      }
    } else {
      if (Array.isArray(key)) {
        commands.push(["del"].concat(key));
      } else {
        commands.push(["del", key]);
      }
    }

    return redis.command(commands);
  }

  function expire(k: string, ttl: number) {
    return redis.command([
      ["expire", `${k}`, `${ttl}`],
    ]);
  }

  function persist(k: string, ttl: number) {
    return redis.command([
      ["persist", `${k}`],
    ]);
  }

  function rename(old: string, n: string) {
    return redis.command([
      ["rename", old, n],
    ]);
  }

  return {
    type,
    ttl,
    value,
    list,
    len,
    del,
    expire,persist,
    rename,
    load,
    set,
  };
}
