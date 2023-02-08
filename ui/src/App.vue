<script setup lang="ts">
import { logLabeled } from "./utils/logger";
import Navbar from "./components/layout/Navbar/index.vue";
import Footer from "./components/layout/Footer/index.vue";
logLabeled("I'm ready!  ⸜(๑'ᵕ'๑)⸝⋆*", "log", "", "color: #66ccff");

</script>

<template>
  <n-message-provider>

  <div class="font-chinese antialiased">
    <div class="min-h-screen" style="min-height: calc(100vh - 0px);min-width:590px">
      <Navbar />
      <n-dialog-provider>
      <RouterView />
      </n-dialog-provider>
    </div>
<!--    <Footer />-->
  </div>

  </n-message-provider>
</template>
