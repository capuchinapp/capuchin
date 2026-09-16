<script lang="ts">
    import Fa from 'svelte-fa';
    import {faTimes, faRubleSign} from '@fortawesome/free-solid-svg-icons';
    import {faClock} from '@fortawesome/free-regular-svg-icons';
    import dayjs from 'dayjs';
    import AirDatepicker from 'air-datepicker';
    import airDatepickerLocaleRu from 'air-datepicker/locale/ru';
    import TomSelect from 'tom-select';
    import Chart from 'chart.js/auto';
    import {onMount} from 'svelte';
    import {t} from '../i18n';
    import PageHeader from '../components/PageHeader.svelte';
    import NothingFound from '../components/NothingFound.svelte';
    import TimelogAPI from '../services/Timelog';
    import ClientAPI from '../services/Clients';
    import ProjectAPI from '../services/Projects';
    import SettingsAPI from '../services/Settings';
    import {moneyFormatted} from '../services/MoneyHelper';
    import {hoursFromSeconds, formatDurationSmart} from '../services/DateTime';
    import {settings} from '../stores';
    import type {Client, Project, Timelog} from '../types';

    type ProjectSummary = {
        title: string;
        clientId: string;
        clientName: string;
        projectId: string;
        projectName: string;
        durationSeconds: number;
        billableAmount: number;
    };

    let calendar!: AirDatepicker<HTMLElement>;
    let selectedDates = $state<dayjs.Dayjs[]>();
    let clientSelect!: TomSelect;
    let projectSelect!: TomSelect;
    let chart!: Chart;

    let workingDays = SettingsAPI.workingDaysFromString($settings.workingDays);

    let pageStarted = $state(false);
    let hasFilter = $state(false);
    let clientSelectedId = $state('');
    let projectSelectedId = $state('');
    let timelogList = $state<Timelog[]>([]);
    let timelogFilteredList = $state<Timelog[]>([]);
    let timelogByProjectList = $state<ProjectSummary[]>([]);
    let clientList: Client[] = [];
    let projectList: Project[] = [];

    $effect(() => {
        if (pageStarted) {
            createChart();
        }
    });

    $effect(() => {
        if (selectedDates) {
            const from = selectedDates[0].format('YYYY-MM-DD');
            const to = selectedDates[1].format('YYYY-MM-DD');

            TimelogAPI.getList(from, to).then((result) => (timelogList = result));
        }
    });

    $effect(() => {
        let tlFilteredList = timelogList;

        hasFilter = false;
        if (clientSelectedId !== '') {
            hasFilter = true;
            tlFilteredList = tlFilteredList.filter((c) => c.clientId === clientSelectedId);
        }
        if (projectSelectedId !== '') {
            hasFilter = true;
            tlFilteredList = tlFilteredList.filter((c) => c.projectId === projectSelectedId);
        }

        timelogFilteredList = tlFilteredList;
    });

    onMount(async () => {
        clientList = await ClientAPI.getList();
        projectList = await ProjectAPI.getList(false);

        clientSelect = new TomSelect('#js-client-select', {
            options: [{id: '', name: `--${$t('clients.all')}--`}, ...clientList],
            items: [''],
            valueField: 'id',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            onChange: (value: string) => (clientSelectedId = String(value)),
            allowEmptyOption: true,
        });

        projectSelect = new TomSelect('#js-project-select', {
            options: [{clientId: '', id: '', name: `--${$t('projects.all')}--`}, ...projectList],
            optgroups: clientList,
            items: [''],
            valueField: 'id',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            optgroupField: 'clientId',
            optgroupValueField: 'id',
            optgroupLabelField: 'name',
            onChange: (value: string) => (projectSelectedId = String(value)),
            allowEmptyOption: true,
            lockOptgroupOrder: true,
        });

        calendar = new AirDatepicker('#js-calendar', {
            locale: airDatepickerLocaleRu,
            dateFormat: getDateFormatForAirDatepicker(),
            autoClose: false,
            isMobile: true,
            range: true,
            multipleDatesSeparator: ' - ',
            buttons: [
                {
                    content: $t('apply'),
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

        const from = dayjs().startOf('month');
        const to = dayjs().endOf('month');

        calendar.selectDate([from.format('YYYY-MM-DD'), to.format('YYYY-MM-DD')], {
            updateTime: true,
            silent: true,
        });

        selectedDates = [from, to];

        pageStarted = true;
    });

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
                return 'dd.MM.yyyy';
        }
    }

    function createChart() {
        console.log('createChart()');

        if (chart) {
            chart.destroy();
        }

        let chartLabels: string[] = [];
        let chartValues: number[] = [];
        let dateKeys: Record<string, number> = {};

        let secondsForAverageValue = 0;
        let workingDates: string[] = [];
        let type = 'months';

        const start = selectedDates![0];
        const end = selectedDates![1];

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

        let tlByProjectList: ProjectSummary[] = [];

        timelogFilteredList.forEach((item) => {
            const key = type === 'days' ? item.date : dayjs(item.date).format('YYYY-MM');

            if (dateKeys.hasOwnProperty(key)) {
                chartValues[dateKeys[key]] += hoursFromSeconds(item.durationSeconds, 2);

                let day = dayjs(item.date).day().toString(10);

                if (workingDays.includes(day)) {
                    secondsForAverageValue += item.durationSeconds;

                    if (!workingDates.includes(item.date)) {
                        workingDates.push(item.date);
                    }
                }

                const logIdx = tlByProjectList.findIndex((e) => e.projectId === item.projectId);
                if (logIdx !== -1) {
                    tlByProjectList[logIdx].durationSeconds += item.durationSeconds;
                    tlByProjectList[logIdx].billableAmount += item.billableAmount;
                } else {
                    tlByProjectList.push({
                        title: `${item.clientName} - ${item.projectName}`,
                        clientId: item.clientId,
                        clientName: item.clientName,
                        projectId: item.projectId,
                        projectName: item.projectName,
                        durationSeconds: item.durationSeconds,
                        billableAmount: item.billableAmount,
                    });
                }
            }
        });

        tlByProjectList.sort((a, b) => (a.title > b.title ? 1 : b.title > a.title ? -1 : 0));

        timelogByProjectList = tlByProjectList;

        const horizontalLinePlugin = {
            id: 'horizontalLine',
            beforeDraw: (chart: Chart) => {
                if (type === 'months') {
                    return;
                }

                const hours = hoursFromSeconds(secondsForAverageValue / workingDates.length, 2);

                const yValue = chart.scales.y.getPixelForValue(hours);
                const ctx = chart.ctx;

                ctx.save();

                // Рисуем линию
                ctx.beginPath();
                ctx.moveTo(chart.chartArea.left, yValue);
                ctx.lineTo(chart.chartArea.right, yValue);
                ctx.strokeStyle = '#ff6421';
                ctx.lineWidth = 2;
                ctx.stroke();

                // Настройки текста
                ctx.fillStyle = '#ff6421';
                ctx.textAlign = 'left';

                // Текст над линией
                const textAbove = hours;
                const textAboveX = chart.chartArea.left + 5; // Отступ слева
                const textAboveY = yValue - 5; // Поднимаем текст немного выше линии
                ctx.font = '11px Arial';
                ctx.fillText(String(textAbove), textAboveX, textAboveY);

                // Текст под линией
                const textBelow = $t('reports.average');
                const textBelowX = chart.chartArea.left + 5; // Отступ слева
                const textBelowY = yValue + 15; // Опускаем текст немного ниже линии
                ctx.font = '10px Arial';
                ctx.fillText(textBelow, textBelowX, textBelowY);

                ctx.restore();
            },
        };

        chart = new Chart(document.getElementById('chart') as HTMLCanvasElement, {
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
                                return `${$t('time')}: ${formatDurationSmart(item.parsed.y * 3600)}`;
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
            plugins: [horizontalLinePlugin],
        });
    }
</script>

<!-- svelte-ignore a11y_label_has_associated_control -->
<PageHeader title={$t('reports.title')} firstColClass="col-auto text-capitalize">
    <div class="col pt2">
        <input type="text" id="js-calendar" class="btn btn-sm btn-primary" />
    </div>
    <div class="col-md pt2">
        <select id="js-client-select" class="form-select form-select-sm"></select>
    </div>
    <div class="col-md pt2">
        <select id="js-project-select" class="form-select form-select-sm"></select>
    </div>
    {#if hasFilter}
        <div class="col-md-auto pt2">
            <button
                type="button"
                title={$t('clearAllFilters')}
                onclick={() => {
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
    <canvas id="chart"></canvas>
</div>

{#if timelogFilteredList.length > 0}
    <table class="table align-middle table-hover mt-3">
        <tbody>
            {#each timelogByProjectList as item (item.title)}
                <tr>
                    <td class="align-middle text-start">
                        {item.clientName} - {item.projectName}
                    </td>
                    <td class="align-middle text-end w1 text-nowrap">
                        <span class="fs-4" title="{$t('hours')}: {hoursFromSeconds(item.durationSeconds, 1)}">
                            {formatDurationSmart(item.durationSeconds)}
                            <Fa icon={faClock} />
                        </span>
                        <br />
                        {moneyFormatted(item.billableAmount / 100)}
                        <Fa icon={faRubleSign} />
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