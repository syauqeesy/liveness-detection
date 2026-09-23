import { BACKEND_URL } from "../config";

import type { InferenceAdapter, InferenceResult } from "./main";

type InferenceData = {
  live: number
  spoof: number
  result: "Live" | "Spoof"
}

type Response<ResponseData> = {
  message: string
  data: ResponseData
}

export default class CloudInference implements InferenceAdapter {
  async initialize(): Promise<void> {
    // No local model to initialize.
  }

  async execute(imageBlob: Blob): Promise<InferenceResult> {
    const body = new FormData();

    body.append("image", imageBlob, "liveness.jpg");

    const response = await fetch(
      `${BACKEND_URL}/inference`,
      {
        method: "POST",
        body,
      },
    );

    const result: Response<InferenceData> = await response.json();

    if (!response.ok) {
      throw new Error(result.message ?? "Inference failed");
    }

    return {
      live: result.data.live,
      spoof: result.data.spoof,
      result: result.data.result
    };
  }
}
