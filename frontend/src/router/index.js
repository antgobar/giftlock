import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../components/Dashboard.vue'
import Home from '../components/Home.vue'
import Login from '../components/Login.vue'
import Register from '../components/Register.vue'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: { requiresGuest: true }
    },
    {
        path: '/register',
        name: 'Register',
        component: Register,
        meta: { requiresGuest: true }
    },
    {
        path: '/dashboard',
        name: 'Dashboard',
        component: Dashboard,
        meta: { requiresAuth: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Route guard to check authentication
router.beforeEach(async (to, from, next) => {
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
    const requiresGuest = to.matched.some(record => record.meta.requiresGuest)

    try {
        // Check authentication status
        const response = await fetch('/api/me', {
            method: 'GET',
            credentials: 'include'
        })

        const isAuthenticated = response.ok

        if (requiresAuth && !isAuthenticated) {
            // Redirect to login if authentication required but user not authenticated
            next('/login')
        } else if (requiresGuest && isAuthenticated) {
            // Redirect to dashboard if guest page but user is authenticated
            next('/dashboard')
        } else {
            next()
        }
    } catch (error) {
        console.error('Authentication check failed:', error)
        if (requiresAuth) {
            next('/login')
        } else {
            next()
        }
    }
})

export default router