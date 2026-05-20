<template>
  <div class="trip-card" @click="$emit('click', trip)">
    <div class="trip-header">
      <div class="driver-info">
        <div class="driver-avatar">
          {{ trip.driver?.first_name?.charAt(0) || '?' }}
        </div>
        <div class="driver-details">
          <span class="driver-name">{{ trip.driver?.first_name || 'Водитель' }}</span>
          <span class="driver-rating">⭐ {{ trip.driver?.rating?.toFixed(1) || '5.0' }}</span>
        </div>
      </div>
      <div class="trip-price">
        {{ trip.price }} ₽
      </div>
    </div>
    
    <div class="trip-route">
      <div class="route-point">
        <span class="point-icon">{{ trip.trip_type === 'from_campus' ? '🏫' : '📍' }}</span>
        <span class="point-text">{{ trip.trip_type === 'from_campus' ? 'Кампус ДВФУ' : trip.point }}</span>
      </div>
      <div class="route-divider">↓</div>
      <div class="route-point">
        <span class="point-icon">{{ trip.trip_type === 'from_campus' ? '📍' : '🏫' }}</span>
        <span class="point-text">{{ trip.trip_type === 'from_campus' ? trip.point : 'Кампус ДВФУ' }}</span>
      </div>
    </div>
    
    <div class="trip-footer">
      <div class="trip-time">
        <span class="time-icon">🕒</span>
        <span class="time-text">{{ formatTime(trip.departure_time) }}</span>
      </div>
      <div class="trip-seats">
        <span class="seats-icon">💺</span>
        <span class="seats-text">{{ trip.seats_available }} из {{ trip.seats_total }}</span>
      </div>
      <div v-if="trip.distance_km" class="trip-distance">
        <span class="distance-icon">📏</span>
        <span class="distance-text">{{ trip.distance_km }} км</span>
      </div>
    </div>
    
    <div v-if="showActions" class="trip-actions">
      <button 
        class="btn-book"
        @click.stop="bookTrip"
        :disabled="trip.seats_available === 0"
      >
        {{ trip.seats_available === 0 ? 'Мест нет' : 'Забронировать' }}
      </button>
    </div>
  </div>
</template>

<script>
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'

export default {
  name: 'TripCard',
  props: {
    trip: {
      type: Object,
      required: true
    },
    showActions: {
      type: Boolean,
      default: true
    }
  },
  emits: ['click', 'book'],
  methods: {
    formatTime(dateString) {
      try {
        const date = new Date(dateString)
        return format(date, 'HH:mm, dd MMM', { locale: ru })
      } catch {
        return dateString
      }
    },
    bookTrip() {
      this.$emit('book', this.trip)
    }
  }
}
</script>

<style scoped>
.trip-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.trip-card:active {
  transform: scale(0.98);
}

.trip-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.driver-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.driver-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 18px;
}

.driver-details {
  display: flex;
  flex-direction: column;
}

.driver-name {
  font-weight: 600;
  font-size: 16px;
}

.driver-rating {
  font-size: 12px;
  color: #666;
}

.trip-price {
  font-size: 20px;
  font-weight: bold;
  color: #007bff;
}

.trip-route {
  margin-bottom: 16px;
}

.route-point {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 4px 0;
}

.point-icon {
  font-size: 20px;
}

.point-text {
  flex: 1;
  font-size: 14px;
}

.route-divider {
  text-align: center;
  color: #999;
  margin: 2px 0;
  font-size: 14px;
}

.trip-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #666;
  margin-bottom: 12px;
}

.trip-time, .trip-seats, .trip-distance {
  display: flex;
  align-items: center;
  gap: 4px;
}

.trip-actions {
  margin-top: 12px;
}

.btn-book {
  width: 100%;
  padding: 12px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-book:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-book:active:not(:disabled) {
  background: #0056b3;
}
</style>