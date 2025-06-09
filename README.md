# Go Concurrency Compendium

This repository serves as a focused resource for understanding and implementing concurrency in Go. It aims to provide clear explanations, practical examples, and effective strategies for leveraging Go's powerful concurrency primitives.

## What's Included:

* **Concurrency Patterns:** Explore common and effective patterns for structuring concurrent Go applications, such as worker pools, generator, fan-in, pipeline and more.
* **Data Race Handling:** Learn how to prevent and resolve data races, a common pitfall in concurrent programming, using Go's built-in mechanisms:
    * **`sync.Mutex`:** Understand the fundamentals of mutual exclusion and how to protect shared resources with mutexes.
    * **`sync/atomic`:** Discover how to perform atomic operations on primitive types for highly efficient and race-free updates.
* **Practical Concurrency Examples:** Dive into hands-on examples demonstrating how to solve real-world problems using Go's concurrency features, including:
    * The dinning philosophers problem
    * The sleeping barber problem
    * Subscription service
