<script lang="ts>
    import { onMount } from 'svelte';
    let {id} = $props()

    let thumbnailWidth = $state(0)
    let thumbnailHeight = $state(0)
    let available = $state(false)
    let videoDuration = $state('')
    let imageFrames = $state(1)
    let assetType = $state('image')

    let preview = $state(false)

    onMount( async()=> {
        const resp = await fetch(`/api/asset/${id}/`)
        const obj = await resp.json()

        const TARGET_HEIGHT = 200
        const ratio = TARGET_HEIGHT / obj.thumbnail_height

        thumbnailWidth = obj.thumbnail_width * ratio
        thumbnailHeight = TARGET_HEIGHT
        available = obj.available

        assetType = obj.type
        videoDuration = obj.video_duration.toString()
        imageFrames = obj.image_frames

        console.log(obj)
    })
</script>

<div 
    role='button'
    tabindex='0'
    class={`block h-[${thumbnailHeight}px] rounded-xl m-1 overflow-hidden`} 
    style={`width: ${thumbnailWidth}px;`}
    onmouseenter={()=>preview = true}
    onmouseleave={()=>preview = false}
>
    <a 
        class='block'
        href={available? `/api/asset/${id}/view/`: ''} 
        class:disabled={!available} 
    >
        <div class='relative w-full h-full'>
            <div 
                class:hidden={preview} 
                class='rounded-xl overflow-hidden w-full h-full box-border'
                style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
                >
                <img 
                    width={thumbnailWidth}
                    height={thumbnailHeight}  
                    src={`/api/asset/${id}/thumbnail`} 
                    alt='{id}'
                    class:hidden={preview}
                >   
            </div>

            <div 
                class:hidden={!preview} 
                class='border-4 rounded-xl overflow-hidden w-full h-full box-border'
                style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
                >
                <img 
                    width={thumbnailWidth}
                    height={thumbnailHeight}  
                    src={`/api/asset/${id}/preview`} 
                    alt='{id}'
                    class='w-full h-full'
                >
            </div>

            <div class='absolute top-1 right-2 place-items-end'>
                {#if assetType === 'video'}
                    <div class='badge' >Video</div>
                {/if}

                {#if imageFrames > 1}
                    <div class='badge'> Animation</div>
                {/if}
            </div>
        </div>
    </a>
</div>