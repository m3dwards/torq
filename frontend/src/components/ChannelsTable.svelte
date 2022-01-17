<script>
    import ChannelColumn from '../components/ChannelColumn.svelte';
    import OtherValues from '../components/cards/OtherValues.svelte';
    import Row from '../components/cards/rows/Row.svelte';

    import FeeCard from '../components/cards/FeeCard.svelte';
    import ForwardsCard from '../components/cards/ForwardsCard.svelte';
    import RebalancingCard from '../components/cards/RebalancingCard.svelte';
    import CapacityCard from '../components/cards/CapacityCard.svelte';
    import SummaryCard from '../components/cards/SummaryCard.svelte';

    export let channels;

    console.log("channels", channels);

</script>

<div class="channels-table-wrapper">
    <div class="channels-table">
        {#each Object.entries(channels.nodes) as node}
            {#each Object.entries(node[1].channels) as channel}
                {#if channel[1].channel.active }
                    <ChannelColumn alias={node[1].alias}>
                        <FeeCard
                            props={{
                                oValue: channel[1].feePerMil,
                                iValue: 100,
                                feeRange: true,
                                feeRangeLower: 100,
                                feeRangeUpper: 400,
                            }}
                        />

                        <ForwardsCard
                            props={{
                                oValue: channel[1].channel.totalSatoshisSent,
                                iValue: channel[1].channel.totalSatoshisReceived,
                                totalRow: true,
                                revenueRow: true,
                                revenueValue: 2980,
                            }}
                        />

                        <RebalancingCard
                            props={{
                                oValue: 30000000,
                                iValue: 0,
                                countRow: true,
                                countValue: 19,
                                costRow: true,
                                costValue: 503,
                            }}
                        />

                        <CapacityCard
                            props={{
                                oValue: channel[1].channel.localBalance,
                                iValue: channel[1].channel.remoteBalance,
                                totalRow: true,
                            }}
                        />

                        <OtherValues>
                            <Row label="Loop Out" value={25000} />
                            <Row label=" Channel Open Cost" value={550} />
                        </OtherValues>

                        <SummaryCard
                            props={{
                                costRow: true,
                                costValue: 1053,
                                revenueRow: true,
                                revenueValue: 2890,
                                profitRow: true,
                                profitValue: 1837,
                                aprRow: true,
                                aprValue: 0.045,
                            }}
                        />
                    </ChannelColumn>
                {/if}
            {/each}
        {/each}
    </div>
</div>

<style>
    .channels-table {
        margin-right: 20px;
        display: flex;
        flex-flow: column nowrap;
        flex-direction: row;
        justify-content: space-between;
        align-items: stretch;
        column-gap: 20px;

        overflow: auto;
        -ms-overflow-style: none; /* Internet Explorer 10+ */
        scrollbar-width: none; /* Firefox */
    }
    /* .channels-table-wrapper {
        position: relative;
    }
    .channels-table-wrapper::before {
        content: " ";
        width: 10px;
        height: 100%;
        position: absolute;
        top: 0;
        left: 0;
        background: linear-gradient(90deg, #F3F4F5 0%, rgba(243, 244, 245, 0) 100%);
    } */
    .channels-table::-webkit-scrollbar {
        width: 0; /* Remove scrollbar space */
        background: transparent; /* Optional: just make scrollbar invisible */
    }
</style>
