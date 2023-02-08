import { action, http } from "./service";

export const redis = {
  command: (commands: string[][]): Promise<any> => {
    return action({ action: "redisCom", params: commands });
  },
  connTest: (connConfig: any) => {
    return action({ action: "redisConnTest", params: connConfig });
  },
  conn: (params: { id: number }): Promise<any> => {
    return action({ action: "redisConn", params });
  },
};

export const apiJson = {
  get: (p: Record<string, any>) => http("/get", "POST", p),
  post: (p: Record<string, any>) => http("/post", "POST", p),
  put: (p: Record<string, any>) => http("/put", "POST", p),
  delete: (p: Record<string, any>) => http("/delete", "POST", p),
};
