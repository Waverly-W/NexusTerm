<template>
  <div class="login-view">
    <div class="login-container animate-fade-in">
        <div class="brand-header">
            <h1 class="brand-title">NEXUS<span class="text-primary">TERM</span></h1>
            <p class="brand-subtitle">SECURE REMOTE ACCESS GATEWAY</p>
        </div>

        <TerminalCard title="AUTHENTICATION REQUIRED" class="login-card">
          <div class="tabs">
            <TerminalButton 
              :variant="authMode === 'login' ? 'primary' : 'secondary'"
              @click="authMode = 'login'"
              class="tab-btn"
            >
              LOGIN
            </TerminalButton>
            <div class="divider"></div>
            <TerminalButton 
               :variant="authMode === 'register' ? 'primary' : 'secondary'"
               @click="authMode = 'register'"
               class="tab-btn"
            >
               REGISTER
            </TerminalButton>
          </div>
            
          <div class="form-body">
            <TerminalInput 
                v-model="authForm.username" 
                prompt="USER" 
                @keyup.enter="performAuth"
                placeholder="Enter username" 
            />
            <TerminalInput 
                type="password" 
                v-model="authForm.password" 
                prompt="PASS" 
                @keyup.enter="performAuth" 
                placeholder="Enter password"
            />
          </div>
    
          <div class="footer">
            <TerminalButton @click="performAuth" :disabled="loading" class="full-width icon-btn">
              <span v-if="loading">PROCESSING...</span>
              <span v-else>{{ authMode === 'login' ? 'ACCESS SYSTEM' : 'INITIALIZE USER' }}</span>
            </TerminalButton>
          </div>
    
          <div v-if="error" class="error">
              <span class="error-icon">!</span>
              {{ error }}
          </div>
        </TerminalCard>
        
        <div class="login-footer text-muted">
            v1.0.0 &bullet; SECURE CONNECTION ESTABLISHED
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import TerminalCard from './ui/TerminalCard.vue';
import TerminalButton from './ui/TerminalButton.vue';
import TerminalInput from './ui/TerminalInput.vue';

const emit = defineEmits(['success']);

const authMode = ref<'login' | 'register'>('login');
const authForm = ref({ username: '', password: '' });
const error = ref('');
const loading = ref(false);

const performAuth = async () => {
    error.value = '';
    
    if (!authForm.value.username || !authForm.value.password) {
        error.value = 'MISSING CREDENTIALS';
        // Add shake animation logic here if easier
        return;
    }

    loading.value = true;
    try {
        const formData = new FormData();
        formData.append('username', authForm.value.username);
        formData.append('password', authForm.value.password);
        
        const endpoint = authMode.value === 'login' ? '/api/login' : '/api/register';
        
        const res = await fetch(endpoint, {
            method: 'POST',
            body: formData
        });
        
        if (!res.ok) throw new Error(await res.text());
        
        if (authMode.value === 'login') {
            const token = await res.text();
            emit('success', token);
        } else {
            authMode.value = 'login';
            error.value = 'REGISTRATION SUCCESS. PLEASE LOGIN.';
            authForm.value.password = ''; // Clear password for safety
        }
    } catch (e: any) {
        error.value = e.message;
    } finally {
        loading.value = false;
    }
};
</script>

<style scoped>
.login-view {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  height: 100dvh;
  padding: 1rem;
  box-sizing: border-box;
}

.login-container {
    width: 100%;
    max-width: 420px;
    display: flex;
    flex-direction: column;
    gap: 2rem;
}

.brand-header {
    text-align: center;
}

.brand-title {
    font-size: 2.5rem;
    margin: 0;
    line-height: 1;
    letter-spacing: -0.05em;
    text-shadow: 0 0 20px rgba(0, 255, 157, 0.2);
}

.brand-subtitle {
    font-size: 0.8rem;
    color: var(--term-text-muted);
    letter-spacing: 0.4em;
    margin-top: 0.5rem;
}

.tabs {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 2rem;
  background: rgba(0,0,0,0.2);
  padding: 4px;
  border-radius: 6px;
  border: 1px solid var(--term-surface-border);
}

.tab-btn {
    flex: 1;
}

.divider {
    width: 1px;
    height: 20px;
    background: var(--term-surface-border);
    margin: 0 8px;
}

.form-body {
    margin-bottom: 2rem;
}

.footer {
    display: flex;
    justify-content: center;
}

.full-width {
    width: 100%;
}

.error {
  color: var(--term-error);
  margin-top: 1rem;
  text-align: center;
  font-size: 0.9rem;
  padding: 0.5rem;
  background: rgba(255, 51, 51, 0.1);
  border: 1px solid rgba(255, 51, 51, 0.2);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.error-icon {
    font-weight: bold;
    display: inline-block;
    width: 20px; height: 20px;
    line-height: 20px;
    background: var(--term-error);
    color: #fff;
    border-radius: 50%;
    font-size: 14px;
}

.login-footer {
    text-align: center;
    font-size: 0.75rem;
    opacity: 0.5;
}
</style>
