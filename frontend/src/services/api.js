import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8000'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true
})

export function getCurrentUser() {
  if (typeof window !== 'undefined' && window.Telegram?.WebApp?.initDataUnsafe?.user) {
    const tgUser = window.Telegram.WebApp.initDataUnsafe.user;
    return {
      id: tgUser.id,
      is_telegram: true,
      name: tgUser.first_name || 'Пользователь Telegram'
    };
  }
  
  const savedUser = localStorage.getItem('user');
  if (savedUser) {
    const user = JSON.parse(savedUser);
    return {
      id: user.id, // Внутренний ID из БД
      is_telegram: false,
      name: user.username || 'Пользователь'
    };
  }
  
  return null;
}

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  register: (data) => api.post('/api/register', data),
  login: (data) => api.post('/api/login', data).then(res => {
    localStorage.setItem('token', res.data.token);
    localStorage.setItem('user', JSON.stringify(res.data.user));
    return res;
  }),
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }
}

export const usersAPI = {
  getOrCreate: (userData) => api.post('/api/users', userData),
  getUserTrips: (tgId) => api.get(`/api/users/${tgId}/trips`)
}

export const tripsAPI = {
  getAll: (params) => api.get('/api/trips', { params }),
  search: (params) => api.get('/api/search_trips', { params }),
  getById: (id) => api.get(`/api/trips/${id}`),
  create: (tripData) => api.post('/api/trips', tripData),
  book: (tripId, bookingData) => api.post(`/api/trips/${tripId}/book`, bookingData),
}

export const bookingsAPI = {
  updateStatus: (bookingId, status) => 
    api.patch(`/api/bookings/${bookingId}`, { status })
}

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
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        action: 'location',
        trip_id: tripId,
        driver_id: driverId,
        lat,
        lon
      }))
    }
  }

  return { ws, sendLocation }
}

export default api