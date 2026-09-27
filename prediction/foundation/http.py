import uvicorn
from fastapi import FastAPI as FastAPIApplication
from foundation.configuration import Configuration


class Http:
    def __init__(self, configuration: Configuration, logger):
        self._logger = logger

        self._app = FastAPIApplication(
            title=configuration.application.service,
        )

        self._server = uvicorn.Server(
            uvicorn.Config(
                self._app,
                host="0.0.0.0",
                port=int(configuration.http.port),
                timeout_graceful_shutdown=15,
            )
        )

    @property
    def app(self) -> FastAPIApplication:
        return self._app

    def start(self) -> None:
        self._logger.info(
            "server running",
            "port",
            self._server.config.port,
        )

        self._server.run()

    def shutdown(self) -> None:
        self._server.should_exit = True


def new_http(configuration: Configuration, logger) -> Http:
    return Http(configuration, logger)
