<script lang="ts">
    import Fa from 'svelte-fa';
    import {
        faStopwatch,
        faBriefcase,
        faUsers,
        faChartBar,
        faRightFromBracket,
        faUser,
        faIdCard,
    } from '@fortawesome/free-solid-svg-icons';
    import capuchinLogo from '../assets/logo.png';
    import {link} from 'svelte-spa-router';
    import active from 'svelte-spa-router/active';
    import {t} from '../i18n';
    import AuthAPI from '../services/Auth';
    import {calculateTime} from '../services/DateTime';
    import {timeTotalSeconds, timeTotalClass, capuchin} from '../stores';

    function onLogout() {
        AuthAPI.logout().then(() => {
            window.location.reload();
        });
    }
</script>

<nav class="navbar navbar-expand-lg navbar-light">
    <a class="navbar-brand logo" href="/" use:link>
        <img src={capuchinLogo} title="Capuchin" alt="Capuchin" />
    </a>
    <!-- svelte-ignore a11y_consider_explicit_label -->
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
        <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNav">
        {#if $capuchin.isAuth}
            <ul class="navbar-nav me-auto">
                <li class="nav-item">
                    <a class="nav-link active" href="/" use:link use:active>
                        <Fa fw icon={faStopwatch} />
                        {$t('track.title')}
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="/projects" use:link use:active>
                        <Fa fw icon={faBriefcase} />
                        {$t('projects.title')}
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="/clients" use:link use:active>
                        <Fa fw icon={faUsers} />
                        {$t('clients.title')}
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="/reports" use:link use:active>
                        <Fa fw icon={faChartBar} />
                        {$t('reports.title')}
                    </a>
                </li>
            </ul>
            {#if $timeTotalSeconds > 0}
                <span class="navbar-text me-3">
                    <span class="fs-3 {$timeTotalClass}">{calculateTime($timeTotalSeconds, true)}</span>
                </span>
            {/if}
            <div class="dropdown">
                <button class="btn btn-outline-dark btn-sm dropdown-toggle" type="button" data-bs-toggle="dropdown">
                    <Fa icon={faUser} />
                </button>
                <ul class="dropdown-menu">
                    <li>
                        <a class="dropdown-item" href="/profile" use:link use:active>
                            <Fa fw icon={faIdCard} />
                            {$t('profile.title')}
                        </a>
                    </li>
                    <li>
                        <button type="button" onclick={onLogout} class="dropdown-item">
                            <Fa fw icon={faRightFromBracket} />
                            {$t('auth.logout.button')}
                        </button>
                    </li>
                </ul>
            </div>
        {/if}
    </div>
</nav>

<style>
    .logo img {
        height: 2.18em;
    }
</style>