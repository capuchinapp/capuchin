<script lang="ts">
    import {t} from '../i18n';
    import AuthAPI from '../services/Auth';
    import Settings from '../services/Settings';

    type Props = {
        params?: Record<string, string>;
    };

    let {params = {}}: Props = $props();

    const payload = {
        // svelte-ignore state_referenced_locally
        userId: params.uid,
        // svelte-ignore state_referenced_locally
        code: params.code,
    };

    AuthAPI.activate(payload).then(
        async () => {
            await Settings.init();

            window.location.assign('/');
        },
        () => {},
    );
</script>

<div class="row justify-content-center">
    <div class="col-md-6">
        <div class="card">
            <div class="card-header fw-bold">
                {$t('auth.activate.title')}
            </div>
            <div class="card-body">
                {$t('auth.activate.waiting')}
            </div>
        </div>
    </div>
</div>