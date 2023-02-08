<script>
export default {
  name: "KeyValue",
  props: {
    value: {},
  },
  emits: ["save"],
  setup(props, { emit }) {
    const val = ref("");

    watch(() => props.value, () => {
      val.value = props.value;
    }, { immediate: true });

    function save() {
      emit("save", val.value);
    }

    return {
      save, val,
    };
  },
};
</script>

<template>
  <codemirror v-model:value="val"/>
  <div class="mt-1 flex justify-end">
    <n-button type="info"  size="small" @click="save">
      save
    </n-button>
  </div>
</template>

<style scoped>

</style>
