<template>
  <button 
    class="term-btn" 
    :class="[`variant-${variant}`, { 'is-pushed': isPushed }]"
    @mousedown="isPushed = true"
    @mouseup="isPushed = false"
    @mouseleave="isPushed = false"
  >
    <span class="bracket">[</span>
    <span class="content"><slot></slot></span>
    <span class="bracket">]</span>
  </button>
</template>

<script setup lang="ts">
import { ref } from 'vue';

defineProps({
  variant: {
    type: String, // 'primary', 'secondary', 'danger'
    default: 'primary'
  }
});

const isPushed = ref(false);
</script>

<style scoped>
.term-btn {
  background: transparent;
  color: var(--term-text);
  border: none;
  font-family: var(--term-font);
  font-size: 1rem;
  cursor: pointer;
  padding: 0 0.5rem;
  text-transform: uppercase;
  transition: all 0.1s;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  user-select: none;
}

.bracket {
  opacity: 0.5;
  transition: opacity 0.1s;
}

.term-btn:hover {
  background: var(--term-text);
  color: var(--term-bg);
  text-shadow: none;
}

.term-btn:hover .bracket {
  opacity: 1;
  color: var(--term-bg);
}

/* Variants */
.variant-secondary {
  color: var(--term-secondary);
}
.variant-secondary:hover {
    background: var(--term-secondary);
    color: var(--term-bg);
}

.variant-danger {
  color: var(--term-error);
}
.variant-danger:hover {
    background: var(--term-error);
    color: var(--term-bg);
}

/* Pushed State */
.is-pushed {
  transform: translateY(1px);
}
</style>
