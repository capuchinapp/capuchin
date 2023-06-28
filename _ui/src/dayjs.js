import 'dayjs/locale/en';
import 'dayjs/locale/ru';

import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import localizedFormat from 'dayjs/plugin/localizedFormat';

import {get} from 'svelte/store';
import {currentLocale} from './stores';

dayjs.extend(duration);
dayjs.extend(localizedFormat);

dayjs.locale(get(currentLocale));
