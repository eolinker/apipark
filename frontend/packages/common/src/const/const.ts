/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-02-27 16:17:37
 * @FilePath: \frontend\packages\core\src\const\const.ts
 */
export type BasicResponse<T> = {
    code:number
    data:T
    msg:string
}


export const STATUS_CODE = {
    SUCCESS:0,
    UNANTHORIZED:401,
    FORBIDDEN:403
}

export const STATUS_COLOR = {
    'done':'text-[#03a9f4]',
    'error':'text-[#ff3b30]'
}
