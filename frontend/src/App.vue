<template>
  <div class="app-container" @mousemove="onActivity" @keydown="onActivity" @touchstart="onActivity">
    <LoginView v-if="!appToken" @success="handleLogin" />
    <TabManager v-else :token="appToken" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import LoginView from './components/LoginView.vue';
import TabManager from './components/TabManager.vue';

const appToken = ref('');
let idleInterval: any = null;

const handleLogin = (token: string) => {
    appToken.value = token;
    localStorage.setItem('nexus_token', token);
    updateActivity();
};

const updateActivity = () => {
    localStorage.setItem('nexus_last_active', Date.now().toString());
};

const onActivity = () => {
    // Throttling could be good, but for now simple update
    // Only update if logged in
    if (appToken.value) {
        // Optimize: only write if > 1 min since last write?
        // Let's do a simple memory check first to avoid spamming I/O
        const now = Date.now();
        const last = parseInt(localStorage.getItem('nexus_last_active') || '0');
        if (now - last > 5000) { // Update every 5s max
            updateActivity();
        }
    }
};

const checkIdle = () => {
    if (!appToken.value) return;
    
    const timeoutMin = parseInt(localStorage.getItem('nexus_login_timeout') || '30');
    if (timeoutMin <= 0) return; // Disabled

    const last = parseInt(localStorage.getItem('nexus_last_active') || '0');
    const diff = Date.now() - last;
    
    if (diff > timeoutMin * 60 * 1000) {
        // Timed out
        console.log("Session timed out");
        logout();
    }
};

const logout = () => {
    appToken.value = '';
    localStorage.removeItem('nexus_token');
    // localStorage.removeItem('nexus_last_active'); // Optional, keep for debugging
};

onMounted(() => {
    // specific logic for init
    const stored = localStorage.getItem('nexus_token');
    if (stored) {
        // Check if expired
        const last = parseInt(localStorage.getItem('nexus_last_active') || '0');
        const timeoutMin = parseInt(localStorage.getItem('nexus_login_timeout') || '30');
        
        if (timeoutMin > 0 && Date.now() - last > timeoutMin * 60 * 1000) {
            console.log("Restored session expired");
            logout();
        } else {
            appToken.value = stored;
            updateActivity(); // Refresh on load
        }
    }
    
    idleInterval = setInterval(checkIdle, 60 * 1000); // Check every minute
});

onUnmounted(() => {
    if (idleInterval) clearInterval(idleInterval);
});
</script>

<style>
/* Global styles are now handled in assets/terminal.css */
.app-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
</style>
