<template>
  <div class="tracker">
    <div id="tracker-map" class="map"></div>
    <div class="status">
      🚗 Водитель в пути
      <span v-if="distance"> ~ {{ distance }} км от вас</span>
    </div>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue'
import L from 'leaflet'
import { createDriverTracker } from '../services/api'

// 📍 Временное решение — функция прямо здесь
const haversine = (lat1, lon1, lat2, lon2) => {
  const R = 6371
  const dLat = (lat2 - lat1) * Math.PI / 180
  const dLon = (lon2 - lon1) * Math.PI / 180
  const a = 
    Math.sin(dLat/2) * Math.sin(dLat/2) +
    Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
    Math.sin(dLon/2) * Math.sin(dLon/2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
  return R * c
}

export default {
  props: {
    tripId: { type: Number, required: true },
    userLat: { type: Number, default: null },
    userLon: { type: Number, default: null }
  },
  setup(props) {
    let map = null
    let driverMarker = null
    const distance = ref(null)
    
    onMounted(() => {
      map = L.map('tracker-map').setView([43.0245, 131.8927], 13)
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png').addTo(map)
      
      const { ws, sendLocation } = createDriverTracker(props.tripId, (data) => {
        const { lat, lon } = data
        if (driverMarker) {
          driverMarker.setLatLng([lat, lon])
        } else {
          driverMarker = L.marker([lat, lon], {
            icon: L.icon({ iconUrl: '/car-icon.png', iconSize: [32, 32] })
          }).addTo(map)
        }
        map.setView([lat, lon], 13)
        
        if (props.userLat && props.userLon) {
          distance.value = haversine(lat, lon, props.userLat, props.userLon).toFixed(1)
        }
      })
    })
    
    return { distance }
  }
}
</script>

<style scoped>
.map {
  height: 300px;
  width: 100%;
  border-radius: 12px;
}
.status {
  padding: 12px;
  background: #f0f9ff;
  border-radius: 8px;
  margin-top: 12px;
}
</style>