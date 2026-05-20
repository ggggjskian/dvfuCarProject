<template>
  <nav class="navbar">
    <router-link 
      v-for="item in menuItems" 
      :key="item.to"
      :to="item.to" 
      class="nav-link"
      :class="{ active: $route.path === item.to }"
    >
      <span class="nav-icon">{{ item.icon }}</span>
      <span class="nav-text">{{ item.text }}</span>
    </router-link>
  </nav>
</template>

<script>
export default {
  name: 'Navbar',
  data() {
    return {
      menuItems: [
        { to: '/', icon: '🔍', text: 'Найти' },
        { to: '/create', icon: '🚗', text: 'Создать' },
        { to: '/my-trips', icon: '📋', text: 'Мои' },
      ]
    }
  }
}
</script>

<style scoped>
.navbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-around;
  align-items: center;
  height: 60px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-top: 1px solid #e0e0e0;
  z-index: 1000;
}

.nav-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #666;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.2s;
}

.nav-link.active {
  color: #007bff;
  background: rgba(0, 123, 255, 0.1);
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 2px;
}

.nav-text {
  font-size: 12px;
  font-weight: 500;
}
</style>