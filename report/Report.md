# Report

## Development plan

1) Create an initial full pipeline for the following components with simple/fake implementations:
    - input
    - output
    - counter
    - timestamp processing
2) Improve components one-by-one with unit tests
    1) Implement working but simple counter with sets
    2) Implement advanced timestamp processing
    3) Test/benchmark ingestion from Kafka (docker-compose)
    4) Output to Kafka
3) Benchmark, profile, improve performance
4) Try a probabilistic cardinality estimator - HyperLogLog++?   
5) Improve understanding of Kafka and document/implement more scalable solution


## Benchmark

### Setup

For benchmarking I used Golang's built-in benchmarking infrastructure and added code for generating
a CPU profile. I decided to only benchmark the input parsing, processing, counting, and output stage,
and in particular not include the possible overhead from Kafka.

### Results

Overall, we achieved an ingestion performance of 200k records/s, or 8MiB/s. This corresponds to about 10M records
per minute and does not seem to be I/O bound, i.e. I would expect Kafka being able to deliver higher data rates.

![ggprof graph](./profile.svg)

### Improvements

## Discussion

### Output format
Eventually we will have to produce statics per hour/day/week/month/year.
In order to solve this requirement we could have a copy of our software collecting continuously data continuously for a
longer period of time. In my view this would not be an optimal solution, as we would have to handle and recover from crashes
of machines.

Instead I will not only return a user count per minute, but also a raw aggregatable counter. This will not only allow
for aggregating the minutely counts to hourly, daily, etc, but also allow for setting up redundant instances of the
counting algorithm for the same partition. 

### JSON serialization
TODO

### Scalability
TODO memory usage
TODO input partitioning - by uid?
TODO output partitioning - by ts?

### Probabilistic cardinality estimation
TODO error

### Error recovery
TODO

### Late/Out-of-order frames
TODO

 