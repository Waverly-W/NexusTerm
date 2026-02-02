<template>
  <div class="app-container">
    <LoginView v-if="!appToken" @success="(t) => appToken = t" />
    <TabManager v-else :token="appToken" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import LoginView from './components/LoginView.vue';
import TabManager from './components/TabManager.vue';

const appToken = ref('');
</script>

<style>
/* Global styles are now handled in assets/terminal.css */
.app-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
</style>
