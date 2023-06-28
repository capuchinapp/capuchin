<script>
    import Fa from 'svelte-fa';
    import {faRuble} from '@fortawesome/free-solid-svg-icons';
    import {moneyFormatted} from '../services/MoneyHelper';

    /**
     * Целое число копеек
     *
     * @type {number}
     */
    export let value;
</script>

<slot>
    {moneyFormatted(value / 100)}
    <Fa icon={faRuble} />
</slot>
