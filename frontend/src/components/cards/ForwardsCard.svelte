<script lang="ts">
    import CardWrapper, { CardWrapperProps } from './CardWrapper.svelte';
    import Bar from './rows/Bar.svelte';
    import Row from './rows/Row.svelte';

    type ForwardsCardProps = Omit<CardWrapperProps, 'type' | 'heading'> & {
        iValue: number;
        oValue: number;
        totalRow?: boolean;
        revenueRow?: boolean;
        revenueValue?: number;
    };

    export let props: ForwardsCardProps;
    
    $: props.iValue = Number(props.iValue)
    $: props.oValue = Number(props.oValue)
    $: totalValue = props.iValue + props.iValue
    
</script>

<CardWrapper props={{ type: 'forwards-card', heading: 'Forwards' }}>
    <Bar oValue={props.oValue} iValue={props.iValue} percent={true} />

    {#if props.totalRow}
        <Row label="Total" value={totalValue} />
    {/if}

    {#if props.revenueRow}
        <Row label="Revenue" value={props.revenueValue} />
    {/if}
</CardWrapper>
