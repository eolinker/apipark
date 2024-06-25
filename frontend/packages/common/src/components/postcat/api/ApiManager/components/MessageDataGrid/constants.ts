/*
 * @Date: 2024-06-04 08:54:16
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 10:07:26
 * @FilePath: \frontend\packages\common\src\components\postcat\api\ApiManager\components\MessageDataGrid\constants.ts
 */
import {generateId} from "@common/utils/postcat.tsx";

type SafeAny = unknown
export function generateRow(data: SafeAny = {}) {
  return Object.assign({
    id: generateId(),
    name: '',
    dataType: null,
    isRequired: 1,
    description: '',
    paramAttr: {
      example: ''
    },
    childList: []
  }, data)
}