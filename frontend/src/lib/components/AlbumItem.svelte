<script lang="ts>
    let {id, title} = $props()

    const IMAGE_WIDTH = 300
    const IMAGE_HEIGHT = 200
</script>

<div class='card bg-base-100 w-[300px] h-[300px]'>
    <figure>
        <a href={`/album/${id}/`}>
            <img src={`/api/album/${id}/cover`} alt={id} width={IMAGE_WIDTH} height={IMAGE_HEIGHT}>
        </a>
    </figure>
    <div class='card-body h-4 overflow-clip'>
        <p>{title}</p>
    </div>
</div>
