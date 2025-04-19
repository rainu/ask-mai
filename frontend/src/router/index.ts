import { createRouter, createWebHashHistory } from 'vue-router'

import Chat from '../views/Chat.vue'
import History from '../views/History.vue'
import Edit from '../views/Edit.vue'
import Profile from '../views/Profile.vue'
import Tool from '../views/Tool.vue'

const router = createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: '/',
			name: 'Chat',
			component: Chat,
		},
		{
			path: '/edit/:idx',
			name: 'Edit',
			component: Edit,
		},
		{
			path: '/history',
			name: 'History',
			component: History,
		},
		{
			path: '/profile',
			name: 'Profile',
			component: Profile,
		},
		{
			path: '/tool',
			name: 'Tool',
			component: Tool,
		},
	],
})

export default router
