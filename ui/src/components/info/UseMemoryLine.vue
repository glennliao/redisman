<template>
  <n-card size="small" title="UsedMemory">
    <div id="usedMemory"></div>
  </n-card>
</template>

<script>
import {Line} from "@antv/g2plot";

export default {
  name: "UseMemoryLine",
  props:{
    data:{}
  },
  setup(props){

    let line = null



    onMounted(()=>{
      line = new Line('usedMemory', {
        data: [],
        xField: 'time',
        yField: 'value',
        xAxis: {
          type: 'time',
          mask: 'HH:MM:ss',
          title: {
            text: "time"
          }
        },
        height:200
      });
      line.render();
      line.changeData(props.data)
    })

    onUnmounted(()=>{
      line.destroy()
    })

    watch(()=>props.data,(val)=>{
      line.changeData(val)
    })
    return {

    }
  }
}
</script>

<style scoped>

</style>
