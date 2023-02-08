export enum RedisTypes {
  String = "string",
  List = "list",
  Set = "set",
  ZSet = "zset",
  Hash = "hash",
  Stream = "stream",

  // Bitmap = "bitmap" // todo bitmap的支持

}

// @ts-expect-error
export const typesOptions = Object.keys(RedisTypes).map(k => RedisTypes[k]).map((item) => {
  return {
    label: item, key: item,
  };
});

interface TypeHandlers {
  del(key: string, field: Record<string, any>): string[];
  list(key: string, pattern: string, cursor: string, count: string): string[][];
  listResultHandler(ret: string[]): any[];

  set(key: string, data: Record<any, any>, oldData?: Record<any, any>): string[][];

}

export const redisTypes: Record<string, TypeHandlers> = {};

redisTypes[RedisTypes.Hash] = {
  set(key: string, data: Record<any, any>, oldData: Record<any, any> = {}): string[][] {
    const command: string[][] = [];
    command.push(["hset", key, data.field, data.value]);
    if (data.field !== oldData.field && oldData.field) {
      command.push(this.del(key, oldData));
    }
    return command;
  },
  del(key: string, data: Record<string, any>): string[] {
    return ["hdel", key, data.field];
  },

  list(key: string, pattern = "*", cursor = "0", count: string) {
    return [
      ["hlen", key],
      ["hscan", key, cursor, "MATCH", pattern, "COUNT", count],
    ];
  },
  listResultHandler(ret: any) {
    const retList = [];
    for (let i = 0; i < ret.length; i += 2) {
      retList.push({
        key: ret[i],
        field: ret[i],
        value: ret[i + 1],
      });
    }
    return retList;
  },

};

redisTypes[RedisTypes.List] = {
  set(key: string, data: Record<any, any>, oldData: Record<any, any> = {}): string[][] {
    const command: string[][] = [];
    command.push(["lpush", key, data.value]);
    if (data.value !== oldData.value && oldData.value) {
      command.push(this.del(key, oldData));
    }
    return command;
  },
  del(key: string, data: Record<string, any>): string[] {
    return ["lrem", key, "1", data.value];
  },
  list(key: string, pattern = "*", cursor = "0", count: string) {
    return [
      ["llen", key],
      ["lrange", key, cursor, count],
    ];
  },
  listResultHandler(list: string[]) {
    return list.map((item) => {
      return {
        value: item,
      };
    });
  },

};

redisTypes[RedisTypes.Set] = {
  set(key: string, data: Record<any, any>, oldData: Record<any, any> = {}): string[][] {
    const command: string[][] = [];
    command.push(["sadd", key, data.value]);
    if (data.value !== oldData.value && oldData.value) {
      command.push(this.del(key, oldData));
    }
    return command;
  },
  del(key: string, data: Record<string, any>): string[] {
    return ["srem", key, data.value];
  },
  list(key: string, pattern = "*", cursor = "0", count: string) {
    return [
      ["scard", key],
      ["sscan", key, cursor, "MATCH", pattern, "COUNT", count],
    ];
  },
  listResultHandler(list: string[]) {
    return list.map((item) => {
      return {
        value: item,
      };
    });
  },

};

redisTypes[RedisTypes.ZSet] = {
  set(key: string, data: Record<any, any>, oldData: Record<any, any> = {}): string[][] {
    const command: string[][] = [];
    command.push(["zadd", key, data.member]);
    if (data.member !== oldData.member && oldData.member) {
      command.push(this.del(key, oldData));
    }
    return command;
  },
  del(key: string, data: Record<string, any>): string[] {
    return ["zrem", key, data.member];
  },
  list(key: string, pattern = "*", cursor = "0", count: string) {
    return [
      ["zcard", key],
      ["zrevrange", key, cursor, count, "WITHSCORES"],
    ];
  },
  listResultHandler(ret: any) {
    const retList = [];
    for (let i = 0; i < ret.length; i += 2) {
      retList.push({
        score: ret[i + 1],
        member: ret[i],
      });
    }
    return retList;
  },

};

redisTypes[RedisTypes.Stream] = {
  set(key: string, data: Record<any, any>, oldData: Record<any, any> = {}): string[][] {
    const command: string[][] = [];
    command.push(["xadd", "*", data.field, data.value]);
    return command;
  },
  del(key: string, data: Record<string, any>): string[] {
    return ["xdel", key, data.id];
  },
  list(key: string, pattern = "*", cursor = "0", count: string) {
    return [
      ["xlen", key],
      ["xrevrange", key, "+", "-", "COUNT", "200"],
    ];
  },
  listResultHandler(ret: any) {
    const retList = [];

    for (let i = 0; i < ret.length; i++) {
      const values = ret[i][1];
      const value = {} as any;
      for (let i = 0; i < values.length; i += 2) {
        value[values[i]] = values[i + 1];
      }
      retList.push({
        id: ret[i][0],
        value: JSON.stringify(value),
      });
    }
    return retList;
  },

};
