import { reactive } from 'vue'

const state = reactive({
	isAuthenticated: localStorage.getItem('auth') === 'true',
	user: JSON.parse(localStorage.getItem('user') || '{}')
})

export default {
	login(user) {
		localStorage.setItem('auth', 'true')
		localStorage.setItem('user', JSON.stringify(user))
		state.isAuthenticated = true
		state.user = user
	},
	logout() {
		localStorage.removeItem('auth')
		localStorage.removeItem('user')
		state.isAuthenticated = false
		state.user = {}
	},
	isAuthenticated() {
		return state.isAuthenticated
	},
	getUser() {
		return state.user
	}
}
