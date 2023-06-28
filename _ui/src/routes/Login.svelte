<script>
    import Fa from 'svelte-fa';
    import {faRightToBracket} from '@fortawesome/free-solid-svg-icons';
    import {onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import AuthAPI from '../services/Auth';
    import Validation from '../services/Validation';
    import Settings from '../services/Settings';

    let /** @type {Validation} */ validation;

    let validateErrors = {};
    let password = '';

    onMount(async () => {
        validation = new Validation();
    });

    async function onLogin() {
        validateErrors = {};

        const data = {
            password: password,
        };

        validation.validateLogin(
            data,
            async () => {
                AuthAPI.login(data).then(
                    async () => {
                        await Settings.init();

                        window.location.assign('/');
                    },
                    () => {}
                );
            },
            (err) => {
                validateErrors = err;
            }
        );
    }
</script>

<div class="row justify-content-center mb-3">
    <div class="col-4">
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <form on:submit|preventDefault={onLogin} class="p-3">
            <div class="card">
                <div class="card-header fw-bold">
                    {$_('auth.login.title')}
                </div>
                <div class="card-body">
                    <div class="mb-3">
                        <label class="form-label">{$_('auth.login.password')}</label>
                        <input
                            type="password"
                            bind:value={password}
                            class="form-control"
                            class:is-invalid={validateErrors.password}
                            autocomplete="off"
                        />
                        {#if validateErrors.password}
                            <div class="invalid-feedback">{validateErrors.password}</div>
                        {/if}
                    </div>
                    <button type="button" on:click={onLogin} class="btn btn-success">
                        <Fa icon={faRightToBracket} />
                        {$_('auth.login.button')}
                    </button>
                </div>
            </div>
        </form>
    </div>
</div>
