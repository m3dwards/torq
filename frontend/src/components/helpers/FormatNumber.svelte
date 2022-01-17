<script lang="ts">
    export let value: number | undefined = undefined;
    export let lang = 'en';
    export let decimals: number | undefined;
    export let notation: 'compact' | 'standard' | undefined = undefined;
    export let style: 'percent' | undefined = undefined;

    // Adding sensible defualt when notation is not set
    if (notation === undefined && value >= 1000000) {
        notation = notation ? notation : 'compact';
    }

    // Adding sensible default when decimals are not set for compact format
    if (notation === 'compact' && decimals === undefined) {
        decimals = 1;
    }

    function format(value: number, language: string, decimals: number): string {
        return value.toLocaleString(language, {
            notation: notation,
            style: style,
            compactDisplay: 'short',
            minimumFractionDigits: 0,
            maximumFractionDigits: decimals,
        });
    }
</script>

{value ? format(value, lang, decimals) : '0'}
