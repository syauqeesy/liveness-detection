from fastapi import FastAPI

from foundation.logger import Logger
from foundation.model import AntiSpoofModel
from service.inference import InferenceService
from service.prediction import PredictionService


def new_service(
    application: FastAPI,
    logger: Logger,
    model: AntiSpoofModel,
):
    prediction_service = PredictionService(
        logger,
        model,
    )

    inference_service = InferenceService(
        prediction_service,
    )

    application.include_router(
        inference_service.router,
    )
