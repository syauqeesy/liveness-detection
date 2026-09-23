import numpy as np
import tensorflow as tf
from io import BytesIO
from PIL import Image
from pathlib import Path


class AntiSpoofModel:
    def __init__(self, model_path: Path):
        self.interpreter = tf.lite.Interpreter(model_path=str(model_path))

        self.interpreter.allocate_tensors()

        self.input_details = self.interpreter.get_input_details()
        self.output_details = self.interpreter.get_output_details()

        self.mean = np.array([151.2405, 119.5950, 107.8395], dtype=np.float32)

        self.scale = np.array([63.0105, 56.4570, 55.0035], dtype=np.float32)

    def execute(self, image_bytes: bytes):
        image = Image.open(BytesIO(image_bytes)).convert("RGB")

        image = np.array(image, dtype=np.float32)

        image = (image - self.mean) / self.scale

        input_tensor = np.expand_dims(image, axis=0)

        self.interpreter.set_tensor(
            self.input_details[0]["index"], input_tensor)

        self.interpreter.invoke()

        output = self.interpreter.get_tensor(self.output_details[0]["index"])

        live = float(output[0][0])
        spoof = float(output[0][1])

        return {
            "result": "Live" if live >= spoof else "Spoof",
            "live": live,
            "spoof": spoof,
        }


def new_model(model_path: Path) -> AntiSpoofModel:
    return AntiSpoofModel(model_path)
