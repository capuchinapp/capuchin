<script lang="ts">
    import {capuchin} from './../stores';
</script>

<div class="row">
    <div class="col-md-4 d-flex align-items-center">
        <!-- svelte-ignore a11y_invalid_attribute -->
        &copy; DimNS {$capuchin.appVersionFront}
    </div>
    <div class="col-md-4 d-md-flex align-items-center justify-content-md-center">
        &nbsp;
    </div>
    <div class="col-md-4 text-md-end">
        <a href="https://www.flaticon.com/free-icons/amazon" title="amazon icons" class="attribution">
            Amazon icons created by Freepik - Flaticon
        </a>
    </div>
</div>

<style>
    .attribution {
        font-size: 0.8em;
    }
</style>
