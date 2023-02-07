import {executeCommand} from "~/api/redis";
import {RedisTypes, redisTypes} from "~/views/redis/redis_types";


export function useValueHook() {
  const type = ref('')
  const ttl = ref("0")

  const value = ref('')
  const list = ref([])
  const len = ref(0)

  function loadMeta(key: string) {
    return executeCommand([
      ["type", key],
      ["ttl", key],
    ]).then(data => {
      type.value = data[0]
      ttl.value = data[1] + ""
    })
  }

  function loadValue(key: string, index: string | number | undefined = undefined) {
    executeCommand([
      ["get", key]
    ]).then(data => {
      value.value = data[0]
    })
  }

  function loadList(key:string,pattern:string = "*",cursor:string="0", count = "1000"){
    let commands = [] as string[][]

    let redisType = redisTypes[type.value]
    if(redisType){
      commands.push(...redisType.list(key,pattern,cursor,count))
    }else{
      console.warn("不支持的redisTypes:",type)
      return Promise.reject("不支持的redisTypes:"+type)
    }


    executeCommand(commands).then(data => {
      len.value = data[0]

      let ret = data[1]
      if(commands[commands.length-1][0].endsWith("scan")){
        ret = ret[1]
        cursor = ret[0] // todo
      }

      ret = redisType.listResultHandler(ret)
      list.value = ret
      return ret
    })
  }

  function load(key: string) {
    if(!key){
      return
    }
    loadMeta(key).then(() => {
      value.value = ""
      len.value = 0
      list.value = []
      if (type.value === RedisTypes.String){
        loadValue(key)
      }else{
        loadList(key,"*","0")
      }
    })
  }

  interface Data{
    value:string,member:string,score:string,id:string;field:string
  }

  function set(key:string,type:string,data:Data,oldData:Data){

    let commands = []
    if(type === RedisTypes.String){
      commands.push(["set",key,data.value])
    }else{
      let redisType = redisTypes[type]
      if(redisType){
        commands.push(...redisType.set(key,data,oldData))
      }else{
        console.warn("不支持的redisTypes:",type)
        return Promise.reject("不支持的redisTypes:"+type)
      }
    }


    return executeCommand(commands).then(()=>{
      load(key)
    })
  }


  function del({key,type,data}:{key:string|string[];type?:string;data?:Record<string, any>}){

    let commands = []
    if (data){

      let redisType = redisTypes[type as string]
      if(redisType){
        commands.push(redisType.del(key as string,data))
      }else{
        console.warn("不支持的redisTypes:",type)
        return Promise.reject("不支持的redisTypes:"+type)
      }
    }else{
      if(Array.isArray(key)){
        commands.push(["del"].concat(key))
      }else{
        commands.push(["del",key])
      }

    }

    return executeCommand(commands)
  }

  function expire(k:string,ttl:number){
    return executeCommand([
      ["expire",k+"",ttl+""]
    ])
  }

  function rename(old:string, n:string){
    return executeCommand([
      ["rename",old,n]
    ])
  }


  return {
    type, ttl, value, list,len,
    del,expire,rename,
    load,set
  }
}
