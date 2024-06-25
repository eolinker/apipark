/*
 * @Date: 2024-02-22 17:52:30
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-02-22 17:52:55
 * @FilePath: \frontend\packages\core\src\utils\validate.ts
 */

export const validateUrlSlash = (_, value) => {
    if (value && value.includes('//')) {
      return Promise.reject(new Error('暂不支持带有双斜杠//的url'));
    }
    return Promise.resolve();
  };