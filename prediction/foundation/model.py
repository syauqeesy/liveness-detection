from dataclasses import dataclass
from io import BytesIO
from pathlib import Path

import numpy as np
import tensorflow as tf
from PIL import Image


@dataclass(frozen=True)
class InferenceResult:
    result: str
    live: float
    spoof: float


class AntiSpoofModel:
    def __init__(self, model_path: Path, memory_limit: int):
        self.configure_gpu(memory_limit)

        self.model = tf.saved_model.load(str(model_path))

        if "serving_default" not in self.model.signatures:
            raise RuntimeError(
                f"SavedModel does not have a 'serving_default' signature. "
                f"Available signatures: {list(self.model.signatures.keys())}"
            )

        self.infer = self.model.signatures["serving_default"]

        self.input_signature = self.infer.structured_input_signature
        self.output_signature = self.infer.structured_outputs

        self.mean = np.array([151.2405, 119.5950, 107.8395], dtype=np.float32)

        self.scale = np.array([63.0105, 56.4570, 55.0035], dtype=np.float32)

    def configure_gpu(self, memory_limit: int) -> None:
        gpus = tf.config.list_physical_devices("GPU")

        if not gpus:
            return

        tf.config.set_logical_device_configuration(
            gpus[0],
            [
                tf.config.LogicalDeviceConfiguration(
                    memory_limit=memory_limit
                )
            ],
        )

    def execute(self, image_bytes: bytes) -> InferenceResult:
        image = Image.open(BytesIO(image_bytes)).convert("RGB")
        image = np.array(image, dtype=np.float32)
        image = (image - self.mean) / self.scale

        input_tensor = np.expand_dims(image, axis=0)
        input_tensor = tf.convert_to_tensor(input_tensor, dtype=tf.float32)

        input_name = list(self.infer.structured_input_signature[1].keys())[0]

        result = self.infer(**{input_name: input_tensor})

        outputs = next(iter(result.values())).numpy()

        live = float(outputs[0][0])
        spoof = float(outputs[0][1])

        return InferenceResult(
            result="Live" if live >= spoof else "Spoof", live=live, spoof=spoof
        )


def new_model(model_path: Path, memory_limit: int) -> AntiSpoofModel:
    return AntiSpoofModel(model_path, memory_limit)
