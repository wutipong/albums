<script lang="ts>
    import notAvailableSvg from '$lib/assets/not-available-small.svg?raw'
    let {album} = $props()
</script>

<div 
    class='card w-[300px] h-[300px] border-base-300 border-1 shadow hover:bg-base-100 hover:shadow-xl m-4'
    id={album.id}
    >
    <figure>
        <a href={`/album/${album.id}/`}>
            {#if album.cover_url == ''}
                {@html notAvailableSvg}
            {:else}
                <img 
                    src={album.cover_url} 
                    alt={album.id} 
                    width='300' 
                    height='200'
                />
            {/if}
        </a>
    </figure>
    <div class='card-body h-4 overflow-clip'>
         <a href={`/album/${album.id}/`}>
            {album.name}
        </a>
    </div>
</div>
