<template>
  <div 
    class="scrubber"
    @touchstart="onTouchStart"
    @touchmove="onTouchMove"
    @touchend="onTouchEnd"
  >
    <div 
      class="scrubber-track" 
      :style="{ transform: `translateX(${visualOffset}px)`, transition: isDragging ? 'none' : 'transform 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275)' }"
    >
      <div class="scrubber-label">SCRUBBER</div>
      <div class="scrubber-pattern"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const emit = defineEmits<{
  (e: 'scroll', direction: 'left' | 'right'): void
}>();

const isDragging = ref(false);
const visualOffset = ref(0);

let lastX = 0;
let accumulator = 0;
const STEP_SIZE = 15; // Pixels per cursor move

const onTouchStart = (e: TouchEvent) => {
  isDragging.value = true;
  if (e.touches[0]) {
    lastX = e.touches[0].clientX;
  }
  accumulator = 0;
  visualOffset.value = 0;
};

const onTouchMove = (e: TouchEvent) => {
  if (e.cancelable) e.preventDefault();
  
  if (!e.touches[0]) return;
  const currentX = e.touches[0].clientX;
  const delta = currentX - lastX;
  lastX = currentX; // Update reference for next frame
  
  accumulator += delta;

  // Visual Feedback: elasticity based on accumulator
  // Clamp visual offset to avoid drifting too far visually
  visualOffset.value = Math.max(-40, Math.min(40, accumulator * 0.5));

  // Event Logic: Consume accumulator
  while (Math.abs(accumulator) >= STEP_SIZE) {
      const direction = accumulator > 0 ? 'right' : 'left';
      emit('scroll', direction);
      
      // Haptic feedback
      if (typeof navigator !== 'undefined' && navigator.vibrate) {
          navigator.vibrate(5);
      }

      // Reduce accumulator by one step, preserving the remainder
      if (accumulator > 0) {
          accumulator -= STEP_SIZE;
      } else {
          accumulator += STEP_SIZE;
      }
  }
};

const onTouchEnd = () => {
  isDragging.value = false;
  visualOffset.value = 0;
  accumulator = 0;
};
</script>

<style scoped>
.scrubber {
  width: 100%;
  height: 48px;
  background: var(--term-bg);
  border-top: 1px solid var(--term-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  touch-action: none;
  user-select: none;
  overflow: hidden;
}

.scrubber-track {
  width: 60%;
  height: 24px;
  background: #111;
  border-radius: 0;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: inset 0 0px 3px rgba(51, 255, 0, 0.2);
  border: 1px solid var(--term-muted);
}

.scrubber-label {
  font-family: var(--term-font);
  font-size: 10px;
  font-weight: bold;
  color: var(--term-muted);
  letter-spacing: 1px;
  z-index: 2;
}

.scrubber-pattern {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 0;
  opacity: 0.1;
  background-image: repeating-linear-gradient(
    90deg,
    transparent,
    transparent 2px,
    var(--term-text) 2px,
    var(--term-text) 3px
  );
}
</style>
