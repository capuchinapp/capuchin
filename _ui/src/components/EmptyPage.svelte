<div class="text-center pb-5">
    <div class="image">
        <slot name="image" />
    </div>
    <h3>
        <slot name="title" />
    </h3>
    <div class="mt-3">
        <slot name="footer" />
    </div>
</div>

<style>
    .image :global(img) {
        height: 30vh;
    }
</style>
