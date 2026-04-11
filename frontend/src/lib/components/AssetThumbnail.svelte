

<script lang="ts>
    import { onMount } from 'svelte';
    let {id} = $props()

    let thumbnailWidth = $state(0)
    let thumbnailHeight = $state(0)
    let available = $state(false)

    let preview = $state(false)

    onMount( async()=> {
        const resp = await fetch(`/api/asset/${id}/`)
        const obj = await resp.json()

        const TARGET_HEIGHT = 200
        const ratio = TARGET_HEIGHT / obj.thumbnail_height

        thumbnailWidth = obj.thumbnail_width * ratio
        thumbnailHeight = TARGET_HEIGHT
        available = obj.available
    })
</script>

<div 
    role='button'
    tabindex='0'
    class={`h-[${thumbnailHeight}px] rounded-xl m-1 overflow-hidden`} 
    style={`width: ${thumbnailWidth}px;`}
    onmouseenter={()=>preview = true}
    onmouseleave={()=>preview = false}
>
    <a 
        href={available? `/api/asset/${id}/view/`: ''} 
        class:disabled={!available} 
    >
        
    <img 
        width={thumbnailWidth}
        height={thumbnailHeight}  
        src={`/api/asset/${id}/thumbnail`} 
        alt='{id}'
        class:hidden={preview}
    >

   
  <div 
    class:hidden={!preview} 
    class='border-4 rounded-xl overflow-hidden w-full h-full box-border'>
    <img 
        width={thumbnailWidth}
        height={thumbnailHeight}  
        src={`/api/asset/${id}/preview`} 
        alt='{id}'
        class='w-full h-full'
    >
    </div>
    </a>
</div>