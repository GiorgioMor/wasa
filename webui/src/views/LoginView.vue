<script>
export default {
	data: function() {
		return {
			errormsg: null,
			username: "",
			token: -1,
		}
	},
	methods: {
		async login() {
			try {
				await this.$axios.post("/session", {
					username:	this.username
				}).then(response => {
        			sessionStorage.setItem('IDToken', response.data);
					this.token = response.data;
					window.dispatchEvent(new CustomEvent('logged', {
						detail: {
							id: sessionStorage.getItem('IDToken')
						}
					}));
					this.$router.replace("/home");
   				});
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async logout() {
			sessionStorage.setItem('IDToken', -1);
			this.token = -1;
			window.dispatchEvent(new CustomEvent('logout', {
						detail: {
							logged: false,
							id: -1							
						}
			}));
			this.$router.replace("/");
		}
	},
	mounted() {
		if (sessionStorage.getItem('IDToken') != null)
		{
			this.token = sessionStorage.getItem('IDToken');
		}
	}
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Login Page</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div v-if="token == -1"  class="row justify-content-md-center">
			<div class="col-sm-10">
				<div class="card border-light m-2">
					<div class="card-header text-center">
						<h3>Insert your useraname (3-20 character) for login or create a new account</h3>
					</div>
					<div class="card-body">
						<form @submit.prevent="login">
							<div class="row justify-content-md-center">
								<div class="col-sm-6">
									<input type="text" class="form-control" v-model="username" placeholder="insert username">
								</div>
							</div>
							<div class="row mt-2 justify-content-md-center">
								<button class="btn no-border" type="submit" :disabled="username.length < 3 || username.length > 20">Send</button>
							</div>
						</form>
					</div>
				</div>
			</div>
		</div>

		<div v-else>
			<p>Do you wanna logout? <a href="javascript:" @click="logout"><i class="fa-solid fa-right-from-bracket"></i></a></p>
		</div>
		
	</div>
</template>

<style>
	.no-border {
		border: none;
		background: transparent;
	}
</style>
