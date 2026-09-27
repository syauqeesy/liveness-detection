import argparse
from pathlib import Path

from foundation.configuration import new_configuration
from foundation.grpc import new_grpc
from foundation.http import new_http
from foundation.logger import new_logger
from foundation.model import new_model
from service.grpc.main import new_service as new_grpc_service
from service.main import new_service as new_http_service

configuration = new_configuration(Path("./config.json"))

model = new_model(
    Path("./anti-spoof-mn3"),
    configuration.application.memory_limit,
)

logger = new_logger(
    configuration.application.service,
    configuration.application.environment,
)


def parse_args():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "mode",
        choices=["grpc", "http"],
        help="Server mode to run",
    )

    return parser.parse_args()


def main():
    args = parse_args()

    if args.mode == "grpc":
        grpc = new_grpc(configuration, logger)

        new_grpc_service(
            grpc.server,
            configuration,
            logger,
            model,
        )

        grpc.start()
        grpc.wait_for_termination()

    elif args.mode == "http":
        fastapi = new_http(configuration, logger)

        new_http_service(
            fastapi.app,
            configuration,
            logger,
            model,
        )

        fastapi.start()


if __name__ == "__main__":
    main()
