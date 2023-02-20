import {computed} from "vue";
import {redis} from "~/api";

const keys = ref([] as string[]);

const treeKeys = computed(() => buildKeysTree(keys.value));
const count = 20000

function _scan(pattern = "*",cur:string): Promise<string[]> {
  return redis.command([
    ["scan", cur, "MATCH", pattern, "COUNT", count+""],
  ])
    //.then((data) => {
    // cursor.value = data[0][0]
    // console.log(data)
    // return data[0][1].sort();
  //});
}


async function scan(pattern = "*") {

  keys.value = []
  let cursor = "0"
  let _keys:string[] = []

  // do {
  //   let data = await _scan(pattern, cursor)
  //   cursor = data[0][0]
  //   _keys = _keys.concat(data[0][1])
  //   keys.value = _keys.sort()
  // }while (cursor != "0")

  let data = await _scan(pattern, cursor)
  cursor = data[0][0]
  _keys = _keys.concat(data[0][1])
  keys.value = _keys.sort()

}


function scanKeys(pattern = "*",cur:string): Promise<string[]> {
  return _scan(pattern,"0")
}

export function useKeysHook() {
  return {
    keys,
    treeKeys,
    scan,
    scanKeys,
  };
}

function buildKeysTree(keys: string[], splitStr = ":") {
  const parentSet = new Set<string>();
  const keyList = keys.map((item) => {
    let label = item;
    let parentId = "root";

    if (item.includes(splitStr)) {
      const splits = item.split(splitStr).filter(item => item);
      if (splits.length > 1) {
        const pid = [];
        for (let i = 0; i < splits.length; i++) {
          if (i === splits.length - 1) {
            label = splits[i];
            if(item.endsWith(splitStr)){
              label += splitStr
            }
            parentId = pid.join(splitStr) + splitStr;
          } else {
            pid.push(splits[i]);
            parentSet.add(pid.join(splitStr) + splitStr);
          }
        }
      }
    }

    return {
      key: item, label, parentId:'pid_'+parentId, leaf: true,
    };
  });
  const pList = Array.from(parentSet.values());



  return toTree(pList.map((k: string) => {
    const splits = k.split(splitStr);
    let parentId = "root";
    let label = k;
    if (splits.length > 2) {
      // 顶目录
      parentId = splits.slice(0, -2).join(splitStr) + splitStr;
    }
    label = splits[splits.length - 2];
    return {
      key: 'pid_' + k, label, parentId: 'pid_' + parentId, leaf: false,
    };
  }).concat(keyList), "pid_root")
}

interface TreeNode {
  label: string;
  children?: TreeNode[];
  childrenCnt?: number;

}

function toTree(list: any[], rootId = "root") {
  if (list.length === 0) {
    return [];
  }
  const itemMap: Record<string, Array<any>> = {};

  for (const item of list) {
    const parentId = item.parentId;
    itemMap[parentId] = itemMap[parentId] || [];
    itemMap[parentId].push(item);
    item.childrenCnt = 1;
  }
  for (const item of list) {
    const id = item.key;
    if (itemMap[id]) {
      if (itemMap[id].length > 0) {
        item.children = itemMap[id];
      }
    }
  }

  const rootTree: TreeNode[] = [];
  itemMap[rootId] && itemMap[rootId].forEach((item: any) => {
    rootTree.push(item);
  });

  function calcChildrenCnt(root: TreeNode[]) {
    let childrenCnt = 0;
    root.forEach((item) => {
      if (item.children) {
        item.childrenCnt = calcChildrenCnt(item.children);
        childrenCnt += item.childrenCnt;
      } else {
        childrenCnt++;
      }
    });

    return childrenCnt;
  }

  calcChildrenCnt(rootTree);

  return rootTree;
}
