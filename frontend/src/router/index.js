import { createRouter, createWebHistory } from 'vue-router'
import { getCurrentUser } from '@/services/api' 

import HomeView from '../views/HomeView.vue'
import FindView from '../views/FindView.vue'
import CreateView from '../views/CreateView.vue'
import TripDetailsView from '../views/TripDetailsView.vue'
import ProfileView from '../views/profileView.vue'
import MyTripsView from '../views/MyTripsView.vue'
import LoginForm from '@/views/LoginForm.vue'
import RegForm from '@/views/RegForm.vue'

const routes = [
  { path: '/', redirect: '/create' },
  
  { path: '/login', name: 'login', component: LoginForm, meta: { guestOnly: true } },
  { path: '/reg', name: 'reg', component: RegForm, meta: { guestOnly: true } },
  
  { path: '/home', name: 'Home', component: HomeView, meta: { requiresAuth: true } },
  { path: '/find', name: 'Find', component: FindView, meta: { requiresAuth: true } },
  { path: '/create', name: 'Create', component: CreateView, meta: { requiresAuth: true } },
  { path: '/trip/:id', name: 'TripDetails', component: TripDetailsView, props: true, meta: { requiresAuth: true } },
  { path: '/profile', name: 'profile', component: ProfileView, meta: { requiresAuth: true } },
  { path: '/my-trips', name: 'MyTrips', component: MyTripsView, meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const user = getCurrentUser()
  const isAuthenticated = !!user

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/reg')
  } 
  else if (to.meta.guestOnly && isAuthenticated) {
    next('/create')
  } 
  else {
    next()
  }
})

export default router