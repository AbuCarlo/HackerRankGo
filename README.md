# HackerRankGo

* Author: Anthony Nassar, anthonyabunassar@gmail.com

This repository contains Go implementations of "interview questions" from https://www.hackerrank.com/interview/interview-preparation-kit. Since people are always posting their complete solutions to the "discussions," and these problems
are merely for practice anyway, I saw no reason to keep this repository private.

You can see that most of this code is under test. I used Go's native support for testing, including fuzzing and benchmarking. I also used property-based testing where appropriate. 

To run all tests:

    go test ./...

...or, for verbose output:

    go test ./... -v

It's more likely that you'll run particular tests or benchmarks directly from IDE.