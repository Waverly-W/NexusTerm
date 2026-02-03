<template>
  <div class="virtual-keyboard">
    <!-- Row 1: Numbers -->
    <div class="kb-row">
      <button v-for="k in ['1','2','3','4','5','6','7','8','9','0']" :key="k" @click="$emit('key', k)">{{k}}</button>
    </div>
    <!-- Row 2: QWERTY -->
    <div class="kb-row">
      <button v-for="k in ['q','w','e','r','t','y','u','i','o','p']" :key="k" @click="$emit('key', k)">{{k}}</button>
    </div>
    <!-- Row 3: ASDF -->
    <div class="kb-row">
      <button v-for="k in ['a','s','d','f','g','h','j','k','l']" :key="k" @click="$emit('key', k)">{{k}}</button>
       <button class="action-btn" @click="$emit('key', 'Enter')">Ent</button>
    </div>
    <!-- Row 4: ZXCV -->
    <div class="kb-row">
       <button class="action-btn" :class="{active: shiftActive}" @click="$emit('toggle-shift')">Shft</button>
       <button v-for="k in ['z','x','c','v','b','n','m']" :key="k" @click="$emit('key', k)">{{k}}</button>
       <button class="action-btn" @click="$emit('key', 'Backspace')">Del</button>
    </div>
    <!-- Row 5: Symbols & Space -->
    <div class="kb-row">
       <button v-for="k in ['-','/','|','_','.',',']" :key="k" @click="$emit('key', k)">{{k}}</button>
       <button class="space-btn" @click="$emit('key', ' ')">SPACE</button>
    </div>
    <!-- Row 6: Terminal Specials -->
    <div class="kb-row specials">
       <button class="special-btn" @click="$emit('key', 'Escape')">Esc</button>
       <button class="special-btn" @click="$emit('key', 'Tab')">Tab</button>
       <button class="special-btn" :class="{active: ctrlActive}" @click="$emit('toggle-ctrl')">Ctrl</button>
       <button class="special-btn" @click="$emit('key', 'Alt')">Alt</button>
       <button class="special-btn" @click="$emit('key', 'ArrowUp')">↑</button>
       <button class="special-btn" @click="$emit('key', 'ArrowDown')">↓</button>
       <button class="special-btn" @click="$emit('key', 'ArrowLeft')">←</button>
       <button class="special-btn" @click="$emit('key', 'ArrowRight')">→</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  ctrlActive: boolean,
  shiftActive: boolean
}>();

defineEmits<{
  (e: 'key', key: string): void,
  (e: 'toggle-ctrl'): void,
  (e: 'toggle-shift'): void
}>();
</script>

<style scoped>
.virtual-keyboard {
  background: var(--term-bg);
  border-top: 1px solid var(--term-muted);
  padding: 4px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  user-select: none;
  /* Prevent touch actions like zooming/scrolling on the keyboard itself */
  touch-action: none; 
}

.kb-row {
  display: flex;
  gap: 4px;
  justify-content: center;
}

button {
  flex: 1;
  background: #111;
  border: 1px solid var(--term-muted);
  color: var(--term-text);
  font-family: var(--term-font);
  font-size: 16px;
  min-height: 40px;
  min-width: 25px;
  border-radius: 2px;
  text-transform: uppercase;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

button:active {
  background: var(--term-text);
  color: var(--term-bg);
}

.active {
  background: var(--term-text);
  color: var(--term-bg);
  box-shadow: 0 0 5px var(--term-text);
}

.action-btn {
  flex: 1.5;
  font-size: 12px;
  background: #222;
}

.space-btn {
  flex: 3;
}

.specials {
  margin-top: 2px;
  padding-top: 4px;
  border-top: 1px dashed var(--term-muted);
}
.special-btn {
  font-size: 12px;
  background: #222;
}
</style>
