import {
  test,
  expect,
  type Page,
  type BrowserContext,
} from '@playwright/test';
import { appendFile, mkdir } from 'node:fs/promises';
import path from 'node:path';

const BASE_URL = 'http://localhost:5173';

const ITERATION = 5;

const RESULT_DIR = path.resolve('./test-results');
const RESULT_FILE = path.join(RESULT_DIR, 'e2e-latency.jsonl');

const CAMERA_TIMEOUT = 30_000;

async function grantCameraPermission(context: BrowserContext) {
  await context.grantPermissions(['camera'], {
    origin: BASE_URL,
  });
}

async function waitForCameraReady(page: Page) {
  const video = page.locator('video').first();

  await expect(video).toBeVisible({
    timeout: CAMERA_TIMEOUT,
  });

  await expect
    .poll(
      async () => {
        return await video.evaluate((element) => {
          const video = element as HTMLVideoElement;

          return {
            readyState: video.readyState,
            videoWidth: video.videoWidth,
            videoHeight: video.videoHeight,
            paused: video.paused,
          };
        });
      },
      {
        timeout: CAMERA_TIMEOUT,
        intervals: [100, 250, 500],
      },
    )
    .toMatchObject({
      readyState: expect.any(Number),
      videoWidth: expect.any(Number),
      videoHeight: expect.any(Number),
      paused: false,
    });

  await page.waitForFunction(
    () => {
      const video = document.querySelector('video');

      return (
        video instanceof HTMLVideoElement &&
        video.readyState >= HTMLMediaElement.HAVE_METADATA &&
        video.videoWidth > 0 &&
        video.videoHeight > 0 &&
        !video.paused
      );
    },
    undefined,
    {
      timeout: CAMERA_TIMEOUT,
    },
  );
}

async function waitForInitialLoadingToFinish(page: Page) {
  const spinner = page.locator('.animate-spin');

  await expect(spinner).toBeHidden({
    timeout: CAMERA_TIMEOUT,
  });
}

async function releaseCamera(page: Page) {
  try {
    await page.evaluate(() => {
      const videos = document.querySelectorAll('video');

      for (const video of videos) {
        if (!(video instanceof HTMLVideoElement)) {
          continue;
        }

        video.pause();

        const stream = video.srcObject;

        if (stream instanceof MediaStream) {
          for (const track of stream.getTracks()) {
            track.stop();
          }

          video.srcObject = null;
        }
      }
    });
  } catch {
    // The page may already be closed/navigated after a test failure.
  }

  try {
    await page.goto('about:blank', {
      waitUntil: 'domcontentloaded',
    });
  } catch {
    // Ignore teardown navigation errors.
  }
}

async function saveResult(result: {
  service: string;
  timestamp: string;
  latency_ms: number;
}) {
  await mkdir(RESULT_DIR, { recursive: true });

  await appendFile(RESULT_FILE, `${JSON.stringify(result)}\n`, 'utf-8');
}

async function testLiveness(page: Page, context: BrowserContext, service: string, buttonName: string) {
  try {
    await grantCameraPermission(context);

    await page.goto(BASE_URL, {
      waitUntil: 'domcontentloaded',
    });

    await waitForCameraReady(page);
    await waitForInitialLoadingToFinish(page);
    await page.waitForTimeout(600);

    const start = performance.now();

    await page
      .getByRole('button', {
        name: buttonName,
        exact: true,
      })
      .click();

    const closeButton = page.getByRole('button', {
      name: 'Close',
      exact: true,
    });

    await expect(closeButton).toBeVisible({
      timeout: CAMERA_TIMEOUT,
    });

    const latency = performance.now() - start;

    await saveResult({
      service,
      timestamp: new Date().toISOString(),
      latency_ms: Number(latency.toFixed(2)),
    });

    await closeButton.click();
  } finally {
    await releaseCamera(page);
  }
}

test('self managed service liveness', async ({ page, context }) => {
  let current = 1
  
  do {
    await testLiveness(
      page,
      context,
      'self-managed',
      'Self Managed Service Liveness Check',
    );

    current++
  } while(current <= ITERATION);
});

test('managed service liveness', async ({ page, context }) => {
  let current = 1
  
  do {
    await testLiveness(
      page,
      context,
      'managed',
      'Managed Service Liveness Check',
    );

    current++
  } while(current <= ITERATION);
});

test('on device service liveness', async ({ page, context }) => {
  let current = 1
  
  do {
    await testLiveness(
      page,
      context,
      'on-device',
      'On Device Service Liveness Check',
    );

    current++
  } while(current <= ITERATION);
});
