# throughput_gradual

##### This Project aims to show how to get the most out of an api without making it go down.

#### The idea is to have a predefined number of tasks to be executed in a period, instead of calculating how many requests would be given per second in the period and doing this until the end,
#### we start the application with the minimum number of requests per second so that all tasks are completed on time in the same way, however, as each minute passes, we validate whether we can increase throughput,
#### this way we are gradually scaling the service called inside the worker, which allows for more calls per second, and of course, less execution time.

#### This algorithm also validates if the increase in thoughput caused a negative impact (number of errors increased), it returns to the previous throughput.

Example of the app being called by the worker:
![graph](https://user-images.githubusercontent.com/32064166/236351459-2bdd4c0f-0b2e-44df-b8f6-c3bcb4008cb5.png)
