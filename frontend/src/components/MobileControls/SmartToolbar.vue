<template>
  <div class="toolbar">
    <div class="btn-group">
      <button class="btn" @click="$emit('key', 'Escape')">[ESC]</button>
      <button class="btn" @click="$emit('key', 'Tab')">[TAB]</button>
    </div>
    <div class="btn-group">
      <button 
        class="btn" 
        :class="{ active: ctrlActive }" 
        @click="$emit('toggle-ctrl')"
      >
        {{ ctrlActive ? '[CTRL*]' : '[CTRL]' }}
      </button>
      <button class="btn" @click="$emit('key', 'Alt')">[ALT]</button>
    </div>
    <div class="btn-group extra-keys">
      <button class="btn" @click="$emit('input', '-')">-</button>
      <button class="btn" @click="$emit('input', '/')">/</button>
      <button class="btn" @click="$emit('input', '|')">|</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  ctrlActive: boolean
}>();

defineEmits<{
  (e: 'key', key: string): void,
  (e: 'input', char: string): void,
  (e: 'toggle-ctrl'): void
}>();
</script>

<style scoped>
.toolbar {
  display: flex;
  height: 48px;
  background: var(--term-bg);
  border-top: 1px solid var(--term-muted);
  overflow-x: auto;
  align-items: center;
  padding: 0 4px;
  gap: 8px;
}

.btn-group {
  display: flex;
  gap: 8px;
}

.btn {
  height: 36px;
  min-width: 44px;
  background: transparent;
  border: 1px solid var(--term-muted);
  border-radius: 0;
  color: var(--term-text);
  font-family: var(--term-font);
  font-size: 13px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 8px;
  user-select: none;
  cursor: pointer;
  white-space: nowrap;
}

.btn:active {
  background: var(--term-text);
  color: var(--term-bg);
  border-color: var(--term-text);
}

.btn.active {
  background: var(--term-text);
  color: var(--term-bg);
  border-color: var(--term-text);
  box-shadow: 0 0 8px var(--term-text);
}

.extra-keys {
  margin-left: auto;
}
</style>
