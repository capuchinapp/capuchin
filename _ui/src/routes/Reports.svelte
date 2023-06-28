<script>
    import Fa from 'svelte-fa';
    import {faChevronLeft, faChevronRight, faTimes, faMoneyBillAlt} from '@fortawesome/free-solid-svg-icons';
    import {faClock} from '@fortawesome/free-regular-svg-icons';
    import dayjs from 'dayjs';
    import AirDatepicker from 'air-datepicker';
    import airDatepickerLocaleEn from 'air-datepicker/locale/en';
    import airDatepickerLocaleRu from 'air-datepicker/locale/ru';
    import TomSelect from 'tom-select';
    import Chart from 'chart.js/auto';
    import {onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import PageHeader from '../components/PageHeader.svelte';
    import NothingFound from '../components/NothingFound.svelte';
    import TimelogAPI from '../services/Timelog';
    import ClientAPI from '../services/Clients';
    import ProjectAPI from '../services/Projects';
    import {moneyFormatted} from '../services/MoneyHelper';
    import {hoursFromSeconds} from '../services/DateTime';
    import {currentLocale, settings} from '../stores';

    let /** @type {AirDatepicker<HTMLElement>} */ calendar;
    let /** @type {Array.<dayjs.Dayjs>} */ selectedDates;
    let /** @type {TomSelect} */ clientSelect;
    let /** @type {TomSelect} */ projectSelect;
    let /** @type {Chart} */ chart;

    let clientSelectedUUID = '';
    let projectSelectedUUID = '';
    let timelogList = [];
    let timelogFilteredList = [];
    let timelogByProjectList = [];
    let clientList = [];
    let projectList = [];
    let hasFilter = false;

    $: if (selectedDates) {
        const from = selectedDates[0].format('YYYY-MM-DD');
        const to = selectedDates[1].format('YYYY-MM-DD');

        calendar.selectDate([from, to], {
            updateTime: true,
            silent: true,
        });

        TimelogAPI.getList(from, to).then((result) => (timelogList = result));
    }

    $: {
        timelogFilteredList = timelogList;
        hasFilter = false;

        if (clientSelectedUUID !== '') {
            hasFilter = true;
            timelogFilteredList = timelogFilteredList.filter((c) => c.clientUUID === clientSelectedUUID);
        }

        if (projectSelectedUUID !== '') {
            hasFilter = true;
            timelogFilteredList = timelogFilteredList.filter((c) => c.projectUUID === projectSelectedUUID);
        }

        createChart();
    }

    onMount(async () => {
        clientList = await ClientAPI.getList();
        projectList = await ProjectAPI.getList(false);

        clientSelect = new TomSelect('#js-client-select', {
            options: [{uuid: '', name: `--${$_('clients.all')}--`}, ...clientList],
            items: [''],
            valueField: 'uuid',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            onChange: (value) => (clientSelectedUUID = String(value)),
            allowEmptyOption: true,
        });

        projectSelect = new TomSelect('#js-project-select', {
            options: [{clientUUID: '', uuid: '', name: `--${$_('projects.all')}--`}, ...projectList],
            optgroups: clientList,
            items: [''],
            valueField: 'uuid',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            optgroupField: 'clientUUID',
            optgroupValueField: 'uuid',
            optgroupLabelField: 'name',
            onChange: (value) => (projectSelectedUUID = String(value)),
            allowEmptyOption: true,
            lockOptgroupOrder: true,
        });

        calendar = new AirDatepicker('#js-calendar', {
            locale: $currentLocale === 'ru' ? airDatepickerLocaleRu : airDatepickerLocaleEn,
            dateFormat: getDateFormatForAirDatepicker(),
            autoClose: false,
            isMobile: true,
            range: true,
            multipleDatesSeparator: ' - ',
            buttons: [
                {
                    content: $_('apply'),
                    onClick(dp) {
                        if (dp.selectedDates.length === 2) {
                            selectedDates = [
                                dayjs(dp.selectedDates[0].toString()),
                                dayjs(dp.selectedDates[1].toString()),
                            ];

                            dp.hide();
                        }
                    },
                },
            ],
        });

        selectedDates = [dayjs().startOf('month'), dayjs().endOf('month')];
    });

    /**
     * Создание графика
     */
    function getDateFormatForAirDatepicker() {
        switch ($settings['dateFormat']) {
            case 'MM/dd/yyyy':
                return '';

            case 'dd/MM/yyyy':
                return '';

            case 'yyyy-MM-dd':
                return '';

            case 'dd.MM.yyyy':
                return '';

            case 'dd-MM-yyyy':
                return '';

            case 'MM-dd-yyyy':
                return '';

            default:
                return $currentLocale === 'ru' ? 'dd.MM.yyyy' : 'MM/dd/yyyy';
        }
    }

    /**
     * Создание графика
     */
    function createChart() {
        if (chart) {
            chart.destroy();
        }

        if (selectedDates) {
            let /** @type {Array.<string>} */ chartLabels = [];
            let /** @type {Array.<number>} */ chartValues = [];
            let /** @type {Object.<string, number>} */ dateKeys = {};

            let type = 'months';

            const start = selectedDates[0];
            const end = selectedDates[1];

            if (start.month() === end.month()) {
                // Выбранный период внутри одного месяца: показываем график по дням
                type = 'days';

                let currDate = start;

                const diff = end.diff(start, 'day');

                for (let i = 0; i <= diff; i++) {
                    chartLabels.push(currDate.format('dd, DD'));
                    chartValues.push(0);
                    dateKeys[currDate.format('YYYY-MM-DD')] = i;

                    currDate = currDate.add(1, 'day');
                }
            } else {
                // Выбранный период охватывает более одного месяца: показываем график по месяцам
                type = 'months';

                let currDate = start.startOf('month');

                const diff = end.diff(start, 'month');

                for (let i = 0; i <= diff; i++) {
                    chartLabels.push(currDate.format('MMM, YYYY'));
                    chartValues.push(0);
                    dateKeys[currDate.format('YYYY-MM')] = i;

                    currDate = currDate.add(1, 'month');
                }
            }

            timelogByProjectList = [];

            timelogFilteredList.forEach((item) => {
                const key = type === 'days' ? item.date : dayjs(item.date).format('YYYY-MM');

                if (dateKeys.hasOwnProperty(key)) {
                    chartValues[dateKeys[key]] += hoursFromSeconds(item.durationSeconds, 2);

                    const logIdx = timelogByProjectList.findIndex((e) => e.projectUUID === item.projectUUID);
                    if (logIdx !== -1) {
                        timelogByProjectList[logIdx].durationSeconds += item.durationSeconds;
                        timelogByProjectList[logIdx].billableAmount += item.billableAmount;
                    } else {
                        timelogByProjectList.push({
                            title: `${item.clientName} - ${item.projectName}`,
                            clientUUID: item.clientUUID,
                            clientName: item.clientName,
                            projectUUID: item.projectUUID,
                            projectName: item.projectName,
                            durationSeconds: item.durationSeconds,
                            billableAmount: item.billableAmount,
                        });
                    }
                }
            });

            timelogFilteredList.sort((a, b) => (a.title > b.title ? 1 : b.title > a.title ? -1 : 0));

            // @ts-ignore
            chart = new Chart(document.getElementById('chart'), {
                type: 'bar',
                data: {
                    labels: chartLabels,
                    datasets: [
                        {
                            data: chartValues,
                            backgroundColor: 'rgba(56, 180, 74, 0.7)',
                        },
                    ],
                },
                options: {
                    animation: {
                        duration: 0,
                    },
                    plugins: {
                        legend: {display: false},
                        tooltip: {
                            displayColors: false,
                            callbacks: {
                                label: (item) => {
                                    return `${$_('hours')}: ${item.parsed.y}`;
                                },
                            },
                        },
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                        },
                    },
                },
            });
        }
    }
</script>

<!-- svelte-ignore a11y-label-has-associated-control -->
<PageHeader title={$_('reports.title')} firstColClass="col-auto text-capitalize">
    <div class="col pt2">
        <input type="text" id="js-calendar" class="btn btn-sm btn-primary" />
    </div>
    <div class="col pt2">
        <select id="js-client-select" class="form-select form-select-sm" />
    </div>
    <div class="col pt2">
        <select id="js-project-select" class="form-select form-select-sm" />
    </div>
    {#if hasFilter}
        <div class="col-auto pt2">
            <button
                type="button"
                title={$_('clearAllFilters')}
                on:click={() => {
                    clientSelect.setValue('', false);
                    projectSelect.setValue('', false);
                }}
                class="btn btn-sm btn-primary"
            >
                <Fa fw icon={faTimes} />
            </button>
        </div>
    {/if}
</PageHeader>

<div class="mt-3" class:d-none={timelogFilteredList.length === 0}>
    <canvas id="chart" />
</div>

{#if timelogFilteredList.length > 0}
    <table class="table table-hover mt-3">
        <tbody>
            {#each timelogByProjectList as item (item.title)}
                <tr>
                    <td class="align-middle text-start">
                        {item.clientName} - {item.projectName}
                    </td>
                    <td class="align-middle text-end w1 text-nowrap">
                        <span class="fs-4" title={$_('hours')}>
                            {hoursFromSeconds(item.durationSeconds, 2)}
                            <Fa icon={faClock} />
                        </span>
                        <br />
                        {moneyFormatted(item.billableAmount / 100)}
                        <Fa icon={faMoneyBillAlt} />
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <NothingFound />
{/if}

<style>
    #chart {
        height: 300px;
    }
</style>
