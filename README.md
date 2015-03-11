## bloom

Just another workaday bloom filter. Use `ComputeOptimal` to get the best parameters (bits per element / number of hash functions) for a given false positive probability.

### Usage

`go get github.com/rtecco/bloom`

See `bf_test.go` for a sample usage.

`optimal.go` is a port of [Cassandra's BloomCalculations.java](https://raw.githubusercontent.com/facebookarchive/cassandra/master/src/org/apache/cassandra/utils/BloomCalculations.java) (MIT licensed)

`murmur32.go` was adapted from [reusee/mmh3](https://github.com/reusee/mmh3) (MIT licensed) to take a seed parameter.

### License

MIT
