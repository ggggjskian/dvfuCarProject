import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import FindView from '../views/FindView.vue'
import CreateView from '../views/CreateView.vue'
import TripDetailsView from '../views/TripDetailsView.vue'
import ProfileView from '../views/profileView.vue'
import MyTripsView from '../views/MyTripsView.vue'

// Создаем простые компоненты прямо здесь
const FindTripView = { 
  template: '<div><h1>🔍 Поиск поездок</h1><p>Раздел в разработке</p></div>',
  name: 'FindTripView'
}

const CreateTripView = {
  template: '<div><h1>🚗 Создать поездку</h1><p>Раздел в разработке</p></div>',
  name: 'CreateTripView'
}



const routes = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/find', name: 'Find', component: FindView },
  { path: '/create', name: 'Create', component: CreateView },
  { path: '/trip/:id', name: 'TripDetails', component: TripDetailsView, props: true },
  { path: '/profile', name: 'profile', component: ProfileView },
  { path: '/my-trips', name: 'MyTrips', component: MyTripsView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router