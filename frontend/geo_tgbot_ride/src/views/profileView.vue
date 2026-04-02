<template>
    <div class="profile-container">
        <h2>Профиль</h2>
        <div class="profile-info">
            <p><strong>ID:</strong> {{ user.id }}</p>
            <p><strong>Имя:</strong> {{ user.name }}</p>
            <p><strong>Выполнено поездок:</strong> {{ user.completedTrips }}</p>
        </div>
    </div>
</template>

<script>
export default {
    name: 'profileView',
    data() {
        return {
            user: {
                id: '',
                name: '',
                completedTrips: 0
            }
        }
    },
    async mounted() {
        await this.fetchProfile()
    },
    methods: {
        async fetchProfile() {
            try {
                const tg_user = window.Telegram.WebApp.initDataUnsafe?.user
                const response = await fetch(`${import.meta.env.VITE_API_URL}/${t.id}`)
                const data = await response.json()
                this.user.id = tg_user.id
                this.user.name = tg_user.first_name
                this.user.completedTrips = data.completedTrips
            } catch (error) {
                console.log(error)
            }
        }
    }
}
</script>

<style scoped>
.profile-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  padding: 16px;
}

.profile-info {
  background-color: #ffffffcc;
  backdrop-filter: blur(8px);
  padding: 16px;
  border-radius: 8px;
  text-align: left;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  margin-top: 16px;
  width: 100%;
  max-width: 320px;
}
</style>