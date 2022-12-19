import http from "k6/http";

export default function () {
  const params = {
    headers: {
      "X-My-Header": "something",
    },
  };

  http.get("http://localhost:3434/latency/50ms", null, params);
}
