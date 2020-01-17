import Vue from 'vue'
import VueRouter from 'vue-router'
import routes from './router/router'
import store from './store/'
import {routerMode} from './config/env'
import './config/rem'

Vue.use(VueRouter)

const router = new VueRouter({
	routes,
	mode: routerMode,
	strict: process.env.NODE_ENV !== 'production',
	scrollBehavior (to, from, savedPosition) {
	    if (savedPosition) {
		    return savedPosition
		} else {
			if (from.meta.keepAlive) {
				from.meta.savedPosition = document.body.scrollTop;
			}
		    return { x: 0, y: to.meta.savedPosition || 0 }
		}
	}
})

/*router.beforeResolve((to, from, next) => {
  alert("sssss2")
  next()
})

router.beforeEach((to, from, next) => {
  alert("sssss")
  next()
})*/


/*router.routes.beforeRouteLeave((to, from, next) => {
  alert("sssssssss")
  next()
})*/

new Vue({
	router,
	store,
}).$mount('#app')






