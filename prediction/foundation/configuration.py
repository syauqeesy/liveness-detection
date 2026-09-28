import json
import os
from dataclasses import dataclass
from pathlib import Path


@dataclass
class ApplicationConfiguration:
    service: str
    environment: str
    memory_limit: int


@dataclass
class GRPCConfiguration:
    port: str


@dataclass
class HTTPConfiguration:
    port: str


@dataclass
class Configuration:
    application: ApplicationConfiguration
    grpc: GRPCConfiguration
    http: HTTPConfiguration


def new_configuration(path: Path) -> Configuration:
    data = {}

    if path.exists():
        with path.open("r", encoding="utf-8") as file:
            data = json.load(file)

    application_data = data.get("application", {})
    grpc_data = data.get("grpc", {})
    http_data = data.get("http", {})

    service = os.getenv(
        "APPLICATION_SERVICE",
        application_data.get("service"),
    )

    environment = os.getenv(
        "APPLICATION_ENVIRONMENT",
        application_data.get("environment"),
    )

    memory_limit = os.getenv(
        "APPLICATION_MEMORY_LIMIT",
        application_data.get("memory_limit"),
    )

    grpc_port = os.getenv(
        "GRPC_PORT",
        grpc_data.get("port"),
    )

    http_port = os.getenv(
        "HTTP_PORT",
        http_data.get("port"),
    )

    if service is None:
        raise ValueError("APPLICATION_SERVICE is not configured")

    if environment is None:
        raise ValueError("APPLICATION_ENVIRONMENT is not configured")

    if memory_limit is None:
        raise ValueError("APPLICATION_MEMORY_LIMIT is not configured")

    if grpc_port is None:
        raise ValueError("GRPC_PORT is not configured")

    if http_port is None:
        raise ValueError("HTTP_PORT is not configured")

    try:
        memory_limit = int(memory_limit)
    except (TypeError, ValueError) as exc:
        raise ValueError("APPLICATION_MEMORY_LIMIT must be an integer") from exc

    return Configuration(
        application=ApplicationConfiguration(
            service=service,
            environment=environment,
            memory_limit=memory_limit,
        ),
        grpc=GRPCConfiguration(
            port=grpc_port,
        ),
        http=HTTPConfiguration(
            port=http_port,
        ),
    )
