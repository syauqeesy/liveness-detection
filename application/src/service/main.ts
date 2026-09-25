export type InferenceResult = {
  live: number;
  spoof: number;
  result: "Live" | "Spoof";
};

export interface InferenceServiceAdapter {
  Initialize(): Promise<void>;
  Predict(imageBlob: Blob): Promise<InferenceResult>;
}

