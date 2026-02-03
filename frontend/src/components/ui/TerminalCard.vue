<template>
  <div class="term-card glass-panel animate-fade-in">
    <div v-if="title" class="term-card-header">
        <div class="header-decoration"></div>
        <span class="title-text">{{ title }}</span>
        <div class="header-decoration"></div>
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
  position: relative;
  margin-bottom: 1.5rem;
  overflow: hidden; /* For glass blur edges */
  transition: transform 0.2s var(--ease-out), box-shadow 0.2s;
}

.term-card:hover {
    box-shadow: 0 8px 32px rgba(0,0,0,0.6);
    border-color: rgba(255,255,255,0.2);
}

.term-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid var(--term-surface-border);
  gap: 1rem;
}

.title-text {
  color: var(--term-primary);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 0.9rem;
  text-shadow: 0 0 5px rgba(0,255,157,0.3);
}

.header-decoration {
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--term-surface-border), transparent);
}

.term-card-body {
  padding: 1.5rem;
}
</style>
