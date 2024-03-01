<script>
import router from '../router';

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			users: null,
		}
	},
	props: ['search'],

	watch: {
		search: function () {
			this.refresh()
		},
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users", {
					params: {
						username: this.search,
					},
				});
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async deleteUser(id) {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/users/" + id);
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		userProfile(id) {
			this.$router.replace('/profile/' + id);
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h2>Users List</h2>
			<div class="btn-toolbar mb-2 mb-md-0">
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="row justify-content-md-center">
			<div class="col-sm-10">
				<div class="card m-2">
					<div class="card-header">
						<h3>List of the users registered</h3>
						<small>you can filter using the search bar at top-right</small>
					</div>
					<div class="card-body">
						<ul class="list-group" v-if="users != null">
							<li class="list-group-item no-border" v-for="u in users" v-bind:key="u.ID"> {{ u.Username }}
								<a href="javascript:" @click="(userProfile(u.ID))">[View Profile]</a>
								<!-- <a href="javascript:" @click="(deleteUser(u.ID))">[Elimina]</a> -->
							</li>
						</ul>
						<ul v-else>
							<li>No user found</li>
						</ul>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style></style>
