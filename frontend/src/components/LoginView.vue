<template>
  <div class="login-view">
    <TerminalCard title="SYSTEM ACCESS" class="login-card">
      <div class="tabs">
        <TerminalButton 
          :variant="authMode === 'login' ? 'primary' : 'secondary'"
          @click="authMode = 'login'"
        >
          AUTH
        </TerminalButton>
        <div class="divider">|</div>
        <TerminalButton 
           :variant="authMode === 'register' ? 'primary' : 'secondary'"
           @click="authMode = 'register'"
        >
           INIT
        </TerminalButton>
      </div>
        
      <div class="form-body">
        <TerminalInput v-model="authForm.username" prompt="USER:" @keyup.enter="performAuth" />
        <TerminalInput type="password" v-model="authForm.password" prompt="PASS:" @keyup.enter="performAuth" />
      </div>

      <div class="footer">
        <TerminalButton @click="performAuth" :disabled="loading" class="full-width">
          {{ loading ? 'PROCESSING...' : (authMode === 'login' ? 'LOGIN' : 'REGISTER') }}
        </TerminalButton>
      </div>

      <div v-if="error" class="error">[ERR] {{ error }}</div>
    </TerminalCard>
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
        return;
    }

    loading.value = true;
    try {
        const formData = new FormData();
        formData.append('username', authForm.value.username);
        formData.append('password', authForm.value.password);
        
        const endpoint = authMode.value === 'login' ? '/api/login' : '/api/register';
        const baseUrl = 'http://localhost:8080';
        
        const res = await fetch(baseUrl + endpoint, {
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
  background: var(--term-bg);
  color: var(--term-text);
}

.login-card {
  width: 100%;
  max-width: 400px;
}

.tabs {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 2rem;
  gap: 1rem;
}

.divider {
    color: var(--term-muted);
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
    justify-content: center;
}

.error {
  color: var(--term-error);
  margin-top: 1rem;
  text-align: center;
}
</style>
