/*
 * @Date: 2024-02-19 18:11:15
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-29 18:38:11
 * @FilePath: \frontend\packages\common\src\components\aoplatform\DatePicker.tsx
 */
import { DatePicker } from 'antd';
import type { Moment } from 'moment';
import momentGenerateConfig from 'rc-picker/lib/generate/moment';

const MyDatePicker = DatePicker.generatePicker<Moment>(momentGenerateConfig);

export default MyDatePicker;