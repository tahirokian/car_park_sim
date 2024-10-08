#!/usr/bin/env python3
# encoding: utf-8

import csv
import json
import os
import time
from pathlib import Path

from flask import Flask, jsonify, request
from flask import render_template
from prometheus_client import Counter, Gauge, start_http_server

app = Flask(__name__)

_HTTP_REQUEST_COUNT = Counter(namespace='web_backend', name='http_requests_count', 
                              documentation='Total number of http requests', labelnames=['method'])
_HTTP_REQUEST_TIME = Gauge(namespace='web_backend', name='http_request_duration_seconds',
                           documentation='Time spent processing http request', 
                           labelnames=['method'])


@app.route("/", methods=['GET'])
def index():
    start_time = time.time()
    rendered = render_template('default.html') if not Path("static/vehicle_parking.csv").exists() \
        else render_template('index.html')
    _HTTP_REQUEST_COUNT.labels(method='GET').inc()
    _HTTP_REQUEST_TIME.labels(method='GET').set(time.time() - start_time)
    return rendered


@app.route('/upload', methods=['POST'])
def upload_to_file():
    start_time = time.time()
    record = json.loads(request.data)
    file_exists = Path("static/vehicle_parking.csv").exists()

    with open('static/vehicle_parking.csv', 'a', newline='', encoding='utf-8') as csv_file:
        csv_writer = csv.DictWriter(csv_file, fieldnames=record.keys())
        if not file_exists:
            csv_writer.writeheader()
        csv_writer.writerow(record)

    _HTTP_REQUEST_COUNT.labels(method='POST').inc()
    _HTTP_REQUEST_TIME.labels(method='POST').set(time.time() - start_time)
    return jsonify({"status": "success"})


if __name__ == "__main__":
    web_addr = os.getenv("WEBBACKEND_ADDR", default = "web_backend")
    web_port = os.getenv("WEBBACKEND_PORT", default = 40000)
    metrics_port = os.getenv("WEBBACKEND_METRICS_PORT", default = 40100)

    start_http_server(port=int(metrics_port))
    print(f"Python web server metrics available at http://{web_addr}:{metrics_port}/metrics")

    print(f"Python web server is available at http://{web_addr}:{web_port}")
    app.run(host="0.0.0.0", port=int(web_port))
