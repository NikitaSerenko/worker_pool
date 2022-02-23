# worker_pool
worker pool

Assumptions: from the task's description it looks like we don't need to think about durability of tasks (don't need tasks to be persistent). So just simple
worker pool was written.
The solution is successfully tested.

### curl examples

curl --location --request POST '127.0.0.1:8000/tasks' --header 'Content-Type: application/json' --data-raw '{"A":3, "B": 4, "C": 5}' -v

curl --location --request GET '127.0.0.1:8000/get_tasks'  -v
