<template>
  <div class="terminal-view">
    <!-- Disconnected State: Show Host List or Manual Connect -->
    <div v-if="!isConnected" class="connect-overlay">
       <HostManager :token="token" @connect="handleHostConnect" />
        
       <div class="divider">
          <!-- Text Divider -->
          <span class="divider-text">--- OR ---</span>
       </div>

       <div class="manual-section">
          <TerminalButton variant="secondary" @click="showManual = !showManual">
            {{ showManual ? 'Hide Manual Connection' : 'Connect Manually' }}
          </TerminalButton>
       </div>

       <TerminalCard v-if="showManual" title="Manual Connection" class="manual-card">
          <div class="row">
            <TerminalInput v-model="form.host" prompt="Host:" class="flex-2" />
            <TerminalInput v-model="form.user" prompt="User:" class="flex-1" />
          </div>
          <div class="row">
            <TerminalInput type="password" v-model="form.password" prompt="Pass:" class="flex-2" />
            <TerminalButton @click="connect()" class="flex-1">Connect</TerminalButton>
          </div>
       </TerminalCard>

       <div v-if="error" class="error">[ERROR] {{ error }}</div>
    </div>
    
    <!-- Connected State: Terminal -->
    <div 
      v-show="isConnected"
      class="terminal-wrapper" 
      :style="{ height: terminalHeight + 'px' }"
    >
      <div ref="terminalContainer" class="xterm-container"></div>
      <Zoomer 
        v-if="isConnected" 
        :current-size="currentFontSize"
        @change="handleFontSizeChange" 
      />
    </div>

    <div v-if="isConnected" class="mobile-controls">
      <SmartToolbar 
        :ctrl-active="ctrlActive"
        @key="handleToolbarKey"
        @input="handleToolbarInput"
        @toggle-ctrl="ctrlActive = !ctrlActive"
      />
      <Scrubber @scroll="handleScrubberLeft" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import 'xterm/css/xterm.css';

import Scrubber from './MobileControls/Scrubber.vue';
import Zoomer from './MobileControls/Zoomer.vue';
import SmartToolbar from './MobileControls/SmartToolbar.vue';
import HostManager from './HostManager.vue';
import TerminalButton from './ui/TerminalButton.vue';
import TerminalCard from './ui/TerminalCard.vue';
import TerminalInput from './ui/TerminalInput.vue';

const props = defineProps({
  token: String,
  visible: Boolean
});

const emit = defineEmits(['update:title']);

const terminalContainer = ref<HTMLElement | null>(null);
const isConnected = ref(false);
const error = ref('');
const showManual = ref(false);
const form = ref({
  host: 'localhost',
  user: '',
  password: ''
});

const terminalHeight = ref(window.innerHeight);
const ctrlActive = ref(false);
const currentFontSize = ref(14);

let term: Terminal | null = null;
let socket: WebSocket | null = null;
let fitAddon: FitAddon | null = null;

// Watch visibility to fit terminal when tab becomes active
watch(() => props.visible, (newVal) => {
    if (newVal) {
        setTimeout(handleViewportResize, 50); // Small delay to allow layout
    }
});

const handleHostConnect = (host: any) => {
    emit('update:title', host.alias || host.hostname);
    connectWithHostID(host.id);
};

const initTerminal = () => {
  if (!terminalContainer.value) return;

  // Use CSS variables for theme
  const style = getComputedStyle(document.body);
  const bg = style.getPropertyValue('--term-bg').trim() || '#0a0a0a';
  const fg = style.getPropertyValue('--term-text').trim() || '#33ff00';

  term = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'JetBrains Mono, Fira Code, VT323, monospace',
    theme: {
      background: bg,
      foreground: fg,
      cursor: fg,
    },
    // Important for mobile: prevent native selection/zoom
    allowProposedApi: true, 
  });

  fitAddon = new FitAddon();
  term.loadAddon(fitAddon);
  term.loadAddon(new WebLinksAddon());

  term.open(terminalContainer.value);
  fitAddon.fit();
  
  // Handle resize
  window.addEventListener('resize', handleResize);
  // Handle visual viewport (keyboard)
  if (window.visualViewport) {
    window.visualViewport.addEventListener('resize', handleViewportResize);
  }
  
  term.onData((data) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      if (ctrlActive.value) {
        // Handle Ctrl modifier
        const charCode = data.toUpperCase().charCodeAt(0);
        if (charCode >= 65 && charCode <= 90) {
           data = String.fromCharCode(charCode - 64);
        }
        ctrlActive.value = false; // Auto turn off
      }
      socket.send(new TextEncoder().encode(data));
    }
  });

  // Initial fit
  handleViewportResize();
};

const handleResize = () => {
  if (!window.visualViewport) {
    terminalHeight.value = window.innerHeight;
    fitTerm();
  }
};

const handleViewportResize = () => {
  if (!window.visualViewport) return;
  const controlsHeight = isConnected.value ? 96 : 0;
  terminalHeight.value = window.visualViewport.height - controlsHeight;
  setTimeout(fitTerm, 50);
  window.scrollTo(0, 0);
};

const fitTerm = () => {
  if (fitAddon) {
    fitAddon.fit();
    if (term && socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({
        t: 'r',
        rows: term.rows,
        cols: term.cols
      }));
    }
  }
};

const handleScrubberLeft = (direction: 'left' | 'right') => {
  const key = direction === 'left' ? '\x1b[D' : '\x1b[C';
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(new TextEncoder().encode(key));
  }
};

const handleFontSizeChange = (newSize: number) => {
  if (!term) return;
  term.options.fontSize = newSize;
  currentFontSize.value = newSize;
  fitTerm();
};

const handleToolbarKey = (key: string) => {
    let data = '';
    if (key === 'Escape') data = '\x1b';
    if (key === 'Tab') data = '\t';
    if (data && socket && socket.readyState === WebSocket.OPEN) {
         socket.send(new TextEncoder().encode(data));
    }
};

const handleToolbarInput = (char: string) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
         socket.send(new TextEncoder().encode(char));
    }
};

const connectWithHostID = (hostID: number) => {
    connect(hostID);
};

const connect = (hostID?: number) => {
  error.value = '';
  
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  socket = new WebSocket(`${protocol}//${window.location.host}/api/ws?token=${props.token}`);
  socket.binaryType = 'arraybuffer'; 

  socket.onopen = () => {
    const payload: any = { t: 'connect' };
    
    if (hostID) {
        payload.host_id = hostID;
        payload.host = '';
        payload.user = '';
        payload.password = '';
    } else {
        payload.host = form.value.host;
        payload.user = form.value.user;
        payload.password = form.value.password;
    }
    
    socket?.send(JSON.stringify(payload));
  };

  socket.onmessage = (event) => {
    if (typeof event.data === 'string') {
      try {
        const msg = JSON.parse(event.data);
        if (msg.status === 'connected') {
          isConnected.value = true;
          setTimeout(handleViewportResize, 100);
        } else if (msg.error) {
          error.value = msg.error;
          socket?.close();
        }
      } catch (e) {
        console.error('Failed to parse:', event.data);
      }
    } else {
      if (term) {
        term.write(new Uint8Array(event.data));
      }
    }
  };

  socket.onclose = () => {
    isConnected.value = false;
    error.value = 'Connection closed';
    handleViewportResize();
  };
  
  socket.onerror = (e) => {
    console.error(e);
    error.value = 'Connection error';
  };
};

onMounted(() => {
  initTerminal();
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
  if (window.visualViewport) {
    window.visualViewport.removeEventListener('resize', handleViewportResize);
  }
  if (socket) socket.close();
  if (term) term.dispose();
});
</script>

<style scoped>
.terminal-view {
  width: 100%;
  height: 100%; /* Changed from 100vh to 100% to fit in TabManager */
  background: var(--term-bg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.terminal-wrapper {
  position: relative;
  width: 100%;
  flex: 1; /* Take remaining space */
  overflow: hidden;
}

.xterm-container {
  width: 100%;
  height: 100%;
}

/* Mobile Controls */
.mobile-controls {
  flex-shrink: 0;
  z-index: 100;
  background: #252526;
  width: 100%;
}

@media (min-width: 768px) {
  .mobile-controls {
     display: none;
  }
}

.connect-overlay {
  max-width: 600px;
  margin: 4rem auto;
  padding: 1rem;
  width: 90%;
}

.manual-section {
    text-align: center;
    margin-bottom: 1rem;
}
.manual-card {
    margin-top: 1rem;
}

.row {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
}
.flex-1 { flex: 1; }
.flex-2 { flex: 2; }

.divider {
  text-align: center;
  margin: 2rem 0;
  position: relative;
}
.divider-text {
  background: var(--term-bg);
  padding: 0 1rem;
  color: var(--term-muted);
}

.error {
  color: var(--term-error);
  margin-top: 1rem;
  text-align: center;
  font-weight: bold;
}
</style>
