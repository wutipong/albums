<script lang="ts>
    let {album} = $props()

    const IMAGE_WIDTH = 300
    const IMAGE_HEIGHT = 200
</script>

<div class='card bg-base-100 w-[300px] h-[300px] hover:bg-base-300 hover:shadow-xl m-4'>
    <figure>
        <a href={`/album/${album.id}/`}>
            <img 
                src={`/api/album/${album.id}/cover`} 
                alt={album.id} 
                width={IMAGE_WIDTH} 
                height={IMAGE_HEIGHT}>
        </a>
    </figure>
    <div class='card-body h-4 overflow-clip'>
         <a href={`/album/${album.id}/`}>
            {album.name}
        </a>
    </div>
</div>
