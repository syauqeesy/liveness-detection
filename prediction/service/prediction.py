import base64
import binascii
import time

from foundation.logger import Logger
from foundation.model import AntiSpoofModel, InferenceResult


class InvalidImageError(Exception):
    pass


class PredictionError(Exception):
    pass


class PredictionService:
    def __init__(
        self,
        logger: Logger,
        model: AntiSpoofModel,
    ):
        self._logger = logger
        self._model = model

    def predict(
        self,
        request_id: str,
        image_base64: str,
    ) -> InferenceResult:
        try:
            image = base64.b64decode(
                image_base64,
                validate=True,
            )

        except (binascii.Error, ValueError) as error:
            self._logger.warn(
                "prediction rejected",
                "request_id",
                request_id,
                "reason",
                "invalid base64 image",
            )

            raise InvalidImageError(
                "Image is not valid base64",
            ) from error

        if not image:
            self._logger.warn(
                "prediction rejected",
                "request_id",
                request_id,
                "reason",
                "empty image",
            )

            raise InvalidImageError(
                "Image is empty",
            )

        started = time.perf_counter()

        try:
            return self._model.execute(image)

        except Exception as error:
            self._logger.error(
                "prediction failed",
                "request_id",
                request_id,
                "error",
                str(error),
            )

            raise PredictionError(
                "Prediction failed",
            ) from error

        finally:
            inference_time_ms = (time.perf_counter() - started) * 1000

            self._logger.info(
                "prediction completed",
                "request_id",
                request_id,
                "duration_ms",
                round(inference_time_ms, 2),
            )
