import type { RouteRecordRaw } from "vue-router";
import {createRouter, createWebHashHistory} from "vue-router";
import NProgress from "nprogress";

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    component: () => import("~/views/redis/index.vue"),
  },
  {
    path: "/404",
    name: "ErrorPage",
    component: () => import("~/views/errorPages/index.vue"),
  },
  // {
  //   path: "/home",
  //   name: "Home",
  //   component: () => import("~/views/home/index.vue"),
  // },
  {
    path: "/redis",
    name: "Redis",
    component: () => import("~/views/redis/index.vue"),
  },
  {
    path: "/redis/conn",
    name: "RedisConn",
    component: () => import("~/views/redis/conn.vue"),
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("~/views/login/index.vue"),
  },
  // {
  //   path: "/hero",
  //   name: "Hero",
  //   component: () => import("~/views/home/components/Hero.vue"),
  // },
  // {
  //   path: "/StoreTest",
  //   name: "StoreTest",
  //   component: () => import("~/views/home/components/StoreTest.vue"),
  // },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/404",
  },
];

const index = createRouter({
  history: createWebHashHistory(),
  routes,
});
index.beforeEach(() => {
  if (!NProgress.isStarted()) {
    NProgress.start();
  }
});

index.afterEach(() => {
  NProgress.done();
});

export default index;
