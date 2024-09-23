#!/usr/bin/env python3
# encoding: utf-8

import csv
import json
import os
import time
from pathlib import Path

from flask import Flask, jsonify, request
from flask import render_template
from prometheus_client import start_http_server, Summary, Counter, Gauge

app = Flask(__name__)

GET_REQUEST_SUMMARY = Summary('get_request_processing_seconds', 'Time spent processing get request')
_GET_REQUEST_COUNT = Counter('get_request_count', 'Total number of get requests')
_GET_REQUEST_TIME = Gauge('get_gauge_request_seconds', 'Time spent processing get request')
POST_REQUEST_SUMMARY = Summary('post_request_processing_seconds', 'Time spent processing post request')
_POST_REQUEST_COUNT = Counter('post_request_count', 'Total number of post requests')
_POST_REQUEST_TIME = Gauge('post_gauge_request_seconds', 'Time spent processing post request')



@GET_REQUEST_SUMMARY.time()
@app.route("/", methods=['GET'])
def index():
    start_time = time.time()
    rendered = render_template('default.html') if not Path("static/vehicle_parking.csv").exists() \
        else render_template('index.html')
    _GET_REQUEST_COUNT.inc()
    _GET_REQUEST_TIME.set(time.time() - start_time)
    return rendered


@POST_REQUEST_SUMMARY.time()
@app.route('/upload', methods=['POST'])
def upload_to_file():
    start_time = time.time()
    record = json.loads(request.data)
    file_exists = Path("static/vehicle_parking.csv").exists()

    with open('static/vehicle_parking.csv', 'a', newline='') as csv_file:
        csv_writer = csv.DictWriter(csv_file, fieldnames=record.keys())
        if not file_exists:
            csv_writer.writeheader()
        csv_writer.writerow(record)

    _POST_REQUEST_COUNT.inc()
    _POST_REQUEST_TIME.set(time.time() - start_time)
    return jsonify({"status": "success"})


if __name__ == "__main__":
    web_addr = os.getenv("WEBBACKEND_ADDR", default = "web_backend")
    web_port = os.getenv("WEBBACKEND_PORT", default = 40000)
    metrics_port = os.getenv("WEBBACKEND_METRICS_PORT", default = 40100)

    start_http_server(port=int(metrics_port))
    print(f"Python web server metrics available at http://{web_addr}:{metrics_port}/metrics")

    print(f"Python web server is available at http://{web_addr}:{web_port}")
    app.run(host="0.0.0.0", port=int(web_port))
