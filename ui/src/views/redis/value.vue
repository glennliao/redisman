<script>
export default {
  name: "KeyValue",
  props: {
    value: {},
  },
  emits: ["save"],
  setup(props, { emit }) {
    const val = ref("");
    const contentMode = ref("text")

    watch(() => props.value, () => {



      let content = props.value
      if(content === ""){
        return
      }

      try{
        content = JSON.parse(content)
        val.value = JSON.stringify(content,null,2);
        contentMode.value = "json"
      }catch (e){
        console.log(e,content)
        val.value = props.value;
      }


    }, { immediate: true });

    function save() {
      let content = val.value
      if(contentMode.value === "json"){
        content = JSON.stringify(JSON.parse(content))
      }
      emit("save", content);
    }


    return {
      save, val,contentMode
    };
  },
};
</script>

<template>
  <codemirror v-model:value="val" :contentMode="contentMode"/>
  <div class="mt-1 flex justify-end">
    <n-button type="info"  size="small" @click="save">
      save
    </n-button>
  </div>
</template>

<style scoped>

</style>
