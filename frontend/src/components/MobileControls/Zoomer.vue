<template>
  <div 
    class="zoomer"
    ref="zoomerRef"
    @touchstart="onTouchStart"
    @touchmove="onTouchMove"
    @touchend="onTouchEnd"
    @mousedown="onMouseDown"
  >
    <div class="zoomer-track">
      <div 
        class="zoomer-handle" 
        :style="{ top: handlePosition + '%' }"
      >
        <span v-if="isActive" class="zoomer-tooltip">{{ previewSize }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

const props = defineProps<{
  currentSize: number
}>();

const emit = defineEmits<{
  (e: 'change', size: number): void
}>();

const MIN_SIZE = 8;
const MAX_SIZE = 36;

const zoomerRef = ref<HTMLElement | null>(null);
const isActive = ref(false);
const localPercent = ref(0); // 0 (Top/Max) to 100 (Bottom/Min)

// Sync prop to local state when not dragging
watch(() => props.currentSize, (newVal) => {
  if (!isActive.value) {
    localPercent.value = sizeToPercent(newVal);
  }
}, { immediate: true });

function sizeToPercent(size: number) {
  // Size 36 (Max) -> 0%
  // Size 8 (Min) -> 100%
  const range = MAX_SIZE - MIN_SIZE;
  return ((MAX_SIZE - size) / range) * 100;
}

function percentToSize(percent: number) {
  // 0% -> 36
  // 100% -> 8
  const range = MAX_SIZE - MIN_SIZE;
  const size = MAX_SIZE - (percent / 100 * range);
  return Math.round(size);
}

const handlePosition = computed(() => {
  return Math.min(100, Math.max(0, localPercent.value));
});

const previewSize = computed(() => {
  return percentToSize(localPercent.value);
});

const onTouchStart = (e: TouchEvent) => {
  isActive.value = true;
  if (e.touches[0]) {
    updateFromClientY(e.touches[0].clientY);
  }
};

const onTouchMove = (e: TouchEvent) => {
  if (e.cancelable) e.preventDefault();
  if (e.touches[0]) {
    updateFromClientY(e.touches[0].clientY);
  }
};

const onTouchEnd = () => {
  isActive.value = false;
  emit('change', previewSize.value);
};

// Mouse Support
const onMouseDown = (e: MouseEvent) => {
  e.preventDefault(); // Prevent text selection
  isActive.value = true;
  updateFromClientY(e.clientY);
  window.addEventListener('mousemove', onMouseMove);
  window.addEventListener('mouseup', onMouseUp);
};

const onMouseMove = (e: MouseEvent) => {
  if (!isActive.value) return;
  e.preventDefault();
  updateFromClientY(e.clientY);
};

const onMouseUp = () => {
  if (isActive.value) {
    isActive.value = false;
    emit('change', previewSize.value);
  }
  window.removeEventListener('mousemove', onMouseMove);
  window.removeEventListener('mouseup', onMouseUp);
};

const updateFromClientY = (clientY: number) => {
  if (!zoomerRef.value) return;
  const rect = zoomerRef.value.getBoundingClientRect();
  
  // Calculate relative Y within the element
  let relativeY = clientY - rect.top;
  // Clamp
  if (relativeY < 0) relativeY = 0;
  if (relativeY > rect.height) relativeY = rect.height;
  
  localPercent.value = (relativeY / rect.height) * 100;
};
</script>

<style scoped>
.zoomer {
  width: 30px; /* Wider hit area */
  height: 200px;
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 100;
  opacity: 0.8;
}

.zoomer-track {
  width: 4px;
  height: 100%;
  background: var(--term-muted);
  margin: 0 auto;
  border-radius: 0;
  position: relative;
  border: 1px solid var(--term-bg);
}

.zoomer-handle {
  width: 16px;
  height: 16px;
  background: var(--term-text);
  border-radius: 0; /* Square/Pixelated */
  position: absolute;
  left: 50%;
  transform: translate(-50%, -50%);
  transition: top 0.1s;
  box-shadow: var(--term-glow);
}

.zoomer-tooltip {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  background: var(--term-bg);
  color: var(--term-text);
  padding: 4px 8px;
  border: 1px solid var(--term-muted);
  font-family: var(--term-font);
  font-size: 12px;
}
</style>
