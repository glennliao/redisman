<template>
  <div class="toolbar mb-1 flex justify-end" >
    <div  class="item" style="width:120px">
      <n-select size="small" v-model:value="viewMode"  :options="viewModeOptions"></n-select>
    </div>
  </div>
  <codemirror
    v-if="!viewModeChangeLoading"
    v-model="code"
    placeholder=""
    :style="{ height: '400px' }"
    :autofocus="true"
    :indent-with-tab="true"
    :tab-size="2"
    :extensions="extensions"
    @ready="handleReady"
    @change="onChange"
  />
</template>

<script>
import { defineComponent, ref, shallowRef } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { json,jsonParseLinter } from '@codemirror/lang-json'
import {linter, lintGutter} from '@codemirror/lint'
import { oneDark } from '@codemirror/theme-one-dark'

export default defineComponent({
  components: {
    Codemirror
  },
  props:{
    value:{

    }
  },
  setup(props,{emit}) {
    const code = ref(``)
    const extensions = computed(()=>{
      let list = [oneDark]
      if(viewMode.value === "json"){
        list = list.concat([json(), linter(jsonParseLinter()), lintGutter()])
      }
      return list
    })

    watch(()=>props.value,(val)=>{
      code.value = val
    },{immediate:true})



    // Codemirror EditorView instance ref
    const view = shallowRef()
    const handleReady = (payload) => {
      view.value = payload.view
    }

    // Status is available at all times via Codemirror EditorView
    const getCodemirrorStates = () => {
      const state = view.value.state
      const ranges = state.selection.ranges
      const selected = ranges.reduce((r, range) => r + range.to - range.from, 0)
      const cursor = ranges[0].anchor
      const length = state.doc.length
      const lines = state.doc.lines
      // more state info ...
      // return ...
    }

    function onChange(val){
      emit("update:value",val)
    }

    const viewMode = ref("text")
    const viewModeOptions = ["text","json"].map(item=>{
      return {
        value:item,label:item
      }
    })

    watch(()=>viewMode.value,()=>{
      viewModeChangeLoading.value = true
      setTimeout(()=>{
        viewModeChangeLoading.value = false
      },10)
    })

    const viewModeChangeLoading = ref(false)

    return {
      code,
      extensions,
      handleReady,
      log: console.log,
      onChange,
      viewMode,viewModeOptions,viewModeChangeLoading
    }
  }
})
</script>
<style scoped lang="less">
.toolbar {
  border: 1px solid #f1f1f1;
  border-bottom: none;
  border-radius: 6px 6px 0 0;
  padding: 4px 4px 0;
}
</style>
