<script lang="ts">
    import {authClient} from '$lib/auth-client'
    async function performLogin() {
        
        await authClient.signIn.oauth2({
            providerId: "generic-oauth",
            callbackURL: "/album",
            errorCallbackURL: "/error",
            // disableRedirect: true,
        })
    
    }
</script>

<div class="w-screen h-screen">
    <div class="flex w-full h-full justify-center items-center"> 
    <button class="btn btn-primary btn-xl "
        onclick={()=>performLogin()}
    >
        Login
    </button>
    </div>
</div>