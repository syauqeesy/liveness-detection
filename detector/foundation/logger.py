import json
import logging
import sys
from datetime import datetime, timezone
from typing import Any


class JsonFormatter(logging.Formatter):
    def format(self, record: logging.LogRecord) -> str:
        payload: dict[str, Any] = {
            "time": datetime.now(timezone.utc).astimezone().isoformat(),
            "level": record.levelname,
            "msg": record.getMessage(),
        }

        fields = getattr(record, "fields", {})
        payload.update(fields)

        return json.dumps(
            payload,
            ensure_ascii=False,
            default=str,
        )


class Logger:
    def __init__(
        self,
        service: str,
        environment: str,
        level: int = logging.INFO,
    ):
        self._logger = logging.getLogger(f"{service}.{environment}")

        self._logger.setLevel(level)
        self._logger.propagate = False

        self._logger.handlers.clear()

        handler = logging.StreamHandler(sys.stdout)
        handler.setLevel(level)
        handler.setFormatter(JsonFormatter())

        self._logger.addHandler(handler)

        self._service = service
        self._environment = environment

    def _log(
        self,
        level: int,
        msg: str,
        args: tuple[Any, ...],
    ) -> None:
        if len(args) % 2 != 0:
            raise ValueError("log arguments must be key/value pairs")

        fields = {str(args[index]): args[index + 1] for index in range(0, len(args), 2)}

        fields["service"] = self._service
        fields["environment"] = self._environment

        self._logger.log(
            level,
            msg,
            extra={"fields": fields},
        )

    def debug(self, msg: str, *args: Any) -> None:
        self._log(logging.DEBUG, msg, args)

    def info(self, msg: str, *args: Any) -> None:
        self._log(logging.INFO, msg, args)

    def warn(self, msg: str, *args: Any) -> None:
        self._log(logging.WARNING, msg, args)

    def error(self, msg: str, *args: Any) -> None:
        self._log(logging.ERROR, msg, args)


def new_logger(service: str, environment: str) -> Logger:
    return Logger(service, environment)
