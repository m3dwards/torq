<script lang="ts">
    import ChannelColumn from '../components/ChannelColumn.svelte';
    import HorizontalBarsCard from './cards/HorizontalBarsCard.svelte';

    export let channels;

    function getGroupName(name) {
        if (name.length > 20) {
            return name.substring(0,20) + "..."
        }
        return name
    }

</script>

<div class="channels-table-wrapper">
    <div class="channels-table">
        {#each channels.aggregatedForwards as fw}
                    <ChannelColumn alias={getGroupName(fw.groupName) || getGroupName(fw.groupId)}>

                        <HorizontalBarsCard
                            props={{
                                type: "forward-amount",
                                heading: "Forwarded amount",
                                oValue: fw.amountOut,
                                iValue: fw.amountIn,
                                totalRow: true
                            }}
                        />

                        <HorizontalBarsCard
                                props={{
                                type: "forward-fee",
                                heading: "Fees earned",
                                oValue: fw.feeOut,
                                iValue: fw.feeIn,
                                totalRow: false,
                                revenueRow: false
                            }}
                        />

                        <HorizontalBarsCard
                                props={{
                                type: "forward-count",
                                heading: "Forwards count",
                                oValue: fw.countOut,
                                iValue: fw.countIn,
                                totalRow: true
                            }}
                        />


                    </ChannelColumn>
        {/each}
    </div>
</div>

<style>
    .channels-table-wrapper {
        display: flex;
        flex-direction: column;
        align-self: stretch;
    }
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
