<script>
    export let alias = '';
</script>

<div class="channel-column">
    <div class="alias">
        {alias}
    </div>
    <div class="cards">
        <slot />
    </div>
</div>

<style>
    .channel-column {
        min-width: 344px;
        min-height: calc(100vh - 80px);
    }

    .alias {
        font-size: 20px;
        margin-bottom: 10px;
    }
</style>
