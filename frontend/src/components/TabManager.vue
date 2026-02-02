<template>
  <div class="tab-manager">
    <!-- Tab Bar -->
    <div class="tab-bar">
      <div 
        v-for="tab in tabs" 
        :key="tab.id" 
        class="tab-item" 
        :class="{ active: activeTabId === tab.id }"
        @click="activeTabId = tab.id"
      >
        <span class="tab-title">{{ tab.title }}</span>
        <span class="tab-close" @click.stop="closeTab(tab.id)">[x]</span>
      </div>
      <button class="add-tab" @click="addTab" title="New Tab">[+]</button>
    </div>

    <!-- Tab Content -->
    <div class="tab-content">
      <!-- We keep all tabs mounted to preserve connection/state, using v-show -->
      <TerminalView 
        v-for="tab in tabs"
        :key="tab.id"
        :token="token" 
        v-show="activeTabId === tab.id"
        :visible="activeTabId === tab.id"
        @update:title="(t) => updateTitle(tab.id, t)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import TerminalView from './TerminalView.vue';

const props = defineProps({
  token: String
});

interface Tab {
  id: number;
  title: string;
}

const tabs = ref<Tab[]>([]);
const activeTabId = ref(0);
let nextId = 1;

const addTab = () => {
  const id = nextId++;
  tabs.value.push({
    id,
    title: 'TERMINAL ' + id
  });
  activeTabId.value = id;
};

const closeTab = (id: number) => {
  const idx = tabs.value.findIndex(t => t.id === id);
  if (idx !== -1) {
    tabs.value.splice(idx, 1);
    // If closed active tab
    if (activeTabId.value === id) {
       // Switch to neighbor
       const newTab = tabs.value[idx] || tabs.value[idx - 1];
       if (newTab) {
         activeTabId.value = newTab.id;
       } else {
         activeTabId.value = 0; // No tabs
       }
    }
  }
};

const updateTitle = (id: number, title: string) => {
  const tab = tabs.value.find(t => t.id === id);
  if (tab) {
    tab.title = title.toUpperCase();
  }
};

onMounted(() => {
  addTab(); // Initial tab
});
</script>

<style scoped>
.tab-manager {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--term-bg);
  overflow: hidden;
}

.tab-bar {
  display: flex;
  background: var(--term-bg);
  height: 35px;
  align-items: center;
  overflow-x: auto;
  border-bottom: 1px solid var(--term-muted);
  flex-shrink: 0;
  padding-left: 0.5rem;
}

.tab-item {
  display: flex;
  align-items: center;
  padding: 0 1rem;
  height: 100%;
  color: var(--term-muted);
  background: transparent;
  cursor: pointer;
  min-width: 120px;
  max-width: 200px;
  user-select: none;
  font-family: var(--term-font);
  border-right: 1px solid var(--term-muted);
  transition: all 0.1s;
}

.tab-item:hover {
  color: var(--term-text);
  background: rgba(51, 255, 0, 0.1);
}

.tab-item.active {
  background: var(--term-text);
  color: var(--term-bg);
}

.tab-title {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 0.9rem;
}

.tab-close {
  margin-left: 8px;
  font-size: 0.8rem;
  opacity: 0.5;
}

.tab-close:hover {
  opacity: 1;
  font-weight: bold;
}

.add-tab {
  background: transparent;
  border: none;
  color: var(--term-muted);
  font-family: var(--term-font);
  font-size: 1rem;
  padding: 0 1rem;
  cursor: pointer;
}
.add-tab:hover {
  color: var(--term-text);
}

.tab-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}
</style>
