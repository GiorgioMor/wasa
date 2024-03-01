<script>
export default {
	data: function () {
		return {
			errormsg: null,
			posts: null,
			followers: null,
			loggedID: sessionStorage.getItem("IDToken"),
			comment: ""
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			let id = sessionStorage.getItem("IDToken");
			try {
				let response = await this.$axios.get("/users/" + id + "/home");
				if (response.data.posts != null) {
					response.data.posts.forEach(post => {
						post['image'] = __API_URL__ + "/users/" + post.User_ID + "/photos/" + post.ID
					});

					response.data.posts.sort((a, b) => { return a.Created_Datetime < b.Created_Datetime; });
				}
				this.posts = response.data.posts;
				this.followers = response.data.followers;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async like(i, userPostid, post_id) {
			try {
				if (i) {
					await this.$axios.put("/users/" + userPostid + "/photos/" + post_id + "/like/" + this.loggedID);
					this.refresh();
				} else {
					await this.$axios.delete("/users/" + userPostid + "/photos/" + post_id + "/like/" + this.loggedID);
					this.refresh();
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async addComment(author_id, post_id) {
			try {
				await this.$axios.post("/users/" + author_id + "/photos/" + post_id + "/comments", {
					comment: this.comment
				});
				this.refresh()
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.comment = ""
			document.getElementById(post_id).value = ""
			document.getElementById(10 + post_id).disabled = true;
		},
		async removeComment(author_id, post_id, comment_id) {
			if (confirm("Do you really want to delete?")) {
				try {
					await this.$axios.delete("/users/" + author_id + "/photos/" + post_id + "/comments/" + comment_id);
					this.refresh()
				} catch (e) {
					this.errormsg = e.toString();
				}
			}
		},
		newPost() {
			this.$router.replace("/newPost");
		},
		userProfile(id) {
			this.$router.replace('/profile/' + id);
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
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newPost()">
						New Post
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div v-if="posts != null" class="row row-cols-1 row-cols-md-3 g-4">
			<div class="col mb-2" v-for="p in posts" v-bind:key="p.ID">
				<div class="card h-100">
					<img class="card-img-top" :src="p.image">
					<div class="card-body">
						<div class="container">
							<div class="d-flex flex-row justify-content-end align-items-center mb-2">
								<button class="btn btn-outline-info btn-info m-0 p-1 me-auto"
									@click="(userProfile(p.User_ID))">
									<i> From {{ followers.find(obj => obj.ID == p.User_ID).Username }}</i>
								</button>
								<button v-if="checkLike(p.Likes)" class="btn btn btn-outline-primary"
									@click="like(false, p.User_ID, p.ID)">
									<i class="fa-solid fa-thumbs-up me-2"></i>{{ countList(p.Likes) }}
								</button>
								<button v-if="!checkLike(p.Likes)" class="btn btn btn-outline-primary"
									@click="like(true, p.User_ID, p.ID)">
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
															v-if="p.Comments[p.Comments.length - 1].User_ID == this.loggedID"
															class="btn btn-outline-danger mx-2 fa-xs "
															@click="removeComment(p.User_ID, p.ID, p.Comments[p.Comments.length - 1].ID)">
															<i class="fa-solid fa-trash"></i>
														</button>
													</div>
												</div>
											</div>
										</div>
										<AllCommentModal :comments=p.Comments :loggedID=this.loggedID :post=p :id = -1
											@refresh="refresh" />
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
														<button v-if="c.User_ID == this.loggedID"
															class="btn btn-outline-danger mx-2 fa-xs "
															@click="removeComment(p.User_ID, p.ID, c.ID)">
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
							<p class="h4 m-0 p-1 me-auto">New comment</p>
						</div>
						<div class="input-group mb-3">
							<textarea class="form-control" rows="2" placeholder="Write something..." :id="p.ID"
								@change="getText(p.ID)"></textarea>
							<button class="btn btn-outline-success" type="button" :id="10 + p.ID" disabled
								@click="addComment(p.User_ID, p.ID)">Add</button>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div v-else>
			No posts yet
		</div>
	</div>
</template>

<style></style>
