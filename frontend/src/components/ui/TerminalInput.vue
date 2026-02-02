<template>
  <div class="term-input-wrapper">
    <span class="prompt">{{ prompt }}</span>
    <div class="input-container">
      <input 
        v-bind="$attrs"
        :type="currentType"
        v-model="modelValue"
        class="real-input"
        :class="{ 'has-toggle': isPasswordType }"
      />
      <button 
        v-if="isPasswordType" 
        class="toggle-btn" 
        @click.prevent="toggleShow"
        tabindex="-1"
        type="button"
      >
        {{ showPassword ? '[HIDE]' : '[SHOW]' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

const props = defineProps({
  modelValue: [String, Number],
  prompt: {
    type: String,
    default: 'user@nexus:~$'
  },
  type: {
    type: String,
    default: 'text'
  }
});

const emit = defineEmits(['update:modelValue']);

const showPassword = ref(false);
const isPasswordType = computed(() => props.type === 'password');

const currentType = computed(() => {
  if (!isPasswordType.value) return props.type;
  return showPassword.value ? 'text' : 'password';
});

const toggleShow = () => {
    showPassword.value = !showPassword.value;
};

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
  display: flex;
  align-items: center;
  position: relative;
  border-bottom: 1px dotted transparent;
  transition: border-color 0.2s;
}

.input-container:focus-within {
  border-bottom-color: var(--term-muted);
}

.real-input {
  flex: 1;
  width: 100%;
  background: transparent;
  border: none;
  color: var(--term-text);
  font-family: var(--term-font);
  font-size: 1rem;
  outline: none;
  caret-color: var(--term-text);
  padding: 0;
}

.real-input.has-toggle {
    padding-right: 4rem;
}

.toggle-btn {
    position: absolute;
    right: 0;
    background: transparent;
    border: none;
    color: var(--term-muted);
    font-family: var(--term-font);
    cursor: pointer;
    font-size: 0.8rem;
    padding: 0;
    user-select: none;
}
.toggle-btn:hover {
    color: var(--term-text);
}

/* Custom placeholder */
.real-input::placeholder {
  color: var(--term-muted);
  opacity: 0.5;
}
</style>
