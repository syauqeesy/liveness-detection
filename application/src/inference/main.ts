export type InferenceResult = {
  live: number;
  spoof: number;
  result: "Live" | "Spoof";
};

export interface InferenceAdapter {
  initialize(): Promise<void>;
  execute(imageBlob: Blob): Promise<InferenceResult>;
}
