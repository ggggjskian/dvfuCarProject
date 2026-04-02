<template>
  <div class="home-view">
    <div class="header">
      <h1>🚗 DVFU Ride</h1>
      <p class="subtitle">Найди попутчика в кампус или из кампуса</p>
    </div>
    
    <div class="main-buttons">
      <router-link to="/find" class="main-button find-button">
        <span class="button-icon">🔍</span>
        <span class="button-text">Найти поездку</span>
        <span class="button-subtext">Пассажирам</span>
      </router-link>
      
      <router-link to="/create" class="main-button create-button">
        <span class="button-icon">🚗</span>
        <span class="button-text">Создать поездку</span>
        <span class="button-subtext">Водителям</span>
      </router-link>
    </div>
    
    <div v-if="recentTrips.length > 0" class="recent-trips">
      <h2>🚀 Недавние поездки</h2>
      <TripCard 
        v-for="trip in recentTrips" 
        :key="trip.id" 
        :trip="trip"
        @click="goToTripDetails(trip.id)"
      />
    </div>
    
    <div v-else class="empty-state">
      <div class="empty-icon">🚗</div>
      <h3>Пока нет активных поездок</h3>
      <p>Будь первым, кто создаст поездку!</p>
      <router-link to="/create" class="btn-primary">
        Создать первую поездку
      </router-link>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TripCard from '../components/TripCardView.vue'
import { tripsAPI } from '../services/api';

export default {
  name: 'HomeView',
  components: { TripCard },
  setup() {
    const router = useRouter()
    const recentTrips = ref([])
    const isLoading = ref(false)

    const loadRecentTrips = async () => {
      try {
        const response = await tripsAPI.getAll({ limit: 3 });
        recentTrips.value = response.data;
      } catch (error) {
        console.error('Ошибка загрузки поездок:', error);
      }
    };
    onMounted(() => {
      loadRecentTrips();
    });

    const goToTripDetails = (tripId) => {
      router.push(`/trip/${tripId}`)
    }

    onMounted(() => {
      loadRecentTrips()
    })

    return {
      recentTrips,
      isLoading,
      goToTripDetails
    }
  }
}
</script>

<style scoped>
.home-view {
  padding: 20px;
  max-width: 500px;
  margin: 0 auto;
  padding-bottom: 80px;
}

.header {
  text-align: center;
  margin-bottom: 30px;
}

.header h1 {
  font-size: 32px;
  margin-bottom: 8px;
  color: #333;
}

.subtitle {
  color: #666;
  font-size: 16px;
}

.main-buttons {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 30px;
}

.main-button {
  display: block;
  text-decoration: none;
  padding: 24px 20px;
  border-radius: 16px;
  transition: transform 0.2s, box-shadow 0.2s;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  cursor: pointer;
}

.main-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}

.main-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.find-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.create-button {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.button-icon {
  display: block;
  font-size: 40px;
  margin-bottom: 10px;
}

.button-text {
  display: block;
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 4px;
}

.button-subtext {
  display: block;
  font-size: 14px;
  opacity: 0.9;
}

.recent-trips h2 {
  font-size: 20px;
  margin-bottom: 16px;
  color: #333;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 60px;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 20px;
  margin-bottom: 8px;
  color: #333;
}

.empty-state p {
  color: #666;
  margin-bottom: 20px;
}

.btn-primary {
  display: inline-block;
  background: #007bff;
  color: white;
  padding: 12px 24px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  border: none;
  font-size: 16px;
}
</style>