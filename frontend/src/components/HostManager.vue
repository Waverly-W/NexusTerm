<template>
  <div class="host-manager animate-fade-in">
    <!-- Top Bar with Settings Icon -->
    <div class="top-bar">
        <div class="user-status">
            <span class="status-indicator online"></span>
            CONNECTED AS <span class="text-primary">ADMIN</span>
        </div>
        <button class="icon-btn settings-btn" @click="openSettings" title="Settings">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="3"></circle>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
            </svg>
        </button>
    </div>

    <TerminalCard title="Network Targets">
      <div class="header">
        <span class="info-text">AVAILABLE HOST ENDPOINTS:</span>
        <div class="header-actions">
           <!-- Refresh Button -->
           <button class="icon-btn" @click="fetchHosts" title="Refresh List">
             <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M23 4v6h-6"></path>
                <path d="M1 20v-6h6"></path>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
             </svg>
           </button>
           
           <!-- Add Host Button -->
           <TerminalButton variant="primary" @click="openAddModal" class="small-btn">
             + NEW HOST
           </TerminalButton>
        </div>
      </div>

      <!-- Host List -->
      <div class="host-list-container">
        <div class="host-list-header">
            <div class="col-icon">#</div>
            <div class="col-alias">ALIAS</div>
            <div class="col-address">ADDRESS</div>
            <div class="col-actions">ACTIONS</div>
        </div>
        
        <div class="host-list scrollbar-custom">
            <div v-if="loading && !hosts.length" class="loading-state">
                <span class="spinner"></span> SCANNING NETWORK...
            </div>
            
            <div 
              v-for="(h, index) in hosts" 
              :key="h.id" 
              class="host-item"
              @click="$emit('connect', h)"
            >
              <div class="col-icon text-muted">{{ index + 1 }}</div>
              <div class="col-alias">
                  <span class="host-alias-text">{{ h.alias }}</span>
              </div>
              <div class="col-address text-muted">
                  {{ h.username }}<span class="at-symbol">@</span>{{ h.hostname }}
              </div>
              <div class="col-actions">
                  <button class="action-btn connect-btn" title="Quick Connect">
                    CONNECT
                  </button>
                 <button class="icon-btn edit-btn" @click.stop="editHost(h)" title="Edit Configuration">
                   <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polygon points="16 3 21 8 8 21 3 21 3 16 16 3"></polygon>
                   </svg>
                 </button>
              </div>
            </div>
            
            <div v-if="!hosts.length && !loading" class="empty-state">
                NO HOSTS CONFIGURED. ADD A NEW ENDPOINT.
            </div>
        </div>
      </div>
    </TerminalCard>

    <!-- Settings Modal -->
    <div v-if="showSettings" class="modal-backdrop glass-panel" @click.self="closeSettings">
        <div class="modal-content glass-panel">
            <div class="modal-header">
                <h3>SYSTEM PREFERENCES</h3>
                <button class="close-btn" @click="closeSettings">×</button>
            </div>
            
            <div class="modal-body">
                <div class="setting-group">
                    <div class="setting-item">
                        <label>VIRTUAL KEYBOARD</label>
                        <div class="toggle-switch">
                            <span class="toggle-label">{{ useVirtualKeyboard ? 'ENABLED' : 'DISABLED' }}</span>
                            <button 
                                class="switch-btn"
                                :class="{ active: useVirtualKeyboard }"
                                @click="toggleVirtualKeyboard"
                            >
                                <div class="switch-knob"></div>
                            </button>
                        </div>
                    </div>
                    <div class="setting-desc">
                        Enable on-screen inputs for mobile devices.
                    </div>
                </div>
                
                <div class="setting-group">
                    <div class="setting-item">
                        <label>BACKGROUND KEEPALIVE (MIN)</label>
                        <TerminalInput 
                            v-model="keepAliveTime" 
                            type="number" 
                            class="setting-input-field" 
                            placeholder="0"
                        />
                    </div>
                    <div class="setting-desc">
                        Session persistence duration (0 = Disabled).
                    </div>
                </div>

                <div class="setting-group">
                    <div class="setting-item">
                        <label>AUTO-LOGOUT (MIN)</label>
                        <TerminalInput 
                            v-model="loginTimeout" 
                            type="number" 
                            class="setting-input-field" 
                            placeholder="30"
                        />
                    </div>
                    <div class="setting-desc">
                        Security timeout for inactivity.
                    </div>
                </div>
            </div>
            
            <div class="modal-footer">
                <TerminalButton variant="danger" @click="handleLogout">TERMINATE SESSION</TerminalButton>
                <div class="spacer"></div>
                <TerminalButton variant="secondary" @click="closeSettings">SAVE & CLOSE</TerminalButton>
            </div>
        </div>
    </div>

    <!-- Add/Edit Host Modal -->
    <div v-if="showModal" class="modal-backdrop glass-panel" @click.self="closeModal">
      <div class="modal-content glass-panel">
        <div class="modal-header">
            <h3>{{ isEditMode ? 'EDIT CONFIGURATION' : 'NEW ENDPOINT' }}</h3>
            <button class="close-btn" @click="closeModal">×</button>
        </div>
        
        <div class="modal-body">
            <div class="form-grid">
                <TerminalInput v-model="currentHost.alias" prompt="ALIAS" class="full-col" placeholder="Friendly Name" />
                
                <TerminalInput v-model="currentHost.hostname" prompt="HOST" class="col-span-2" placeholder="IP / Domain" />
                <TerminalInput v-model.number="currentHost.port" prompt="PORT" class="col-span-1" type="number" />
                
                <TerminalInput v-model="currentHost.username" prompt="USER" class="col-span-2" placeholder="ssh user" />
                <TerminalInput type="password" v-model="currentHost.password" prompt="PASS" class="col-span-1" placeholder="******" />
            </div>
        </div>
          
        <div v-if="testStatus" class="test-status" :class="testStatus.type">
            <span class="status-icon">
                {{ testStatus.type === 'success' ? '✓' : '!' }}
            </span>
            {{ testStatus.msg }}
        </div>

        <div class="modal-footer">
            <TerminalButton variant="secondary" @click="testConnection" :disabled="testing || !canTest">
                 {{ testing ? 'PINGING...' : 'TEST CONNECTION' }}
            </TerminalButton>
            <div class="spacer"></div>
            <TerminalButton @click="saveHost" :disabled="loading">
                {{ loading ? 'SAVING...' : 'CONFIRM CONFIG' }}
            </TerminalButton>
        </div>
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

const handleLogout = () => {
    if (confirm("Confirm System Logout?")) {
        localStorage.removeItem('nexus_token');
        window.location.reload();
    }
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
        hosts.value = Array.isArray(data) ? data : [];
    } else if (res.status === 401) {
        window.location.reload();
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
    currentHost.value = { ...host, password: '' }; 
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
            testStatus.value = { type: 'success', msg: 'Connection Successful' };
        } else {
             testStatus.value = { type: 'error', msg: 'Failed: ' + await res.text() };
        }
    } catch(e) {
         testStatus.value = { type: 'error', msg: 'Error: ' + e };
    } finally {
        testing.value = false;
    }
};

const saveHost = async () => {
  if (!currentHost.value.hostname || !currentHost.value.username) return;
  if (!isEditMode.value && !currentHost.value.password) {
      alert("Password required");
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
      if (res.status === 401) window.location.reload();
      else alert(await res.text());
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

.top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding: 0 0.5rem;
}

.user-status {
    font-size: 0.8rem;
    color: var(--term-text-muted);
    display: flex;
    align-items: center;
    gap: 8px;
}

.status-indicator {
    width: 8px; height: 8px;
    background: var(--term-error);
    border-radius: 50%;
    box-shadow: 0 0 5px var(--term-error);
}
.status-indicator.online {
    background: var(--term-success);
    box-shadow: 0 0 5px var(--term-success);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}
.info-text {
  color: var(--term-text-muted);
  font-size: 0.85rem;
  letter-spacing: 0.05em;
}
.header-actions {
    display: flex;
    gap: 12px;
    align-items: center;
}

/* Host List Grid */
.host-list-container {
    border: 1px solid var(--term-surface-border);
    border-radius: 4px;
    background: rgba(0,0,0,0.2);
    overflow: hidden;
}

.host-list-header {
    display: flex;
    padding: 0.75rem 1rem;
    background: rgba(255,255,255,0.05);
    border-bottom: 1px solid var(--term-surface-border);
    font-size: 0.75rem;
    font-weight: bold;
    color: var(--term-text-muted);
}

.host-item {
    display: flex;
    padding: 1rem;
    border-bottom: 1px solid var(--term-surface-border);
    cursor: pointer;
    align-items: center;
    transition: background 0.1s;
}

.host-item:hover {
    background: rgba(255,255,255,0.05);
}

.host-item:last-child {
    border-bottom: none;
}

/* Columns */
.col-icon { width: 40px; text-align: center; }
.col-alias { flex: 2; font-weight: bold; color: var(--term-primary); }
.col-address { flex: 3; font-family: monospace; font-size: 0.9rem; }
.col-actions { flex: 2; display: flex; justify-content: flex-end; gap: 8px; }

.at-symbol { color: var(--term-text-muted); margin: 0 2px; }

.action-btn {
    background: transparent;
    border: 1px solid var(--term-primary);
    color: var(--term-primary);
    font-family: var(--term-font);
    font-size: 0.7rem;
    padding: 4px 8px;
    cursor: pointer;
    border-radius: 2px;
    text-transform: uppercase;
}
.action-btn:hover {
    background: var(--term-primary);
    color: #000;
}

.icon-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    color: var(--term-text-muted);
    padding: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
}
.icon-btn:hover {
    color: var(--term-text);
}

/* Modals */
.modal-backdrop {
    position: fixed;
    top: 0; left: 0; width: 100%; height: 100%;
    z-index: 2000;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0,0,0,0.85); /* Darker backdrop */
    backdrop-filter: blur(5px);
}

.modal-content {
    width: 90%;
    max-width: 500px;
    background: #0f0f0f;
    border: 1px solid var(--term-surface-border);
    box-shadow: 0 20px 50px rgba(0,0,0,0.8);
    display: flex;
    flex-direction: column;
    max-height: 90vh;
}

.modal-header {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--term-surface-border);
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.modal-header h3 { margin: 0; font-size: 1rem; }

.close-btn {
    background: transparent;
    border: none;
    color: var(--term-text-muted);
    font-size: 1.5rem;
    cursor: pointer;
    line-height: 1;
}
.close-btn:hover { color: var(--term-text); }

.modal-body {
    padding: 1.5rem;
    overflow-y: auto;
}

.modal-footer {
    padding: 1rem 1.5rem;
    border-top: 1px solid var(--term-surface-border);
    background: rgba(255,255,255,0.02);
    display: flex;
    align-items: center;
    gap: 1rem;
}

.form-grid {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 1rem;
}
.full-col { grid-column: span 3; }
.col-span-2 { grid-column: span 2; }
.col-span-1 { grid-column: span 1; }

.spacer { flex: 1; }

/* Switch */
.toggle-switch {
    display: flex;
    align-items: center;
    gap: 12px;
}
.switch-btn {
    width: 40px;
    height: 20px;
    background: var(--term-surface-border);
    border: 1px solid var(--term-text-muted);
    border-radius: 10px;
    position: relative;
    cursor: pointer;
    transition: all 0.2s;
}
.switch-knob {
    width: 14px; height: 14px;
    background: var(--term-text-muted);
    border-radius: 50%;
    position: absolute;
    top: 2px; left: 2px;
    transition: all 0.2s;
}
.switch-btn.active {
    border-color: var(--term-primary);
    background: rgba(0, 255, 157, 0.1);
}
.switch-btn.active .switch-knob {
    left: 22px;
    background: var(--term-primary);
    box-shadow: 0 0 5px var(--term-primary);
}

.test-status {
    margin: 0 1.5rem 1rem;
    padding: 0.75rem;
    border: 1px solid;
    border-radius: 4px;
    font-size: 0.9rem;
    display: flex;
    align-items: center;
    gap: 8px;
}
.test-status.success { color: var(--term-success); border-color: var(--term-success); background: rgba(0,255,157,0.1); }
.test-status.error { color: var(--term-error); border-color: var(--term-error); background: rgba(255,51,51,0.1); }

@media (max-width: 600px) {
    .host-list-header { display: none; }
    .host-item {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
    }
    .col-icon { display: none; }
    .col-actions { width: 100%; margin-top: 8px; }
}
</style>
