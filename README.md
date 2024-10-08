# Car Park Simulator

## Task
Most of the shopping malls nowadays have automated parking metering. There are cameras that track
the entry and exit of the vehicle, and provide a realtime summary of vehicle time spent in parking
lot. This in turn is used to send parking invoices to the vehicle owner. The task here is to
simulate this use case using GoLang, Python, Redis, RabbitMQ and Docker.

## Requirements
1. Configure RabbitMQ with two queues - one for tracking vehicle entry events and another one for
tracking vehicle exit events.

2. To simulate vehile entry, create a generator service in GoLang that generates events with a json
payload having at least below mentioned fields:

 ```
{
  "id": <identifier for the event>,
  "vehicle_plate": <alphanumeric registration id of the vehicle>,
  "entry_date_time": <date time in UTC>
}
```

3. Publish above events to RabbitMQ queue for tracking entry events.

4. To simulate vehile exit, create a generator service in GoLang that generates events with a json
payload having at least below mentioned fields:

```
{
  "id": <identifier for the event>,
  "vehicle_plate": <alphanumeric registration id of the vehicle>,
  "exit_date_time": <date time in UTC>,
}
```

5. Publish above events to RabbitMQ queue for tracking exit events.

6. 80% of the exit events generated should match a vehicle plate that has a corresponding entry
event. Remianing exit events should have no correspnding entry events. This is to simulate the case
where the mall camera did not record vehicle entry event, perhaps because one or more of the vehicle
plates were not clean/clear enough!

7. Design and implement a backend service in GoLang that consumes events from above services and
maintains a record of vehicle entry and exit times.

8. Once the exit event is triggered, backend service invokes a REST API to a Python server that
writes summary of vehicle, entry, exit and duration into a local file.

9. For all storage needs, use a single redis instance for all services. DBs in redis should be
different for each service, if used.

10. Include basic statistics in services to measure latency of event processing. Use open source
tools such as prometheus, grafana, signoz etc.

## Deployment

1. Deploy all the services mentioned above (generator services, backend service, rest api server,
redis and rabbitmq) as docker containers.

2. All events and services are asynchronous.

## Implementation Details
Implementation is tested on Ubuntu LTS Server 20 and Docker version 24.0.5.

* `go_backend_service` consumes events from RabbitMQ and maintain vehicle entry/exit records
in Redis. It also calls a Python Flask server to write the vehicle summary into a local file
once an exit event is processed.

* `python_server` receives the vehicle summary from `go_backend_service` and writes it to a local
file. The file format is csv.

* `queue_configurator` configures RabbitMQ with two queues: one for tracking vehicle entry
events (vehicle_entry) and another for tracking vehicle exit events (vehicle_exit).

* `vehicle_entry_service` generates vehicle entry events and publishes them to RabbitMQ queue
`vehicle_entry`. Each event is a JSON payload with fields: `id`, `vehicle_plate`, and
`entry_date_time`.

* `vehicle_exit_service` generates vehicle exit events and publishes them to RabbitMQ queue
`vehicle_exit`. Each event is a JSON payload with fields: `id`, `vehicle_plate`, and
`exit_date_time`.

* `prometheus` and `grafana` provides the instrumentation.

* All data to `redis` is saved in db0 by the `go_backend_service`. Following queries can
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
