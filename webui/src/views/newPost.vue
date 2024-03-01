<script>
import router from '../router';

export default {
	data: function () {
		return {
			errormsg: null,
			previewImage: undefined,
			fileName: "",
			file: "",
			caption: "",
		}
	},
	methods: {
		async submitFile() {
			let id = sessionStorage.getItem('IDToken');
			let formData = new FormData();

			formData.append('file', this.file);
			formData.append('caption', this.caption);
			formData.append('fileName', this.fileName);
			try {
				await this.$axios.post("/users/" + id + "/photos",
					formData, {
					headers: {
						'Content-Type': 'multipart/form-data'
					}
				}).then(
					this.$router.replace("/profile/" + id)
				);
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		handleFileUpload() {
			this.file = this.$refs.file.files[0];
			this.previewImage = URL.createObjectURL(this.file);
			this.fileName = this.file.name;
		}
	},
	mounted() {

	},
};
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h2>New Post</h2>
			<div class="btn-toolbar mb-2 mb-md-0">
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="row justify-content-md-center">
			<div class="col-sm-10">
				<div class="card m-2">
					<div class="card-header">
						<h5>Pick a photo, write something and create a post!</h5>
					</div>
					<div class="card-body">
						<div class="col-md-6 m-2 p-2">
							<label class="btn btn-default p-0">
								<input type="file" id="file" class="form-control" accept="image/*" ref="file" @change="handleFileUpload" />
							</label>
						</div>

						<div class="col-md-6  m-2 p-2">
							<input type="text" class="form-control" v-model="caption" placeholder="insert caption">
						</div>

						<button class="btn btn-success btn-sm m-2 p-2" :disabled="!file || !caption"
							@click="submitFile">
							Upload
						</button>

						<div v-if="previewImage">
							<div>
								<img class="preview my-3" :src="previewImage" alt="" />
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.preview {
	max-width: 500px;
}
</style>