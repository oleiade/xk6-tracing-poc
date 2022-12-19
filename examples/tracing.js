import http from "k6/http";
import { check } from "k6";
import tracing from "k6/x/tracing";

const options = {
  vus: 1,
  iterations: 1,
};

tracing.instrumentHTTP({
  sampling: 12,
  propagator: "w3c",
  baggage: { "X-My-baggage": "some other thing" },
});

export default () => {
  const params = {
    headers: {
      "X-My-Header": "something",
    },
  };

  let res = http.get("http://localhost:3434/latency/50ms", params);
  check(res, {
    "status is 200": (r) => r.status === 200,
  });

  let data = { name: "Bert" };

  res = http.post("http://httpbin.org/post", JSON.stringify(data), params);
  check(res, {
    "status is 200": (r) => r.status === 200,
  });
};
