import http from "k6/http";
import tracing from "k6/x/tracing";

tracing.instrumentHTTP({
  sampling: 12,
  propagator: "w3c",
  baggage: { "X-My-baggage": "something" },
});

// let client = new tracing.Client();

export default () => {
  tracing.get("coucou");
  global.get("https://test-api.k6.io");
  console.log("http: ", JSON.stringify(http));

  //   console.log(`http from the script: ${JSON.stringify(http)}`);

  //   console.log(intrumentedHTTPGet("https://test-api.k6.io", {}));
  //   console.log(http);
  //   console.log(Object.keys(http));
  //   http.get = function (url, params) {
  //     console.log("GETTING");
  //   };
  //   http.get("https://test-api.k6.io", {});
  //   console.log(http.get("https://test-api.k6.io", {}));
  //   let someurl = "https://test.k6.io";
  //   http.get(someurl); // this will be instrumented
};
