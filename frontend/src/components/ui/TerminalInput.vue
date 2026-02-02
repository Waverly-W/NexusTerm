<template>
  <div class="term-input-wrapper">
    <span class="prompt">{{ prompt }}</span>
    <div class="input-container">
      <input 
        v-bind="$attrs"
        v-model="modelValue"
        class="real-input"
        type="text"
      />
      <!-- Visual cursor logic could be complex, for now we rely on the native caret and styling -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps({
  modelValue: [String, Number],
  prompt: {
    type: String,
    default: 'user@nexus:~$'
  }
});

const emit = defineEmits(['update:modelValue']);

const modelValue = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
});
</script>

<style scoped>
.term-input-wrapper {
  display: flex;
  align-items: center;
  font-family: var(--term-font);
  margin-bottom: 0.5rem;
  width: 100%;
}

.prompt {
  color: var(--term-secondary);
  margin-right: 0.75rem;
  white-space: nowrap;
}

.input-container {
  flex: 1;
  position: relative;
}

.real-input {
  width: 100%;
  background: transparent;
  border: none;
  color: var(--term-text);
  font-family: var(--term-font);
  font-size: 1rem;
  outline: none;
  caret-color: var(--term-text); /* The blinking cursor */
}

/* Custom placeholder */
.real-input::placeholder {
  color: var(--term-muted);
  opacity: 0.5;
}
</style>
