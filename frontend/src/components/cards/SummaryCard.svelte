<script lang="ts">
    import CardWrapper, { CardWrapperProps } from './CardWrapper.svelte';
    import Row from './rows/Row.svelte';

    type SummaryCardProps = Omit<CardWrapperProps, 'type' | 'heading'> & {
        costRow?: boolean;
        costValue?: number;
        revenueRow?: boolean;
        revenueValue?: number;
        profitRow?: boolean;
        profitValue?: number;
        aprRow?: boolean;
        aprValue?: number;
    };

    export let props: SummaryCardProps;
</script>

<CardWrapper props={{ type: 'summary-card', heading: 'Summary' }}>
    {#if props.revenueRow}
        <Row label="Revenue" value={props.revenueValue} />
    {/if}

    {#if props.costRow}
        <Row label="Cost" value={props.costValue} />
    {/if}

    {#if props.profitRow}
        <Row label="Profit" value={props.profitValue} />
    {/if}

    {#if props.aprRow}
        <Row label="APR" value={props.aprValue} style={'percent'} decimals={2} />
    {/if}
</CardWrapper>
