<script>
    import './app.css';
    import './dayjs';
    import './i18n';
    import {onMount} from 'svelte';
    import {push} from 'svelte-spa-router';
    import {capuchin, preloaderCount} from './stores';
    import {SvelteToast} from '@zerodevx/svelte-toast';
    import Router from 'svelte-spa-router';
    import {routes} from './routes';
    import Header from './components/Header.svelte';
    import Footer from './components/Footer.svelte';
    import Settings from './services/Settings';

    const svelteToastOptions = {
        duration: 5000,
        pausable: true,
    };

    let /** @type {HTMLElement} */ preloaderContainer;
    let /** @type {number} */ preloaderTimer;

    $: if ($preloaderCount > 0 && typeof preloaderTimer === 'undefined') {
        if (preloaderContainer) {
            preloaderContainer.style.display = '';
        }

        preloaderTimer = setInterval(() => {
            if ($preloaderCount <= 0) {
                clearInterval(preloaderTimer);
                preloaderTimer = undefined;
                $preloaderCount = 0;

                if (preloaderContainer) {
                    preloaderContainer.style.display = 'none';
                }
            }
        }, 100);
    }

    onMount(async () => {
        preloaderContainer = document.getElementById('preloader-container');

        if ($capuchin.authMode && !$capuchin.isAuth) {
            push('/login');

            return;
        }

        await Settings.init();
    });
</script>

<header class="container border-warning">
    <Header />
</header>

<main class="container">
    <Router {routes} restoreScrollState={true} />
</main>

<footer class="container border-warning">
    <Footer />
</footer>

<SvelteToast options={svelteToastOptions} />

<div id="preloader-container" class="preloader-container">
    <div class="preloader-overlay" />
    <div class="preloader-spinner">
        <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
</div>

<style>
    header {
        border-bottom-width: 4px;
        border-bottom-style: solid;
    }

    footer {
        border-top-width: 2px;
        border-top-style: solid;
        padding: 0.5em 1em;
        font-size: 0.9em;
        text-align: center;
    }

    main {
        padding: 1em 0 0 0;
    }

    .container {
        max-width: 960px;
    }

    .preloader-container {
        display: none;
    }

    .preloader-overlay {
        position: fixed;
        z-index: 500;
        background: #fff;
        opacity: 0.8;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }

    .preloader-spinner {
        z-index: 501;
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
</style>
