<script lang="ts">
    import FormatNumber from '../../helpers/FormatNumber.svelte';

    export let postfix = '';
    export let iValue: number;
    export let oValue: number;

    export let values = true;
    export let percent = false;

    $: iValue = Number(iValue);
    $: oValue = Number(oValue);

    $: iValuePercent = (Number(iValue) / (Number(oValue) + Number(iValue))) * 100;
    $: oValuePercent = (Number(oValue) /  (Number(oValue) + Number(iValue))) * 100;
</script>

<div class="bar-row">
    <div class="bar-heading">
        <div class="outbound">Outbound</div>
        <div class="inbound">Inbound </div>
    </div>
    <div class="bar">
        <div class="bar-value" style="width: {oValuePercent}%" />
    </div>
    <div class="bar-values">
        <div class="outbound">
            {#if values}
                <FormatNumber value={oValue} decimals={0} notation="standard" />
                {postfix}
            {/if}
        </div>
        <div class="inbound">
            {#if values}
                <FormatNumber value={iValue} decimals={0} notation="standard" />
                {postfix}
            {/if}
        </div>
    </div>
    <div class="bar-percent">
        <span>
            {#if percent}
                <span class="percent">
                    <FormatNumber value={oValuePercent} decimals={0} /> %
                </span>
            {/if}
        </span>
        <span>
            {#if percent}
                <span class="percent">
                    <FormatNumber value={iValuePercent} decimals={0} /> %
                </span>
            {/if}
        </span>
    </div>
</div>

<style>
    .bar-values,
    .bar-percent {
        display: flex;
        justify-content: space-between;
        padding-top: 10px;
    }
    .bar-percent {
        padding-top: 5px;
        padding-bottom: 5px;
    }
    .bar-heading {
        display: flex;
        justify-content: space-between;
        font-size: 13px;
        color: rgba(0, 0, 0, 0.35);
        margin-bottom: 7px;
    }
    .bar-row {
        padding: 20px 15px 10px;
    }
    .percent {
        color: rgba(0, 0, 0, 0.35);
    }
    .bar {
        position: relative;
        background-color: #66786a;
        width: 100%;
        height: 5px;
        border-radius: 5px;
    }
    .bar-value {
        background-color: #fa9401;
        height: 5px;
        width: 0%;
        border-radius: 2px;
    }
    /* .bar-text {
        width: 100%;
        display: flex;
        justify-content: space-between;
        font-size: 16px;
        vertical-align: bottom;
        
        color: rgba(0, 0, 0, 0.35);
        opacity: 1;
        transition: opacity 180ms linear;
        transition-delay: 0s;
    }
    .bar-row:hover .bar-text {
        opacity: 1;
        transition: opacity 150ms linear;
        transition-delay: 300ms;
    } */
</style>
