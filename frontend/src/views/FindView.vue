<template>
  <div class="find-view">
    <h2>🔍 Поиск поездок</h2>
    
    <div class="search-form">
      <div class="trip-type">
        <label>
          <input type="radio" value="from_campus" v-model="tripType" />
          🏫 Из кампуса
        </label>
        <label>
          <input type="radio" value="to_campus" v-model="tripType" />
          📍 В кампус
        </label>
      </div>
      
      <div class="address-input">
        <input 
          v-model="address" 
          placeholder="Введите адрес или кликните на карте" 
          @keyup.enter="searchTrips"
        />
        <button @click="searchTrips" :disabled="loading">
          {{ loading ? 'Поиск...' : 'Найти' }}
        </button>
      </div>
      
      <div class="filters">
        <label>
          Радиус (км):
          <input type="number" v-model="maxDistance" min="1" max="20" />
        </label>
        <label>
          Отклонение (мин):
          <input type="number" v-model="maxDeviation" min="0" max="60" />
        </label>
      </div>
    </div>
    
    <!-- Карта с точками -->
    <div id="map" class="map-container"></div>
    
    <!-- Результаты поиска -->
    <div class="results">
      <h3>Найдено поездок: {{ trips.length }}</h3>
      <div v-if="loading" class="loading">Загрузка...</div>
      <div v-else-if="trips.length === 0" class="empty">
        Нет подходящих поездок
      </div>
      <TripCard
        v-for="trip in trips"
        :key="trip.id"
        :trip="trip"
        @click="goToTrip(trip.id)"
        showActions
        @book="bookTrip(trip)"
      />
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TripCard from '../components/TripCardView.vue'
import { tripsAPI } from '../services/api.js'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

export default {
  name: 'FindView',
  components: { TripCard },
  setup() {
    const router = useRouter()
    const tripType = ref('from_campus')
    const address = ref('')
    const maxDistance = ref(5)
    const maxDeviation = ref(30)
    const loading = ref(false)
    const trips = ref([])
    let map = null
    let marker = null

    const initMap = () => {
      map = L.map('map').setView([43.0245, 131.8927], 12)
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap'
      }).addTo(map)

      map.on('click', async (e) => {
        const { lat, lng } = e.latlng
        if (marker) marker.remove()
        marker = L.marker([lat, lng]).addTo(map)
        address.value = `${lat.toFixed(5)}, ${lng.toFixed(5)}`
      })
    }

    const searchTrips = async () => {
      const trimmedAddress = address.value.trim()
      if (!trimmedAddress) {
        alert("Введите адрес или выберите точку на карте")
        return
      }
      loading.value = true
      try {
        const params = {
          trip_type: tripType.value,
          address: trimmedAddress,
          max_distance_km: maxDistance.value,
          max_deviation_minutes: maxDeviation.value
        }
        console.log("Search params:", params)
        const result = await tripsAPI.search(params)
        trips.value = result
        // Центрируем карту на первом результате, если есть
        if (result.length > 0 && result[0].point_lat && result[0].point_lon) {
          map.setView([result[0].point_lat, result[0].point_lon], 13)
        }
      } catch (error) {
        console.error("Ошибка поиска:", error)
        alert("Не удалось выполнить поиск")
      } finally {
        loading.value = false
      }
    }

    const goToTrip = (id) => {
      router.push(`/trip/${id}`)
    }

    const bookTrip = async (trip) => {
      try {
        await tripsAPI.book(trip.id, {})
        alert('✅ Заявка отправлена водителю')
      } catch (error) {
        console.error("Ошибка бронирования:", error)
        alert('❌ Не удалось забронировать')
      }
    }

    onMounted(() => {
      initMap()
    })

    return {
      tripType,
      address,
      maxDistance,
      maxDeviation,
      loading,
      trips,
      searchTrips,
      goToTrip,
      bookTrip
    }
  }
}
</script>

<style scoped>
.map-container {
  height: 300px;
  margin: 20px 0;
  border-radius: 12px;
  overflow: hidden;
}
.search-form {
  background: white;
  padding: 16px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.trip-type {
  display: flex;
  gap: 20px;
  margin-bottom: 12px;
}
.address-input {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}
.address-input input {
  flex: 1;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
}
.address-input button {
  padding: 12px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
}
.filters {
  display: flex;
  gap: 20px;
}
.filters label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}
.filters input {
  width: 60px;
  padding: 6px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>