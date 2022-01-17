<script lang="ts">
    import CardWrapper, { CardWrapperProps } from './CardWrapper.svelte';
    import Bar from './rows/Bar.svelte';
    import Row from './rows/Row.svelte';

    type FeeCardProps = Omit<CardWrapperProps, 'type' | 'heading'> & {
        iValue: number | undefined;
        oValue: number | undefined;
        feeRange: boolean | undefined;
        feeRangeLower: number | undefined;
        feeRangeUpper: number | undefined;
    };

    $: props.iValue = Number(props.iValue)
    $: props.oValue = Number(props.oValue)

    export let props: FeeCardProps;
</script>

<CardWrapper props={{ type: 'fee-card', heading: 'Fee' }}>
    <Bar oValue={props.iValue} iValue={props.oValue} percent={false} postfix="ppm" />

    {#if props.feeRange}
        <Row
            label="Range"
            fromValue={props.feeRangeLower}
            toValue={props.feeRangeUpper}
            postfix="ppm"
        />
    {/if}
</CardWrapper>
