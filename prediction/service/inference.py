from typing import Any

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel, Field

from service.prediction import InvalidImageError, PredictionError, PredictionService


class PredictionInstance(BaseModel):
    request_id: str = Field(min_length=1)
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
        prediction_service: PredictionService,
    ):
        self._prediction_service = prediction_service

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
                result = self._prediction_service.predict(
                    request_id=instance.request_id,
                    image_base64=instance.image,
                )

            except InvalidImageError as error:
                raise HTTPException(
                    status_code=400,
                    detail=str(error),
                ) from error

            except PredictionError as error:
                raise HTTPException(
                    status_code=500,
                    detail=str(error),
                ) from error

            predictions.append(
                Prediction(
                    result=result.result,
                    live=result.live,
                    spoof=result.spoof,
                )
            )

        return PredictionResponse(
            predictions=predictions,
        )
