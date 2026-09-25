import { BACKEND_URL } from "../config";

import type { InferenceServiceAdapter, InferenceResult } from "./main";

type InferenceData = {
  live: number
  spoof: number
  result: "Live" | "Spoof"
}

type Response<ResponseData> = {
  message: string
  data: ResponseData
}

export default class SelfManagedService implements InferenceServiceAdapter {
  async Initialize(): Promise<void> {
    // No local model to initialize.
  }

  async Predict(imageBlob: Blob): Promise<InferenceResult> {
    const body = new FormData();

    body.append("mode", "self_managed_service")
    body.append("image", imageBlob, "liveness.jpg");

    const response = await fetch(
      `${BACKEND_URL}/inference/predict`,
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


