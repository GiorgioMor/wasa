<script>
export default {	
	data(){
		return{
			
		}
	},
    props: ['id', 'comments', 'loggedID', 'post'],
    computed: {
		checkUser() {
            if (this.id != -1)
			    return this.id == this.loggedID
            else return false
		},
	},
	methods: {
        async removeComment(author_id, post_id, comment_id){
			if(confirm("Do you really want to delete?")) {
				try{
					await this.$axios.delete("/users/"+ author_id +"/photos/"+ post_id +"/comments/"+ comment_id)
                    .then(response => {
                        this.$emit('refresh');
   				});
				}catch(e){
					this.errormsg = e.toString();
				}

                document.getElementById('close').click();
			}
		}
    }
}
</script>

<template>
    <div class="modal fade my-modal-disp-none" :id="'AllComment' + post.ID" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog modal-dialog-scrollable ">
            <div class="modal-content">

                <div class="modal-header">
                    <h1 class="modal-title fs-5" :id="'AllComment'">Comment List</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close" id="close"></button>
                </div>

                <div class="modal-body justify-content-center w-100">
                        <div v-for="c in comments" v-bind:key="c.ID"  class="card mb-4">
                            <div class="card-body">
                                <p>{{ c.Text }}</p>
                                <div class="d-flex justify-content-between">
                                    <div class="d-flex flex-row align-items-center">
                                        <p class="small mb-0 ms-1">{{ c.Username }}</p>
                                    </div>
                                    <div class="d-flex flex-row align-items-center">
                                        <button v-if=" c.User_ID == loggedID || checkUser" class="btn btn-outline-danger mx-2 fa-xs " @click="removeComment(post.User_ID, post.ID, c.ID)">
                                            <i class="fa-solid fa-trash"></i>
                                        </button>
                                    </div>
                                </div>
                            </div>
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
