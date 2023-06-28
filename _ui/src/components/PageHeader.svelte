<script>
    /**
     * Размер модального окна (xl, lg, md, sm)
     *
     * @type {string}
     */
    export let title;

    /**
     * Класс для первой колонки (по умолчанию: "col")
     *
     * @type {string}
     */
    export let firstColClass = 'col';
</script>

<div class="row">
    <div class="{firstColClass}">
        <h3 class="m-0">{title}</h3>
    </div>
    <slot />
</div>
