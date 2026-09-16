<script lang="ts">
    import Fa from 'svelte-fa';
    import {faRightToBracket, faUserPlus, faUserCheck} from '@fortawesome/free-solid-svg-icons';
    import {onMount} from 'svelte';
    import {t} from '../i18n';
    import AuthAPI from '../services/Auth';
    import Validation from '../services/Validation';
    import Settings from '../services/Settings';
    import {toastSuccess} from '../services/Toast';

    type LoginValidateErrors = {
        email: string;
    };

    type RegisterValidateErrors = {
        email: string;
    };

    type ApplyValidateErrors = {
        email: string;
        code: string;
    };

    let validation!: Validation;

    let loginValidateErrors = $state<LoginValidateErrors>({
        email: '',
    });

    let registerValidateErrors = $state<RegisterValidateErrors>({
        email: '',
    });

    let applyValidateErrors = $state<ApplyValidateErrors>({
        email: '',
        code: '',
    });

    let loginEmail = $state('');
    let registerEmail = $state('');
    let applyEmail = $state('');
    let applyCode = $state('');

    let applyCodeSent = $state(false);
    let registered = $state(false);

    onMount(async () => {
        validation = new Validation();
    });

    async function onLogin() {
        loginValidateErrors = {
            email: '',
        };

        const data = {
            email: loginEmail,
        };

        validation.validateLogin(
            data,
            async () => {
                AuthAPI.login(data).then(
                    async () => {
                        applyEmail = loginEmail;
                        applyCode = '';
                        applyCodeSent = true;
                    },
                    () => {},
                );
            },
            (err) => {
                loginValidateErrors = err;
            },
        );
    }

    async function onRegister() {
        registerValidateErrors = {
            email: '',
        };

        const data = {
            email: registerEmail,
        };

        validation.validateRegister(
            data,
            async () => {
                AuthAPI.register(data).then(
                    async () => {
                        toastSuccess($t('auth.register.success'), true);
                        registered = true;
                    },
                    () => {},
                );
            },
            (err) => {
                registerValidateErrors = err;
            },
        );
    }

    async function onApplyCode() {
        applyValidateErrors = {
            email: '',
            code: '',
        };

        const data = {
            email: applyEmail,
            code: applyCode,
        };

        validation.validateApplyCode(
            data,
            async () => {
                AuthAPI.applyCode(data).then(
                    async () => {
                        await Settings.init();

                        window.location.assign('/');
                    },
                    () => {},
                );
            },
            (err) => {
                applyValidateErrors = err;
            },
        );
    }
</script>

{#if !applyCodeSent}
    {#if registered}
        <div class="row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header fw-bold">
                        {$t('auth.register.registered.title')}
                    </div>
                    <div class="card-body">
                        {$t('auth.register.registered.message')}
                    </div>
                </div>
            </div>
        </div>
    {:else}
        <div class="row">
            <div class="col-md-6">
                <ul class="nav nav-tabs" id="authTab" role="tablist">
                    <li class="nav-item" role="presentation">
                        <button
                            class="nav-link active"
                            id="login-tab"
                            data-bs-toggle="tab"
                            data-bs-target="#login-tab-pane"
                            type="button"
                            role="tab"
                        >
                            {$t('auth.login.title')}
                        </button>
                    </li>
                    <li class="nav-item" role="presentation">
                        <button
                            class="nav-link"
                            id="register-tab"
                            data-bs-toggle="tab"
                            data-bs-target="#register-tab-pane"
                            type="button"
                            role="tab"
                        >
                            {$t('auth.register.title')}
                        </button>
                    </li>
                </ul>
                <div class="tab-content border border-top-0" id="authTabContent">
                    <div class="tab-pane fade show active" id="login-tab-pane" role="tabpanel" tabindex="0">
                        <!-- svelte-ignore a11y_label_has_associated_control -->
                        <form
                            onsubmit={(e) => {
                                e.preventDefault();
                                onLogin();
                            }}
                            class="p-3"
                        >
                            <div class="mb-3">
                                <label class="form-label">{$t('auth.login.email')}</label>
                                <input
                                    type="email"
                                    bind:value={loginEmail}
                                    class="form-control"
                                    class:is-invalid={loginValidateErrors.email}
                                    autocomplete="off"
                                />
                                {#if loginValidateErrors.email}
                                    <div class="invalid-feedback">{loginValidateErrors.email}</div>
                                {/if}
                            </div>
                            <div class="mb-3">
                                <button type="button" onclick={onLogin} class="btn btn-success">
                                    <Fa icon={faRightToBracket} />
                                    {$t('auth.login.button')}
                                </button>
                            </div>
                        </form>
                    </div>
                    <div class="tab-pane fade" id="register-tab-pane" role="tabpanel" tabindex="0">
                        <!-- svelte-ignore a11y_label_has_associated_control -->
                        <form
                            onsubmit={(e) => {
                                e.preventDefault();
                                onRegister();
                            }}
                            class="p-3"
                        >
                            <div class="mb-3">
                                <label class="form-label">{$t('auth.register.email')}</label>
                                <input
                                    type="email"
                                    bind:value={registerEmail}
                                    class="form-control"
                                    class:is-invalid={registerValidateErrors.email}
                                    autocomplete="off"
                                />
                                {#if registerValidateErrors.email}
                                    <div class="invalid-feedback">{registerValidateErrors.email}</div>
                                {/if}
                            </div>
                            <div class="mb-0">
                                <button type="button" onclick={onRegister} class="btn btn-success">
                                    <Fa icon={faUserPlus} />
                                    {$t('auth.register.button')}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    {/if}
{/if}
{#if applyCodeSent}
    <div class="row">
        <div class="col-md-6">
            <!-- svelte-ignore a11y_label_has_associated_control -->
            <form
                onsubmit={(e) => {
                    e.preventDefault();
                    onApplyCode();
                }}
                class="p-3"
            >
                <div class="card">
                    <div class="card-header fw-bold">
                        {$t('auth.applycode.title')}
                    </div>
                    <div class="card-body">
                        <div class="mb-3 d-none">
                            <label class="form-label">{$t('auth.applycode.email')}</label>
                            <input
                                type="email"
                                bind:value={applyEmail}
                                class="form-control"
                                class:is-invalid={applyValidateErrors.email}
                                autocomplete="off"
                            />
                            {#if applyValidateErrors.email}
                                <div class="invalid-feedback">{applyValidateErrors.email}</div>
                            {/if}
                        </div>
                        <div class="mb-3">
                            <label class="form-label">{$t('auth.applycode.code')}</label>
                            <input
                                type="text"
                                bind:value={applyCode}
                                class="form-control"
                                class:is-invalid={applyValidateErrors.code}
                                autocomplete="off"
                            />
                            {#if applyValidateErrors.code}
                                <div class="invalid-feedback">{applyValidateErrors.code}</div>
                            {/if}
                        </div>
                        <button type="button" onclick={onApplyCode} class="btn btn-success">
                            <Fa icon={faUserCheck} />
                            {$t('auth.applycode.button')}
                        </button>
                    </div>
                </div>
            </form>
        </div>
    </div>
{/if}