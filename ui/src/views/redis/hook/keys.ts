import { computed } from "vue";
import { redis } from "~/api";

const keys = ref([]);
const cursor = ref(0);

const treeKeys = computed(() => buildKeysTree(keys.value));

function scan(pattern = "*") {
  redis.command([
    ["scan", `${cursor.value}`, "MATCH", pattern, "COUNT", "10000"],
  ]).then((data) => {
    keys.value = data[0][1].sort();
  });
}

function scanKeys(pattern = "*"): Promise<string[]> {
  return redis.command([
    ["scan", `${cursor.value}`, "MATCH", pattern, "COUNT", "10000"],
  ]).then((data) => {
    return data[0][1].sort();
  });
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



  let treeKeys = toTree(pList.map((k: string) => {
    const splits = k.split(splitStr);
    let parentId = "root";
    let label = k;
    if (splits.length > 2) {
      // 顶目录
      parentId = splits.slice(0, -2).join(splitStr) + splitStr;
    }
    label = splits[splits.length - 2];
    return {
      key: 'pid_'+k, label, parentId:'pid_'+parentId, leaf: false,
    };
  }).concat(keyList), "pid_root");

  return treeKeys
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
