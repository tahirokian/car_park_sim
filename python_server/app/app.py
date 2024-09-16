#!/usr/bin/env python3
# encoding: utf-8

import csv
import json
import os
from pathlib import Path

from flask import Flask, jsonify, request
from flask import render_template

app = Flask(__name__)


@app.route("/", methods=['GET'])
def index():
    if not Path("static/vehicle_parking.csv").exists():
        return render_template('default.html')
    return render_template('index.html')


@app.route('/upload', methods=['POST'])
def upload_to_file():
    record = json.loads(request.data)
    file_exists = Path("static/vehicle_parking.csv").exists()
    
    with open('static/vehicle_parking.csv', 'a', newline='') as csv_file:
        csv_writer = csv.DictWriter(csv_file, fieldnames=record.keys())
        if not file_exists:
            csv_writer.writeheader()
        csv_writer.writerow(record)

    return jsonify({"status": "success"})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=os.getenv("WEB_BACKEND_PORT", default = 40000))
