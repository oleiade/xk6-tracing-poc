## TODO

- [X] Experiment with using require in the context of the `InstrumenHTTP` function
- [X] Set the headers, even if none were passed in
- [X] Handle http request with, and without body parameter
- [X] Add output metadata per request
- [X] What happens with span_id? do we set one? do we emit it as a non-indexed tag? :: We need to set it, but we don't actually actually use it; we could simply attach a random number every time. TURNS OUT: we already add one in our propagator (those random number generators)
- [ ] Check various propagators and look into how they handle sampling
- [ ] Add checks that the http object exposes the expected methods
- [ ] What metadata should the k6 output emit? Same as the HTTP header? 
- [ ] how to handle batch?
- [ ] how to handle request?
- [ ] Nice to have: Implement the baggage W3C specification?
- [ ] Nice to have: sampling, we have a proposal, but we could wait for a feature request 
- [ ] Sample at request / sending a sampling bit / no sampling at all


## Questions and remarks

- [ ] Should we fail the whole HTTP request if something related to tracing is not correct, like: arguments are not correct, or batch arg is not an array etc? We rely on the underlying function working correctly, so we should probably not fail the whole request, but we should probably log an error?