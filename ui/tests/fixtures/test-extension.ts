import type { Page } from "@playwright/test";
import { test as base } from "@playwright/test";
import { finishCoverage, startCoverage } from "./coverage";

export interface TestOptions {
  hasCoverage: boolean;
  runCoverage: Page;
}

export const test = base.extend<TestOptions>({
  hasCoverage: [false, { option: true }],
  runCoverage: [
    async ({ page, hasCoverage }, use) => {
      if (hasCoverage) {
        await startCoverage(page);
        await use(page);
        await finishCoverage(page);
      } else {
        await use(page);
      }
    },
    { auto: true },
  ],
});

export const expect = test.expect;
