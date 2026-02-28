import { createRouter, createWebHistory } from 'vue-router'
import auth from '@/auth'

const routes = [
	{
		path: '/',
		name: 'home',
		component: () => import(/* webpackChunkName: "home" */ '@/views/Home.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/login',
		name: 'login',
		component: () => import(/* webpackChunkName: "auth" */ '@/views/Login.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/cadastro',
		name: 'cadastro',
		component: () => import(/* webpackChunkName: "auth" */ '@/views/Cadastro.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/esqueci-senha',
		name: 'esqueci-senha',
		component: () => import(/* webpackChunkName: "auth" */ '@/views/EsqueciSenha.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/sobre-nos',
		name: 'sobre-nos',
		component: () => import(/* webpackChunkName: "about" */ '@/views/SobreNos.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/explorar',
		name: 'explorar',
		component: () => import(/* webpackChunkName: "explore" */ '@/views/Explore.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/dashboard',
		name: 'dashboard',
		component: () => import(/* webpackChunkName: "secure" */ '@/views/Dashboard.vue'),
		beforeEnter: (to, from, next) => {
			auth.isAuthenticated() ? next() : next('/login')
		}
	},
	{
		path: '/perfil',
		name: 'perfil',
		component: () => import(/* webpackChunkName: "secure" */ '@/views/ProfileEdit.vue'),
		beforeEnter: (to, from, next) => {
			auth.isAuthenticated() ? next() : next('/login')
		}
	},
	{
		path: '/estudo/:id',
		name: 'estudo',
		component: () => import(/* webpackChunkName: "study" */ '@/views/Estudo.vue'),
		meta: { requiresAuth: false }
	},
	{
		path: '/flashcards',
		name: 'flashcards',
		component: () => import(/* webpackChunkName: "secure" */ '@/views/Flashcards.vue'),
		beforeEnter: (to, from, next) => {
			auth.isAuthenticated() ? next() : next('/login')
		}
	},
	{
		path: '/admin',
		name: 'admin',
		component: () => import(/* webpackChunkName: "secure" */ '@/views/AdminDashboard.vue'),
		beforeEnter: (to, from, next) => {
			const user = auth.getUser();
			if (auth.isAuthenticated() && (user.tipo === 1 || user.tipo === 3)) { // 1: Curador, 3: Admin (se houver)
				next();
			} else {
				next('/dashboard');
			}
		}
	}
]

const router = createRouter({
	history: createWebHistory(),
	routes: routes,
})

export default router
