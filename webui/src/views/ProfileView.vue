<script>
import router from '../router';

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			id: this.$route.params.id,
			resp: null,
			posts: null,
			nPhoto: null,
			nFollower: null,
			nFollowing: null,
			banState: false,
			followState: false,
			loggedID: sessionStorage.getItem('IDToken'),
			comment: ""
		}
	},
	created() {
		this.$watch(
			() => this.$route.params.id,
			(toParams, previousParams) => {
				this.id = this.$route.params.id,
					this.refresh()
			}
		)
	},
	computed: {
		checkUser() {
			return this.id == this.loggedID
		},
	},
	methods: {
		async refresh() {
			if (this.$route.params.id === undefined) {
				return
			}
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + this.id);
				if (response.data.posts != null) {
					response.data.posts.forEach(post => {
						post['image'] = __API_URL__ + "/users/" + this.id + "/photos/" + post.ID
					});
				}
				this.resp = response.data;
				this.posts = this.resp.posts;
				this.nPhoto = this.resp.posts != null ? this.resp.posts.length : 0
				this.nFollower = this.resp.followers != null ? this.resp.followers.length : 0
				this.nFollowing = this.resp.following != null ? this.resp.following.length : 0
				this.followState = this.resp.followers != null ? response.data.followers.find(obj => obj.ID == this.loggedID) : false
				this.banState = this.resp.isBanned
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async follow() {
			try {
				if (this.followState) {
					await this.$axios.delete("/users/" + this.id + "/follow/" + this.loggedID);
					this.nFollower -= 1
				} else {
					await this.$axios.put("/users/" + this.id + "/follow/" + this.loggedID);
					this.nFollower += 1
				}
				this.followState = !this.followState
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async ban() {
			try {
				if (this.banState) {
					await this.$axios.delete("/users/" + this.id + "/ban/" + this.loggedID);
				} else {
					await this.$axios.put("/users/" + this.id + "/ban/" + this.loggedID);
					this.followState = false
				}
				this.banState = !this.banState
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async like(i, post_id) {
			try {
				if (i) {
					await this.$axios.put("/users/" + this.id + "/photos/" + post_id + "/like/" + this.loggedID);
					this.refresh();
				} else {
					await this.$axios.delete("/users/" + this.id + "/photos/" + post_id + "/like/" + this.loggedID);
					this.refresh();
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async removePost(post_id) {
			if (confirm("Do you really want to delete?")) {
				try {
					await this.$axios.delete("/users/" + this.id + "/photos/" + post_id);
					this.refresh()
				} catch (e) {
					this.errormsg = e.toString();
				}
			}
		},
		async addComment(post_id) {
			try {
				await this.$axios.post("/users/" + this.id + "/photos/" + post_id + "/comments", {
					comment: this.comment
				});
				this.refresh()
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.comment = "";
			document.getElementById(post_id).value = "";
			document.getElementById(10 + post_id).disabled = true;
		},
		async removeComment(post_id, comment_id) {
			if (confirm("Do you really want to delete?")) {
				try {
					await this.$axios.delete("/users/" + this.id + "/photos/" + post_id + "/comments/" + comment_id);
					this.refresh()
				} catch (e) {
					this.errormsg = e.toString();
				}
			}
		},
		newPost() {
			this.$router.replace("/newPost");
		},
		checkLike(list) {
			if (list == null) {
				return false
			} else {
				return list.find(obj => obj.User_ID == this.loggedID)
			}
		},
		countList(list) {
			if (list == null) {
				return 0
			} else {
				return list.length
			}
		},
		getText(id) {
			this.comment = document.getElementById(id).value;
			if (this.comment.length > 1)
				document.getElementById(10 + id).disabled = false;
			else
				document.getElementById(10 + id).disabled = true;
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
			<h1 class="h1" v-if="resp">{{ resp.username }}'s profile</h1>
			<div v-if="checkUser" class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newPost()">
						New Post
					</button>
					<button type="button" class="btn btn-sm btn-outline-primary" data-bs-toggle="modal"
						:data-bs-target="'#EditUsername'">
						Change Username
					</button>

					<EditUsernameModal :id=this.id @refresh="refresh" />
				</div>
			</div>
			<div v-else class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button v-if="banState == false" type="button" class="btn btn-sm btn-outline-success" @click="follow()">
						{{ followState ? "Unfollow" : "Follow" }}
					</button>
					<button type="button" class="btn btn-sm btn-outline-danger" @click="ban()">
						{{ banState ? "Unban" : "Ban" }}
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div v-if="!this.banState">
			<div class="row pb-2 mb-3 border-bottom">
				<div class="col-md-4 h3">Post: {{ nPhoto }}</div>
				<div class="col-md-4 h3">Follower: {{ nFollower }} </div>
				<div class="col-md-4 h3">Following: {{ nFollowing }}</div>
			</div>

			<div v-if="posts != null" class="row row-cols-1 row-cols-md-3 g-4">
				<div class="col mb-2" v-for="p in posts" v-bind:key="p.ID">
					<div class="card h-100">
						<img class="card-img-top" :src="p.image">
						<div class="card-body">
							<div class="container">
								<div class="d-flex flex-row justify-content-end align-items-center mb-2">
									<button class="btn btn-outline-info btn-info m-0 p-1 me-auto">
										<i> From {{ resp.username }}</i> </button>
									<button v-if="checkUser" class="btn btn-outline-danger me-2 " @click="removePost(p.ID)">
										<i class="fa-solid fa-trash w-100 h-100"></i>
									</button>
									<button v-if="checkLike(p.Likes)" class="btn btn btn-outline-primary" :disabled="checkUser"
										@click="like(false, p.ID)">
										<i class="fa-solid fa-thumbs-up me-2"></i>{{ countList(p.Likes) }}
									</button>
									<button v-if="!checkLike(p.Likes)" class="btn btn btn-outline-primary" :disabled="checkUser"
										@click="like(true, p.ID)">
										<i class="fa-regular fa-thumbs-up me-2"></i>{{ countList(p.Likes) }}
									</button>
								</div>
								<figure>
									<blockquote class="blockquote">
										<p>{{ p.Caption }}</p>
									</blockquote>
									<figcaption class="blockquote-footer">
										Uploaded on <cite title="date">{{ p.Created_Datetime }} </cite>
									</figcaption>
								</figure>
								<div class="container">
									<div class="d-flex flex-row justify-content-end align-items-center mb-2">
										<p class="h4 m-0 p-1 me-auto">Comments {{ countList(p.Comments) }}</p>
										<button v-if="p.Comments != null && p.Comments.length > 1" type="button"
											class="btn btn-sm btn-outline-secondary" data-bs-toggle="modal"
											:data-bs-target="'#AllComment' + p.ID">
											See all
										</button>
									</div>
									<div v-if="p.Comments != null">
										<div v-if="p.Comments.length > 1">
											<div class="card mb-4">
												<div class="card-body">
													<p>{{ p.Comments[p.Comments.length - 1].Text }}</p>
													<div class="d-flex justify-content-between">
														<div class="d-flex flex-row align-items-center">
															<p class="small mb-0 ms-1">{{ p.Comments[p.Comments.length -
																1].Username }}</p>
														</div>
														<div class="d-flex flex-row align-items-center">
															<button
																v-if="p.Comments[p.Comments.length - 1].User_ID == this.loggedID || checkUser"
																class="btn btn-outline-danger mx-2 fa-xs "
																@click="removeComment(p.ID, p.Comments[p.Comments.length - 1].ID)">
																<i class="fa-solid fa-trash"></i>
															</button>
														</div>
													</div>
												</div>
											</div>
											<AllCommentModal
												:comments=p.Comments
												:loggedID=this.loggedID
												:post=p
												:id = this.id

												@refresh="refresh"
											/>
										</div>
										<div v-else>
											<div v-for="c in p.Comments" v-bind:key="c.ID" class="card mb-4">
												<div class="card-body">
													<p>{{ c.Text }}</p>
													<div class="d-flex justify-content-between">
														<div class="d-flex flex-row align-items-center">
															<p class="small mb-0 ms-1">{{ c.Username }}</p>
														</div>
														<div class="d-flex flex-row align-items-center">
															<button v-if="c.User_ID == this.loggedID || checkUser"
																class="btn btn-outline-danger mx-2 fa-xs "
																@click="removeComment(p.ID, c.ID)">
																<i class="fa-solid fa-trash"></i>
															</button>
														</div>
													</div>
												</div>
											</div>
										</div>
									</div>
									<div v-else>
										<div class="card mb-4">
											<div class="card-body">
												<p> This post has no comment, go add the first one below! </p>
											</div>
										</div>
									</div>
								</div>
							</div>
						</div>
						<div class="card-footer">
							<div class="d-flex flex-row justify-content-end align-items-center mb-2">
								<p class="h4 m-0 p-1 me-auto">New Comment</p>
							</div>
							<div class="input-group mb-3">
								<textarea class="form-control" rows="2" placeholder="Write something..." :id="p.ID"
									@change="getText(p.ID)"></textarea>
								<button class="btn btn-outline-success" type="button" :id="10 + p.ID" disabled
									@click="addComment(p.ID)">Add</button>
							</div>
						</div>
					</div>
				</div>
			</div>
			<div v-else>
				No post yet
			</div>
		</div>
		<div v-else>
			This user is banned from you, you must unban him to see his profile
		</div>
	</div>
</template>

<style>
textarea {
	resize: none;
}
</style>
