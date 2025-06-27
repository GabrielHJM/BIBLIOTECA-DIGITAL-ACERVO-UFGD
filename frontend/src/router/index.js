import { createRouter, createWebHistory } from 'vue-router'
import auth from '@/auth'

const routes = [
	{
		path: '/',
		name: 'home',
		component: () => import('@/views/Home.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/login',
		name: 'login',
		component: () => import('@/views/Login.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/cadastro',
		name: 'cadastro',
		component: () => import('@/views/Cadastro.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/dashboard',
		name: 'dashboard',
		component: () => import('@/views/Dashboard.vue'),
		beforeEnter: (to, from, next) => {
			auth.isAuthenticated() ? next() : next('/login')
		}
	},
	{
		path: '/estudo',
		name: 'estudo',
		component: () => import('@/views/Estudo.vue'),
		meta: { requiresAuth: false }
	}
]

const router = createRouter({
	history: createWebHistory(),
	routes: routes,
})

export default router
