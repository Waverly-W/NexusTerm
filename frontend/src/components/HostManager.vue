<template>
  <div class="host-manager">
    <!-- Top Bar with Settings Icon -->
    <div class="top-bar">
        <button class="icon-btn settings-btn" @click="openSettings" title="Settings">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="3"></circle>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
            </svg>
        </button>
    </div>

    <TerminalCard title="Available Hosts">
      <div class="header">
        <span class="info-text">SELECT TARGET SYSTEM:</span>
        <div class="header-actions">
           <!-- Refresh Button (Icon) -->
           <button class="icon-btn" @click="fetchHosts" title="Refresh List">
             <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M23 4v6h-6"></path>
                <path d="M1 20v-6h6"></path>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
             </svg>
           </button>
           
           <!-- Add Host Button (Icon, Right of Refresh) -->
           <button class="icon-btn" @click="openAddModal" title="Add Host">
             <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
             </svg>
           </button>
        </div>
      </div>

      <!-- Host List -->
      <div class="host-list">
        <div v-if="loading && !hosts.length" class="loading">SCANNING NETWORKS...</div>
        <div 
          v-for="h in hosts" 
          :key="h.id" 
          class="host-item"
          @click="$emit('connect', h)"
        >
          <div class="host-icon">
            >_
          </div>
          <div class="host-info">
            <div class="host-alias">{{ h.alias }}</div>
            <div class="host-detail">{{ h.username }}@{{ h.hostname }}</div>
          </div>
          <div class="host-actions">
             <button class="icon-btn edit-btn" @click.stop="editHost(h)" title="Edit">
               [EDIT]
             </button>
          </div>
        </div>
        <div v-if="!hosts.length && !loading" class="empty">NO HOSTS FOUND IN REGISTRY.</div>
      </div>
    </TerminalCard>

    <!-- Settings Modal -->
    <div v-if="showSettings" class="modal-backdrop" @click.self="closeSettings">
                <div class="modal-body">
                    <div class="setting-group">
                        <div class="setting-item">
                            <label>Virtual Keyboard</label>
                            <div class="toggle-switch">
                                <span class="toggle-label">{{ useVirtualKeyboard ? 'ON' : 'OFF' }}</span>
                                <TerminalButton 
                                   :variant="useVirtualKeyboard ? 'primary' : 'secondary'"
                                   @click="toggleVirtualKeyboard"
                                   class="mini-btn"
                                >
                                   {{ useVirtualKeyboard ? '[X]' : '[ ]' }}
                                </TerminalButton>
                            </div>
                        </div>
                        <div class="setting-desc">
                            Use built-in virtual keyboard for mobile input.
                        </div>
                    </div>
					
                    <div class="setting-group">
                        <div class="setting-item">
                            <label>Max Keep-Alive (min)</label>
                            <input 
                                v-model="keepAliveTime" 
                                class="simple-input" 
                                type="number" 
                                placeholder="0"
                            />
                        </div>
                        <div class="setting-desc">
                            Session active time after backgrounding (0 = disable).
                        </div>
                    </div>

                    <div class="setting-group">
                        <div class="setting-item">
                            <label>Login Timeout (min)</label>
                            <input 
                                v-model="loginTimeout" 
                                class="simple-input" 
                                type="number" 
                                placeholder="30"
                            />
                        </div>
                        <div class="setting-desc">
                            Auto-logout after inactivity (0 = disable).
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <TerminalButton @click="closeSettings">CLOSE</TerminalButton>
                </div>
            </TerminalCard>
        </div>
    </div>

    <!-- Add/Edit Host Modal -->
    <div v-if="showModal" class="modal-backdrop" @click.self="closeModal">
      <div class="modal-content">
        <TerminalCard :title="isEditMode ? 'Edit System Config' : 'New System Registration'">
          <div class="modal-body">
            <TerminalInput v-model="currentHost.alias" prompt="Alias:" />
            <div class="row">
              <TerminalInput v-model="currentHost.hostname" prompt="Host:" class="flex-2" />
              <TerminalInput v-model.number="currentHost.port" prompt="Port:" class="flex-1" />
            </div>
            <div class="row">
              <TerminalInput v-model="currentHost.username" prompt="User:" class="flex-1" />
              <TerminalInput type="password" v-model="currentHost.password" prompt="Pass:" class="flex-1" />
            </div>
          </div>
          
          <div v-if="testStatus" class="test-status" :class="testStatus.type">
              [{{ testStatus.type.toUpperCase() }}] {{ testStatus.msg }}
          </div>

          <div class="modal-footer">
              <TerminalButton variant="secondary" @click="testConnection" :disabled="testing || !canTest">
                   {{ testing ? 'PINGING...' : 'TEST CONN' }}
              </TerminalButton>
              <div class="spacer"></div>
              <TerminalButton variant="danger" @click="closeModal">CANCEL</TerminalButton>
              <TerminalButton @click="saveHost" :disabled="loading">SAVE CONFIG</TerminalButton>
          </div>
        </TerminalCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import TerminalCard from './ui/TerminalCard.vue';
import TerminalButton from './ui/TerminalButton.vue';
import TerminalInput from './ui/TerminalInput.vue';

const props = defineProps({
  token: String
});

const emit = defineEmits(['connect']);

const hosts = ref<any[]>([]);
const showModal = ref(false);
const showSettings = ref(false);
const isEditMode = ref(false);
const loading = ref(false);
const testing = ref(false);
const testStatus = ref<{type: string, msg: string} | null>(null);

// Settings
const useVirtualKeyboard = ref(localStorage.getItem('nexus_use_vk') !== 'false'); // Default TRUE
const keepAliveTime = ref(localStorage.getItem('nexus_keep_alive') || '0');
const loginTimeout = ref(localStorage.getItem('nexus_login_timeout') || '30');

const openSettings = () => {
    showSettings.value = true;
};

const closeSettings = () => {
    localStorage.setItem('nexus_keep_alive', keepAliveTime.value);
    localStorage.setItem('nexus_login_timeout', loginTimeout.value);
    localStorage.setItem('nexus_use_vk', String(useVirtualKeyboard.value));
    showSettings.value = false;
};

const toggleVirtualKeyboard = () => {
    useVirtualKeyboard.value = !useVirtualKeyboard.value;
    localStorage.setItem('nexus_use_vk', String(useVirtualKeyboard.value));
};

const currentHost = ref({
  id: 0,
  alias: '',
  hostname: '',
  port: 22,
  username: '',
  password: ''
});

const canTest = computed(() => {
    return currentHost.value.hostname && currentHost.value.username && currentHost.value.password;
});

const fetchHosts = async () => {
  loading.value = true;
  try {
    const res = await fetch('/api/hosts', {
      headers: { 'Authorization': `Bearer ${props.token}` }
    });
    if (res.ok) {
      const data = await res.json();
      if (Array.isArray(data)) {
        hosts.value = data;
      } else {
        hosts.value = [];
      }
    }
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

const openAddModal = () => {
    isEditMode.value = false;
    currentHost.value = { id: 0, alias: '', hostname: '', port: 22, username: '', password: '' };
    testStatus.value = null;
    showModal.value = true;
};

const editHost = (host: any) => {
    isEditMode.value = true;
    // Clone to avoid reactive mess
    currentHost.value = { ...host, password: '' }; // Password empty on edit
    testStatus.value = null;
    showModal.value = true;
};

const closeModal = () => {
    showModal.value = false;
};

const testConnection = async () => {
    testing.value = true;
    testStatus.value = null;
    try {
        const res = await fetch('/api/test-connection', {
             method: 'POST',
             headers: { 'Content-Type': 'application/json' },
             body: JSON.stringify({
                 host: currentHost.value.hostname,
                 port: currentHost.value.port,
                 user: currentHost.value.username,
                 password: currentHost.value.password
             })
        });
        if (res.ok) {
            testStatus.value = { type: 'success', msg: 'Connection Successful!' };
        } else {
             testStatus.value = { type: 'error', msg: 'Connection Failed: ' + await res.text() };
        }
    } catch(e) {
         testStatus.value = { type: 'error', msg: 'Error: ' + e };
    } finally {
        testing.value = false;
    }
};

const saveHost = async () => {
  if (!currentHost.value.hostname || !currentHost.value.username) return;
  
  // Create mode requires password
  if (!isEditMode.value && !currentHost.value.password) {
      alert("Password required for new host");
      return;
  }
  
  loading.value = true;
  const method = isEditMode.value ? 'PUT' : 'POST';
  
  try {
    const res = await fetch('/api/hosts', {
      method: method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${props.token}`
      },
      body: JSON.stringify(currentHost.value)
    });
    
    if (res.ok) {
      closeModal();
      fetchHosts();
    } else {
      alert(await res.text());
    }
  } catch (e) {
    alert('Failed to save host');
  } finally {
    loading.value = false;
  }
};

onMounted(fetchHosts);
</script>

<style scoped>
.host-manager {
  width: 100%;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}
.info-text {
  color: var(--term-muted);
}
.header-actions {
    display: flex;
    gap: 8px;
}
.setting-group {
    margin-bottom: 1.5rem;
}

.setting-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.25rem;
}

.setting-item label {
    color: var(--term-text);
    font-weight: bold;
}

.setting-desc {
    font-size: 0.8rem;
    color: var(--term-muted);
    line-height: 1.4;
}

.simple-input {
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--term-muted);
    color: var(--term-secondary);
    font-family: var(--term-font);
    font-size: 1rem;
    width: 60px;
    text-align: right;
    padding: 4px;
    outline: none;
}
.simple-input:focus {
    border-bottom-color: var(--term-text);
    color: var(--term-text);
}

.toggle-switch {
    display: flex;
    align-items: center;
    gap: 8px;
}
.toggle-label {
    font-size: 0.9rem;
    color: var(--term-secondary);
}
.mini-btn {
    padding: 2px 8px;
    min-width: 40px;
}

.host-list {
  display: flex;
  flex-direction: column;
  gap: 0;
  max-height: 400px;
  overflow-y: auto;
  border-top: 1px dashed var(--term-muted);
}

.host-item {
  padding: 0.75rem;
  border-bottom: 1px dashed var(--term-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.1s;
}
.host-item:hover {
  background: var(--term-text);
  color: var(--term-bg);
}
.host-item:hover .host-detail {
  color: var(--term-bg);
  opacity: 0.8;
}

.host-icon {
    font-weight: bold;
}

.host-info {
    flex: 1;
}
.host-alias {
  font-weight: bold;
  text-transform: uppercase;
}
.host-detail {
  color: var(--term-secondary);
  font-size: 0.9em;
}

.host-actions {
    opacity: 0;
}
.host-item:hover .host-actions {
    opacity: 1;
    color: var(--term-bg);
}

@media (max-width: 600px) {
    .header {
        /* Keep header horizontal on mobile now that buttons are icons (smaller) */
        flex-direction: row; 
        align-items: center;
        gap: 0;
    }
    .header-actions {
        width: auto;
        justify-content: flex-end;
    }

    .host-actions {
        opacity: 1; 
        color: var(--term-text);
    }
    .host-item:hover .host-actions {
        color: var(--term-bg);
    }
    .host-item {
        padding: 1rem;
        gap: 12px;
    }
}

.top-bar {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 8px;
}
.settings-btn {
    opacity: 0.7;
    transition: opacity 0.2s;
}
.settings-btn:hover {
    opacity: 1;
}

.icon-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    font-family: var(--term-font);
    color: inherit;
    padding: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
}
.icon-btn:hover {
    background: rgba(255, 255, 255, 0.1);
}

/* Modal */
.modal-backdrop {
    position: fixed;
    top: 0; 
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0,0,0,0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}
.modal-content {
    width: 90%;
    max-width: 500px;
}

.modal-body {
    margin-bottom: 1rem;
}

.test-status {
    font-size: 13px;
    margin-bottom: 10px;
    padding: 8px;
    border: 1px solid currentColor;
}
.test-status.success {
    color: var(--term-text);
}
.test-status.error {
    color: var(--term-error);
}

.modal-footer {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px dashed var(--term-muted);
}
.spacer { flex: 1; }

.row {
  display: flex;
  gap: 12px;
}
.flex-1 { flex: 1; }
.flex-2 { flex: 2; }
.setting-input { width: 80px; text-align: right; }
</style>
