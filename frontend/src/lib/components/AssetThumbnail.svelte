

<script lang="ts>
    import { onMount } from 'svelte';
    let {id} = $props()

    let thumbnailWidth = $state(0)
    let thumbnailHeight = $state(0)

    onMount( async()=> {
        const resp = await fetch(`/api/asset/${id}/thumbnail/meta/`)
        const obj = await resp.json()

        const TARGET_HEIGHT = 200
        const ratio = TARGET_HEIGHT / obj.thumbnail_height

        thumbnailWidth = obj.thumbnail_width * ratio
        thumbnailHeight = TARGET_HEIGHT
    })
</script>

<img 
    width={thumbnailWidth}
    height={thumbnailHeight}  
    src={`/api/asset/${id}/thumbnail`} alt='{id}' class='h-[200px] max-w-full rounded-xl'>