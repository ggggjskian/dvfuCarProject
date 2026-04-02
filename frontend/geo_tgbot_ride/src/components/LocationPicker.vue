<template>
  <div class="location-picker">
    <div id="picker-map" class="map"></div>
    <div class="address-bar" v-if="selectedAddress">
      <span>📍 {{ selectedAddress }}</span>
      <button @click="confirmLocation" class="confirm-btn">Подтвердить</button>
    </div>
    <div v-else class="hint">
      Нажмите на карту, чтобы выбрать точку назначения
    </div>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

export default {
  name: 'LocationPicker',
  emits: ['select'],
  setup(props, { emit }) {
    const selectedAddress = ref('')
    const selectedLat = ref(null)
    const selectedLon = ref(null)
    let map = null
    let marker = null

    // Ваш API-ключ Яндекс.Карт (получите бесплатно в Яндекс.Облаке)
    const YANDEX_API_KEY = 'a1e4868e-80d2-4c26-a2d6-460f06388489'  // ⚠️ В продакшене вынесите в .env

    onMounted(() => {
      map = L.map('picker-map').setView([43.0245, 131.8927], 13)
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap'
      }).addTo(map)

      map.on('click', async (e) => {
        const { lat, lng } = e.latlng
        if (marker) marker.remove()
        marker = L.marker([lat, lng]).addTo(map)

        try {
          // Яндекс.Обратное геокодирование
          const url = `https://geocode-maps.yandex.ru/1.x/?apikey=${YANDEX_API_KEY}&geocode=${lng},${lat}&format=json&results=1`
          const res = await fetch(url)
          const data = await res.json()
          
          const feature = data.response.GeoObjectCollection.featureMember[0]?.GeoObject
          if (feature) {
            selectedAddress.value = feature.metaDataProperty.GeocoderMetaData.text
          } else {
            selectedAddress.value = `Координаты: ${lat.toFixed(4)}, ${lng.toFixed(4)}`
          }
          selectedLat.value = lat
          selectedLon.value = lng
        } catch (error) {
          console.error('Ошибка геокодирования', error)
          selectedAddress.value = `Координаты: ${lat.toFixed(4)}, ${lng.toFixed(4)}`
          selectedLat.value = lat
          selectedLon.value = lng
        }
      })
    })

    const confirmLocation = () => {
      emit('select', {
        address: selectedAddress.value,
        lat: selectedLat.value,
        lon: selectedLon.value
      })
    }

    return { selectedAddress, confirmLocation }
  }
}
</script>

<style scoped>
.location-picker {
  width: 100%;
  height: 400px;
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 20px;
}
.map {
  width: 100%;
  height: 100%;
}
.address-bar {
  position: absolute;
  bottom: 20px;
  left: 20px;
  right: 20px;
  background: white;
  padding: 12px 20px;
  border-radius: 30px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.hint {
  position: absolute;
  bottom: 20px;
  left: 20px;
  right: 20px;
  background: rgba(0,0,0,0.7);
  color: white;
  padding: 12px;
  border-radius: 30px;
  text-align: center;
}
.confirm-btn {
  background: #007bff;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 20px;
  font-weight: 600;
  cursor: pointer;
}
</style>