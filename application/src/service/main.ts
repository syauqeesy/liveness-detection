export type InferenceResult = {
  live: number;
  spoof: number;
  result: "Live" | "Spoof";
};

export interface InferenceServiceAdapter {
  Initialize(): Promise<void>;
  Predict(imageBlob: Blob, requestId: string): Promise<InferenceResult>;
}

