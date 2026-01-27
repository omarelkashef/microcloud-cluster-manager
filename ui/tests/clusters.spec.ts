import { test } from "@playwright/test";
import { expect } from "./fixtures/test-extension";

test.describe("Cluster Bulk Configuration", () => {
  test("should bulk configure clusters ", async ({ page }) => {
    await page.goto("/ui/clusters");
    await page.waitForSelector(".scrollable-table");
    const rows = page
      .locator(".clusterlist-table")
      .locator("tbody")
      .locator("tr");
    const rowCount = await rows.count();

    await page.getByLabel("Select all").click();
    await page
      .getByText(
        `All ${rowCount} cluster${rowCount !== 1 ? "s" : ""} selected.`,
      )
      .click();

    await page.getByRole("button", { name: "Configure clusters" }).click();
    await page.getByTitle("Set Disk threshold").click();
    await page.getByLabel("Disk threshold").fill("100");
    await page.getByTitle("Set Memory threshold").click();
    await page.getByLabel("Memory threshold").fill("100");
    await page.getByRole("button", { name: "Save changes" }).click();

    await expect(
      page.getByText(`Updated ${rowCount} cluster${rowCount !== 1 ? "s" : ""}`),
    ).toBeVisible();
  });
});
