import { createRouter, createWebHashHistory } from 'vue-router'

import Home from '../views/Home.vue'
import History from '../views/History.vue'

const router = createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: '/',
			name: 'Home',
			component: Home,
		},
		{
			path: '/history',
			name: 'History',
			component: History,
		},
	],
})

export default router
