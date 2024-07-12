// vite.config.ts
import { defineConfig } from "file:///Users/liujian/work/golang/src/github.com/eolinker/apipark/frontend/node_modules/.pnpm/vite@5.0.12_@types+node@20.11.10_less@4.2.0/node_modules/vite/dist/node/index.js";
import react from "file:///Users/liujian/work/golang/src/github.com/eolinker/apipark/frontend/node_modules/.pnpm/@vitejs+plugin-react@4.2.1_vite@5.0.12/node_modules/@vitejs/plugin-react/dist/index.mjs";
import path from "path";
import dynamicImportVars from "file:///Users/liujian/work/golang/src/github.com/eolinker/apipark/frontend/node_modules/.pnpm/@rollup+plugin-dynamic-import-vars@2.1.2/node_modules/@rollup/plugin-dynamic-import-vars/dist/es/index.js";
import tailwindcss from "file:///Users/liujian/work/golang/src/github.com/eolinker/apipark/frontend/node_modules/.pnpm/tailwindcss@3.4.1/node_modules/tailwindcss/lib/index.js";
import autoprefixer from "file:///Users/liujian/work/golang/src/github.com/eolinker/apipark/frontend/node_modules/.pnpm/autoprefixer@10.4.17_postcss@8.4.33/node_modules/autoprefixer/lib/autoprefixer.js";
var __vite_injected_original_dirname = "/Users/liujian/work/golang/src/github.com/eolinker/apipark/frontend/packages/core";
var vite_config_default = defineConfig({
  cacheDir: "./node_modules/.vite",
  build: {
    outDir: "../../dist",
    sourcemap: false,
    chunkSizeWarningLimit: 5e4,
    cacheDir: "./node_modules/.vite",
    output: {
      manualChunks(id) {
        if (id.includes("node_modules")) {
          return id.toString().split("node_modules/")[1].split("/")[0].toString();
        }
        if (id.includes(".pnpm")) {
          const segments = id.split(path.sep);
          const packageName = segments[segments.indexOf(".pnpm") + 1].split("@")[0];
          return packageName;
        }
      }
    }
  },
  css: {
    postcss: {
      plugins: [
        tailwindcss(path.resolve(__vite_injected_original_dirname, "../common/tailwind.config.js")),
        autoprefixer
      ]
    },
    preprocessorOptions: {
      less: {
        javascriptEnabled: true
      }
    },
    modules: {
      localsConvention: "camelCase",
      generateScopedName: "[local]_[hash:base64:2]"
    }
  },
  plugins: [
    react(),
    dynamicImportVars({
      include: ["src"],
      exclude: [],
      warnOnError: false
    })
    //  MillionLint.vite()
    //  visualizer({
    //    filename: 'stats.html', // 生成的可视化报告文件名
    //    sourcemap: true, // 使用 sourcemap，以便更准确地显示原始源文件
    //   //  gzip: true, // 显示 gzip 压缩后的大小
    //    brotliSize: false, // 显示 brotli 压缩后的大小
    //    open: true, // 自动生成报告后自动打开浏览器
    //    // 其他可视化选项...
    //  }),
  ],
  resolve: {
    alias: [
      { find: /^~/, replacement: "" },
      { find: "@common", replacement: path.resolve(__vite_injected_original_dirname, "../common/src") },
      { find: "@market", replacement: path.resolve(__vite_injected_original_dirname, "../market/src") },
      { find: "@core", replacement: path.resolve(__vite_injected_original_dirname, "./src") }
    ]
  },
  server: {
    proxy: {
      "/api/v1": {
        // target: 'http://uat.apikit.com:11204/mockApi/aoplatform/',
        target: "http://172.18.166.219:8288/",
        changeOrigin: true
      },
      "/api2/v1": {
        // target: 'http://uat.apikit.com:11204/mockApi/aoplatform/',
        target: "http://172.18.166.219:8288/",
        changeOrigin: true
      }
    }
  },
  logLevel: "info"
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvVXNlcnMvbGl1amlhbi93b3JrL2dvbGFuZy9zcmMvZ2l0aHViLmNvbS9lb2xpbmtlci9hcGlwYXJrL2Zyb250ZW5kL3BhY2thZ2VzL2NvcmVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIi9Vc2Vycy9saXVqaWFuL3dvcmsvZ29sYW5nL3NyYy9naXRodWIuY29tL2VvbGlua2VyL2FwaXBhcmsvZnJvbnRlbmQvcGFja2FnZXMvY29yZS92aXRlLmNvbmZpZy50c1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vVXNlcnMvbGl1amlhbi93b3JrL2dvbGFuZy9zcmMvZ2l0aHViLmNvbS9lb2xpbmtlci9hcGlwYXJrL2Zyb250ZW5kL3BhY2thZ2VzL2NvcmUvdml0ZS5jb25maWcudHNcIjsvKlxuICogQERhdGU6IDIwMjQtMDEtMzEgMTU6MDA6MzlcbiAqIEBMYXN0RWRpdG9yczogbWFnZ2lleXl5XG4gKiBATGFzdEVkaXRUaW1lOiAyMDI0LTA2LTA3IDE3OjA0OjQxXG4gKiBARmlsZVBhdGg6IFxcZnJvbnRlbmRcXHBhY2thZ2VzXFxjb3JlXFx2aXRlLmNvbmZpZy50c1xuICovXG5pbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tICd2aXRlJ1xuaW1wb3J0IHJlYWN0IGZyb20gJ0B2aXRlanMvcGx1Z2luLXJlYWN0J1xuaW1wb3J0IHBhdGggZnJvbSAncGF0aCdcbmltcG9ydCBkeW5hbWljSW1wb3J0VmFycyBmcm9tICdAcm9sbHVwL3BsdWdpbi1keW5hbWljLWltcG9ydC12YXJzJztcbmltcG9ydCB0YWlsd2luZGNzcyBmcm9tICd0YWlsd2luZGNzcyc7XG5pbXBvcnQgYXV0b3ByZWZpeGVyIGZyb20gJ2F1dG9wcmVmaXhlcic7XG4vLyBpbXBvcnQge3Zpc3VhbGl6ZXJ9IGZyb20gJ3JvbGx1cC1wbHVnaW4tdmlzdWFsaXplcic7XG4vLyBpbXBvcnQgTWlsbGlvbkxpbnQgZnJvbSAnQG1pbGxpb24vbGludCdcblxuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKHtcbiAgY2FjaGVEaXI6ICcuL25vZGVfbW9kdWxlcy8udml0ZScsXG4gIGJ1aWxkOntcbiAgICBvdXREaXI6Jy4uLy4uL2Rpc3QnLFxuICAgIHNvdXJjZW1hcDogZmFsc2UsXG4gICAgY2h1bmtTaXplV2FybmluZ0xpbWl0OiA1MDAwMCxcbiAgICBjYWNoZURpcjogJy4vbm9kZV9tb2R1bGVzLy52aXRlJywgXG4gICAgICBvdXRwdXQ6IHtcbiAgICAgICAgbWFudWFsQ2h1bmtzKGlkKSB7XG4gICAgICAgICAgaWYgKGlkLmluY2x1ZGVzKCdub2RlX21vZHVsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuIGlkLnRvU3RyaW5nKCkuc3BsaXQoJ25vZGVfbW9kdWxlcy8nKVsxXS5zcGxpdCgnLycpWzBdLnRvU3RyaW5nKCk7XG4gICAgICAgICAgfVxuICAgICAgICAgIC8vIFx1OTQ4OFx1NUJGOSBwbnBtIFx1NTQ4QyBNb25vcmVwbyBcdTcyNzlcdTZCOEFcdTU5MDRcdTc0MDZcbiAgICAgICAgICBpZiAoaWQuaW5jbHVkZXMoJy5wbnBtJykpIHtcbiAgICAgICAgICAgIGNvbnN0IHNlZ21lbnRzID0gaWQuc3BsaXQocGF0aC5zZXApO1xuICAgICAgICAgICAgY29uc3QgcGFja2FnZU5hbWUgPSBzZWdtZW50c1tzZWdtZW50cy5pbmRleE9mKCcucG5wbScpICsgMV0uc3BsaXQoJ0AnKVswXTtcbiAgICAgICAgICAgIHJldHVybiBwYWNrYWdlTmFtZTtcbiAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICAgIH0sXG4gICAgfSxcbiAgY3NzOiB7XG4gICAgcG9zdGNzczoge1xuICAgICAgcGx1Z2luczogW1xuICAgICAgICB0YWlsd2luZGNzcyhwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnLi4vY29tbW9uL3RhaWx3aW5kLmNvbmZpZy5qcycpKSwgXG4gICAgICAgIGF1dG9wcmVmaXhlclxuICAgICAgXSxcbiAgICB9LFxuICAgIHByZXByb2Nlc3Nvck9wdGlvbnM6IHtcbiAgICAgIGxlc3M6IHtcbiAgICAgICAgamF2YXNjcmlwdEVuYWJsZWQ6IHRydWUsXG4gICAgICB9LFxuICAgIH0sXG4gICAgbW9kdWxlczp7XG4gICAgICBsb2NhbHNDb252ZW50aW9uOlwiY2FtZWxDYXNlXCIsXG4gICAgICBnZW5lcmF0ZVNjb3BlZE5hbWU6XCJbbG9jYWxdX1toYXNoOmJhc2U2NDoyXVwiXG4gICAgfVxuICB9LFxuICBwbHVnaW5zOiBbcmVhY3QoKSxcbiAgICAgIGR5bmFtaWNJbXBvcnRWYXJzKHtcbiAgICAgICAgaW5jbHVkZTpbXCJzcmNcIl0sXG4gICAgICAgIGV4Y2x1ZGU6W10sXG4gICAgICAgIHdhcm5PbkVycm9yOmZhbHNlXG4gICAgICAgfSksXG4gICAgICAvLyAgTWlsbGlvbkxpbnQudml0ZSgpXG4gICAgICAvLyAgdmlzdWFsaXplcih7XG4gICAgICAvLyAgICBmaWxlbmFtZTogJ3N0YXRzLmh0bWwnLCAvLyBcdTc1MUZcdTYyMTBcdTc2ODRcdTUzRUZcdTg5QzZcdTUzMTZcdTYyQTVcdTU0NEFcdTY1ODdcdTRFRjZcdTU0MERcbiAgICAgIC8vICAgIHNvdXJjZW1hcDogdHJ1ZSwgLy8gXHU0RjdGXHU3NTI4IHNvdXJjZW1hcFx1RkYwQ1x1NEVFNVx1NEZCRlx1NjZGNFx1NTFDNlx1Nzg2RVx1NTczMFx1NjYzRVx1NzkzQVx1NTM5Rlx1NTlDQlx1NkU5MFx1NjU4N1x1NEVGNlxuICAgICAgLy8gICAvLyAgZ3ppcDogdHJ1ZSwgLy8gXHU2NjNFXHU3OTNBIGd6aXAgXHU1MzhCXHU3RjI5XHU1NDBFXHU3Njg0XHU1OTI3XHU1QzBGXG4gICAgICAvLyAgICBicm90bGlTaXplOiBmYWxzZSwgLy8gXHU2NjNFXHU3OTNBIGJyb3RsaSBcdTUzOEJcdTdGMjlcdTU0MEVcdTc2ODRcdTU5MjdcdTVDMEZcbiAgICAgIC8vICAgIG9wZW46IHRydWUsIC8vIFx1ODFFQVx1NTJBOFx1NzUxRlx1NjIxMFx1NjJBNVx1NTQ0QVx1NTQwRVx1ODFFQVx1NTJBOFx1NjI1M1x1NUYwMFx1NkQ0Rlx1ODlDOFx1NTY2OFxuICAgICAgLy8gICAgLy8gXHU1MTc2XHU0RUQ2XHU1M0VGXHU4OUM2XHU1MzE2XHU5MDA5XHU5ODc5Li4uXG4gICAgICAvLyAgfSksXG4gICAgXSxcbiAgcmVzb2x2ZToge1xuICAgIGFsaWFzOiBbXG4gICAgICB7IGZpbmQ6IC9efi8sIHJlcGxhY2VtZW50OiAnJyB9LFxuICAgICAgeyBmaW5kOiAnQGNvbW1vbicsIHJlcGxhY2VtZW50OiBwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnLi4vY29tbW9uL3NyYycpIH0sXG4gICAgICB7IGZpbmQ6ICdAbWFya2V0JywgcmVwbGFjZW1lbnQ6IHBhdGgucmVzb2x2ZShfX2Rpcm5hbWUsICcuLi9tYXJrZXQvc3JjJykgfSxcbiAgICAgIHsgZmluZDogJ0Bjb3JlJywgcmVwbGFjZW1lbnQ6IHBhdGgucmVzb2x2ZShfX2Rpcm5hbWUsICcuL3NyYycpIH0sXG4gICAgXVxuICB9LFxuICBzZXJ2ZXI6IHtcbiAgICBwcm94eToge1xuICAgICAgJy9hcGkvdjEnOiB7XG4gICAgICAgIC8vIHRhcmdldDogJ2h0dHA6Ly91YXQuYXBpa2l0LmNvbToxMTIwNC9tb2NrQXBpL2FvcGxhdGZvcm0vJyxcbiAgICAgICAgdGFyZ2V0OiAnaHR0cDovLzE3Mi4xOC4xNjYuMjE5OjgyODgvJyxcbiAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgfSxcbiAgICAgICcvYXBpMi92MSc6IHtcbiAgICAgICAgLy8gdGFyZ2V0OiAnaHR0cDovL3VhdC5hcGlraXQuY29tOjExMjA0L21vY2tBcGkvYW9wbGF0Zm9ybS8nLFxuICAgICAgICB0YXJnZXQ6ICdodHRwOi8vMTcyLjE4LjE2Ni4yMTk6ODI4OC8nLFxuICAgICAgICBjaGFuZ2VPcmlnaW46IHRydWUsXG4gICAgICB9XG4gICAgfVxuICB9LFxuICBsb2dMZXZlbDonaW5mbydcbn0pXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBTUEsU0FBUyxvQkFBb0I7QUFDN0IsT0FBTyxXQUFXO0FBQ2xCLE9BQU8sVUFBVTtBQUNqQixPQUFPLHVCQUF1QjtBQUM5QixPQUFPLGlCQUFpQjtBQUN4QixPQUFPLGtCQUFrQjtBQVh6QixJQUFNLG1DQUFtQztBQWV6QyxJQUFPLHNCQUFRLGFBQWE7QUFBQSxFQUMxQixVQUFVO0FBQUEsRUFDVixPQUFNO0FBQUEsSUFDSixRQUFPO0FBQUEsSUFDUCxXQUFXO0FBQUEsSUFDWCx1QkFBdUI7QUFBQSxJQUN2QixVQUFVO0FBQUEsSUFDUixRQUFRO0FBQUEsTUFDTixhQUFhLElBQUk7QUFDZixZQUFJLEdBQUcsU0FBUyxjQUFjLEdBQUc7QUFDL0IsaUJBQU8sR0FBRyxTQUFTLEVBQUUsTUFBTSxlQUFlLEVBQUUsQ0FBQyxFQUFFLE1BQU0sR0FBRyxFQUFFLENBQUMsRUFBRSxTQUFTO0FBQUEsUUFDeEU7QUFFQSxZQUFJLEdBQUcsU0FBUyxPQUFPLEdBQUc7QUFDeEIsZ0JBQU0sV0FBVyxHQUFHLE1BQU0sS0FBSyxHQUFHO0FBQ2xDLGdCQUFNLGNBQWMsU0FBUyxTQUFTLFFBQVEsT0FBTyxJQUFJLENBQUMsRUFBRSxNQUFNLEdBQUcsRUFBRSxDQUFDO0FBQ3hFLGlCQUFPO0FBQUEsUUFDVDtBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUEsRUFDRjtBQUFBLEVBQ0YsS0FBSztBQUFBLElBQ0gsU0FBUztBQUFBLE1BQ1AsU0FBUztBQUFBLFFBQ1AsWUFBWSxLQUFLLFFBQVEsa0NBQVcsOEJBQThCLENBQUM7QUFBQSxRQUNuRTtBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUEsSUFDQSxxQkFBcUI7QUFBQSxNQUNuQixNQUFNO0FBQUEsUUFDSixtQkFBbUI7QUFBQSxNQUNyQjtBQUFBLElBQ0Y7QUFBQSxJQUNBLFNBQVE7QUFBQSxNQUNOLGtCQUFpQjtBQUFBLE1BQ2pCLG9CQUFtQjtBQUFBLElBQ3JCO0FBQUEsRUFDRjtBQUFBLEVBQ0EsU0FBUztBQUFBLElBQUMsTUFBTTtBQUFBLElBQ1osa0JBQWtCO0FBQUEsTUFDaEIsU0FBUSxDQUFDLEtBQUs7QUFBQSxNQUNkLFNBQVEsQ0FBQztBQUFBLE1BQ1QsYUFBWTtBQUFBLElBQ2IsQ0FBQztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBVUo7QUFBQSxFQUNGLFNBQVM7QUFBQSxJQUNQLE9BQU87QUFBQSxNQUNMLEVBQUUsTUFBTSxNQUFNLGFBQWEsR0FBRztBQUFBLE1BQzlCLEVBQUUsTUFBTSxXQUFXLGFBQWEsS0FBSyxRQUFRLGtDQUFXLGVBQWUsRUFBRTtBQUFBLE1BQ3pFLEVBQUUsTUFBTSxXQUFXLGFBQWEsS0FBSyxRQUFRLGtDQUFXLGVBQWUsRUFBRTtBQUFBLE1BQ3pFLEVBQUUsTUFBTSxTQUFTLGFBQWEsS0FBSyxRQUFRLGtDQUFXLE9BQU8sRUFBRTtBQUFBLElBQ2pFO0FBQUEsRUFDRjtBQUFBLEVBQ0EsUUFBUTtBQUFBLElBQ04sT0FBTztBQUFBLE1BQ0wsV0FBVztBQUFBO0FBQUEsUUFFVCxRQUFRO0FBQUEsUUFDUixjQUFjO0FBQUEsTUFDaEI7QUFBQSxNQUNBLFlBQVk7QUFBQTtBQUFBLFFBRVYsUUFBUTtBQUFBLFFBQ1IsY0FBYztBQUFBLE1BQ2hCO0FBQUEsSUFDRjtBQUFBLEVBQ0Y7QUFBQSxFQUNBLFVBQVM7QUFDWCxDQUFDOyIsCiAgIm5hbWVzIjogW10KfQo=
