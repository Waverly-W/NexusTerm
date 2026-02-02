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

let startX = 0;
let lastEmittedStep = 0;
const STEP_SIZE = 12; // Pixels per cursor move

const onTouchStart = (e: TouchEvent) => {
  isDragging.value = true;
  if (e.touches[0]) {
    startX = e.touches[0].clientX;
  }
  lastEmittedStep = 0;
};

const onTouchMove = (e: TouchEvent) => {
  if (e.cancelable) e.preventDefault();
  
  if (!e.touches[0]) return;
  const currentX = e.touches[0].clientX;
  const totalDelta = currentX - startX;
  
  // Visual Feedback: Apply resistance/damping as you pull further
  // Logarithmic damping for "rubber band" feel
  const sign = Math.sign(totalDelta);
  const absDelta = Math.abs(totalDelta);
  // Damping function: y = x / (1 + x/limit) * damping_factor
  // Or simple sqrt damping: sign * sqrt(abs) * factor
  visualOffset.value = sign * Math.pow(absDelta, 0.8);

  // Event Logic: Discrete steps based on raw linear distance
  const currentStep = Math.floor(Math.abs(totalDelta) / STEP_SIZE);
  
  if (currentStep > lastEmittedStep) {
    const stepsToEmit = currentStep - lastEmittedStep;
    const direction = totalDelta > 0 ? 'right' : 'left';
    
    for (let i = 0; i < stepsToEmit; i++) {
        emit('scroll', direction);
        // Haptic feedback if available (subtle tick)
        if (typeof navigator !== 'undefined' && navigator.vibrate) {
            navigator.vibrate(5);
        }
    }
    
    lastEmittedStep = currentStep;
  }
};

const onTouchEnd = () => {
  isDragging.value = false;
  visualOffset.value = 0;
  startX = 0;
  lastEmittedStep = 0;
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
