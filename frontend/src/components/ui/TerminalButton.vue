<template>
  <button 
    class="term-btn" 
    :class="[`variant-${variant}`, { 'is-pushed': isPushed }]"
    @mousedown="isPushed = true"
    @mouseup="isPushed = false"
    @mouseleave="isPushed = false"
    @touchstart="isPushed = true"
    @touchend="isPushed = false"
  >
    <span class="btn-content"><slot></slot></span>
    <div class="btn-glitch" aria-hidden="true"><slot></slot></div>
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
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  color: var(--term-text);
  border: 1px solid var(--term-surface-border);
  font-family: var(--term-font);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  padding: 0.6rem 1.2rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  transition: all 0.2s var(--ease-out);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  user-select: none;
  overflow: hidden;
  border-radius: 4px; /* Slight rounding */
}

.term-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: var(--term-primary);
  color: var(--term-primary);
  box-shadow: 0 0 15px var(--term-primary-dim);
  text-shadow: 0 0 8px var(--term-primary-dim);
  transform: translateY(-2px);
}

/* Active State (Click) */
.is-pushed {
  transform: translateY(1px) !important;
  box-shadow: inset 0 0 10px rgba(0,0,0,0.5);
  opacity: 0.8;
}

/* Variants */
.variant-primary {
  border-color: var(--term-primary);
  color: var(--term-primary);
}

.variant-primary:hover {
    background: var(--term-primary);
    color: #000; /* Contrast text on bright bg */
    box-shadow: 0 0 20px rgba(0, 255, 157, 0.4);
    text-shadow: none;
}

.variant-secondary {
  color: var(--term-text-muted);
  border-color: transparent;
}
.variant-secondary:hover {
    color: var(--term-text);
    border-color: var(--term-text);
     box-shadow: none;
}

.variant-danger {
  color: var(--term-error);
  border-color: rgba(255, 51, 51, 0.3);
}
.variant-danger:hover {
    background: var(--term-error);
    color: #fff;
    box-shadow: 0 0 15px rgba(255, 51, 51, 0.4);
}

/* Glitch Effect (Text Duplicate) */
.btn-glitch {
    position: absolute;
    top: 0; left: 0; width: 100%; height: 100%;
    display: flex; align-items: center; justify-content: center;
    background: var(--term-bg);
    color: var(--term-primary);
    opacity: 0;
    z-index: -1;
    transform: translate(-2px, 0);
    clip-path: polygon(0 0, 100% 0, 100% 45%, 0 45%);
    transition: opacity 0.1s;
}

.term-btn:hover .btn-glitch {
    opacity: 0.1;
    animation: glitch-anim 0.3s infinite linear alternate-reverse;
}

@keyframes glitch-anim {
  0% { transform: translate(-2px, -1px); }
  20% { transform: translate(2px, 1px); }
  40% { transform: translate(-2px, 1px); }
  60% { transform: translate(2px, -1px); }
  80% { transform: translate(-1px, 2px); }
  100% { transform: translate(1px, -2px); }
}

/* Mobile Touch Optimization */
@media (hover: none) {
    .term-btn:hover {
        transform: none;
        box-shadow: none;
    }
}
</style>
