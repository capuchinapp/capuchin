<script>
    import Fa from 'svelte-fa';
    import {
        faStopwatch,
        faBriefcase,
        faUsers,
        faChartBar,
        faRightFromBracket,
        faUser,
    } from '@fortawesome/free-solid-svg-icons';
    import capuchinLogo from '../assets/logo.png';
    import {link} from 'svelte-spa-router';
    import active from 'svelte-spa-router/active';
    import {_} from 'svelte-i18n';
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
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
        <span class="navbar-toggler-icon" />
    </button>
    <div class="collapse navbar-collapse" id="navbarNav">
        <ul class="navbar-nav me-auto">
            <li class="nav-item">
                <a class="nav-link active" href="/" use:link use:active>
                    <Fa icon={faStopwatch} />
                    {$_('track.title')}
                </a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="/projects" use:link use:active>
                    <Fa icon={faBriefcase} />
                    {$_('projects.title')}
                </a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="/clients" use:link use:active>
                    <Fa icon={faUsers} />
                    {$_('clients.title')}
                </a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="/reports" use:link use:active>
                    <Fa icon={faChartBar} />
                    {$_('reports.title')}
                </a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="/profile" use:link use:active>
                    <Fa icon={faUser} />
                    {$_('profile.title')}
                </a>
            </li>
        </ul>
        {#if $timeTotalSeconds > 0}
            <span class="navbar-text me-3">
                <span class="fs-3 {$timeTotalClass}">{calculateTime($timeTotalSeconds, true)}</span>
            </span>
        {/if}
        {#if $capuchin.authMode && $capuchin.isAuth}
            <button type="button" on:click={onLogout} class="btn btn-link nav-link">
                <Fa icon={faRightFromBracket} />
                {$_('auth.logout.button')}
            </button>
        {/if}
    </div>
</nav>

<style>
    .logo img {
        height: 2.18em;
    }

    .navbar-nav .nav-link.active {
        color: #e95420;
    }
</style>
