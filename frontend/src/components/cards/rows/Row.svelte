<script lang="ts">
    import FormatNumber from '../../helpers/FormatNumber.svelte';

    export let label: string;
    export let postfix: string | undefined = undefined;
    export let style: 'percent' | undefined = undefined;
    export let notation: 'compact' | 'standard' | undefined = undefined;
    export let decimals: number | undefined = undefined;
    export let value: number | undefined = 0;
    export let fromValue: number | undefined = undefined;
    export let toValue: number | undefined = undefined;

    if (fromValue) {
        value = fromValue;
    }
</script>

<div class="row">
    <span class="label">
        {label}
    </span>
    <span class="value">
        <FormatNumber {value} {style} {decimals} {notation} />
        {#if toValue}
            -
            <FormatNumber value={toValue} {style} {decimals} {notation} />
        {/if}
        {postfix ? postfix : ' '}
    </span>
</div>

<style>
    .row {
        display: flex;
        justify-content: space-between;
        padding: 5.5px 15px 5.5px;
        margin-top: 0px;
        border-top: 1px solid #f3f4f5;
        font-size: 14px;
        line-height: 200%;
        color: #3A463C;
    }
</style>
