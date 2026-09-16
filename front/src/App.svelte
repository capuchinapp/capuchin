<script lang="ts">
    import './app.css';
    import './dayjs';
    import {onMount} from 'svelte';
    import Router, {push} from 'svelte-spa-router';
    import {capuchin, preloaderCount} from './stores';
    import {SvelteToast} from '@zerodevx/svelte-toast';
    import {routes} from './routes';
    import Header from './components/Header.svelte';
    import Footer from './components/Footer.svelte';
    import Upgrade from './components/Upgrade.svelte';
    import ExistRunningTimelog from './components/ExistRunningTimelog.svelte';
    import Settings from './services/Settings';
    import {getDate, nowDate} from './services/DateTime';

    const svelteToastOptions = {
        duration: 5000,
        pausable: true,
    };

    let preloaderContainer = $state<HTMLElement | null>(null);
    let preloaderTimer = $state<number | undefined>();

    let pageLoaded = $state(false);

    $effect(() => {
        if (pageLoaded && $preloaderCount > 0 && typeof preloaderTimer === 'undefined') {
            if (preloaderContainer) {
                preloaderContainer.style.display = 'block';
            }

            preloaderTimer = setInterval(() => {
                if ($preloaderCount <= 0) {
                    if (preloaderTimer !== undefined) {
                        clearInterval(preloaderTimer);
                    }
                    preloaderTimer = undefined;
                    $preloaderCount = 0;

                    if (preloaderContainer) {
                        preloaderContainer.style.display = 'none';
                    }
                }
            }, 100);
        }
    });

    onMount(async () => {
        if ('ontouchstart' in window || navigator.maxTouchPoints > 0) {
            document.body.classList.add('touch-device');
        } else {
            document.body.classList.add('non-touch-device');
        }

        preloaderContainer = document.getElementById('preloader-container');
        pageLoaded = true;

        if (!$capuchin.isAuth) {
            push('/login');

            return;
        }

        await Settings.init();
    });
</script>

<header>
    <div class="container border-warning">
        <Header />
    </div>
</header>

<main class="py-3">
    <div class="container">
        {#if $capuchin.appVersionFront !== $capuchin.appVersionBack}
            <Upgrade />
        {/if}
        {#if $capuchin.runningTimelogDatetime !== null && getDate($capuchin.runningTimelogDatetime) !== nowDate()}
            <ExistRunningTimelog date={getDate($capuchin.runningTimelogDatetime)} />
        {/if}
        <Router {routes} restoreScrollState={true} />
    </div>
</main>

<footer>
    <div class="container border-warning">
        <Footer />
    </div>
</footer>

<div class="toast-container-top-right">
    <SvelteToast target="general" options={svelteToastOptions} />
</div>
<div class="toast-container-top-center">
    <SvelteToast target="errors" options={svelteToastOptions} />
</div>

<div id="preloader-container" class="preloader-container">
    <div class="preloader-overlay"></div>
    <div class="preloader-spinner">
        <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
</div>

<style>
    header .container {
        border-bottom-width: 4px;
        border-bottom-style: solid;
    }

    footer .container {
        border-top-width: 2px;
        border-top-style: solid;
        font-size: 0.9em;
    }

    .container {
        max-width: 960px;
    }

    .preloader-container {
        display: none;
    }

    .preloader-overlay {
        position: fixed;
        z-index: 50000;
        background: #fff;
        opacity: 0.8;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }

    .preloader-spinner {
        z-index: 50001;
        width: 100%;
        height: 100%;
        position: fixed;
        top: 0;
        left: 0;
        display: flex;
        align-items: center;
        align-content: center;
        justify-content: center;
        overflow: auto;
    }

    .toast-container-top-right {
        --toastContainerTop: 1rem;
        --toastContainerRight: 1rem;
    }

    .toast-container-top-center {
        --toastContainerTop: 1rem;
        --toastContainerRight: auto;
        --toastContainerBottom: auto;
        --toastContainerLeft: calc(50vw - 8rem);
    }

    @media (min-width: 440px) {
        .toast-container-top-right {
            --toastWidth: 26rem !important;
        }

        .toast-container-top-center {
            --toastWidth: 26rem !important;
            --toastContainerLeft: calc(50vw - 13rem) !important;
        }
    }
</style>