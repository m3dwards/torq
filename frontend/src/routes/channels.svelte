
<script lang="ts">
    import ChannelsTable from '../components/ChannelsTable.svelte';
    //import { grpc } from '@improbable-eng/grpc-web';
    import { onMount, beforeUpdate} from 'svelte';
    //import {lncrpcClientImpl, GrpcWebImpl} from '../lncrpc/lncrpc'
    // let channels;
    
    // const rpc = new GrpcWebImpl('https://localhost:50051', {
    //     debug: false,
    //     metadata: new grpc.Metadata({}),
    // });
    
    // const client = new lncrpcClientImpl(rpc);
    // let getChannels = async function() {
    //     return client.GetNodeChannels({})
    // }
	
</script>

<svelte:head>
    <title>LNC | Channels</title>
</svelte:head>

<div class="index-page">
    {#await getChannels()}
        Waiting for channels
    {:then channels} 
        <ChannelsTable channels={channels} />    
    {/await}
    
</div>
