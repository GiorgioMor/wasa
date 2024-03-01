import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import usersList from '../views/usersList.vue'
import newPost from '../views/newPost.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/usersList', component: usersList},
		{path: '/newPost', component: newPost},
		{path: '/home', component: HomeView},
		{path: '/profile/:id', component: ProfileView},
		{path: '/some/:id/link', component: HomeView},
		
	]
})

export default router
