import base64
import binascii
import time
from typing import Any

from fastapi import APIRouter, HTTPException
from foundation.configuration import Configuration
from foundation.logger import Logger
from foundation.model import AntiSpoofModel
from pydantic import BaseModel, Field


class PredictionInstance(BaseModel):
    image: str = Field(min_length=1)


class PredictionRequest(BaseModel):
    instances: list[PredictionInstance] = Field(min_length=1)
    parameters: dict[str, Any] | None = None


class Prediction(BaseModel):
    result: str
    live: float
    spoof: float


class PredictionResponse(BaseModel):
    predictions: list[Prediction]


class InferenceService:
    def __init__(
        self,
        configuration: Configuration,
        logger: Logger,
        model: AntiSpoofModel,
    ):
        self._configuration = configuration
        self._logger = logger
        self._model = model

        self.router = APIRouter()

        self._register_routes()

    def _register_routes(self) -> None:
        self.router.add_api_route(
            "/health",
            self.health,
            methods=["GET"],
        )

        self.router.add_api_route(
            "/predict",
            self.predict,
            methods=["POST"],
            response_model=PredictionResponse,
        )

    def health(self) -> dict[str, str]:
        return {
            "status": "ok",
        }

    async def predict(
        self,
        request: PredictionRequest,
    ) -> PredictionResponse:
        predictions: list[Prediction] = []

        for instance in request.instances:
            try:
                image = base64.b64decode(
                    instance.image,
                    validate=True,
                )

            except (binascii.Error, ValueError):
                self._logger.warn(
                    "prediction rejected",
                    "reason",
                    "invalid base64 image",
                )

                raise HTTPException(
                    status_code=400,
                    detail="Image is not valid base64",
                )

            if not image:
                self._logger.warn(
                    "prediction rejected",
                    "reason",
                    "empty image",
                )

                raise HTTPException(
                    status_code=400,
                    detail="Image is empty",
                )

            started = time.perf_counter()

            try:
                result = self._model.execute(image)

                predictions.append(
                    Prediction(
                        result=result.result,
                        live=result.live,
                        spoof=result.spoof,
                    )
                )

            except Exception as error:
                self._logger.error(
                    "prediction failed",
                    "error",
                    str(error),
                )

                raise HTTPException(
                    status_code=500,
                    detail="Prediction failed",
                )

            finally:
                inference_time_ms = (time.perf_counter() - started) * 1000

                self._logger.info(
                    "prediction completed",
                    "duration_ms",
                    round(inference_time_ms, 2),
                )

        return PredictionResponse(
            predictions=predictions,
        )
