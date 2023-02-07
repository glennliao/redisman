import {executeCommand} from "~/api/redis";
import {computed} from "vue";

const keys = ref([])
const cursor = ref(0)

const treeKeys = computed(() => buildKeysTree(keys.value))


function scan(pattern: string = "*") {
  executeCommand([
    ["scan", cursor.value + "", "MATCH", pattern, "COUNT", "10000"],
  ]).then(data => {
    keys.value = data[0][1].sort()
  })
}

function scanKeys(pattern: string = "*"): Promise<string[]> {
  return executeCommand([
    ["scan", cursor.value + "", "MATCH", pattern, "COUNT", "10000"],
  ]).then(data => {
    return data[0][1].sort()
  })
}



export function useKeysHook() {
  return {
    keys,
    treeKeys,
    scan,
    scanKeys
  }
}

function buildKeysTree(keys: string[], splitStr: string = ":") {
  let parentSet = new Set<string>()
  let keyList = keys.map(item => {
    let label = item
    let parentId = "root"
    if (item.includes(splitStr)) {
      let splits = item.split(splitStr).filter(item => item)
      if (splits.length > 1) {
        let pid = []
        for (let i = 0; i < splits.length; i++) {
          if (i === splits.length - 1) {
            label = splits[i]
            parentId = pid.join(splitStr) + splitStr
          } else {
            pid.push(splits[i])
            parentSet.add(pid.join(splitStr) + splitStr)
          }
        }

      }
    }

    return {
      key: item, label: label, parentId, leaf: true
    }
  })
  let pList = Array.from(parentSet.values())


  return toTree(pList.map((k: string) => {
    const splits = k.split(splitStr)
    let parentId = 'root'
    let label = k
    if (splits.length > 2) {
      // 顶目录
      parentId = splits.slice(0, -2).join(splitStr) + splitStr
    }
    label = splits[splits.length - 2]
    return {
      key: k, label: label, parentId, leaf: false
    }
  }).concat(keyList), 'root')
}


interface TreeNode {
  label: string
  children?: TreeNode[]
  childrenCnt?: number

}

function toTree(list: any[], rootId: string = "root") {
  if (list.length === 0) return []
  const itemMap: Record<string, Array<any>> = {}

  for (let item of list) {
    let parentId = item['parentId']
    itemMap[parentId] = itemMap[parentId] || []
    itemMap[parentId].push(item)
    item.childrenCnt = 1
  }
  for (let item of list) {
    const id = item.key
    if (itemMap[id]) {
      if (itemMap[id].length > 0) {
        item.children = itemMap[id]
      }
    }
  }


  const rootTree: TreeNode[] = []
  itemMap[rootId] && itemMap[rootId].forEach((item: any) => {
    rootTree.push(item)
  });


  function calcChildrenCnt(root: TreeNode[]) {
    let childrenCnt = 0
    root.forEach(item => {
      if (item.children) {
        item.childrenCnt = calcChildrenCnt(item.children)
        childrenCnt += item.childrenCnt
      } else {
        childrenCnt++
      }
    })

    return childrenCnt
  }

  calcChildrenCnt(rootTree)

  return rootTree
}
