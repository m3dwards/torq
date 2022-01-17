<script lang="ts">
    import CardWrapper, { CardWrapperProps } from './CardWrapper.svelte';
    import Bar from './rows/Bar.svelte';
    import Row from './rows/Row.svelte';

    type RebalancingCardProps = Omit<CardWrapperProps, 'type' | 'heading'> & {
        iValue: number;
        oValue: number;
        countRow?: boolean;
        countValue?: number;
        costRow?: boolean;
        costValue?: number;
    };

    export let props: RebalancingCardProps;
</script>

<CardWrapper props={{ type: 'rebalancing-card', heading: 'Rebalancing' }}>
    <Bar oValue={props.oValue} iValue={props.iValue} percent={true} />

    {#if props.countRow}
        <Row label="Count" value={props.countValue} />
    {/if}

    {#if props.costRow}
        <Row label="Cost" value={props.costValue} />
    {/if}
</CardWrapper>
