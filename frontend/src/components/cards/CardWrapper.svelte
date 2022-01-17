<script context="module" lang="ts">
    // This might seem overly complicated, but prepares for drag/drop
    // functionality etc.
    export type CardWrapperProps = {
        type: string;
        heading: string;
    };
</script>

<script lang="ts">
    export let props: CardWrapperProps = {
        type: '',
        heading: '',
    };
</script>

<div class="card-wrapper">
    <div class="heading">
        <div class="value">{props.heading}</div>
        <div class="actions" />
    </div>
    <div class="card {props.type}">
        <slot />
    </div>
</div>

<style>
    .card {
        background-color: white;
        border-radius: 5px;
        margin-bottom: 10px;
    }
    .heading {
        display: flex;
        justify-content: space-between;
        font-size: 13px;
        color: rgba(0, 0, 0, 0.2);
    }
    .heading .value,
    .heading .actions {
        padding: 5px;
    }
    .actions {
        opacity: 0;
    }
    .card-wrapper:hover .actions {
        opacity: 1;
    }
</style>
