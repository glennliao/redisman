<template>
  <n-input v-model:value="val" type="textarea"> </n-input>
  <div class="mt-1 justify-end flex">
    <n-button @click="save" type="info" ghost size="small">save</n-button>
  </div>
</template>

<script>
export default {
  name: "KeyValue",
  props:{
    value:{}
  },
  emits:["save"] ,
  setup(props,{emit}){

    const val = ref('')

    watch(()=>props.value,()=>{
      val.value = props.value
    },{immediate:true})

    function save(){
      emit("save",val.value)
    }

    return {
      save,val
    }
  }
}
</script>

<style scoped>

</style>
