<template>
  <div class="my-trips">
    <h2>🚗 Мои поездки</h2>

    <h3>Как водитель</h3>
    <TripCard
      v-for="trip in driverTrips"
      :key="trip.id"
      :trip="trip"
      @click="goToTrip(trip.id)"
      :showActions="false"
    />
    <div v-if="!driverTrips.length" class="empty">
      Вы ещё не создали ни одной поездки
    </div>

    <h3>Как пассажир (бронирования)</h3>
    <div v-for="booking in passengerBookings" :key="booking.id" class="booking-item">
      <TripCard
        :trip="booking.trip"
        @click="goToTrip(booking.trip.id)"
        :showActions="false"
      />
      <div class="booking-status">Статус: {{ booking.status }}</div>
    </div>
    <div v-if="!passengerBookings.length" class="empty">
      У вас нет бронирований
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import TripCard from '../components/TripCardView.vue';
import { usersAPI } from '../services/api';

export default {
  name: 'MyTripsView',
  components: { TripCard },
  setup() {
    const router = useRouter();
    const driverTrips = ref([]);
    const passengerBookings = ref([]);

    const loadMyTrips = async () => {
      const user = window.Telegram?.WebApp?.initDataUnsafe?.user || { id: 1 };
      try {
        const response = await usersAPI.getUserTrips(user.id);
        driverTrips.value = response.data.as_driver;
        passengerBookings.value = response.data.as_passenger;
      } catch (error) {
        console.error('Ошибка загрузки поездок:', error);
      }
    };

    const goToTrip = (id) => {
      router.push(`/trip/${id}`);
    };

    onMounted(() => {
      loadMyTrips();
    });

    return { driverTrips, passengerBookings, goToTrip };
  },
};
</script>

<style scoped>
.my-trips {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
  padding-bottom: 80px;
}
h2, h3 {
  color: #333;
  margin: 20px 0 10px;
}
.empty {
  text-align: center;
  padding: 40px;
  color: #666;
}
.booking-status {
  margin: -10px 0 20px 0;
  padding: 8px;
  background: #f0f0f0;
  border-radius: 8px;
  text-align: center;
  font-size: 14px;
}
</style>