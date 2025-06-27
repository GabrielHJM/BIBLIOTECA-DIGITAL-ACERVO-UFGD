// src/auth.js
export default {
	login(user) {
		localStorage.setItem('auth', 'true')
		localStorage.setItem('user', JSON.stringify(user))
	},
	logout() {
		localStorage.removeItem('auth')
		localStorage.removeItem('user')
	},
	isAuthenticated() {
		return localStorage.getItem('auth') === 'true'
	},
	getUser() {
		return JSON.parse(localStorage.getItem('user') || '{}')
	}
}
