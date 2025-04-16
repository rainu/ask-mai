import { createRouter, createWebHashHistory } from 'vue-router'

import Chat from '../views/Chat.vue'
import History from '../views/History.vue'

const router = createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: '/',
			name: 'Chat',
			component: Chat,
		},
		{
			path: '/history',
			name: 'History',
			component: History,
		},
	],
})

export default router
