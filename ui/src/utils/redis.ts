

export function parseInfo(infoStr:string) {
  let info:Record<string, any> = {}
  let curSection = ""
  infoStr.split("\r\n").forEach(item=>{
    if(!item.trim()){
      return
    }
    if (item.startsWith("# ")){
      curSection = item.substring(2).trim()
      info[curSection] = {}
    }else{
      let sps = item.split(":")
      info[curSection][sps[0]] = sps[1]
    }
  })


  return info
}
