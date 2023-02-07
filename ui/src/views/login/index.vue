<template>
  <div>
    <n-button @click="login"> login</n-button>
  </div>
</template>

<script setup>
const router = useRouter()
function login(){
  router.push("/redis")
}
</script>

<style scoped>

</style>
