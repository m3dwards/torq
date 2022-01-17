<script lang="ts">
    import CardWrapper, { CardWrapperProps } from './CardWrapper.svelte';
    import Bar from './rows/Bar.svelte';
    import Row from './rows/Row.svelte';

    type CapacityCardProps = Omit<CardWrapperProps, 'type' | 'heading'> & {
        iValue: number;
        oValue: number;
        totalRow?: boolean;
    };

    export let props: CapacityCardProps;
    let totalCapacity = Number(props.oValue) + Number(props.iValue);
    
    
</script>

<CardWrapper props={{ type: 'capacity-card', heading: 'Capacity' }}>
    <Bar oValue={props.iValue} iValue={props.oValue} percent={true} />

    {#if props.totalRow}
        <Row label="Total" value={totalCapacity} />
    {/if}
</CardWrapper>
