import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: '/',
			name: 'home',
			component: () => import('../views/Home.vue'),
		},
		{
			path: '/history',
			name: 'history',
			component: () => import('../views/History.vue'),
		},
	],
})

export default router
