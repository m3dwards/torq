<script lang="ts">
    import CardWrapper, { CardWrapperProps } from './CardWrapper.svelte';
    import Bar from './rows/Bar.svelte';
    import Row from './rows/Row.svelte';

    type BarCardProps = CardWrapperProps & {
        iValue: number;
        oValue: number;
        totalRow?: boolean;
    };

    export let props: BarCardProps;
    
    $: props.iValue = Number(props.iValue)
    $: props.oValue = Number(props.oValue)
    $: totalValue = props.iValue + props.oValue

</script>

<CardWrapper props={{ type: props.type, heading: props.heading }}>
    <Bar oValue={props.oValue} iValue={props.iValue} percent={true} />

    {#if props.totalRow}
        <Row label="Total" value={totalValue} notation="standard" />
    {/if}

</CardWrapper>
