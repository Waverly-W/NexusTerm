<template>
  <div class="term-card">
    <div v-if="title" class="term-card-header">
      <span class="corner">+</span>
      <span class="title-text">--- {{ title }} ---</span>
      <span class="corner">+</span>
    </div>
    <div class="term-card-body">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  title: String,
});
</script>

<style scoped>
.term-card {
  border: 1px solid var(--term-muted);
  background: rgba(10, 10, 10, 0.9);
  position: relative;
  margin-bottom: 1rem;
}

.term-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--term-muted);
  background: var(--term-muted);
  color: var(--term-bg);
  font-weight: bold;
  padding: 0 0.5rem;
  height: 2rem;
}

.title-text {
  text-transform: uppercase;
}

.term-card-body {
  padding: 1rem;
}

/* Optional: Make it look like a window */
</style>
