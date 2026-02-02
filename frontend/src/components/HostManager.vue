<template>
  <div class="host-manager">
    <TerminalCard title="Available Hosts">
      <div class="header">
        <span class="info-text">SELECT TARGET SYSTEM:</span>
        <TerminalButton @click="openAddModal">New Host</TerminalButton>
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
const isEditMode = ref(false);
const loading = ref(false);
const testing = ref(false);
const testStatus = ref<{type: string, msg: string} | null>(null);

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

.icon-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    font-family: var(--term-font);
    color: inherit;
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
</style>
