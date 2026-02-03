<template>
  <div class="term-input-wrapper">
    <span v-if="prompt" class="prompt">{{ prompt }}</span>
    <div class="input-container" :class="{ 'focused': isFocused }">
      <input 
        v-bind="$attrs"
        :type="currentType"
        v-model="modelValue"
        class="real-input"
        @focus="isFocused = true"
        @blur="isFocused = false"
      />
      <div v-if="isFocused" class="input-glow"></div>
      
      <button 
        v-if="isPasswordType" 
        class="toggle-btn" 
        @click.prevent="toggleShow"
        tabindex="-1"
        type="button"
      >
        <span class="icon">{{ showPassword ? '👁️' : '🔒' }}</span>
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
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  }
});

const emit = defineEmits(['update:modelValue']);

const showPassword = ref(false);
const isFocused = ref(false);
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
  margin-bottom: 1.2rem;
  width: 100%;
}

.prompt {
  color: var(--term-primary);
  margin-right: 1rem;
  white-space: nowrap;
  font-weight: bold;
  font-size: 0.9rem;
  text-shadow: 0 0 5px rgba(0, 255, 157, 0.2);
  min-width: 60px;
  text-align: right;
}

.input-container {
  flex: 1;
  display: flex;
  align-items: center;
  position: relative;
  border-bottom: 2px solid var(--term-surface-border);
  transition: all 0.3s var(--ease-out);
  background: rgba(255, 255, 255, 0.02);
  border-radius: 4px 4px 0 0;
}

.input-container.focused {
  border-bottom-color: var(--term-primary);
  background: rgba(255, 255, 255, 0.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.real-input {
  flex: 1;
  width: 100%;
  background: transparent;
  border: none;
  color: var(--term-text);
  font-family: var(--term-font);
  font-size: 1rem;
  padding: 0.5rem 0.75rem;
  outline: none;
  caret-color: var(--term-primary);
  z-index: 1;
}

.input-glow {
    position: absolute;
    bottom: -2px; left: 0; width: 100%; height: 2px;
    background: var(--term-primary);
    box-shadow: 0 0 10px var(--term-primary);
    animation: glow-pulse 2s infinite alternate;
}

@keyframes glow-pulse {
    from { opacity: 0.5; box-shadow: 0 0 5px var(--term-primary); }
    to { opacity: 1; box-shadow: 0 0 15px var(--term-primary), 0 -5px 10px rgba(0,255,157,0.2); }
}

.toggle-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0 0.75rem;
    font-size: 1rem;
    color: var(--term-text-muted);
    transition: color 0.2s;
    height: 100%;
    display: flex;
    align-items: center;
}

.toggle-btn:hover {
    color: var(--term-primary);
}

.real-input::placeholder {
  color: var(--term-text-muted);
  opacity: 0.4;
  font-style: italic;
}
</style>
