import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Инициализация Telegram WebApp
function initTelegramWebApp() {
  if (window.Telegram?.WebApp) {
    const tg = window.Telegram.WebApp
    
    // Расширяем на весь экран
    tg.expand()
    
    // Настраиваем тему
    if (tg.colorScheme === 'dark') {
      document.documentElement.classList.add('dark-theme')
    }
    
    // Включаем кнопку "Назад"
    tg.BackButton.show()
    tg.BackButton.onClick(() => {
      if (window.history.length > 1) {
        router.go(-1)
      } else {
        tg.close()
      }
    })
    
    // Скрываем кнопку "Назад" на главной
    router.afterEach((to) => {
      if (to.path === '/') {
        tg.BackButton.hide()
      } else {
        tg.BackButton.show()
      }
    })
  }
}

// Создаём и монтируем приложение
const app = createApp(App)
app.use(router)

// Инициализируем Telegram перед монтированием
initTelegramWebApp()

app.mount('#app')
