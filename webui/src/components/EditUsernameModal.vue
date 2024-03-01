<script>
export default {	
	data(){
		return{
			username:"",
            errormsg: null
		}
	},
    props: ['id'],

	methods: {
		async changeUsername() {
			this.errormsg = null;
			try {
				let response = await this.$axios.put("/users/"+ this.id, {
                    username:	this.username
                }).then(response => {
                    if(response.data.Message != null) {
                        this.errormsg = response.data.Message;
                    } else {
                        this.$emit('refresh');
                        document.getElementById('close').click();
                    }
                    this.username = "";
   				});;
			} catch (e) {
				this.errormsg = e.toString();
			}

            
		}
	},
}
</script>

<template>
    <div class="modal fade my-modal-disp-none" :id="'EditUsername'" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog modal-dialog-scrollable ">
            <div class="modal-content">

                <div class="modal-header">
                    <h1 class="modal-title fs-5" :id="'EditUsername'">Change Username</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close" id="close"></button>
                </div>

                <div class="modal-body d-flex justify-content-center w-100">
                    <div class="row w-100">
                        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
                        <input type="text" class="form-control" v-model="username">
                    </div>
                </div>
                <div class="modal-footer d-flex justify-content-center w-100">
                    <div class="row w-100 ">
                            <button type="button" class="btn btn-primary" @click.prevent="changeUsername" :disabled="username.length < 3 || username.length > 20">
							    Save
							</button>
                    </div>
                </div>
            </div>
        </div>
    </div>

</template>

<style> 
.my-modal-disp-none{
	display: none;
}
</style>
