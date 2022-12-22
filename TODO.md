## TODO

- [X] Experiment with using require in the context of the `InstrumenHTTP` function
- [X] Set the headers, even if none were passed in
- [X] Handle http request with, and without body parameter
- [X] Add output metadata per request
- [ ] What metadata should the k6 output emit? Same as the HTTP header? 
- [ ] What happens with span_id? do we set one? do we emit it as a non-indexed tag?
- [ ] how to handle batch?
- [ ] how to handle request?
- [ ] Implement the baggage W3C specification?