<template>
  <div class="register-form">
    <h2>📝 Регистрация</h2>
    
    <form @submit.prevent="handleRegister">
      <div class="form-group">
        <label for="name">Имя</label>
        <input
          id="name"
          type="text"
          v-model="form.name"
          required
          placeholder="Иван"
        />
      </div>

      <div class="form-group">
        <label for="phone">Номер телефона</label>
        <input
          id="phone"
          type="text"
          v-model="form.phone"
          required
          placeholder="+79991112233"
        />
      </div>

      <div class="form-group">
        <label for="password">Пароль</label>
        <input
          id="password"
          type="password"
          v-model="form.password"
          required
          placeholder="••••••"
        />
      </div>

      <button type="submit" class="submit-btn" :disabled="isSubmitting">
        {{ isSubmitting ? "Регистрация..." : "Создать аккаунт" }}
      </button>
    </form>

    <div v-if="error" class="error-message">{{ error }}</div>

    <div class="switch-auth">
      Уже есть аккаунт? 
      <button type="button" @click="$emit('switch-to-login')" class="link-btn">Войти</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue"
import { authAPI } from "../services/api"

const emit = defineEmits(["success", "switch-to-login"])

const isSubmitting = ref(false)
const error = ref("")
const form = ref({
  name: "",
  phone: "",
  password: "",
})

const handleRegister = async () => {
  error.value = ""
  isSubmitting.value = true

  try {
    await authAPI.register(form.value)
    await authAPI.login({
      phone: form.value.phone,
      password: form.value.password,
    })
    emit("success")
  } catch (err) {
    console.error("Registration error:", err)
    error.value = err.response?.data?.error || "Не удалось зарегистрироваться. Возможно, номер телефона уже занят."
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.register-form {
  background: #fff;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
h2 {
  margin-top: 0;
  margin-bottom: 24px;
  color: #333;
  text-align: center;
}
.form-group {
  margin-bottom: 20px;
}
.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #555;
}
input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  box-sizing: border-box;
}
.submit-btn {
  width: 100%;
  padding: 16px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.submit-btn:hover:not(:disabled) {
  background: #218838;
}
.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
.error-message {
  margin-top: 16px;
  padding: 12px;
  background: #f8d7da;
  color: #721c24;
  border-radius: 8px;
  font-size: 14px;
}
.switch-auth {
  margin-top: 24px;
  text-align: center;
  font-size: 14px;
  color: #666;
}
.link-btn {
  background: none;
  border: none;
  color: #007bff;
  font-weight: bold;
  cursor: pointer;
  text-decoration: underline;
  padding: 0;
  font-size: 14px;
  margin-left: 5px;
}
</style>