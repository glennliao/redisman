<script setup lang="ts">
import ThemeChange from "./components/ThemeChange/index.vue";
import LocalesChange from "./components/LocalesChange/index.vue";
</script>

<template>
  <div
    class="sticky top-0 z-30 flex h-10 w-full justify-center text-base-content opacity-90 backdrop-blur transition-all duration-100"
  >
    <nav class="navbar w-full" style="min-height:2.5rem" >
      <div class="flex flex-1 md:gap-1 lg:gap-2">
        <router-link
          to="/"
          aria-current="page"
          aria-label="Homepage"
          class="btn-ghost btn px-2"
        >
          <div
            class="inline-flex text-lg text-primary transition-all duration-200 md:text-3xl"
          >
            <span class="text-accent">R<span class="lowercase">edis</span></span>
            <span class="text-base-content">M<span class="lowercase">an</span></span>
          </div>
        </router-link>
      </div>
      <div>
        <ThemeChange />
        <LocalesChange />
        <span
          class="tooltip tooltip-bottom before:text-xs before:content-[attr(data-tip)]"
          data-tip="GitHub"
        >
          <span class="flex-none items-center">
            <a
              aria-label="Github"
              target="_blank"
              href="https://github.com/glennliao/redisman"
              class="btn-ghost drawer-button btn-square btn normal-case"
            >
              <IMdiGithub />
            </a>
          </span>
        </span>
      </div>
    </nav>
  </div>
</template>
