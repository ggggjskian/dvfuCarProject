import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8000'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true
})

// ----- Пользователь Telegram -----
function getTelegramUser() {
  // Если мы не в Telegram (например, localhost), возвращаем тестового пользователя
  if (!window.Telegram?.WebApp) {
    return { id: 1 }; // тестовый ID
  }
  return window.Telegram.WebApp.initDataUnsafe?.user || null
}

// ----- Интерцептор для добавления tg_id -----
api.interceptors.request.use((config) => {
  config.params = config.params || {};
  // Добавляем параметры только для POST/PUT/PATCH и для GET списка, но не для GET конкретной поездки
  if (!config.url.match(/\/api\/trips\/\d+$/)) {  // если это не GET /api/trips/число
    config.params.driver_tg_id = 1;
    config.params.passenger_tg_id = 1;
  }
  return config;
});

// ----- API для поездок -----
export const tripsAPI = {
  getAll: (params) => api.get('/api/trips', { params }),
  search: (params) => api.get('/api/search_trips', { params }),
  getById: (id) => api.get(`/api/trips/${id}`),
  create: (tripData) => api.post('/api/trips', tripData),
  book: (tripId, bookingData) => {
    console.log("Booking request data:", bookingData);
    return api.post(`/api/trips/${tripId}/book`, bookingData);
  },
 } 

// ----- API для бронирований -----
export const bookingsAPI = {
  updateStatus: (bookingId, status) => 
    api.patch(`/api/bookings/${bookingId}`, { status })
}

// ----- API для пользователей -----
export const usersAPI = {
  getOrCreate: (userData) => api.post('/api/users', userData),
  getUserTrips: (tgId) => api.get(`/api/users/${tgId}/trips`)
}

// ----- WebSocket трекер -----
export const createDriverTracker = (tripId, onLocationUpdate) => {
  const wsUrl = API_BASE_URL.replace('http', 'ws') + `/ws/trip/${tripId}`
  const ws = new WebSocket(wsUrl)

  ws.onopen = () => console.log('WebSocket connected')
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    if (data.type === 'driver_location') {
      onLocationUpdate(data)
    }
  }
  ws.onerror = (error) => console.error('WebSocket error', error)

  const sendLocation = (lat, lon, driverId) => {
    ws.send(JSON.stringify({
      action: 'location',
      trip_id: tripId,
      driver_id: driverId,
      lat,
      lon
    }))
  }

  return { ws, sendLocation }
}

// ----- Экспорт по умолчанию -----
export default api