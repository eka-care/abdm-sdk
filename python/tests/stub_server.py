"""Hermetic HTTP stub server for tests: http.server on port 0 in a background
thread. Subclasses/handlers provide the routing; no real network is touched.
"""

import json
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer


class StubServer:
    """Wraps an HTTPServer running a given handler class in a daemon thread.

    Usage:
        server = StubServer(MyHandler)
        server.start()
        ...  # server.url
        server.stop()

    Or as a context manager.
    """

    def __init__(self, handler_class):
        self._httpd = HTTPServer(("127.0.0.1", 0), handler_class)
        self._thread = threading.Thread(target=self._httpd.serve_forever, daemon=True)

    @property
    def url(self) -> str:
        host, port = self._httpd.server_address[:2]
        return "http://%s:%d" % (host, port)

    def start(self):
        self._thread.start()
        return self

    def stop(self):
        self._httpd.shutdown()
        self._httpd.server_close()
        self._thread.join(timeout=5)

    def __enter__(self):
        return self.start()

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.stop()


class JSONHandler(BaseHTTPRequestHandler):
    """Base handler with helpers for reading a JSON body and writing a JSON
    response. Silences default request logging."""

    def log_message(self, fmt, *args):
        pass

    def read_json(self):
        length = int(self.headers.get("Content-Length", 0))
        raw = self.rfile.read(length) if length else b""
        return json.loads(raw) if raw else {}

    def write_json(self, status: int, payload: dict):
        body = json.dumps(payload).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)
