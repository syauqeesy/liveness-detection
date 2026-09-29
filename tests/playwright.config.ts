import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: '.',

  use: {
    baseURL: 'http://localhost:5173',
    channel: 'chrome',
    headless: false,
    permissions: ['camera'],
  },

  workers: 1,
});

