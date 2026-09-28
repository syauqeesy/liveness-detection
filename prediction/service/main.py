from fastapi import FastAPI

from foundation.configuration import Configuration
from foundation.logger import Logger
from foundation.model import AntiSpoofModel
from service.inference import InferenceService


def new_service(
    application: FastAPI,
    configuration: Configuration,
    logger: Logger,
    model: AntiSpoofModel,
):
    inference_service = InferenceService(
        configuration,
        logger,
        model,
    )

    application.include_router(inference_service.router)
