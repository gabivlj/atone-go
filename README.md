# atone in Go

### Implementation of atone (rust) in Go

This is my implementation of atone in Go, I think it can be faster but it's still under development

### What is atone?

Atone is an array data structure implementation in Rust created by @jonhoo [Jon Gjengset](https://github.com/jonhoo),
you can see it in [here](https://github.com/jonhoo/atone).

- As explained by its creator:

```

Most vector-like implementations, such as Vec and VecDeque, must occasionally "resize" the backing memory for the vector as the number of elements grows. This means allocating a new vector (usually of twice the size), and moving all the elements from the old vector to the new one. As your vector gets larger, this process takes longer and longer.

For most applications, this behavior is fine — if some very small number of pushes take longer than others, the application won't even notice. And if the vector is relatively small anyway, even those "slow" pushes are quite fast. Similarly, if your vector grow for a while, and then stops growing, the "steady state" of your application won't see any resizing pauses at all.

Where resizing becomes a problem is in applications that use vectors to keep ever-growing state where tail latency is important. At large scale, it is simply not okay for one push to take 30 milliseconds when most take double-digit nanoseconds. Worse yet, these resize pauses can compound to create significant spikes in tail latency.

This crate implements a technique referred to as "incremental resizing", in contrast to the common "all-at-once" approached outlined above. At its core, the idea is pretty simple: instead of moving all the elements to the resized vector immediately, move a couple each time a push happens. This spreads the cost of moving the elements so that each push becomes a little slower until the resize has finished, instead of one push becoming a lot slower.

This approach isn't free, however. While the resize is going on, the old vector must be kept around (so memory isn't reclaimed immediately), and iterators and other vector-wide operations must access both vectors, which makes them slower. Only once the resize completes is the old vector reclaimed and full performance restored.


```

## What is this implementation struggling with

- We still need to implement in a way VecDeque, so when we shift we dont move all of the elements to the right or make big memory allocations when doing `carry()`, still, Go slices are very powerful and right now there is a workaround in the code. Still we need to do that!

## Why Go implementation?

Because it's fun and I like Jon's work a lot. Also generics are coming to Go 🥳.

## License

Licensed under either of

- Apache License, Version 2.0
  ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
- MIT license
  ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)
