

<script lang="ts>
    import { onMount } from 'svelte';
    let {id} = $props()

    let thumbnailWidth = $state(0)
    let thumbnailHeight = $state(0)
    let available = $state(false)

    onMount( async()=> {
        const resp = await fetch(`/api/asset/${id}/thumbnail/meta/`)
        const obj = await resp.json()

        const TARGET_HEIGHT = 200
        const ratio = TARGET_HEIGHT / obj.thumbnail_height

        thumbnailWidth = obj.thumbnail_width * ratio
        thumbnailHeight = TARGET_HEIGHT
        available = obj.available
    })
</script>

<div 
    class={`h-[${thumbnailHeight}px] rounded-xl m-1 overflow-hidden`} 
    style={`width: ${thumbnailWidth}px;`}
>
    <a 
        href={available? `/api/asset/${id}/view/`: ''} 
        class:disabled={!available} 
    >
        
    <img 
        width={thumbnailWidth}
        height={thumbnailHeight}  
        src={`/api/asset/${id}/thumbnail`} 
        alt='{id}'>

    </a>
</div>