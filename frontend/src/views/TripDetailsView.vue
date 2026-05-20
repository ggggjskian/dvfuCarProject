<template>
  <div class="trip-details">
    <div v-if="loading">Загрузка...</div>
    <div v-else-if="!trip">Поездка не найдена</div>
    <div v-else>
      <div class="header">
        <h2>🚗 Поездка #{{ trip.id }}</h2>
        <span class="status" :class="trip.status">{{ trip.status }}</span>
      </div>
      
      <TripCard :trip="trip" :showActions="false" />
      
      <div class="actions">
        <button 
          v-if="canBook && trip.seats_available > 0"
          @click="bookTrip" 
          class="btn-book"
          :disabled="bookingInProgress"
        >
          {{ bookingInProgress ? 'Бронирование...' : 'Забронировать место' }}
        </button>
        <button 
          v-if="isDriver && trip.status === 'active'"
          @click="startTrip" 
          class="btn-start"
        >
          🚦 Начать поездку
        </button>
        <button 
          v-if="isDriver && trip.status === 'in_progress'"
          @click="completeTrip" 
          class="btn-complete"
        >
          ✅ Завершить
        </button>
      </div>
      
      <!-- Блок отслеживания, если поездка началась -->
      <div v-if="trip.status === 'in_progress' && driverLocation" class="tracking-section">
        <h3>📍 Водитель в пути</h3>
        <DriverTracker 
          :tripId="trip.id" 
          :userLat="userLat" 
          :userLon="userLon"
          @locationUpdate="onDriverLocationUpdate"
        />
        <div v-if="isPassenger && arrivalTime" class="eta">
          ⏱ Примерное прибытие: {{ arrivalTime }}
        </div>
      </div>
      
      <!-- Список пассажиров (для водителя) -->
      <div v-if="isDriver && trip.bookings?.length" class="passengers-list">
        <h3>👥 Пассажиры</h3>
        <div v-for="booking in trip.bookings" :key="booking.id" class="passenger-item">
          <span>{{ booking.passenger.first_name }}</span>
          <span :class="'status ' + booking.status">{{ booking.status }}</span>
          <div class="passenger-actions" v-if="booking.status === 'pending'">
            <button @click="updateBooking(booking.id, 'confirmed')" class="btn-confirm">✅</button>
            <button @click="updateBooking(booking.id, 'rejected')" class="btn-reject">❌</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TripCard from '../components/TripCardView.vue'
import DriverTracker from '../components/DriveTracker.vue'
import { tripsAPI, bookingsAPI, createDriverTracker } from '../services/api'

export default {
  name: 'TripDetailsView',
  components: { TripCard, DriverTracker },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const trip = ref(null)
    const loading = ref(true)
    const bookingInProgress = ref(false)
    const driverLocation = ref(null)
    const arrivalTime = ref(null)
    
    // Текущий пользователь Telegram
    const tgUser = window.Telegram?.WebApp?.initDataUnsafe?.user || null
    const isDriver = computed(() => trip.value?.driver?.tg_id === tgUser?.id)
    const isPassenger = computed(() => {
      return trip.value?.bookings?.some(b => b.passenger.tg_id === tgUser?.id)
    })
    const canBook = computed(() => {
      return !isDriver.value && !isPassenger.value && trip.value?.status === 'active'
    })
    
    // Координаты пользователя (для трекера) — можно запросить
    const userLat = ref(null)
    const userLon = ref(null)
    
    const loadTrip = async () => {
      try {
        const id = route.params.id;
        console.log("Loading trip with id:", id);
        const response = await tripsAPI.getById(id);
        console.log("Response:", response);
        trip.value = response.data;   // <-- важно: response.data
      } catch (err) {
        console.error("Error loading trip:", err);
      } finally {
        loading.value = false;
      }
    };
    
    const bookTrip = async () => {
      bookingInProgress.value = true
      try {
        await tripsAPI.book(trip.value.id, {})
        alert('Заявка отправлена! Ожидайте подтверждения водителя.')
        await loadTrip() // обновить данные
      } catch (error) {
        alert('Ошибка бронирования')
      } finally {
        bookingInProgress.value = false
      }
    }
    
    const startTrip = async () => {
      if (!confirm('Начать поездку? Пассажиры будут видеть ваше местоположение.')) return
      // Здесь должен быть API-вызов для изменения статуса поездки на 'in_progress'
      // Пока заглушка
      alert('Функция старта поездки в разработке (нужен отдельный эндпоинт)')
    }
    
    const completeTrip = async () => {
      // Аналогично
      alert('Завершение поездки')
    }
    
    const updateBooking = async (bookingId, status) => {
      try {
        await bookingsAPI.updateStatus(bookingId, status)
        await loadTrip()
      } catch (error) {
        alert('Ошибка обновления')
      }
    }
    
    const onDriverLocationUpdate = (location) => {
      driverLocation.value = location
      // Здесь можно рассчитать ETA, если есть координаты пассажира
    }
    
    onMounted(() => {
      loadTrip()
      // Попросить геолокацию, если нужно для трекера
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition((pos) => {
          userLat.value = pos.coords.latitude
          userLon.value = pos.coords.longitude
        })
      }
    })
    
    return {
      trip,
      loading,
      isDriver,
      isPassenger,
      canBook,
      bookingInProgress,
      userLat,
      userLon,
      driverLocation,
      arrivalTime,
      bookTrip,
      startTrip,
      completeTrip,
      updateBooking,
      onDriverLocationUpdate
    }
  }
}
</script>

<style scoped>
.trip-details {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
  padding-bottom: 80px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.status {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  text-transform: uppercase;
}

.status.active { background: #d4edda; color: #155724; }
.status.in_progress { background: #fff3cd; color: #856404; }
.status.completed { background: #e2e3e5; color: #383d41; }
.status.cancelled { background: #f8d7da; color: #721c24; }

.actions {
  margin: 24px 0;
}

.btn-book, .btn-start, .btn-complete {
  width: 100%;
  padding: 16px;
  border: none;
  border-radius: 12px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
}

.btn-book {
  background: #007bff;
  color: white;
}

.btn-start {
  background: #28a745;
  color: white;
}

.btn-complete {
  background: #6c757d;
  color: white;
}

.passengers-list {
  margin-top: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 16px;
}

.passenger-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #eee;
}

.passenger-item:last-child {
  border-bottom: none;
}

.passenger-actions {
  display: flex;
  gap: 8px;
}

.btn-confirm, .btn-reject {
  border: none;
  background: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
}

.btn-confirm:hover { background: #d4edda; }
.btn-reject:hover { background: #f8d7da; }

.tracking-section {
  margin-top: 30px;
}

.eta {
  margin-top: 12px;
  padding: 12px;
  background: #e8f4fd;
  border-radius: 8px;
  text-align: center;
}
</style>