<script>
import { RouterLink, RouterView } from 'vue-router'

export default {
	data: function() {
		return {
			id: -1,
			search: "",
		}
	},
	methods: {
		async searchUser() {
			this.$router.replace("/usersList");
		}
	},
	mounted() {
		if (sessionStorage.getItem('IDToken') != null)
		{
			this.id = sessionStorage.getItem('IDToken');
		}
		window.addEventListener('logged', (event) => {
    		this.id = event.detail.id;
  		});
		window.addEventListener('logout', (event) => {
    		this.logged = event.detail.logged;
			this.id = event.detail.id;
  		});
	}
}
</script>

<template>

<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container-fluid">
      <a class="navbar-brand" href="#">WASA Photo</a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarColor01" aria-controls="navbarColor01" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
	
      <div v-if="id != -1" class="collapse navbar-collapse" id="navbarColor01">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
			<li class="nav-item">
				<RouterLink to="/home" class="nav-link">
					<i class="fa-solid fa-house"></i>
					Home
				</RouterLink>
			</li>
			<li class="nav-item">
				<RouterLink to="/usersList" class="nav-link">
					<i class="fa-solid fa-address-book"></i>
					Users List
				</RouterLink>
			</li>
			<li class="nav-item">
				<RouterLink :to="'/profile/' + id" class="nav-link">
					<i class="fa-solid fa-user"></i>
					Profile
				</RouterLink>
			</li>
			<li class="nav-item">
				<RouterLink to="/newPost" class="nav-link">
					<i class="fa-solid fa-square-plus"></i>
					New Post
				</RouterLink>
			</li>
        </ul>
        <form class="d-flex" @submit.prevent="searchUser">
          <input class="form-control me-2" placeholder="Search" v-model="search">
          <button class="btn btn-outline-light no-border" type="submit" :disabled="search.length > 20"><i class="fa-solid fa-magnifying-glass"></i></button>
        </form>
      </div>
    </div>
  </nav>

  <main class="col-md-10 justify-content-center m-auto">
		<RouterView
		:search = "search"/>
  </main>
</template>

<style>
</style>
