import http from "k6/http";
import { check, sleep, group } from "k6";
import { generateGetAd, generatePostAd } from "./generateData.js";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 5000,
      timeUnit: "1s",
      duration: "20s",
      preAllocatedVUs: 100,
      maxVUs: 5000,
    },
  },
};

// export let options = {
//   stages: [{ duration: "10s", target: 100 }],
//   thresholds: {
//     http_req_duration: ["p(95)<500"],
//     http_req_failed: ["rate<0.01"],
//   },
// };

export default function () {
  group("GET Ad", function () {
    //let url = "http://127.0.0.1:8080/api/v1/ad/test";
    const { url } = generateGetAd();
    const res = http.get(url);
    check(res, {
      "GET status was 200": (r) => r.status === 200,
    });
  });
  sleep(1);
}
