<template>
  <div class="app" :class="{ 'dark-theme': isDark }">
    <div class="app-content">
      <router-view />
    </div>
    <Navbar v-if="showNavbar" />
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import Navbar from './components/Navbar.vue'

export default {
  name: 'App',
  components: { Navbar },
  setup() {
    const route = useRoute()
    const isDark = ref(false)
    
    // Проверяем тему Telegram
    if (window.Telegram?.WebApp) {
      isDark.value = window.Telegram.WebApp.colorScheme === 'dark'
      
      // Слушаем изменения темы
      window.Telegram.WebApp.onEvent('themeChanged', () => {
        isDark.value = window.Telegram.WebApp.colorScheme === 'dark'
      })
    }
    
    // Показывать навбар не на всех страницах
    const showNavbar = computed(() => {
      const hideOnPages = ['TripDetails']
      return !hideOnPages.includes(route.name)
    })
    
    return {
      isDark,
      showNavbar
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app {
  min-height: 100vh;
  background: #f5f5f7;
  transition: background 0.3s;
}

.app.dark-theme {
  background: #000;
  color: #fff;
}

.app-content {
  padding-bottom: 60px; /* Для навбара */
}
</style>