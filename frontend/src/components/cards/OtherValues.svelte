<script lang="ts">
    import CardWrapper from './CardWrapper.svelte';
</script>

<CardWrapper props={{ type: 'other-card', heading: 'Other' }}>
    <slot />
</CardWrapper>
