import { request } from "./http";
import type { Request } from "./ws";
import { WebSocketService } from "./ws";

import {
  createDiscreteApi,
  ConfigProviderProps,
  darkTheme,
  lightTheme
} from 'naive-ui'

const {  dialog } = createDiscreteApi(
  [ 'dialog'],
  {}
)


const isDesktop = "go" in window;

function appHttp(url: string, method: string, data: Record<string, any>) {
  return window.go.main.App.Http({ method, url, data });
}

export const http = function (url: string, method: string, data: Record<string, any>){
  return (isDesktop ? appHttp : request)(url,method,data).then((data:{code:Number;data:Record<any, any>;msg:string})=>{
    if (data.code === 200){
      if(!data.data){
        return data // apijson 暂与code同级
      }
      return data.data
    }
    if(data.code === 500){
      //@ts-ignore

      dialog.error({
        title: 'api err : '+`[${method}] `+url,
        content: data.msg,
        positiveText: '啊',
      })
    }
    throw data
  })
};

function appAction(req: Request): Promise<any> {
  return window.go.main.App.Action(JSON.stringify(req)).then(msg=>{
    let resp = JSON.parse(msg)
    if (resp.code === 200){
      return resp.data
    }
    throw  resp
  });
}

let _action = appAction;

if (!isDesktop) {
  window.basicURL = process.env.NODE_ENV === "development" ? "localhost:16379" : `${window.location.host}`;
  const ws = new WebSocketService();
  ws.init(`ws://${window.basicURL}/ws`);
  _action = ws.request;
}

export const action = function (req: Request): Promise<any> {
  let beginAt = Date.now()
  return _action(req).then(data=>{
    let span = Date.now() - beginAt
    console.log(`[${span}]`, req.action, JSON.stringify(req.params), data)
    return data
  })
};
