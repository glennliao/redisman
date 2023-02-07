import {defaultWebSocketService} from "~/api/index";


// @ts-ignore
window.basicURL = "//localhost:16379"

export const executeCommand = (commands:string[][]): Promise<any> => {
  return defaultWebSocketService.request({ action: 'redisCom', params:commands })
}

export const testConn = (connConfig:any)=>{
  return defaultWebSocketService.request({action:"redisTest", params:connConfig})
}

export const conn = (params: { id:number }): Promise<any> => {
  return defaultWebSocketService.request({ action: 'redisConn', params })
}

function request(url:string, method:string, data:Record<string, any>){
  // @ts-ignore
  url = "http://"+window.basicURL+url
  return fetch(url,{credentials:'include', method:method, body:JSON.stringify(data)})
    .then(resp=>resp.json())
}

export const apijson = {
  get : (p:any)=> request("/get","POST", p),
  post : (p:any)=> request("/post","POST", p),
  put : (p:any)=> request("/put","POST", p),
  delete : (p:any)=> request("/delete","POST", p)
}
