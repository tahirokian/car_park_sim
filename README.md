# car_park_simulator
This code is tested on Ubuntu LTS Server 24 and Docker version 24.0.5.

* `go_backend_service` consumes events from RabbitMQ and maintain vehicle entry/exit records
in Redis. It also call a Python Flask server to write the vehicle summary into a local file
once an exit event is processed.

* `python_server` receives the vehicle summary from `go_backend_service` and writes it to a local
file. The file format is csv.

* `rabbitmq_config_service` configures RabbitMQ with two queues: one for tracking vehicle entry
events (vehicle_entry) and another for tracking vehicle exit events (vehicle_exit).

* `vehicle_entry_service` generates vehicle entry events and publishes them to RabbitMQ queue
`vehicle_entry`. Each event is a JSON payload with fields: `id`, `vehicle_plate`, and
`entry_date_time`.

* `vehicle_exit_service` generates vehicle exit events and publishes them to RabbitMQ queue
`vehicle_exit`. Each event is a JSON payload with fields: `id`, `vehicle_plate`, and
`exit_date_time`.

* `prometheus` provides the instrumentation.

* `redis` All data to redis is saved in db0 by the `go_backend_service`. Following queries can
be made after connecting to db0.

- List all keys
```
KEYS *
```
- List set members
```
SMEMBERS entry_plates
```
- List value associated with field in the hash stored at key
```
HGET vehicle_palte entry_time
HGET vehicle_palte parking_enter_data
HGET vehicle_palte parking_exit_data
```

* `python_server` is available at `http://ip:40000`

* `rabbitmq` management interface is available at `http://ip:15762`

* `prometheus` server is available at `http://ip:9090`

* `grafana` server is available at `http://ip:3000`

- To query metrics from different services, use following namespaces to list available metrics and
  build grafana dashboards:
    - go_backend
    - web_backend
    - queue_configurator
    - vehicle_entry
    - vehicle_exit

## Use below command to see all possible options to interact with the services
```
make help
```
