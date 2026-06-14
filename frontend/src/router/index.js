import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/Homeview.vue'
import MovieDetailView from '../views/MovieDetailView.vue'
import SeatView from '../views/SeatView.vue'
import BookingConfirmView from '../views/BookingConfirmView.vue'
import AdminHomeView from '../views/AdminHomeView.vue'
import HistoryView from '../views/HistoryView.vue'
import AdminBookingsView from '../views/AdminBookingsView.vue'
import AdminUsersView from '../views/AdminUsersView.vue'
import AdminLogsView from '../views/AdminLogsView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/login', component: LoginView },
    { path: '/home', component: HomeView, meta: { requiresAuth: true } },
    { path: '/movies/:id', component: MovieDetailView, meta: { requiresAuth: true } },
    { path: '/showtimes/:id/seats', component: SeatView, meta: { requiresAuth: true } },
    { path: '/bookings/:id', component: BookingConfirmView, meta: { requiresAuth: true } },
    { path: '/admin/home', component: AdminHomeView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/movies', component: AdminHomeView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/history', component: HistoryView, meta: { requiresAuth: true } },
    { path: '/admin/bookings', component: AdminBookingsView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/users', component: AdminUsersView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/logs', component: AdminLogsView, meta: { requiresAuth: true, requiresAdmin: true } },
{ path: '/profile', component: ProfileView, meta: { requiresAuth: true } },

  ],
})

router.beforeEach((to, from) => {
  const token = localStorage.getItem('token')
  const role = localStorage.getItem('role')
  if (to.meta.requiresAuth && !token) return '/login'
  if (to.meta.requiresAdmin && role !== 'admin') return '/home'
})

export default router