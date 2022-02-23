# worker_pool
worker pool

curl --location --request POST '127.0.0.1:8000/tasks' --header 'Content-Type: application/json' --data-raw '{"A":3, "B": 4, "C": 5}' -v
curl --location --request GET '127.0.0.1:8000/get_tasks'  -v
