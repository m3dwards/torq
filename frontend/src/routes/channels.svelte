
<script lang="ts">
    import ChannelsTable from '../components/ChannelsTable.svelte';
    import { grpc } from '@improbable-eng/grpc-web';
    import { onMount, beforeUpdate} from 'svelte';
    import {torqrpcClientImpl, GrpcWebImpl} from '../torqrpc/torq'
    let channels;

    const rpc = new GrpcWebImpl('https://localhost:50051', {
        debug: false,
        metadata: new grpc.Metadata({}),
    });

    const client = new torqrpcClientImpl(rpc);
    let getChannels = async function() {
        return client.GetChannelFlow({
            chanIds: [772511372188909569],
            fromTime: 0,
            toTime: 1643555040,
        })
    }
	
</script>

<svelte:head>
    <title>Torq | Channels</title>
</svelte:head>

<div class="index-page">
    {#await getChannels()}
        Waiting for channels
    {:then channels}
        {console.log(channels)}
<!--        <ChannelsTable channels={channels} />    -->
    {/await}
    
</div>
