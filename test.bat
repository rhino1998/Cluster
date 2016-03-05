
for /L %%A in (1,1,50000) do (
    curl -X POST -m 1 -d @testtask.json http://108.56.251.125:2002/api/task
)