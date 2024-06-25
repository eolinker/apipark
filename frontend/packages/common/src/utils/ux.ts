/*
 * @Date: 2024-04-29 14:28:07
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-04-29 14:45:47
 * @FilePath: \frontend\packages\core\src\utils\ux.ts
 */
export const withMinimumDelay = <T>(fn: () => Promise<T>, delay: number = 100): Promise<T> => {
    const startTime = Date.now();
    return fn().then(async result => {
      const endTime = Date.now();
      const elapsed = endTime - startTime;
      if (elapsed < delay) {
        await new Promise(resolve => setTimeout(resolve, delay - elapsed));
      }
      return result;
    });
  };