import { loadAndCompile, loadLiteRt, Tensor } from "@litertjs/core";

import type { InferenceAdapter, InferenceResult } from "./main";

type LiteRtModel = Awaited<ReturnType<typeof loadAndCompile>>;

export default class LocalInference implements InferenceAdapter {
  private model: LiteRtModel | null = null;

  async initialize(): Promise<void> {
    if (this.model) {
      return;
    }

    await loadLiteRt("/litert-wasm/");

    this.model = await loadAndCompile(
      "/models/anti-spoof-mn3_float32.tflite",
      {
        accelerator: "wasm",
      },
    );

    console.log("Liveness model loaded");
  }

  async execute(imageBlob: Blob): Promise<InferenceResult> {
    if (!this.model) {
      throw new Error("Model is not initialized")
    }

    const image = await createImageBitmap(imageBlob);

    try {
      if (image.width !== 128 || image.height !== 128) {
        throw new Error(`Expected 128x128 image, got ${image.width}x${image.height}`);
      }

      const canvas = document.createElement("canvas");

      canvas.width = 128;
      canvas.height = 128;

      const context = canvas.getContext("2d");

      if (!context) {
        throw new Error("Could not get canvas context");
      }

      context.drawImage(image, 0, 0);

      const imageData =
        context.getImageData(0, 0, 128, 128);

      const mean = [151.2405, 119.5950, 107.8395];

      const scale = [63.0105, 56.4570, 55.0035];

      const inputData = new Float32Array(128 * 128 * 3);

      for (let i = 0; i < 128 * 128; i++) {
        const pixelIndex = i * 4;
        const inputIndex = i * 3;

        inputData[inputIndex] = (imageData.data[pixelIndex] - mean[0]) / scale[0];

        inputData[inputIndex + 1] = (imageData.data[pixelIndex + 1] - mean[1]) / scale[1];

        inputData[inputIndex + 2] = (imageData.data[pixelIndex + 2] - mean[2]) / scale[2];
      }

      const inputTensor = new Tensor(inputData, [1, 128, 128, 3]);

      const results = await this.model.run(inputTensor);

      const output = await results[0].data();

      const live = output[0];
      const spoof = output[1];

      inputTensor.delete();
      results[0].delete();

      return {
        live,
        spoof,
        result:
          live >= spoof
            ? "Live"
            : "Spoof",
      };
    } finally {
      image.close();
    }
  }
}
