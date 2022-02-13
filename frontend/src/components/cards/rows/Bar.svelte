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
    <div class="bar-values">
        <div class="outbound">
            <FormatNumber value={oValue} decimals={0} notation="standard" />
            <span class="percent">
                &ensp;
                (<FormatNumber value={oValuePercent} decimals={0} />%)
            </span>
        </div>
    </div>
    <div class="bars">
        <div class="bar-value out" style="width: {oValuePercent}%" />
        <div class="bar-value inn" style="width: {iValuePercent}%" />
    </div>
    <div class="bar-values">
        <div class="inbound">
            <FormatNumber value={iValue} decimals={0} notation="standard" />
            <span class="percent">
                &ensp;
                (<FormatNumber value={iValuePercent} decimals={0} />%)
            </span>
        </div>
    </div>
</div>

<style>
    .bars {
        margin-bottom: 4px;
        margin-top: 2px;
    }
    .bar-values {
        display: flex;
        justify-content: space-between;
    }
    .bar-row {
        padding: 8px 15px 5.5px;
        font-size: 14px;
        align: bottom;
        line-height: 28px;
        color: #3A463C;
    }
    .percent {
        color: rgba(58, 70, 60, 0.35);
    }
    .bar-value {
        height: 10px;
        width: 0%;
        min-width: 2px;
    }

    .bar-value.out {
        background-color: #66786A;
        margin-bottom: 2px;
    }
    .bar-value.inn {
        background-color: #C7D1C9;
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
